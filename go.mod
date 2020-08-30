module go2o

go 1.15

replace github.com/ixre/gof => ../github.com/ixre/gof

replace github.com/ixre/tto => ../github.com/ixre/tto

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/apache/thrift v0.13.0
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo v1.8.2
	github.com/ixre/alidayu v0.0.0-20160128071321-7eadea36c79c
	github.com/ixre/gof v1.8.5
	github.com/ixre/tto v0.0.0-00010101000000-000000000000
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v3 v3.0.0-beta
	github.com/micro/go-plugins/registry/consul v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/micro/v3 v3.0.0-beta
	github.com/nats-io/nats.go v1.10.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/smartwalle/resize v1.0.0
	github.com/valyala/fasttemplate v1.2.1 // indirect
	golang.org/x/text v0.3.3
	google.golang.org/protobuf v1.25.0
)

//exclude github.com/coreos/etcd v3.3.18+incompatible

//exclude github.com/coreos/etcd v3.3.19+incompatible
