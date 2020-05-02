GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
DEMO_SERVICE=demoservice
DUMMY_IMAGES=dummyimages

all: test build
build: 
	$(GOBUILD) -o ./bin/$(DEMO_SERVICE) -v ./cmd/${DEMO_SERVICE}
	$(GOBUILD) -o ./bin/$(DUMMY_IMAGES) -v ./cmd/${DUMMY_IMAGES}
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
