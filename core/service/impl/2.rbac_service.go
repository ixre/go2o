package impl

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : perm_dept_service.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020/12/02 13:02:38
 * description :
 * history :
 */

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/ixre/go2o/core/dao"
	"github.com/ixre/go2o/core/dao/impl"
	"github.com/ixre/go2o/core/dao/model"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types"
	"github.com/ixre/gof/types/typeconv"
	"github.com/ixre/gof/util"
	"strings"
	"time"
)

var _ proto.RbacServiceServer = new(rbacServiceImpl)

// 基于角色的权限服务
type rbacServiceImpl struct {
	dao          dao.IRbacDao
	registryRepo registry.IRegistryRepo
	s            storage.Interface
	serviceUtil
}

func NewRbacService(s storage.Interface, o orm.Orm, registryRepo registry.IRegistryRepo) *rbacServiceImpl {
	return &rbacServiceImpl{
		s:            s,
		registryRepo: registryRepo,
		dao:          impl.NewRbacDao(o),
	}
}

func (p *rbacServiceImpl) UserLogin(_ context.Context, r *proto.RbacLoginRequest) (*proto.RbacLoginResponse, error) {
	if len(r.Pwd) != 32 {
		return &proto.RbacLoginResponse{
			ErrCode: 1,
			ErrMsg:  "密码无效",
		}, nil
	}
	// 超级管理员登录
	if r.User == "master" {
		superPwd, _ := p.registryRepo.GetValue(registry.SysSuperLoginToken)
		encPwd := domain.Sha1Pwd(r.User+r.Pwd, "")
		if superPwd != encPwd {
			return &proto.RbacLoginResponse{
				ErrCode: 3,
				ErrMsg:  "密码不正确",
			}, nil
		}
		dst := &proto.RbacLoginResponse{
			UserId:      0,
			Permissions: []string{"master", "admin"},
		}
		return p.withAccessToken(0, "master", dst, r.Expires)
	}
	// 普通系统用户登录
	usr := p.dao.GetPermUserBy("usr=$1", r.User)
	if usr == nil {
		return &proto.RbacLoginResponse{
			ErrCode: 2,
			ErrMsg:  "用户不存在",
		}, nil
	}
	decPwd := crypto.Sha1([]byte(r.Pwd + usr.Salt))
	if usr.Pwd != decPwd {
		return &proto.RbacLoginResponse{
			ErrCode: 3,
			ErrMsg:  "密码不正确",
		}, nil
	}
	if usr.Enabled != 1 {
		return &proto.RbacLoginResponse{
			ErrCode: 4,
			ErrMsg:  "用户已停用",
		}, nil
	}
	dst := &proto.RbacLoginResponse{
		UserId: usr.Id,
	}
	dst.Roles, dst.Permissions = p.getUserRolesPerm(usr.Id)
	return p.withAccessToken(usr.Id, usr.Usr, dst, r.Expires)
}

// 返回带有令牌的结果
func (p *rbacServiceImpl) withAccessToken(userId int64, userName string,
	dst *proto.RbacLoginResponse, expires int32) (*proto.RbacLoginResponse, error) {
	accessToken, err := p.createAccessToken(userId, userName,
		strings.Join(dst.Permissions, ","), expires)
	dst.AccessToken = accessToken
	if err != nil {
		dst.ErrCode = 2
		dst.ErrMsg = err.Error()
	}
	return dst, nil
}

// 创建令牌
func (p *rbacServiceImpl) createAccessToken(userId int64, userName string, perm string, exp int32) (string, error) {
	if exp <= 0 {
		exp = int32((time.Hour * 24 * 365).Seconds())
	}
	var claims = jwt.MapClaims{
		"exp":    time.Now().Add(time.Second * time.Duration(exp)).Unix(),
		"aud":    userId,
		"iss":    "Go2o",
		"sub":    "Go2o-RBAC-Token",
		"name":   userName,
		"x-perm": perm,
	}
	key, _ := p.registryRepo.GetValue(registry.SysJWTSecret)
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	return token.SignedString([]byte(key))
}

// 获取JWT密钥
func (p *rbacServiceImpl) GetJwtToken(_ context.Context, empty *proto.Empty) (*proto.String, error) {
	key, _ := p.registryRepo.GetValue(registry.SysJWTSecret)
	return &proto.String{Value: key}, nil
}

