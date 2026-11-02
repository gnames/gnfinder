# Install just: https://github.com/casey/just
APP := "gnfinder"
ORG := "github.com/gnames/"

# Project variables
VERSION := `git describe --tags 2>/dev/null || echo "dev"`
VER := `git describe --tags --abbrev=0 2>/dev/null || echo "dev"`
DATE := `date -u '+%Y-%m-%d_%H:%M:%S'`

NO_C := "CGO_ENABLED=0"
FLAGS_LD := "-trimpath -ldflags '-w -s -X " + ORG + APP + \
            "/pkg.Build=" + DATE + " -X " + ORG + APP + \
            "/pkg.Version=" + VERSION + "'"
FLAGS_REL := "-trimpath -ldflags '-s -w -X " + ORG + APP + \
             "/pkg.Build=" + DATE + "'"

X86 := "GOARCH=amd64"
ARM := "GOARCH=arm64"
LINUX := "GOOS=linux"
MAC := "GOOS=darwin"
WIN := "GOOS=windows"

# Default recipe - install to ~/go/bin
default: install

# Run all tests
test:
    go test -shuffle=on -count=1 -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run tests with coverage report
test-coverage:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

# Download dependencies
deps:
    go mod download

# Install development tools listed in tools.go
tools: deps
    @grep _ tools.go | awk -F'"' '{print $2}' | xargs -tI % go install %

# Build the binary (development build with timestamp and git version)
build:
    @mkdir -p bin
    {{NO_C}} go build {{FLAGS_LD}} -o bin/{{APP}}
    @echo "✅ {{APP}} built to bin/{{APP}}"

# Build release binary (uses version.go for Version, timestamp for Build)
build-release:
    @mkdir -p bin
    {{NO_C}} go build {{FLAGS_REL}} -o bin/{{APP}}
    @echo "✅ {{APP}} release binary built to bin/{{APP}}"

# Install to ~/go/bin (development build with timestamp and git version)
install:
    {{NO_C}} go install {{FLAGS_LD}}
    @echo "✅ {{APP}} installed to ~/go/bin/{{APP}}"

# Build releases for Linux, Mac (Intel, ARM) and Windows
release: dockerhub
    @mkdir -p bin/releases

    {{NO_C}} {{LINUX}} {{X86}} go build {{FLAGS_REL}} -o bin/releases/{{APP}}
    tar zcvf bin/releases/{{APP}}-{{VER}}-linux.tar.gz -C bin/releases {{APP}}
    rm bin/releases/{{APP}}

    {{NO_C}} {{MAC}} {{X86}} go build {{FLAGS_REL}} -o bin/releases/{{APP}}
    tar zcvf bin/releases/{{APP}}-{{VER}}-mac-amd64.tar.gz -C bin/releases {{APP}}
    rm bin/releases/{{APP}}

    {{NO_C}} {{MAC}} {{ARM}} go build {{FLAGS_REL}} -o bin/releases/{{APP}}
    tar zcvf bin/releases/{{APP}}-{{VER}}-mac-arm64.tar.gz -C bin/releases {{APP}}
    rm bin/releases/{{APP}}

    {{NO_C}} {{WIN}} {{X86}} go build {{FLAGS_REL}} -o bin/releases/{{APP}}.exe
    cd bin/releases && zip -9 {{APP}}-{{VER}}-win-64.zip {{APP}}.exe
    rm bin/releases/{{APP}}.exe

    @echo "✅ Release binaries created in bin/releases/"

# Build docker image (Dockerfile expects the binary in the project root)
docker:
    {{NO_C}} {{LINUX}} {{X86}} go build {{FLAGS_REL}} -o {{APP}}
    docker buildx build -t gnames/{{APP}}:latest -t gnames/{{APP}}:{{VERSION}} .
    rm {{APP}}

# Build and push docker image to Docker Hub
dockerhub: docker
    docker push gnames/{{APP}}
    docker push gnames/{{APP}}:{{VERSION}}

# Clean build artifacts
clean:
    rm -rf bin coverage.out coverage.html
    go clean
    @echo "✅ Cleaned build artifacts"

# Format all Go code
fmt:
    go fmt ./...

# Run linter (requires golangci-lint)
lint:
    golangci-lint run

# Tidy dependencies
tidy:
    go mod tidy

# Verify project builds and all tests pass
verify: fmt tidy test build
    @echo "✅ Verification complete: code formatted, dependencies tidied, tests passing, build successful"
