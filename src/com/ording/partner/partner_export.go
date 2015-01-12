package partner

import (
	"com/share/glob"
	"database/sql"
	"encoding/json"
	"github.com/atnet/gof/data/transfer"
	"net/http"
	"strconv"
	"strings"
)

//==================  数据导出 ===============//

type PartnerDbGetter struct{}

func (dbGetter *PartnerDbGetter) GetDB() *sql.DB {
	return glob.CurrContext().Db().GetDb()
}

var EXP_Manager *transfer.ExportItemManager = &transfer.ExportItemManager{DbGetter: &PartnerDbGetter{}}

//================== 导出控制器 ==============//

// 获取导出数据
func GetExportData(w http.ResponseWriter, r *http.Request, ptId int) {
	query := r.URL.Query()
	r.ParseForm()
	var exportItm transfer.IDataExportPortal = EXP_Manager.GetExportItem(query.Get("portal"))
	//var exportItm *ExportItem = GetExportItem(query.Get("portal"))

	//fmt.Println(">>>"+strconv.FormatBool(exportItm != nil))

	if exportItm != nil {
		page, rows := r.Form.Get("page"), r.Form.Get("rows")
		var parameter *transfer.ExportParams = transfer.GetExportParams(query.Get("params"), nil)

		parameter.Parameters["partnerId"] = strconv.Itoa(ptId)

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
