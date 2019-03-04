NAME = $(notdir $(PWD))

VERSION = $(shell printf "%s.%s" \
	$$(git rev-list --count HEAD) \
	$$(git rev-parse --short HEAD) \
)

version:
	@echo $(VERSION)

test:
	@go test -failfast -v ./...
