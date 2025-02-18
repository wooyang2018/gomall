#!/bin/bash

. scripts/list_app.sh
get_app_list

set -ex

readonly root_path=$(pwd)
for app_path in ${app_list[*]}; do
  #nohup 适合简单的后台任务，功能单一但轻量易用。
  #pm2 适合生产环境的复杂应用，功能强大但需要额外安装。
  #>output.log：将标准输出重定向到 output.log 文件。
  #2>&1：将标准错误重定向到标准输出（即 output.log 文件）。
  cd ${root_path}/${app_path}
  # 检查 log 目录是否存在
  if [ ! -d "log" ]; then
    mkdir log
  fi
  #nohup 命令不能直接处理环境变量，需要使用 env 命令设置。
  nohup env GO_ENV=dev go run . >log/nohup.out 2>&1 &
  cd ../..
done
