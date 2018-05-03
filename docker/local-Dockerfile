# Go2o Docker Image
# Version 1.0
# Author : jarrysix(jarrysix@gmail.com)
# Date : 2018-05-02 23:20

FROM alpine

MAINTAINER jarrysix
LABEL vendor="github.com/jsix/go2o"
LABEL version="1.0.0"

WORKDIR /app
COPY docker/go2o-* ./bin/
COPY conf ./conf
COPY docker/app_docker.conf ./conf/app.conf
RUN ln -s /app/bin/* /bin/ && ls -l .

RUN echo "if [ ! -f '/conf/app.conf' ];then cp -r /app/app.conf /conf/app.conf &&\
     cp -r /app/conf /conf;fi;"\
    "go2o-rpc -conf /conf/app.conf -conf-dir /conf/conf"> /docker-boot.sh

VOLUME ["/data","/conf"]
EXPOSE 14280 14281

#CMD ["sh","/docker-boot.sh"]

