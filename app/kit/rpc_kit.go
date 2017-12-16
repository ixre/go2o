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
	cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer cli.Transport.Close()
		r, _ := cli.GetRegistryV1(keys)
		return r
	}
	return []string{}
}

func (r *RpcToolkit) RegistryMap(keys ...string) map[string]string {
	cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer cli.Transport.Close()
		r, _ := cli.GetRegistryMapV1(keys)
		return r
	}
	return map[string]string{}
}

func (r *RpcToolkit) GetComplexMember(memberId int64) *define.ComplexMember {
	cli, err := thrift.MemberServeClient()
	if err == nil {
		defer cli.Transport.Close()
		mc, _ := cli.Complex(memberId)
		return mc
	}
	return nil
}

func (r *RpcToolkit) InviterArray(memberId int64, depth int32) []int64 {
	cli, err := thrift.MemberServeClient()
	if err == nil {
		defer cli.Transport.Close()
		mc, _ := cli.InviterArray(memberId, depth)
		return mc
	}
	return nil
}

func (r *RpcToolkit) GetMerchant(mchId int32) *define.ComplexMerchant {
	cli, err := thrift.MerchantServeClient()
	if err == nil {
		defer cli.Transport.Close()
		mc, _ := cli.Complex(mchId)
		return mc
	}
	return nil
}

func (r *RpcToolkit) GetLevel(levelId int32) *define.Level {
	cli, err := thrift.MemberServeClient()
	if err == nil {
		defer cli.Transport.Close()
		mc, _ := cli.GetLevel(levelId)
		return mc
	}
	return nil
}

// 获取订单
func (r *RpcToolkit) GetOrder(orderNo string, sub bool) *define.ComplexOrder {
	cli, err := thrift.OrderServeClient()
	if err == nil {
		defer cli.Transport.Close()
		o, _ := cli.GetOrder(orderNo, sub)
		return o
	}
	return nil
}

// 获取订单和商品项信息
func (r *RpcToolkit) GetOrderAndItems(orderNo string, sub bool) *define.ComplexOrder {
	cli, err := thrift.OrderServeClient()
	if err == nil {
		defer cli.Transport.Close()
		o, _ := cli.GetOrderAndItems(orderNo, sub)
		return o
	}
	return nil
}

// 获取店铺
func (r *RpcToolkit) GetStore(vendorId int32) *define.Store {
	cli, err := thrift.ShopServeClient()
	if err == nil {
		defer cli.Transport.Close()
		o, _ := cli.GetStore(vendorId)
		return o
	}
	return nil
}

// 获取店铺
func (r *RpcToolkit) GetStoreById(shopId int32) *define.Store {
	cli, err := thrift.ShopServeClient()
	if err == nil {
		defer cli.Transport.Close()
		o, _ := cli.GetStoreById(shopId)
		return o
	}
	return nil
}
