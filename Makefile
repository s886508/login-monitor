.PHONY: build
build:
	go mod tidy
	go build -v -o ./dist/login-monitor ./cmd/login-monitor/main.go


.PHONY: run
run: build
	./dist/login-monitor --filePath "$(FILEPATH)"

.PHONY: test
test:
	go test -v ./...
