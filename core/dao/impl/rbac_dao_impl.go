package impl

import (
	"database/sql"
	"fmt"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"go2o/core/dao"
	"go2o/core/dao/model"
	"log"
)

var _ dao.IRbacDao = new(rbacDaoImpl)

type rbacDaoImpl struct {
	_orm orm.Orm
}

var rbacDaoImplMapped = false

// Create new PermDeptDao
func NewRbacDao(o orm.Orm) dao.IRbacDao {
	if !rbacDaoImplMapped {
		_ = o.Mapping(model.PermDept{}, "perm_dept")
		rbacDaoImplMapped = true
	}
	return &rbacDaoImpl{
		_orm: o,
	}
}

// Get 部门
func (t *rbacDaoImpl) GetPermDept(primary interface{}) *model.PermDept {
	e := model.PermDept{}
	err := t._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return nil
}

// GetBy 部门
func (t *rbacDaoImpl) GetPermDeptBy(where string, v ...interface{}) *model.PermDept {
	e := model.PermDept{}
	err := t._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return nil
}

// Count 部门 by condition
func (t *rbacDaoImpl) CountPermDept(where string, v ...interface{}) (int, error) {
	return t._orm.Count(model.PermDept{}, where, v...)
}

// Select 部门
func (t *rbacDaoImpl) SelectPermDept(where string, v ...interface{}) []*model.PermDept {
	list := make([]*model.PermDept, 0)
	err := t._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return list
}

// Save 部门
func (t *rbacDaoImpl) SavePermDept(v *model.PermDept) (int, error) {
	id, err := orm.Save(t._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return id, err
}

// Delete 部门
func (t *rbacDaoImpl) DeletePermDept(primary interface{}) error {
	err := t._orm.DeleteByPk(model.PermDept{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return err
}

// Batch Delete 部门
func (t *rbacDaoImpl) BatchDeletePermDept(where string, v ...interface{}) (int64, error) {
	r, err := t._orm.Delete(model.PermDept{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return r, err
}

// Query paging data
func (t *rbacDaoImpl) PagingQueryPermDept(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(0) FROM perm_dept WHERE %s`, where)
	t._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM perm_dept WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		t._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}
