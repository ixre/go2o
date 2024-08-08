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
		_ = o.Mapping(model.RbacDepart{}, "rbac_depart")
		_ = o.Mapping(model.RbacJob{}, "rbac_job")
		_ = o.Mapping(model.RbacUser{}, "rbac_user")
		_ = o.Mapping(model.RbacRole{}, "rbac_role")
		_ = o.Mapping(model.RbacRes{}, "rbac_res")
		_ = o.Mapping(model.RbacRoleDept{}, "rbac_role_dept")
		_ = o.Mapping(model.RbacRoleRes{}, "rbac_role_res")
		_ = o.Mapping(model.RbacUserRole{}, "rbac_user_role")

		rbacDaoImplMapped = true
	}
	return &rbacDaoImpl{
		_orm: o,
	}
}

func (p *rbacDaoImpl) GetRoleResList(roles []int) []*model.RbacRoleRes {
	where := fmt.Sprintf("role_id IN (%s)", util.JoinIntArray(roles, ","))
	return p.SelectPermRoleRes(where)
}

// Get 部门
func (p *rbacDaoImpl) GetDepart(primary interface{}) *model.RbacDepart {
	e := model.RbacDepart{}
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
func (p *rbacDaoImpl) GetDepartBy(where string, v ...interface{}) *model.RbacDepart {
	e := model.RbacDepart{}
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
	return p._orm.Count(model.RbacDepart{}, where, v...)
}

// Select 部门
func (p *rbacDaoImpl) SelectPermDept(where string, v ...interface{}) []*model.RbacDepart {
	list := make([]*model.RbacDepart, 0)
	err := p._orm.Select(&list, where+" 1=1 ORDER BY id ASC", v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return list
}

// Save 部门
func (p *rbacDaoImpl) SaveDepart(v *model.RbacDepart) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return id, err
}

// Delete 部门
func (p *rbacDaoImpl) DeleteDepart(primary interface{}) error {
	err := p._orm.DeleteByPk(model.RbacDepart{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return err
}

// Batch Delete 部门
func (p *rbacDaoImpl) BatchDeleteDepart(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.RbacDepart{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermDept")
	}
	return r, err
}

// Get 岗位
func (p *rbacDaoImpl) GetJob(primary interface{}) *model.RbacJob {
	e := model.RbacJob{}
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
func (p *rbacDaoImpl) GetJobBy(where string, v ...interface{}) *model.RbacJob {
	e := model.RbacJob{}
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
	return p._orm.Count(model.RbacJob{}, where, v...)
}

// Select 岗位
func (p *rbacDaoImpl) SelectPermJob(where string, v ...interface{}) []*model.RbacJob {
	list := make([]*model.RbacJob, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermJob")
	}
	return list
}

// Save 岗位
func (p *rbacDaoImpl) SaveJob(v *model.RbacJob) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermJob")
	}
	return id, err
}

// Delete 岗位
func (p *rbacDaoImpl) DeleteJob(primary interface{}) error {
	err := p._orm.DeleteByPk(model.RbacJob{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermJob")
	}
	return err
}

// Batch Delete 岗位
func (p *rbacDaoImpl) BatchDeleteJob(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.RbacJob{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermJob")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) QueryPagingJob(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM rbac_job WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT *,
			(SELECT name FROM rbac_depart WHERE id=dept_id) as dept_name
			 FROM rbac_job WHERE %s %s LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Printf("[ dao][ error]: %s (table:rbac_job) \n", err.Error())
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

// Get 系统用户
func (p *rbacDaoImpl) GetUser(primary interface{}) *model.RbacUser {
	e := model.RbacUser{}
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
func (p *rbacDaoImpl) GetUserBy(where string, v ...interface{}) *model.RbacUser {
	e := model.RbacUser{}
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
	return p._orm.Count(model.RbacUser{}, where, v...)
}

// Select 系统用户
func (p *rbacDaoImpl) SelectPermUser(where string, v ...interface{}) []*model.RbacUser {
	list := make([]*model.RbacUser, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return list
}

// Save 系统用户
func (p *rbacDaoImpl) SaveUser(v *model.RbacUser) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return id, err
}

// Delete 系统用户
func (p *rbacDaoImpl) DeleteUser(primary interface{}) error {
	err := p._orm.DeleteByPk(model.RbacUser{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return err
}

// Batch Delete 系统用户
func (p *rbacDaoImpl) BatchDeleteUser(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.RbacUser{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUser")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) QueryPagingPermUser(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM rbac_user WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM rbac_user WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:rbac_user)", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

// Get 角色
func (p *rbacDaoImpl) GetRole(primary interface{}) *model.RbacRole {
	e := model.RbacRole{}
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
func (p *rbacDaoImpl) GetRoleBy(where string, v ...interface{}) *model.RbacRole {
	e := model.RbacRole{}
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
	return p._orm.Count(model.RbacRole{}, where, v...)
}

// Select 角色
func (p *rbacDaoImpl) SelectPermRole(where string, v ...interface{}) []*model.RbacRole {
	list := make([]*model.RbacRole, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRole")
	}
	return list
}

// Save 角色
func (p *rbacDaoImpl) SavePermRole(v *model.RbacRole) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRole")
	}
	return id, err
}

// Delete 角色
func (p *rbacDaoImpl) DeletePermRole(primary interface{}) error {
	err := p._orm.DeleteByPk(model.RbacRole{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRole")
	}
	return err
}

// Batch Delete 角色
func (p *rbacDaoImpl) BatchDeletePermRole(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.RbacRole{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRole")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) QueryPagingPermRole(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM rbac_role WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM rbac_role WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:rbac_role)", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

// Get PermRes
func (p *rbacDaoImpl) GetRbacResource(primary interface{}) *model.RbacRes {
	e := model.RbacRes{}
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
func (p *rbacDaoImpl) GetRbacResourceBy(where string, v ...interface{}) *model.RbacRes {
	e := model.RbacRes{}
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
	return p._orm.Count(model.RbacRes{}, where, v...)
}

// Select PermRes
func (p *rbacDaoImpl) SelectPermRes(where string, v ...interface{}) []*model.RbacRes {
	list := make([]*model.RbacRes, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return list
}

func (p *rbacDaoImpl) GetMaxResourceSortNum(parentId int) int {
	var i int
	p._orm.Connector().ExecScalar(
		`SELECT MAX(sort_num) FROM rbac_res
 		  WHERE pid = $1`, &i, parentId)
	return i
}

// GetMaxResourceSortNum 获取最大的Key
func (p *rbacDaoImpl) GetMaxResouceKey(parentId int) string {
	var s string
	p._orm.Connector().ExecScalar(
		`SELECT MAX(res_key) FROM rbac_res
 		  WHERE pid = $1`, &s, parentId)
	return s
}

// Save PermRes
func (p *rbacDaoImpl) SaveRbacResource(v *model.RbacRes) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return id, err
}

// Delete PermRes
func (p *rbacDaoImpl) DeleteRbacResource(primary interface{}) error {
	err := p._orm.DeleteByPk(model.RbacRes{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return err
}

// Batch Delete PermRes
func (p *rbacDaoImpl) BatchDeleteRbacResource(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.RbacRes{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return r, err
}

// Get 用户角色关联
func (p *rbacDaoImpl) GetUserRole(primary interface{}) *model.RbacUserRole {
	e := model.RbacUserRole{}
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
func (p *rbacDaoImpl) GetUserRoleBy(where string, v ...interface{}) *model.RbacUserRole {
	e := model.RbacUserRole{}
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
	return p._orm.Count(model.RbacUserRole{}, where, v...)
}

// Select 用户角色关联
func (p *rbacDaoImpl) SelectPermUserRole(where string, v ...interface{}) []*model.RbacUserRole {
	list := make([]*model.RbacUserRole, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUserRole")
	}
	return list
}

// Save 用户角色关联
func (p *rbacDaoImpl) SaveUserRole(v *model.RbacUserRole) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUserRole")
	}
	return id, err
}

// Delete 用户角色关联
func (p *rbacDaoImpl) DeleteUserRole(primary interface{}) error {
	err := p._orm.DeleteByPk(model.RbacUserRole{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUserRole")
	}
	return err
}

// Batch Delete 用户角色关联
func (p *rbacDaoImpl) BatchDeleteUserRole(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.RbacUserRole{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermUserRole")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) QueryPagingPermUserRole(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM rbac_user_role WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM rbac_user_role WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:rbac_user_role)", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

// Get 角色菜单关联
func (p *rbacDaoImpl) GetRoleRes(primary interface{}) *model.RbacRoleRes {
	e := model.RbacRoleRes{}
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
func (p *rbacDaoImpl) GetRoleResBy(where string, v ...interface{}) *model.RbacRoleRes {
	e := model.RbacRoleRes{}
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
	return p._orm.Count(model.RbacRoleRes{}, where, v...)
}

// Select 角色菜单关联
func (p *rbacDaoImpl) SelectPermRoleRes(where string, v ...interface{}) []*model.RbacRoleRes {
	list := make([]*model.RbacRoleRes, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleRes")
	}
	return list
}

// Save 角色菜单关联
func (p *rbacDaoImpl) SavePermRoleRes(v *model.RbacRoleRes) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleRes")
	}
	return id, err
}

// Delete 角色菜单关联
func (p *rbacDaoImpl) DeletePermRoleRes(primary interface{}) error {
	err := p._orm.DeleteByPk(model.RbacRoleRes{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleRes")
	}
	return err
}

// Batch Delete 角色菜单关联
func (p *rbacDaoImpl) BatchDeletePermRoleRes(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.RbacRoleRes{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleRes")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) QueryPagingPermRoleRes(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM rbac_role_res WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM rbac_role_res WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:rbac_role_res)", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

// Get 角色部门关联
func (p *rbacDaoImpl) GetRoleDept(primary interface{}) *model.RbacRoleDept {
	e := model.RbacRoleDept{}
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
func (p *rbacDaoImpl) GetRoleDeptBy(where string, v ...interface{}) *model.RbacRoleDept {
	e := model.RbacRoleDept{}
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
	return p._orm.Count(model.RbacRoleDept{}, where, v...)
}

// Select 角色部门关联
func (p *rbacDaoImpl) SelectPermRoleDept(where string, v ...interface{}) []*model.RbacRoleDept {
	list := make([]*model.RbacRoleDept, 0)
	err := p._orm.Select(&list, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleDept")
	}
	return list
}

// Save 角色部门关联
func (p *rbacDaoImpl) SavePermRoleDept(v *model.RbacRoleDept) (int, error) {
	id, err := orm.Save(p._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleDept")
	}
	return id, err
}

// Delete 角色部门关联
func (p *rbacDaoImpl) DeletePermRoleDept(primary interface{}) error {
	err := p._orm.DeleteByPk(model.RbacRoleDept{}, primary)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleDept")
	}
	return err
}

// Batch Delete 角色部门关联
func (p *rbacDaoImpl) BatchDeletePermRoleDept(where string, v ...interface{}) (int64, error) {
	r, err := p._orm.Delete(model.RbacRoleDept{}, where, v...)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRoleDept")
	}
	return r, err
}

// Params paging data
func (p *rbacDaoImpl) QueryPagingPermRoleDept(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	s := fmt.Sprintf(`SELECT COUNT(1) FROM rbac_role_dept WHERE %s`, where)
	p._orm.Connector().ExecScalar(s, &total)
	if total > 0 {
		s = fmt.Sprintf(`SELECT * FROM rbac_role_dept WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(s, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Println(fmt.Sprintf("[ Orm][ Error]: %s (table:rbac_role_dept)", err.Error()))
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}

func (p *rbacDaoImpl) GetUserRoles(id int) []*model.RbacUserRole {
	return p.SelectPermUserRole("user_id = $1", id)
}

func (p *rbacDaoImpl) GetRoleResources(roles []int) []*model.RbacRes {
	where := fmt.Sprintf("role_id IN (%s)", util.JoinIntArray(roles, ","))
	var arr []*model.RbacRes
	err := p._orm.SelectByQuery(&arr, `SELECT * FROM rbac_res 
			INNER JOIN rbac_role_res ON rbac_role_res.res_id = rbac_res.id
			WHERE is_forbidden <> 1 AND is_enabled = 1 AND `+where)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ Orm][ Error]:", err.Error(), "; Entity:PermRes")
	}
	return arr
}

// QueryPagingLoginLog Query paging data
func (p *rbacDaoImpl) QueryPagingLoginLog(begin, end int, where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
		where = "1=1"
	}
	query := fmt.Sprintf(`SELECT COUNT(1) FROM rbac_login_log WHERE %s`, where)
	_ = p._orm.Connector().ExecScalar(query, &total)
	if total > 0 {
		query = fmt.Sprintf(`SELECT * FROM rbac_login_log WHERE %s %s
	        LIMIT $2 OFFSET $1`,
			where, orderBy)
		err := p._orm.Connector().Query(query, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, begin, end-begin)
		if err != nil {
			log.Printf("[ Orm][ Error]: %s (table:rbac_login_log)\n", err.Error())
		}
	} else {
		rows = make([]map[string]interface{}, 0)
	}
	return total, rows
}
