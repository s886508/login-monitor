package consumer

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/s886508/ruckus-assignment/pkg/metric"
	"github.com/s886508/ruckus-assignment/pkg/model"
)

type LoginEventConsumer struct {
	Buffer                    chan model.LoginEvent
	LoginFailEvents           map[string][]model.LoginEvent
	loginFailFirstTimestamp   map[string]time.Time
	loginSuccessLastTimestamp map[string]time.Time
	loginFailCount            map[string]int
	alertBuffer               chan model.Alert
	ctxCancel                 context.Context
}

func (c *LoginEventConsumer) Init(ctx context.Context, alertBuffer chan model.Alert) {
	c.Buffer = make(chan model.LoginEvent, 100) // TODO: Make it configurable
	c.LoginFailEvents = make(map[string][]model.LoginEvent)
	c.loginFailFirstTimestamp = make(map[string]time.Time)
	c.loginSuccessLastTimestamp = make(map[string]time.Time)
	c.loginFailCount = make(map[string]int)
	c.alertBuffer = alertBuffer
	c.ctxCancel = ctx
}

func (c *LoginEventConsumer) Run() {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			start := time.Now()
			select {
			case e, ok := <-c.Buffer:
				if !ok {
					return
				}

				// Record for every event
				metric.TotalEventProcessed++

				// Check if the event is valid for processing
				if !e.IsValid() {
					log.Println("Invalid event")
					metric.TotalInvalidEvents++
					continue
				}

				log.Printf("Record event: User: %s, Timestmap: %v, Success: %v\n", e.UserID, e.Timestamp, e.Success)

				// ignore the successfully login if the last one is later than cureent one
				if c.loginSuccessLastTimestamp[e.UserID].After(e.Timestamp) {
					continue
				}

				firstFailTS, failSessionExist := c.loginFailFirstTimestamp[e.UserID]
				if e.Success {
					// ignore the successfully log in if it was from earlier times
					if failSessionExist && firstFailTS.After(e.Timestamp) {
						continue
					}
					// reset failed login information once login successfully
					c.loginFailCount[e.UserID] = 0
					c.loginSuccessLastTimestamp[e.UserID] = e.Timestamp
					delete(c.loginFailFirstTimestamp, e.UserID)
					delete(c.LoginFailEvents, e.UserID)
					continue
				}

				metric.TotalFailedLoginEvents++

				// first seen the UserID within failed log in
				if !failSessionExist {
					c.loginFailFirstTimestamp[e.UserID] = e.Timestamp
					c.loginFailCount[e.UserID] = 1
					c.LoginFailEvents[e.UserID] = []model.LoginEvent{e}
					continue
				}

				// failed to log in, record fail events
				c.LoginFailEvents[e.UserID] = append(c.LoginFailEvents[e.UserID], e)

				// The log in attempt is over 30 seconds, re-counting
				if e.Timestamp.After(firstFailTS) {
					if e.Timestamp.Sub(firstFailTS) > 30*time.Second {
						log.Printf("Fail to login > 30 secs: %v\n", firstFailTS)
						// Record the latest failed log in
						c.loginFailFirstTimestamp[e.UserID] = e.Timestamp
						c.loginFailCount[e.UserID] = 1
					} else {
						c.loginFailCount[e.UserID]++
					}
				} else if e.Timestamp.Before(firstFailTS) && e.Timestamp.Add(30*time.Second).After(firstFailTS) {
					c.loginFailFirstTimestamp[e.UserID] = e.Timestamp
					c.loginFailCount[e.UserID]++
				}

				// send alert while fail to log in 3 times within 30 seconds
				if c.loginFailCount[e.UserID] == 3 {
					sort.Slice(c.LoginFailEvents[e.UserID], func(i, j int) bool {
						return c.LoginFailEvents[e.UserID][i].Timestamp.Before(c.LoginFailEvents[e.UserID][j].Timestamp)
					})

					failCount := len(c.LoginFailEvents[e.UserID])

					// calculate time window for the last 3 consecutive fail log in
					start := c.LoginFailEvents[e.UserID][failCount-3].Timestamp
					end := c.LoginFailEvents[e.UserID][failCount-1].Timestamp
					duration := end.Sub(start)

					alert := model.Alert{
						UserID:      e.UserID,
						FailedCount: failCount,
						TimeWindow:  fmt.Sprintf("%.3f seconds", duration.Seconds()),
						Events:      c.LoginFailEvents[e.UserID],
					}
					c.alertBuffer <- alert
					metric.TotalAlertSent++
					//log.Printf("Alert sent: %v\n", alert)
				}
			case <-c.ctxCancel.Done():
				log.Println("Consumer close")
				return
			}
			duration := time.Since(start)
			metric.TotalProcessingDuration += duration.Nanoseconds()
		}
	}()

	wg.Wait()
}
