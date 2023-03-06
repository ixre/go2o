package regex

import (
	"strings"
	"testing"
)

func TestContainInvalidChars(t *testing.T) {
	txt := "测试商品%$@#*[哈哈哈]"
	b,arr := ContainInvalidChars(txt)
	if b{
		t.Logf("特殊字符包含%s",strings.Join(arr,","))
	}else{
		t.Log("未包含特殊字符")
	}
}
