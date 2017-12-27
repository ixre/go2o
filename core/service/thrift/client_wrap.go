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
	"git.apache.org/thrift.git/lib/go/thrift"
	"go2o/gen-code/thrift/define"
)

// 商户客户端
func MerchantServeClient() (thrift.TTransport, *define.MerchantServiceClient, error) {
	trans, cli, err := getClient("merchant")
	return trans, define.NewMerchantServiceClient(cli), err
}

// 会员客户端
func MemberServeClient() (thrift.TTransport, *define.MemberServiceClient, error) {
	trans, cli, err := getClient("member")
	return trans, define.NewMemberServiceClient(cli), err
}

// 基础服务
func FoundationServeClient() (thrift.TTransport, *define.FoundationServiceClient, error) {
	trans, cli, err := getClient("foundation")
	return trans, define.NewFoundationServiceClient(cli), err
}

// 基础服务
func PaymentServeClient() (thrift.TTransport, *define.PaymentServiceClient, error) {
	trans, cli, err := getClient("payment")
	return trans, define.NewPaymentServiceClient(cli), err
}

// 基础服务
func WalletClient() (thrift.TTransport, *define.WalletServiceClient, error) {
	trans, cli, err := getClient("wallet")
	return trans, define.NewWalletServiceClient(cli), err
}

// 订单服务
func OrderServeClient() (thrift.TTransport, *define.OrderServiceClient, error) {
	trans, cli, err := getClient("order")
	return trans, define.NewOrderServiceClient(cli), err
}

// 基础服务
func ShipmentServeClient() (thrift.TTransport, *define.ShipmentServiceClient, error) {
	trans, cli, err := getClient("shipment")
	return trans, define.NewShipmentServiceClient(cli), err
}

// 商品服务
func ItemServeClient() (thrift.TTransport, *define.ItemServiceClient, error) {
	trans, cli, err := getClient("item")
	return trans, define.NewItemServiceClient(cli), err
}

// 商店服务
func ShopServeClient() (thrift.TTransport, *define.ShopServiceClient, error) {
	trans, cli, err := getClient("shop")
	return trans, define.NewShopServiceClient(cli), err
}

// 商店服务
func FinanceServeClient() (thrift.TTransport, *define.FinanceServiceClient, error) {
	trans, cli, err := getClient("finance")
	return trans, define.NewFinanceServiceClient(cli), err
}
