#!/bin/bash
#set -ex：这行代码设置了两个Shell选项：
#e：如果任何命令执行失败（返回非零状态码），脚本将立即退出。
#x：在执行每个命令之前，先打印出命令本身及其参数。
set -ex

source scripts/list_app.sh

get_app_list

#安装golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
for app_path in ${app_list[*]}; do
  #-E gofumpt
  #-E 表示启用指定的 linter（代码检查工具）。
  #gofumpt 是一个严格的 Go 代码格式化工具，基于 gofmt，但提供了更多的格式化规则。
  #
  #--path-prefix=.
  #指定路径前缀为当前目录（.）。
  #
  #--fix
  #自动修复代码中可修复的问题。
  #
  #--timeout=5m
  #设置 golangci-lint 的运行超时时间为 5 分钟。
  cd ${app_path} && golangci-lint run -E gofumpt --path-prefix=. --fix --timeout=5m && cd ../../
done
