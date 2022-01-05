package model
#!target:{{.global.pkg}}/dao/model/{{.table.Name}}_model.go
{{$shortTitle := .table.ShortTitle}}

// {{$shortTitle}} {{.table.Comment}}
type {{$shortTitle}} struct{
    {{range $i,$c := .columns}} \
    // {{$c.Comment}}
    {{$c.Prop}} {{type "go" $c.Type}} `db:"{{$c.Name}}"\
    {{if $c.IsPk}} pk:"yes"{{end}}\
    {{if $c.IsAuto}} auto:"yes"{{end}}`
    {{end}}
}
