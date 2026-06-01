APP := daryaft
VERSION ?= 0.5.0-dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo local)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BUILT_BY ?= make
VERSION_PKG := github.com/he8um/daryaft/pkg/version
LD_FLAGS := -X $(VERSION_PKG).Version=$(VERSION) -X $(VERSION_PKG).Commit=$(COMMIT) -X $(VERSION_PKG).Date=$(DATE) -X $(VERSION_PKG).BuiltBy=$(BUILT_BY)

.PHONY: help test lint build build-local version release-check run clean

help:
	@echo "Daryaft development commands:"
	@echo "  make test         Run Go tests"
	@echo "  make build        Build ./bin/$(APP)"
	@echo "  make build-local  Build ./bin/$(APP) with local version ldflags"
	@echo "  make version      Run go run . version"
	@echo "  make release-check  Run local GoReleaser snapshot check without publishing"
	@echo "  make run          Run the local CLI"
	@echo "  make clean        Remove local build artifacts"

test:
	go test ./...

lint:
	go vet ./...

build:
	go build -o bin/$(APP) .

build-local:
	go build -ldflags "$(LD_FLAGS)" -o bin/$(APP) .

version:
	go run . version

release-check:
	@command -v goreleaser >/dev/null 2>&1 || { echo "GoReleaser is required. Install it with: brew install goreleaser"; exit 1; }
	goreleaser release --snapshot --clean --skip=publish

run:
	go run .

clean:
	rm -rf bin dist build coverage.out coverage.html
