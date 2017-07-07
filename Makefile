GIT ?= git
GO ?= go
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --abbrev=0 --tags 2>/dev/null)
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
TARGET := gitlab.com/NagByte/Palette
LD_FLAGS := -X $(TARGET)/common.Version=$(VERSION) -X $(TARGET)/common.BuildTime=$(BUILD_TIME)
FORMAT := '{{ join .Deps " " }}'

.PHONY: help clean dependencies update kill run pull deploy
help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  Palette   to build the main binary for current platform"
	@echo "  clean     to remove generated files"
	@echo "  kill      to stop service"
	@echo "  run       to start service"
	@echo "  deploy    to clone and start"

clean:
	rm -f Palette

dependencies:
	$(GO) get -v ./...

update:
	$(GO) get -t -v -u ./...

Palette: dependencies
	$(GO) build -o="Palette" -ldflags="$(LD_FLAGS)" $(TARGET)

kill:
	pkill Palette || ``

pull: 
	$(GIT) pull origin master

run: Palette 
	./Palette serve 2>> log.txt &

deploy: kill clean pull run
