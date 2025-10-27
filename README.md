# Description
The program is simulating a alert system which receive login events and record the events. An alert will be trigger and being sent if there are 3 consecutive failed log in with 30 seconds.

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
echo "testData/test2.txt"
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
