package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/typeconv"
)

// 测试查询树形数据
func TestInitialTreeNode(t *testing.T) {
	list, err := inject.GetRbacService().QueryRbacResourceList(context.TODO(), &proto.QueryRbacResRequest{
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
	ret, _ := inject.GetRbacService().CheckRBACToken(context.TODO(), &proto.RbacCheckTokenRequest{
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
	ret, _ := inject.GetRbacService().GetDepart(context.TODO(), &proto.RbacDepartId{
		Value: 4,
	})
	t.Log(typeconv.MustJson(ret))
}

// 测试获取部门
func TestGetJoinList(t *testing.T) {
	ret, _ := inject.GetRbacService().PagingJobList(context.TODO(), &proto.RbacJobPagingRequest{
		Params: &proto.SPagingParams{
			Begin: 0,
			End:   30,
		},
	})
	t.Log(typeconv.MustJson(ret))
}

// 测试创建新的资源Key
func TestGenerateResourceKey(t *testing.T) {
	// gk := inject.GetRbacService().GenerateResourceKey
	// ret := gk(model.PermRes{Id: 0})
	// t.Log("新建一级:", ret)
	// ret = gk(model.PermRes{Id: 2328, ResKey: "D"})
	// t.Log("新建商户二级:", ret)
	// ret = gk(model.PermRes{Id: 2321, ResKey: "B0101"})
	// t.Log("新建商户三级:", ret)
}

// 测试保存资源
func TestSaveRbacResource(t *testing.T) {
	s := inject.GetRbacService()
	r, _ := s.GetRbacRes(context.TODO(), &proto.PermResId{
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
	})
	t.Logf(typeconv.MustJson(ret))
}

// 测试获取用户资源
func TestGetUserResources(t *testing.T) {
	s := inject.GetRbacService()
	ret, _ := s.GetUserResource(context.TODO(), &proto.RbacUserResourceRequest{
		UserId:   1,
		AppIndex: 0,
	})
	t.Logf(typeconv.MustJson(ret))
}
