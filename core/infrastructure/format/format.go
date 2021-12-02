/**
 * Copyright 2015 @ 56x.net.
 * name : format
 * author : jarryliu
 * date : 2016-05-23 19:42
 * description :
 * history :
 */
package format

// 获取性别
func GetGender(gender int32) string {
	switch gender {
	case 1:
		return "男性"
	case 2:
		return "女性"
	}
	return "-"
}
