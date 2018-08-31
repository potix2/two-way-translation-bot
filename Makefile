GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=action

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v .

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).zip

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v

dist: build-linux
	zip $(BINARY_NAME).zip $(BINARY_NAME)
