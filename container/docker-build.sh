#!/usr/bin/env sh

docker build -t go2o -f ../Dockerfile ../

# 打包生成镜像并发布
docker tag go2o docker-base.56x.net:5020/go2o:latest
docker push docker-base.56x.net:5020/go2o:latest
