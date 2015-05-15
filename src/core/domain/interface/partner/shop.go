/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-11-22 20:01
 * description :
 * history :
 */

package partner

type IShop interface {
	GetDomainId() int

	GetValue() ValueShop

	SetValue(*ValueShop) error

	//	// 获取经维度
	//	GetLngLat() (float64, float64)
	//
	//	// 是否可以配送
	//	// 返回是否可以配送，以及距离(米)
	//	CanDeliver(lng, lat float64) (bool, int)
	//
	//	// 是否可以配送
	//	// 返回是否可以配送，以及距离(米)
	//	CanDeliverTo(address string) (bool, int)

	Save() (int, error)
}
