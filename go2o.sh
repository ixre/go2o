#!/bin/bash

  ###   ###   ###   ###
 #     #  ##    #  #  ##
#     #    #    # #    #
#  #  #   #   ##  #   #
#  #  #   #  #    #   #
 ###   ###   ###   ###

prod="go2o"
conf="conf"
debug=0
log_dir="./logs"
cmd=""


while getopts "s:c:dev" args
do
    case ${args} in
        v) ./go2o-serve -v && exit 1;;
        e) echo "extra update files ..." && \
           tar xvzf ${prod}-update.tar.gz;;
        d) debug=1;;
        s) cmd=${OPTARG};;
        c) conf=${OPTARG};;
        \?) echo "usage args [-s][-v] [-e] [-c] [-d]"
            echo " -s : exec command,[stop|start|clean|restart]"
            echo " -e : extra update zip file"
            echo " -c : config file"
            echo " -v : print version"
            exit 1;;
    esac
done

# clean cache
if [[ ${cmd} = "clean" ]];then
    ./master-serve -conf=app.conf -clean
fi
# stop service
if [[ ${cmd} = "stop" ]] || [[ ${cmd} = "restart" ]];then
   proc_arr=("go2o-serve" "master-serve" "go2o-tcp" "go2o-rpc")
    for proc in ${proc_arr[@]}
    do
       ps -ef| grep ${proc}| awk '{print $2}'| xargs kill -15
    done
fi
# start service
if [[ ${cmd} = "start" ]] || [[ ${cmd} = "restart" ]];then
   mkdir -p ${log_dir}
   chmod a+x go2o-serve master-serve go2o-tcpserve go2o.sh
   if [[ ${debug} = 1 ]];then
        nohup ./go2o-serve -debug -conf=app.conf -d -r>${log_dir}/go2o.log 2>&1 &
        nohup ./master-serve -debug -conf=app.conf>${log_dir}/master.log 2>&1 &
        nohup ./go2o-tcpserve -debug -conf=app.conf>${log_dir}/tcp.log 2>&1 &
   else
        nohup ./go2o-serve -conf=app.conf -d -r>${log_dir}/go2o.log 2>&1 &
        nohup ./master-serve -conf=app.conf>${log_dir}/master.log 2>&1 &
        nohup ./go2o-tcpserve -conf=app.conf>${log_dir}/tcp.log 2>&1 &
   fi
   sleep 2 && echo "boot success" && exit 1
fi

