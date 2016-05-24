cd proto/
protoc master_*.proto --go_out=plugins=grpc:master
