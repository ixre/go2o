#!target:sql/clickhouse/{{.table.Name}}.sql
{{$columns := .columns}}

CREATE TABLE IF NOT EXISTS {{.table.Name}}
(
    {{range $i,$c := .columns}}
     {{ $goType := $c.Type}}\
    `{{$c.Name}}` \
     {{if eq $goType 1}} String \
     {{else if eq $goType 2}} Int8 \
     {{else if eq $goType 3}} Int16 \
     {{else if eq $goType 4}} Int32 \
     {{else if eq $goType 5}} Int64 \
     {{else if eq $goType 6}} Float32 \
     {{else if eq $goType 7}} Float64 \
     {{else if eq $goType 14}} String \
     {{else if eq $goType 15}} Date \
     {{else if eq $goType 16}}D ecimal64 \
     {{else}} {{$goType}} {{end}} COMMENT '{{$c.Comment}}'{{if not (is_last $i $columns)}},{{end}}\
    {{end}}
) ENGINE = MergeTree
ORDER BY {{.table.Pk}}
SETTINGS index_granularity= 8192 ;
