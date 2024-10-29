# spinit <!-- omit in toc -->
spinit is a command-line tool to build and deploy apps to local Kubernetes clusters. It enables the user to avoid writing `yaml` configurations for the app and seamlessly load docker images to every node in the cluster.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Prerequisites](#prerequisites)
- [Supported Platforms](#supported-platforms)
- [Installation](#installation)
- [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Prerequisites
`spinit` uses [`buildkit`](https://github.com/moby/buildkit/tree/master) to build docker images. You can follow the instructions [here](https://github.com/moby/buildkit/tree/master?tab=readme-ov-file#quick-start) on how to set `buildkit` up.

## Supported Platforms
Currently `spinit` has only been tested to work on Linux.

Additionally, I am using [`k3s`](https://github.com/k3s-io/k3s) as the Kubernetes distribution with the default `containerd` container runtime.

## Installation


## Contributing
See [CONTRIBUTING.md](./CONTRIBUTING.md)