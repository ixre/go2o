/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package partner

import (
	"database/sql"
	"encoding/json"
	"github.com/atnet/gof/data/report"
	"github.com/atnet/gof/web"
	"go2o/src/core"
	"strconv"
	"strings"
)

//==================  数据导出 ===============//

type PartnerDbGetter struct{}

func (dbGetter *PartnerDbGetter) GetDB() *sql.DB {
	return core.GlobalApp.Db().GetDb()
}

var EXP_Manager *report.ExportItemManager = &report.ExportItemManager{DbGetter: &PartnerDbGetter{}}

//================== 导出控制器 ==============//

// 获取导出数据
func GetExportData(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.ResponseWriter
	query := r.URL.Query()
	r.ParseForm()
	var exportItm report.IDataExportPortal = EXP_Manager.GetExportItem(query.Get("portal"))
	//var exportItm *ExportItem = GetExportItem(query.Get("portal"))

	//fmt.Println(">>>"+strconv.FormatBool(exportItm != nil))

	if exportItm != nil {
		page, rows := r.Form.Get("page"), r.Form.Get("rows")
		var parameter *report.ExportParams = report.GetExportParams(query.Get("params"), nil)

		parameter.Parameters["partnerId"] = strconv.Itoa(partnerId)

		if page != "" {
			parameter.Parameters["pageIndex"] = page
		}
		if rows != "" {
			parameter.Parameters["pageSize"] = rows
		}

		w.Header().Add("Content-Type", "application/json")

		_rows, total, err := exportItm.GetSchemaAndData(parameter.Parameters)
		if err == nil {
			var arr []string = []string{"{\"total\":", strconv.Itoa(total), ",\"rows\":", "", "}"}
			json, _ := json.Marshal(_rows)
			arr[3] = string(json)
			w.Write([]byte(strings.Join(arr, "")))
		} else {
			w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		}
	}
}
