all:
	go build -o bin/speedtest cmd/speedtest/main.go
clean:
	rm -r bin
test:
	go test ./...
format:
	gofmt -w ./..
