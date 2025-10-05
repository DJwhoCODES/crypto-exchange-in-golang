# Binary name
BINARY = exchange


# Default target
all: build


# Build the project
build:
	go build -o $(BINARY) ./...


# Run the project
run:
	go run main.go


# Run tests (verbose, no cache)
test:
	go test -v -count=1 ./...


# Clean build files
clean:
	rm -f $(BINARY)


# Format code
fmt:
	go fmt ./...


# Tidy up modules
tidy:
	go mod tidy