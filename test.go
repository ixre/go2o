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
	"go2o/src/core/infrastructure/domain"
)

func main() {

	fmt.Println(domain.EncodePartnerPwd("wzo2o", "12345"))

	fmt.Println(domain.NewApiId(105))
	fmt.Println(domain.EncodeMemberPwd("u1000", "123456"))

}
