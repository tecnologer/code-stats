.SILENT:
.PHONY: *

VERSION=$(shell git describe --tags --always)

install: build
	go install
	@printf '\033[32m\033[1m%s\033[0m\n' ":: Install complete"

build:
	go build
	@printf '\033[32m\033[1m%s\033[0m\n' ":: Build complete"
