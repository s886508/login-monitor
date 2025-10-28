# Description
The program is simulating an alert system which receive log in events and record the events. An alert will be trigger and being sent if there are 3 consecutive failed log in with 30 seconds.

# How to Build
Run the following command and the binary will be build under `/dist` directory
```shell
$ make build
go mod tidy
go: downloading gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405
go build -v -o ./dist/login-monitor ./cmd/login-monitor/main.go
github.com/s886508/ruckus-assignment/pkg/consumer
github.com/s886508/ruckus-assignment/pkg/alerting
```

# Run the Tool
It is only support reading events in JSON format line by line as following from a file.

**Note** The `timestamp` is using default format `RFC3339` that can be load into time.Time package.
```json
{"user_id":"A","timestamp":"2025-10-26T14:30:00+00:00","success":false}
{"user_id":"A","timestamp":"2025-10-26T14:30:01+00:00","success":false}
{"user_id":"A","timestamp":"2025-10-26T14:30:32+00:00","success":false}
```

## Command to Test with File
```shell
$ make run FILEPATH=testData/test2.txt
go mod tidy
go build -v -o ./dist/login-monitor ./cmd/login-monitor/main.go
testData/test2.txt
./dist/login-monitor --filePath "testData/test2.txt"
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:30:00 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:30:01 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:30:32 +0000 +0000
2025/10/27 21:15:28 > 30 secs: 2025-10-26 14:30:00 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:32:01 +0000 +0000
2025/10/27 21:15:28 > 30 secs: 2025-10-26 14:30:32 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:32:05 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:32:11 +0000 +0000
2025/10/27 21:15:28 Alert sent
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:35:00 +0000 +0000
2025/10/27 21:15:28 > 30 secs: 2025-10-26 14:32:01 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:35:01 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:35:32 +0000 +0000
2025/10/27 21:15:28 > 30 secs: 2025-10-26 14:35:00 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:36:01 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:36:11 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:40:01 +0000 +0000
2025/10/27 21:15:28 > 30 secs: 2025-10-26 14:36:11 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:40:11 +0000 +0000
2025/10/27 21:15:28 Record fail event: 2025-10-26 14:40:05 +0000 +0000
2025/10/27 21:15:28 Alert sent
2025/10/27 21:15:28 Received alert: {"user_id":"A","failed_count":6,"time_window":"10.000 seconds","events":[{"user_id":"A","timestamp":"2025-10-26T14:30:00Z","success":false},{"user_id":"A","timestamp":"2025-10-26T14:30:01Z","success":false},{"user_id":"A","timestamp":"2025-10-26T14:30:32Z","success":false},{"user_id":"A","timestamp":"2025-10-26T14:32:01Z","success":false},{"user_id":"A","timestamp":"2025-10-26T14:32:05Z","success":false},{"user_id":"A","timestamp":"2025-10-26T14:32:11Z","success":false}]}
2025/10/27 21:15:28 Received alert: {"user_id":"A","failed_count":4,"time_window":"10.000 seconds","events":[{"user_id":"A","timestamp":"2025-10-26T14:36:11Z","success":false},{"user_id":"A","timestamp":"2025-10-26T14:40:01Z","success":false},{"user_id":"A","timestamp":"2025-10-26T14:40:05Z","success":false},{"user_id":"A","timestamp":"2025-10-26T14:40:11Z","success":false}]}
```

# Simulation Tool
The tool is able to generate testing input file with given arguments. Please see below for the usage:
```shell
$ go run testData/simulator.go --help
Usage of /home/will/.cache/go-build/fa/fa68702c24857c823d0fe26b345102c559f27fbd5088eef8f96f884b664b7b3d-d/simulator:
  -filePath string
        File paht to store the simulation events (default "simulateTestFile.txt")
  -nums int
        Number of events to generate log in events (default 20)
  -timeOffset int
        Time offset to simulate per event in seconds, it will generate timestamp with -timeOffset < time < timeOffset (default 10)
```

Example to generate 500K data to input:
```shell
$ go run testData/simulator.go --nums 500000 -timeOffset 40
2025/10/28 13:10:18 Simulate file generated: simulateTestFile.txt
```

