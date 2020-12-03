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
		_ = o.Mapping(model.PermJob{}, "perm_job")
		rbacDaoImplMapped = true
	}
	return &rbacDaoImpl{
		_orm: o,
	}
}

// Get 部门
func (p *rbacDaoImpl) GetPermDept(primary interface{}) *model.PermDept {
	e := model.PermDept{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return nil
}

// GetBy 部门
func (p *rbacDaoImpl) GetPermDeptBy(where string, v ...interface{}) *model.PermDept {
	e := model.PermDept{}
	err := p._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return nil
}

// Count 部门 by condition
func (p *rbacDaoImpl) CountPermDept(where string, v ...interface{}) (int, error) {
	return p._orm.Count(model.PermDept{}, where, v...)
}

// Select 部门
func (p *rbacDaoImpl) SelectPermDept(where string, v ...interface{}) []*model.PermDept {
	list := make([]*model.PermDept, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return list
}

// Save 部门
func (p *rbacDaoImpl) SavePermDept(v *model.PermDept) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return id, err
}

// Delete 部门
func (p *rbacDaoImpl) DeletePermDept(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PermDept{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return err
}

// Batch Delete 部门
func (p *rbacDaoImpl) BatchDeletePermDept(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PermDept{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return r, err
}

// Get 岗位
func (p *rbacDaoImpl) GetPermJob(primary interface{}) *model.PermJob {
	e := model.PermJob{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermJob")
	}
	return nil
}

// GetBy 岗位
func (p *rbacDaoImpl) GetPermJobBy(where string, v ...interface{}) *model.PermJob {
	e := model.PermJob{}
	err := p._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermJob")
	}
	return nil
}

// Count 岗位 by condition
func (p *rbacDaoImpl) CountPermJob(where string, v ...interface{}) (int, error) {
	return p._orm.Count(model.PermJob{}, where, v...)
}

// Select 岗位
func (p *rbacDaoImpl) SelectPermJob(where string, v ...interface{}) []*model.PermJob {
	list := make([]*model.PermJob, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermJob")
	}
	return list
}

// Save 岗位
func (p *rbacDaoImpl) SavePermJob(v *model.PermJob) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermJob")
	}
	return id, err
}

// Delete 岗位
func (p *rbacDaoImpl) DeletePermJob(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PermJob{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermJob")
	}
	return err
}

// Batch Delete 岗位
func (p *rbacDaoImpl) BatchDeletePermJob(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PermJob{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermJob")
	}
	return r, err
}

// Query paging data
func (p *rbacDaoImpl) PagingQueryPermJob(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(0) FROM perm_job WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM perm_job WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}
