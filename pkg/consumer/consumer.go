package consumer

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/s886508/ruckus-assignment/pkg/model"
)

type LoginEventConsumer struct {
	Buffer                    chan model.LoginEvent
	LoginEvents               map[string][]model.LoginEvent
	LoginFailEvents           map[string][]model.LoginEvent
	loginFailFirstTimestamp   map[string]time.Time
	loginSuccessLastTimestamp map[string]time.Time
	loginFailCount            int
	alertBuffer               chan model.Alert
	ctxCancel                 context.Context
}

func (c *LoginEventConsumer) Init(ctx context.Context, alertBuffer chan model.Alert) {
	c.Buffer = make(chan model.LoginEvent, 100)         // TODO: Make it configurable
	c.LoginEvents = make(map[string][]model.LoginEvent) // TODO: Just for records at the moment, no real use case now
	c.LoginFailEvents = make(map[string][]model.LoginEvent)
	c.loginFailFirstTimestamp = make(map[string]time.Time)
	c.loginSuccessLastTimestamp = make(map[string]time.Time)
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
			select {
			case e, ok := <-c.Buffer:
				if !ok {
					return
				}
				if e.Timestamp.IsZero() {
					log.Println("Invalid timestamp")
					return
				}

				//log.Printf("%v, %v", e.Timestamp, c.loginFailFirstTimestamp[e.UserID])
				// add login evnets to storage every time
				if _, ok := c.LoginEvents[e.UserID]; ok {
					c.LoginEvents[e.UserID] = append(c.LoginEvents[e.UserID], e)
				} else {
					c.LoginEvents[e.UserID] = []model.LoginEvent{e}
				}

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
					c.loginFailCount = 0
					c.loginSuccessLastTimestamp[e.UserID] = e.Timestamp
					delete(c.loginFailFirstTimestamp, e.UserID)
					delete(c.LoginFailEvents, e.UserID)
					continue
				}

				// skip if the failed events timestamp is from the past time
				/*if e.Timestamp.Before(c.loginFailLastTimestamp[e.UserID]) {
					continue
				}*/

				log.Printf("Record fail event: %v\n", e.Timestamp)

				// first seen the UserID within failed log in
				if !failSessionExist {
					//log.Println("First log in failed")
					c.loginFailFirstTimestamp[e.UserID] = e.Timestamp
					c.loginFailCount = 1
					c.LoginFailEvents[e.UserID] = []model.LoginEvent{e}
					continue
				}

				// failed to log in, record fail events
				c.LoginFailEvents[e.UserID] = append(c.LoginFailEvents[e.UserID], e)

				// The log in attempt is over 30 seconds, re-counting
				if e.Timestamp.After(firstFailTS) {
					if e.Timestamp.Sub(firstFailTS) > 30*time.Second {
						log.Printf("> 30 secs: %v\n", firstFailTS)
						// Record the latest failed log in
						c.loginFailFirstTimestamp[e.UserID] = e.Timestamp
						c.loginFailCount = 1
					} else {
						c.loginFailCount++
					}
				} else if e.Timestamp.Before(firstFailTS) && e.Timestamp.Add(30*time.Second).After(firstFailTS) {
					c.loginFailFirstTimestamp[e.UserID] = e.Timestamp
					c.loginFailCount++
				}

				// send alert while fail to log in 3 times within 30 seconds
				if c.loginFailCount == 3 {

					sort.Slice(c.LoginFailEvents[e.UserID], func(i, j int) bool {
						return c.LoginFailEvents[e.UserID][i].Timestamp.Before(c.LoginFailEvents[e.UserID][j].Timestamp)
					})
					failCount := len(c.LoginFailEvents[e.UserID])

					// calculate time window for the last 3 consecutive fail log in
					start := c.LoginFailEvents[e.UserID][failCount-3].Timestamp
					end := c.LoginFailEvents[e.UserID][failCount-1].Timestamp
					duration := end.Sub(start)

					c.alertBuffer <- model.Alert{
						UserID:      e.UserID,
						FailedCount: failCount,
						TimeWindow:  fmt.Sprintf("%.3f seconds", duration.Seconds()),
						Events:      c.LoginFailEvents[e.UserID],
					}
					log.Println("Alert sent")
				}
			case <-c.ctxCancel.Done():
				log.Println("Consumer close")
				return
			}
		}
	}()

	wg.Wait()
}
