package daemon

import (
	"testing"
	"time"

	"github.com/ixre/go2o/core/infrastructure/util"
)

// 测试生成商户的报表
func TestGenerateMchDayChart(t *testing.T) {
	dt := time.Now().Add(time.Hour * -24 * 15)
	for i := 0; i < 15; i++ {
		st, et := util.GetStartEndUnix(dt.Add(time.Hour * 24 * time.Duration(i)))
		generateMchDayChart(st, et)
	}
}
