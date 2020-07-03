VERSION=$(shell git tag --contains HEAD | head)
EXTERNAL_TOOLS = \
    github.com/Songmu/goxz/cmd/goxz \
    github.com/tcnksm/ghr \
    github.com/Songmu/ghch/cmd/ghch

BINDIR:=$(CURDIR)/bin
BINNAME?=jarvis

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
GOTEST=$(GOCMD) test
GOBUILDFLAGS= "-trimpath"

.PHONY: devel-deps
devel-deps:
	@for tool in $(EXTERNAL_TOOLS) ; do \
      echo "Installing $$tool" ; \
      GO111MODULE=off go get $$tool; \
    done

.PHONY: build
build:
	$(GOBUILD) -o $(BINDIR)/$(BINNAME) -v $(GOBUILDFLAGS)

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: generate
generate:
	$(GOGENERATE) -v ./...

.PHONY: clean
clean:
	$(GOCLEAN)

# release part

.PHONY: validate-version
validate-version:
ifeq ($(strip $(VERSION)),)
	$(error Version must be set, please add a tag)
endif

.PHONY: upload
upload: validate-version devel-deps
	ghr -v
	ghr -body="$$(ghch --latest -F markdown)" v${VERSION} pkg/dist/v${VERSION}

.PHONY: validate-version crossbuild
crossbuild:
	goxz -pv=v${VERSION} \
        -arch=386,amd64 -d=./pkg/dist/v${VERSION}
	cd pkg/dist/v${VERSION} && shasum -a 256 * > ./v${VERSION}_SHASUMS

.PHONY: release
release: validate-version devel-deps crossbuild upload
