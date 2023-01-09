/**
 * Copyright 2015 @ 56x.net.
 * name : format
 * author : jarryliu
 * date : 2016-05-23 19:42
 * description :
 * history :
 */
package format

import "strings"

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

// 屏蔽昵称
func MaskNickname(nickname string) string {
	if strings.HasPrefix(nickname,"USER"){
		nickname = nickname[4:]
	}

	l :=len(nickname)
	// 手机号
	if l == 11{
		return nickname[:3]+strings.Repeat("*",4)+nickname[8:]
	}
	if l > 1{
		return nickname[1:]+strings.Repeat("*",l-1)
	}
	return "用*"
}