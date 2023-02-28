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
	list, err := impl.RbacService.QueryResList(context.TODO(), &proto.QueryPermResRequest{
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
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOjAsImV4cCI6MTY3NzY2MTAwOCwiaXNzIjoiZ28ybyIsIm5hbWUiOiJtYXN0ZXIiLCJzdWIiOiJnbzJvLXJiYWMtdG9rZW4iLCJ4LXBlcm0iOiJtYXN0ZXIsYWRtaW4ifQ.wuPIxCjW8OjvG_n9RQhsZMxHCxO0okjQZwB3KMfna_4"
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
