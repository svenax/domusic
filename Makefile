all: build test

deps:
	go mod download
	go mod tidy

build:
	$(eval now := $(shell date +'%Y-%m-%dT%T'))
	$(eval sha := $(shell git rev-parse --short HEAD))
	$(eval base := github.com/svenax/domusic/cmd)
	go install -ldflags "-X $(base).GitSha1=$(sha) -X $(base).BuildTime=$(now)"

test:
	domusic version