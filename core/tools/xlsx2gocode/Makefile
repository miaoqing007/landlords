# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=xlsx2gocode

all: build
build: 
		$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) -v 
clean: 
		$(GOCLEAN)
		rm -f $(GOBIN)/$(BINARY_NAME)
debug: 
		$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) -v 
