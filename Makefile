.PHONY: build run test clean

# Build the application
build:
	go build -o bin/micro-agent cmd/micro-agent/main.go

# Run the application
run:
	go run cmd/micro-agent/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
