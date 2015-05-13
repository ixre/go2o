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
    "go2o/src/core/infrastructure/domain"
    "fmt"
)

func main(){

    fmt.Println(domain.EncodePartnerPwd("wzo2o","12345"))

    fmt.Println(domain.NewApiId(105))
}
