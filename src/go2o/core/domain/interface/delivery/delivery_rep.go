/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-14 15:37
 * description :
 * history :
 */
package delivery

type IDeliveryRep interface {
	// 获取配送
	GetDelivery(int) IDelivery

	// 保存覆盖区域
	SaveCoverageArea(*CoverageValue) (int, error)
}
