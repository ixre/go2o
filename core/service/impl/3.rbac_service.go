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
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ixre/go2o/core/dao"
	"github.com/ixre/go2o/core/dao/impl"
	rbac "github.com/ixre/go2o/core/domain/interface/rabc"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	"github.com/ixre/go2o/core/service/proto"
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
	_repo        rbac.IRbacRepository
	serviceUtil
	proto.UnimplementedRbacServiceServer
}

func NewRbacService(s storage.Interface, o orm.Orm,
	repo rbac.IRbacRepository,
	registryRepo registry.IRegistryRepo) proto.RbacServiceServer {
	return &rbacServiceImpl{
		s:            s,
		_repo:        repo,
		registryRepo: registryRepo,
		dao:          impl.NewRbacDao(o),
	}
}

func (p *rbacServiceImpl) createLoginLog(userId int, ipAddress string, isSuccess int) {
	p._repo.LoginLogRepo().Save(&rbac.RbacLoginLog{
		Id:         userId,
		UserId:     userId,
		Ip:         ipAddress,
		IsSuccess:  isSuccess,
		CreateTime: int(time.Now().Unix()),
	})
}

func (p *rbacServiceImpl) UserLogin(_ context.Context, r *proto.RbacLoginRequest) (*proto.RbacLoginResponse, error) {
	if len(r.Password) != 32 {
		return &proto.RbacLoginResponse{
			Code:    1,
			Message: "密码长度不正确，应该为32位长度的md5字符",
		}, nil
	}
	expires := 3600 * 24
	// 超级管理员登录
	if r.Username == "master" {
		superPwd, _ := p.registryRepo.GetValue(registry.SysSuperLoginToken)
		encPwd := domain.SuperPassword(r.Username+r.Password, "")
		if superPwd != encPwd {
			p.createLoginLog(0, r.IpAddress, 1) // 登录失败
			return &proto.RbacLoginResponse{
				Code:    3,
				Message: "密码不正确",
			}, nil
		}

		p.createLoginLog(0, r.IpAddress, 0) // 登录失败
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
			Code:    2,
			Message: "用户不存在",
		}, nil
	}
	decPwd := domain.RbacPassword(r.Password, usr.Salt)
	if usr.Password != decPwd {
		p.createLoginLog(usr.Id, r.IpAddress, 3) // 登录失败
		return &proto.RbacLoginResponse{
			Code:    3,
			Message: "密码不正确",
		}, nil
	}
	if usr.Enabled != 1 {
		p.createLoginLog(usr.Id, r.IpAddress, 4) // 登录失败
		return &proto.RbacLoginResponse{
			Code:    4,
			Message: "用户已停用",
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
		dst.Code = 2
		dst.Message = err.Error()
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
	key, _ := p.registryRepo.GetValue(registry.SysPrivateKey)
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	return token.SignedString([]byte(key))
}

// 检查令牌是否有效并返回新的令牌
func (p *rbacServiceImpl) CheckRBACToken(_ context.Context, request *proto.RbacCheckTokenRequest) (*proto.RbacCheckTokenResponse, error) {
	if len(request.AccessToken) == 0 {
		return &proto.RbacCheckTokenResponse{Code: 1001, Message: "令牌不能为空"}, nil
	}
	jwtSecret, err := p.registryRepo.GetValue(registry.SysPrivateKey)
	if err != nil {
		log.Println("[ GO2O][ ERROR]: check access token error ", err.Error())
		return &proto.RbacCheckTokenResponse{Code: 1002, Message: err.Error()}, nil
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
		return &proto.RbacCheckTokenResponse{
			Code: 1003, Message: "令牌无效"}, nil
	}
	if !dstClaims.VerifyIssuer("go2o", true) ||
		dstClaims["sub"] != "go2o-rbac-token" {
		return &proto.RbacCheckTokenResponse{
			Code: 1004, Message: "未知颁发者的令牌"}, nil
	}
	// 令牌过期时间
	exp := int64(dstClaims["exp"].(float64))
	// 判断是否有效
	if !tk.Valid {
		ve, _ := err.(*jwt.ValidationError)
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return &proto.RbacCheckTokenResponse{
				Code:             1005,
				Message:          "令牌已过期",
				IsExpires:        true,
				TokenExpiresTime: exp,
			}, nil
		}
		return &proto.RbacCheckTokenResponse{
			Code:    1006,
			Message: "令牌无效:" + ve.Error()}, nil
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
	key, _ := p.registryRepo.GetValue(registry.SysPrivateKey)
	return &proto.String{Value: key}, nil
}

// 获取用户的角色和权限
func (p *rbacServiceImpl) getUserRoles(userId int) ([]int64, []*rbac.RbacRole) {
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
func (p *rbacServiceImpl) MoveResourceOrdinal(_ context.Context, r *proto.MoveResourceOrdinalRequest) (*proto.TxResult, error) {
	res := p.dao.GetRbacResource(r.ResourceId)
	if res == nil {
		return p.errorV2(errors.New("no such data")), nil
	}
	// 获取交换的对象
	var swapRes *rbac.RbacRes
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
	return p.successV2(nil), nil
}

// 　如果上级菜单未加入,则加入上级菜单
func (p *rbacServiceImpl) appendParentResource(arr *[]*rbac.RbacRes) {
	mp := make(map[int]*rbac.RbacRes)
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
	var resList []*rbac.RbacRes
	rolePermMap := make(map[int]int, 0)
	if r.UserId <= 0 { // master为超级管理员,拥有权限管理权限,admin为管理员
		dst.Roles = []string{"master", "admin"}
		resList = p.dao.SelectPermRes("is_forbidden <> 1 AND is_enabled = 1")
	} else {
		usr := p.dao.GetUser(r.UserId)
		if usr == nil {
			return nil, fmt.Errorf("no such user %v", r.UserId)
		}
		//log.Println("ss1", time.Now().UnixMilli())
		// 获取用户角色
		_, userRoles := p.getUserRoles(int(r.UserId))
		roleList := make([]int, len(userRoles))
		for i, v := range userRoles {
			roleList[i] = int(v.Id)
			dst.Roles = append(dst.Roles, v.Code)
		}
		// 获取角色所有对应的资源, 多个角色可能存在资源重复,且权限有差异
		resList = p.dao.GetRoleResources(roleList)
		p.appendParentResource(&resList)
		roleResList := p.dao.GetRoleResList(roleList)
		for _, v := range roleResList {
			rolePermMap[int(v.ResId)] = v.PermFlag
		}
		log.Println(typeconv.MustJson(roleResList))
	}
	// 准备菜单数据
	//log.Println("ss2", time.Now().UnixMilli())
	parents := make(map[int][]int, 0)
	resMap := make(map[int]*rbac.RbacRes, 0)
	for _, v := range resList {
		if v.AppIndex != int(r.AppIndex) {
			// 其他应用的资源排除
			continue
		}
		if v.IsMenu == 0 {
			// 非菜单资源排除
			continue
		}
		if _, ok := parents[v.Pid]; !ok {
			parents[v.Pid] = []int{int(v.Id)}
		} else {
			if !collections.InArray(parents[v.Pid], int(v.Id)) {
				// 已经已包含(多个角色存在重复加载的情况),则跳过
				parents[v.Pid] = append(parents[v.Pid], int(v.Id))
			}
		}
		resMap[v.Id] = v
	}

	// 获取菜单
	root := proto.SUserMenu{}
	wg := sync.WaitGroup{}

	// 遍历函数
	//log.Println("sss", time.Now().UnixMilli())
	var f func(w *sync.WaitGroup, root *proto.SUserMenu, arr map[int][]int)
	f = func(w *sync.WaitGroup, root *proto.SUserMenu, arr map[int][]int) {
		root.Children = []*proto.SUserMenu{}
		children := arr[int(root.Id)]
		for _, c := range children {
			if root.Id == 0 && len(arr[c]) == 0 {
				// 如果一级栏目不包含下级,则不添加上级菜单
				continue
			}
			// 绑定下级菜单
			v := resMap[c]
			m := &proto.SUserMenu{
				Id:      int64(v.Id),
				Key:     v.ResKey,
				Name:    v.Name,
				Path:    v.Path,
				Icon:    v.Icon,
				SortNum: int32(v.SortNum),
			}
			root.Children = append(root.Children, m)
			w.Add(1)
			go f(w, m, arr)
		}
		w.Done()
	}
	wg.Add(1)
	f(&wg, &root, parents)
	wg.Wait()

	log.Println("ssa", time.Now().UnixMilli())
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

func walkDepartTree(node *proto.SRbacTree, nodeList []*rbac.RbacDepart) {
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
func (p *rbacServiceImpl) DepartTree(_ context.Context, req *proto.DepartTreeRequest) (*proto.SRbacTree, error) {
	root := &proto.SRbacTree{
		Id:       0,
		Label:    "根节点",
		Children: make([]*proto.SRbacTree, 0),
	}
	list := p.dao.SelectPermDept("")
	if !req.Lazy {
		walkDepartTree(root, list)
	} else {
		arr := collections.FilterArray(list, func(v *rbac.RbacDepart) bool {
			return v.Pid == int(req.ParentId)
		})
		dstArr := collections.MapList(arr, func(v *rbac.RbacDepart) *proto.SRbacTree {
			return &proto.SRbacTree{
				Id:    int64(v.Id),
				Label: v.Name,
				IsLeaf: !collections.AnyArray(list, func(r *rbac.RbacDepart) bool {
					return r.Pid == v.Id
				}),
				Children: []*proto.SRbacTree{},
			}
		})
		root.Children = dstArr
	}
	return root, nil
}

// 保存部门
func (p *rbacServiceImpl) SaveDepart(_ context.Context, r *proto.SaveDepartRequest) (*proto.SaveDepartResponse, error) {
	var dst *rbac.RbacDepart
	if r.Id > 0 {
		dst = p.dao.GetDepart(r.Id)
	} else {
		dst = &rbac.RbacDepart{}
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
		ret.Code = 1
		ret.Message = err.Error()
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

func (p *rbacServiceImpl) DeleteDepart(_ context.Context, id *proto.RbacDepartId) (*proto.TxResult, error) {
	err := p.dao.DeleteDepart(id.Value)
	return p.errorV2(err), nil
}

func (p *rbacServiceImpl) parsePermDept(v *rbac.RbacDepart) *proto.SPermDept {
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
	var dst *rbac.RbacJob
	if r.Id > 0 {
		dst = p.dao.GetJob(r.Id)
	} else {
		dst = &rbac.RbacJob{}
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
		ret.Code = 1
		ret.Message = err.Error()
	}
	return ret, nil
}

func (p *rbacServiceImpl) parsePermJob(v *rbac.RbacJob) *proto.SRbacJob {
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

func (p *rbacServiceImpl) DeleteJob(_ context.Context, id *proto.RbacJobId) (*proto.TxResult, error) {
	err := p.dao.DeleteJob(id.Value)
	return p.errorV2(err), nil
}

func (p *rbacServiceImpl) PagingJobList(_ context.Context, r *proto.RbacJobPagingRequest) (*proto.PagingRbacJobResponse, error) {
	total, rows := p.dao.QueryPagingJob(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.PagingRbacJobResponse{
		Total: int64(total),
		Rows:  make([]*proto.PagingJobList, len(rows)),
	}
	for i, v := range rows {
		ret.Rows[i] = &proto.PagingJobList{
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
	var dst *rbac.RbacUser
	if r.Id > 0 {
		dst = p.dao.GetUser(r.Id)
		if dst == nil {
			return &proto.SaveRbacUserResponse{
				Code:    2,
				Message: "no such record",
			}, nil
		}
	} else {
		dst = &rbac.RbacUser{}
		dst.Salt = util.RandString(4)
		dst.Username = r.Username
		dst.CreateTime = int(time.Now().Unix())
	}
	if l := len(r.Password); l > 0 {
		if l != 32 {
			return &proto.SaveRbacUserResponse{
				Code:    1,
				Message: "非32位md5密码",
			}, nil
		}
		dst.Password = domain.RbacPassword(r.Password, dst.Salt)
	}
	if len(r.ProfilePhoto) == 0 {
		filePath, _ := p.registryRepo.GetValue(registry.FileServerUrl)
		dst.ProfilePhoto, _ = url.JoinPath(filePath, "static/init/avatar.jpg")
	}
	dst.Flag = int(r.Flag)
	dst.ProfilePhoto = r.ProfilePhoto
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
		ret.Code = 1
		ret.Message = err.Error()
	}
	return ret, nil
}

// UpdateUserPassword implements proto.RbacServiceServer.
func (p *rbacServiceImpl) UpdateUserPassword(_ context.Context, req *proto.RbacPasswordRequest) (*proto.TxResult, error) {
	iu := p.dao.GetUser(req.UserId)
	if iu == nil {
		return p.errorV2(errors.New("no such user")), nil
	}
	if len(req.NewPassword) != 32 {
		return p.errorV2(errors.New("非32位md5密码")), nil
	}
	if l := len(req.OldPassword); l > 0 {
		if l != 32 {
			return p.errorV2(errors.New("非32位md5密码")), nil
		}
		origin := domain.RbacPassword(req.OldPassword, iu.Salt)
		if origin != iu.Password {
			return p.errorV2(errors.New("原密码不正确")), nil
		}
	}
	iu.Password = domain.RbacPassword(req.NewPassword, iu.Salt)
	_, err := p.dao.SaveUser(iu)
	return p.errorV2(err), nil
}
func (p *rbacServiceImpl) parsePermUser(v *rbac.RbacUser) *proto.SRbacUser {
	if len(v.ProfilePhoto) == 0 {
		// 如果未设置,则用系统内置头像
		prefix, _ := p.registryRepo.GetValue(registry.FileServerUrl)
		v.ProfilePhoto, _ = url.JoinPath(prefix, "static/init/avatar.jpg")
	}
	return &proto.SRbacUser{
		Id:           int64(v.Id),
		Username:     v.Username,
		Password:     v.Password,
		Flag:         int32(v.Flag),
		ProfilePhoto: v.ProfilePhoto,
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

func (p *rbacServiceImpl) DeleteUser(_ context.Context, id *proto.RbacUserId) (*proto.TxResult, error) {
	err := p.dao.DeleteUser(id.Value)
	return p.errorV2(err), nil
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
	total, rows := p.dao.QueryPagingPermUser(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.PagingRbacUserResponse{
		Total: int64(total),
		Rows:  make([]*proto.PagingUser, len(rows)),
	}
	for i, v := range rows {
		ret.Rows[i] = &proto.PagingUser{
			Id:           int64(typeconv.MustInt(v["id"])),
			Username:     typeconv.Stringify(v["username"]),
			Password:     typeconv.Stringify(v["password"]),
			Flag:         int32(typeconv.MustInt(v["flag"])),
			ProfilePhoto: typeconv.Stringify(v["profilePhoto"]),
			Nickname:     typeconv.Stringify(v["nickname"]),
			Gender:       int32(typeconv.Int(v["gender"])),
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
	var dst *rbac.RbacRole
	if r.Id > 0 {
		dst = p.dao.GetRole(r.Id)
	} else {
		dst = &rbac.RbacRole{}
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
		ret.Code = 1
		ret.Message = err.Error()
	}
	return ret, nil
}

// 更新角色资源
func (p *rbacServiceImpl) UpdateRoleResource(_ context.Context, r *proto.UpdateRoleResRequest) (*proto.TxResult, error) {
	role := p.dao.GetRole(r.RoleId)
	if role == nil {
		return p.errorV2(errors.New("角色不存在")), nil
	}
	if role.Code == "master" {
		return p.errorV2(errors.New("超级管理员已拥有全部权限")), nil
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
	mp := make(map[int]*rbac.RbacRoleRes, 0)
	for i, v := range dataList {
		old[i] = int(v.ResId)
		mp[int(v.ResId)] = v
	}
	_, deleted := util.IntArrayDiff(old, arr, func(resId int, add bool) {
		var id int64
		if !add {
			id = int64(mp[resId].Id)
		}
		p.dao.SavePermRoleRes(&rbac.RbacRoleRes{
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
	return p.errorV2(nil), nil
}

func (p *rbacServiceImpl) parsePermRole(v *rbac.RbacRole) *proto.SRbacRole {
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
func (p *rbacServiceImpl) DeletePermRole(_ context.Context, id *proto.RbacRoleId) (*proto.TxResult, error) {
	err := p.dao.DeletePermRole(id.Value)
	return p.errorV2(err), nil
}

// PagingPermRole 角色分页信息
func (p *rbacServiceImpl) PagingPermRole(_ context.Context, r *proto.RbacRolePagingRequest) (*proto.PagingRbacRoleResponse, error) {
	total, rows := p.dao.QueryPagingPermRole(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where)
	ret := &proto.PagingRbacRoleResponse{
		Total: int64(total),
		Rows:  make([]*proto.PagingPermRole, len(rows)),
	}
	for i, v := range rows {
		ret.Rows[i] = &proto.PagingPermRole{
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
func (p *rbacServiceImpl) checkParentResource(id int, currentPid int, pid int) (rbac.RbacRes, error) {
	var parentRes rbac.RbacRes
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
		var parent *rbac.RbacRes = p.dao.GetRbacResource(pid)
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
func (p *rbacServiceImpl) GenerateResourceKey(parent rbac.RbacRes) string {
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
	var dst *rbac.RbacRes
	if r.Id > 0 {
		if dst = p.dao.GetRbacResource(r.Id); dst == nil {
			return &proto.SaveRbacResResponse{
				Code:    2,
				Message: "no such data",
			}, nil
		}
	} else {
		dst = &rbac.RbacRes{}
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
				Code:    3,
				Message: "该资源(包含上级资源)路径不能以'/'开头",
			}, nil
		}
	}
	parent, err := p.checkParentResource(int(dst.Id), int(dst.Pid), int(r.Pid))
	if err != nil {
		return &proto.SaveRbacResResponse{
			Code:    2,
			Message: err.Error(),
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
		ret.Code = 1
		ret.Message = err.Error()
	} else {
		if parentChanged {
			depth := p.getResDepth(dst.Pid)
			p.updateResDepth(dst, depth)
		}
	}
	return ret, nil
}

func (p *rbacServiceImpl) parseRbacRes(v *rbac.RbacRes) *proto.SRbacRes {
	return &proto.SRbacRes{
		Id:         int64(v.Id),
		Name:       v.Name,
		ResType:    int32(v.ResType),
		Pid:        int64(v.Pid),
		Key:        v.ResKey,
		Path:       v.Path,
		Icon:       v.Icon,
		SortNum:    int32(v.SortNum),
		IsMenu:     v.IsMenu == 1,
		IsEnabled:  v.IsEnabled == 1,
		CreateTime: int64(v.CreateTime),
		AppIndex:   int32(v.AppIndex),
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

// GetResourceSQL implements proto.RbacServiceServer.
func (p *rbacServiceImpl) GetResourceSQL(_ context.Context, req *proto.PermResId) (*proto.ResourcesSQLResponse, error) {
	v := p.dao.GetRbacResource(req.Value)
	if v == nil {
		return &proto.ResourcesSQLResponse{}, nil
	}
	s := fmt.Sprintf(`INSERT INTO rbac_res `+
		`(id, name, res_type, pid, res_key, path, icon, sort_num, is_menu, is_enabled,create_time,depth,is_forbidden,app_index)`+
		` VALUES (%d,'%s',%d,%d,'%s','%s','%s',%d,%d,%d,%d,%d,%d,%d);`, v.Id,
		v.Name,
		v.ResType,
		v.Pid,
		v.ResKey,
		v.Path,
		v.Icon,
		v.SortNum,
		v.IsMenu,
		v.IsEnabled,
		v.CreateTime,
		v.Depth,
		v.IsForbidden,
		v.AppIndex,
	)
	return &proto.ResourcesSQLResponse{
		ResourceId: int64(req.Value),
		Sql:        s,
	}, nil
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
		findParent := func(pid int, arr []*rbac.RbacRes) int {
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

func (p *rbacServiceImpl) queryResChildren(parentId int, arr []*rbac.RbacRes) []*proto.SRbacRes {
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
func (p *rbacServiceImpl) DeleteRbacResource(_ context.Context, id *proto.PermResId) (*proto.TxResult, error) {
	res := p.dao.GetRbacResource(id.Value)
	if res == nil {
		return p.errorV2(errors.New("资源不存在")), nil
	}
	res.IsForbidden = 1
	_, err := p.dao.SaveRbacResource(res)
	return p.errorV2(err), nil
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
			p.dao.SaveUserRole(&rbac.RbacUserRole{
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
func (p *rbacServiceImpl) updateResDepth(dst *rbac.RbacRes, depth int) {
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
	total, rows := p.dao.QueryPagingLoginLog(int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret := &proto.LoginLogPagingResponse{
		Total: int64(total),
		Rows:  make([]*proto.PagingLoginLog, len(rows)),
	}
	for i, v := range rows {
		ret.Rows[i] = &proto.PagingLoginLog{
			Id:         int64(typeconv.MustInt(v["id"])),
			UserId:     int64(typeconv.MustInt(v["userId"])),
			Username:   typeconv.Stringify(v["username"]),
			Nickname:   typeconv.Stringify(v["nickname"]),
			Ip:         typeconv.Stringify(v["ip"]),
			IsSuccess:  int32(typeconv.MustInt(v["isSuccess"])),
			CreateTime: int64(typeconv.MustInt(v["createTime"])),
		}
		if r := ret.Rows[i]; r.UserId == 0 {
			r.Username = "master"
			r.Nickname = "超级管理员"
		}
	}
	return ret, nil
}
