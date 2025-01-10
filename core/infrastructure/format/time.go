/**
 * Copyright 2015 @ 56x.net.
 * name : time
 * author : jarryliu
 * date : 2015-11-07 22:52
 * description :
 * history :
 */
package format

import (
	"fmt"
	"time"
)

func HanDateTime(t time.Time) string {
	return t.Format("2006年01月02日 15:04")
}

func HanUnixDateTime(unix int64) string {
	return HanDateTime(time.Unix(unix, 0))
}

func UnixTimeStr(unix int64) string {
	return time.Unix(unix, 0).Format("2006-01-02 15:04")
}

// 将秒数转换为分钟和秒
func DurationString(seconds int) string {
	minutes := seconds / 60
	seconds = seconds % 60
	if seconds == 0 {
		return fmt.Sprintf("%d分", minutes)
	}
	return fmt.Sprintf("%d分%d秒", minutes, seconds)
}
