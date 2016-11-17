/**
 * Copyright 2015 @ at3.net.
 * name : generator_test.go
 * author : jarryliu
 * date : 2016-11-17 13:58
 * description :
 * history :
 */
package idl

import (
	"fmt"
	"github.com/jsix/gof/shell"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func getFiles() []string {
	list := []string{}
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			if strings.HasSuffix(path, ".thrift") {
				list = append(list, path)
			}
		}
		return err
	})
	return list
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

// 生成Golang的Thrift代码
func TestGo(t *testing.T) {
	for _, v := range getFiles() {
		genCode(t, "go", v)
	}
}

// 生成Golang的Thrift代码
func TestJava(t *testing.T) {
	for _, v := range getFiles() {
		genCode(t, "java", v)
	}
}
