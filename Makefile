.PHONY: staticcheck clean build release all test
.DEFAULT_GOAL := build

PKGS       := $(shell go list ./...)
REPO       := github.com/guessi/cloudtrail-cli
BUILDTIME  := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GITVERSION := $(shell git describe --tags --abbrev=8)
GOVERSION  := $(shell go version | cut -d' ' -f3)
LDFLAGS    := -s -w -X "$(REPO)/pkg/constants.GitVersion=$(GITVERSION)" -X "$(REPO)/pkg/constants.GoVersion=$(GOVERSION)" -X "$(REPO)/pkg/constants.BuildTime=$(BUILDTIME)"

staticcheck: STATICCHECK_VERSION := 2025.1.1
staticcheck:
	@echo "Ensuring staticcheck version ($(STATICCHECK_VERSION))..."
	@current_version=$$(staticcheck -version 2>/dev/null | grep -o '[0-9]\{4\}\.[0-9]\+\.[0-9]\+' || echo "none"); \
	if [ "$$current_version" != "$(STATICCHECK_VERSION)" ]; then \
		echo "Installing staticcheck ($(STATICCHECK_VERSION))..."; \
		go install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION) || exit 1; \
	else \
		echo "staticcheck ($(STATICCHECK_VERSION)) already installed"; \
	fi
	@echo "Running staticcheck..."
	@staticcheck ./...

test:
	@go fmt ./...
	@go vet ./...
	@go test ./...

build-linux-x86_64:
	@echo "Creating Build for Linux (x86_64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./releases/$(GITVERSION)/Linux-x86_64/cloudtrail-cli || exit 1
	@cp ./LICENSE ./releases/$(GITVERSION)/Linux-x86_64/LICENSE
	@tar zcf ./releases/$(GITVERSION)/cloudtrail-cli-Linux-x86_64.tar.gz -C releases/$(GITVERSION)/Linux-x86_64 cloudtrail-cli LICENSE

build-linux-arm64:
	@echo "Creating Build for Linux (arm64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o ./releases/$(GITVERSION)/Linux-arm64/cloudtrail-cli || exit 1
	@cp ./LICENSE ./releases/$(GITVERSION)/Linux-arm64/LICENSE
	@tar zcf ./releases/$(GITVERSION)/cloudtrail-cli-Linux-arm64.tar.gz -C releases/$(GITVERSION)/Linux-arm64 cloudtrail-cli LICENSE

build-darwin-x86_64:
	@echo "Creating Build for macOS (x86_64)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./releases/$(GITVERSION)/Darwin-x86_64/cloudtrail-cli || exit 1
	@cp ./LICENSE ./releases/$(GITVERSION)/Darwin-x86_64/LICENSE
	@tar zcf ./releases/$(GITVERSION)/cloudtrail-cli-Darwin-x86_64.tar.gz -C releases/$(GITVERSION)/Darwin-x86_64 cloudtrail-cli LICENSE

build-darwin-arm64:
	@echo "Creating Build for macOS (arm64)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o ./releases/$(GITVERSION)/Darwin-arm64/cloudtrail-cli || exit 1
	@cp ./LICENSE ./releases/$(GITVERSION)/Darwin-arm64/LICENSE
	@tar zcf ./releases/$(GITVERSION)/cloudtrail-cli-Darwin-arm64.tar.gz -C releases/$(GITVERSION)/Darwin-arm64 cloudtrail-cli LICENSE

build-windows-x86_64:
	@echo "Creating Build for Windows (x86_64)..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./releases/$(GITVERSION)/Windows-x86_64/cloudtrail-cli.exe || exit 1
	@cp ./LICENSE ./releases/$(GITVERSION)/Windows-x86_64/LICENSE.txt
	@tar zcf ./releases/$(GITVERSION)/cloudtrail-cli-Windows-x86_64.tar.gz -C releases/$(GITVERSION)/Windows-x86_64 cloudtrail-cli.exe LICENSE.txt

build-linux: build-linux-x86_64 build-linux-arm64
build-darwin: build-darwin-x86_64 build-darwin-arm64
build-windows: build-windows-x86_64

build: build-linux build-darwin build-windows

clean:
	@echo "Cleanup Releases..."
	@rm -rf ./releases/*
	@rm -f ghr
	@rm -f ghr_*.tar.gz

release: GHR_VERSION := v0.17.0
release:
	@echo "Creating Releases..."
	@if [ -z "$$GITHUB_TOKEN" ]; then \
		echo "Error: GITHUB_TOKEN environment variable is required for releases."; \
		echo "To build binaries locally without releasing, use 'make build' instead."; \
		exit 1; \
	fi
	@if [ ! -f ghr ]; then \
		echo "Downloading ghr..."; \
		curl -fsSL https://github.com/tcnksm/ghr/releases/download/$(GHR_VERSION)/ghr_$(GHR_VERSION)_linux_amd64.tar.gz -o ghr_$(GHR_VERSION)_linux_amd64.tar.gz && \
		tar --strip-components=1 -xf ghr_$(GHR_VERSION)_linux_amd64.tar.gz ghr_$(GHR_VERSION)_linux_amd64/ghr; \
	fi
	@./ghr -version
	@./ghr -replace -recreate -token $$GITHUB_TOKEN $(GITVERSION) releases/$(GITVERSION)/
	@if command -v sha256sum >/dev/null 2>&1; then \
		sha256sum releases/$(GITVERSION)/*.tar.gz > releases/$(GITVERSION)/SHA256SUM; \
	elif command -v shasum >/dev/null 2>&1; then \
		shasum -a 256 releases/$(GITVERSION)/*.tar.gz > releases/$(GITVERSION)/SHA256SUM; \
	else \
		echo "Warning: No SHA256 tool found, skipping checksum generation"; \
	fi

all: staticcheck test clean build
