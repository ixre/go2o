package util

import (
	"os"
	"testing"

	"github.com/ixre/go2o/core/infrastructure/fw/types"
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

func TestResizeImage(t *testing.T) {
	bytes, err := MakeThumbnail("./resize.png", 200, 200)
	if err != nil {
		t.Error(err)
	}
	fi, err := os.Create("test_resize.jpg")
	if err != nil {
		t.Error(err)
	}
	defer fi.Close()
	fi.Write(bytes)
}
