package testing

import (
	"bytes"
	"github.com/ixre/gof/algorithm"
	"github.com/ixre/gof/log"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/testing/ti"
	"strconv"
	"testing"
)

var (
	modelId int32 = 5
	itemId  int64 = 1
)

// 测试商品模型设置品牌
func TestSetBrand(t *testing.T) {
	rep := ti.Factory.GetProModelRepo()
	brands := rep.SelectProBrand("")
	list := []int32{}
	for i, v := range brands {
		if true || i%2 == 0 || i == 0 {
			list = append(list, v.ID)
		}
	}
	if len(list) == 0 {
		t.Fatal("还没有产品品牌")
	}
	m := rep.CreateModel(&promodel.ProModel{
		Name:    "测试商品模型",
		Enabled: 1,
	})

	m = rep.GetModel(modelId)
	err := m.SetBrands(list)
	if err == nil {
		_, err = m.Save()
	}

	if err != nil {
		t.Error(err)
	}
	t.Log("添加模型成功")
}

// 测试商品模型添加属性
func TestModelSaveAttrs(t *testing.T) {
	rep := ti.Factory.GetProModelRepo()
	m := rep.GetModel(modelId)
	attrs := []*promodel.Attr{
		{
			Name:     "材质",
			IsFilter: 1,
			MultiChk: 1,
			Items: []*promodel.AttrItem{
				{
					Value: "棉",
				},
				{
					Value: "涤纶",
				},
				{
					Value: "布",
				},
				{
					Value: "鸭绒",
				},
			},
		},
	}
	err := m.SetAttrs(attrs)
	if err == nil {
		_, err = m.Save()
	}

	if err != nil {
		t.Error(err)
	}
	t.Log("保存属性成功")
}

// 测试商品模型添加规格
func TestModelSaveSpecs(t *testing.T) {
	rep := ti.Factory.GetProModelRepo()
	m := rep.GetModel(modelId)
	specs := []*promodel.Spec{
		{
			Name: "尺寸",
			Items: []*promodel.SpecItem{
				{
					Value: "1.5*1.8",
				},
				{
					Value: "1.8*2",
				},
				{
					Value: "2*2.3",
				},
			},
		},
	}
	err := m.SetSpecs(specs)
	if err == nil {
		_, err = m.Save()
	}

	if err != nil {
		t.Error(err)
	}
	t.Log("保存规格成功")
}

// 测试商品保存SKU
func TestItemSaveSku(t *testing.T) {
	itemRepo := ti.Factory.GetItemRepo()
	catRepo := ti.Factory.GetCategoryRepo()
	proMRepo := ti.Factory.GetProModelRepo()
	it := itemRepo.GetItem(itemId)
	if it == nil {
		t.Errorf("编号为%d的商品不存在", itemId)
	}
	iv := it.GetValue()
	catId := iv.CatId
	cat := catRepo.GetCategory(0, catId)
	if cat == nil {
		t.Errorf("编号为%d的分类不存在", catId)
	}
	//生成的规格组合
	specs := proMRepo.GetModel(cat.ProModel).Specs()
	//最多只使用2个规格
	if len(specs) > 2 {
		specs = specs[:2]
	}
	//组合SKU
	dim := make([][]interface{}, len(specs))
	result := [][]interface{}{}
	for i, v := range specs {
		dim[i] = make([]interface{}, len(v.Items))
		for iv, vv := range v.Items {
			dim[i][iv] = vv
		}
	}
	algorithm.Descartes(dim, &result)
	arr := []*item.Sku{}
	buf := bytes.NewBuffer(nil)
	for _, v := range result {
		for j, k := range v {
			if j != 0 {
				buf.WriteString(";")
			}
			ov := k.(*promodel.SpecItem)
			buf.WriteString(strconv.Itoa(int(ov.SpecId)))
			buf.WriteString(":")
			buf.WriteString(strconv.Itoa(int(ov.ID)))
		}
		arr = append(arr, &item.Sku{
			ProductId: iv.ProductId,
			// 商品编号
			ItemId: iv.ID,
			// 图片
			Image: "",
			// 规格数据
			SpecData: buf.String(),
			SpecWord: "口径2:1.2 批次数量2:50件批发",

			// 价格（分)
			Price: 100,
		})
		log.Println(buf.String())
		buf.Reset()
	}

	t.Logf("最多SKU个数:%d", len(result))

	// 保存规格
	err := it.SetSku(arr)
	if err == nil {
		_, err = it.Save()
	}
	if err != nil {
		t.Errorf("保存规格出错：%s", err.Error())
	} else {
		t.Log("保存规格成功")
	}
}
