package impl

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : rbac_dept_service.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020/12/02 13:02:38
 * description :
 * history :
 */

//todo: 用户可以添加禁用权限, 关联权限时可以选择菜单的增删改查

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
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
	"github.com/ixre/gof/typeconv"
	"github.com/ixre/gof/types"
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

func NewRbacService(s storage.Interface, o orm.Orm, registryRepo registry.IRegistryRepo) proto.RbacServiceServer {
	return &rbacServiceImpl{
		s:            s,
		registryRepo: registryRepo,
		dao:          impl.NewRbacDao(o),
	}
}

func (p *rbacServiceImpl) createLoginLog(userId int, ipAddress string, isSuccess int) {

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
			p.createLoginLog(0, r.IpAddress, 1) // 登录失败
			return &proto.RbacLoginResponse{
				ErrCode: 3,
				ErrMsg:  "密码不正确",
			}, nil
		}
		dst := &proto.RbacLoginResponse{
			UserId: 0,
			Roles:  []string{"master", "admin"},
		}
		return p.withAccessToken("master", dst, expires)
	}
	// 普通系统用户登录
	usr := p.dao.GetUserBy("username = $1", r.Username)
	if usr == nil {
		return &proto.RbacLoginResponse{
			ErrCode: 2,
			ErrMsg:  "用户不存在",
		}, nil
	}
	decPwd := crypto.Sha1([]byte(r.Password + usr.Salt))
	if usr.Password != decPwd {
		p.createLoginLog(usr.Id, r.IpAddress, 3) // 登录失败
		return &proto.RbacLoginResponse{
			ErrCode: 3,
			ErrMsg:  "密码不正确",
		}, nil
	}
	if usr.Enabled != 1 {
		p.createLoginLog(usr.Id, r.IpAddress, 4) // 登录失败
		return &proto.RbacLoginResponse{
			ErrCode: 4,
			ErrMsg:  "用户已停用",
		}, nil
	}
	p.createLoginLog(usr.Id, r.IpAddress, 0) // 登录成功
	dst := &proto.RbacLoginResponse{
		UserId: int64(usr.Id),
	}
	_, roles := p.getUserRoles(usr.Id)
	for _, v := range roles {
		dst.Roles = append(dst.Roles, v.Code)
	}
	return p.withAccessToken(usr.Username, dst, expires)
}

// 返回带有令牌的结果
func (p *rbacServiceImpl) withAccessToken(username string,
	dst *proto.RbacLoginResponse, expires int) (*proto.RbacLoginResponse, error) {
	accessToken, err := p.createAccessToken(dst.UserId, username,
		strings.Join(dst.Roles, ","), expires)
	dst.AccessToken = accessToken
	if err != nil {
		dst.ErrCode = 2
		dst.ErrMsg = err.Error()
	}
	return dst, nil
}

// 创建令牌
func (p *rbacServiceImpl) createAccessToken(userId int64, username string, perm string, exp int) (string, error) {
	if exp <= 0 {
		exp = int((time.Hour * 24 * 365).Seconds())
	}
	var claims = jwt.MapClaims{
		"exp":    time.Now().Add(time.Second * time.Duration(exp)).Unix(),
		"aud":    userId,
		"iss":    "go2o",
		"sub":    "go2o-rbac-token",
		"name":   username,
		"x-perm": perm,
	}
	key, _ := p.registryRepo.GetValue(registry.SysJWTSecret)
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	return token.SignedString([]byte(key))
}

