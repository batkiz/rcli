#!/usr/bin/env -S just --justfile

set windows-shell := ["powershell.exe", "-c"]
set shell := ["bash"]

build:
    go build

lint:
    golangci-lint run
