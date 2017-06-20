package kit

import (
	"github.com/jsix/goex/report"
	"strconv"
	"strings"
	"time"
)

var DefaultExportFormatter report.IExportFormatter = &exportFormatter{}

type exportFormatter struct {
}

func (e *exportFormatter) Format(field, name string, data interface{}) interface{} {

	// 格式化是否
	if strings.HasPrefix(field, "is_") || strings.Index(name, "是否") != -1 {
		if data == "1" || data == "true" || data == 1 {
			return "是"
		}
		return "否"
	}
	// 格式化时间
	if strings.Index(field, "_time") != -1 || strings.Index(name, "时间") != -1 {
		s := data.(string)
		if len(s) == 10 {
			i64, err := strconv.ParseInt(s, 0, 64)
			if err == nil {
				dt := time.Unix(i64, 0)
				return dt.Format("2006-01-02 15:03:04")
			}
		}
		return data
	}

	return data
}
