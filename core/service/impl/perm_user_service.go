package impl

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : perm_user_service.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020/11/21 20:48:11
 * description :
 * history :
 */

import (
	"context"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types"
	"github.com/ixre/gof/types/typeconv"
	"go2o/core/dao"
	"go2o/core/dao/impl"
	"go2o/core/dao/model"
	"go2o/core/service/proto"
	"time"
)

var _ proto.PermUserServiceServer = new(permUserServiceImpl)

type permUserServiceImpl struct {
	dao dao.IPermUserDao
	s   storage.Interface
	serviceUtil
}

func NewPermUserService(s storage.Interface, o orm.Orm) *permUserServiceImpl {
	return &permUserServiceImpl{
		s:   s,
		dao: impl.NewPermUserDao(o),
	}
}

// 保存系统用户
func (a *permUserServiceImpl) SavePermUser(_ context.Context, r *proto.SavePermUserRequest) (*proto.SavePermUserResponse, error) {
	var dst *model.PermUser
	if r.Id > 0 {
		dst = a.dao.Get(r.Id)
	} else {
		dst = &model.PermUser{}
		dst.CreateTime = time.Now().Unix()
	}

	dst.User = r.User
	dst.Pwd = r.Pwd
	dst.Flag = int(r.Flag)
	dst.Avatar = r.Avatar
	dst.NickName = r.NickName
	dst.Sex = r.Sex
	dst.Email = r.Email
	dst.Phone = r.Phone
	dst.DeptId = r.DeptId
	dst.JobId = r.JobId
	dst.Enabled = int16(r.Enabled)
	dst.LastLogin = r.LastLogin

	id, err := a.dao.Save(dst)
	ret := &proto.SavePermUserResponse{Id: int64(id)}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

// 获取系统用户
func (a *permUserServiceImpl) GetPermUser(_ context.Context, id *proto.PermUserId) (*proto.SPermUser, error) {
	v := a.dao.Get(id.Value)
	if v == nil {
		return nil, nil
	}
	return &proto.SPermUser{
		Id:         v.Id,
		User:       v.User,
		Pwd:        v.Pwd,
		Flag:       int32(v.Flag),
		Avatar:     v.Avatar,
		NickName:   v.NickName,
		Sex:        v.Sex,
		Email:      v.Email,
		Phone:      v.Phone,
		DeptId:     v.DeptId,
		JobId:      v.JobId,
		Enabled:    int32(v.Enabled),
		LastLogin:  v.LastLogin,
		CreateTime: v.CreateTime,
	}, nil
}

func (a *permUserServiceImpl) DeletePermUser(_ context.Context, id *proto.PermUserId) (*proto.Result, error) {
	err := a.dao.Delete(id.Value)
	return a.error(err), nil
}

func (a *permUserServiceImpl) PagingShops(_ context.Context, r *proto.PermUserPagingRequest) (*proto.PermUserPagingResponse, error) {
	total, rows := a.dao.PagingQuery(int(r.Params.Begin),
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
			User:       types.Stringify(v["user"]),
			Pwd:        types.Stringify(v["pwd"]),
			Flag:       int32(typeconv.MustInt(v["flag"])),
			Avatar:     types.Stringify(v["avatar"]),
			NickName:   types.Stringify(v["nick_name"]),
			Sex:        types.Stringify(v["sex"]),
			Email:      types.Stringify(v["email"]),
			Phone:      types.Stringify(v["phone"]),
			DeptId:     int64(typeconv.MustInt(v["dept_id"])),
			JobId:      int64(typeconv.MustInt(v["job_id"])),
			Enabled:    int32(typeconv.MustInt(v["enabled"])),
			LastLogin:  int64(typeconv.MustInt(v["last_login"])),
			CreateTime: int64(typeconv.MustInt(v["create_time"])),
		}
	}
	return ret, nil
}
