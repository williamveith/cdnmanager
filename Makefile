BINARY_NAME := cdnmanager
APP_BUNDLE  := build/bin/$(BINARY_NAME).app/Contents/MacOS/$(BINARY_NAME)

# Wails build flags
# -ldflags="-s -w" removes debugging information, making the binary smaller
# -trimpath removes file system paths
WAILS_BUILD_FLAGS := -ldflags "-s -w" -trimpath

# If your project requires specific Go build tags, you can add them here:
# Example: WAILS_BUILD_FLAGS += -tags "production"
# WAILS_BUILD_FLAGS += -tags "production"

# Default target: build the app
all: build

# Build the Wails application
build:
	@echo "Building Wails application..."
	@if [ "$(shell uname -s)" = "Darwin" ] && [ "$(shell uname -m)" = "arm64" ]; then \
		echo "Skipping UPX compression for macOS arm64"; \
		wails build -clean -ldflags "-s -w" -trimpath -o $(BINARY_NAME); \
	else \
		wails build -clean -ldflags "-s -w" -trimpath -upx -upxflags "--lzma" -o $(BINARY_NAME); \
	fi
	@echo "Build complete: $(BINARY_NAME)"

# Run the Wails dev server for testing in a live environment
test:
	@echo "Starting Wails dev server for testing..."
	wails dev

# Run the application after building
run: build
	@echo "Running application..."
	./$(BINARY_NAME)

# PHONY targets are not associated with real files
.PHONY: all build run
