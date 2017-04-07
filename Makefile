PROJECT_NAME ?= BlogggerBot
GO_VARS ?= GOOS=linux GOARCH=amd64
COMMIT := $(shell git rev-parse --short HEAD)
VERSION ?= $(shell git describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
LD_FLAGS := -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)

deps:
	go get -v -u ./...

build:
	$(GO_VARS) go build -o $(PROJECT_NAME).o -ldflags="$(LD_FLAGS)" github.com/mehdy/BlogggerBot/cmd/BlogggerBot

test:
	$(GO_VARS) go test -v -cover -race ./...

clean:
	rm -rf $(PROJECT_NAME).o