// 获取用户的角色和权限
func (p *rbacServiceImpl) getUserRolesPerm(userId int64) ([]int64, []string) {
	userRoles := p.dao.GetUserRoles(userId)
	// 绑定角色ID
	roles := make([]int64, len(userRoles))
	roleList := make([]int, len(userRoles))
	for i, v := range userRoles {
		roles[i] = v.RoleId
		roleList[i] = int(v.RoleId)
	}
	// 获取角色的权限
	rolesList := p.dao.SelectPermRole(
		fmt.Sprintf("id IN (%_s)", util.JoinIntArray(roleList, ",")))
	permissions := make([]string, len(roles))
	for i, v := range rolesList {
		permissions[i] = v.Permission
	}
	return roles, permissions
}

// 移动资源顺序
func (p *rbacServiceImpl) MoveResOrdinal(_ context.Context, r *proto.MoveResOrdinalRequest) (*proto.Result, error) {
	res := p.dao.GetPermRes(r.ResourceId)
	if res == nil {
		return p.error(errors.New("no such data")), nil
	}
	// 获取交换的对象
	var swapRes *model.PermRes
	if r.Direction == 0 { // 向上移,获取上一个
		swapRes = p.dao.GetPermResBy(
			`sort_num < $1 AND pid = $2 AND depth=$3 ORDER BY sort_num DESC`,
			res.SortNum, res.Pid, res.Depth)
	} else {
		swapRes = p.dao.GetPermResBy(
			`sort_num > $1 AND pid = $2 AND depth=$3 ORDER BY sort_num ASC`,
			res.SortNum, res.Pid, res.Depth)
	}
	// 交换顺序
	if swapRes != nil {
		sortNum := swapRes.SortNum
		swapRes.SortNum = res.SortNum
		res.SortNum = sortNum
		p.dao.SavePermRes(res)
		p.dao.SavePermRes(swapRes)
	}
	return p.success(nil), nil
}

func (p *rbacServiceImpl) GetUserResource(_ context.Context, r *proto.GetUserResRequest) (*proto.RbacUserResourceResponse, error) {
	dst := &proto.RbacUserResourceResponse{}
	var resList []*model.PermRes
	if r.UserId <= 0 { // 管理员
		dst.Roles = []int64{}
		dst.Permissions = []string{"master", "admin"}
		resList = p.dao.SelectPermRes("")
		// 获取管理员
		for _, v := range resList {
			v.Permission = "master,admin"
		}
	} else {
		dst.Roles, dst.Permissions = p.getUserRolesPerm(r.UserId)
		usr := p.dao.GetPermUser(r.UserId)
		if usr == nil {
			return nil, nil
		}
		roleList := make([]int, len(dst.Roles))
		for i, v := range dst.Roles {
			roleList[i] = int(v)
		}
		resList = p.dao.GetRoleResources(roleList)
	}
	root := proto.SUserRes{}
	var f func(root *proto.SUserRes, arr []*model.PermRes)
	f = func(root *proto.SUserRes, arr []*model.PermRes) {
		root.Children = []*proto.SUserRes{}
		for _, v := range arr {
			if r.OnlyMenu && (v.ResType != 0 && v.ResType != 2) {
				continue // 只显示菜单
			}
			if v.Pid == root.Id {
				c := &proto.SUserRes{
					Id:            v.Id,
					Key:           v.Key,
					Name:          v.Name,
					ResType:       int32(v.ResType),
					Path:          v.Path,
					Icon:          v.Icon,
					Permission:    v.Permission,
					SortNum:       int32(v.SortNum),
					IsHidden:      v.IsHidden == 1,
					ComponentName: v.ComponentName,
				}
				c.Children = make([]*proto.SUserRes, 0)
				root.Children = append(root.Children, c)
				f(c, arr)
			}
		}
	}
	f(&root, resList)
	dst.Resources = root.Children
	return dst, nil
}

func walkDepartTree(node *proto.RbacTree, nodeList []*model.PermDept) {
	node.Children = []*proto.RbacTree{}
	for _, v := range nodeList {
		if v.Pid == node.Id {
			v := &proto.RbacTree{
				Id:       v.Id,
				Label:    v.Name,
				Children: make([]*proto.RbacTree, 0),
			}
			node.Children = append(node.Children, v)
			walkDepartTree(v, nodeList)
		}
	}
}

