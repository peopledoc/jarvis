BINDIR:=$(CURDIR)/bin
BINNAME?=jarvis

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
GOTEST=$(GOCMD) test

#------
#	all
.PHONY: all

all: test build

#------
#	build
.PHONY: build

build:
	$(GOBUILD) -o $(BINDIR)/$(BINNAME) -v 

#------
#	test
.PHONY: test

test:
	$(GOTEST) -v ./...

#------
# generate codes (ex:gomock)
.PHONY: generate

generate:
	$(GOGENERATE) -v ./...

#------
# clean	
.PHONY: clean

clean:
	$(GOCLEAN)
