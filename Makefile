GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run

BINARY_NAME=chat
BASEDIR=/Users/adrianbrad/workspace/repos/chat-v2

build:
	cd cmd/chat-database;$(GOBUILD) -o $(BASEDIR)/$(BINARY_NAME) -v -mod vendor
	
run:
	$(GORUN) ./cmd/chat-database -b=$(BASEDIR) -d=test-db-config.yaml -a=application-config.yaml