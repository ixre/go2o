package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func TestProductCategoryTree(t *testing.T) {
	parentId := 33
	node, err := inject.GetProductService().GetCategoryTreeNode(context.TODO(), &proto.CategoryTreeRequest{
		ParentId:      int64(parentId),
		ExcludeIdList: nil,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(len(node.Value))
}

func TestCategoryInitialTreeNode(t *testing.T) {
	list, err := inject.GetProductService().GetCategoryTreeNode(context.TODO(), &proto.CategoryTreeRequest{
		ParentId: 0,
		Lazy:     true,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(typeconv.MustJson(list.Value))
}

func TestSourceCategories(t *testing.T) {
	list, err := inject.GetProductService().FindParentCategory(context.TODO(), &proto.CategoryIdRequest{
		CategoryId: 2041,
	})
	if err != nil {
		t.Error(err)
	}
	bytes, _ := json.Marshal(list.List)
	t.Log(string(bytes))
}

func TestGetCategoryBrands(t *testing.T) {
	list, err := inject.GetProductService().GetCategory(context.TODO(), &proto.GetCategoryRequest{
		CategoryId: 2185,
		WithBrand:  true,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(typeconv.MustJson(list))
}

func TestUpdateProductDescription(t *testing.T) {
	prod, err := inject.GetProductService().GetProduct(context.TODO(), &proto.ProductId{
		Value: 1,
	})
	if err != nil {
		t.Error(err)
	}
	prod.Description = "111" + prod.Description
	ret, _ := inject.GetProductService().SaveProduct(context.TODO(), &proto.SaveProductRequest{
		Id:                prod.Id,
		CategoryId:        prod.CategoryId,
		Name:              prod.Name,
		VendorId:          prod.VendorId,
		BrandId:           prod.BrandId,
		Code:              prod.Code,
		Image:             prod.Image,
		Description:       prod.Description,
		Remark:            prod.Remark,
		State:             prod.State,
		SortNum:           prod.SortNum,
		Attrs:             prod.Attrs,
		UpdateDescription: true,
	})
	if ret.ErrCode > 0 {
		t.Error(ret.ErrMsg)
		t.FailNow()
	}
}

func TestGetProductCategory(t *testing.T) {
	ret, _ := inject.GetProductService().GetCategory(context.TODO(), &proto.GetCategoryRequest{
		CategoryId: 2095,
		WithBrand:  true,
		WithModel:  true,
	})
	t.Log(typeconv.MustJson(ret))
}
