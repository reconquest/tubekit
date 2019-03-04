NAME = $(notdir $(PWD))

VERSION = $(shell printf "%s.%s" \
	$$(git rev-list --count HEAD) \
	$$(git rev-parse --short HEAD) \
)

version:
	@echo $(VERSION)

test:
	@go test -failfast -v ./...

build:
	GO111MODULES=on CGO_ENABLED=0 go build \
		 -ldflags="-s -w -X main.version=$(VERSION)" \
		 -gcflags="-trimpath=$(GOPATH)" \
		 ./cmd/...
