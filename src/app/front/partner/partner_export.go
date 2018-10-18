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
	"github.com/jsix/goex/report"
	"net/http"
	"strconv"
	"strings"
)

//==================  数据导出 ===============//

type PartnerDbGetter struct{}

func (dbGetter *PartnerDbGetter) GetDB() *sql.DB {
	return gof.CurrentApp.Db().Raw()
}

//var ExpManager *report.ItemManager = &report.ItemManager{DbGetter: &PartnerDbGetter{}}
var EXP_Manager = report.NewItemManager(&PartnerDbGetter{},"./conf/query",true)

//================== 导出控制器 ==============//

// 获取导出数据
func GetExportData(r *http.Request, partnerId int) []byte {
	query := r.URL.Query()
	r.ParseForm()
	var exportItem report.IDataExportPortal = EXP_Manager.GetItem(query.Get("portal"))
	//var exportItm *ExportItem = GetExportItem(query.Get("portal"))

	//fmt.Println(">>>"+strconv.FormatBool(exportItm != nil))

	if exportItem != nil {
		page, rows := r.Form.Get("page"), r.Form.Get("rows")
		//var parameter *report.Params = report.ParseParams(query.Get("params"))
		var parameter report.Params = report.ParseParams(query.Get("params"))

		//parameter.Parameters["partner_id"] = strconv.Itoa(partnerId)
		parameter["partner_id"] = strconv.Itoa(partnerId)

		if page != "" {
			//parameter.Parameters["pageIndex"] = page
			parameter["pageIndex"] = page
		}
		if rows != "" {
			//parameter.Parameters["pageSize"] = rows
			parameter["pageSize"] = rows
		}

		//_rows, total, err := exportItem.GetSchemaAndData(parameter.Parameters)
		_rows, total, err := exportItem.GetSchemaAndData(parameter)
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
