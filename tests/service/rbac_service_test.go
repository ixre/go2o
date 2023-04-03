package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ixre/go2o/core/dao/model"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func TestInitialTreeNode(t *testing.T) {
	list, err := impl.RbacService.QueryResList(context.TODO(), &proto.QueryRbacResRequest{
		Keyword:   "",
		OnlyMenu:  true,
		ParentId:  0,
		InitialId: 2654,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(len(list.List))
	bytes, _ := json.Marshal(list.List)
	t.Log(string(bytes))
}

func TestCheckRBACToken(t *testing.T) {
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOi0xLCJleHAiOjE2Nzc2NjMxODYsImlzcyI6ImdvMm8iLCJuYW1lIjoibWFzdGVyIiwic3ViIjoiZ28yby1yYmFjLXRva2VuIiwieC1wZXJtIjoibWFzdGVyLGFkbWluIn0.VMUGah8mgG8fbVicb4K45K83wbvnUZccImWMH9-vehs"
	ret, _ := impl.RbacService.CheckRBACToken(context.TODO(), &proto.CheckRBACTokenRequest{
		AccessToken: accessToken,
	})
	if len(ret.Error) > 0 {
		t.Log(ret.Error)
		t.FailNow()
	}
	t.Log(typeconv.MustJson(ret))
	t.Log("用户Id", ret.UserId)
}

// 测试获取部门
func TestGetDepart(t *testing.T) {
	ret, _ := impl.RbacService.GetDepart(context.TODO(), &proto.RbacDepartId{
		Value: 4,
	})
	t.Log(typeconv.MustJson(ret))
}

// 测试获取部门
func TestGetJoinList(t *testing.T) {
	ret, _ := impl.RbacService.PagingJobList(context.TODO(), &proto.RbacJobPagingRequest{
		Params: &proto.SPagingParams{
			Begin: 0,
			End:   30,
		},
	})
	t.Log(typeconv.MustJson(ret))
}

// 测试创建新的资源Key
func TestGenerateResourceKey(t *testing.T) {
	gk := impl.RbacService.GenerateResourceKey
	ret := gk(model.PermRes{Id: 0})
	t.Log("新建一级:", ret)
	ret = gk(model.PermRes{Id: 2328, Key: "D"})
	t.Log("新建商户二级:", ret)
	ret = gk(model.PermRes{Id: 2383, Key: "D01"})
	t.Log("新建商户二级:", ret)
}

func TestSaveRbacResource(t *testing.T) {
	s := impl.RbacService
	r, _ := s.GetPermRes(context.TODO(), &proto.PermResId{
		Value: 2383,
	})
	ret, _ := s.SaveRbacResource(context.TODO(), &proto.SaveRbacResRequest{
		Id:            r.Id,
		Name:          r.Name,
		ResType:       r.ResType,
		Pid:           r.Pid,
		Key:           r.Key,
		Path:          r.Path,
		Icon:          r.Icon,
		Permission:    r.Permission,
		SortNum:       r.SortNum,
		IsExternal:    r.IsExternal,
		IsHidden:      r.IsHidden,
		CreateTime:    r.CreateTime,
		ComponentName: r.ComponentName,
		Cache:         r.Cache,
	})
	t.Logf(typeconv.MustJson(ret))
}
