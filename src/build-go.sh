#!/bin/bash

dir=$(
   cd "$(dirname "$0")"
   pwd
)

cd ${dir}/main/go/

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build go-nacos.go

chmod +x go-nacos