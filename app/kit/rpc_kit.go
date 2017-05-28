package kit

import (
	"go2o/core/service/thrift"
	"go2o/core/service/thrift/idl/gen-go/define"
)

type rpcToolkit struct {
}

func (r *rpcToolkit) Registry(keys ...string) map[string]string {
	cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer cli.Transport.Close()
		r, _ := cli.GetRegistryV1(keys)
		return r
	}
	return map[string]string{}
}

func (r *rpcToolkit) GetComplexMember(memberId int64) *define.ComplexMember {
	cli, err := thrift.MemberServeClient()
	if err == nil {
		defer cli.Transport.Close()
		mc, _ := cli.Complex(memberId)
		return mc
	}
	return nil
}

func (r *rpcToolkit) InviterArray(memberId int64, depth int32) []int64 {
	cli, err := thrift.MemberServeClient()
	if err == nil {
		defer cli.Transport.Close()
		mc, _ := cli.InviterArray(memberId, depth)
		return mc
	}
	return nil
}

func (r *rpcToolkit) GetMerchant(mchId int32) *define.ComplexMerchant {
	cli, err := thrift.MerchantServeClient()
	if err == nil {
		defer cli.Transport.Close()
		mc, _ := cli.Complex(mchId)
		return mc
	}
	return nil
}

func (r *rpcToolkit) GetLevel(levelId int32) *define.Level {
	cli, err := thrift.MemberServeClient()
	if err == nil {
		defer cli.Transport.Close()
		mc, _ := cli.GetLevel(levelId)
		return mc
	}
	return nil
}

// 获取订单
func (r *rpcToolkit) GetOrder(orderNo string, sub bool) *define.ComplexOrder {
	cli, err := thrift.SaleServeClient()
	if err == nil {
		defer cli.Transport.Close()
		o, _ := cli.GetOrder(orderNo, sub)
		return o
	}
	return nil
}

// 获取订单和商品项信息
func (r *rpcToolkit) GetOrderAndItems(orderNo string, sub bool) *define.ComplexOrder {
	cli, err := thrift.SaleServeClient()
	if err == nil {
		defer cli.Transport.Close()
		o, _ := cli.GetOrderAndItems(orderNo, sub)
		return o
	}
	return nil
}
