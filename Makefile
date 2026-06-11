APP := daryaft
VERSION ?= 1.13.0-dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo local)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BUILT_BY ?= make
VERSION_PKG := github.com/he8um/daryaft/pkg/version
LD_FLAGS := -X $(VERSION_PKG).Version=$(VERSION) -X $(VERSION_PKG).Commit=$(COMMIT) -X $(VERSION_PKG).Date=$(DATE) -X $(VERSION_PKG).BuiltBy=$(BUILT_BY)
GO_BIN := $(shell go env GOPATH 2>/dev/null)/bin
TOOL_PATH := $(GO_BIN):$(PATH)

.PHONY: help test lint security rc-check rc-info build build-local version ci release-check run clean homebrew-formula-update homebrew-formula-update-dry-run release-preflight release-preflight-allow-skip

help:
	@echo "Daryaft development commands:"
	@echo "  make test         Run Go tests"
	@echo "  make lint         Run golangci-lint"
	@echo "  make security     Run govulncheck and gosec"
	@echo "  make rc-check     Run release-candidate checks"
	@echo "  make rc-info      Print local release-candidate state"
	@echo "  make build        Build ./bin/$(APP)"
	@echo "  make build-local  Build ./bin/$(APP) with local version ldflags"
	@echo "  make version      Run go run . version"
	@echo "  make ci           Run local CI checks"
	@echo "  make release-check  Run local GoReleaser snapshot check without publishing"
	@echo "  make run          Run the local CLI"
	@echo "  make clean        Remove local build artifacts"
	@echo "  make homebrew-formula-update VERSION=X.Y.Z TAP_DIR=/path  Update tap formula (no push)"
	@echo "  make homebrew-formula-update-dry-run VERSION=X.Y.Z TAP_DIR=/path  Preview tap update"
	@echo "  make release-preflight VERSION=X.Y.Z  Run release version guardrail checks"
	@echo "  make release-preflight-allow-skip VERSION=X.Y.Z  Same but allow version skips"

test:
	go test ./...

lint:
	@PATH="$(TOOL_PATH)" command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is required. Install it with: brew install golangci-lint"; exit 1; }
	@PATH="$(TOOL_PATH)" golangci-lint run

security:
	@PATH="$(TOOL_PATH)"; \
	missing=0; \
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
	@PATH="$(TOOL_PATH)" govulncheck ./...
	@PATH="$(TOOL_PATH)" gosec ./...

rc-check:
	go test ./...
	go build ./...
	go test -race ./internal/downloader
	go test -race ./internal/tui
	$(MAKE) lint
	@PATH="$(TOOL_PATH)"; \
	missing=0; \
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
	@PATH="$(TOOL_PATH)" govulncheck ./...
	@PATH="$(TOOL_PATH)" gosec ./...
	@command -v goreleaser >/dev/null 2>&1 || { echo "GoReleaser is required. Install it with: brew install goreleaser"; exit 1; }
	goreleaser check
	git diff --check
	sh -n scripts/manual-qa-server.sh
	bash -n scripts/update-homebrew-formula.sh
	bash -n scripts/release-preflight.sh

rc-info:
	@echo "Git describe:"
	@git describe --tags --always 2>/dev/null || echo "unavailable"
	@echo
	@echo "RC tags:"
	@git tag --list "v*-rc.*" 2>/dev/null || true
	@echo
	@echo "Version:"
	@go run . version
	@echo
	@echo "Next checks:"
	@echo "  make rc-check"
	@echo "  make release-check"

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

homebrew-formula-update:
	@test -n "$(VERSION)" || (echo "VERSION is required. Example: make homebrew-formula-update VERSION=1.2.0 TAP_DIR=/tmp/homebrew-tap" && exit 1)
	@test -n "$(TAP_DIR)" || (echo "TAP_DIR is required. Example: TAP_DIR=/tmp/homebrew-tap" && exit 1)
	./scripts/update-homebrew-formula.sh --version "$(VERSION)" --tap-dir "$(TAP_DIR)"

homebrew-formula-update-dry-run:
	@test -n "$(VERSION)" || (echo "VERSION is required. Example: make homebrew-formula-update-dry-run VERSION=1.2.0 TAP_DIR=/tmp/homebrew-tap" && exit 1)
	@test -n "$(TAP_DIR)" || (echo "TAP_DIR is required. Example: TAP_DIR=/tmp/homebrew-tap" && exit 1)
	./scripts/update-homebrew-formula.sh --version "$(VERSION)" --tap-dir "$(TAP_DIR)" --dry-run

release-preflight:
	@test -n "$(VERSION)" || (echo "VERSION is required, e.g. make release-preflight VERSION=1.5.0" && exit 1)
	./scripts/release-preflight.sh "$(VERSION)"

release-preflight-allow-skip:
	@test -n "$(VERSION)" || (echo "VERSION is required, e.g. make release-preflight-allow-skip VERSION=1.5.0" && exit 1)
	./scripts/release-preflight.sh "$(VERSION)" --allow-skip
