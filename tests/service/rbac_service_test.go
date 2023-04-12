package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/dao/model"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

// 测试查询树形数据
func TestInitialTreeNode(t *testing.T) {
	list, err := impl.RbacService.QueryRbacResourceList(context.TODO(), &proto.QueryRbacResRequest{
		Keyword:   "",
		OnlyMenu:  false,
		ParentId:  0,
		InitialId: 2654,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(typeconv.MustJson(list))
}

func TestCheckRBACToken(t *testing.T) {
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOi0xLCJleHAiOjE2Nzc2NjMxODYsImlzcyI6ImdvMm8iLCJuYW1lIjoibWFzdGVyIiwic3ViIjoiZ28yby1yYmFjLXRva2VuIiwieC1wZXJtIjoibWFzdGVyLGFkbWluIn0.VMUGah8mgG8fbVicb4K45K83wbvnUZccImWMH9-vehs"
	ret, _ := impl.RbacService.CheckRBACToken(context.TODO(), &proto.RbacCheckTokenRequest{
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
	ret = gk(model.PermRes{Id: 2321, Key: "B0101"})
	t.Log("新建商户三级:", ret)
}

// 测试保存资源
func TestSaveRbacResource(t *testing.T) {
	s := impl.RbacService
	r, _ := s.GetPermRes(context.TODO(), &proto.PermResId{
		Value: 2383,
	})
	ret, _ := s.SaveRbacResource(context.TODO(), &proto.SaveRbacResRequest{
		Id:            0,
		Name:          r.Name,
		ResType:       r.ResType,
		Pid:           2321,
		Path:          r.Path,
		Icon:          r.Icon,
		SortNum:       r.SortNum,
		IsMenu:        r.IsMenu,
		IsEnabled:     r.IsEnabled,
		CreateTime:    r.CreateTime,
		ComponentName: r.ComponentName,
		Cache:         r.Cache,
	})
	t.Logf(typeconv.MustJson(ret))
}

// 测试获取用户资源
func TestGetUserResources(t *testing.T) {
	s := impl.RbacService
	ret, _ := s.GetUserResource(context.TODO(), &proto.RbacUserResourceRequest{
		UserId:   1,
		OnlyMenu: true,
	})
	t.Logf(typeconv.MustJson(ret))
}
