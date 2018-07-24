all:
	go build -o bin/speedtest main.go
clean:
	rm -r bin
