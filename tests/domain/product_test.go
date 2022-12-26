package domain

import (
	"encoding/json"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/tests/ti"
	"github.com/ixre/gof/types/typeconv"
)

func TestSaveProduct(t *testing.T){
	var productId int64 = 1
	repo := ti.Factory.GetProductRepo()
	prod := repo.GetProduct(productId)
	if prod == nil{
		t.Error(product.ErrNoSuchProduct)
		t.FailNow()
	}
	value := prod.GetValue()
	value.Attrs = prod.Attr()
	t.Log(typeconv.MustJson(prod.Attr()))

	json.Unmarshal([]byte(`[ { "id": 16, "attrId": 7, "attrData": "NaN,20,22", "attrWord": "雷电接口,Type-C,标准USB" }, { "id": 17, "attrId": 6, "attrData": "18", "attrWord": "电脑" } ]`),
	&value.Attrs)

	err := prod.SetAttr(value.Attrs)
	if err == nil{
		_,err = prod.Save()
	}
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	t.Log(typeconv.MustJson(prod.Attr()))
}