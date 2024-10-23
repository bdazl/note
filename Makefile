.PHONY: build-amd64

# On Arch Linux, the package is extra/mingw-w64-gcc
win_cc := /usr/bin/x86_64-w64-mingw32-gcc
module := github.com/bdazl/note
release_version := $(shell git describe --tag)

install-prereq:
	# https://github.com/mattn/go-sqlite3?tab=readme-ov-file#installation
	CGO_ENABLED=1 go install github.com/mattn/go-sqlite3

build-linux:
	mkdir -p build/amd64/linux
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-X '$(module)/cmd.Version=$(release_version)'" -o build/amd64/linux/note

cross-build-windows:
	mkdir -p build/amd64/windows
	CC=$(win_cc) CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-X '$(module)/cmd.Version=$(release_version)'" -o build/amd64/windows/note.exe
