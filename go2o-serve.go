/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-16 21:45
 * description :
 * history :
 */

package main

import (
	"github.com/ixre/go2o/app"
)

var _ = `

	 ####   ####  #######  ####  
	#    # #    #       # #    # 
	#      #    #  #####  #    # 
	#  ### #    # #       #    # 
	#    # #    # #       #    # 
	 ####   ####  #######  #### 


Go2o is Google Go language binding domain-driven design (DDD) O2O open source implementation. Support Online Store
, Offline stores; multi-channel (businesses), multi-store, merchandise, snapshots, orders, sales, payment, distribution and other functions.

Project by a management center (including platform management center, business background, store background), online store (PC shop,
Handheld shops, micro-channel), the member center, open API in four parts.

Go2o using domain-driven design for business depth abstract, theoretical support in most sectors O2O scenarios.
Through open API, you can seamlessly integrate into legacy systems.


Email: jarrysix#gmail.com

`

// - GO2O_SERVER_HOST: 当前节点的主机头或IP,用于指定固定的服务发现IP

func main() {
	app.ParseFlags()
	ch := make(chan bool)
	app.Run(ch, nil)
	<-ch
}

/*
// todo: v3 还是测试版本
func runGoMicro() {
	r := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	grpc.NewServer(
		server.Name("Greeter"),
		server.Registry(NewRegisterV3(r)))
	s := service.New(
		service.Name("Greeter"),
		service.Address(":1081"),
		)
	//service := micro.NewService(
	//	micro.Name("Greeter"),
	//	//micro.Address(":1081"),
	//	micro.Registry(r),
	//	)
	//service.Init()
	s.Handle(new(grpc.TestServiceImpl))
	//proto.RegisterGreeterServiceHandler(service,new(grpc.TestServiceImpl))
	service.Run()
}

func NewRegisterV3(r registry.Registry) registry2.Registry {
	return &RegisterV3{r}
}
*/
