#!/bin/bash

app_list=() #初始化为空数组

get_app_list() {
  #local关键字用于声明局部变量，这意味着idx只在get_app_list函数内部可见
  local idx=0
  #app/*是一个通配符表达式，表示app目录下的所有文件和子目录。
  for d in app/*; do
    #-d是一个测试操作符，用于检查文件是否是一个目录
    if [ -d "$d" ]; then
      app_list[idx]=$d
      idx+=1
    fi
  done
}
