/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package merchant

import (
	"database/sql"
	"encoding/json"
	"github.com/jsix/gof"
	"github.com/jsix/gof/data/report"
	"net/http"
	"strconv"
	"strings"
)

//==================  数据导出 ===============//

type PartnerDbGetter struct{}

func (dbGetter *PartnerDbGetter) GetDB() *sql.DB {
	return gof.CurrentApp.Db().GetDb()
}

var ExpManager *report.ExportItemManager = &report.ExportItemManager{DbGetter: &PartnerDbGetter{}}

//================== 导出控制器 ==============//

// 获取导出数据
func GetExportData(r *http.Request, partnerId int) []byte {
	query := r.URL.Query()
	r.ParseForm()
	var exportItem report.IDataExportPortal = ExpManager.GetExportItem(query.Get("portal"))
	//var exportItm *ExportItem = GetExportItem(query.Get("portal"))

	//fmt.Println(">>>"+strconv.FormatBool(exportItm != nil))

	if exportItem != nil {
		page, rows := r.Form.Get("page"), r.Form.Get("rows")
		var parameter *report.ExportParams = report.GetExportParams(query.Get("params"), nil)

		parameter.Parameters["merchant_id"] = strconv.Itoa(partnerId)

		if page != "" {
			parameter.Parameters["pageIndex"] = page
		}
		if rows != "" {
			parameter.Parameters["pageSize"] = rows
		}

		_rows, total, err := exportItem.GetSchemaAndData(parameter.Parameters)
		if err == nil {
			var arr []string = []string{"{\"total\":", strconv.Itoa(total), ",\"rows\":", "", "}"}
			json, _ := json.Marshal(_rows)
			arr[3] = string(json)
			return []byte(strings.Join(arr, ""))
		}
		return []byte(`{"error":"` + err.Error() + `"}`)
	}

	return []byte(`{"error":"no such export item"}`)
}
