module github.com/ixre/go2o

go 1.18

replace github.com/ixre/gof => ../github.com/ixre/gof

replace github.com/ixre/tto => ../github.com/ixre/tto

// replace google.golang.org/grpc => google.golang.org/grpc v1.28.0

//replace github.com/lucas-clemente/quic-go => github.com/lucas-clemente/quic-go v0.14.1
replace golang.org/x/text => github.com/golang/text v0.3.4

replace github.com/gomodule/redigo/redis => github.com/gomodule/redigo v1.8.2

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.3.0 // indirect
	github.com/ixre/alidayu v0.0.0-20160128071321-7eadea36c79c
	github.com/ixre/gof v1.13.2
	github.com/labstack/echo/v4 v4.6.1
	github.com/lib/pq v1.10.7 // indirect
	github.com/nats-io/nats-server/v2 v2.2.0 // indirect
	github.com/nats-io/nats.go v1.13.0
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/smartwalle/resize v1.0.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.2.0
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	google.golang.org/grpc v1.51.0
)

require google.golang.org/protobuf v1.28.1

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.0.12
	go.etcd.io/etcd/client/v3 v3.5.6
)

require (
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/paulmach/orb v0.4.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	go.etcd.io/etcd/api/v3 v3.5.6 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.6 // indirect
	go.opentelemetry.io/otel v1.4.1 // indirect
	go.opentelemetry.io/otel/trace v1.4.1 // indirect
)
