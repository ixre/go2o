#!/bin/bash


# generate api docs
rm -rf api_docs && \
apidoc -i $(find . -name 'apidoc.json' -print -quit|xargs dirname) -o api_docs/


# 生成服务代码
# author : liuming
# data   : 2020-09-03 11:00


if [[ $GO2O_JAVA_HOME != "" ]];then java_target_path=$GO2O_JAVA_HOME;fi;

thrift_path=$(find . -name "service.thrift" -print -quit)

cmd=$1

rm -rf ./go2o/core/service/proto
protoc -I ./proto  --go_out=plugins=grpc:proto ./proto/*.proto


#if [[ ${cmd} = "all" ]] || [[ ${cmd} = "format" ]];then
#	cd ${go_target_path}
#fi
