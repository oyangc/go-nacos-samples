#!/bin/bash

dir=$(
   cd "$(dirname "$0")"
   pwd
)

projectJarFilePath=${dir}/go-nacos

if [ ! -f ${projectJarFilePath} ] ; then
  go build go-nacos.go
fi
#go build go-nacos.go
chmod +x go-nacos