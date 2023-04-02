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
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
)

var _ proto.RbacServiceServer = new(rbacServiceImpl)

// 基于角色的权限服务
type rbacServiceImpl struct {
	dao          dao.IRbacDao
	registryRepo registry.IRegistryRepo
	s            storage.Interface
	serviceUtil
	proto.UnimplementedRbacServiceServer
}

func NewRbacService(s storage.Interface, o orm.Orm, registryRepo registry.IRegistryRepo) *rbacServiceImpl {
	return &rbacServiceImpl{
		s:            s,
		registryRepo: registryRepo,
		dao:          impl.NewRbacDao(o),
	}
}

func (p *rbacServiceImpl) UserLogin(_ context.Context, r *proto.RbacLoginRequest) (*proto.RbacLoginResponse, error) {
	if len(r.Password) != 32 {
		return &proto.RbacLoginResponse{
			ErrCode: 1,
			ErrMsg:  "密码长度不正确，应该为32位长度的md5字符",
		}, nil
	}
	expires := 3600 * 24
	// 超级管理员登录
	if r.Username == "master" {
		superPwd, _ := p.registryRepo.GetValue(registry.SysSuperLoginToken)
		encPwd := domain.Sha1Pwd(r.Username+r.Password, "")
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
		return p.withAccessToken("master", dst, expires)
	}
	// 普通系统用户登录
	usr := p.dao.GetUserBy("usr=$1", r.Username)
	if usr == nil {
		return &proto.RbacLoginResponse{
			ErrCode: 2,
			ErrMsg:  "用户不存在",
		}, nil
	}
	decPwd := crypto.Sha1([]byte(r.Password + usr.Salt))
	if usr.Password != decPwd {
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
	return p.withAccessToken(usr.Username, dst, expires)
}

// 返回带有令牌的结果
func (p *rbacServiceImpl) withAccessToken(userName string,
	dst *proto.RbacLoginResponse, expires int) (*proto.RbacLoginResponse, error) {
	accessToken, err := p.createAccessToken(dst.UserId, userName,
		strings.Join(dst.Permissions, ","), expires)
	dst.AccessToken = accessToken
	if err != nil {
		dst.ErrCode = 2
		dst.ErrMsg = err.Error()
	}
	return dst, nil
}

// 创建令牌
func (p *rbacServiceImpl) createAccessToken(userId int64, userName string, perm string, exp int) (string, error) {
	if exp <= 0 {
		exp = int((time.Hour * 24 * 365).Seconds())
	}
	var claims = jwt.MapClaims{
		"exp":    time.Now().Add(time.Second * time.Duration(exp)).Unix(),
		"aud":    userId,
		"iss":    "go2o",
		"sub":    "go2o-rbac-token",
		"name":   userName,
		"x-perm": perm,
	}
	key, _ := p.registryRepo.GetValue(registry.SysJWTSecret)
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	return token.SignedString([]byte(key))
}

// 检查令牌是否有效并返回新的令牌
func (p *rbacServiceImpl) CheckRBACToken(_ context.Context, request *proto.CheckRBACTokenRequest) (*proto.CheckRBACTokenResponse, error) {
	if len(request.AccessToken) == 0 {
		return &proto.CheckRBACTokenResponse{Error: "令牌不能为空"}, nil
	}
	jwtSecret, err := p.registryRepo.GetValue(registry.SysJWTSecret)
	if err != nil {
		log.Println("[ GO2O][ ERROR]: check access token error ", err.Error())
		return &proto.CheckRBACTokenResponse{Error: err.Error()}, nil
	}
	// 去掉"Bearer "
	if len(request.AccessToken) > 6 &&
		strings.HasPrefix(request.AccessToken, "Bearer") {
		request.AccessToken = request.AccessToken[7:]
	}
	// 转换token
	dstClaims := jwt.MapClaims{} // 可以用实现了Claim接口的自定义结构
	tk, err := jwt.ParseWithClaims(request.AccessToken, &dstClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if tk == nil {
		return &proto.CheckRBACTokenResponse{Error: "令牌无效"}, nil
	}
	if !dstClaims.VerifyIssuer("go2o", true) ||
		dstClaims["sub"] != "go2o-rbac-token" {
		return &proto.CheckRBACTokenResponse{Error: "未知颁发者的令牌"}, nil
	}
	// 令牌过期时间
	exp := int64(dstClaims["exp"].(float64))
	// 判断是否有效
	if !tk.Valid {
		ve, _ := err.(*jwt.ValidationError)
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return &proto.CheckRBACTokenResponse{
				Error:            "令牌已过期",
				IsExpires:        true,
				TokenExpiresTime: exp,
			}, nil
		}
		return &proto.CheckRBACTokenResponse{Error: "令牌无效:" + ve.Error()}, nil
	}
	aud := int64(typeconv.MustInt(dstClaims["aud"]))
	// 如果设置了续期参数
	if exp <= request.CheckExpireTime {
		exp := int((time.Hour * 24 * 365).Seconds())
		dstClaims["exp"] = exp
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, dstClaims)
		accessToken, _ := token.SignedString([]byte(jwtSecret))
		return &proto.CheckRBACTokenResponse{
			UserId:           aud,
			TokenExpiresTime: int64(exp),
			RenewAccessToken: accessToken,
		}, nil
	}
	return &proto.CheckRBACTokenResponse{
		UserId:           aud,
		TokenExpiresTime: exp,
	}, nil
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
		fmt.Sprintf("id IN (%s)", util.JoinIntArray(roleList, ",")))
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
		resList = p.dao.SelectPermRes("is_forbidden <> 1 AND is_hidden <> 1")
		// 获取管理员
		for _, v := range resList {
			v.Permission = "master,admin"
		}
	} else {
		dst.Roles, dst.Permissions = p.getUserRolesPerm(r.UserId)
		usr := p.dao.GetUser(r.UserId)
		if usr == nil {
			return nil, fmt.Errorf("no such user %v", r.UserId)
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
func (p *rbacServiceImpl) SaveDepart(_ context.Context, r *proto.SaveDepartRequest) (*proto.SaveDepartResponse, error) {
	var dst *model.PermDept
	if r.Id > 0 {
		dst = p.dao.GetDepart(r.Id)
	} else {
		dst = &model.PermDept{}
		dst.CreateTime = time.Now().Unix()
	}

	dst.Name = r.Name
	dst.Code = r.Code
	dst.Pid = r.Pid
	dst.Enabled = int16(r.Enabled)

	id, err := p.dao.SaveDepart(dst)
	ret := &proto.SaveDepartResponse{
		Id: int64(id),
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

// 获取部门
func (p *rbacServiceImpl) GetDepart(_ context.Context, id *proto.RbacDepartId) (*proto.SPermDept, error) {
	v := p.dao.GetDepart(id.Value)
	if v == nil {
		return nil, fmt.Errorf("no such dept: %v", id)
	}
	return p.parsePermDept(v), nil
}

func (p *rbacServiceImpl) DeleteDepart(_ context.Context, id *proto.RbacDepartId) (*proto.Result, error) {
	err := p.dao.DeleteDepart(id.Value)
	return p.error(err), nil
}

func (p *rbacServiceImpl) parsePermDept(v *model.PermDept) *proto.SPermDept {
	return &proto.SPermDept{
		Id:         v.Id,
		Name:       v.Name,
		Code:       v.Code,
		Pid:        v.Pid,
		Enabled:    int32(v.Enabled),
		CreateTime: v.CreateTime,
	}
}

// 保存岗位
func (p *rbacServiceImpl) SaveJob(_ context.Context, r *proto.SaveJobRequest) (*proto.SaveJobResponse, error) {
	var dst *model.PermJob
	if r.Id > 0 {
		dst = p.dao.GetJob(r.Id)
	} else {
		dst = &model.PermJob{}
		dst.CreateTime = time.Now().Unix()
	}

	dst.Name = r.Name
	dst.Enabled = int16(r.Enabled)
	dst.Sort = int(r.Sort)
	dst.DeptId = r.DeptId

	id, err := p.dao.SaveJob(dst)
	ret := &proto.SaveJobResponse{
		Id: int64(id),
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

func (p *rbacServiceImpl) parsePermJob(v *model.PermJob) *proto.SRbacJob {
	return &proto.SRbacJob{
		Id:         v.Id,
		Name:       v.Name,
		Enabled:    int32(v.Enabled),
		Sort:       int32(v.Sort),
		DeptId:     v.DeptId,
		CreateTime: v.CreateTime,
	}
}

// 获取岗位
func (p *rbacServiceImpl) GetJob(_ context.Context, id *proto.RbacJobId) (*proto.SRbacJob, error) {
	v := p.dao.GetJob(id.Value)
	if v == nil {
		return nil, fmt.Errorf("no such job: %v", id.Value)
	}
	return p.parsePermJob(v), nil
}

// 获取岗位列表
func (p *rbacServiceImpl) QueryJobList(_ context.Context, r *proto.QueryJobRequest) (*proto.QueryJobResponse, error) {
	where := ""
	if r.DepartId > 0 {
		arr := p.walkDepartArray(int(r.DepartId))
		if len(where) > 0 {
			where += " AND "
		}
		where += " dept_id IN (" + util.JoinIntArray(arr, ",") + ")"
	}
	arr := p.dao.SelectPermJob(where)
	ret := &proto.QueryJobResponse{
		List: make([]*proto.SRbacJob, len(arr)),
	}
	for i, v := range arr {
		ret.List[i] = p.parsePermJob(v)
	}
	return ret, nil
}

func (p *rbacServiceImpl) DeleteJob(_ context.Context, id *proto.RbacJobId) (*proto.Result, error) {
	err := p.dao.DeleteJob(id.Value)
	return p.error(err), nil
}

func (p *rbacServiceImpl) PagingJobList(_ context.Context, r *proto.RbacJobPagingRequest) (*proto.PagingRbacJobResponse, error) {
	total, rows := p.dao.PagingQueryJob(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.PagingRbacJobResponse{
		Total: int64(total),
		Value: make([]*proto.PagingJobList, len(rows)),
	}
	for i, v := range rows {
		ret.Value[i] = &proto.PagingJobList{
			Id:         int64(typeconv.MustInt(v["id"])),
			Name:       typeconv.Stringify(v["name"]),
			Enabled:    int32(typeconv.MustInt(v["enabled"])),
			Sort:       int32(typeconv.MustInt(v["sort"])),
			DeptName:   typeconv.Stringify(v["dept_name"]),
			CreateTime: int64(typeconv.MustInt(v["create_time"])),
		}
	}
	return ret, nil
}

// 保存系统用户
func (p *rbacServiceImpl) SaveUser(_ context.Context, r *proto.SaveRbacUserRequest) (*proto.SaveRbacUserResponse, error) {
	var dst *model.PermUser
	if r.Id > 0 {
		dst = p.dao.GetUser(r.Id)
		if dst == nil {
			return &proto.SaveRbacUserResponse{
				ErrCode: 2,
				ErrMsg:  "no such record",
			}, nil
		}
	} else {
		dst = &model.PermUser{}
		dst.Salt = util.RandString(4)
		dst.CreateTime = time.Now().Unix()
	}
	if l := len(r.Password); l > 0 {
		if l != 32 {
			return &proto.SaveRbacUserResponse{
				ErrCode: 1,
				ErrMsg:  "非32位md5密码",
			}, nil
		}
		dst.Password = crypto.Sha1([]byte(r.Password + dst.Salt))
	}
	dst.Flag = int(r.Flag)
	dst.Avatar = r.Portrait
	dst.Nickname = r.Nickname
	dst.Gender = r.Gender
	dst.Email = r.Email
	dst.Phone = r.Phone
	dst.DeptId = r.DeptId
	dst.JobId = r.JobId
	dst.Enabled = int16(r.Enabled)
	dst.LastLogin = r.LastLogin
	id, err := p.dao.SaveUser(dst)
	ret := &proto.SaveRbacUserResponse{
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

func (p *rbacServiceImpl) parsePermUser(v *model.PermUser) *proto.SRbacUser {
	return &proto.SRbacUser{
		Id:         v.Id,
		Username:   v.Username,
		Password:   v.Password,
		Flag:       int32(v.Flag),
		Portrait:   v.Avatar,
		Nickname:   v.Nickname,
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
func (p *rbacServiceImpl) GetUser(_ context.Context, id *proto.RbacUserId) (*proto.SRbacUser, error) {
	v := p.dao.GetUser(id.Value)
	if v == nil {
		return nil, fmt.Errorf("no such user %v", id.Value)
	}
	dst := p.parsePermUser(v)
	dst.Roles, dst.Permissions = p.getUserRolesPerm(v.Id)
	return dst, nil
}

func (p *rbacServiceImpl) DeleteUser(_ context.Context, id *proto.RbacUserId) (*proto.Result, error) {
	err := p.dao.DeleteUser(id.Value)
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

func (p *rbacServiceImpl) PagingUser(_ context.Context, r *proto.PagingRbacUserRequest) (*proto.PagingRbacUserResponse, error) {
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
	ret := &proto.PagingRbacUserResponse{
		Total: int64(total),
		Value: make([]*proto.PagingUser, len(rows)),
	}
	for i, v := range rows {
		ret.Value[i] = &proto.PagingUser{
			Id:         int64(typeconv.MustInt(v["id"])),
			Username:   typeconv.Stringify(v["username"]),
			Password:   typeconv.Stringify(v["password"]),
			Flag:       int32(typeconv.MustInt(v["flag"])),
			Portrait:   typeconv.Stringify(v["avatar"]),
			Nickname:   typeconv.Stringify(v["nickname"]),
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
func (p *rbacServiceImpl) SavePermRole(_ context.Context, r *proto.SaveRbacRoleRequest) (*proto.SaveRbacRoleResponse, error) {
	var dst *model.PermRole
	if r.Id > 0 {
		dst = p.dao.GetRole(r.Id)
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
	ret := &proto.SaveRbacRoleResponse{
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

func (p *rbacServiceImpl) parsePermRole(v *model.PermRole) *proto.SRbacRole {
	return &proto.SRbacRole{
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
func (p *rbacServiceImpl) GetRole(_ context.Context, id *proto.RbacRoleId) (*proto.SRbacRole, error) {
	v := p.dao.GetRole(id.Value)
	if v == nil {
		return nil, fmt.Errorf("no such role: %v", id.Value)
	}
	dst := p.parsePermRole(v)
	// 绑定资源ID
	dst.RelatedResIdList = p.dao.GetRoleResList(v.Id)
	return dst, nil
}

// 获取角色列表
func (p *rbacServiceImpl) QueryPermRoleList(_ context.Context, r *proto.QueryRbacRoleRequest) (*proto.QueryRbacRoleResponse, error) {
	arr := p.dao.SelectPermRole("")
	ret := &proto.QueryRbacRoleResponse{
		List: make([]*proto.SRbacRole, len(arr)),
	}
	for i, v := range arr {
		ret.List[i] = p.parsePermRole(v)
	}
	return ret, nil
}

func (p *rbacServiceImpl) DeletePermRole(_ context.Context, id *proto.RbacRoleId) (*proto.Result, error) {
	err := p.dao.DeletePermRole(id.Value)
	return p.error(err), nil
}

func (p *rbacServiceImpl) PagingPermRole(_ context.Context, r *proto.RbacRolePagingRequest) (*proto.PagingRbacRoleResponse, error) {
	total, rows := p.dao.PagingQueryPermRole(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.PagingRbacRoleResponse{
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
func (p *rbacServiceImpl) SavePermRes(_ context.Context, r *proto.SaveRbacResRequest) (*proto.SaveRbacResResponse, error) {
	var dst *model.PermRes
	if r.Id > 0 {
		if dst = p.dao.GetPermRes(r.Id); dst == nil {
			return &proto.SaveRbacResResponse{
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
			return &proto.SaveRbacResResponse{
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
					return &proto.SaveRbacResResponse{
						ErrCode: 2,
						ErrMsg:  "不能选择下级作为上级资源",
					}, nil
				}
			}
		}
		// 限制下级资源路径不能以'/'开头,以避免无法找到资源的情况
		if len(r.Path) > 0 && r.Path[0] == '/' {
			return &proto.SaveRbacResResponse{
				ErrCode: 3,
				ErrMsg:  "该资源(包含上级资源)路径不能以'/'开头",
			}, nil
		}
	}

	// 上级是否改变
	var parentChanged = r.Id > 0 && (dst.Pid != r.Pid || (r.Pid > 0 && dst.Depth == 0))
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
	ret := &proto.SaveRbacResResponse{
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
		return nil, fmt.Errorf("no such resource %v", id.Value)
	}
	return p.parsePermRes(v), nil
}

// 获取PermRes列表
func (p *rbacServiceImpl) QueryResList(_ context.Context, r *proto.QueryRbacResRequest) (*proto.QueryPermResResponse, error) {
	var where string = "is_forbidden <> 1"
	if r.Keyword != "" {
		where += " AND name LIKE '%" + r.Keyword + "%'"
	}
	if r.OnlyMenu {
		where += " AND res_type IN(0,2)"
	}
	arr := p.dao.SelectPermRes(where + " ORDER BY sort_num ASC")
	// 获取第一级分类
	roots := p.queryResChildren(r.ParentId, arr)
	initial := make([]int64, 0)

	// 初始化已选择的节点
	if r.ParentId <= 0 && r.InitialId > 0 {
		findParent := func(pid int64, arr []*model.PermRes) int64 {
			for _, v := range arr {
				if v.Id == pid && v.Pid > 0 {
					return v.Pid
				}
			}
			return pid
		}
		for pid := r.InitialId; pid > 0; {
			id := findParent(pid, arr)
			if id == pid {
				break
			}
			initial = append([]int64{int64(id)}, initial...)
			pid = id
		}
	}
	ret := &proto.QueryPermResResponse{
		List:        roots,
		InitialList: initial,
	}
	return ret, nil
}

func (p *rbacServiceImpl) queryResChildren(parentId int64, arr []*model.PermRes) []*proto.SPermRes {
	var list []*proto.SPermRes
	for _, v := range arr {
		if v.Pid != parentId {
			continue
		}
		c := p.parsePermRes(v)
		c.IsLeaf = true
		for _, r := range arr {
			if r.Pid == v.Id {
				c.IsLeaf = false
				break
			}
		}
		list = append(list, c)
	}
	return list
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
			p.dao.SaveUserRole(&model.PermUserRole{
				RoleId: int64(v),
				UserId: userId,
			})
		}
	})
	if len(deleted) > 0 {
		p.dao.BatchDeleteUserRole(
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
