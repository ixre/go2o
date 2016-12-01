/**
 * Copyright 2015 @ at3.net.
 * name : thrift_test.go
 * author : jarryliu
 * date : 2016-11-17 13:37
 * description :
 * history :
 */
package idl

import (
	"github.com/jsix/gof/generator"
	"go2o/core/dto"
	"testing"
)

var (
	//v = &member.Member{}
	//v = &member.Profile{}
	v = dto.MemberSummary{}
)

// 生成Thrift结构
func TestThriftStruct(t *testing.T) {
	data, err := generator.ThriftStruct(v)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("生成代码如下:\n\n" + string(data) + "\n\n")
	}
}

// 生成结构赋值代码
func TestStructAssignCode(t *testing.T) {
	data, err := generator.StructAssignCode(v)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("生成代码如下:\n\n" + string(data) + "\n\n")
	}
}
