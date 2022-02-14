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

// New{{$title}}Dao Create new {{$title}}Dao
func New{{$title}}Dao(o orm.Orm) dao.I{{$title}}Dao{
    if !{{$structName}}Mapped{
        _ = o.Mapping(model.{{$shortTitle}}{},"{{.table.Name}}")
        {{$structName}}Mapped = true
    }
    return &{{$structName}}{
        _orm:o,
    }
}
// Get{{$shortTitle}} Get {{.table.Comment}}
func ({{$p}} *{{$structName}}) Get{{$shortTitle}}(primary interface{})*model.{{$shortTitle}}{
    e := model.{{$shortTitle}}{}
    err := {{$p}}._orm.Get(primary,&e)
    if err == nil{
        return &e
    }
    if err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$shortTitle}}")
    }
    return nil
}

// Get{{$shortTitle}}By GetBy {{.table.Comment}}
func ({{$p}} *{{$structName}}) Get{{$shortTitle}}By(where string,v ...interface{})*model.{{$shortTitle}}{
    e := model.{{$shortTitle}}{}
    err := {{$p}}._orm.GetBy(&e,where,v...)
    if err == nil{
        return &e
    }
    if err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$shortTitle}}")
    }
    return nil
}

// Count{{$shortTitle}} Count {{.table.Comment}} by condition
func ({{$p}} *{{$structName}}) Count{{$shortTitle}}(where string,v ...interface{})(int,error){
   return {{$p}}._orm.Count(model.{{$shortTitle}}{},where,v...)
}

// Select{{$shortTitle}} Select {{.table.Comment}}
func ({{$p}} *{{$structName}}) Select{{$shortTitle}}(where string,v ...interface{})[]*model.{{$shortTitle}} {
    list := make([]*model.{{$shortTitle}},0)
    err := {{$p}}._orm.Select(&list,where,v...)
    if err != nil && err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$shortTitle}}")
    }
    return list
}

// Save{{$shortTitle}} Save {{.table.Comment}}
func ({{$p}} *{{$structName}}) Save{{$shortTitle}}(v *model.{{$shortTitle}})(int,error){
    id,err := orm.Save({{$p}}._orm,v,int(v.{{title .table.Pk}}))
    if err != nil && err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$shortTitle}}")
    }
    return id,err
}

// Delete{{$shortTitle}} Delete {{.table.Comment}}
func ({{$p}} *{{$structName}}) Delete{{$shortTitle}}(primary interface{}) error {
    err := {{$p}}._orm.DeleteByPk(model.{{$shortTitle}}{}, primary)
    if err != nil && err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$shortTitle}}")
    }
    return err
}

// BatchDelete{{$shortTitle}} Batch Delete {{.table.Comment}}
func ({{$p}} *{{$structName}}) BatchDelete{{$shortTitle}}(where string,v ...interface{})(int64,error) {
    r,err := {{$p}}._orm.Delete(model.{{$shortTitle}}{},where,v...)
    if err != nil && err != sql.ErrNoRows{
      log.Println("[ Orm][ Error]:",err.Error(),"; Entity:{{$shortTitle}}")
    }
    return r,err
}

// PagingQuery{{$shortTitle}} Query paging data
func ({{$p}} *{{$structName}}) PagingQuery{{$shortTitle}}(begin, end int,where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
	    where = "1=1"
	}
	query := fmt.Sprintf(`SELECT COUNT(1) FROM {{.table.Name}} WHERE %s`, where)
	_ = {{$p}}._orm.Connector().ExecScalar(query,&total)
	if total > 0{
	    query = fmt.Sprintf(`SELECT * FROM {{.table.Name}} WHERE %s %s
	        {{if eq .global.db "pgsql"}}LIMIT $2 OFFSET $1{{else}}LIMIT $1,$2{{end}}`,
            where, orderBy)
        err := {{$p}}._orm.Connector().Query(query, func(_rows *sql.Rows) {
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
