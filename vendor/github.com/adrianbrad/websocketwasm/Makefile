GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
TESTDIRECTORY=./test/test

# Binary names
BINARY_NAME=main.wasm

compiletest:
	GOOS=js GOARCH=wasm $(GOTEST) $(TESTDIRECTORY) -c -o $(TESTDIRECTORY)/$(BINARY_NAME)