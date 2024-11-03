MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
RELEASE_VERSION := $(shell git describe --tags --always 2>/dev/null || git rev-parse --short HEAD)

MODULE := github.com/bdazl/note

# On Arch Linux, the package is extra/mingw-w64-gcc
WIN_CC ?= /usr/bin/x86_64-w64-mingw32-gcc
DOCKER ?= podman
DOCKER_IMAGE ?= note:latest
VHS_IMAGE ?= note-vhs:latest

clean:
	rm -rf $(MAKEFILE_DIR)build

install-prereq:
	# https://github.com/mattn/go-sqlite3?tab=readme-ov-file#installation
	CGO_ENABLED=1 go install github.com/mattn/go-sqlite3

build-all: build-linux cross-build-windows

build-linux:
	mkdir -p build/amd64/linux
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X '$(MODULE)/cmd.Version=$(RELEASE_VERSION)'" -o build/amd64/linux/note

build-cross-windows:
	mkdir -p build/amd64/windows
	CC=$(WIN_CC) CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w -X '$(MODULE)/cmd.Version=$(RELEASE_VERSION)'" -o build/amd64/windows/note.exe

build-docker-all: build-docker build-docker-vhs generate-gifs

build-docker:
	$(DOCKER) build -t $(DOCKER_IMAGE) $(MAKEFILE_DIR)

build-vhs: build-docker-vhs generate-gifs

build-docker-vhs:
	$(DOCKER) build -t $(VHS_IMAGE) $(MAKEFILE_DIR)docs

generate-gifs:
	cd $(MAKEFILE_DIR)docs
	DOCKER="$(DOCKER)" VHS_IMAGE="$(VHS_IMAGE)" $(MAKEFILE_DIR)docs/generate $(MAKEFILE_DIR)docs/gifs
