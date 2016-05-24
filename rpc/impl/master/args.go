/**
 * Copyright 2015 @ z3q.net.
 * name : args.go
 * author : jarryliu
 * date : 2016-05-21 13:12
 * description :
 * history :
 */
package impl

import (
	"errors"
)

const (
	token string = "go2o-master-comm-key"
)

type Args map[string]interface{}

// 检查参数是否缺失
func checkArgs(args *Args, p ...string) error {
	for _, k := range p {
		if _, ok := args[k]; !ok {
			return errors.New("miss params '" + k + "'")
		}
	}
	return nil
}
