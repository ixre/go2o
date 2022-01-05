#!/usr/bin/env bash

TABLE_PREFIX=$1

if [[ $(whereis tto) = 'tto:' ]]; then
  echo '未安装tto客户端,请运行安装命令： curl -L https://raw.githubusercontent.com/ixre/tto/master/install | sh'
fi
tto -m go -conf tto.conf -table ${TABLE_PREFIX} -clean
find output -name "*.go" -print0 |  xargs -0 sed -i ':label;N;s/This.*Copy/Copy/g;b label'
