#!/bin/bash

package_name="go2o_deploy_tmp.tar.gz"

goods=linux
arch=amd64
server="127.0.0.1"
user=root
root_user=root
ssh_port=22
env = "0"

echo "---------------------------"
echo "*** GO2O PUBLISH SCRIPT ***"
echo "---------------------------"
echo " Select publish target environment:"
echo " 1): Development"
echo " 2): Nightly & Beta"
echo " 3): Release"


# 获取发布环境
while true
do
   # echo -n 不换行
   echo -n "Please input : "
   # 读取发布环境
   read env
    if [ ${env} = "1" ] || [ ${env} = "2" ] || [ ${env} = "3" ];then
        break
    else
        echo "Please retype correct index:"
    fi
done

# 开发环境
if [ ${env} = "1" ];then
    arch=arm
    server="www.flm-dev.redvp.com"
    user=pi
    root_user=pi
    ssh_port=22
    app_dir="/home/${user}/www/flm-dev/"
    echo "Selected : Development"
fi

# 测试或每夜环境
if [ ${env} = "2" ];then
     server="192.168.4.201"
     user=flm
     root_user=flm
     app_dir="/home/${user}/www/flm/"
     echo "Selected : Nightly & Beta"
fi

# 发布版
if [ ${env} = "3" ];then
     server="official.ns.fanleme.net"
     user=flm
     app_dir="/home/${user}/www/flm/"
     echo "Selected : Release"
fi

boot_sh="cd ${app_dir} &&tar xvzf /home/${user}/${package_name} && sudo ./go2o.sh restart"

echo ""
echo "[ Setup 1 ]: compile program "
echo "Please confirm compile : [Y/N]"

read compile
if [ ${compile} = "Y" ]||[ ${compile} = "y" ];then
    echo "program compiling ..."
    CGO_ENABLED=0 GOOS=${goods} GOARCH=${arch} go build -v go2o-serve.go
    CGO_ENABLED=0 GOOS=${goods} GOARCH=${arch} go build -v go2o-tcpserve.go
    CGO_ENABLED=0 GOOS=${goods} GOARCH=${arch} go build -v master-serve.go
    CGO_ENABLED=0 GOOS=${goods} GOARCH=${arch} go build -v go2o-rpc.go
    #CGO_ENABLED=0 GOOS=${goods} GOARCH=${arch} go build -v merchant-serve.go
    #CGO_ENABLED=0 GOOS=${goods} GOARCH=${arch} go build -v pub-serve.go
else
    echo "  skipping compile"
fi



echo "[ Setup 2 ]: zipping tar package ..."
sleep 1

zipRes="go2o-serve go2o-tcpserve master-serve public/views conf/query uploads/res go2o.sh"

echo "Include conf folder [Y/N] ? "
read zipConf

if [ ${zipConf} = "Y" ]||[ ${zipConf} = "y" ];then
    zipRes=${zipRes}" conf/core"
fi

echo "Include public/static folder [Y/N] ? "
read zipStatic
if [ ${zipStatic} = "Y" ]||[ ${zipStatic} = "y" ];then
    zipRes=${zipRes}" public/static"
fi

echo ${zipRes}
tar cvzf ../${package_name} ${zipRes}

echo "[ Setup 3 ]: upload tar package to server ..."

scp ../${package_name} ${user}@${server}:/home/${user}/

echo "[ Setup 4 ]: restart server"

ssh -t -p ${ssh_port} ${root_user}@${server} "${boot_sh}"

echo "[ Setup 5 ]: cleaning ..."
rm go2o-serve go2o-tcpserve master-serve ../${package_name}

echo "Configurations, publish successfully!"

exit 0
