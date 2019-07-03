FROM golang:1.11.4 as build

ARG SRC_REPO=github.com/bitnami-labs/kube-custodian
ARG SRC_TAG=master
ARG ARCH=amd64

RUN apt-get upgrade

COPY . /go/src/${SRC_REPO}
RUN GOARCH=${ARCH} go get ${SRC_REPO}
RUN find /go/bin -name kube-custodian -type f | xargs -I@ install @ /

FROM bitnami/minideb:stretch
MAINTAINER JuanJo Ciarlante <juanjosec@gmail.com>

USER 1001
COPY --from=build /kube-custodian /opt/kube-custodian/bin/kube-custodian
COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /usr/share/ca-certificates /usr/share/ca-certificates

ENTRYPOINT [ "/opt/kube-custodian/bin/kube-custodian" ]
