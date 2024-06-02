package util

import "testing"

// 测试获取首字母
func TestGetFirstLetter(t *testing.T) {
	var str string = "上海市"
	v := GetHansFirstLetter(str)
	if v != "S" {
		t.Errorf("GetHansFirstLetter error , acture: %s", v)
	}
}
