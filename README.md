
![Go2o](https://raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif "GO2O")




 

## What's Go2o? ##

Go2o is Google Go language binding domain-driven design (DDD) O2O open source implementation. Support Online Store
, Offline stores; multi-channel (businesses), multi-store, merchandise, snapshots, orders, sales, payment, distribution and other functions.

Project by a management center (including platform management center, business background, store background), online store (PC shop,
Handheld shops, micro-channel), the member center, open API in four parts.

Go2o using domain-driven design for business depth abstract, theoretical support in most sectors O2O scenarios.
Through open API, you can seamlessly integrate into legacy systems.

## Go2o 介绍 ##

Go2o是使用Golang语言结合领域驱动设计（DDD)的O2O实现。支持线上商店，线下门店；多渠道

（商户)、多门店、商品、规格SKU、快照、订单、促销、支付、配送等功能。


项目由管理中心(包括平台管理中心、商户后台、门店后台）、线上商店(PC商店、手持设备商店、微信)、

会员中心和通行证、服务四部分组成。


Go2o使用领域驱动设计对业务深度抽象，支持常见的O2O应用场景。通过Socket服务,可以开发安卓和IOS,

使用Rpc服务可以扩展子系统。



![Go2o](https://raw.githubusercontent.com/jsix/go2o/master/snapshot/merchant.png "GO2O-Merchant")


更多系统截图见:#snapshot#目录

## 项目说明 ##


__最后提交时间: 2016-12-20 __

__代码已重构完毕,见develop分支, 新的代码库不在包含UI, UI见分支v0.1.1__



------------------------
贡献代码请看： [todo list](https://github.com/jsix/go2o/tree/master/docs/dev/todo.md) |
[bug list](https://github.com/atnet/go2o/tree/master/docs/dev/bug.md)


请支持开源，不做伸手党，不拿来主意！

========================================


感谢以下哥们和匿名捐助的朋友：

*巍
zhu***@126.com
职业码农
奋斗富三代


QQ群：**338164725**


#### MAC下运行请先设置最大连接数:

    sudo sysctl -w kern.ipc.somaxconn=4096

## Deploy ##
### 1. Import database ###
> Create new mysql db instance named "go2o"
 and import data use mysql utility.
 Database backup file is here : [go2o.sql](https://github.com/jsix/go2o/blob/master/docs/data/go2o.sql)

### 2.Complied ###
	git clone https://github.com/jsix/go2o.git /home/usr/go/src/go2o
	export GOPATH=$GOPATH:/home/usr/go/
	cd /home/usr/go/src/go2o
	go build go2o-serve.go
	go build go2o-daemon.go
	go build go2o-tcpserve.go

### 2.Running Service ###
	Usage of ./go2o-serve:
		 -conf string
             	 (default "app.conf")
           -d	
                run daemon
           -r   
                run rpc server
           -debug
             	enable debug
           -trace
             	enable trace
           -help
             	command usage
           -port int
             	web server port (default 14190)
           -restport int
             	rest api port (default 14191)

	Usage of ./go2o-daemon:
		-debug = false : enable debug

	Usage of ./go2o-tcpserve:
	  -conf string
        	 (default "app.conf")
      -l	log output
      -port int
        	 (default 14197)

### 3.Add http proxy pass for Nginx ###
	server {
            listen          80;
            server_name     static.ts.com;
            root    /home/usr/go/src/go2o/static;
    	location ~* \.(eot|ttf|woff|woff2|svg)$ {
          		add_header Access-Control-Allow-Origin *;
      	}
    }

    server {
            listen          80;
            server_name     img.ts.com;
            root            /home/usr/go/src/go2o/uploads;
    	location ~* \.(eot|ttf|woff|woff2|svg)$ {
          		add_header Access-Control-Allow-Origin *;
      	}
    }

    server {
            listen          80;
            server_name     *.ts.com;
            client_max_body_size    10m;  
            location / {
                    proxy_pass   http://localhost:14190;
                    proxy_set_header X-Real-IP $remote_addr;
                    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                    proxy_set_header Host $host;
                    proxy_redirect   off;
            }
    }



### 4.Add test hosts ###
> echo   127.0.0.1    go2o.ts.com static.ts.com img.ts.com mch.ts.com hapi.ts.com 
u.ts.com mu.ts.com passport.ts.com mpp.ts.com
 master.ts.com zy.ts.com whs.ts.com >> /etc/hosts

## Access Entry ##

### WebMaster ##
master.ts.com

account: go2o / 123456

### Merchant Management ###
mch.ts.com

account: go2o / 123456

### Member Center ###
u.ts.com

### Merchant Sales ###
go2o.ts.com


