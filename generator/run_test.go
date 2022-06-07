/**
 * Copyright 2015 @ at3.net.
 * name : orm_test
 * author : jarryliu
 * date : 2016-11-11 15:26
 * description :
 * history :
 */
package tool

import (
	"github.com/ixre/gof/shell"
	"testing"
)

// 按模板生成数据库所有的代码文件
func TestGenAll(t *testing.T) {
	tablePrefix := "sale_order"
	_, _, err := shell.Run("bash tto.sh "+tablePrefix, true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
