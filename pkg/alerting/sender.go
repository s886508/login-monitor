package alerting

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/s886508/ruckus-assignment/pkg/model"
)

type Sender struct {
	Buffer    chan model.Alert
	ctxCancel context.Context
}

func (s *Sender) Init(ctx context.Context) {
	s.Buffer = make(chan model.Alert, 100) // TODO: Make it configurable
	s.ctxCancel = ctx
}

func (s *Sender) Run() {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case a, ok := <-s.Buffer:
				if !ok {
					return
				}
				data, err := json.Marshal(a)
				if err != nil {
					log.Println("Fail to load as JSON")
				}
				log.Println("Received alert: " + string(data))
			case <-s.ctxCancel.Done():
				log.Println("Sender close")
				return
			}
		}
	}()
}
