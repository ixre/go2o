package daemon

import (
	"github.com/ixre/go2o/core/infrastructure/tool"
	"testing"
	"time"
)

// 测试生成商户的报表
func TestGenerateMchDayChart(t *testing.T) {
	dt := time.Now().Add(time.Hour * -24 * 15)
	for i := 0; i < 15; i++ {
		st, et := tool.GetStartEndUnix(dt.Add(time.Hour * 24 * time.Duration(i)))
		generateMchDayChart(st, et)
	}
}
