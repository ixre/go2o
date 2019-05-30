namespace java {{pkg "thrift" .global.Pkg}}.rpc
namespace csharp {{pkg "thrift" .global.Pkg}}.rpc

/** {{.table.Comment}} */
struct S{{.table.Title}}{
    {{range $i,$c:=.T.Columns}}
    /** {{$c.Comment}} */
    {{plus $c.Ordinal 1}}:{{type "thrift" $c.TypeId}} {{$c.Title}}{{end}}
}