Test with the testing file:
```shell
$ make run FILEPATH=simulateTestFile.txt
...
^C2025/10/28 13:23:00 Received signal: interrupt
2025/10/28 13:23:00 Sender close
2025/10/28 13:23:00 Consumer close
2025/10/28 13:23:00 [Metrics]
  TotalEventProcessed: 500000
  TotalInvalidEvents 0
  TotalFailedLoginEvents: 3084
  TotalAlertSent: 389
  AvgEventProcessingDuration(nanoSeconds): 80
2025/10/28 13:23:00 Main exit
```

# Run the Unit Test
```shell
$ make test
go test -v ./...
?       github.com/s886508/ruckus-assignment/cmd/login-monitor  [no test files]
=== RUN   TestAlertSender
--- PASS: TestAlertSender (0.00s)
PASS
ok      github.com/s886508/ruckus-assignment/pkg/alerting       (cached)
=== RUN   TestConsumerRunSequentialEvents
...
--- PASS: TestConsumerRunSequentialEvents (0.00s)
=== RUN   TestConsumerRunEventsOutOfOrder
...
--- PASS: TestConsumerRunEventsOutOfOrder (0.00s)
PASS
ok      github.com/s886508/ruckus-assignment/pkg/consumer       (cached)
?       github.com/s886508/ruckus-assignment/pkg/input  [no test files]
?       github.com/s886508/ruckus-assignment/pkg/metric [no test files]
=== RUN   TestLoginEventIsValid
--- PASS: TestLoginEventIsValid (0.00s)
PASS
ok      github.com/s886508/ruckus-assignment/pkg/model  0.008s
?       github.com/s886508/ruckus-assignment/testData   [no test files]
```

# Graceful Shutdown
The tool is listening to few signal and will lead graceful shutdown for each goroutines before exising the tool.
```shell
$ make run FILEPATH=testData/test2.txt
...
^C2025/10/27 21:15:56 Received signal: interrupt
2025/10/27 21:15:56 Consumer close
2025/10/27 21:15:56 Sender close
2025/10/27 21:15:56 Main exit
make: *** [Makefile:10: run] Interrupt
```

# Design
## Events Orderings
This is one of the major design for my tool that the events could be out of orders. Meanwhile the tool checkes the events with ealier timestamp to make corresponding handlings.

## Alert
The alert from the design document does not mention the exact context of `TimeWindow` and `Events`. So I mades assumption that what the user would like to see while receiving the alert.
1. TimeWindow: The duration of 3 consecutive failed login within 30 seconds. The unit is seoncds as well.
2. Events: The failed login events in the past history before a successfully log in. So the `FailedCount` is counted for the length of `Events` as well.

# Future Integration
## Message Queue (Input strategy)
The events are better leverage either one of the message queue system for consitency, fault tolerance and system resume as a consideration. For example, using Kafka as an upstream to pass in the log in events or a downstream as an alert notification system. The message queue is able to connect the current monitoring service with other serivce to expand its capabilities and usages. Also, the message queue usaually do replicas to avoid data loss and help the service to restore from specific time point.

The serivce has defined an interface to different approaches to be implemented. Please check [input.go](https://github.com/s886508/ruckus-assignment/blob/main/pkg/input/input.go) for more details. The interface can also be expanded if needed.

## Storage
Right now, the service is using in-memory to store the events, alerts. This can be achieved in such way if the data volume is extremely small. In the real world, log in events should be a high data volume events and will need other storage, such as cache, persistent storage (database etc...). It can be both to speed up the service while querying log in events. There are couple advantage of the design.
1. Persisten storage: The table schema and index can be levarage here to improve the overall query effieciency compared to in-memory storage if the data volume of log in events is high.
2. Retention: The data can set retention per users or other strategy to keep better query performance.
3. Service Down and Recovery: With persisten storage, the log in events will less likely to be loss compared to in-memory storage and able to keep the alerting service working as expected robustly.

## Scalability
To improve overall system performance, the service shoudl be able to be deployed as a microservice and can be deployed as multiple instance. With the message queue and persisten storage, this can be a stateless service and scale up as needed. In the mean time, take kafka as an example, the partitions and retention settings can be different or configurable to tolerate more high data volume and to be more real-time processing to send alert in time to the end users.

Also, there are few couple go channel buffer that can be setup with a configurable settings to achiever hight throughpyt overall. This can be implemented later on as needed. 
