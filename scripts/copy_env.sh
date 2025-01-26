#!/bin/bash
#使用.命令（也称为source命令）来执行另一个脚本文件scripts/list_app.sh。
#这意味着脚本文件中的命令将在当前Shell进程中执行，而不是在一个新的子Shell进程中执行。
. scripts/list_app.sh

get_app_list
#定义root_path只读变量
readonly root_path=$(pwd)
#${app_list[*]} 是一种数组扩展的语法，用于获取数组 app_list 中的所有元素，
#并将它们作为一个字符串返回。这个字符串中的元素之间默认以空格分隔。建议使用 app_list[@]，
#这样可以确保每个元素都被正确地处理，即使元素中包含空格或其他特殊字符。
for app_path in ${app_list[*]}; do
  if [[ "${app_path}" = "app/common" ]]; then
    continue
  fi
  #-e用于检查文件是否存在
  if [[ -e "${app_path}/.env" ]]; then
    continue
  fi
  echo "copy ${app_path} env file"
  cp "${app_path}/.env.example" "${app_path}/.env"
  echo "Done! Please replace the real value"
done
