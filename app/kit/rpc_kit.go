package kit

import (
	"go2o/core/service/auto_gen/rpc/mch_service"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/auto_gen/rpc/order_service"
	"go2o/core/service/auto_gen/rpc/shop_service"
	"go2o/core/service/thrift"
)

// RPC服务
var RPC = NewRpcToolkit()

type RpcToolkit struct {
}

func NewRpcToolkit() *RpcToolkit {
	return &RpcToolkit{}
}

func (r *RpcToolkit) Registry(keys ...string) map[string]string {
	trans, cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.GetRegistries(thrift.Context, keys)
		return r
	}
	return make(map[string]string, 0)
}

func (r *RpcToolkit) RegistryMap(keys ...string) map[string]string {
	trans, cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.GetRegistries(thrift.Context, keys)
		return r
	}
	return map[string]string{}
}

func (r *RpcToolkit) GetComplexMember(memberId int64) *member_service.SComplexMember {
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		mc, _ := cli.Complex(thrift.Context, memberId)
		return mc
	}
	return nil
}

func (r *RpcToolkit) InviterArray(memberId int64, depth int32) []int64 {
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		mc, _ := cli.InviterArray(thrift.Context, memberId, depth)
		return mc
	}
	return nil
}

func (r *RpcToolkit) GetMerchant(mchId int32) *mch_service.SComplexMerchant {
	trans, cli, err := thrift.MerchantServeClient()
	if err == nil {
		defer trans.Close()
		mc, _ := cli.Complex(thrift.Context, mchId)
		return mc
	}
	return nil
}

func (r *RpcToolkit) GetLevel(levelId int32) *member_service.SLevel {
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		mc, _ := cli.GetLevel(thrift.Context, levelId)
		return mc
	}
	return nil
}

// 获取订单
func (r *RpcToolkit) GetOrder(orderNo string, sub bool) *order_service.SComplexOrder {
	trans, cli, err := thrift.OrderServeClient()
	if err == nil {
		defer trans.Close()
		o, _ := cli.GetOrder(thrift.Context, orderNo, sub)
		return o
	}
	return nil
}

// 获取订单和商品项信息
func (r *RpcToolkit) GetOrderAndItems(orderNo string, sub bool) *order_service.SComplexOrder {
	trans, cli, err := thrift.OrderServeClient()
	if err == nil {
		defer trans.Close()
		o, _ := cli.GetOrderAndItems(thrift.Context, orderNo, sub)
		return o
	}
	return nil
}

// 获取店铺
func (r *RpcToolkit) GetStore(vendorId int32) *shop_service.SStore {
	trans, cli, err := thrift.ShopServeClient()
	if err == nil {
		defer trans.Close()
		o, _ := cli.GetStore(thrift.Context, vendorId)
		return o
	}
	return nil
}

// 获取店铺
func (r *RpcToolkit) GetStoreById(shopId int32) *shop_service.SStore {
	trans, cli, err := thrift.ShopServeClient()
	if err == nil {
		defer trans.Close()
		o, _ := cli.GetStoreById(thrift.Context, shopId)
		return o
	}
	return nil
}
