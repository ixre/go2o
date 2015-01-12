package transfer

import (
	_ "database/sql"
	"encoding/xml"
	"io/ioutil"
	"regexp"
	"strings"
)

//列映射
type DataColumnMapping struct {
	//列的字段
	Field string
	//列的名称
	Name string
}

//导入导出项目配置
type ExportItemConfig struct {
	ColumnMappingString string
	Query               string
	Total               string
	Import              string
}

//导出参数
type ExportParams struct {
	Parameters map[string]string
	//要到导出的列(对应IDataExportPortal的ColumnNames或DataTable的Shema
	ColumnNames []string
}

//数据导出入口
type IDataExportPortal interface {
	//导出的列名(比如：数据表是因为列，这里我需要列出中文列)
	//ColumnNames() (names []DataColumnMapping)
	//获取要导出的数据及表结构
	GetSchemaAndData(ht map[string]string) (rows []map[string]interface{}, total int, err error)
	//获取要导出的数据Json格式
	GetJsonData(ht map[string]string) string
	//获取统计数据
	GetTotalView(ht map[string]string) (row map[string]interface{})
	//根据参数获取导出列名及导出名称
	GetExportColumnIndexAndName(exportColumnNames []string) (dict map[string]string)
}

//导出
type IDataExportProvider interface {
	//导出
	Export(rows []map[string]interface{}, columns map[string]string) (binary []byte)
}

type DataExportPortal struct {
	//导出的列名(比如：数据表是因为列，这里我需要列出中文列)
	ColumnNames []DataColumnMapping
}

//根据参数获取导出列名及导出名称
func (portal *DataExportPortal) GetExportColumnIndexAndName(
	exportColumnNames []string) (dict map[string]string) {
	dict = make(map[string]string)
	for _, cName := range exportColumnNames {
		for _, cMap := range portal.ColumnNames {
			if cMap.Name == cName {
				dict[cMap.Field] = cMap.Name
				break
			}
		}
	}
	return dict
}

//获取列映射
func GetColumnMappings(columnMappingString string) (
	columnsMapping []DataColumnMapping, err error) {
	re, err := regexp.Compile("([^:]+):([;]*)")
	if err != nil {
		return nil, err
	}

	var matches [][]string = re.FindAllStringSubmatch(columnMappingString, 0)
	if matches == nil {
		return nil, nil
	}
	columnsMapping = make([]DataColumnMapping, 0, len(matches))
	for i, v := range matches {
		columnsMapping[i] = DataColumnMapping{Field: v[1], Name: v[2]}
	}
	return columnsMapping, nil
}

//获取列映射数组
func LoadExportConfigFromXml(xmlFilePath string) (*ExportItemConfig, error) {
	var cfg ExportItemConfig
	content, _err := ioutil.ReadFile(xmlFilePath)
	if _err != nil {
		return &ExportItemConfig{}, _err
	}
	err := xml.Unmarshal(content, &cfg)
	return &cfg, err
}

func Export(portal IDataExportPortal, parameters ExportParams,
	provider IDataExportProvider) []byte {
	rows, _, _ := portal.GetSchemaAndData(parameters.Parameters)
	return provider.Export(rows, portal.GetExportColumnIndexAndName(
		parameters.ColumnNames))
}

func GetExportParams(paramMappings string, columnNames []string) *ExportParams {

	var parameters map[string]string = make(map[string]string)

	if paramMappings != "" {

		var paramsArr, splitArr []string

		paramsArr = strings.Split(paramMappings, ";")

		//添加传入的参数
		for _, v := range paramsArr {
			splitArr = strings.Split(v, ":")
			parameters[splitArr[0]] = v[len(splitArr[0])+1:]
		}

	}
	return &ExportParams{ColumnNames: columnNames, Parameters: parameters}

}

// 格式化sql语句
func SqlFormat(sql string, ht map[string]string) (formatted string) {
	formatted = sql
	for k, v := range ht {
		formatted = strings.Replace(formatted, "{"+k+"}", v, 20)
	}
	return formatted
}
