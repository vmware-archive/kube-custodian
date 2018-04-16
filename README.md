[![Build Status](https://travis-ci.org/jjo/kube-custodian.svg?branch=master)](https://travis-ci.org/jjo/kube-custodian)
[![Go Report Card](https://goreportcard.com/badge/github.com/jjo/kube-custodian)](https://goreportcard.com/report/github.com/jjo/kube-custodian)

kube-custodian - Kubernetes cleanup tool

## Purpose

On Kubernetes clusters for development, it's pretty common to have
workloads that become forgotten by developers, holding resources thus
potentially voiding new workloads from scheduling.

`kube-custodian` will delete those workloads (Deployments,
StatefulSets, Jobs, Pods) based on some labelling "condition".

For example, to delete all workloads not having the "created_by"
label:

```bash
$ kube-custodian -v --namespace=default --dry-run delete --required-labels created_by
```

## Install

Install it from source with:

```bash
$ go get github.com/bitnami-labs/kube-custodian
```

## Docker image

You can build your own docker image (`docker build -t YOU/kube-custodian .`)
or use pre-built as:

```
docker run -it -v $HOME/.kube:/.kube quay.io/jjo/kube-custodian \
  -v --namespace=default --dry-run delete --required-labels created_by
```

## Source

Based on source code from https://github.com/ksonnet/kubecfg,
https://github.com/bitnami-labs/kubewatch.
