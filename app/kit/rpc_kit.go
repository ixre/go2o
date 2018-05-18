package kit

import (
	"go2o/core/service/thrift"
	"go2o/gen-code/thrift/define"
)

type RpcToolkit struct {
}

func NewRpcToolkit() *RpcToolkit {
	return &RpcToolkit{}
}

func (r *RpcToolkit) Registry(keys ...string) []string {
	trans, cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.GetRegistryV1(thrift.Context, keys)
		return r
	}
	return []string{}
}

func (r *RpcToolkit) RegistryMap(keys ...string) map[string]string {
	trans, cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.GetRegistryMapV1(thrift.Context, keys)
		return r
	}
	return map[string]string{}
}

func (r *RpcToolkit) GetComplexMember(memberId int64) *define.ComplexMember {
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

func (r *RpcToolkit) GetMerchant(mchId int32) *define.ComplexMerchant {
	trans, cli, err := thrift.MerchantServeClient()
	if err == nil {
		defer trans.Close()
		mc, _ := cli.Complex(thrift.Context, mchId)
		return mc
	}
	return nil
}

func (r *RpcToolkit) GetLevel(levelId int32) *define.Level {
	trans, cli, err := thrift.MemberServeClient()
	if err == nil {
		defer trans.Close()
		mc, _ := cli.GetLevel(thrift.Context, levelId)
		return mc
	}
	return nil
}

// 获取订单
func (r *RpcToolkit) GetOrder(orderNo string, sub bool) *define.SComplexOrder {
	trans, cli, err := thrift.OrderServeClient()
	if err == nil {
		defer trans.Close()
		o, _ := cli.GetOrder(thrift.Context, orderNo, sub)
		return o
	}
	return nil
}

// 获取订单和商品项信息
func (r *RpcToolkit) GetOrderAndItems(orderNo string, sub bool) *define.SComplexOrder {
	trans, cli, err := thrift.OrderServeClient()
	if err == nil {
		defer trans.Close()
		o, _ := cli.GetOrderAndItems(thrift.Context, orderNo, sub)
		return o
	}
	return nil
}

// 获取店铺
func (r *RpcToolkit) GetStore(vendorId int32) *define.Store {
	trans, cli, err := thrift.ShopServeClient()
	if err == nil {
		defer trans.Close()
		o, _ := cli.GetStore(thrift.Context, vendorId)
		return o
	}
	return nil
}

// 获取店铺
func (r *RpcToolkit) GetStoreById(shopId int32) *define.Store {
	trans, cli, err := thrift.ShopServeClient()
	if err == nil {
		defer trans.Close()
		o, _ := cli.GetStoreById(thrift.Context, shopId)
		return o
	}
	return nil
}
