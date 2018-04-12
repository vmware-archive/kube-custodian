FROM golang:1.9.4-alpine as build

ARG SRC_REPO=github.com/jjo/kube-custodian
ARG SRC_TAG=master
ARG ARCH=amd64

RUN apk update && apk add git ca-certificates

# Unfortunately bitly/oauth2_proxy is not vendored - FYI built OK on 2018-03-29
# Checkout specific version
RUN go get -d ${SRC_REPO}
RUN git -C ${GOPATH}/src/${SRC_REPO} checkout -b build-${ARCH}-${SRC_TAG} ${SRC_TAG}
RUN GOARCH=${ARCH} go get ${SRC_REPO}
RUN find /go/bin -name kube-custodian -type f | xargs -I@ install @ /

FROM alpine:3.7
MAINTAINER JuanJo Ciarlante <juanjosec@gmail.com>

USER 1001
COPY --from=build /kube-custodian /opt/kube-custodian/bin/kube-custodian
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /usr/share/ca-certificates /usr/share/ca-certificates

ENTRYPOINT [ "/opt/kube-custodian/bin/kube-custodian" ]
