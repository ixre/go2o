module go2o

go 1.16

replace github.com/ixre/gof => ../github.com/ixre/gof

replace github.com/ixre/tto => ../github.com/ixre/tto

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

//replace github.com/lucas-clemente/quic-go => github.com/lucas-clemente/quic-go v0.14.1
replace golang.org/x/text => github.com/golang/text v0.3.4

replace github.com/gomodule/redigo/redis => github.com/gomodule/redigo v1.8.2

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

exclude github.com/coreos/etcd v3.3.18+incompatible

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

require (
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.2.0 // indirect
	github.com/ixre/alidayu v0.0.0-20160128071321-7eadea36c79c
	github.com/ixre/gof v1.10.8
	github.com/labstack/echo/v4 v4.3.0
	github.com/lib/pq v1.10.1 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	github.com/nats-io/jwt v1.2.2 // indirect
	github.com/nats-io/nats-server/v2 v2.1.8 // indirect
	github.com/nats-io/nats.go v1.11.0
	github.com/pelletier/go-toml v1.9.1 // indirect
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/smartwalle/resize v1.0.0
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed
	golang.org/x/sys v0.0.0-20210511113859-b0526f3d8744 // indirect
	golang.org/x/text v0.3.6
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	google.golang.org/genproto v0.0.0-20210510173355-fb37daa5cd7a // indirect
	google.golang.org/grpc v1.37.0
)
