package alerting

import (
	"context"
	"sync"
	"testing"

	"github.com/s886508/ruckus-assignment/pkg/model"
	"github.com/stretchr/testify/assert"
)

func createTestSender() *AlertSender {
	ctx, _ := context.WithCancel(context.Background())
	sender := &AlertSender{}
	sender.Init(ctx)
	return sender

}

func TestAlertSender(t *testing.T) {
	var wg sync.WaitGroup

	// Case 1, success login before send alert
	sender := createTestSender()
	wg.Add(1)
	go func() {
		defer wg.Done()
		sender.Run()
	}()
	sender.Buffer <- model.Alert{UserID: "TestUserA", FailedCount: 3, TimeWindow: "30 seconds", Events: []model.LoginEvent{}}
	close(sender.Buffer)
	wg.Wait()

	assert.Empty(t, sender.Buffer)
	assert.Len(t, sender.outputBuffer, 1)
}
