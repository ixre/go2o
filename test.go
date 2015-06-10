/**
 * Copyright 2015 @ S1N1 Team.
 * name : test.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package main

import (
	"fmt"
	"github.com/atnet/gof/util"
	"go2o/src/core/infrastructure/domain"
)

func main() {

	fmt.Println(domain.Md5PartnerPwd("wzo2o", "12345"))

	fmt.Println(domain.NewApiId(105))
	fmt.Println(domain.Md5MemberPwd("u1000", "123456"))
	fmt.Println(1 << 2)
	fmt.Println(util.IsMobileAgent("Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"))
}
