module go2o

go 1.15

replace github.com/ixre/gof => ../github.com/ixre/gof

replace github.com/ixre/tto => ../github.com/ixre/tto

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

//replace github.com/lucas-clemente/quic-go => github.com/lucas-clemente/quic-go v0.14.1

replace github.com/gomodule/redigo/redis => github.com/gomodule/redigo v1.8.2

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

//exclude github.com/coreos/etcd v3.3.18+incompatible

//exclude github.com/coreos/etcd v3.3.19+incompatible

require (
	github.com/apache/thrift v0.13.0
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo v1.8.2
	github.com/google/uuid v1.1.2 // indirect
	github.com/ixre/alidayu v0.0.0-20160128071321-7eadea36c79c
	github.com/ixre/gof v0.0.0-00010101000000-000000000000
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/smartwalle/resize v1.0.0
	github.com/valyala/fasttemplate v1.2.1 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/text v0.3.3
	google.golang.org/grpc v1.23.0
)
