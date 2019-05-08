GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test

BINARY_NAME=chat
BASEDIR=/Users/adrianbrad/workspace/repos/chat-v2

DATABASE_CONFIG_FILE=test-db-config.yaml
APPLICATION_CONFIG_FILE=application-config.yaml

build:
	cd cmd/chat-database;$(GOBUILD) -v -o $(BASEDIR)/$(BINARY_NAME) -mod vendor
	
run:
	$(GORUN) ./cmd/chat-database -b=$(BASEDIR) -d=$(DATABASE_CONFIG_FILE) -a=$(APPLICATION_CONFIG_FILE)

runtest:
	$(GOTEST) {./test/...,./pkg/...,./cmd/...,./configs/...,./internal/...}
