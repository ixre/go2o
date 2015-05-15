/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-14 15:37
 * description :
 * history :
 */
package delivery

type IDeliveryRep interface {
	// 获取配送
	GetDelivery(int) IDelivery

	// 根据区名获取区域
	GetAreaByArea(name string) []*AreaValue

	// 保存覆盖区域
	SaveCoverageArea(*CoverageValue) (int, error)

	// 获取覆盖区域
	GetCoverageArea(areaId, id int) *CoverageValue

	// 获取所有的覆盖区域
	GetAllCoverageAreas(areaId int) []*CoverageValue

	// 获取配送绑定
	GetDeliveryBind(partnerId, coverageId int) *PartnerDeliverBind
}
