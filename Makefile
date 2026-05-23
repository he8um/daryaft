APP := daryaft

.PHONY: test lint build run clean

test:
	go test ./...

lint:
	go vet ./...

build:
	go build -o bin/$(APP) .

run:
	go run .

clean:
	rm -rf bin dist build coverage.out coverage.html
