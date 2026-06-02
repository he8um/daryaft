APP := daryaft
VERSION ?= 0.5.0-dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo local)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BUILT_BY ?= make
VERSION_PKG := github.com/he8um/daryaft/pkg/version
LD_FLAGS := -X $(VERSION_PKG).Version=$(VERSION) -X $(VERSION_PKG).Commit=$(COMMIT) -X $(VERSION_PKG).Date=$(DATE) -X $(VERSION_PKG).BuiltBy=$(BUILT_BY)

.PHONY: help test lint security build build-local version ci release-check run clean

help:
	@echo "Daryaft development commands:"
	@echo "  make test         Run Go tests"
	@echo "  make lint         Run golangci-lint"
	@echo "  make security     Run govulncheck and gosec"
	@echo "  make build        Build ./bin/$(APP)"
	@echo "  make build-local  Build ./bin/$(APP) with local version ldflags"
	@echo "  make version      Run go run . version"
	@echo "  make ci           Run local CI checks"
	@echo "  make release-check  Run local GoReleaser snapshot check without publishing"
	@echo "  make run          Run the local CLI"
	@echo "  make clean        Remove local build artifacts"

test:
	go test ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is required. Install it with: brew install golangci-lint"; exit 1; }
	golangci-lint run

security:
	@missing=0; \
	if ! command -v govulncheck >/dev/null 2>&1; then \
		echo "govulncheck is required. Install it with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
		missing=1; \
	fi; \
	if ! command -v gosec >/dev/null 2>&1; then \
		echo "gosec is required. Install it with: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
		missing=1; \
	fi; \
	if [ "$$missing" -ne 0 ]; then \
		exit 1; \
	fi
	govulncheck ./...
	gosec ./...

build:
	go build -o bin/$(APP) .

build-local:
	go build -ldflags "$(LD_FLAGS)" -o bin/$(APP) .

version:
	go run . version

ci:
	go mod tidy
	git diff --exit-code go.mod go.sum
	go test ./...
	go build ./...
	go test -race ./internal/tui
	git diff --check
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser check; \
	else \
		echo "GoReleaser not found; skipping goreleaser check. Install with: brew install goreleaser"; \
	fi

release-check:
	@command -v goreleaser >/dev/null 2>&1 || { echo "GoReleaser is required. Install it with: brew install goreleaser"; exit 1; }
	goreleaser release --snapshot --clean --skip=publish

run:
	go run .

clean:
	rm -rf bin dist build coverage.out coverage.html
