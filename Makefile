BINDIR		:=	$(CURDIR)/bin
BINNAME		?=	jarvis

GOCMD			=	go
GOBUILD		=	$(GOCMD) build
GOCLEAN		=	$(GOCMD) clean
GOTEST		=	$(GOCMD) test

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
#	test
.PHONY: clean

clean:
	$(GOCLEAN)
