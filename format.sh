#!/bin/bash

go env -w GO111MODULE=on 
go env -w GOPROXY=https://goproxy.cn,direct

go mod tidy
goimports -w .
gofmt -w .
