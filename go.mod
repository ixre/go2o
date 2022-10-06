module github.com/ixre/go2o

go 1.18

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
	github.com/coreos/etcd v3.3.27+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.3.0 // indirect
	github.com/ixre/alidayu v0.0.0-20160128071321-7eadea36c79c
	github.com/ixre/gof v1.13.0
	github.com/labstack/echo/v4 v4.6.1
	github.com/lib/pq v1.10.4 // indirect
	github.com/nats-io/nats-server/v2 v2.2.0 // indirect
	github.com/nats-io/nats.go v1.13.0
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/smartwalle/resize v1.0.0
	go.etcd.io/etcd v3.3.27+incompatible
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20211123203042-d83791d6bcd9
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/text v0.3.7
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	google.golang.org/genproto v0.0.0-20211016002631-37fc39342514 // indirect
	google.golang.org/grpc v1.41.0
)

require google.golang.org/protobuf v1.27.1 // indirect

require github.com/ClickHouse/clickhouse-go/v2 v2.0.12

require (
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20220709002704-04386ae12ed0 // indirect
	github.com/denisenkom/go-mssqldb v0.12.2 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/jonboulle/clockwork v0.3.0 // indirect
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/paulmach/orb v0.4.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.opentelemetry.io/otel v1.4.1 // indirect
	go.opentelemetry.io/otel/trace v1.4.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
