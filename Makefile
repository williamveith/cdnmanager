BINARY_NAME := cdnmanager
APP_BUNDLE  := build/bin/$(BINARY_NAME).app/Contents/MacOS/$(BINARY_NAME)

# Wails build flags
# -ldflags="-s -w" removes debugging information, making the binary smaller
# -trimpath removes file system paths
WAILS_BUILD_FLAGS := -ldflags "-s -w" -trimpath

# Default target: build the app
all: check build

check:
	@clear
	@echo "Prebuild Check:"
	@if [ -f .env ]; then \
		echo "  .env File Found      | Success"; \
	else \
		if [ -f template.env ]; then \
			echo "Build Failed: Fill out template.env with your Cloudflare credentials and rename template.env to .env"; \
		else \
			echo "Build Failed: Missing .env and template.env. Create .env file or use template from repository"; \
		fi; \
		exit 1; \
	fi

	@if which go > /dev/null; then \
		echo "  Go Installed         | Success (Version $$(go version | awk '{print $$3}'))"; \
	else \
		echo "Build Failed: You are missing Go, which is required to build the app. Follow this link to learn how to install it: https://go.dev/dl/"; \
		exit 1; \
	fi

	@if which npm > /dev/null; then \
		echo "  NPM Installed        | Success (Version $$(npm --version))"; \
	else \
		echo "Build Failed: You are missing NPM, which is required to build the app. Follow this link to learn how to install it: https://nodejs.org/en/download/package-manager"; \
		exit 1; \
	fi

	@if which wails > /dev/null; then \
		WAILS_VERSION=$$(wails version | head -n 1 | awk '{print $$1}'); \
		echo "  Wails Installed      | Success (Version $$WAILS_VERSION)"; \
	else \
		echo "Build Failed: You are missing Wails, which is required to build the app. Follow this link to learn how to install it: https://wails.io/docs/gettingstarted/installation#installing-wails"; \
		exit 1; \
	fi
	@echo "_______________________________________________________"
	@echo "\n"

# Build the Wails application
build: check
	@echo "Building Wails application..."
	@if [ "$(shell uname -s)" = "Darwin" ] && [ "$(shell uname -m)" = "arm64" ]; then \
		echo "Skipping UPX compression for macOS arm64"; \
		wails build -clean -ldflags "-s -w" -trimpath -o $(BINARY_NAME); \
	else \
		wails build -clean -ldflags "-s -w" -trimpath -upx -upxflags "--lzma" -o $(BINARY_NAME); \
	fi
	@echo "_______________________________________________________"
	@echo "\n"
	@echo "Results:"
	@echo "  Build                | Success"
	@echo "  Application          | $(shell pwd)/build/bin/$(BINARY_NAME).app"
	@echo "\n"

# Run the Wails dev server for testing in a live environment
test: check
	@echo "Starting Wails dev server for testing..."
	wails dev

# Run the application after building
run: build
	@echo "Running application..."
	./$(BINARY_NAME)

# PHONY targets are not associated with real files
.PHONY: all check build test run
