#!/usr/bin/env -S just --justfile

set windows-shell := ["powershell.exe", "-c"]
set shell := ["bash"]

set dotenv-load

fmt:
    go fmt

lint: fmt
    golangci-lint run

build: lint
    go build

dev: fmt
    go build

