package domain

import (
	"encoding/json"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/gof/types/typeconv"
)

func TestSaveProduct(t *testing.T) {
	var productId int64 = 1
	repo := inject.GetProductRepo()
	prod := repo.GetProduct(productId)
	if prod == nil {
		t.Error(product.ErrNoSuchProduct)
		t.FailNow()
	}
	value := prod.GetValue()
	value.Attrs = prod.Attr()
	t.Log(typeconv.MustJson(prod.Attr()))

	json.Unmarshal([]byte(`[{"Id":17,"ProductId":1,"AttrName":"","AttrId":6,"AttrData":"17","AttrWord":"电脑"},{"Id":16,"ProductId":1,"AttrName":"","AttrId":7,"AttrData":"20,21","AttrWord":"Type-C,标准USB"}]`),
		&value.Attrs)

	err := prod.SetAttr(value.Attrs)
	if err == nil {
		_, err = prod.Save()
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(typeconv.MustJson(prod.Attr()))
}

func TestGetModelSortedSpecItems(t *testing.T) {
	var modelId int = 1
	repo := inject.GetProModelRepo()
	im := repo.GetModel(modelId)
	for _, spec := range im.Specs() {
		t.Log(spec.SortNum, spec.Name)
		//sort.Sort(spec.Items)
		for _, item := range spec.Items {
			t.Log("---", item.SortNum, item.Value)

		}
	}
}

// 测试获取产品模型
func TestGetProductModel(t *testing.T) {
	var modelId int = 4
	repo := inject.GetProModelRepo()
	im := repo.GetModel(modelId)
	arr := im.Brands()
	t.Log(typeconv.MustJson(arr))
}

// 测试销毁产品模型
func TestDestoryProductModel(t *testing.T) {
	var modelId int = 8
	repo := inject.GetProModelRepo()
	im := repo.GetModel(modelId)
	err := im.Destroy()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
