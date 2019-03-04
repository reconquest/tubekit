NAME = $(notdir $(PWD))

VERSION = $(shell printf "%s.%s" \
	$$(git rev-list --count HEAD) \
	$$(git rev-parse --short HEAD) \
)

GOFLAGS = GO111MODULE=on CGO_ENABLED=0

version:
	@echo $(VERSION)

test:
	@go test -failfast -v ./...

get:
	$(GOFLAGS) go get -v -d

build:
	$(GOFLAGS) go build \
		 -ldflags="-s -w -X main.version=$(VERSION)" \
		 -gcflags="-trimpath=$(GOPATH)" \
		 ./cmd/...
