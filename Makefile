PROJ_NAME = gnfinder

VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')

NO_C = CGO_ENABLED=0
FLAGS_SHARED = $(NO_C) GOARCH=amd64
FLAGS_MAC_ARM = $(NO_C) $GOARCH=arm64 GOOS=darwin
FLAGS_LD = -trimpath -ldflags "-w -s \
					 -X github.com/gnames/$(PROJ_NAME)/pkg.Build=$(DATE) \
           -X github.com/gnames/$(PROJ_NAME)/pkg.Version=$(VERSION)"
FLAGS_REL = -trimpath -ldflags "-s -w \
						-X github.com/gnames/$(PROJ_NAME)/pkg.Build=$(DATE)"

GOCMD=go
GOINSTALL = $(GOCMD) install $(FLAGS_LD)
GOBUILD = $(GOCMD) build $(FLAGS_LD)
GORELEASE = $(GOCMD) build $(FLAGS_REL)
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get

all: install

test:
	go test -shuffle=on -count=1 -race -coverprofile=coverage.txt -covermode=atomic ./...

tools: deps
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

deps:
	@echo Download go.mod dependencies
	$(GOCMD) mod download;

build:
	$(GOCLEAN); \
	$(NO_C) $(GOBUILD);

buildrel:
	$(GOCLEAN); \
	$(NO_C) $(GORELEASE);

install:
	$(NO_C) $(GOINSTALL);

release: dockerhub
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=linux $(GORELEASE); \
	tar zcvf /tmp/$(PROJ_NAME)-$(VER)-linux.tar.gz $(PROJ_NAME); \
	$(GOCLEAN); \
	$(FLAGS_MAC_ARM) $(GORELEASE); \
	tar zcvf /tmp/$(PROJ_NAME)-$(VER)-mac.tar.gz $(PROJ_NAME); \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=windows $(GORELEASE); \
	zip -9 /tmp/$(PROJ_NAME)-$(VER)-win-64.zip $(PROJ_NAME).exe; \
	$(GOCLEAN);

docker: buildrel
	docker buildx build -t gnames/$(PROJ_NAME):latest -t gnames/$(PROJ_NAME):$(VERSION) .; \
	$(GOCLEAN);

dockerhub: docker
	docker push gnames/$(PROJ_NAME); \
	docker push gnames/$(PROJ_NAME):$(VERSION)
