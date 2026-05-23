APP := daryaft

.PHONY: help test lint build run clean

help:
	@echo "Daryaft development commands:"
	@echo "  make test   Run Go tests"
	@echo "  make build  Build ./bin/$(APP)"
	@echo "  make run    Run the local CLI"
	@echo "  make clean  Remove local build artifacts"

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
