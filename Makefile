BINARY_NAME := cdnmanager
BUILD_DIR := build
APP_BUNDLE  := build/bin/$(BINARY_NAME).app/Contents/MacOS/$(BINARY_NAME)

# Wails build flags
# -ldflags="-s -w" removes debugging information, making the binary smaller
# -trimpath removes file system paths
WAILS_BUILD_FLAGS := -ldflags "-s -w" -trimpath

# Formatting variables
BOLD := \033[1m
RESET := \033[0m
RED := \033[1;31m
GREEN := \033[1;32m
YELLOW := \033[1;33m
HEADER := \033[1;34m

# Default target: build the app
all: check build

start-section:
	@printf '\n'
	@printf '%*s\n' "$(shell tput cols)" '' | tr ' ' '─'

check:
	@clear
	@echo "$(HEADER)Running Checks$(RESET)"
	@$(MAKE) check-wails
	@$(MAKE) check-env

check-env:
	@$(MAKE) start-section
	@echo "$(HEADER)Prebuild Check: Checking .env File...\n$(RESET)"
	
	@if [ -f .env ]; then \
		echo "$(GREEN)Found$(RESET) | .env"; \
	else \
		if [ -f template.env ]; then \
			cp template.env .env; \
			echo "$(RED)Build Failed:$(RESET) No .env file found. New .env file create. Fill in the new .env file with your Cloudflare credentials$(RESET)"; \
		else \
			echo "$(RED)Build Failed:$(RESET) Missing .env and template.env. Create .env file or use template from repository"; \
		fi; \
		exit 1; \
	fi

	@EMAIL=$$(grep -E "^cloudflare_email[ ]{0,1}=[ ]{0,1}['\"].{1,64}@.{2,255}['\"]{0,1}$$" .env); \
	if [ -z "$$EMAIL" ]; then \
		echo "$(RED)Build Failed:$(RESET) No valid cloudflare_email in .env"; \
		exit 1; \
	else \
		echo "${GREEN}Valid${RESET} | $$EMAIL"; \
	fi

	@API_KEY=$$(grep -E "^cloudflare_api_key[ ]{0,1}=[ ]{0,1}['\"].{1,}['\"]{0,1}$$" .env); \
	if [ -z "$$API_KEY" ]; then \
		echo "${RED}Build Failed:${RESET} Enter valid cloudflare_api_key in .env"; \
		exit 1; \
	else \
		echo "${GREEN}Valid${RESET} | $${API_KEY}"; \
	fi

	@ACCOUNT_ID=$$(grep -E "^account_id[ ]{0,1}=[ ]{0,1}['\"].{1,}['\"]{0,1}$$" .env); \
	if [ -z "$$ACCOUNT_ID" ]; then \
		echo "${RED}Build Failed:${RESET} Enter valid account_id in .env"; \
		exit 1; \
	else \
		echo "${GREEN}Valid${RESET} | $${ACCOUNT_ID}"; \
	fi

	@NAMESPACE_ID=$$(grep -E "^namespace_id[ ]{0,1}=[ ]{0,1}['\"].{1,}['\"]{0,1}$$" .env); \
	if [ -z "$$NAMESPACE_ID" ]; then \
		echo "${RED}Build Failed:${RESET} Enter valid namespace_id in .env"; \
		exit 1; \
	else \
		echo "${GREEN}Valid${RESET} | $${NAMESPACE_ID}"; \
	fi

	@DOMAIN=$$(grep -E "^domain[ ]{0,1}=[ ]{0,1}['\"].{1,}\.{1}[a-zA-Z]{2,63}['\"]{0,1}$$" .env); \
	if [ -z "$$DOMAIN" ]; then \
		echo "${RED}Build Failed:${RESET} Enter valid domain in .env"; \
		exit 1; \
	else \
		echo "${GREEN}Valid${RESET} | $${DOMAIN}"; \
	fi
	
check-wails:
	@$(MAKE) start-section
	@echo "$(HEADER)Prebuild Check: Checking Wails & Dependencies...\n$(RESET)"
	@wails doctor

# Build the Wails application
build: check
	@$(MAKE) start-section
	@echo "$(HEADER)Building Wails application...\n$(RESET)"
	@if [ "$(shell uname -s)" = "Darwin" ] && [ "$(shell uname -m)" = "arm64" ]; then \
		echo "Skipping UPX compression for macOS arm64"; \
		wails build -clean -ldflags "-s -w" -trimpath -o $(BINARY_NAME); \
	else \
		wails build -clean -ldflags "-s -w" -trimpath -upx -upxflags "--lzma" -o $(BINARY_NAME); \
	fi

	@$(MAKE) start-section
	@echo "Results:"
	@echo "  Build                | Success"
	@echo "  Application          | $(shell pwd)/build/bin/$(BINARY_NAME).app"
	@open $(shell pwd)/build/bin

# Run the Wails dev server for testing in a live environment
test: check
	@$(MAKE) start-section
	@echo "$(HEADER)Starting Wails dev server for testing...\n$(RESET)"
	@wails dev

# Run the application after building
run: build
	@$(MAKE) start-section
	@echo "$(HEADER)Running application...\n$(RESET)"
	./$(BINARY_NAME)

# Deletes all previous builds
clean:
	@clear
	@$(MAKE) start-section
	@echo "${HEADER}Cleaning Build Directory...\n$(RESET)"
	@find build -mindepth 1 -type d -exec rm -rf {} +
	@echo "Build Directory Clean"

# PHONY targets are not associated with real files
.PHONY: all start-section check-wails check-env build test run clean
