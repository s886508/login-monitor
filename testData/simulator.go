package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"
)

var testUsers = []string{"UserA", "UserB", "UserC", "UserD", "UserE"}

type LoginEvent struct {
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
}

func main() {
	numEvents := flag.Int("nums", 20, "Number of events to generate log in events")
	filePath := flag.String("filePath", "simulateTestFile.txt", "File paht to store the simulation events")
	timeOffset := flag.Int("timeOffset", 10, "Time offset to simulate per event in seconds, it will generate timestamp with -timeOffset < time < timeOffset")
	flag.Parse()

	timeOffsetMaxRange := *timeOffset  // seconds
	timeOffsetMinRange := -*timeOffset // seconds

	rand.Seed(time.Now().UnixNano())
	file, err := os.Create(*filePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}

	// Ensure the file is closed after the function completes, even if errors occur.
	defer file.Close()

	// Write a string to the file.

	// Goroutine to handle the signal
	for i := 0; i < *numEvents; i++ {
		// generate log in events into file
		user := testUsers[rand.Intn(5)]
		d := rand.Intn(timeOffsetMaxRange - timeOffsetMinRange + 1)
		time := time.Now().Add(time.Duration(d) * time.Second)
		success := rand.Intn(2)

		event := &LoginEvent{UserID: user, Timestamp: time, Success: success == 1}
		data, err := json.Marshal(event)
		if err != nil {
			log.Println("Fail to load as JSON")
		}

		_, err = file.WriteString(string(data) + "\n")
		if err != nil {
			log.Println("Error writing to file:", err)
			continue
		}
	}
	log.Printf("Simulate file generated: %v", *filePath)
}
