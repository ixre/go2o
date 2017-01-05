package rsi

import (
	"errors"
	"github.com/jsix/gof/web/ui/tree"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/service/thrift/idl/gen-go/define"
	"go2o/core/service/thrift/parser"
	"strconv"
)

// 产品服务
type productService struct {
	pmRepo promodel.IProModelRepo
	catRep product.ICategoryRepo
	proRep product.IProductRepo
}

func NewProService(pmRep promodel.IProModelRepo,
	catRep product.ICategoryRepo,
	proRep product.IProductRepo) *productService {
	return &productService{
		pmRepo: pmRep,
		catRep: catRep,
		proRep: proRep,
	}
}

// 获取产品模型
func (p *productService) GetModel(id int32) *promodel.ProModel {
	return p.pmRepo.GetProModel(id)
}

// 获取产品模型
func (p *productService) GetModels() []*promodel.ProModel {
	return p.pmRepo.SelectProModel("enabled=1")
}

// 获取属性
func (p *productService) GetAttr(id int32) *promodel.Attr {
	return p.pmRepo.GetAttr(id)
}

// 获取属性项
func (p *productService) GetAttrItem(id int32) *promodel.AttrItem {
	return p.pmRepo.GetAttrItem(id)
}

// 获取模型属性
func (p *productService) GetModelAttrs(proModel int32) []*promodel.Attr {
	m := p.pmRepo.CreateModel(&promodel.ProModel{Id: proModel})
	return m.Attrs()
}

// 获取模型属性Html
func (p *productService) GetModelAttrsHtml(proModel int32) string {
	m := p.pmRepo.CreateModel(&promodel.ProModel{Id: proModel})
	attrs := m.Attrs()
	return p.pmRepo.AttrService().AttrsHtml(attrs)
}

// 获取模型规格
func (p *productService) GetModelSpecs(proModel int32) []*promodel.Spec {
	m := p.pmRepo.CreateModel(&promodel.ProModel{Id: proModel})
	return m.Specs()
}

// 保存产品模型
func (p *productService) SaveModel(v *promodel.ProModel) (*define.Result_, error) {
	var pm promodel.IModel
	var err error
	if v.Id > 0 {
		ev := p.GetModel(v.Id)
		if ev == nil {
			err = errors.New("模型不存在")
			goto R
		}
		ev.Name = v.Name
		ev.Enabled = v.Enabled
		pm = p.pmRepo.CreateModel(ev)
	} else {
		pm = p.pmRepo.CreateModel(v)
	}
	// 保存属性
	if err == nil && v.Attrs != nil {
		err = pm.SetAttrs(v.Attrs)
	}
	// 保存规格
	if err == nil && v.Specs != nil {
		err = pm.SetSpecs(v.Specs)
	}
	// 保存品牌
	if err == nil && v.BrandArray != nil {
		err = pm.SetBrands(v.BrandArray)
	}
	// 保存模型
	if err == nil {
		v.Id, err = pm.Save()
	}
R:
	return parser.Result(v.Id, err), nil
}

// 删除产品模型
func (p *productService) DeleteProModel_(id int32) (*define.Result_, error) {
	return &define.Result_{Result_: true}, nil
}

/***** 品牌  *****/

// Get 产品品牌
func (p *productService) GetProBrand_(id int32) *promodel.ProBrand {
	return p.pmRepo.BrandService().Get(id)
}

// Save 产品品牌
func (p *productService) SaveProBrand_(v *promodel.ProBrand) (*define.Result_, error) {
	id, err := p.pmRepo.BrandService().SaveBrand(v)
	return parser.Result(id, err), nil
}

// Delete 产品品牌
func (p *productService) DeleteProBrand_(id int32) (*define.Result_, error) {
	err := p.pmRepo.BrandService().DeleteBrand(id)
	return parser.Result(0, err), nil
}

// 获取所有产品品牌
func (p *productService) GetBrands() []*promodel.ProBrand {
	return p.pmRepo.BrandService().AllBrands()
}

