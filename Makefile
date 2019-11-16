NAME = $(notdir $(PWD))

VERSION = $(shell printf "%s.%s.%s" \
	$$(git describe --tags) \
	$$(git rev-list --count HEAD) \
	$$(git rev-parse --short HEAD) \
)

GOFLAGS = GO111MODULE=on CGO_ENABLED=0

version:
	@echo $(VERSION)

test:
	$(GOFLAGS) go test -failfast -v ./...

get:
	$(GOFLAGS) go get -v -d ./cmd/...

build:
	$(GOFLAGS) go build \
		 -ldflags="-s -w -X main.version=$(VERSION)" \
		 -gcflags="-trimpath=$(GOPATH)" \
		 ./cmd/...
