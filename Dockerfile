# Go2o Docker Image
# Version 1.0
# Author : jarrysix(jarrysix@gmail.com)
# Date : 2018-05-02 23:20


FROM golang:latest AS build
ENV BUILD=/gobuild
ENV GOPATH=$BUILD
WORKDIR $BUILD/src/go2o
COPY . ./

RUN go get -u github.com/golang/dep/cmd/dep && ln -s $BUILD/bin/* /bin
RUN dep ensure -v; rm -f $BUILD/bin/dep ;return 0
# 释放其他引用的包
#RUN cd ../ && tar xvzf go2o/gopkg.tar.gz
RUN CGO_ENABLED=0 GOOS=linux ARCH=amd64 GOFLAGS='-ldflags="-s -w"' go build bin/go2o-pub.go &&\
    CGO_ENABLED=0 GOOS=linux ARCH=amd64 GOFLAGS='-ldflags="-s -w"' go build bin/go2o-rpc.go
RUN ls -l bin



#RUN mkdir -p /opt/dist && cp -r mzld conf /opt/dist && ls -l /opt/dist
#
#FROM alpine
#MAINTAINER jarrysix
#LABEL vendor="YQTech"
#LABEL version="1.0.0"
#
#WORKDIR /app
#COPY --from=build /opt/dist/mzld /app/
#COPY --from=build /opt/dist/conf /app/conf
#RUN ln -s /app/mzl* /bin
#RUN apk update && apk add ca-certificates
#RUN echo "if [ ! -d '/data/conf' ];then cp -r /app/conf /data/;fi;"\
#    "mzld -conf /data/conf"> /docker-boot.sh
#VOLUME ["/data"]
#EXPOSE 7020
#CMD ["sh","/docker-boot.sh"]


