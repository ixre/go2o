#!/usr/bin/env sh

docker build -t go2o -f ../Dockerfile ../

# 打包生成镜像并发布
docker tag go2o uhub.service.ucloud.cn/fze-registry/go2o:latest
docker push uhub.service.ucloud.cn/fze-registry/go2o:latest
