module go2o

go 1.15

replace github.com/ixre/gof => ../github.com/ixre/gof

replace github.com/ixre/tto => ../github.com/ixre/tto

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/lucas-clemente/quic-go => github.com/lucas-clemente/quic-go v0.14.1

require (
	github.com/apache/thrift v0.13.0
	github.com/armon/go-metrics v0.3.4 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/go-log/log v0.2.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo v1.8.2
	github.com/google/uuid v1.1.2 // indirect
	github.com/hashicorp/consul/api v1.6.0 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.2.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.9.4 // indirect
	github.com/ixre/alidayu v0.0.0-20160128071321-7eadea36c79c
	github.com/ixre/gof v1.8.8
	github.com/ixre/tto v0.0.0-00010101000000-000000000000
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/lucas-clemente/quic-go v0.18.0 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v3 v3.0.0-beta.0.20200825081046-bf8b3aeac796
	github.com/micro/go-plugins/registry/consul v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/micro/v3 v3.0.0-beta.3
	github.com/miekg/dns v1.1.31 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/nats-io/jwt v1.0.1 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-io/nkeys v0.2.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/smartwalle/resize v1.0.0
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200828194041-157a740278f4 // indirect
	golang.org/x/text v0.3.3
	golang.org/x/tools v0.0.0-20200828161849-5deb26317202 // indirect
	google.golang.org/genproto v0.0.0-20200829155447-2bf3329a0021 // indirect
	google.golang.org/grpc v1.31.1 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
)

//exclude github.com/coreos/etcd v3.3.18+incompatible

//exclude github.com/coreos/etcd v3.3.19+incompatible
