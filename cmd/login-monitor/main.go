package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/s886508/ruckus-assignment/pkg/alerting"
	"github.com/s886508/ruckus-assignment/pkg/consumer"
	"github.com/s886508/ruckus-assignment/pkg/input"
	"github.com/s886508/ruckus-assignment/pkg/metric"
	"github.com/s886508/ruckus-assignment/pkg/model"
)

func main() {
	filePath := flag.String("filePath", "", "Input file name")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(*filePath) == 0 {
		log.Fatal("Empty file name")
		os.Exit(1)
	}
	_, err := os.Stat(*filePath)
	if os.IsNotExist(err) {
		log.Fatal("File does not exist")
		os.Exit(1)
	}

	sender := &alerting.AlertSender{}
	sender.Init(ctx)
	consumer := &consumer.LoginEventConsumer{}
	consumer.Init(ctx, sender.Buffer)
	fileInput := &input.FileInput{FilePath: *filePath}
	err = fileInput.Init()
	if err != nil {
		log.Fatal("File does not exist")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	// Start the alert sender
	go func() {
		defer wg.Done()
		sender.Run()
	}()

	// Start the event consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run()
	}()

	// Start the data input, could be file or other message queues
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			payload, err := fileInput.FetchEvent()
			if err != nil {
				break
			}
			var event model.LoginEvent
			err = json.Unmarshal([]byte(payload), &event)
			if err != nil {
				log.Println("Invalid payload format")
				continue
			}
			//log.Printf("%v\n", event)
			consumer.Buffer <- event
		}
	}()

	// Create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)

	// Notify the channel on SIGTERM or SIGINT
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	// Goroutine to handle the signal
	go func() {
		for {
			select {
			case alert, ok := <-sender.OutputBuffer:
				if !ok {
					cancel()
					return
				}
				log.Println("Received alert: " + string(alert))
			case sig := <-signalChan:
				log.Printf("Received signal: %s\n", sig)
				cancel()
				return
			}
		}
	}()

	wg.Wait()
	avgProcessingDuration := metric.TotalProcessingDuration / float64(metric.TotalEventProcessed)
	log.Printf("[Metrics] \n  TotalEventProcessed: %d\n  TotalInvalidEvents %d\n  TotalFailedLoginEvents: %d\n  TotalAlertSent: %d\n  AvgEventProcessingDuration: %.3f\n",
		metric.TotalEventProcessed,
		metric.TotalInvalidEvents,
		metric.TotalFailedLoginEvents,
		metric.TotalAlertSent,
		avgProcessingDuration)

	log.Println("Main exit")
}
