package impl

#!type:0#!lang:go
#!target:{{.global.pkg}}/dao/impl/{{.table.Name}}_dao_impl.go
{{$title := .table.Title}}
{{$shortTitle := .table.ShortTitle}}
{{$p := substr .table.Name 0 1 }}
{{$structName := join (lower_title $title) "DaoImpl"}}
import(
	"database/sql"
	"fmt"
    "{{pkg "go" .global.pkg}}/dao/model"
    "{{pkg "go" .global.pkg}}/dao"
    "github.com/ixre/gof/db"
    "github.com/ixre/gof/db/orm"
    "log"
)

var _ dao.I{{$title}}Dao = new({{$structName}})
type {{$structName}} struct{
    _orm orm.Orm
}

var {{$structName}}Mapped = false

// Create new {{$title}}Dao
func New{{$title}}Dao(o orm.Orm) dao.I{{$title}}Dao{
    if !{{$structName}}Mapped{
        _ = o.Mapping(model.{{$title}}{},"{{.table.Name}}")
        {{$structName}}Mapped = true
    }
    return &{{$structName}}{
        _orm:o,
    }
}
// Get {{.table.Comment}}
func ({{$p}} *{{$structName}}) Get{{$shortTitle}}(primary interface{})*model.{{$title}}{
    e := model.{{$title}}{}
    err := {{$p}}._orm.Get(primary,&e)
    if err == nil{
        return &e
    }
    if err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$title}}")
    }
    return nil
}

// GetBy {{.table.Comment}}
func ({{$p}} *{{$structName}}) Get{{$shortTitle}}By(where string,v ...interface{})*model.{{$title}}{
    e := model.{{$title}}{}
    err := {{$p}}._orm.GetBy(&e,where,v...)
    if err == nil{
        return &e
    }
    if err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$title}}")
    }
    return nil
}

// Count {{.table.Comment}} by condition
func ({{$p}} *{{$structName}}) Count{{$shortTitle}}(where string,v ...interface{})(int,error){
   return {{$p}}._orm.Count(model.{{$title}}{},where,v...)
}

// Select {{.table.Comment}}
func ({{$p}} *{{$structName}}) Select{{$shortTitle}}(where string,v ...interface{})[]*model.{{$title}} {
    list := make([]*model.{{$title}},0)
    err := {{$p}}._orm.Select(&list,where,v...)
    if err != nil && err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$title}}")
    }
    return list
}

// Save {{.table.Comment}}
func ({{$p}} *{{$structName}}) Save{{$shortTitle}}(v *model.{{$title}})(int,error){
    id,err := orm.Save({{$p}}._orm,v,int(v.{{title .table.Pk}}))
    if err != nil && err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$title}}")
    }
    return id,err
}

// Delete {{.table.Comment}}
func ({{$p}} *{{$structName}}) Delete{{$shortTitle}}(primary interface{}) error {
    err := {{$p}}._orm.DeleteByPk(model.{{$title}}{}, primary)
    if err != nil && err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$title}}")
    }
    return err
}

// Batch Delete {{.table.Comment}}
func ({{$p}} *{{$structName}}) BatchDelete{{$shortTitle}}(where string,v ...interface{})(int64,error) {
    r,err := {{$p}}._orm.Delete(model.{{$title}}{},where,v...)
    if err != nil && err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$title}}")
    }
    return r,err
}

// Query paging data
func ({{$p}} *{{$structName}}) PagingQuery{{$shortTitle}}(begin, end int,where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
	    where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(0) FROM {{.table.Name}} WHERE %s`, where)
	_ = {{$p}}._orm.Connector().ExecScalar(s,&total)
	if total > 0{
	    s = fmt.Sprintf(`SELECT * FROM {{.table.Name}} WHERE %s %s
	        {{if eq .global.db "pgsql"}}LIMIT $2 OFFSET $1{{else}}LIMIT $1,$2{{end}}`,
            where, orderBy)
        err := {{$p}}._orm.Connector().Query(s, func(_rows *sql.Rows) {
            rows = db.RowsToMarshalMap(_rows)
        }, begin, end-begin)
        if err != nil{
            log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:{{.table.Name}})", err.Error()))
        }
	}else{
	    rows = make([]map[string]interface{},0)
	}
	return total, rows
}