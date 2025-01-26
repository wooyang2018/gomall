#!/bin/bash

. scripts/list_app.sh

get_app_list

readonly root_path=$(pwd)
for app_path in ${app_list[*]}; do
  # go vet 是 Go 语言官方提供的一个静态分析工具，用于检查 Go 代码中的常见错误和潜在问题。
  # 它不会检查代码风格或格式化问题，而是专注于发现代码中的逻辑错误、可疑的代码模式。
  go vet ${root_path}/${app_path}
done
