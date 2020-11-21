#!type:1
#!target:{{.global.pkg}}/service/{{.table.Name}}_service.go
{{$title := .table.Title}}
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
	"github.com/ixre/gof/types"
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
func (a *{{$structName}}) Save{{$title}}(_ context.Context, r *proto.Save{{$title}}Request) (*proto.Save{{$title}}Response, error) {
	var dst *model.{{$title}}
    {{if equal_any .table.PkType 3 4 5}}\
    if r.{{.table.PkProp}} > 0 {
    {{else}}
    if r.{{.table.PkProp}} != "" {
    {{end}}
        dst = a.dao.Get(r.{{.table.PkProp}})
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
	id, err := a.dao.Save(dst)
    ret := &proto.Save{{$title}}Response{
        {{.table.PkProp}}: {{type "go" .table.PkType}}(id),
    }
    if err != nil{
        ret.ErrCode = 1
        ret.ErrMsg = err.Error()
    }
    return ret,nil
}

// 获取{{$comment}}
func (a *{{$structName}}) Get{{$title}}(_ context.Context, id *proto.{{$pkType}}) (*proto.S{{$title}}, error) {
	v := a.dao.Get(id.Value)
	if v == nil {
		return nil, nil
	}
	return &proto.S{{$title}}{
		 {{range $i,$c :=  .columns }}
		 {{ $goType := type "go" $c.Type}}\
         {{if eq $goType "int"}}{{$c.Prop}} : int32(v.{{$c.Prop}}),\
         {{else if eq $goType "int16"}}{{$c.Prop}} : int32(v.{{$c.Prop}}),\
         {{else}}{{$c.Prop}} : v.{{$c.Prop}},{{end}}{{end}}
	}, nil
}

func (a *{{$structName}}) Delete{{$title}}(_ context.Context, id *proto.{{$pkType}}) (*proto.Result, error) {
	err := a.dao.Delete(id.Value)
	return a.error(err), nil
}

func (a *{{$structName}}) PagingShops(_ context.Context, r *proto.{{$title}}PagingRequest) (*proto.{{$title}}PagingResponse, error) {
	total, rows := a.dao.PagingQuery(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.{{$title}}PagingResponse{
		Total:                int64(total),
		Value:                make([]*proto.Paging{{$title}},len(rows)),
	}
	for i,v := range rows{
	    /** #! 直接将数据库字端转换值 */
		ret.Value[i] = &proto.Paging{{$title}}{
	         {{range $i,$c := .columns }}
    		 {{ $goType := type "go" $c.Type}}\
             {{if eq $goType "int"}}{{$c.Prop}} : int32(typeconv.MustInt(v["{{$c.Name}}"])),\
             {{else if eq $goType "int16"}}{{$c.Prop}} : int32(typeconv.MustInt(v["{{$c.Name}}"])),\
             {{else if eq $goType "int64"}}{{$c.Prop}} : int64(typeconv.MustInt(v["{{$c.Name}}"])),\
             {{else if eq $goType "bool"}}{{$c.Prop}} : typeconv.MustBool(v["{{$c.Name}}"]),\
             {{else if eq $goType "float32"}}{{$c.Prop}} : float32(typeconv.MustFloat64(v["{{$c.Name}}"])),\
             {{else if eq $goType "float64"}}{{$c.Prop}} : typeconv.MustFloat64(v["{{$c.Name}}"]),\
             {{else}}{{$c.Prop}} : types.Stringify(v["{{$c.Name}}"]),{{end}}{{end}}
		}
	}
	return ret,nil
}