NAME ?= kube-custodian
LDFLAGS ?= -ldflags="-s -w -X main.version=$(VERSION) -X main.revision=$(REVISION)"
GOSRC = $(shell find -name '*.go'|fgrep -v vendor/)
GOPKGS = $(shell glide novendor)

DOCKER_REPO ?= quay.io
DOCKER_IMG ?= $(DOCKER_REPO)/jjo/kube-custodian


all: build

build: bin/$(NAME)

bin/$(NAME): $(GOSRC)
	go build $(LDFLAGS) -o bin/$(NAME)

lint:
	golint $(GOPKGS)

test:
	go test -v -cover $(GOPKGS)

clean:
	rm -fv bin/$(NAME)


docker-build:
	docker build -t $(DOCKER_IMG) .

docker-push:
	docker push $(DOCKER_IMG)



.PHONY: all build lint test clean
