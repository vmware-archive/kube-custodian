[![Build Status](https://travis-ci.org/bitnami-labs/kube-custodian.svg?branch=master)](https://travis-ci.org/bitnami-labs/kube-custodian)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitnami-labs/kube-custodian)](https://goreportcard.com/report/github.com/bitnami-labs/kube-custodian)

kube-custodian - Kubernetes cleanup tool

## Purpose

On Kubernetes clusters for development, it's pretty common to have
workloads that become forgotten by developers, holding resources thus
potentially voiding new workloads from scheduling.

`kube-custodian` will mark for later deletion those workloads (Deployments,
StatefulSets, Jobs, Pods) lacking `--skip-labels` unless their
namespace has them.

The main subcommand is `run`, usage help from its `Flags` section:
```bash
$ kube-custodian run --help
Scan Kubernetes objects, mark for deletion (via annotation), delete those already "expired"

Usage:
  kube-custodian run [flags]

Flags:
      --cleanup-tag                Untag resources from later deletion
      --delete-tagged              Delete tagged resources, after their Tag TTL has passed (default true)
      --skip-labels strings        Labels required for resources to be skipped from scanning (default [created_by])
      --skip-namespace-re string   Regex of namespaces to skip, typically 'system' ones and alike (default "kube-.*|.*(-system|monitoring|logging|ingress)")
      --tag-for-deletion           Tag resources for later deletion (default true)
      --tag-ttl string             Time to live after marked, before deletion (default "24h")
[...]
```

```bash
For example, to mark for later deletion all workloads not having the `created_by`
label, run:

```bash
$ kube-custodian -v --namespace=default --dry-run run --tag-ttl 24h --skip-labels created_by

```

Obviously, remove `--dry-run` to _actually_ mark them :), it'll add an
annotation as

  `kube-custodian.bitnami.com/expiration-mark: <current epoch secs>`

Then, 24h later same run as above will:
- Update any new workload without this above annotation
- Delete all workloads for which:

  (`kube-custodian.bitnami.com/expiration-mark` + `tag-ttl`) >= `now`


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
  -v --namespace=default --dry-run run --required-labels created_by
```

## Source

Based on source code from https://github.com/ksonnet/kubecfg,
https://github.com/bitnami-labs/kubewatch.
