.PHONY: all host target clean

# Default target builds for both host and target
all: host target

# Build for the host system
host:
	@echo "Building for host system..."
	go build -o dbc-updater main.go
	@echo "Host build complete: ./dbc-updater"

# Build for the armv7l target with size optimizations and no debug symbols
target:
	@echo "Building for armv7l target..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags='-s -w' -o dbc-updater-armv7l main.go
	@echo "Target build complete: ./dbc-updater-armv7l"

# Clean up build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f dbc-updater dbc-updater-armv7l
	@echo "Clean complete."
