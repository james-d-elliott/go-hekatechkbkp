EXECUTABLE=hekatechkbkp
WINDOWS=$(EXECUTABLE).exe
LINUX=$(EXECUTABLE)
DARWIN=$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags --always)
COMMIT=$(shell git rev-parse HEAD)

LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Executable=$(EXECUTABLE) -X main.Build=$(TRAVIS_BUILD_NUMBER)"
.PHONY: all test clean

all: test build ## Build and run tests

test: ## Run unit tests
	go vet ./src/
	go test -v -race ./src/

build: windows linux ## Build Windows and Linux binaries
	@echo version: $(VERSION)

windows: $(WINDOWS) ## Build Windows binary

linux: $(LINUX) ## Build Linux binary

darwin: $(DARWIN) ## Build Darwin (macOS) binary

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -i -v -o ./out/$(WINDOWS) $(LDFLAGS)  ./src/

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -i -v -o ./out/$(LINUX) $(LDFLAGS)  ./src/

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -i -v -o ./out/$(DARWIN) $(LDFLAGS)  ./src/

clean: ## Remove previous builds
	rm -fr ./out

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'