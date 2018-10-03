.RECIPEPREFIX +=

GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=ginkgo

VERSION=`git describe --tags`
VER=`git describe --tags --abbrev=0`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
LDFLAGS=-ldflags "-X main.buildDate=${DATE} \
                  -X main.buildVersion=${VERSION}"


all: grpc install

test:
  ginkgo

build: grpc
  cd gnfinder; \
  $(GOCLEAN); \
  GOOS=linux GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
  tar zcvf /tmp/gnfinder-${VER}-linux.tar.gz gnfinder; \
  $(GOCLEAN); \
  GOOS=darwin GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
  tar zcvf /tmp/gnfinder-${VER}-mac.tar.gz gnfinder; \
  $(GOCLEAN); \
  GOOS=windows GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
  zip -9 /tmp/gnfinder-${VER}-win-64.zip gnfinder.exe; \
  $(GOCLEAN);

install:
  cd gnfinder; \
  $(GOINSTALL) ${LDFLAGS};

grpc:
  cd protob; \
  protoc -I . ./protob.proto --go_out=plugins=grpc:.;