GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
TESTDIRECTORY=./test/test
TESTPORT=:3000
PID=$(shell lsof -t -i :3000 -s TCP:LISTEN)


# Binary names
BINARY_NAME=chatclientapi.wasm
TEST_BINARY_NAME=testchatclientapi.wasm

build:
	GOOS=js GOARCH=wasm $(GOBUILD) -o ./wasm/$(BINARY_NAME) -v

starttest:
	GOOS=js GOARCH=wasm $(GOTEST) $(TESTDIRECTORY) -c -o $(TESTDIRECTORY)/$(BINARY_NAME)
ifneq ($(PID),) 
	kill $(PID)
endif 
	go run ./test/server.go