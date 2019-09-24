GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
GOGET = $(GOCMD) get
FLAG_MODULE = GO111MODULE=on
VERSION=`git describe --tags`
VER=`git describe --tags --abbrev=0`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
LDFLAGS=-ldflags "-X main.buildDate=${DATE} \
                  -X main.buildVersion=${VERSION}"


all: install

test: deps install
	$(FLAG_MODULE) go test ./...

deps:
	$(FLAG_MODULE) $(GOGET) github.com/spf13/cobra/cobra@7547e83; \
	$(FLAG_MODULE) $(GOGET) github.com/onsi/ginkgo/ginkgo@505cc35; \
	$(FLAG_MODULE) $(GOGET) github.com/onsi/gomega@ce690c5; \
	$(FLAG_MODULE) $(GOGET) github.com/golang/protobuf/protoc-gen-go@347cf4a; \
	$(GOGENERATE)

build: grpc
	$(GOGENERATE)
	cd gnfinder; \
	$(GOCLEAN); \
	GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) ${LDFLAGS};

release: grpc dockerhub
	cd gnfinder; \
	$(GOCLEAN); \
	GO111MODULE=on GOOS=linux GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
	tar zcvf /tmp/gnfinder-${VER}-linux.tar.gz gnfinder; \
	$(GOCLEAN); \
	GO111MODULE=on GOOS=darwin GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
	tar zcvf /tmp/gnfinder-${VER}-mac.tar.gz gnfinder; \
	$(GOCLEAN); \
	GO111MODULE=on GOOS=windows GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
	zip -9 /tmp/gnfinder-${VER}-win-64.zip gnfinder.exe; \
	$(GOCLEAN);

install: grpc
	$(GOGENERATE)
	cd gnfinder; \
	GO111MODULE=on $(GOINSTALL) ${LDFLAGS};

.PHONY:grpc
grpc:
	cd protob; \
	protoc -I . ./protob.proto --go_out=plugins=grpc:.;

docker: build
	docker build -t gnames/gnfinder:latest -t gnames/gnfinder:${VERSION} .; \
	cd gnfinder; \
	$(GOCLEAN);

dockerhub: docker
	docker push gnames/gnfinder; \
	docker push gnames/gnfinder:${VERSION}
