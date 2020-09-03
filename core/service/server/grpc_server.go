package server

import (
	grpc2 "go2o/core/service/grpc"
	"go2o/core/service/proto"
	"google.golang.org/grpc"
	"net"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : grpc_server
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-03 17:03
 * description :
 * history :
 */

func ServeRPC(ch chan bool, addr string){
	s := grpc.NewServer()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	proto.RegisterGreeterServiceServer(s, &grpc2.TestServiceImpl{})
	if err = s.Serve(l); err != nil {
		ch <- false
		panic(err)
	}
}
