.RECIPEPREFIX +=

GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=ginkgo

VERSION=`git describe --tags`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
LDFLAGS=-ldflags "-X main.buildDate=${DATE} \
                  -X main.buildVersion=${VERSION}"



all: install
test:
  ginkgo
build:
  cd gnfinder; \
  $(GOCLEAN); \
  $(GOBUILD) ${LDFLAGS};
install:
  cd gnfinder; \
  $(GOINSTALL) ${LDFLAGS};

