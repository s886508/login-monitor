package consumer

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/s886508/ruckus-assignment/pkg/model"
	"github.com/stretchr/testify/assert"
)

func createTestConsumer() (*LoginEventConsumer, chan model.Alert) {
	ctx, _ := context.WithCancel(context.Background())
	alertBuffer := make(chan model.Alert, 10)
	consumer := &LoginEventConsumer{}
	consumer.Init(ctx, alertBuffer)
	return consumer, alertBuffer

}

func TestConsumerRunSequentialEvents(t *testing.T) {
	var wg sync.WaitGroup

	// Case 1, success login before send alert
	consumer, alertBuffer := createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: time.Now(), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: time.Now(), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: time.Now(), Success: true}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	_, ok := consumer.loginFailFirstTimestamp["TestUserA"]
	assert.False(t, ok)
	assert.Empty(t, consumer.LoginFailEvents)
	assert.Empty(t, alertBuffer)

	// Case 2, success login without fail
	consumer, alertBuffer = createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: time.Now(), Success: true}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: time.Now(), Success: true}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: time.Now(), Success: true}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	_, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.False(t, ok)
	assert.Empty(t, consumer.LoginFailEvents)
	assert.Empty(t, alertBuffer)

	// Case 3, success login after 2 fail within 30 seconds and 1 fail > 30 seconds
	consumer, alertBuffer = createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	initTime := time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(35 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(40 * time.Second), Success: true}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	_, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.False(t, ok)
	assert.Empty(t, consumer.LoginFailEvents)
	assert.Empty(t, alertBuffer)

	// Case 5, fail login 3 times within 30 seconds
	consumer, alertBuffer = createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(10 * time.Second), Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	ts, ok := consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	assert.Equal(t, initTime, ts)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 3)
	assert.NotEmpty(t, alertBuffer)

	// Case 6, fail login 2 times within 30 seconds and 1 time > 30 seconds
	consumer, alertBuffer = createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(40 * time.Second), Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	ts, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	assert.Equal(t, initTime.Add(40*time.Second), ts)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 3)
	assert.Empty(t, alertBuffer)

	// Case 7. fail login 1 times within 30 seconds, 1 time > 30 seconds and 1 time > 60 seoncds
	consumer, alertBuffer = createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(40 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(80 * time.Second), Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	ts, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	assert.Equal(t, initTime.Add(80*time.Second), ts)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 3)
	assert.Empty(t, alertBuffer)

	// Case 8. fail login 2 times within 30 seconds, 1 time > 30 seconds and 2 times between 30 to 60 seoncds
	consumer, alertBuffer = createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(40 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(45 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(55 * time.Second), Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	ts, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	assert.Equal(t, initTime.Add(40*time.Second), ts)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 4)
	assert.NotEmpty(t, alertBuffer)

	// Case 9. invalid user or timestamp in between will be ignored
	consumer, alertBuffer = createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{Timestamp: initTime.Add(40 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(25 * time.Second), Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	ts, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	assert.Equal(t, initTime, ts)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 2)
	assert.Empty(t, alertBuffer)

	// Case 10. single event
	consumer, alertBuffer = createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	ts, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	assert.Equal(t, initTime, ts)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 1)
	assert.Empty(t, alertBuffer)

	// Case 8. mixed users
	consumer, alertBuffer = createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(45 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(55 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(56 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(57 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserB", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserB", Timestamp: initTime.Add(10 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserB", Timestamp: initTime.Add(15 * time.Second), Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	ts1, ok := consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	ts2, ok := consumer.loginFailFirstTimestamp["TestUserB"]
	assert.True(t, ok)
	assert.Equal(t, initTime.Add(45*time.Second), ts1)
	assert.Equal(t, initTime, ts2)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 5)
	assert.Len(t, consumer.LoginFailEvents["TestUserB"], 3)
	assert.Len(t, alertBuffer, 2)
}

func TestConsumerRunEventsOutOfOrder(t *testing.T) {
	var wg sync.WaitGroup

	// Case 1, success login before send alert
	consumer, alertBuffer := createTestConsumer()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	initTime := time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(-10 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(-20 * time.Second), Success: true}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	_, ok := consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 2)
	assert.Empty(t, alertBuffer)

	// Case 2, success login without fail
	consumer, alertBuffer = createTestConsumer()
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(-10 * time.Second), Success: true}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: true}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(-5 * time.Second), Success: true}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	_, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.False(t, ok)
	assert.Empty(t, consumer.LoginFailEvents)
	assert.Empty(t, alertBuffer)

	// Case 3, success login after 2 fail within 30 seconds and 1 fail > 30 seconds
	consumer, alertBuffer = createTestConsumer()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(40 * time.Second), Success: true}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(35 * time.Second), Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	_, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.False(t, ok)
	assert.Empty(t, consumer.LoginFailEvents)
	assert.Empty(t, alertBuffer)

	// Case 5, fail login 3 times within 30 seconds
	consumer, alertBuffer = createTestConsumer()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(10 * time.Second), Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	ts, ok := consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	assert.Equal(t, initTime, ts)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 3)
	assert.NotEmpty(t, alertBuffer)

	// Case 6, fail login 2 times within 30 seconds and 1 time > 30 seconds
	consumer, alertBuffer = createTestConsumer()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(40 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	ts, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	assert.Equal(t, initTime.Add(40*time.Second), ts)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 3)
	assert.Empty(t, alertBuffer)

	// Case 7. fail login 1 times within 30 seconds, 1 time > 30 seconds and 1 time > 60 seoncds
	consumer, alertBuffer = createTestConsumer()
	initTime = time.Date(2025, time.October, 26, 8, 0, 0, 0, time.UTC)
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(40 * time.Second), Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime, Success: false}
	consumer.Buffer <- model.LoginEvent{UserID: "TestUserA", Timestamp: initTime.Add(80 * time.Second), Success: false}
	close(consumer.Buffer)
	wg.Wait()

	assert.Empty(t, consumer.Buffer)
	ts, ok = consumer.loginFailFirstTimestamp["TestUserA"]
	assert.True(t, ok)
	assert.Equal(t, initTime.Add(80*time.Second), ts)
	assert.Len(t, consumer.LoginFailEvents["TestUserA"], 3)
	assert.Empty(t, alertBuffer)
}