// 获取模型关联的产品品牌
func (p *productService) GetModelBrands(id int32) []*promodel.ProBrand {
	pm := p.pmRepo.CreateModel(&promodel.ProModel{Id: id})
	return pm.Brands()
}

/***** 分类 *****/

// 获取商品分类
func (p *productService) GetCategory(mchId, id int32) *product.Category {
	c := p.catRep.GlobCatService().GetCategory(id)
	if c != nil {
		return c.GetValue()
	}
	return nil
}

// 获取商品分类和选项
func (p *productService) GetCategoryAndOptions(mchId, id int32) (*product.Category,
	domain.IOptionStore) {
	c := p.catRep.GlobCatService().GetCategory(id)
	if c != nil {
		return c.GetValue(), c.GetOption()
	}
	return nil, nil
}

func (p *productService) DeleteCategory(mchId, id int32) error {
	return p.catRep.GlobCatService().DeleteCategory(id)
}

func (p *productService) SaveCategory(mchId int32, v *product.Category) (int32, error) {
	sl := p.catRep.GlobCatService()
	var ca product.ICategory
	if v.ID > 0 {
		ca = sl.GetCategory(v.ID)
	} else {
		ca = sl.CreateCategory(v)
	}
	if err := ca.SetValue(v); err != nil {
		return 0, err
	}
	return ca.Save()
}

func (p *productService) GetCategoryTreeNode(mchId int32) *tree.TreeNode {
	cats := p.catRep.GlobCatService().GetCategories()
	rootNode := &tree.TreeNode{
		Text:   "根节点",
		Value:  "",
		Url:    "",
		Icon:   "",
		Open:   true,
		Childs: nil}
	p.walkCategoryTree(rootNode, 0, cats)
	return rootNode
}

// 分类树形
func (p *productService) CategoryTree(parentId int32) *product.Category {
	return p.catRep.GlobCatService().CategoryTree(parentId)
}

// 获取分类关联的品牌
func (p *productService) GetCatBrands(catId int32) []*promodel.ProBrand {
	arr := p.catRep.GlobCatService().RelationBrands(catId)
	for _, v := range arr {
		v.Image = format.GetResUrl(v.Image)
	}
	return arr
}

func (p *productService) walkCategoryTree(node *tree.TreeNode, parentId int32, categories []product.ICategory) {
	node.Childs = []*tree.TreeNode{}
	for _, v := range categories {
		cate := v.GetValue()
		if cate.ParentId == parentId {
			cNode := &tree.TreeNode{
				Text:   cate.Name,
				Value:  strconv.Itoa(int(cate.ID)),
				Url:    "",
				Icon:   "",
				Open:   true,
				Childs: nil}
			node.Childs = append(node.Childs, cNode)
			p.walkCategoryTree(cNode, cate.ID, categories)
		}
	}
}

func (p *productService) GetCategories(mchId int32) []*product.Category {
	cats := p.catRep.GlobCatService().GetCategories()
	var list []*product.Category = make([]*product.Category, len(cats))
	for i, v := range cats {
		vv := v.GetValue()
		vv.Icon = format.GetResUrl(vv.Icon)
		list[i] = vv
	}
	return list
}

// 根据上级编号获取分类列表
func (p *productService) GetCategoriesByParentId(mchId, parentId int32) []*product.Category {
	cats := p.catRep.GlobCatService().GetCategories()
	list := []*product.Category{}
	for _, v := range cats {
		if vv := v.GetValue(); vv.ParentId == parentId && vv.Enabled == 1 {
			v2 := *vv
			v2.Icon = format.GetResUrl(v2.Icon)
			list = append(list, &v2)
		}
	}
	return list
}

func (p *productService) getCategoryManager(mchId int32) product.IGlobCatService {
	return p.catRep.GlobCatService()
}

func (p *productService) GetBigCategories(mchId int32) []*define.Category {
	cats := p.catRep.GlobCatService().GetCategories()
	list := []*define.Category{}
	for _, v := range cats {
		if v2 := v.GetValue(); v2.ParentId == 0 && v2.Enabled == 1 {
			v2.Icon = format.GetResUrl(v2.Icon)
			list = append(list, parser.CategoryDto(v2))
		}
	}
	return list
}

