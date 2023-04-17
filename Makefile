# Build the binaries
build:
	go build -o ./dist/vault2env ./cmd/vault2env.go

# Run the orders command and reboot when changes to code are detected.
run:
	go run ./cmd/vault2env.go

# Ensure code is properly formatted and doesn't have any errors/bad practices.
lint:
	gofmt -l . && \
		golangci-lint run

# Run all tests.
test:
	go test ./...
