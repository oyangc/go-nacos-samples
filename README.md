
# go-nacos-samples
go版nacos初始化配置脚本

## go发布nacos

    go拉取nacos配置，如果远程配置不存在配置文件直接post配置文件到远程nacos。

## maven编译打包

    windows环境下测试，测试通过之后，使用maven编译go脚本成linux上可运行的脚本。

## go脚本编译、运行

    Nacos Open API 指南
    https://nacos.io/zh-cn/docs/open-api.html

    #打包成linux版本
    1. 编译
    set CGO_ENABLED=0
    set GOARCH=amd64
    set GOOS=linux

    go build go-nacos.go

    CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build go-nacos.go

    2.运行
	chmod +x go-nacos
    ./go-nacos
