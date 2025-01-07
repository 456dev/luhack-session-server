# Variables
BINARY_NAME=session-server
BINARY_PATH=./$(BINARY_NAME)
TAR_NAME=release.tar.gz
STATIC_DIR=./static
TEMPLATE_DIR=./templates

# Default target
all: build bundle

# Build the static binary
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_PATH) .

# Bundle the static binary and directories into a tar.gz
bundle:
	tar -czvf $(TAR_NAME) $(BINARY_PATH) $(STATIC_DIR) $(TEMPLATE_DIR) README.md LICENSE

# Clean up
clean:
	rm -f $(BINARY_PATH) $(TAR_NAME)

.PHONY: all build bundle clean