func (p *productService) GetChildCategories(mchId, parentId int32) []*define.Category {
	cats := p.catRep.GlobCatService().GetCategories()
	list := []*define.Category{}
	for _, v := range cats {
		if vv := v.GetValue(); vv.ParentId == parentId && vv.Enabled == 1 {
			vv.Icon = format.GetResUrl(vv.Icon)
			p.setChild(cats, vv)
			list = append(list, parser.CategoryDto(vv))
		}
	}
	return list
}

//
//func CopyCategory(src *product.Category, dst *dto.Category) {
//	dst.Id = src.ID
//	dst.Name = src.Name
//	dst.Level = src.Level
//	dst.Icon = src.Icon
//	dst.Url = src.CatUrl
//}

func (p *productService) setChild(list []product.ICategory, dst *product.Category) {
	for _, v := range list {
		if vv := v.GetValue(); vv.ParentId == dst.ID && vv.Enabled == 1 {
			if dst.Children == nil {
				dst.Children = []*product.Category{}
			}
			vv.Icon = format.GetResUrl(vv.Icon)
			dst.Children = append(dst.Children, vv)
		}
	}
}

/***** 产品 *****/

// 获取产品值
func (p *productService) GetProductValue(productId int32) *product.Product {
	pro := p.proRep.GetProduct(productId)
	if pro != nil {
		v := pro.GetValue()
		return &v
	}
	return nil
}

// 保存产品
func (p *productService) SaveProduct(v *product.Product) (r *define.Result_, err error) {
	var pro product.IProduct
	if v.Id > 0 {
		pro = p.proRep.GetProduct(v.Id)
		if pro == nil || pro.GetValue().VendorId != v.VendorId {
			err = product.ErrNoSuchProduct
			goto R
		}
		// 修改货品时，不会修改详情
		v.Description = pro.GetValue().Description
	} else {
		pro = p.proRep.CreateProduct(v)
	}
	// 保存
	err = pro.SetValue(v)
	if err == nil {
		// 保存属性
		if v.Attr != nil {
			err = pro.SetAttr(v.Attr)
			if err != nil {
				goto R
			}
		}
		v.Id, err = pro.Save()
	}
R:
	return parser.Result(v.Id, err), nil
}

// 保存货品描述
func (p *productService) SaveProductInfo(supplierId int32,
	productId int32, info string) error {
	pro := p.proRep.GetProduct(productId)
	if pro == nil || pro.GetValue().VendorId != supplierId {
		return product.ErrNoSuchProduct
	}
	return pro.SetDescribe(info)
}

// 删除货品
func (p *productService) DeleteItem(supplierId int32, productId int32) error {
	pro := p.proRep.GetProduct(productId)
	if pro == nil || pro.GetValue().VendorId != supplierId {
		return product.ErrNoSuchProduct
	}
	return pro.Destroy()
}

// 获取产品属性
func (p *productService) GetAttrArray(productId int32) []*product.Attr {
	pro := p.proRep.CreateProduct(&product.Product{Id: productId})
	return pro.Attr()
}

// 获取商品的销售标签
func (p *productService) GetItemSaleLabels(mchId, itemId int32) []*item.Label {
	var list = make([]*item.Label, 0)
	//todo: refactor

	//sl := s._rep.GetSale(mchId)
	//if goods := sl.ItemManager().GetItem(itemId); goods != nil {
	//	list = goods.GetSaleLabels()
	//}
	return list
}

// 保存商品的销售标签
func (p *productService) SaveItemSaleLabels(mchId, itemId int32, tagIds []int) error {
	var err error

	//todo: refactor

	//sl := s._rep.GetSale(mchId)
	//if goods := sl.ItemManager().GetItem(itemId); goods != nil {
	//	err = goods.SaveSaleLabels(tagIds)
	//} else {
	//	err = errors.New("商品不存在")
	//}
	return err
}
