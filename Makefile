.PHONY : dev build
.DEFAULT_GOAL := dev

BUILD_PATH := ./build
dev:
	@ go run .

build:
	@rm -rf $(BUILD_PATH)/*
	@go build -o "$(BUILD_PATH)/finchd"
	@echo "$(BUILD_PATH)/finchd"