#!target:{{.global.pkg}}/grpc-restful/{{.table.Name}}_restful.go
package restful
{{$title := .table.Title}}
{{$shortTitle := .table.ShortTitle}}
{{$structName := join (lower_title $title) "Resource"}}
{{$p := substr .table.Name 0 1 }}
{{$namePath := name_path .table.Name}}
{{$pk := .table.Pk}}
{{$pkType := join .table.Title .table.PkProp ""}}

import (
  "context"
  "github.com/ixre/gof/types/typeconv"
  "github.com/labstack/echo/v4"
  "{{pkg "go" .global.pkg}}/service"
  "{{pkg "go" .global.pkg}}/service/proto"
  "net/http"
)

// {{.table.Comment}}
type {{$structName}} struct{
}

func ({{$p}} {{$structName}}) Routes(g *echo.Group) {
  // {{.table.Name}} router
  g.GET("/{{$namePath}}/paging",{{$p}}.paging{{$shortTitle}})
  g.GET("/{{$namePath}}/:id",{{$p}}.get{{$shortTitle}})
  g.GET("/{{$namePath}}",{{$p}}.query{{$shortTitle}})
  g.POST("/{{$namePath}}",{{$p}}.create{{$shortTitle}})
  g.DELETE("/{{$namePath}}/:id",{{$p}}.delete{{$shortTitle}})
  g.PUT("/{{$namePath}}/:id",{{$p}}.update{{$shortTitle}})
}

func ({{$p}} *{{$structName}}) get{{$shortTitle}}(ctx echo.Context) error {
  /** #! 转换主键 */
  {{ $goType := type "protobuf" .table.PkType}}\
  {{if eq $goType "int32"}}{{$pk}} := int32(typeconv.MustInt(ctx.Param("id")))\
  {{else if eq $goType "int64"}}{{$pk}} := int64(typeconv.MustInt(ctx.Param("id")))\
  {{else}}{{$pk}} := ctx.Param("id"){{end}}
  trans,cli,_ := service.{{$title}}ServiceClient()
  defer trans.Close()
  ret, _ := cli.Get{{$shortTitle}}(context.TODO(), &proto.{{$pkType}}{Value:{{$pk}}})
  return ctx.JSON(http.StatusOK,ret)
}

func ({{$p}} *{{$structName}}) create{{$shortTitle}}(ctx echo.Context) error {
  dst := proto.Save{{$shortTitle}}Request{}
  err := ctx.Bind(&dst)
  if err == nil{
    trans,cli,_ := service.{{$title}}ServiceClient()
    defer trans.Close()
    ret, _ := cli.Save{{$shortTitle}}(context.TODO(), &dst)
    return ctx.JSON(http.StatusOK,ret)
  }
  return err
}

func ({{$p}} *{{$structName}}) delete{{$shortTitle}}(ctx echo.Context) error {
  /** #! 转换主键 */
  {{ $goType := type "protobuf" .table.PkType}}\
  {{if eq $goType "int32"}}{{$pk}} := int32(typeconv.MustInt(ctx.Param("id")))\
  {{else if eq $goType "int64"}}{{$pk}} := int64(typeconv.MustInt(ctx.Param("id")))\
  {{else}}{{$pk}} := ctx.Param("id"){{end}}
  trans, cli, _ := service.{{$title}}ServiceClient()
  defer trans.Close()
  ret, _ := cli.Delete{{$shortTitle}}(context.TODO(), &proto.{{$pkType}}{Value: {{$pk}}})
  return ctx.JSON(http.StatusOK, ret)
}

func ({{$p}} *{{$structName}}) update{{$shortTitle}}(ctx echo.Context) error {
  dst := proto.Save{{$shortTitle}}Request{}
  err := ctx.Bind(&dst)
  if err == nil{
    /** #! 转换主键 */
    {{ $goType := type "protobuf" .table.PkType}}\
    {{if eq $goType "int32"}}{{$pk}} := int32(typeconv.MustInt(ctx.Param("id")))\
    {{else if eq $goType "int64"}}{{$pk}} := int64(typeconv.MustInt(ctx.Param("id")))\
    {{else}}{{$pk}} := ctx.Param("id"){{end}}
    dst.{{.table.PkProp}} = {{$pk}}
    trans, cli, _ := service.{{$title}}ServiceClient()
    defer trans.Close()
    ret, _ := cli.Save{{$shortTitle}}(context.TODO(), &dst)
    return ctx.JSON(http.StatusOK, ret)
  }
  return err
}

func ({{$p}} *{{$structName}}) paging{{$shortTitle}}(ctx echo.Context) error {
    /** #! 直接使用分页请求对象 */
    //dst := proto.{{$shortTitle}}PagingRequest{}
    //ctx.Bind(&dst) \
    /** #! 使用分页参数(GET)来分页 */
    page := typeconv.MustInt(ctx.QueryParam("page"))
    size := typeconv.MustInt(ctx.QueryParam("rows"))
    mp := make(map[string]interface{}, 0)
    _ = json.Unmarshal([]byte(ctx.QueryParam("params")), &mp)
    params := &proto.SPagingParams{
        Begin:  int64((page - 1) * size),
        End:    int64(page * size),
        Where:  ctx.QueryParam("where"),
        SortBy: typeconv.Stringify(mp["order_by"]),
        Parameters: map[string]string{
            "keyword": typeconv.Stringify(mp["keyword"]),
            "state":   typeconv.Stringify(mp["state"]),
        },
    }
    dst := proto.{{$shortTitle}}PagingRequest{
        Params:params,
    }
    trans, cli, _ := service.{{$title}}ServiceClient()
    defer trans.Close()
    ret, _ := cli.Paging{{$shortTitle}}(context.TODO(), &dst)
    if ret.Value == nil{
        ret.Value = make([]*proto.Paging{{$shortTitle}},0)
    }
    /** 转换为小写的字段 */
    rows :=make([]map[string]interface{}, len(ret.Value))
    data := map[string]interface{}{
        "total": ret.Total,
        "rows": rows,
    }
    for i,v := range ret.Value{
        rows[i] = map[string]interface{}{
            {{range $i,$c := .columns}}
            "{{$c.Name}}":v.{{$c.Prop}},
            {{end}}
        }
    }
    return ctx.JSON(http.StatusOK, data)
}

func ({{$p}} *{{$structName}}) query{{$shortTitle}}(ctx echo.Context) error {
  dst := &proto.Query{{$shortTitle}}Request{}
  trans, cli, _ := service.{{$title}}ServiceClient()
  defer trans.Close()
  ret, _ := cli.Query{{$shortTitle}}List(context.TODO(), dst)
  return ctx.JSON(http.StatusOK,ret.Value)
}
