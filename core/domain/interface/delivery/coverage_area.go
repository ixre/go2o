/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-12 16:38
 * description :
 * history :
 */
package delivery

type ICoverageArea interface {
	GetDomainId() int32

	GetValue() CoverageValue

	SetValue(*CoverageValue) error

	// 是否可以配送
	// 返回是否可以配送，以及距离(米)
	CanDeliver(lng, lat float64) (bool, int)

	// 是否可以配送
	// 返回是否可以配送，以及距离(米)
	CanDeliverTo(address string) (bool, int)

	Save() (int32, error)
}
