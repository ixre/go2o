/**
 * Copyright 2015 @ at3.net.
 * name : generator_test.go
 * author : jarryliu
 * date : 2016-11-17 13:58
 * description :
 * history :
 */
package test

import (
	"fmt"
	"github.com/ixre/gof/shell"
	"testing"
)

const entryFile = "../idl/service.thrift"

// 生成Golang的Thrift代码
func TestGo(t *testing.T) {
	genCode(t, "go", entryFile)
}

// 生成Golang的Thrift代码
func TestJava(t *testing.T) {
	genCode(t, "java", entryFile)
}

func genCode(t *testing.T, lang string, file string) {
	cmd := fmt.Sprintf("thrift -r -gen %s %s", lang, file)
	_, output, err := shell.Run(cmd)
	if err == nil {
		t.Log("生成成功!")
		return
	}
	t.Log(output + "\n")
	t.Fail()
}
