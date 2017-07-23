#!/bin/bash

action=$1

# tar deploy package
if [[ ${action} = "tar" ]];then
    file=$2
    if [[ ${action} = "" ]];then
        echo "no such tar file"
    else
        tar xvzf ${file}
    fi
    exit 0
fi

# clean cache
if [[ ${action} = "clean" ]];then
    ./master-serve -conf=app.conf -clean
fi

# stop service
if [[ ${action} = "stop" ]];then
   killall go2o-serve
   killall master-serve
   killall go2o-tcpserve
   killall go2o-rpc
fi

# start service
if [[ ${action} = "start" ]];then
   nohup  ./go2o-serve -conf=app.conf -d -r>logs/go2o.log 2>&1 &
   nohup ./master-serve -conf=app.conf>logs/master.log 2>&1 &
   nohup ./go2o-tcpserve -conf=app.conf>logs/tcp.log 2>&1 &
   # 暂停2s等待服务启动成功
   sleep 2
   echo "success"
fi

# reboot service
if [[ ${action} = "restart" ]];then
    # 停止服务
    killall go2o-serve
    killall master-serve
    killall go2o-tcpserve
    # 修改权限
    chmod o+x go2o-serve master-serve go2o-tcpserve go2o.sh
    # 启动服务
    nohup  ./go2o-serve -conf=app.conf -d -r>logs/go2o.log 2>&1 &
    nohup ./master-serve -conf=app.conf>logs/master.log 2>&1 &
    nohup ./go2o-tcpserve -conf=app.conf>logs/tcp.log 1>&1 &
    # 暂停2s等待服务启动成功
    sleep 2
    echo "success"
fi

