/**
 * Copyright 2015 @ z3q.net.
 * name : time
 * author : jarryliu
 * date : 2015-11-07 22:52
 * description :
 * history :
 */
package format

import "time"

func HanDateTime(t time.Time) string {
	return t.Format("2006年01月02日 15:04")
}

func HanUnixDateTime(unix int64) string {
	return HanDateTime(time.Unix(unix, 0))
}
