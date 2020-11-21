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

var _ dao.IPermUserDao = new(permUserDaoImpl)

type permUserDaoImpl struct {
	_orm orm.Orm
}

var permUserDaoImplMapped = false

// Create new PermUserDao
func NewPermUserDao(o orm.Orm) dao.IPermUserDao {
	if !permUserDaoImplMapped {
		_ = o.Mapping(model.PermUser{}, "perm_user")
		permUserDaoImplMapped = true
	}
	return &permUserDaoImpl{_orm: o}
}

// Get 系统用户
func (t *permUserDaoImpl) Get(primary interface{}) *model.PermUser {
	e := model.PermUser{}
	err := t._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return nil
}

// GetBy 系统用户
func (t *permUserDaoImpl) GetBy(where string, v ...interface{}) *model.PermUser {
	e := model.PermUser{}
	err := t._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return nil
}

// Count 系统用户 by condition
func (t *permUserDaoImpl) Count(where string, v ...interface{}) (int, error) {
	return t._orm.Count(model.PermUser{}, where, v...)
}

// Select 系统用户
func (t *permUserDaoImpl) Select(where string, v ...interface{}) []*model.PermUser {
	list := make([]*model.PermUser, 0)
	err := t._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return list
}

// Save 系统用户
func (t *permUserDaoImpl) Save(v *model.PermUser) (int, error) {
	id, err := orm.Save(t._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return id, err
}

// Delete 系统用户
func (t *permUserDaoImpl) Delete(primary interface{}) error {
	err := t._orm.DeleteByPk(model.PermUser{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return err
}

// Batch Delete 系统用户
func (t *permUserDaoImpl) BatchDelete(where string, v ...interface{}) (int64, error) {
	r, err := t._orm.Delete(model.PermUser{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return r, err
}

// Query paging data
func (t *permUserDaoImpl) PagingQuery(begin, end int, where, orderBy string) (num int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(0) FROM perm_user WHERE %s`, where)
	t._orm.Connector().ExecScalar(s, &num)
	if num > 0 {
		s = fmt.Sprintf(`SELECT * FROM perm_user WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		t._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return num, rows
}
