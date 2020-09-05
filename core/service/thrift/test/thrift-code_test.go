/**
 * Copyright 2015 @ at3.net.
 * name : thrift_test.go
 * author : jarryliu
 * date : 2016-11-17 13:37
 * description :
 * history :
 */
package test

import (
	"github.com/ixre/tto"
	"github.com/ixre/tto/utils"

	//"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/member"
	"testing"
)

var (
	//v = &member.Member{}
	//v = &member.Profile{}
	v = member.Level{}
)

// 生成Thrift结构
func TestThriftStruct(t *testing.T) {
	data, err := utils.ThriftStruct(v)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("生成代码如下:\n\n" + string(data) + "\n\n")
	}
}

// 生成结构赋值代码
func TestStructAssignCode(t *testing.T) {
	data, err := tto.GoStructAssignCode(v)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("生成代码如下:\n\n" + string(data) + "\n\n")
	}
}
