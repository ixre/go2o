# Go2o Docker Image
# Version 1.0
# Author : jarrysix(jarrysix@gmail.com)
# Date : 2019-06-12 13:20

FROM golang:latest AS build
ENV BUILD=/gobuild
ENV GOPATH=$BUILD
WORKDIR $BUILD/src/go2o
COPY . ./

ENV GO111MODULE=on
ENV GOPROXY=https://athens.azurefd.net
RUN rm -rf go.sum && sed -i 's/replace/\/\/replace/' go.mod && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -ldflags='-s -w' go2o-serve.go && \
    CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -ldflags='-s -w' go2o-rpc.go

RUN mkdir -p /opt/go2o/dist && \
    cp -r go2o-serve go2o-rpc LICENSE README.md app.conf /opt/go2o/dist

FROM alpine
MAINTAINER jarrysix
LABEL vendor="Go2o"
LABEL version="1.0.0"

WORKDIR /app
COPY --from=build /opt/go2o/dist/* /app/

RUN ln -s /app/go2o-* /bin && \
    apk update && apk add ca-certificates && \
    echo "if [ ! -f '/data/app.conf' ];then cp -r /app/app.conf /data;fi;"\
    "go2o-serve -conf /data/app.conf  -r -d"> /docker-boot.sh
VOLUME ["/data"]
EXPOSE 1427 1428
CMD ["sh","/docker-boot.sh"]


