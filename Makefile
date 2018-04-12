NAME ?= kube-custodian
LDFLAGS ?= -ldflags="-s -w -X main.version=$(VERSION) -X main.revision=$(REVISION)"
GOSRC = $(shell find -name '*.go'|fgrep -v vendor/)
GOPKGS = $(shell glide novendor)


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


.PHONY: all build lint test clean
