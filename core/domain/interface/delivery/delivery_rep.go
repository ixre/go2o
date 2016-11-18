/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 15:37
 * description :
 * history :
 */
package delivery

type IDeliveryRep interface {
	// 获取配送
	GetDelivery(id int64) IDelivery

	// 根据区名获取区域
	GetAreaByArea(name string) []*AreaValue

	// 保存覆盖区域
	SaveCoverageArea(*CoverageValue) (int64, error)

	// 获取覆盖区域
	GetCoverageArea(areaId, id int64) *CoverageValue

	// 获取所有的覆盖区域
	GetAllCoverageAreas(areaId int64) []*CoverageValue

	// 获取配送绑定
	GetDeliveryBind(mchId, coverageId int64) *MerchantDeliverBind
}
