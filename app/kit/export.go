package kit

import (
	"github.com/jsix/goex/report"
	"go2o/core/domain/interface/enum"
	"strconv"
	"strings"
	"time"
)

var DefaultExportFormatter report.IExportFormatter = &exportFormatter{}

type exportFormatter struct {
}

func (e *exportFormatter) Format(field, name string, rowNum int, data interface{}) interface{} {

	// 格式化时间
	if strings.Index(field, "_time") != -1 || strings.Index(name, "时间") != -1 {
		s := data.(string)
		if len(s) == 10 {
			i64, err := strconv.ParseInt(s, 0, 64)
			if err == nil {
				dt := time.Unix(i64, 0)
				return dt.Format("2006-01-02 15:04:05")
			}
		}
		return data
	}
	// 格式化是否
	if strings.HasPrefix(field, "is_") || strings.Index(name, "是否") != -1 {
		if data == "1" || data == "true" || data == 1 {
			return "是"
		}
		return "否"
	}
	// 性别
	if field == "sex" || strings.Index(name, "性别") != -1 {
		if data == "1" || data == "female" {
			return "男"
		}
		return "女"
	}
	// 审核状态
	if name == "review_state" || strings.HasPrefix(name, "审核") ||
		strings.HasPrefix(name, "review_") {
		i, err := strconv.Atoi(data.(string))
		if err == nil {
			return enum.ReviewString(int32(i))
		}
		return "-"
	}
	// 上架状态
	if field == "shelve_state" || strings.HasPrefix(name, "上架") {
		switch data.(string) {
		case "1":
			return "已下架"
		case "2":
			return "已上架"
		case "3":
			return "已违规下架"
		}
		return "待上架"
	}
	return data
}
