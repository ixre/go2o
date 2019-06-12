/**
 * Copyright 2015 @ at3.net.
 * name : server.go
 * author : jarryliu
 * date : 2016-11-12 18:52
 * description :
 * history :
 */
package thrift

import (
	"crypto/tls"
	"github.com/apache/thrift/lib/go/thrift"
	"go2o/core/service/auto_gen/rpc/finance_service"
	"go2o/core/service/auto_gen/rpc/foundation_service"
	"go2o/core/service/auto_gen/rpc/item_service"
	"go2o/core/service/auto_gen/rpc/mch_service"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/auto_gen/rpc/order_service"
	"go2o/core/service/auto_gen/rpc/payment_service"
	"go2o/core/service/auto_gen/rpc/shipment_service"
	"go2o/core/service/auto_gen/rpc/shop_service"
	"go2o/core/service/auto_gen/rpc/status_service"
	"go2o/core/service/auto_gen/rpc/wallet_service"
	"go2o/core/service/rsi"
	"log"
)

// 运行Thrift服务
func ListenAndServe(addr string, secure bool) error {
	var err error
	var transport thrift.TServerTransport
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(addr)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	if err == nil {
		processor := thrift.NewTMultiplexedProcessor()
		processor.RegisterProcessor("status", status_service.NewStatusServiceProcessor(rsi.StatusService))
		processor.RegisterProcessor("merchant", mch_service.NewMerchantServiceProcessor(rsi.MerchantService))
		processor.RegisterProcessor("member", member_service.NewMemberServiceProcessor(rsi.MemberService))
		processor.RegisterProcessor("foundation", foundation_service.NewFoundationServiceProcessor(rsi.FoundationService))
		processor.RegisterProcessor("payment", payment_service.NewPaymentServiceProcessor(rsi.PaymentService))
		processor.RegisterProcessor("order", order_service.NewOrderServiceProcessor(rsi.ShoppingService))
		processor.RegisterProcessor("shipment", shipment_service.NewShipmentServiceProcessor(rsi.ShipmentService))
		processor.RegisterProcessor("item", item_service.NewItemServiceProcessor(rsi.ItemService))
		processor.RegisterProcessor("shop", shop_service.NewShopServiceProcessor(rsi.ShopService))
		processor.RegisterProcessor("finance", finance_service.NewFinanceServiceProcessor(rsi.PersonFinanceService))
		processor.RegisterProcessor("wallet", wallet_service.NewWalletServiceProcessor(rsi.WalletService))
		server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
		log.Println("** [ Go2o][ RPC]: Starting thrift server on port ", addr)
		err = server.Serve()
	}
	return err
}
