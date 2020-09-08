/**
 * Copyright 2015 @ at3.net.
 * name : client.go
 * author : jarryliu
 * date : 2016-11-13 12:32
 * description :
 * history :
 */
package  service

import (
	"go.etcd.io/etcd/clientv3"
	"go2o/core/etcd"
	"go2o/core/service/proto"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)


func init() {
}

//var cfg  clientv3.Config
var selector etcd.Selector

// 设置Thrift地址
func ConfigureClient(c clientv3.Config){
	//cfg = c
	log.Println("[ Go2o][ RPC]: connecting go2o rpc server...")
	s,err := etcd.NewSelector(service,c,etcd.AlgRoundRobin)
	if err != nil{
		log.Println("[ Go2o][ RPC]: can't connect go2o rpc server! ", err.Error())
		os.Exit(1)
	}
	selector = s
	tryConnect(10)
}

// 尝试连接服务,如果连接不成功,则退出
func tryConnect(retryTimes int) {
	for i := 0; i < retryTimes; i++ {
		trans, _, err := StatusServeClient()
		if err == nil {
			trans.Close()
			break
		} else if i == retryTimes-1 {
			log.Println("[ Go2o][ RPC]: can't connect go2o rpc server! ", err.Error())
			os.Exit(1)
		}
		time.Sleep(time.Second * 2)
	}
}

// 获取连接
func getConn(selector etcd.Selector) (*grpc.ClientConn, error) {
	next, err := selector.Next()
	if err != nil{
		log.Printf("[ go2o][ rpc]: select node error %s\n",err.Error())
		return nil,err
	}
	return grpc.Dial(next.Addr, grpc.WithInsecure(), grpc.WithBlock())
}

// 状态客户端
func StatusServeClient() (*grpc.ClientConn, proto.StatusServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewStatusServiceClient(conn),err
	}
	return conn, nil, err
}

// 基础服务
func RegistryServeClient() (*grpc.ClientConn, proto.RegistryServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewRegistryServiceClient(conn),err
	}
	return conn, nil, err
}

// 商户客户端
func MerchantServeClient() (*grpc.ClientConn, proto.MerchantServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewMerchantServiceClient(conn),err
	}
	return conn, nil, err
}

// 会员客户端
func MemberServeClient() (*grpc.ClientConn, proto.MemberServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewMemberServiceClient(conn),err
	}
	return conn, nil, err
}

// 基础服务
func FoundationServeClient() (*grpc.ClientConn, proto.FoundationServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewFoundationServiceClient(conn),err
	}
	return conn, nil, err
}

// 消息客户端
func MessageServeClient() (*grpc.ClientConn, proto.MessageServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewMessageServiceClient(conn),err
	}
	return conn, nil, err
}

// 消息客户端
func ContentServeClient() (*grpc.ClientConn, proto.ContentServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewContentServiceClient(conn),err
	}
	return conn, nil, err
}

// 支付服务
func PaymentServeClient() (*grpc.ClientConn, proto.PaymentServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewPaymentServiceClient(conn),err
	}
	return conn, nil, err
}

// 钱包服务
func WalletClient() (*grpc.ClientConn, proto.WalletServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewWalletServiceClient(conn),err
	}
	return conn, nil, err
}

// 订单服务
func OrderServeClient() (*grpc.ClientConn, proto.OrderServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewOrderServiceClient(conn),err
	}
	return conn, nil, err
}

// 物流服务
func ShipmentServeClient() (*grpc.ClientConn, proto.ShipmentServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewShipmentServiceClient(conn),err
	}
	return conn, nil, err
}

// 商品服务
func ItemServeClient() (*grpc.ClientConn, proto.ItemServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewItemServiceClient(conn),err
	}
	return conn, nil, err
}

// 商店服务
func ShopServeClient() (*grpc.ClientConn, proto.ShopServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewShopServiceClient(conn),err
	}
	return conn, nil, err
}

// 财务服务
func FinanceServeClient() (*grpc.ClientConn, proto.FinanceServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil{
		return conn,proto.NewFinanceServiceClient(conn),err
	}
	return conn, nil, err
}
