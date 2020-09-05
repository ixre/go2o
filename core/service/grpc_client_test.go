package service

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"go2o/core/etcd"
	"go2o/core/service/proto"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : grpc_client_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-05 02:29
 * description :
 * history :
 */

var service = "Go2oService"
var ttl int64 = 3
var cfg = clientv3.Config{
	Endpoints:   []string{"http://localhost:2379/"},
	DialTimeout: 5 * time.Second,
	Username: "",
	Password: "",
}
func TestSelector(t *testing.T) {
	selector,_ := etcd.NewSelector(service,cfg,etcd.AlgRoundRobin)
	for {
		next, err := selector.Next()
		if err != nil{
			t.Logf("select node error %s",err.Error())
		}else {
			t.Log("selected:" + next.Addr)
			requestRPC(next.Addr)
		}
		time.Sleep(time.Second*5)
	}
}

func requestRPC(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewGreeterServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Hello(ctx, &proto.User1{Name: "jarrysix"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Name)
}
