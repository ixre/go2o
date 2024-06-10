package util

import (
	"testing"

	"github.com/ixre/go2o/core/infrastructure/util/types"
)

// 测试获取首字母
func TestGetFirstLetter(t *testing.T) {
	var str string = "上海市"
	v := GetHansFirstLetter(str)
	if v != "S" {
		t.Errorf("GetHansFirstLetter error , acture: %s", v)
	}
}

func TestOrValue(t *testing.T) {
	v := types.OrValue("str", "1")
	v2 := types.OrValue("", "1")
	t.Log(v, v2)
}
