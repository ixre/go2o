package service

import (
	"context"
	"encoding/json"
	"testing"

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
