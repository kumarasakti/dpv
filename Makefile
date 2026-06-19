BINARY  := dpv
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)
GOFLAGS := -trimpath

.PHONY: build install test lint clean release

build:
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o bin/$(BINARY) .

install:
	go install $(GOFLAGS) -ldflags '$(LDFLAGS)' .

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ dist/

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

release: clean
	@mkdir -p dist
	@$(foreach platform,$(PLATFORMS),\
		$(eval OS := $(word 1,$(subst /, ,$(platform))))\
		$(eval ARCH := $(word 2,$(subst /, ,$(platform))))\
		$(eval EXT := $(if $(filter windows,$(OS)),.exe,))\
		echo "Building $(OS)/$(ARCH)..." && \
		GOOS=$(OS) GOARCH=$(ARCH) go build $(GOFLAGS) -ldflags '$(LDFLAGS)' \
			-o dist/$(BINARY)-$(OS)-$(ARCH)$(EXT) . && \
	) true
