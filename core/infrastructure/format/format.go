/**
 * Copyright 2015 @ 56x.net.
 * name : format
 * author : jarryliu
 * date : 2016-05-23 19:42
 * description :
 * history :
 */
package format

import (
	"regexp"
	"strings"
)

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

// 屏蔽手机号
func MaskPhone(phone string) string {
	l := len(strings.TrimSpace(phone))
	// 手机号
	if l == 11 {
		return phone[:3] + strings.Repeat("*", 4) + phone[7:]
	}
	if l > 1 {
		return phone[1:] + strings.Repeat("*", l-1)
	}
	return "****"
}

// 屏蔽昵称
func MaskNickname(nickname string) string {
	nickname = strings.TrimPrefix(nickname, "USER")
	l := len(nickname)
	// 手机号
	if l == 11 {
		return nickname[:3] + strings.Repeat("*", 4) + nickname[7:]
	}
	if l > 1 {
		return nickname[1:] + strings.Repeat("*", l-1)
	}
	return "用*"
}

var _cssElementRegexp = regexp.MustCompile(`(line-height|margin|text-indent)[^;]+;*`)

// 移除HTML样式
func RemoveHtmlStyle(html string) string {
	html = strings.ReplaceAll(html, "<style>", "")
	html = strings.ReplaceAll(html, "</style>", "")
	html = strings.ReplaceAll(html, "<script>", "")
	html = strings.ReplaceAll(html, "</script>", "")
	html = _cssElementRegexp.ReplaceAllString(html, "")
	return html
}
