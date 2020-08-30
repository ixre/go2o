package micro

import (
	"context"
	"fmt"
	"go2o/core/service/proto"
	"testing"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reseved.
 *
 * name : micro_client_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-08-30 17:04
 * description :
 * history :
 */

func  TestMicroClient(t *testing.T)  {
	// 创建一个新的服务
	//service := micro.NewService(micro.Name("Greeter.Client"))
	// 初始化
	//service.Init()

	// 创建 Greeter 客户端
	greeter := proto.NewGreeterService("Greeter")

	// 远程调用 Greeter 服务的 Hello 方法
	ret, err := greeter.Hello(context.TODO(), &proto.User{
		Name:    "jarry",
		GroupId: 0,
		Extra:   map[string]string{},
	})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", ret)
}
