# Building Kubernetes on Windows

## Overview
This directory is for building Kubernetes on Windows. The build process is currently on host because container is not available on any released versions of Windows yet. The script is written in python so that later can be used on Linux as well.

This is part of an effort to port Kubernetes on Windows Server 2016 which supports containers.

## Requirements

1. [Go 1.5+](https://golang.org/) installed and included in $PATH
2. [Godep tool](https://github.com/tools/godep) installed and included in $PATH
3. [Python 3.5+](https://www.python.org) installed
4. Kubernetes project to be under $GOPATH/src/k8s.io/kubernetes (To be revisited)

## Build

Run `build.py`