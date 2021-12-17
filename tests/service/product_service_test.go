package service

import (
	"context"
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
