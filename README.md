Go2o
================
# What's Go2o  #
Golang combine simple o2o DDD domain-driven design realization,
including multi-channel (businesses), multi-store, multi-member commodity,
Promotions, orders, coupons implementation also includes a micro framework
in package "ops/cf", providing ORM, Reporting, Web Framework,Rpc Framework.

Golang 结合DDD领域驱动设计的简单o2o实现，包含多渠道(商家),多门店,多会员.商品，
促销，订单，优惠券的实现，同时包含一个微型框架ops/cf,提供ORM,报表,Web Framework,Rpc Framework.

# Deploy #

## 1.Complied ##
  git clone https://github.com/newmin/go2o.git /home/usr/go2o
  export GOPATH=$GOPATH:/home/usr/go2o/
  cd /home/usr/go2o
  go build server.go

## 2.Running Service ##
  Usage of ./server:
    -debug=false: enable debug
    -help=false: command usage
    -mode="sh": boot mode.'h'- boot http service,'s'- boot socket service
    -port=1001: web server port
    -port2=1002: socket server port

## 3.Add http proxy by nginx ##
server {
        listen  80;
        server_name     *.ts.com;

        location / {
           proxy_pass   http://localhost:1001;
           proxy_set_header Host $host;
        }
  }


## 4.Add test hosts ##
  vi /etc/hosts
  127.0.0.1   wly.ts.com static.ts.com img.ts.com partner.ts.com
  member.ts.com www.ts1.com www.ts2.com api.ts.com wsapi.ts.com


# Access Entry #
## Partner Management ##
partner.ts.com

## Member Center ##
member.ts.com

## Partner Sales ##
wly.ts.com

you can add host to table "pt_host" use MySql Workbench.

