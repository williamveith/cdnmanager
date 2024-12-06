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
	wails build $(WAILS_BUILD_FLAGS) -o $(BINARY_NAME)
	@echo "Build complete: $(BINARY_NAME)"

	@if command -v strip >/dev/null 2>&1; then \
		echo "Stripping binary..."; \
		strip $(APP_BUNDLE); \
	else \
		echo "No 'strip' tool found. Skipping binary stripping."; \
	fi

# Clean the build artifacts
clean:
	@echo "Cleaning build artifacts..."
	wails build -clean
	@rm -f $(BINARY_NAME)
	@echo "Clean complete."

# Run the Wails dev server for testing in a live environment
test:
	@echo "Starting Wails dev server for testing..."
	wails dev

# Run the application after building
run: build
	@echo "Running application..."
	./$(BINARY_NAME)

# PHONY targets are not associated with real files
.PHONY: all build clean run
