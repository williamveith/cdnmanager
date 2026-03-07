BINARY_NAME := cdnmanager
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin

VERSION := 2.0.0
VOL_NAME := CDN Manager
APP_NAME := CDN Manager.app
APP_SOURCE := $(BIN_DIR)/$(BINARY_NAME).app
APP_EXEC := $(APP_SOURCE)/Contents/MacOS/$(BINARY_NAME)
DMG_DIR := $(BUILD_DIR)/dmg
DMG_NAME := CDN-Manager-v$(VERSION).dmg
DMG_PATH := $(BIN_DIR)/$(DMG_NAME)

WAILS_BUILD_FLAGS := -ldflags "-s -w" -trimpath

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

BOLD := \033[1m
RESET := \033[0m
RED := \033[1;31m
GREEN := \033[1;32m
YELLOW := \033[1;33m
HEADER := \033[1;34m

.PHONY: all check check-wails build stage-dmg dmg release test run clean start-section

all: build

release: clean dmg

start-section:
	@printf '\n'
	@printf '%*s\n' "$$(tput cols 2>/dev/null || echo 80)" '' | tr ' ' '─'

check: check-wails

check-wails:
	@$(MAKE) start-section
	@echo "$(HEADER)Prebuild Check: Checking Wails & Dependencies...\n$(RESET)"
	@wails doctor

build: check
	@$(MAKE) start-section
	@echo "$(HEADER)Building Wails application...\n$(RESET)"
	@if [ "$(UNAME_S)" = "Darwin" ] && [ "$(UNAME_M)" = "arm64" ]; then \
		echo "Skipping UPX compression for macOS arm64"; \
		wails build -clean $(WAILS_BUILD_FLAGS) -o $(BINARY_NAME); \
	else \
		wails build -clean $(WAILS_BUILD_FLAGS) -upx -upxflags "--lzma" -o $(BINARY_NAME); \
	fi
	@$(MAKE) start-section
	@echo "Results:"
	@echo "  Build                | Success"
	@echo "  Application          | $(shell pwd)/$(APP_SOURCE)"
	@open "$(shell pwd)/$(BIN_DIR)"

stage-dmg: build
	@$(MAKE) start-section
	@echo "$(HEADER)Staging DMG contents...\n$(RESET)"
	@rm -rf "$(DMG_DIR)"
	@mkdir -p "$(DMG_DIR)"
	@cp -R "$(APP_SOURCE)" "$(DMG_DIR)/$(APP_NAME)"
	@ln -s /Applications "$(DMG_DIR)/Applications"
	@echo "$(GREEN)Staged$(RESET) | $(shell pwd)/$(DMG_DIR)"

dmg: stage-dmg
	@$(MAKE) start-section
	@echo "$(HEADER)Creating DMG...\n$(RESET)"
	@rm -f "$(DMG_PATH)"
	@hdiutil create \
		-volname "$(VOL_NAME)" \
		-srcfolder "$(DMG_DIR)" \
		-ov \
		-format UDZO \
		"$(DMG_PATH)"
	@echo "$(GREEN)DMG Created$(RESET) | $(shell pwd)/$(DMG_PATH)"
	@open "$(shell pwd)/$(BIN_DIR)"

test: check
	@$(MAKE) start-section
	@echo "$(HEADER)Starting Wails dev server for testing...\n$(RESET)"
	@wails dev

run: build
	@$(MAKE) start-section
	@echo "$(HEADER)Running application...\n$(RESET)"
	@"$(APP_EXEC)"

clean:
	@$(MAKE) start-section
	@echo "$(HEADER)Cleaning Build Directory...\n$(RESET)"
	@rm -rf "$(BUILD_DIR)"
	@echo "Build Directory Clean"