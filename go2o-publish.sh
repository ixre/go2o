#!/bin/bash


  ###   ###   ###   ###
 #     #  ##    #  #  ##
#     #    #    # #    #
#  #  #   #   ##  #   #
#  #  #   #  #    #   #
 ###   ###   ###   ###


goos=linux
arch=amd64
env="0"
server="127.0.0.1"
ssh_port=22
user=root
root_user=root
app_dir="/data/www/go2o/"

echo "---------------------------"
echo "*** GO2O PUBLISH SCRIPT ***"
echo "---------------------------"
echo " Select publish target environment:"
echo " 1): Development"
echo " 2): Nightly & Beta"
echo " 3): Release"

package_name="update.tar.gz"
zip_bin="go2o-serve go2o-tcpserve master-serve go2o-rpc"
zip_res="public mobile conf/query uploads/res go2o.sh"
boot_sh="sudo ./go2o.sh restart"

# 获取发布环境
while true
do
   # echo -n 不换行
   echo -n "Please input : "
   # 读取发布环境
   read env
    if [[ ${env} = "1" ]] || [[ ${env} = "2" ]] || [[ ${env} = "3" ]];then
        break
    else
        echo "Please retype correct index:"
    fi
done

# 开发环境
if [[ ${env} = "1" ]];then
    arch=arm
    server="dev.go2o.to2.net"
    user=pi
    root_user=pi
    ssh_port=22
    app_dir="/data/www/go2o"
    echo "Selected : Development"
fi

# 测试或每夜环境
if [[ ${env} = "2" ]];then
     server="192.168.4.201"
     user=flm
     root_user=flm
     app_dir="/data/www/go2o"
     echo "Selected : Nightly & Beta"
fi

# 发布版
if [[ ${env} = "3" ]];then
     server="official.go2o.to2.net"
     user=flm
     app_dir="/data/www/go2o"
     echo "Selected : Release"
fi

echo ""
echo "[ Setup 1 ]: compile program "
echo "Please confirm compile : [Y/N]"

read compile
if [[ ${compile} = "Y" ]] || [[ ${compile} = "y" ]];then
    echo "compiling ..."
    CGO_ENABLED=0 GOOS=${goos} GOARCH=${arch} go build -ldflags "-w" go2o-serve.go
    CGO_ENABLED=0 GOOS=${goos} GOARCH=${arch} go build -ldflags "-w" go2o-tcpserve.go
    CGO_ENABLED=0 GOOS=${goos} GOARCH=${arch} go build -ldflags "-w" master-serve.go
    CGO_ENABLED=0 GOOS=${goos} GOARCH=${arch} go build -ldflags "-w" go2o-rpc.go
    #CGO_ENABLED=0 GOOS=${goos} GOARCH=${arch} go build -ldflags "-w" merchant-serve.go
    #CGO_ENABLED=0 GOOS=${goos} GOARCH=${arch} go build -ldflags "-w" pub-serve.go
else
    echo "  skipping compile" && echo ""
fi

echo "[ Setup 2 ]: zipping tar package ..."
sleep 1


echo "Include ""conf"" folder [Y/N] ? "
read zipConf

if [[ ${zipConf} = "Y" ]] || [[ ${zipConf} = "y" ]];then
    zip_res=${zip_res}" conf/core"
else
    echo "  skipping conf folder" && echo ""
fi

echo "Include ""static"" folder [Y/N] ? "
read zipStatic
if [[ ${zipStatic} = "Y" ]] || [[ ${zipStatic} = "y" ]];then
    zip_res=${zip_res}" static"
else
    echo "  skipping static folder" && echo ""
fi

zip_res=${zip_bin}" "${zip_res}
echo ${zip_res}
tar cvzf ../${package_name} ${zip_res}


echo "[ Setup 3 ]: upload tar package to server ..."

scp ../${package_name} ${user}@${server}:${app_dir}

echo "[ Setup 4 ]: restart server"
ssh -t -p ${ssh_port} ${root_user}@${server} "cd ${app_dir} && tar xvzf ${package_name} && ${boot_sh}"

echo "[ Setup 5 ]: cleaning ..."
rm  ${zip_bin} ../${package_name}

echo "Configurations, publish successfully!"

exit 0
