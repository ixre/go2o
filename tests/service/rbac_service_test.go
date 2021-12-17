package service

import (
	"context"
	"encoding/json"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"testing"
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
