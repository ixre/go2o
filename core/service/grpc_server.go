package service

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
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

// 启动RPC服务
func ServeRPC(ch chan bool, cfg *clientv3.Config, port int) {
	log.Println("[ GO2O][ INFO]: start grpc server")
	// 启动RPC服务
	s := grpc.NewServer()
	proto.RegisterStatusServiceServer(s, inject.GetStatusService())
	proto.RegisterRegistryServiceServer(s, inject.GetRegistryService())
	proto.RegisterCheckServiceServer(s, inject.GetCheckService())
	proto.RegisterMerchantServiceServer(s, inject.GetMerchantService())
	proto.RegisterMemberServiceServer(s, inject.GetMemberService())
	proto.RegisterSystemServiceServer(s, inject.GetSystemService())
	proto.RegisterMessageServiceServer(s, inject.GetMessageService())
	proto.RegisterContentServiceServer(s, inject.GetContentService())
	proto.RegisterPaymentServiceServer(s, inject.GetPaymentService())
	proto.RegisterWalletServiceServer(s, inject.GetWalletService())
	proto.RegisterCartServiceServer(s, inject.GetCartService())
	proto.RegisterOrderServiceServer(s, inject.GetOrderService())

	proto.RegisterShopServiceServer(s, inject.GetShopService())
	proto.RegisterShipmentServiceServer(s, inject.GetShipmentService())
	proto.RegisterItemServiceServer(s, inject.GetItemService())
	proto.RegisterFinanceServiceServer(s, inject.GetPersonFinanceService())
	proto.RegisterQueryServiceServer(s, inject.GetQueryService())
	proto.RegisterProductServiceServer(s, inject.GetProductService())

	proto.RegisterAfterSalesServiceServer(s, inject.GetAfterSalesService())
	proto.RegisterExpressServiceServer(s, inject.GetExpressService())
	proto.RegisterAdvertisementServiceServer(s, inject.GetAdService())
	proto.RegisterPortalServiceServer(s, inject.GetPortalService())
	proto.RegisterExecutionServiceServer(s, inject.GetExecuteService())

	// standalone service
	proto.RegisterQuickPayServiceServer(s, inject.GetQuickPayService())
	proto.RegisterAppServiceServer(s, inject.GetAppService())
	proto.RegisterRbacServiceServer(s, inject.GetRbacService())
	go serveRPC(ch, s, port)
}

func serveRPC(ch chan bool, s *grpc.Server, port int) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	log.Println("[ GO2O][ INFO]: grpc node serve on port :" + strconv.Itoa(port))
	if err = s.Serve(l); err != nil {
		ch <- false
		panic(err)
	}
}
