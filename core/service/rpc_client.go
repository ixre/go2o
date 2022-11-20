/**
 * Copyright 2015 @ at3.net.
 * name : restful_client.go
 * author : jarryliu
 * date : 2016-11-13 12:32
 * description :
 * history :
 */
package service

import (
	"context"
	"os"
	"time"

	"github.com/ixre/go2o/core/etcd"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/log"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
)

var addr string
var selector etcd.Selector

// ConfigureClient 设置RPC地址,defaultAddr为默认地址,当未指定clientv3.Config时使用
func ConfigureClient(c *clientv3.Config, defaultAddr string) {
	addr = defaultAddr
	if c != nil {
		log.Println("[ Go2o][ INFO]: connecting go2o rpc server...")
		s, err := etcd.NewSelector(service, *c, etcd.AlgRoundRobin)
		if err != nil {
			log.Println("[ Go2o][ ERROR]: can't connect go2o rpc server! ", err.Error())
			os.Exit(1)
		}
		selector = s
		tryConnect(30)
	}
}

// 尝试连接服务,如果连接不成功,则退出
func tryConnect(retryTimes int) {
	for i := 0; i < retryTimes; i++ {
		trans, _, err := StatusServiceClient()
		if err == nil {
			trans.Close()
			break
		}
		time.Sleep(time.Second)
		if i >= retryTimes-1 {
			log.Println("[ Go2o][ Fatal]: Can not connect go2o rpc server")
			os.Exit(1)
		}
	}
}

// 获取连接
func getConn(selector etcd.Selector) (*grpc.ClientConn, error) {
	addr := addr
	if selector != nil {
		next, err := selector.Next()
		if err != nil {
			return nil, err
		}
		addr = next.Addr
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	cancel()
	return conn, err
}

// StatusServiceClient 状态客户端
func StatusServiceClient() (*grpc.ClientConn, proto.StatusServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewStatusServiceClient(conn), err
	}
	return conn, nil, err
}

// RegistryServiceClient 基础服务
func RegistryServiceClient() (*grpc.ClientConn, proto.RegistryServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewRegistryServiceClient(conn), err
	}
	return conn, nil, err
}

// MerchantServiceClient 商户客户端
func MerchantServiceClient() (*grpc.ClientConn, proto.MerchantServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewMerchantServiceClient(conn), err
	}
	return conn, nil, err
}

// MemberServiceClient 会员客户端
func MemberServiceClient() (*grpc.ClientConn, proto.MemberServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewMemberServiceClient(conn), err
	}
	return conn, nil, err
}

// FoundationServiceClient 基础服务
func FoundationServiceClient() (*grpc.ClientConn, proto.FoundationServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewFoundationServiceClient(conn), err
	}
	return conn, nil, err
}

// MessageServiceClient 消息客户端
func MessageServiceClient() (*grpc.ClientConn, proto.MessageServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewMessageServiceClient(conn), err
	}
	return conn, nil, err
}

// ContentServiceClient 消息客户端
func ContentServiceClient() (*grpc.ClientConn, proto.ContentServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewContentServiceClient(conn), err
	}
	return conn, nil, err
}

// PaymentServiceClient 支付服务
func PaymentServiceClient() (*grpc.ClientConn, proto.PaymentServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewPaymentServiceClient(conn), err
	}
	return conn, nil, err
}

// QuickPaymentServiceClient 快捷支付服务
func QuickPaymentServiceClient() (*grpc.ClientConn, proto.QuickPayServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewQuickPayServiceClient(conn), err
	}
	return conn, nil, err
}

// WalletClient 钱包服务
func WalletClient() (*grpc.ClientConn, proto.WalletServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewWalletServiceClient(conn), err
	}
	return conn, nil, err
}

// OrderServiceClient 订单服务
func OrderServiceClient() (*grpc.ClientConn, proto.OrderServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewOrderServiceClient(conn), err
	}
	return conn, nil, err
}

// CartServiceClient 购物车服务
func CartServiceClient() (*grpc.ClientConn, proto.CartServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewCartServiceClient(conn), err
	}
	return conn, nil, err
}

// ExpressServiceClient 快递服务
func ExpressServiceClient() (*grpc.ClientConn, proto.ExpressServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewExpressServiceClient(conn), err
	}
	return conn, nil, err
}

// ShipmentServiceClient 物流服务
func ShipmentServiceClient() (*grpc.ClientConn, proto.ShipmentServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewShipmentServiceClient(conn), err
	}
	return conn, nil, err
}

// ItemServiceClient 商品服务
func ItemServiceClient() (*grpc.ClientConn, proto.ItemServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewItemServiceClient(conn), err
	}
	return conn, nil, err
}

// ProductServiceClient 产品服务
func ProductServiceClient() (*grpc.ClientConn, proto.ProductServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewProductServiceClient(conn), err
	}
	return conn, nil, err
}

// ShopServiceClient 商店服务
func ShopServiceClient() (*grpc.ClientConn, proto.ShopServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewShopServiceClient(conn), err
	}
	return conn, nil, err
}

// FinanceServiceClient 财务服务
func FinanceServiceClient() (*grpc.ClientConn, proto.FinanceServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewFinanceServiceClient(conn), err
	}
	return conn, nil, err
}

// QueryServiceClient 查询服务
func QueryServiceClient() (*grpc.ClientConn, proto.QueryServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewQueryServiceClient(conn), err
	}
	return conn, nil, err
}

// AfterSalesServiceClient 售后服务
func AfterSalesServiceClient() (*grpc.ClientConn, proto.AfterSalesServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewAfterSalesServiceClient(conn), err
	}
	return conn, nil, err
}

// AdvertisementServiceClient 广告服务
func AdvertisementServiceClient() (*grpc.ClientConn, proto.AdvertisementServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewAdvertisementServiceClient(conn), err
	}
	return conn, nil, err
}

// PortalServiceClient 门户服务
func PortalServiceClient() (*grpc.ClientConn, proto.PortalServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewPortalServiceClient(conn), err
	}
	return conn, nil, err
}

// ExecutionServiceClient 任务执行服务
func ExecutionServiceClient() (*grpc.ClientConn, proto.ExecutionServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewExecutionServiceClient(conn), err
	}
	return conn, nil, err
}

// AppServiceClient APP服务
func AppServiceClient() (*grpc.ClientConn, proto.AppServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewAppServiceClient(conn), err
	}
	return conn, nil, err
}

// RbacServiceClient RBAC服务
func RbacServiceClient() (*grpc.ClientConn, proto.RbacServiceClient, error) {
	conn, err := getConn(selector)
	if err == nil {
		return conn, proto.NewRbacServiceClient(conn), err
	}
	return conn, nil, err
}
