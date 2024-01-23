SHELL=/usr/bin/env bash
GO_BUILD_IMAGE?=golang:1.20

.PHONY: all
all: build

.PHONY: build
build:
	go build -o spade-tenant-svc

.PHONE: clean
clean:
	rm -f spade-tenant-svc

install:
	install -C -m 0755 spade-tenant-svc /usr/local/bin