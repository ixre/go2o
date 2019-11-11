/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-01-09 21:45
 * description :
 * history :
 */

package enum

const (
	ORDER_LOG_SETUP        OrderLogType = 1
	ORDER_LOG_CHANGE_PRICE OrderLogType = 2
)

type OrderLogType int

func (this OrderLogType) String() string {
	switch this {
	case ORDER_LOG_SETUP:
		return "流程"
	case ORDER_LOG_CHANGE_PRICE:
		return "调价"
	}
	return ""
}
