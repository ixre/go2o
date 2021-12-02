#!target:{{.global.pkg}}/restful/{{.table.Name}}_restful.go
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
  "encoding/json"
  "github.com/ixre/goex/echox"
  "github.com/ixre/gof/types/typeconv"
  "github.com/labstack/echo/v4"
  "{{pkg "go" .global.pkg}}/service"
  "{{pkg "go" .global.pkg}}/service/proto"
  "net/http"
)

var _ echox.GroupHandler = new({{$structName}});

// {{.table.Comment}}
type {{$structName}} struct{
}

func ({{$p}} *{{$structName}}) Handle(g *echo.Group) {
  // {{.table.Name}} router
  g.GET("/{{$namePath}}/paging",{{$p}}.paging{{$shortTitle}})
  g.GET("/{{$namePath}}/:id",{{$p}}.get{{$shortTitle}})
  g.GET("/{{$namePath}}",{{$p}}.query{{$shortTitle}})
  g.POST("/{{$namePath}}",{{$p}}.create{{$shortTitle}})
  g.PUT("/{{$namePath}}/:id",{{$p}}.update{{$shortTitle}})
  g.DELETE("/{{$namePath}}/:id",{{$p}}.delete{{$shortTitle}})
}

func ({{$p}} *{{$structName}}) paging{{$shortTitle}}(ctx echo.Context) error {
    /** #! 直接使用分页请求对象 */
    //dst := proto.{{$shortTitle}}PagingRequest{}
    //ctx.Bind(&dst) \
    /** #! 使用分页参数(GET)来分页 */
    bytes := GetExportDataV2(ctx, "/{{name_path .table.Name}}")
    return ctx.Blob(http.StatusOK, "application/json", bytes)
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


func ({{$p}} *{{$structName}}) query{{$shortTitle}}(ctx echo.Context) error {
  dst := &proto.Query{{$shortTitle}}Request{}
  trans, cli, _ := service.{{$title}}ServiceClient()
  defer trans.Close()
  ret, _ := cli.Query{{$shortTitle}}List(context.TODO(), dst)
  return ctx.JSON(http.StatusOK,ret.List)
}
