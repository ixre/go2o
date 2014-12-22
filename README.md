 Go2o 
================
# What's Go2o  #
Golang combine simple o2o DDD domain-driven design realization, including multi-channel (businesses), multi-store, multi-member commodity,Promotions, orders, coupons implementation also includes a mini-framework in package "ops/cf", providing ORM, Reporting, Web Framework,Rpc Framework.

# Deploy #

## 1.Complied ##
  git clone https://github.com/newmin/go2o.git /home/usr/go2o
  export GOPATH=$GOPATH:/home/usr/go2o
  cd /home/usr/go2o
  go build server.go

## 2.Running Service ##
  Usage of ./server:
    -debug=false: enable debug
    -help=false: command usage
    -mode="sh": boot mode.'h'- boot http service,'s'- boot socket service
    -port=1001: web server port
    -port2=1002: socket server port

## 3.Add http proxy for nginx ##

## 4.Add test hosts ##
  vi /etc/hosts
  127.0.0.1   wly.ts.com static.ts.com img.ts.com partner.ts.com
              member.ts.com www.ts1.com www.ts2.com api.ts.com
              wsapi.ts.com

