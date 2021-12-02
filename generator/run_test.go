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
	tablePrefix := "comm_"
	_, output, err := shell.Run("tto -m go -conf tto.conf -table " + tablePrefix + " -clean")
	if err != nil {
		if strings.Index(err.Error(), "not found") != -1 {
			t.Log("未安装tto客户端,请运行命令：curl -L https://raw.githubusercontent.com/ixre/tto/master/install | sh 安装")
		} else {
			t.Error(err)
		}
		t.FailNow()
	}
	t.Log(output)
}
