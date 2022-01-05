package dao
#!type:0#!lang:go
#!target:{{.global.pkg}}/dao/{{.table.Name}}_dao.go
{{$title := .table.Title}}
{{$shortTitle := .table.ShortTitle}}

import(
    "{{pkg "go" .global.pkg}}/dao/model"
)

type I{{$title}}Dao interface{
    // Get{{$shortTitle}} Get {{.table.Comment}}
    Get{{$shortTitle}}(primary interface{})*model.{{$shortTitle}}
    // Get{{$shortTitle}}By GetBy {{.table.Comment}}
    Get{{$shortTitle}}By(where string,v ...interface{})*model.{{$shortTitle}}
    // Count{{$shortTitle}} Count {{.table.Comment}} by condition
    Count{{$shortTitle}}(where string,v ...interface{})(int,error)
    // Select{{$shortTitle}} Select {{.table.Comment}}
    Select{{$shortTitle}}(where string,v ...interface{})[]*model.{{$shortTitle}}
    // Save{{$shortTitle}} Save {{.table.Comment}}
    Save{{$shortTitle}}(v *model.{{$shortTitle}})(int,error)
    // Delete{{$shortTitle}} Delete {{.table.Comment}}
    Delete{{$shortTitle}}(primary interface{}) error
    // BatchDelete{{$shortTitle}} Batch Delete {{.table.Comment}}
    BatchDelete{{$shortTitle}}(where string,v ...interface{})(int64,error)
    // PagingQuery{{$shortTitle}} Query paging data
    PagingQuery{{$shortTitle}}(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})
}
