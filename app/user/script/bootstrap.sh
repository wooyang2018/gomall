#! /usr/bin/env bash
#dirname $0 命令会提取 $0 中文件路径的目录部分。
#获取当前脚本文件所在的目录路径，并将其赋值给变量 CURDIR
CURDIR=$(cd $(dirname $0); pwd)
echo "$CURDIR/bin/user"
exec "$CURDIR/bin/user"
