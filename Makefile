NAME ?= kube-custodian
LDFLAGS ?= -ldflags="-s -w -X main.version=$(VERSION) -X main.revision=$(REVISION)"
VERSION ?= master
GOARCH ?= amd64
GOFLAGS = -mod=vendor # ensure we always honour vendor/

GOPKGS = ./...

ifeq ($(GOARCH), arm)
DOCKERFILE_SED_EXPR?=s,FROM bitnami/minideb:stretch,FROM armhf/debian:stretch-slim,
DOCKER_IMG_FULL=$(DOCKER_IMG):arm-$(VERSION)
else ifeq ($(GOARCH), arm64)
DOCKERFILE_SED_EXPR?=s,FROM bitnami/minideb:stretch,FROM aarch64/debian:stretch-slim,
DOCKER_IMG_FULL=$(DOCKER_IMG):arm64-$(VERSION)
else
DOCKERFILE_SED_EXPR?=
DOCKER_IMG_FULL=$(DOCKER_IMG):$(VERSION)
endif

DOCKER_REPO ?= quay.io
DOCKER_IMG ?= $(DOCKER_REPO)/bitnami-labs/kube-custodian

# -mod=vendor is accepted only if Go Modules are turned on
export GO111MODULE = on

all: build

build: bin/$(NAME)

bin/$(NAME):
	GOARCH=$(GOARCH) go build $(GOFLAGS) $(LDFLAGS) -o bin/$(NAME)

check: lint vet inef

dep:
	go mod tidy
	go mod vendor

lint:
	golint $(GOPKGS)

vet:
	go vet $(GOFLAGS) $(GOPKGS)

inef:
	ineffassign .

fmt:
	go fmt $(GOPKGS)

test:
	go test -v $(GOFLAGS) -cover -count=1 $(GOPKGS)

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

.PHONY: all build bin/* check lint vet inef test clean docker-build docker-push docker-clean
