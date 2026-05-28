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
    go install -ldflags "{{flags}}"

# Test the build by running version command
test: build
    go test ./...
    ./domusic version
