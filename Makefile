# Build configuration
REPO        := github.com/guessi/cloudtrail-cli
BINARY      := cloudtrail-cli

# Variables
PKGS        := $(shell test -f go.mod && go list ./... 2>/dev/null || echo "")
BUILDTIME   := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GITVERSION  := $(shell git describe --tags --abbrev=8 2>/dev/null || echo "dev-$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)")
GOVERSION   := $(shell go version | cut -d' ' -f3)
NPROC       := $(shell nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 2)
SHASUM      := $(shell command -v sha1sum >/dev/null 2>&1 && echo "sha1sum" || echo "shasum -a 1")

# Build flags
LDFLAGS     := -s -w \
               -X "$(REPO)/pkg/constants.GitVersion=$(GITVERSION)" \
               -X "$(REPO)/pkg/constants.GoVersion=$(GOVERSION)" \
               -X "$(REPO)/pkg/constants.BuildTime=$(BUILDTIME)"

RELEASE_DIR := releases/$(GITVERSION)

.PHONY: default \
         check-tools \
         staticcheck \
         test \
         dependency \
         build-linux-x86_64 build-linux-arm64 build-darwin-x86_64 build-darwin-arm64 build-windows-x86_64 \
         build-linux build-darwin build-windows \
         build-parallel \
         clean release all

default: build-parallel

check-tools:
	@for tool in go git tar curl; do \
		command -v $$tool >/dev/null 2>&1 || { echo "$$tool is required but not installed"; exit 1; }; \
	done

staticcheck: check-tools
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@if [ -n "$(PKGS)" ]; then staticcheck $(PKGS); else echo "No packages found, skipping staticcheck"; fi

test: check-tools
	@go fmt ./...
	@go vet ./...
	@go test ./...

dependency: check-tools
	@test -f go.mod || { echo "go.mod not found"; exit 1; }
	@go mod download

clean:
	@rm -rf ./releases/* ghr

define build_binary
	@echo "Building for $(1)/$(2)..."
	@mkdir -p $(RELEASE_DIR)/$(3)
	@CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build \
		-ldflags="$(LDFLAGS)" \
		-o $(RELEASE_DIR)/$(3)/$(BINARY)$(4) || exit 1
	@test -f LICENSE || { echo "LICENSE file not found"; exit 1; }
	@cp LICENSE $(RELEASE_DIR)/$(3)/LICENSE$(5)
	@tar zcf $(RELEASE_DIR)/$(BINARY)-$(3).tar.gz \
		-C $(RELEASE_DIR)/$(3) $(BINARY)$(4) LICENSE$(5)
	@test -f $(RELEASE_DIR)/$(BINARY)-$(3).tar.gz || { echo "Failed to create archive"; exit 1; }
endef

build-linux-x86_64: check-tools
	$(call build_binary,linux,amd64,Linux-x86_64,,)

build-linux-arm64: check-tools
	$(call build_binary,linux,arm64,Linux-arm64,,)

build-darwin-x86_64: check-tools
	$(call build_binary,darwin,amd64,Darwin-x86_64,,)

build-darwin-arm64: check-tools
	$(call build_binary,darwin,arm64,Darwin-arm64,,)

build-windows-x86_64: check-tools
	$(call build_binary,windows,amd64,Windows-x86_64,.exe,.txt)

build-linux: check-tools dependency
	@mkdir -p $(RELEASE_DIR)
	@$(MAKE) -j$(NPROC) build-linux-x86_64 build-linux-arm64

build-darwin: check-tools dependency
	@mkdir -p $(RELEASE_DIR)
	@$(MAKE) -j$(NPROC) build-darwin-x86_64 build-darwin-arm64

build-windows: check-tools dependency
	@mkdir -p $(RELEASE_DIR)
	@$(MAKE) -j$(NPROC) build-windows-x86_64

build-parallel: check-tools dependency
	@mkdir -p $(RELEASE_DIR)
	@$(MAKE) -j$(NPROC) build-linux-x86_64 build-linux-arm64 build-darwin-x86_64 build-darwin-arm64 build-windows-x86_64

release: check-tools build-parallel
	@echo "Creating release..."
	@test -n "$(GITHUB_TOKEN)" || { echo "GITHUB_TOKEN is required"; exit 1; }
	@test -d "$(RELEASE_DIR)" || { echo "Release directory not found"; exit 1; }
	@curl -fsSL "https://github.com/tcnksm/ghr/releases/download/v0.17.0/ghr_v0.17.0_linux_amd64.tar.gz" -O || { echo "Failed to download ghr"; exit 1; }
	@tar --strip-components=1 -xvf "ghr_v0.17.0_linux_amd64.tar.gz" "ghr_v0.17.0_linux_amd64/ghr" || { echo "Failed to extract ghr"; exit 1; }
	@rm -f "ghr_v0.17.0_linux_amd64.tar.gz"
	@chmod +x ./ghr
	@./ghr -replace -recreate -token $(GITHUB_TOKEN) $(GITVERSION) $(RELEASE_DIR)/ || { echo "Failed to create release"; exit 1; }
	@test -n "$$(ls $(RELEASE_DIR)/*.tar.gz 2>/dev/null)" || { echo "No tar.gz files found for checksum"; exit 1; }
	@$(SHASUM) $(RELEASE_DIR)/*.tar.gz > $(RELEASE_DIR)/SHA1SUM || { echo "Failed to generate checksums"; exit 1; }

release-only: check-tools
	@echo "Creating release..."
	@test -n "$(GITHUB_TOKEN)" || { echo "GITHUB_TOKEN is required"; exit 1; }
	@test -d "$(RELEASE_DIR)" || { echo "Release directory not found"; exit 1; }
	@curl -fsSL "https://github.com/tcnksm/ghr/releases/download/v0.17.0/ghr_v0.17.0_linux_amd64.tar.gz" -O || { echo "Failed to download ghr"; exit 1; }
	@tar --strip-components=1 -xvf "ghr_v0.17.0_linux_amd64.tar.gz" "ghr_v0.17.0_linux_amd64/ghr" || { echo "Failed to extract ghr"; exit 1; }
	@rm -f "ghr_v0.17.0_linux_amd64.tar.gz"
	@chmod +x ./ghr
	@./ghr -replace -recreate -token $(GITHUB_TOKEN) $(GITVERSION) $(RELEASE_DIR)/ || { echo "Failed to create release"; exit 1; }
	@test -n "$$(ls $(RELEASE_DIR)/*.tar.gz 2>/dev/null)" || { echo "No tar.gz files found for checksum"; exit 1; }
	@$(SHASUM) $(RELEASE_DIR)/*.tar.gz > $(RELEASE_DIR)/SHA1SUM || { echo "Failed to generate checksums"; exit 1; }
	@rm -f ghr

all: staticcheck dependency clean build
