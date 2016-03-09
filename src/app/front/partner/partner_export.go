/**
 * Copyright 2014 @ z3q.net.
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
	"github.com/jsix/gof"
	"github.com/jsix/gof/data/report"
<<<<<<< HEAD
	"net/http"
=======
	"github.com/jsix/gof/web"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"strconv"
	"strings"
)

//==================  数据导出 ===============//

type PartnerDbGetter struct{}

func (dbGetter *PartnerDbGetter) GetDB() *sql.DB {
	return gof.CurrentApp.Db().GetDb()
}

<<<<<<< HEAD
var ExpManager *report.ExportItemManager = &report.ExportItemManager{DbGetter: &PartnerDbGetter{}}
=======
var EXP_Manager *report.ExportItemManager = &report.ExportItemManager{DbGetter: &PartnerDbGetter{}}
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d

//================== 导出控制器 ==============//

// 获取导出数据
<<<<<<< HEAD
func GetExportData(r *http.Request, partnerId int) []byte {
	query := r.URL.Query()
	r.ParseForm()
	var exportItem report.IDataExportPortal = ExpManager.GetExportItem(query.Get("portal"))
=======
func GetExportData(ctx *web.Context, partnerId int) {
	r, w := ctx.Request, ctx.Response
	query := r.URL.Query()
	r.ParseForm()
	var exportItm report.IDataExportPortal = EXP_Manager.GetExportItem(query.Get("portal"))
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	//var exportItm *ExportItem = GetExportItem(query.Get("portal"))

	//fmt.Println(">>>"+strconv.FormatBool(exportItm != nil))

<<<<<<< HEAD
	if exportItem != nil {
=======
	if exportItm != nil {
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		page, rows := r.Form.Get("page"), r.Form.Get("rows")
		var parameter *report.ExportParams = report.GetExportParams(query.Get("params"), nil)

		parameter.Parameters["partner_id"] = strconv.Itoa(partnerId)

		if page != "" {
			parameter.Parameters["pageIndex"] = page
		}
		if rows != "" {
			parameter.Parameters["pageSize"] = rows
		}

<<<<<<< HEAD
		_rows, total, err := exportItem.GetSchemaAndData(parameter.Parameters)
=======
		w.Header().Add("Content-Type", "application/json")

		_rows, total, err := exportItm.GetSchemaAndData(parameter.Parameters)
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		if err == nil {
			var arr []string = []string{"{\"total\":", strconv.Itoa(total), ",\"rows\":", "", "}"}
			json, _ := json.Marshal(_rows)
			arr[3] = string(json)
<<<<<<< HEAD
			return []byte(strings.Join(arr, ""))
		}
		return []byte(`{"error":"` + err.Error() + `"}`)
	}

	return []byte(`{"error":"no such export item"}`)
=======
			w.Write([]byte(strings.Join(arr, "")))
		} else {
			w.Write([]byte(`{"error":"` + err.Error() + `"}`))
		}
	}
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
}
