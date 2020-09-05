/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 15:37
 * description :
 * history :
 */
package delivery

type IDeliveryRepo interface {
	// 获取配送
	GetDelivery(id int32) IDelivery

	// 根据区名获取区域
	GetAreaByArea(name string) []*AreaValue

	// 保存覆盖区域
	SaveCoverageArea(*CoverageValue) (int32, error)

	// 获取覆盖区域
	GetCoverageArea(areaId, id int32) *CoverageValue

	// 获取所有的覆盖区域
	GetAllCoverageAreas(areaId int32) []*CoverageValue

	// 获取配送绑定
	GetDeliveryBind(mchId, coverageId int32) *MerchantDeliverBind
}
