package service

import (
	"context"
	"encoding/json"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"testing"
)

func TestProductCategoryTree(t *testing.T) {
	parentId := 33
	node, err := impl.ProductService.GetCategoryTreeNode(context.TODO(), &proto.CategoryTreeRequest{
		ParentId:      int64(parentId),
		ExcludeIdList: nil,
		Depth:         0,
	})
	if err != nil {
		t.Error(err)
	}
	println(len(node.Value))
}

func TestCategoryInitialTreeNode(t *testing.T) {
	list, err := impl.ProductService.GetCategoryTreeNode(context.TODO(), &proto.CategoryTreeRequest{
		ParentId:  0,
		InitialId: 3,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(list.InitialList)
	bytes, _ := json.Marshal(list.Value)
	t.Log(string(bytes))
}
