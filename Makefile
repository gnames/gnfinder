GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
GOGET = $(GOCMD) get
FLAG_MODULE = GO111MODULE=on
FLAGS_SHARED = $(FLAG_MODULE) CGO_ENABLED=0 GOARCH=amd64
VERSION=`git describe --tags`
VER=`git describe --tags --abbrev=0`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`

all: install

test: deps install
	$(FLAG_MODULE) go test ./...

deps:
	$(FLAG_MODULE) $(GOGET) github.com/spf13/cobra/cobra@7547e83; \
	$(FLAG_MODULE) $(GOGET) github.com/onsi/ginkgo/ginkgo@505cc35; \
	$(FLAG_MODULE) $(GOGET) github.com/onsi/gomega@ce690c5; \
	$(FLAG_MODULE) $(GOGET) github.com/golang/protobuf/protoc-gen-go@347cf4a; \
	$(GOGENERATE)

version:
	echo "package gnfinder" > version.go
	echo "" >> version.go
	echo "const Version = \"$(VERSION)"\" >> version.go
	echo "const Build = \"$(DATE)\"" >> version.go

asset:
	cd fs; \
	$(FLAGS_SHARED) go run -tags=dev assets_gen.go

build: grpc version asset
	$(GOGENERATE)
	cd gnfinder; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=linux $(GOBUILD);

release: grpc version asset dockerhub
	cd gnfinder; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=linux $(GOBUILD); \
	tar zcvf /tmp/gnfinder-${VER}-linux.tar.gz gnfinder; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=darwin $(GOBUILD); \
	tar zcvf /tmp/gnfinder-${VER}-mac.tar.gz gnfinder; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=windows $(GOBUILD); \
	zip -9 /tmp/gnfinder-${VER}-win-64.zip gnfinder.exe; \
	$(GOCLEAN);

install: grpc version asset
	$(GOGENERATE)
	cd gnfinder; \
	$(FLAGS_SHARED) $(GOINSTALL);

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
