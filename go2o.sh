go build go2o-server.go
go build go2o-daemon.go
go build go2o-tcpserve.go

./go2o-server -conf=app.conf -d &
./go2o-tcpserve.go -conf=app.conf &