// 部门树形数据
func (p *rbacServiceImpl) DepartTree(_ context.Context, empty *proto.Empty) (*proto.RbacTree, error) {
	root := &proto.RbacTree{
		Id:       0,
		Label:    "根节点",
		Children: make([]*proto.RbacTree, 0),
	}
	list := p.dao.SelectPermDept("")
	walkDepartTree(root, list)
	return root, nil
}

// 保存部门
func (p *rbacServiceImpl) SavePermDept(_ context.Context, r *proto.SavePermDeptRequest) (*proto.SavePermDeptResponse, error) {
	var dst *model.PermDept
	if r.Id > 0 {
		dst = p.dao.GetPermDept(r.Id)
	} else {
		dst = &model.PermDept{}
		dst.CreateTime = time.Now().Unix()
	}

	dst.Name = r.Name
	dst.Pid = r.Pid
	dst.Enabled = int16(r.Enabled)

	id, err := p.dao.SavePermDept(dst)
	ret := &proto.SavePermDeptResponse{
		Id: int64(id),
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

// 获取部门
func (p *rbacServiceImpl) GetPermDept(_ context.Context, id *proto.PermDeptId) (*proto.SPermDept, error) {
	v := p.dao.GetPermDept(id.Value)
	if v == nil {
		return nil, nil
	}
	return p.parsePermDept(v), nil
}

func (p *rbacServiceImpl) DeletePermDept(_ context.Context, id *proto.PermDeptId) (*proto.Result, error) {
	err := p.dao.DeletePermDept(id.Value)
	return p.error(err), nil
}

func (p *rbacServiceImpl) parsePermDept(v *model.PermDept) *proto.SPermDept {
	return &proto.SPermDept{
		Id:         v.Id,
		Name:       v.Name,
		Pid:        v.Pid,
		Enabled:    int32(v.Enabled),
		CreateTime: v.CreateTime,
	}
}

// 保存岗位
func (p *rbacServiceImpl) SavePermJob(_ context.Context, r *proto.SavePermJobRequest) (*proto.SavePermJobResponse, error) {
	var dst *model.PermJob
	if r.Id > 0 {
		dst = p.dao.GetPermJob(r.Id)
	} else {
		dst = &model.PermJob{}
		dst.CreateTime = time.Now().Unix()
	}

	dst.Name = r.Name
	dst.Enabled = int16(r.Enabled)
	dst.Sort = int(r.Sort)
	dst.DeptId = r.DeptId

	id, err := p.dao.SavePermJob(dst)
	ret := &proto.SavePermJobResponse{
		Id: int64(id),
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

func (p *rbacServiceImpl) parsePermJob(v *model.PermJob) *proto.SPermJob {
	return &proto.SPermJob{
		Id:         v.Id,
		Name:       v.Name,
		Enabled:    int32(v.Enabled),
		Sort:       int32(v.Sort),
		DeptId:     v.DeptId,
		CreateTime: v.CreateTime,
	}
}

// 获取岗位
func (p *rbacServiceImpl) GetPermJob(_ context.Context, id *proto.PermJobId) (*proto.SPermJob, error) {
	v := p.dao.GetPermJob(id.Value)
	if v == nil {
		return nil, nil
	}
	return p.parsePermJob(v), nil
}

// 获取岗位列表
func (p *rbacServiceImpl) QueryPermJobList(_ context.Context, r *proto.QueryPermJobRequest) (*proto.QueryPermJobResponse, error) {
	where := ""
	if r.DepartId > 0 {
		arr := p.walkDepartArray(int(r.DepartId))
		if len(where) > 0 {
			where += " AND "
		}
		where += " dept_id IN (" + util.JoinIntArray(arr, ",") + ")"
	}
	arr := p.dao.SelectPermJob(where)
	ret := &proto.QueryPermJobResponse{
		List: make([]*proto.SPermJob, len(arr)),
	}
	for i, v := range arr {
		ret.List[i] = p.parsePermJob(v)
	}
	return ret, nil
}

func (p *rbacServiceImpl) DeletePermJob(_ context.Context, id *proto.PermJobId) (*proto.Result, error) {
	err := p.dao.DeletePermJob(id.Value)
	return p.error(err), nil
}

func (p *rbacServiceImpl) PagingPermJob(_ context.Context, r *proto.PermJobPagingRequest) (*proto.PermJobPagingResponse, error) {
	total, rows := p.dao.PagingQueryPermJob(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.PermJobPagingResponse{
		Total: int64(total),
		Value: make([]*proto.PagingPermJob, len(rows)),
	}
	for i, v := range rows {
		ret.Value[i] = &proto.PagingPermJob{
			Id:         int64(typeconv.MustInt(v["id"])),
			Name:       typeconv.Stringify(v["name"]),
			Enabled:    int32(typeconv.MustInt(v["enabled"])),
			Sort:       int32(typeconv.MustInt(v["sort"])),
			DeptId:     int64(typeconv.MustInt(v["dept_id"])),
			CreateTime: int64(typeconv.MustInt(v["create_time"])),
		}
	}
	return ret, nil
}

// 保存系统用户
func (p *rbacServiceImpl) SavePermUser(_ context.Context, r *proto.SavePermUserRequest) (*proto.SavePermUserResponse, error) {
	var dst *model.PermUser
	if r.Id > 0 {
		dst = p.dao.GetPermUser(r.Id)
		if dst == nil {
			return &proto.SavePermUserResponse{
				ErrCode: 2,
				ErrMsg:  "no such record",
			}, nil
		}
	} else {
		dst = &model.PermUser{}
		dst.Salt = util.RandString(4)
		dst.CreateTime = time.Now().Unix()
	}
	if l := len(r.Pwd); l > 0 {
		if l != 32 {
			return &proto.SavePermUserResponse{
				ErrCode: 1,
				ErrMsg:  "非32位md5密码",
			}, nil
		}
		dst.Pwd = crypto.Sha1([]byte(r.Pwd + dst.Salt))
	}
	dst.Flag = int(r.Flag)
	dst.Avatar = r.Avatar
	dst.NickName = r.NickName
	dst.Gender = r.Gender
	dst.Email = r.Email
	dst.Phone = r.Phone
	dst.DeptId = r.DeptId
	dst.JobId = r.JobId
	dst.Enabled = int16(r.Enabled)
	dst.LastLogin = r.LastLogin
	id, err := p.dao.SavePermUser(dst)
	ret := &proto.SavePermUserResponse{
		Id: int64(id),
	}
	if err == nil {
		err = p.updateUserRoles(int64(id), r.Roles)
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

func (p *rbacServiceImpl) parsePermUser(v *model.PermUser) *proto.SPermUser {
	return &proto.SPermUser{
		Id:         v.Id,
		Usr:        v.Usr,
		Pwd:        v.Pwd,
		Flag:       int32(v.Flag),
		Avatar:     v.Avatar,
		NickName:   v.NickName,
		Gender:     v.Gender,
		Email:      v.Email,
		Phone:      v.Phone,
		DeptId:     v.DeptId,
		JobId:      v.JobId,
		Enabled:    int32(v.Enabled),
		LastLogin:  v.LastLogin,
		CreateTime: v.CreateTime,
	}
}

// 获取系统用户
func (p *rbacServiceImpl) GetPermUser(_ context.Context, id *proto.PermUserId) (*proto.SPermUser, error) {
	v := p.dao.GetPermUser(id.Value)
	if v == nil {
		return nil, nil
	}
	dst := p.parsePermUser(v)
	dst.Roles, dst.Permissions = p.getUserRolesPerm(v.Id)
	return dst, nil
}

func (p *rbacServiceImpl) DeletePermUser(_ context.Context, id *proto.PermUserId) (*proto.Result, error) {
	err := p.dao.DeletePermUser(id.Value)
	return p.error(err), nil
}

func (p *rbacServiceImpl) walkDepartArray(pid int) []int {
	var arr = make([]int, 0)
	if pid > 0 {
		var f func(pid int, arr *[]int)
		f = func(pid int, arr *[]int) {
			*arr = append(*arr, pid)
			for _, v := range p.dao.SelectPermDept("pid = $1", pid) {
				f(int(v.Id), arr)
			}
		}
		f(pid, &arr)
	}
	return arr
}

func (p *rbacServiceImpl) PagingPermUser(_ context.Context, r *proto.PermUserPagingRequest) (*proto.PermUserPagingResponse, error) {
	if r.DepartId > 0 {
		arr := p.walkDepartArray(int(r.DepartId))
		if len(r.Params.Where) > 0 {
			r.Params.Where += " AND "
		}
		r.Params.Where += " dept_id IN (" + util.JoinIntArray(arr, ",") + ")"
	}
	total, rows := p.dao.PagingQueryPermUser(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.PermUserPagingResponse{
		Total: int64(total),
		Value: make([]*proto.PagingPermUser, len(rows)),
	}
	for i, v := range rows {
		ret.Value[i] = &proto.PagingPermUser{
			Id:         int64(typeconv.MustInt(v["id"])),
			Usr:        typeconv.Stringify(v["usr"]),
			Pwd:        typeconv.Stringify(v["pwd"]),
			Flag:       int32(typeconv.MustInt(v["flag"])),
			Avatar:     typeconv.Stringify(v["avatar"]),
			NickName:   typeconv.Stringify(v["nick_name"]),
			Gender:     typeconv.Stringify(v["gender"]),
			Email:      typeconv.Stringify(v["email"]),
			Phone:      typeconv.Stringify(v["phone"]),
			DeptId:     int64(typeconv.MustInt(v["dept_id"])),
			JobId:      int64(typeconv.MustInt(v["job_id"])),
			Enabled:    int32(typeconv.MustInt(v["enabled"])),
			LastLogin:  int64(typeconv.MustInt(v["last_login"])),
			CreateTime: int64(typeconv.MustInt(v["create_time"])),
		}
	}
	return ret, nil
}

// 保存角色
func (p *rbacServiceImpl) SavePermRole(_ context.Context, r *proto.SavePermRoleRequest) (*proto.SavePermRoleResponse, error) {
	var dst *model.PermRole
	if r.Id > 0 {
		dst = p.dao.GetPermRole(r.Id)
	} else {
		dst = &model.PermRole{}
		dst.CreateTime = time.Now().Unix()
	}

	dst.Name = r.Name
	dst.Level = int(r.Level)
	dst.DataScope = r.DataScope
	dst.Permission = r.Permission
	dst.Remark = r.Remark

	id, err := p.dao.SavePermRole(dst)
	ret := &proto.SavePermRoleResponse{
		Id: int64(id),
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

// 更新角色资源
func (p *rbacServiceImpl) UpdateRoleResource(_ context.Context, r *proto.UpdateRoleResRequest) (*proto.Result, error) {
	dataList := p.dao.SelectPermRoleRes("role_id=$1", r.RoleId)
	old := make([]int, len(dataList))
	arr := typeconv.Int64Array(r.Resources)
	mp := make(map[int]*model.PermRoleRes, 0)
	// 旧数组
	for i, v := range dataList {
		old[i] = int(v.ResId)
		mp[int(v.ResId)] = v
	}
	_, deleted := util.IntArrayDiff(old, arr, func(v int, add bool) {
		if add {
			p.dao.SavePermRoleRes(&model.PermRoleRes{
				ResId:  int64(v),
				RoleId: r.RoleId,
			})
		}
	})
	if len(deleted) > 0 {
		p.dao.BatchDeletePermRoleRes(
			fmt.Sprintf("role_id = %d AND res_id IN (%_s)",
				r.RoleId, util.JoinIntArray(deleted, ",")))
	}
	return p.error(nil), nil
}

func (p *rbacServiceImpl) parsePermRole(v *model.PermRole) *proto.SPermRole {
	return &proto.SPermRole{
		Id:         v.Id,
		Name:       v.Name,
		Level:      int32(v.Level),
		DataScope:  v.DataScope,
		Permission: v.Permission,
		Remark:     v.Remark,
		CreateTime: v.CreateTime,
	}
}

// 获取角色
func (p *rbacServiceImpl) GetPermRole(_ context.Context, id *proto.PermRoleId) (*proto.SPermRole, error) {
	v := p.dao.GetPermRole(id.Value)
	if v == nil {
		return nil, nil
	}
	dst := p.parsePermRole(v)
	// 绑定资源ID
	dst.RelatedResIdList = p.dao.GetRoleResList(v.Id)
	return dst, nil
}

// 获取角色列表
func (p *rbacServiceImpl) QueryPermRoleList(_ context.Context, r *proto.QueryPermRoleRequest) (*proto.QueryPermRoleResponse, error) {
	arr := p.dao.SelectPermRole("")
	ret := &proto.QueryPermRoleResponse{
		List: make([]*proto.SPermRole, len(arr)),
	}
	for i, v := range arr {
		ret.List[i] = p.parsePermRole(v)
	}
	return ret, nil
}

func (p *rbacServiceImpl) DeletePermRole(_ context.Context, id *proto.PermRoleId) (*proto.Result, error) {
	err := p.dao.DeletePermRole(id.Value)
	return p.error(err), nil
}

func (p *rbacServiceImpl) PagingPermRole(_ context.Context, r *proto.PermRolePagingRequest) (*proto.PermRolePagingResponse, error) {
	total, rows := p.dao.PagingQueryPermRole(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.PermRolePagingResponse{
		Total: int64(total),
		Value: make([]*proto.PagingPermRole, len(rows)),
	}
	for i, v := range rows {
		ret.Value[i] = &proto.PagingPermRole{
			Id:         int64(typeconv.MustInt(v["id"])),
			Name:       typeconv.Stringify(v["name"]),
			Level:      int32(typeconv.MustInt(v["level"])),
			DataScope:  typeconv.Stringify(v["data_scope"]),
			Permission: typeconv.Stringify(v["permission"]),
			Remark:     typeconv.Stringify(v["remark"]),
			CreateTime: int64(typeconv.MustInt(v["create_time"])),
		}
	}
	return ret, nil
}

// 保存PermRes
func (p *rbacServiceImpl) SavePermRes(_ context.Context, r *proto.SavePermResRequest) (*proto.SavePermResResponse, error) {
	var dst *model.PermRes
	if r.Id > 0 {
		if dst = p.dao.GetPermRes(r.Id); dst == nil {
			return &proto.SavePermResResponse{
				ErrCode: 2,
				ErrMsg:  "no such data",
			}, nil
		}
	} else {
		dst = &model.PermRes{}
		dst.CreateTime = time.Now().Unix()
		dst.Depth = 0
		// 如果首次没有填写ResKey, 则默认通过Path生成
		if r.Key == "" {
			r.Key = strings.Replace(r.Path, "/", ":", -1)
		}
	}
	if r.Pid > 0 {
		// 检测上级是否为自己
		if dst.Id == r.Pid {
			return &proto.SavePermResResponse{
				ErrCode: 2,
				ErrMsg:  "不能将自己指定为上级资源",
			}, nil
		}
		// 检测上级是否为下级
		if dst.Id > 0 {
			var parent *model.PermRes = p.dao.GetPermRes(r.Pid)
			for parent != nil && parent.Pid > 0 {
				parent = p.dao.GetPermRes(parent.Pid)
				if parent != nil && parent.Id == r.Id {
					return &proto.SavePermResResponse{
						ErrCode: 2,
						ErrMsg:  "不能选择下级作为上级资源",
					}, nil
				}
			}
		}
		// 限制下级资源路径不能以'/'开头,以避免无法找到资源的情况
		if len(r.Path) > 0 && r.Path[0] == '/' {
			return &proto.SavePermResResponse{
				ErrCode: 3,
				ErrMsg:  "该资源(包含上级资源)路径不能以'/'开头",
			}, nil
		}
	}

	// 上级是否改变
	var parentChanged = dst.Pid != r.Pid || (r.Pid != 0 && dst.Depth == 0)
	dst.Name = r.Name
	dst.ResType = int16(r.ResType)
	dst.Pid = r.Pid
	dst.Key = r.Key
	dst.Path = r.Path
	dst.Icon = r.Icon
	dst.Permission = r.Permission
	dst.SortNum = int(r.SortNum)
	dst.IsExternal = int16(types.ElseInt(r.IsExternal, 1, 0))
	dst.IsHidden = int16(types.ElseInt(r.IsHidden, 1, 0))
	dst.ComponentName = r.ComponentName
	dst.Cache = r.Cache
	// 如果未设置排列序号,或者更改了上级,则需系统自动编号
	if dst.SortNum <= 0 || parentChanged {
		dst.SortNum = p.dao.GetMaxResourceSortNum(int(dst.Pid)) + 1
	}
	id, err := p.dao.SavePermRes(dst)
	ret := &proto.SavePermResResponse{
		Id: int64(id),
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	} else {
		if parentChanged {
			depth := p.getResDepth(dst.Pid)
			p.updateResDepth(dst, int16(depth))
		}
	}
	return ret, nil
}

func (p *rbacServiceImpl) parsePermRes(v *model.PermRes) *proto.SPermRes {
	return &proto.SPermRes{
		Id:            v.Id,
		Name:          v.Name,
		ResType:       int32(v.ResType),
		Pid:           v.Pid,
		Key:           v.Key,
		Path:          v.Path,
		Icon:          v.Icon,
		Permission:    v.Permission,
		SortNum:       int32(v.SortNum),
		IsExternal:    v.IsExternal == 1,
		IsHidden:      v.IsHidden == 1,
		CreateTime:    v.CreateTime,
		ComponentName: v.ComponentName,
		Cache:         v.Cache,
	}
}

// 获取PermRes
func (p *rbacServiceImpl) GetPermRes(_ context.Context, id *proto.PermResId) (*proto.SPermRes, error) {
	v := p.dao.GetPermRes(id.Value)
	if v == nil {
		return nil, nil
	}
	return p.parsePermRes(v), nil
}

func (p *rbacServiceImpl) walkPermRes(root *proto.SPermRes, arr []*model.PermRes) {
	root.Children = []*proto.SPermRes{}
	for _, v := range arr {
		if v.Pid == root.Id {
			c := p.parsePermRes(v)
			c.Children = make([]*proto.SPermRes, 0)
			root.Children = append(root.Children, c)
			p.walkPermRes(c, arr)
		}
	}
}

// 获取PermRes列表
func (p *rbacServiceImpl) QueryResList(_ context.Context, r *proto.QueryPermResRequest) (*proto.QueryPermResResponse, error) {
	var where string = "1=1"
	if r.Keyword != "" {
		where += " AND name LIKE '%" + r.Keyword + "%'"
	}
	if r.OnlyMenu {
		where += " AND res_type IN(0,2)"
	}
	//todo: 搜索结果不为pid
	arr := p.dao.SelectPermRes(where + " ORDER BY sort_num ASC")
	root := proto.SPermRes{}
	p.walkPermRes(&root, arr)
	ret := &proto.QueryPermResResponse{
		List: root.Children,
	}
	return ret, nil
}

func (p *rbacServiceImpl) DeletePermRes(_ context.Context, id *proto.PermResId) (*proto.Result, error) {
	err := p.dao.DeletePermRes(id.Value)
	return p.error(err), nil
}

func (p *rbacServiceImpl) updateUserRoles(userId int64, roles []int64) error {
	dataList := p.dao.GetUserRoles(userId)
	old := make([]int, len(dataList))
	arr := make([]int, len(roles))
	// 旧数组
	for i, v := range dataList {
		old[i] = int(v.RoleId)
	}
	for i, v := range roles {
		arr[i] = int(v)
	}
	_, deleted := util.IntArrayDiff(old, arr, func(v int, add bool) {
		if add {
			p.dao.SavePermUserRole(&model.PermUserRole{
				RoleId: int64(v),
				UserId: userId,
			})
		}
	})
	if len(deleted) > 0 {
		p.dao.BatchDeletePermUserRole(
			fmt.Sprintf("user_id = %d AND role_id IN (%_s)",
				userId, util.JoinIntArray(deleted, ",")))
	}
	return nil
}

// 获取资源的深度
func (p *rbacServiceImpl) getResDepth(pid int64) int {
	depth := 0
	for pid > 0 {
		v := p.dao.GetPermResBy("id=$1", pid)
		if v != nil {
			pid = v.Pid
			depth++
		}
	}
	return depth
}

// 更新资源及下级资源的深度
func (p *rbacServiceImpl) updateResDepth(dst *model.PermRes, depth int16) {
	dst.Depth = depth
	_, _ = p.dao.SavePermRes(dst)
	list := p.dao.SelectPermRes("pid=$1", dst.Id)
	for _, v := range list {
		p.updateResDepth(v, depth+1)
	}
}
