/**
 * Copyright 2014 @ ops.
 * name :
 * author : newmin
 * date : 2013-11-10 22:10
 * description :
 * history :
 */

package domain

import (
	"bytes"
	"fmt"
	"github.com/atnet/gof/crypto"
	"regexp"
)

var (
	cartFmtRegex = regexp.MustCompile("[^\\*]+\\*(\\d+\\*\\d+)\\*\\d")
)

//todo: will removed

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

// 生成16位唯一的md5购物车码
func GenerateCartKey(unix int64, nano int) string {
	str := fmt.Sprintf("%d-%d*%d", unix, nano, unix%int64(nano))
	result := crypto.Md5([]byte(str))
	return result[8:24]
}
