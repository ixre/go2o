package impl

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ixre/go2o/core/dao"
	"github.com/ixre/go2o/core/dao/model"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/util"
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
		_ = o.Mapping(model.PermUser{}, "perm_user")
		_ = o.Mapping(model.PermRole{}, "perm_role")
		_ = o.Mapping(model.PermRes{}, "perm_res")
		_ = o.Mapping(model.PermRoleDept{}, "perm_role_dept")
		_ = o.Mapping(model.PermRoleRes{}, "perm_role_res")
		_ = o.Mapping(model.PermUserRole{}, "perm_user_role")

		rbacDaoImplMapped = true
	}
	return &rbacDaoImpl{
		_orm: o,
	}
}

func (p *rbacDaoImpl) GetRoleResList(roleId int64) []int64 {
	// 绑定资源ID
	roles := p.SelectPermRoleRes("role_id=$1", roleId)
	arr := make([]int64, len(roles))
	for i, v := range roles {
		arr[i] = v.ResId
	}
	return arr
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
	err := p._orm.Select(&list, where+" 1=1 ORDER BY id ASC", v...)
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

// Params paging data
func (p *rbacDaoImpl) PagingQueryPermJob(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM perm_job WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM perm_job WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ dao][ error]: %s (table:perm_job) ", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

// Get 系统用户
func (p *rbacDaoImpl) GetPermUser(primary interface{}) *model.PermUser {
	e := model.PermUser{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return nil
}

// GetBy 系统用户
func (p *rbacDaoImpl) GetPermUserBy(where string, v ...interface{}) *model.PermUser {
	e := model.PermUser{}
	err := p._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return nil
}

// Count 系统用户 by condition
func (p *rbacDaoImpl) CountPermUser(where string, v ...interface{}) (int, error) {
	return p._orm.Count(model.PermUser{}, where, v...)
}

// Select 系统用户
func (p *rbacDaoImpl) SelectPermUser(where string, v ...interface{}) []*model.PermUser {
	list := make([]*model.PermUser, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return list
}

// Save 系统用户
func (p *rbacDaoImpl) SavePermUser(v *model.PermUser) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return id, err
}

// Delete 系统用户
func (p *rbacDaoImpl) DeletePermUser(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PermUser{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return err
}

// Batch Delete 系统用户
func (p *rbacDaoImpl) BatchDeletePermUser(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PermUser{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) PagingQueryPermUser(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM perm_user WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM perm_user WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:perm_user)", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

// Get 角色
func (p *rbacDaoImpl) GetPermRole(primary interface{}) *model.PermRole {
	e := model.PermRole{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRole")
	}
	return nil
}

// GetBy 角色
func (p *rbacDaoImpl) GetPermRoleBy(where string, v ...interface{}) *model.PermRole {
	e := model.PermRole{}
	err := p._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRole")
	}
	return nil
}

// Count 角色 by condition
func (p *rbacDaoImpl) CountPermRole(where string, v ...interface{}) (int, error) {
	return p._orm.Count(model.PermRole{}, where, v...)
}

// Select 角色
func (p *rbacDaoImpl) SelectPermRole(where string, v ...interface{}) []*model.PermRole {
	list := make([]*model.PermRole, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRole")
	}
	return list
}

// Save 角色
func (p *rbacDaoImpl) SavePermRole(v *model.PermRole) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRole")
	}
	return id, err
}

// Delete 角色
func (p *rbacDaoImpl) DeletePermRole(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PermRole{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRole")
	}
	return err
}

// Batch Delete 角色
func (p *rbacDaoImpl) BatchDeletePermRole(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PermRole{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRole")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) PagingQueryPermRole(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM perm_role WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM perm_role WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:perm_role)", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

// Get PermRes
func (p *rbacDaoImpl) GetPermRes(primary interface{}) *model.PermRes {
	e := model.PermRes{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return nil
}

// GetBy PermRes
func (p *rbacDaoImpl) GetPermResBy(where string, v ...interface{}) *model.PermRes {
	e := model.PermRes{}
	err := p._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return nil
}

// Count PermRes by condition
func (p *rbacDaoImpl) CountPermRes(where string, v ...interface{}) (int, error) {
	return p._orm.Count(model.PermRes{}, where, v...)
}

// Select PermRes
func (p *rbacDaoImpl) SelectPermRes(where string, v ...interface{}) []*model.PermRes {
	list := make([]*model.PermRes, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return list
}

func (p *rbacDaoImpl) GetMaxResourceSortNum(parentId int) int {
	var i int
	p._orm.Connector().ExecScalar(
		`SELECT MAX(sort_num) FROM perm_res
 		  WHERE pid = $1`, &i, parentId)
	return i
}

// Save PermRes
func (p *rbacDaoImpl) SavePermRes(v *model.PermRes) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return id, err
}

// Delete PermRes
func (p *rbacDaoImpl) DeletePermRes(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PermRes{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return err
}

// Batch Delete PermRes
func (p *rbacDaoImpl) BatchDeletePermRes(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PermRes{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return r, err
}

// Get 用户角色关联
func (p *rbacDaoImpl) GetPermUserRole(primary interface{}) *model.PermUserRole {
	e := model.PermUserRole{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUserRole")
	}
	return nil
}

// GetBy 用户角色关联
func (p *rbacDaoImpl) GetPermUserRoleBy(where string, v ...interface{}) *model.PermUserRole {
	e := model.PermUserRole{}
	err := p._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUserRole")
	}
	return nil
}

// Count 用户角色关联 by condition
func (p *rbacDaoImpl) CountPermUserRole(where string, v ...interface{}) (int, error) {
	return p._orm.Count(model.PermUserRole{}, where, v...)
}

// Select 用户角色关联
func (p *rbacDaoImpl) SelectPermUserRole(where string, v ...interface{}) []*model.PermUserRole {
	list := make([]*model.PermUserRole, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUserRole")
	}
	return list
}

// Save 用户角色关联
func (p *rbacDaoImpl) SavePermUserRole(v *model.PermUserRole) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUserRole")
	}
	return id, err
}

// Delete 用户角色关联
func (p *rbacDaoImpl) DeletePermUserRole(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PermUserRole{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUserRole")
	}
	return err
}

// Batch Delete 用户角色关联
func (p *rbacDaoImpl) BatchDeletePermUserRole(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PermUserRole{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUserRole")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) PagingQueryPermUserRole(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM perm_user_role WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM perm_user_role WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:perm_user_role)", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

// Get 角色菜单关联
func (p *rbacDaoImpl) GetPermRoleRes(primary interface{}) *model.PermRoleRes {
	e := model.PermRoleRes{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleRes")
	}
	return nil
}

// GetBy 角色菜单关联
func (p *rbacDaoImpl) GetPermRoleResBy(where string, v ...interface{}) *model.PermRoleRes {
	e := model.PermRoleRes{}
	err := p._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleRes")
	}
	return nil
}

// Count 角色菜单关联 by condition
func (p *rbacDaoImpl) CountPermRoleRes(where string, v ...interface{}) (int, error) {
	return p._orm.Count(model.PermRoleRes{}, where, v...)
}

// Select 角色菜单关联
func (p *rbacDaoImpl) SelectPermRoleRes(where string, v ...interface{}) []*model.PermRoleRes {
	list := make([]*model.PermRoleRes, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleRes")
	}
	return list
}

// Save 角色菜单关联
func (p *rbacDaoImpl) SavePermRoleRes(v *model.PermRoleRes) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleRes")
	}
	return id, err
}

// Delete 角色菜单关联
func (p *rbacDaoImpl) DeletePermRoleRes(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PermRoleRes{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleRes")
	}
	return err
}

// Batch Delete 角色菜单关联
func (p *rbacDaoImpl) BatchDeletePermRoleRes(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PermRoleRes{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleRes")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) PagingQueryPermRoleRes(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM perm_role_res WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM perm_role_res WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:perm_role_res)", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

// Get 角色部门关联
func (p *rbacDaoImpl) GetPermRoleDept(primary interface{}) *model.PermRoleDept {
	e := model.PermRoleDept{}
	err := p._orm.Get(primary, &e)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleDept")
	}
	return nil
}

// GetBy 角色部门关联
func (p *rbacDaoImpl) GetPermRoleDeptBy(where string, v ...interface{}) *model.PermRoleDept {
	e := model.PermRoleDept{}
	err := p._orm.GetBy(&e, where, v...)
	if err == nil {
		return &e
	}
	if err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleDept")
	}
	return nil
}

// Count 角色部门关联 by condition
func (p *rbacDaoImpl) CountPermRoleDept(where string, v ...interface{}) (int, error) {
	return p._orm.Count(model.PermRoleDept{}, where, v...)
}

// Select 角色部门关联
func (p *rbacDaoImpl) SelectPermRoleDept(where string, v ...interface{}) []*model.PermRoleDept {
	list := make([]*model.PermRoleDept, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleDept")
	}
	return list
}

// Save 角色部门关联
func (p *rbacDaoImpl) SavePermRoleDept(v *model.PermRoleDept) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleDept")
	}
	return id, err
}

// Delete 角色部门关联
func (p *rbacDaoImpl) DeletePermRoleDept(primary interface{}) error {
	err := p._orm.DeleteByPk(model.PermRoleDept{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleDept")
	}
	return err
}

// Batch Delete 角色部门关联
func (p *rbacDaoImpl) BatchDeletePermRoleDept(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.PermRoleDept{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleDept")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) PagingQueryPermRoleDept(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM perm_role_dept WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM perm_role_dept WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:perm_role_dept)", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

func (p *rbacDaoImpl) GetUserRoles(id int64) []*model.PermUserRole {
	return p.SelectPermUserRole("user_id = $1", id)
}

func (p *rbacDaoImpl) GetRoleResources(roles []int) []*model.PermRes {
	where := fmt.Sprintf("role_id IN (%s)", util.JoinIntArray(roles, ","))
	var arr []*model.PermRes
	err := p._orm.SelectByQuery(&arr, `SELECT * FROM perm_res 
			INNER JOIN perm_role_res ON perm_role_res.res_id = perm_res.id
			WHERE `+where)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return arr
}
