module go2o

go 1.12

require (
	github.com/Shopify/sarama v1.22.1
	github.com/afocus/captcha v0.0.0-20190403092343-1e99620393ea
	github.com/apache/thrift v0.12.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/ixre/goex v1.5.1
	github.com/ixre/gof v1.3.2
	github.com/ixre/tto v0.0.0-00010101000000-000000000000
	github.com/jsix/alidayu v0.0.0-20160128071321-7eadea36c79c
	github.com/labstack/echo v3.3.10+incompatible
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/robfig/cron v0.0.0-20180505203441-b41be1df6967
	golang.org/x/image v0.0.0-20190424155947-59b11bec70c7 // indirect
	golang.org/x/text v0.3.2
	gopkg.in/square/go-jose.v1 v1.1.2
)

replace github.com/ixre/gof => ../github.com/ixre/gof

replace github.com/ixre/tto => ../github.com/ixre/tto

replace github.com/ixre/goex => ../github.com/ixre/goex
