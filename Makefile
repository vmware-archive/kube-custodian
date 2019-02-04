NAME ?= kube-custodian
LDFLAGS ?= -ldflags="-s -w -X main.version=$(VERSION) -X main.revision=$(REVISION)"
VERSION ?= master
GOARCH ?= amd64

GOPKGS = . ./cmd/... ./pkg/...
GOSRC = $(shell go list -f '{{$$d := .Dir}}{{range .GoFiles}}{{$$d}}/{{.}} {{end}}' $(GOPKGS))

ifeq ($(GOARCH), arm)
DOCKERFILE_SED_EXPR?=s,FROM alpine:,FROM multiarch/alpine:armhf-v,
DOCKER_IMG_FULL=$(DOCKER_IMG):arm-$(VERSION)
else ifeq ($(GOARCH), arm64)
DOCKERFILE_SED_EXPR?=s,FROM alpine:,FROM multiarch/alpine:aarch64-v,
DOCKER_IMG_FULL=$(DOCKER_IMG):arm64-$(VERSION)
else
DOCKERFILE_SED_EXPR?=
DOCKER_IMG_FULL=$(DOCKER_IMG):$(VERSION)
endif

DOCKER_REPO ?= quay.io
DOCKER_IMG ?= $(DOCKER_REPO)/bitnami-labs/kube-custodian


all: build

build: bin/$(NAME)

bin/$(NAME): $(GOSRC)
	GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/$(NAME)

check: lint vet inef

lint:
	golint $(GOPKGS)

vet:
	go vet $(GOPKGS)

inef:
	ineffassign .

fmt:
	gofmt -s -w $(GOSRC)

test:
	go test -v -cover -count=1 $(GOPKGS)

clean:
	rm -fv bin/$(NAME)


docker-build: Dockerfile.$(GOARCH).run
	docker build --build-arg SRC_TAG=$(VERSION) --build-arg ARCH=$(GOARCH) -t $(DOCKER_IMG_FULL) -f $(^) .

Dockerfile.%.run: Dockerfile
	@sed -e "$(DOCKERFILE_SED_EXPR)" Dockerfile > $(@)


docker-push:
	docker push $(DOCKER_IMG)

docker-clean:
	docker image rm $(DOCKER_IMG)


multiarch-setup:
	docker run --rm --privileged multiarch/qemu-user-static:register
	dpkg -l qemu-user-static

.PHONY: all build check lint vet inef test clean docker-build docker-push docker-clean
