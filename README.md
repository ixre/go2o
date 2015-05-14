
![Go2o](https://raw.githubusercontent.com/atnet/go2o/master/docs/mark.gif "GO2O")

## What's Go2o? ##

Go2o is Google Go language binding domain-driven design (DDD) O2O open source implementation. Support Online Store
, Offline stores; multi-channel (businesses), multi-store, merchandise, snapshots, orders, sales, payment, distribution and other functions.

Project by a management center (including platform management center, business background, store background), online store (PC shop,
Handheld shops, micro-channel), the member center, open API in four parts.

Go2o using domain-driven design for business depth abstract, theoretical support in most sectors O2O scenarios.
Through open API, you can seamlessly integrate into legacy systems.

## Go2o 介绍 ##

Go2o是Google Go语言结合领域驱动设计（DDD)的开源O2O实现。支持线上商店，线下门店；多渠道

（商户)、多门店、商品、快照、订单、促销、支付、配送等功能。


项目由管理中心(包括平台管理中心、商户后台、门店后台）、线上商店(PC商店、手持设备商店、微信)、

会员中心、开放API四部分组成。


Go2o使用领域驱动设计对业务深度抽象，理论上支持大部分行业的O2O应用场景。通过开放API,可以无缝

集成到旧有系统。


## 项目说明 ##


__项目最新版本: v 0.2__ （因为要求高，v1.0还有距离)

由 #刘铭#, #大鹏# （全职开发,无收入 -_- )，

------------------------

代码保证每周更新2-5次(___因为项目本身，希望你能接受不定时可能较为频繁的代码改动和DB结构调整___ )
同时也欢迎朋友加入我们（___代码可能写的不好，但我们会改进，请大牛不要鄙视，开源需要支持和鼓励___）。

请支持开源，不做伸手党，不拿来主意！

**链接:** [TODO列表](https://github.com/atnet/go2o/tree/master/docs/dev/TODO.md)

========================================

捐赠支付宝:newmin.net@gmail.com (金额随意)
**如有定制需求可邮件联系< lm#s1n1.com >。**


## 演示地址 ##
********* (目前还不完善，暂不提供)

珠三Go技术群：**102629585** (___珠三是珠三角 , 是QQ群，不是其他群___)


## Deploy ##
### 1. Import database ###
> Create new mysql db instance named "go2o"
 and import data use mysql utility.
 Database backup file is here : [go2o.sql](https://github.com/atnet/go2o/blob/master/docs/data/go2o.sql)

### 2.Complied ###
	git clone https://github.com/atnet/go2o.git /home/usr/go/src/go2o
	export GOPATH=$GOPATH:/home/usr/go/
	cd /home/usr/go
	go build server.go

### 2.Running Service ###
	Usage of ./server:
		-debug=false: enable debug
		-help=false: command usage
		-mode="sha": boot mode.'h'- boot http service,'s'- boot socket service
	    -port=1001: web server port
		-port2=1002: socket server port

### 3.Add http proxy by nginx ###
	server {
		listen          80;
		server_name     *.ts.com;
		location / {
			proxy_pass   http://localhost:1002;
			proxy_set_header Host $host;
		}
	}



### 4.Add test hosts ###
> echo   127.0.0.1    go2o.ts.com static.ts.com img.ts.com partner.ts.com member.ts.com www.ts1.com www.ts2.com api.ts.com webmaster.ts.com >> /etc/hosts

## Access Entry ##

### WebMaster ##
webmaster.ts.com

account: go2o / 123456

### Partner Management ###
partner.ts.com

account: go2o / 123456

### Member Center ###
member.ts.com

### Partner Sales ###
go2o.ts.com

you can add host to table "pt_host" use MySql Workbench.