// 检查令牌是否有效并返回新的令牌
func (p *rbacServiceImpl) CheckRBACToken(_ context.Context, request *proto.RbacCheckTokenRequest) (*proto.RbacCheckTokenResponse, error) {
	if len(request.AccessToken) == 0 {
		return &proto.RbacCheckTokenResponse{Error: "令牌不能为空"}, nil
	}
	jwtSecret, err := p.registryRepo.GetValue(registry.SysJWTSecret)
	if err != nil {
		log.Println("[ GO2O][ ERROR]: check access token error ", err.Error())
		return &proto.RbacCheckTokenResponse{Error: err.Error()}, nil
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
		return &proto.RbacCheckTokenResponse{Error: "令牌无效"}, nil
	}
	if !dstClaims.VerifyIssuer("go2o", true) ||
		dstClaims["sub"] != "go2o-rbac-token" {
		return &proto.RbacCheckTokenResponse{Error: "未知颁发者的令牌"}, nil
	}
	// 令牌过期时间
	exp := int64(dstClaims["exp"].(float64))
	// 判断是否有效
	if !tk.Valid {
		ve, _ := err.(*jwt.ValidationError)
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return &proto.RbacCheckTokenResponse{
				Error:            "令牌已过期",
				IsExpires:        true,
				TokenExpiresTime: exp,
			}, nil
		}
		return &proto.RbacCheckTokenResponse{Error: "令牌无效:" + ve.Error()}, nil
	}
	aud := int64(typeconv.MustInt(dstClaims["aud"]))
	// 如果设置了续期参数
	if exp <= request.CheckExpireTime {
		exp := int((time.Hour * 24 * 365).Seconds())
		dstClaims["exp"] = exp
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, dstClaims)
		accessToken, _ := token.SignedString([]byte(jwtSecret))
		return &proto.RbacCheckTokenResponse{
			UserId:           aud,
			TokenExpiresTime: int64(exp),
			RenewAccessToken: accessToken,
		}, nil
	}
	return &proto.RbacCheckTokenResponse{
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
func (p *rbacServiceImpl) getUserRoles(userId int) ([]int64, []*model.RbacRole) {
	userRoles := p.dao.GetUserRoles(userId)
	// 绑定角色ID
	roles := make([]int64, len(userRoles))
	roleList := make([]int, len(userRoles))
	for i, v := range userRoles {
		roles[i] = int64(v.RoleId)
		roleList[i] = int(v.RoleId)
	}
	// 获取角色的权限
	rolesList := p.dao.SelectPermRole(
		fmt.Sprintf("id IN (%s)", util.JoinIntArray(roleList, ",")))
	return roles, rolesList
}

// 移动资源顺序
func (p *rbacServiceImpl) MoveResourceOrdinal(_ context.Context, r *proto.MoveResourceOrdinalRequest) (*proto.Result, error) {
	res := p.dao.GetRbacResource(r.ResourceId)
	if res == nil {
		return p.error(errors.New("no such data")), nil
	}
	// 获取交换的对象
	var swapRes *model.RbacRes
	if r.Direction == 0 { // 向上移,获取上一个
		swapRes = p.dao.GetRbacResourceBy(
			`sort_num < $1 AND pid = $2 AND depth=$3 WHERE is_forbidden <> 1 ORDER BY sort_num DESC`,
			res.SortNum, res.Pid, res.Depth)
	} else {
		swapRes = p.dao.GetRbacResourceBy(
			`sort_num > $1 AND pid = $2 AND depth=$3 WHERE is_forbidden <> 1 ORDER BY sort_num ASC`,
			res.SortNum, res.Pid, res.Depth)
	}
	// 交换顺序
	if swapRes != nil {
		sortNum := swapRes.SortNum
		swapRes.SortNum = res.SortNum
		res.SortNum = sortNum
		p.dao.SaveRbacResource(res)
		p.dao.SaveRbacResource(swapRes)
	}
	return p.success(nil), nil
}

// 　如果上级菜单未加入,则加入上级菜单
func (p *rbacServiceImpl) appendParentResource(arr *[]*model.RbacRes) {
	mp := make(map[int]*model.RbacRes)
	for _, v := range *arr {
		mp[int(v.Id)] = v
	}
	for _, v := range *arr {
		if _, ok := mp[v.Pid]; !ok && v.Pid > 0 {
			pid := v.Pid
			for pid > 0 {
				d := p.dao.GetRbacResource(pid)
				if d == nil {
					break
				}
				mp[pid] = d
				pid = d.Pid
				*arr = append(*arr, d)
			}
		}
	}
}

// GetUserResource 获取用户的资源,在前端处理排序问题
func (p *rbacServiceImpl) GetUserResource(_ context.Context, r *proto.RbacUserResourceRequest) (*proto.RbacUserResourceResponse, error) {
	dst := &proto.RbacUserResourceResponse{}
	var resList []*model.RbacRes
	rolePermMap := make(map[int]int, 0)
	if r.UserId <= 0 { // master为超级管理员,拥有权限管理权限,admin为管理员
		dst.Roles = []string{"master", "admin"}
		resList = p.dao.SelectPermRes("is_forbidden <> 1 AND is_enabled = 1")
	} else {
		usr := p.dao.GetUser(r.UserId)
		if usr == nil {
			return nil, fmt.Errorf("no such user %v", r.UserId)
		}
		_, userRoles := p.getUserRoles(int(r.UserId))
		roleList := make([]int, len(userRoles))
		for i, v := range userRoles {
			roleList[i] = int(v.Id)
			dst.Roles = append(dst.Roles, v.Code)
		}
		resList = p.dao.GetRoleResources(roleList)
		p.appendParentResource(&resList) //　添加上级资源
		roleResList := p.dao.GetRoleResList(roleList)
		for _, v := range roleResList {
			rolePermMap[int(v.ResId)] = v.PermFlag
		}
	}
	// 获取菜单
	root := proto.SUserMenu{}
	wg := sync.WaitGroup{}
	var f func(*sync.WaitGroup, *proto.SUserMenu, []*model.RbacRes)
	f = func(w *sync.WaitGroup, root *proto.SUserMenu, arr []*model.RbacRes) {
		root.Children = []*proto.SUserMenu{}
		for _, v := range arr {
			if v.AppIndex != int(r.AppIndex) {
				// 其他应用的资源排除
				continue
			}
			if v.IsMenu == 0 {
				// 非菜单资源排除
				continue
			}
			if v.Pid == int(root.Id) {
				c := &proto.SUserMenu{
					Id:            int64(v.Id),
					Key:           v.ResKey,
					Name:          v.Name,
					Path:          v.Path,
					Icon:          v.Icon,
					SortNum:       int32(v.SortNum),
					ComponentName: v.ComponentName,
				}
				c.Children = make([]*proto.SUserMenu, 0)
				root.Children = append(root.Children, c)
				w.Add(1)
				go f(w, c, arr)
			}
		}
		w.Done()
	}
	wg.Add(1)
	f(&wg, &root, resList)
	wg.Wait()
	dst.Menu = root.Children
	// 普通用户返回权限Keys,格式如:["A0101","A010102+7"],不用区分应用
	if r.UserId > 0 {
		for _, v := range resList {
			if len(v.ResKey) > 0 {
				// 添加权限flag到key中
				flag := rolePermMap[int(v.Id)]
				if flag > 0 {
					v.ResKey = fmt.Sprintf("%s+%d", v.ResKey, flag)
				}
				dst.ResourceKeys = append(dst.ResourceKeys, v.ResKey)
			}
		}
	}
	return dst, nil
}

func walkDepartTree(node *proto.SRbacTree, nodeList []*model.RbacDepart) {
	node.Children = []*proto.SRbacTree{}
	for _, v := range nodeList {
		if v.Pid == int(node.Id) {
			v := &proto.SRbacTree{
				Id:       int64(v.Id),
				Label:    v.Name,
				Children: make([]*proto.SRbacTree, 0),
			}
			node.Children = append(node.Children, v)
			walkDepartTree(v, nodeList)
		}
	}
}

// 部门树形数据
func (p *rbacServiceImpl) DepartTree(_ context.Context, empty *proto.Empty) (*proto.SRbacTree, error) {
	root := &proto.SRbacTree{
		Id:       0,
		Label:    "根节点",
		Children: make([]*proto.SRbacTree, 0),
	}
	list := p.dao.SelectPermDept("")
	walkDepartTree(root, list)
	return root, nil
}

// 保存部门
func (p *rbacServiceImpl) SaveDepart(_ context.Context, r *proto.SaveDepartRequest) (*proto.SaveDepartResponse, error) {
	var dst *model.RbacDepart
	if r.Id > 0 {
		dst = p.dao.GetDepart(r.Id)
	} else {
		dst = &model.RbacDepart{}
		dst.CreateTime = int(time.Now().Unix())
	}

	dst.Name = r.Name
	dst.Code = r.Code
	dst.Pid = int(r.Pid)
	dst.Enabled = int(r.Enabled)

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

func (p *rbacServiceImpl) parsePermDept(v *model.RbacDepart) *proto.SPermDept {
	return &proto.SPermDept{
		Id:         int64(v.Id),
		Name:       v.Name,
		Code:       v.Code,
		Pid:        int64(v.Pid),
		Enabled:    int32(v.Enabled),
		CreateTime: int64(v.CreateTime),
	}
}

// 保存岗位
func (p *rbacServiceImpl) SaveJob(_ context.Context, r *proto.SaveJobRequest) (*proto.SaveJobResponse, error) {
	var dst *model.RbacJob
	if r.Id > 0 {
		dst = p.dao.GetJob(r.Id)
	} else {
		dst = &model.RbacJob{}
		dst.CreateTime = int(time.Now().Unix())
	}

	dst.Name = r.Name
	dst.Enabled = int(r.Enabled)
	dst.Sort = int(r.Sort)
	dst.DeptId = int(r.DeptId)

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

func (p *rbacServiceImpl) parsePermJob(v *model.RbacJob) *proto.SRbacJob {
	return &proto.SRbacJob{
		Id:         int64(v.Id),
		Name:       v.Name,
		Enabled:    int32(v.Enabled),
		Sort:       int32(v.Sort),
		DeptId:     int64(v.DeptId),
		CreateTime: int64(v.CreateTime),
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
			DeptName:   typeconv.Stringify(v["deptName"]),
			CreateTime: int64(typeconv.MustInt(v["createTime"])),
		}
	}
	return ret, nil
}

// 保存系统用户
func (p *rbacServiceImpl) SaveUser(_ context.Context, r *proto.SaveRbacUserRequest) (*proto.SaveRbacUserResponse, error) {
	var dst *model.RbacUser
	if r.Id > 0 {
		dst = p.dao.GetUser(r.Id)
		if dst == nil {
			return &proto.SaveRbacUserResponse{
				ErrCode: 2,
				ErrMsg:  "no such record",
			}, nil
		}
	} else {
		dst = &model.RbacUser{}
		dst.Salt = util.RandString(4)
		dst.CreateTime = int(time.Now().Unix())
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
	dst.Avatar = r.ProfilePhoto
	dst.Nickname = r.Nickname
	dst.Gender = int(r.Gender)
	dst.Email = r.Email
	dst.Phone = r.Phone
	dst.DeptId = int(r.DeptId)
	dst.JobId = int(r.JobId)
	dst.Enabled = int(r.Enabled)
	dst.LastLogin = int(r.LastLogin)
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

func (p *rbacServiceImpl) parsePermUser(v *model.RbacUser) *proto.SRbacUser {
	return &proto.SRbacUser{
		Id:           int64(v.Id),
		Username:     v.Username,
		Password:     v.Password,
		Flag:         int32(v.Flag),
		ProfilePhoto: v.Avatar,
		Nickname:     v.Nickname,
		Gender:       int32(v.Gender),
		Email:        v.Email,
		Phone:        v.Phone,
		DeptId:       int64(v.DeptId),
		JobId:        int64(v.JobId),
		Enabled:      int32(v.Enabled),
		LastLogin:    int64(v.LastLogin),
		CreateTime:   int64(v.CreateTime),
	}
}

// 获取系统用户
func (p *rbacServiceImpl) GetUser(_ context.Context, id *proto.RbacUserId) (*proto.SRbacUser, error) {
	v := p.dao.GetUser(id.Value)
	if v == nil {
		return nil, fmt.Errorf("no such user %v", id.Value)
	}
	dst := p.parsePermUser(v)
	dst.Roles, _ = p.getUserRoles(v.Id)
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
			Id:           int64(typeconv.MustInt(v["id"])),
			Username:     typeconv.Stringify(v["username"]),
			Password:     typeconv.Stringify(v["password"]),
			Flag:         int32(typeconv.MustInt(v["flag"])),
			ProfilePhoto: typeconv.Stringify(v["avatar"]),
			Nickname:     typeconv.Stringify(v["nickname"]),
			Gender:       typeconv.Stringify(v["gender"]),
			Email:        typeconv.Stringify(v["email"]),
			Phone:        typeconv.Stringify(v["phone"]),
			DeptId:       int64(typeconv.MustInt(v["deptId"])),
			JobId:        int64(typeconv.MustInt(v["jobId"])),
			Enabled:      int32(typeconv.MustInt(v["enabled"])),
			LastLogin:    int64(typeconv.MustInt(v["lastLogin"])),
			CreateTime:   int64(typeconv.MustInt(v["createTime"])),
		}
	}
	return ret, nil
}

// 保存角色
func (p *rbacServiceImpl) SavePermRole(_ context.Context, r *proto.SaveRbacRoleRequest) (*proto.SaveRbacRoleResponse, error) {
	var dst *model.RbacRole
	if r.Id > 0 {
		dst = p.dao.GetRole(r.Id)
	} else {
		dst = &model.RbacRole{}
		dst.CreateTime = int(time.Now().Unix())
		dst.Code = r.Code
	}
	dst.Name = r.Name
	dst.Level = int(r.Level)
	dst.DataScope = r.DataScope
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
	role := p.dao.GetRole(r.RoleId)
	if role == nil {
		return p.error(errors.New("角色不存在")), nil
	}
	if role.Code == "admin" {
		return p.error(errors.New("管理员已拥有全部权限")), nil
	}
	arr := make([]int, 0)
	permMap := make(map[int]int32, 0)
	for _, v := range r.Resources {
		arr = append(arr, int(v.ResId))
		permMap[int(v.ResId)] = v.PermFlag
	}
	// 旧数组
	dataList := p.dao.GetRoleResList([]int{int(r.RoleId)})
	old := make([]int, len(dataList))
	//　更新数组
	mp := make(map[int]*model.RbacRoleRes, 0)
	for i, v := range dataList {
		old[i] = int(v.ResId)
		mp[int(v.ResId)] = v
	}
	_, deleted := util.IntArrayDiff(old, arr, func(resId int, add bool) {
		var id int64
		if !add {
			id = int64(mp[resId].Id)
		}
		p.dao.SavePermRoleRes(&model.RbacRoleRes{
			Id:       int(id),
			ResId:    resId,
			RoleId:   int(r.RoleId),
			PermFlag: int(permMap[int(resId)]),
		})
	})
	if len(deleted) > 0 {
		p.dao.BatchDeletePermRoleRes(
			fmt.Sprintf("role_id = %d AND res_id IN (%s)",
				r.RoleId, util.JoinIntArray(deleted, ",")))
	}
	return p.error(nil), nil
}

func (p *rbacServiceImpl) parsePermRole(v *model.RbacRole) *proto.SRbacRole {
	return &proto.SRbacRole{
		Id:         int64(v.Id),
		Code:       v.Code,
		Name:       v.Name,
		Level:      int32(v.Level),
		DataScope:  v.DataScope,
		Remark:     v.Remark,
		CreateTime: int64(v.CreateTime),
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
	res := p.dao.GetRoleResList([]int{int(v.Id)})
	dst.ResourceList = make([]*proto.SRolePermPair, 0)
	for _, v := range res {
		dst.ResourceList = append(dst.ResourceList, &proto.SRolePermPair{
			ResId:    int64(v.ResId),
			PermFlag: int32(v.PermFlag),
		})
	}
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

// DeletePermRole 删除角色
func (p *rbacServiceImpl) DeletePermRole(_ context.Context, id *proto.RbacRoleId) (*proto.Result, error) {
	err := p.dao.DeletePermRole(id.Value)
	return p.error(err), nil
}

// PagingPermRole 角色分页信息
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
			DataScope:  typeconv.Stringify(v["dataScope"]),
			Code:       typeconv.Stringify(v["code"]),
			Remark:     typeconv.Stringify(v["remark"]),
			CreateTime: int64(typeconv.MustInt(v["createTime"])),
		}
	}
	return ret, nil
}

// 验证上级资源是否合法
func (p *rbacServiceImpl) checkParentResource(id int, currentPid int, pid int) (model.RbacRes, error) {
	var parentRes model.RbacRes
	// 检测上级
	if pid <= 0 {
		return parentRes, nil
	}
	// 检测上级是否为自己
	if pid == id && id > 0 {
		return parentRes, errors.New("不能将自己指定为上级资源")
	}
	// 检测上级是否为下级
	if currentPid > 0 {
		var parent *model.RbacRes = p.dao.GetRbacResource(pid)
		// 获取上级的Key,用于获取
		if parent != nil {
			parentRes = *parent
		}
		for parent != nil && parent.Pid > 0 {
			parent = p.dao.GetRbacResource(parent.Pid)
			if parent != nil && int(parent.Id) == pid {
				return parentRes, errors.New("不能选择下级作为上级资源")
			}
		}
	}
	return parentRes, nil
}

// 生成资源标识
func (p *rbacServiceImpl) GenerateResourceKey(parent model.RbacRes) string {
	maxKey := p.dao.GetMaxResouceKey(int(parent.Id))
	l := len(maxKey)
	// 一级资源,如果未以字母命名,则启用规则,并依次命名
	if parent.Id == 0 {
		if l == 0 || l > 1 {
			return "A"
		}
		return strings.ToUpper(string(rune(maxKey[0]) + 1))
	}
	// 下级资源采用数字编号
	if l == 0 {
		return fmt.Sprintf("%s01", parent.ResKey)
	}
	// 获取末尾编号,如:05,并进行累加
	var v = 0
	if len(maxKey) >= 2 {
		// 如果存在一级菜单挪动到二级,则计数会报错
		v, _ = strconv.Atoi(maxKey[l-2:])
	}
	v += 1
	if v < 10 {
		return fmt.Sprintf("%s0%d", parent.ResKey, v)
	}
	return fmt.Sprintf("%s%d", parent.ResKey, v)
}

// 保存PermRes
func (p *rbacServiceImpl) SaveRbacResource(_ context.Context, r *proto.SaveRbacResRequest) (*proto.SaveRbacResResponse, error) {
	var dst *model.RbacRes
	if r.Id > 0 {
		if dst = p.dao.GetRbacResource(r.Id); dst == nil {
			return &proto.SaveRbacResResponse{
				ErrCode: 2,
				ErrMsg:  "no such data",
			}, nil
		}
	} else {
		dst = &model.RbacRes{}
		dst.CreateTime = int(time.Now().Unix())
		dst.Pid = int(r.Pid) // 设置上级,用于生成资源key
		dst.Depth = 0
	}

	// 如果pid传入小于0,则强制为0,以避免数据无法显示
	if r.Pid < 0 {
		r.Pid = 0
	}
	if r.Pid > 0 {
		// 限制下级资源路径不能以'/'开头,以避免无法找到资源的情况
		if len(r.Path) > 0 && r.Path[0] == '/' {
			return &proto.SaveRbacResResponse{
				ErrCode: 3,
				ErrMsg:  "该资源(包含上级资源)路径不能以'/'开头",
			}, nil
		}
	}
	parent, err := p.checkParentResource(int(dst.Id), int(dst.Pid), int(r.Pid))
	if err != nil {
		return &proto.SaveRbacResResponse{
			ErrCode: 2,
			ErrMsg:  err.Error(),
		}, nil
	}
	// 如果新增, 则生成key
	if r.Id <= 0 || len(dst.ResKey) == 0 {
		dst.ResKey = p.GenerateResourceKey(parent)
		// 新增时设置应用序号
		dst.AppIndex = int(r.AppIndex)
		// 如果包含上级,则与上级的应用序号保持一致
		if parent.Id > 0 {
			dst.AppIndex = parent.AppIndex
		}
	}

	// 上级是否改变
	var parentChanged = r.Id > 0 && (dst.Pid != int(r.Pid) || (r.Pid > 0 && dst.Depth == 0))
	dst.Name = r.Name
	dst.ResType = int(r.ResType)
	dst.Pid = int(r.Pid)
	dst.Path = r.Path
	dst.Icon = r.Icon
	dst.SortNum = int(r.SortNum)
	dst.IsMenu = types.ElseInt(r.IsMenu, 1, 0)
	dst.IsEnabled = types.ElseInt(r.IsEnabled, 1, 0)
	dst.ComponentName = r.ComponentName
	dst.AppIndex = int(r.AppIndex)
	// 如果未设置排列序号,或者更改了上级,则需系统自动编号
	if dst.SortNum <= 0 || parentChanged {
		dst.SortNum = p.dao.GetMaxResourceSortNum(int(dst.Pid)) + 1
	}
	id, err := p.dao.SaveRbacResource(dst)
	ret := &proto.SaveRbacResResponse{
		Id: int64(id),
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	} else {
		if parentChanged {
			depth := p.getResDepth(dst.Pid)
			p.updateResDepth(dst, depth)
		}
	}
	return ret, nil
}

func (p *rbacServiceImpl) parseRbacRes(v *model.RbacRes) *proto.SRbacRes {
	return &proto.SRbacRes{
		Id:            int64(v.Id),
		Name:          v.Name,
		ResType:       int32(v.ResType),
		Pid:           int64(v.Pid),
		Key:           v.ResKey,
		Path:          v.Path,
		Icon:          v.Icon,
		SortNum:       int32(v.SortNum),
		IsMenu:        v.IsMenu == 1,
		IsEnabled:     v.IsEnabled == 1,
		CreateTime:    int64(v.CreateTime),
		ComponentName: v.ComponentName,
		AppIndex:      int32(v.AppIndex),
	}
}

// 获取PermRes
func (p *rbacServiceImpl) GetRbacRes(_ context.Context, id *proto.PermResId) (*proto.SRbacRes, error) {
	v := p.dao.GetRbacResource(id.Value)
	if v == nil {
		return nil, fmt.Errorf("no such resource %v", id.Value)
	}
	return p.parseRbacRes(v), nil
}

// 获取PermRes列表
func (p *rbacServiceImpl) QueryRbacResourceList(_ context.Context, r *proto.QueryRbacResRequest) (*proto.QueryRbacResourceResponse, error) {
	var where string = "is_forbidden <> 1"
	if r.Keyword != "" {
		where += " AND name LIKE '%" + r.Keyword + "%'"
	}
	if r.OnlyMenu {
		where += " AND is_menu = 1"
	}
	arr := p.dao.SelectPermRes(where + " ORDER BY sort_num ASC")
	// 获取第一级分类
	roots := p.queryResChildren(int(r.ParentId), arr)
	initial := make([]int64, 0)

	// 初始化已选择的节点
	if r.ParentId <= 0 && r.InitialId > 0 {
		findParent := func(pid int, arr []*model.RbacRes) int {
			for _, v := range arr {
				if v.Id == pid && v.Pid > 0 {
					return int(v.Pid)
				}
			}
			return pid
		}
		for pid := int(r.InitialId); pid > 0; {
			id := findParent(pid, arr)
			if id == pid {
				break
			}
			initial = append([]int64{int64(id)}, initial...)
			pid = id
		}
	}
	ret := &proto.QueryRbacResourceResponse{
		List:        roots,
		InitialList: initial,
	}
	return ret, nil
}

func (p *rbacServiceImpl) queryResChildren(parentId int, arr []*model.RbacRes) []*proto.SRbacRes {
	var list []*proto.SRbacRes
	for _, v := range arr {
		if v.Pid != parentId {
			continue
		}
		c := p.parseRbacRes(v)
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

// 删除资源
func (p *rbacServiceImpl) DeleteRbacResource(_ context.Context, id *proto.PermResId) (*proto.Result, error) {
	res := p.dao.GetRbacResource(id.Value)
	if res == nil {
		return p.error(errors.New("资源不存在")), nil
	}
	res.IsForbidden = 1
	_, err := p.dao.SaveRbacResource(res)
	return p.error(err), nil
}

func (p *rbacServiceImpl) updateUserRoles(userId int64, roles []int64) error {
	dataList := p.dao.GetUserRoles(int(userId))
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
			p.dao.SaveUserRole(&model.RbacUserRole{
				RoleId: int(v),
				UserId: int(userId),
			})
		}
	})
	if len(deleted) > 0 {
		p.dao.BatchDeleteUserRole(
			fmt.Sprintf("user_id = %d AND role_id IN (%s)",
				userId, util.JoinIntArray(deleted, ",")))
	}
	return nil
}

