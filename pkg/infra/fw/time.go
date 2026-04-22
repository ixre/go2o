package fw

import (
	"time"
)

// getMonthStartEndUnix 返回本月开始和结束的Unix时间戳（秒）
func GetMonthStartEndUnix(s int64) (int64, int64) {
	// 获取当前时间
	now := time.Unix(s, 0)
	// 获取本月第一天（即当前月份的第一天，时间设为00:00:00 UTC）
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	// 获取本月最后一天（可以通过当前月份+1月，然后减去1天来得到）
	// 注意：AddDate不会改变原始时间的时间（小时、分钟、秒），因此不需要担心时间偏移
	endOfMonth := startOfMonth.AddDate(0, 1, -1).Add(24*time.Hour - time.Second)
	// 转换为Unix时间戳（秒）
	return startOfMonth.Unix(), endOfMonth.Unix()
}
