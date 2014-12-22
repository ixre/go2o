/**
 * Copyright 2014 @ ops.
 * name :
 * author : newmin
 * date : 2013-11-10 16:10
 * description :
 * history :
 */

package ording

import (
	"bytes"
	"regexp"
)

var (
	cartFmtRegex = regexp.MustCompile("[^\\*]+\\*(\\d+\\*\\d+)\\*\\d")
)

//传唤客户端的购物车cookie为服务端的方式.
func CartCookieFmt(s string) string {
	//cart=%u91CE%u5C71%u6912%u7092%u8089*5*1*2|2
	matches := cartFmtRegex.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return ""
	}
	b := bytes.NewBufferString("")
	for i, k := range matches {
		if i != 0 {
			b.WriteString("|")
		}
		b.WriteString(k[1])
	}
	return b.String()
}
