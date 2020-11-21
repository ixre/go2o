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
	"strings"
	"testing"
)


// 按模板生成数据库所有的代码文件
func TestGenAll(t *testing.T) {
	_, output, err := shell.Run("tto -conf tto.conf")
	if err != nil{
		if strings.Index(err.Error(),"not found") != -1{
			t.Log("未安装tto客户端,下载地址：https://github.com/ixre/tto/releases/")
		}else {
			t.Error(err)
		}
		t.FailNow()
	}
	t.Log(output)
}

