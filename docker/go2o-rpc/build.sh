#!/usr/bin/env sh
goos=linux
arch=amd64
cd ../../
CGO_ENABLED=0 GOOS=${goos} GOARCH=${arch} go build -ldflags "-w" ../bin/go2o-pub.go
#CGO_ENABLED=0 GOOS=${goos} GOARCH=${arch} go build -ldflags "-w" ../bin/go2o-tcpserve.go
CGO_ENABLED=0 GOOS=${goos} GOARCH=${arch} go build -ldflags "-w" ../bin/go2o-rpc.go
docker build -t go2o -f ./local-Dockerfile ../
#rm -rf go2o*

# 打包生成镜像并发布
#docker build -t mzl-api-hub . &&\
#docker tag mzl-api-hub docker.tech.meizhuli.net:5000/mzl-api-hub &&\
#docker push docker.tech.meizhuli.net:5000/mzl-api-hub