// 获取资源的深度
func (p *rbacServiceImpl) getResDepth(pid int) int {
	depth := 0
	for pid > 0 {
		v := p.dao.GetRbacResourceBy("id = $1", pid)
		if v != nil {
			pid = v.Pid
			depth++
		}
	}
	return depth
}

// 更新资源及下级资源的深度
func (p *rbacServiceImpl) updateResDepth(dst *model.RbacRes, depth int) {
	dst.Depth = depth
	_, _ = p.dao.SaveRbacResource(dst)
	list := p.dao.SelectPermRes("pid=$1", dst.Id)
	for _, v := range list {
		p.updateResDepth(v, depth+1)
	}
}

// PagingLoginLog implements proto.RbacServiceServer.
func (p *rbacServiceImpl) PagingLoginLog(_ context.Context, r *proto.LoginLogPagingRequest) (*proto.LoginLogPagingResponse, error) {
	//todo:  keyword
	total, rows := p.dao.PagingQueryLoginLog(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.LoginLogPagingResponse{
		Total: int64(total),
		Value: make([]*proto.PagingLoginLog, len(rows)),
	}
	for i, v := range rows {
		ret.Value[i] = &proto.PagingLoginLog{
			Id:         int64(typeconv.MustInt(v["id"])),
			UserId:     int64(typeconv.MustInt(v["userId"])),
			Username:   typeconv.Stringify(v["username"]),
			Nickname:   typeconv.Stringify(v["nickname"]),
			Ip:         typeconv.Stringify(v["ip"]),
			IsSuccess:  int32(typeconv.MustInt(v["isSuccess"])),
			CreateTime: int64(typeconv.MustInt(v["createTime"])),
		}
		if r := ret.Value[i]; r.UserId == 0 {
			r.Username = "master"
			r.Nickname = "超级管理员"
		}
	}
	return ret, nil
}
