#!type:1
#!target:{{.global.pkg}}/service/{{.table.Name}}_service.go
{{$title := .table.Title}}
{{$shortTitle := .table.ShortTitle}}
{{$p := substr .table.Name 0 1 }}
{{$pkName := .table.Pk}}
{{$comment := .table.Comment}}
/** #! 主键对象类型 */
{{$pkType := join .table.Title .table.PkProp }}
/** #! 服务结构名称 */
{{$structName := join (lower_title .table.Title) "ServiceImpl"}}
package impl

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : {{.table.Name}}_service.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : {{.global.time}}
 * description :
 * history :
 */

import (
	"context"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types/typeconv"
	"{{.global.pkg}}/dao"
	"{{.global.pkg}}/dao/impl"
	"{{.global.pkg}}/dao/model"
	"{{.global.pkg}}/proto"
	"time"
)

var _ proto.{{$title}}ServiceServer = new({{$structName}})

type {{$structName}} struct {
	dao dao.I{{.table.Title}}Dao
	s   storage.Interface
	serviceUtil


}

func New{{$title}}Service(s storage.Interface, o orm.Orm) *{{$structName}} {
	return &{{$structName}}{
		s:   s,
		dao: impl.New{{.table.Title}}Dao(o),
	}
}

// 保存{{$comment}}
func ({{$p}} *{{$structName}}) Save{{$shortTitle}}(_ context.Context, r *proto.Save{{$title}}Request) (*proto.Save{{$title}}Response, error) {
	var dst *model.{{$title}}
    {{if equal_any .table.PkType 3 4 5}}\
    if r.{{.table.PkProp}} > 0 {
    {{else}}
    if r.{{.table.PkProp}} != "" {
    {{end}}
        if dst = {{$p}}.dao.Get{{$shortTitle}}(r.{{.table.PkProp}}); dst == nil{
            return &proto.Save{{$shortTitle}}Response{
                ErrCode: 2,
                ErrMsg:  "no such record",
            }, nil
        }
    } else {
        dst = &model.{{$title}}{}
        {{$c := try_get .columns "create_time"}} \
        {{if $c}}dst.CreateTime = time.Now().Unix(){{end}}
    }
    /** #! 为对象赋值 */
    {{range $i,$c := exclude .columns $pkName "create_time" "update_time"}}
    {{ $goType := type "go" $c.Type}}\
    {{if eq $goType "int"}}dst.{{$c.Prop}} = int(r.{{$c.Prop}})\
    {{else if eq $goType "int16"}}dst.{{$c.Prop}} = int16(r.{{$c.Prop}})\
    {{else if eq $goType "int32"}}dst.{{$c.Prop}} = int32(r.{{$c.Prop}})\
    {{else}}dst.{{$c.Prop}} = r.{{$c.Prop}}{{end}}{{end}}

    {{$c := try_get .columns "update_time"}}
    {{if $c}}dst.UpdateTime = time.Now().Unix(){{end}}\
	id, err := {{$p}}.dao.Save{{$shortTitle}}(dst)
    ret := &proto.Save{{$shortTitle}}Response{
        {{.table.PkProp}}: {{type "go" .table.PkType}}(id),
    }
    if err != nil{
        ret.ErrCode = 1
        ret.ErrMsg = err.Error()
    }
    return ret,nil
}

func ({{$p}} *{{$structName}}) parse{{$shortTitle}}(v *model.{{$title}}) *proto.S{{$title}} {
	return &proto.S{{$shortTitle}}{
	 {{range $i,$c :=  .columns }}
     {{ $goType := type "go" $c.Type}}\
     {{if eq $goType "int"}}{{$c.Prop}} : int32(v.{{$c.Prop}}),\
     {{else if eq $goType "int16"}}{{$c.Prop}} : int32(v.{{$c.Prop}}),\
     {{else}}{{$c.Prop}} : v.{{$c.Prop}},{{end}}{{end}}
	}
}

// 获取{{$comment}}
func ({{$p}} *{{$structName}}) Get{{$shortTitle}}(_ context.Context, id *proto.{{$pkType}}) (*proto.S{{$title}}, error) {
	v := {{$p}}.dao.Get{{$shortTitle}}(id.Value)
	if v == nil {
		return nil, nil
	}
	return {{$p}}.parse{{$shortTitle}}(v), nil
}

// 获取{{$comment}}列表
func ({{$p}} *{{$structName}}) Query{{$shortTitle}}List(_ context.Context, r *proto.Query{{$title}}Request) (*proto.Query{{$title}}Response, error) {
	arr := {{$p}}.dao.Select{{$shortTitle}}("")
	ret := &proto.Query{{$shortTitle}}Response{
		List:make([]*proto.S{{$shortTitle}},len(arr)),
	}
	for i,v := range arr{
		ret.List[i] = {{$p}}.parse{{$shortTitle}}(v)
	}
	return ret,nil
}

func ({{$p}} *{{$structName}}) Delete{{$shortTitle}}(_ context.Context, id *proto.{{$pkType}}) (*proto.Result, error) {
	err := {{$p}}.dao.Delete{{$shortTitle}}(id.Value)
	return {{$p}}.error(err), nil
}

func ({{$p}} *{{$structName}}) Paging{{$shortTitle}}(_ context.Context, r *proto.{{$title}}PagingRequest) (*proto.{{$title}}PagingResponse, error) {
	total, rows := {{$p}}.dao.PagingQuery{{$shortTitle}}(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.{{$shortTitle}}PagingResponse{
		Total:                int64(total),
		Value:                make([]*proto.Paging{{$shortTitle}},len(rows)),
	}
	for i,v := range rows{
	    /** #! 直接将数据库字端转换值 */
		ret.Value[i] = &proto.Paging{{$shortTitle}}{
	         {{range $i,$c := .columns }}
    		 {{ $goType := type "go" $c.Type}}\
             {{if eq $goType "int"}}{{$c.Prop}} : int32(typeconv.MustInt(v["{{$c.Name}}"])),\
             {{else if eq $goType "int16"}}{{$c.Prop}} : int32(typeconv.MustInt(v["{{$c.Name}}"])),\
             {{else if eq $goType "int64"}}{{$c.Prop}} : int64(typeconv.MustInt(v["{{$c.Name}}"])),\
             {{else if eq $goType "bool"}}{{$c.Prop}} : typeconv.MustBool(v["{{$c.Name}}"]),\
             {{else if eq $goType "float32"}}{{$c.Prop}} : float32(typeconv.MustFloat64(v["{{$c.Name}}"])),\
             {{else if eq $goType "float64"}}{{$c.Prop}} : typeconv.MustFloat64(v["{{$c.Name}}"]),\
             {{else}}{{$c.Prop}} : typeconv.Stringify(v["{{$c.Name}}"]),{{end}}{{end}}
		}
	}
	return ret,nil
}