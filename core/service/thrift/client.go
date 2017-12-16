/**
 * Copyright 2015 @ at3.net.
 * name : client.go
 * author : jarryliu
 * date : 2016-11-13 12:32
 * description :
 * history :
 */
package thrift

import (
	"crypto/tls"
	"git.apache.org/thrift.git/lib/go/thrift"
	"go2o/gen-code/thrift/define"
)

var (
	cliHostPort      = "localhost:14288"
	transportFactory = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory  = thrift.NewTCompactProtocolFactory()
)

// 客户端初始化
func CliInit(hostPort string) {
	cliHostPort = hostPort
}

func getTransportAndProtocol() (thrift.TTransport, thrift.TProtocolFactory, error) {
	var err error
	var serveTransport thrift.TTransport
	secure := false
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return nil, protocolFactory, err
		}
		serveTransport, err = thrift.NewTSSLSocket(cliHostPort, cfg)
	} else {
		serveTransport, err = thrift.NewTSocket(cliHostPort)
	}
	transport := transportFactory.GetTransport(serveTransport)
	return transport, protocolFactory, err
}

// 商户客户端
func MerchantServeClient() (*define.MerchantServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			//多个服务
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "merchant")
			return define.NewMerchantServiceClientProtocol(transport, proto, opProto), err
		}
	}
	return nil, err
}

// 会员客户端
func MemberServeClient() (*define.MemberServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			//多个服务
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "member")
			return define.NewMemberServiceClientProtocol(transport, proto, opProto), err
			//单个服务
			//return define.NewMemberServiceClientFactory(transport, protocol), err
		}
	}
	return nil, err
}

// 基础服务
func FoundationServeClient() (*define.FoundationServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			//多个服务
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "foundation")
			return define.NewFoundationServiceClientProtocol(transport, proto, opProto), err
			//单个服务
			//return define.NewMemberServiceClientFactory(transport, protocol), err
		}
	}
	return nil, err
}

// 基础服务
func PaymentServeClient() (*define.PaymentServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "payment")
			return define.NewPaymentServiceClientProtocol(transport, proto, opProto), err
		}
	}
	return nil, err
}

// 基础服务
func WalletClient() (*define.WalletServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "wallet")
			return define.NewWalletServiceClientProtocol(transport, proto, opProto), err
		}
	}
	return nil, err
}

// 订单服务
func OrderServeClient() (*define.OrderServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "order")
			return define.NewOrderServiceClientProtocol(transport, proto, opProto), err
		}
	}
	return nil, err
}

// 基础服务
func ShipmentServeClient() (*define.ShipmentServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "shipment")
			return define.NewShipmentServiceClientProtocol(transport, proto, opProto), err
		}
	}
	return nil, err
}

// 商品服务
func ItemServeClient() (*define.ItemServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "item")
			return define.NewItemServiceClientProtocol(transport, proto, opProto), err
		}
	}
	return nil, err
}

// 商店服务
func ShopServeClient() (*define.ShopServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "shop")
			return define.NewShopServiceClientProtocol(transport, proto, opProto), err
		}
	}
	return nil, err
}

// 商店服务
func FinanceServeClient() (*define.FinanceServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "finance")
			return define.NewFinanceServiceClientProtocol(transport, proto, opProto), err
		}
	}
	return nil, err
}
