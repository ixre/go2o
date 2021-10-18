![Go2o](https://raw.githubusercontent.com/ixre/go2o/master/docs/mark.gif "GO2O")

[![Build Status](https://cloud.drone.io/api/badges/ixre/cms/status.svg)](https://cloud.drone.io/ixre/cms)

## What's Go2o? ##

Go2o is Google Go language binding domain-driven design (DDD) O2O open source implementation. Support Online Store ,
Offline stores; multi-channel (businesses), multi-store, merchandise, snapshots, orders, sales, payment, distribution
and other functions.

Project by a management center (including platform management center, business background, store background), online
store (PC shop, Handheld shops, micro-channel), the member center, open API in four parts.

Go2o using domain-driven design for business depth abstract, theoretical support in most sectors O2O scenarios. Through
open API, you can seamlessly integrate into legacy systems.

## Go2o 介绍 ##

Go2o是使用Golang语言结合领域驱动设计（DDD)的O2O实现。支持线上商店，线下门店；多渠道

（商户)、多门店、商品、规格SKU、快照、订单、促销、支付、配送等功能。

项目由管理中心(包括平台管理中心、商户后台、门店后台）、线上商店(PC商店、手持设备商店、微信)、

会员中心和通行证、服务四部分组成。

Go2o使用领域驱动设计对业务深度抽象，支持常见的O2O应用场景。通过Socket服务,可以开发安卓和IOS,

使用RPC服务可以方便与其他系统进行集成。

![Go2o](https://raw.githubusercontent.com/ixre/go2o/master/snapshot/dashboard.png "GO2O-DASHBOARD")


贡献代码请看： [todo list](https://github.com/ixre/go2o/tree/master/docs/dev/todo.md) |
[bug list](https://github.com/ixre/go2o/tree/master/docs/dev/bug.md)


========================================

感谢以下哥们和匿名捐助的朋友：

- *巍
- zhu***@126.com 
- 职业码农 
- 奋斗富三代

QQ群：**338164725**

**特别感谢: 领域驱动设计的专家-(腾讯)王立老师,我的良师益友;没有他,就没有这个项目!


## 运行 

### 准备运行环境

- 安装`PostgreSQL`并创建名为`go2o`的数据库, 下载数据备份文件:[go2o.sql](https://github.com/ixre/go2o/blob/master/docs/data/go2o.sql)进行还原
- 安装nats
- 安装etcd,单机创建单节点既可

###　编译运行
```
git clone https://github.com/ixre/go2o.git ./go2o
cd go2o && go mod tidy 
go run go2o-serve.go
```
指定参数,请参考:
```
Usage of go2o-serve:
  -apiport int
        api service port (default 1428)
  -conf string
         (default "app.conf")
  -d    run daemon
  -debug
        enable debug
  -endpoint etcd endpoints
  -help
        command usage
  -mqs string
        mq cluster address, like: 192.168.1.1:4222,192.168.1.2:4222 (default "127.0.0.1:4222")
  -port int
        gRPC service port (default 1427)
  -trace
        enable trace
  -v    print version
```

推荐使用[docker-compose](container/docker-compose.yaml)一键运行
```
docker-compose up -f container/docker-compose.yaml
```


