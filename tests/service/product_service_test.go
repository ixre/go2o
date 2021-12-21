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
	})
	if err != nil {
		t.Error(err)
	}
	println(len(node.Value))
}

func TestCategoryInitialTreeNode(t *testing.T) {
	list, err := impl.ProductService.GetCategoryTreeNode(context.TODO(), &proto.CategoryTreeRequest{
		ParentId: 0,
	})
	if err != nil {
		t.Error(err)
	}
	bytes, _ := json.Marshal(list.Value)
	t.Log(string(bytes))
}

func TestSourceCategories(t *testing.T) {
	list, err := impl.ProductService.GetSourceCategories(context.TODO(), &proto.CategoryIdRequest{
		Id: 2174,
	})
	if err != nil {
		t.Error(err)
	}
	bytes, _ := json.Marshal(list.List)
	t.Log(string(bytes))
}
