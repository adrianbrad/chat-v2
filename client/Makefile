GOCMD=go
GOBUILD=$(GOCMD) build

# Binary names
BINARY_NAME=chatclientapi.wasm

build:
	GOOS=js GOARCH=wasm $(GOBUILD) -o ./wasm/$(BINARY_NAME) -v