now := `date +'%Y-%m-%dT%T'`
sha := `git rev-parse --short HEAD`
base := "github.com/svenax/domusic/cmd"

default: build test

# Download and tidy Go modules
deps:
    go mod download
    go mod tidy

# Build the project with git version info
build:
    go build -ldflags "-X {{base}}.gitSha1={{sha}} -X {{base}}.buildTime={{now}}"

# Install the project with git version info
install:
    go clean
    go install -ldflags "-X {{base}}.gitSha1={{sha}} -X {{base}}.buildTime={{now}}"

# Test the build by running version command
test: build
    go test ./...
    ./domusic version