VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')
FLAG_MODULE = GO111MODULE=on
FLAGS_SHARED = $(FLAG_MODULE) CGO_ENABLED=0 GOARCH=amd64
FLAGS_LD=-ldflags "-w -s -X github.com/gnames/gnfinder.Build=${DATE} \
                  -X github.com/gnames/gnfinder.Version=${VERSION}"
NO_C = CGO_ENABLED=0
GOCMD=go
GOINSTALL=$(GOCMD) install $(FLAGS_LD)
GOBUILD=$(GOCMD) build $(FLAGS_LD)
GOCLEAN=$(GOCMD) clean
GOGET = $(GOCMD) get

all: install

test: deps install
	$(FLAG_MODULE) go test ./...

tools: deps
	@echo Installing tools from tools.go
	@cat gnfinder/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %


deps:
	@echo Download go.mod dependencies
	$(GOCMD) mod download; \

build:
	cd gnfinder; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(NO_C) $(GOBUILD);

release: dockerhub
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

install:
	cd gnfinder; \
	$(FLAGS_SHARED) $(GOINSTALL);

docker: build
	docker build -t gnames/gnfinder:latest -t gnames/gnfinder:${VERSION} .; \
	cd gnfinder; \
	$(GOCLEAN);

dockerhub: docker
	docker push gnames/gnfinder; \
	docker push gnames/gnfinder:${VERSION}

clib:
	cd binding; \
	$(GOBUILD) -buildmode=c-shared -o libgnfinder.so;
