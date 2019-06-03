/**
 * Copyright 2015 @ to2.net.
 * name : format
 * author : jarryliu
 * date : 2016-05-23 19:42
 * description :
 * history :
 */
package format

// 获取性别
func GetSex(sex int32) string {
	switch sex {
	case 1:
		return "男性"
	case 2:
		return "女性"
	}
	return "-"
}
