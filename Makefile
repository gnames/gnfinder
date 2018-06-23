.RECIPEPREFIX +=

GOCMD=go
GOINSTALL=$(GOCMD) install
GOTEST=ginkgo

VERSION=`git describe --tags`
GITHASH=`git rev-parse HEAD | cut -c1-7`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
LDFLAGS=-ldflags "-X main.buildDate=${DATE} \
                  -X main.buildVersion=${VERSION}-${GITHASH}"



all: install
test:
  ginkgo
install:
  cd gnfinder; \
  $(GOINSTALL) ${LDFLAGS};

