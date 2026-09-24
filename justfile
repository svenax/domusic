now := `date +'%Y-%m-%dT%T'`
sha := `git rev-parse --short HEAD`
base := "github.com/svenax/domusic/cmd"
flags := f"-X {{base}}.gitSha1={{sha}} -X {{base}}.buildTime={{now}}"

default: build test

# Download and tidy Go modules
deps:
    go mod download
    go mod tidy

# Build the project with git version info
build:
    go build -ldflags "{{flags}}"

# Cross-compile the project for Windows with git version info
win-build:
    GOOS=windows GOARCH=amd64 go build -o domusic.exe -ldflags "{{flags}}"

# Install the project with git version info
install:
    go clean
    GOBIN="$(go env GOPATH)/bin" go install -ldflags "{{flags}}"
    @echo "Installed domusic to $(go env GOPATH)/bin/domusic"
    @if [ -n "$(go env GOBIN)" ] && [ "$(go env GOBIN)" != "$(go env GOPATH)/bin" ]; then GOBIN="$(go env GOBIN)" go install -ldflags "{{flags}}" && echo "Also installed domusic to $(go env GOBIN)/domusic"; fi

# Test the build by running version command
test: build
    go test ./...
    ./domusic version
