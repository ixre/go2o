package dao
#!type:0#!lang:go
#!target:{{.global.pkg}}/dao/{{.table.Name}}_dao.go
{{$title := .table.Title}}
{{$shortTitle := .table.ShortTitle}}

import(
    "{{pkg "go" .global.pkg}}/dao/model"
)

type I{{$title}}Dao interface{
    // Get {{.table.Comment}}
    Get{{$shortTitle}}(primary interface{})*model.{{$shortTitle}}
    // GetBy {{.table.Comment}}
    Get{{$shortTitle}}By(where string,v ...interface{})*model.{{$shortTitle}}
    // Count {{.table.Comment}} by condition
    Count{{$shortTitle}}(where string,v ...interface{})(int,error)
    // Select {{.table.Comment}}
    Select{{$shortTitle}}(where string,v ...interface{})[]*model.{{$shortTitle}}
    // Save {{.table.Comment}}
    Save{{$shortTitle}}(v *model.{{$shortTitle}})(int,error)
    // Delete {{.table.Comment}}
    Delete{{$shortTitle}}(primary interface{}) error
    // Batch Delete {{.table.Comment}}
    BatchDelete{{$shortTitle}}(where string,v ...interface{})(int64,error)
    // Query paging data
    PagingQuery{{$shortTitle}}(begin, end int, where, orderBy string) (total int, rows []map[string]interface{})
}
