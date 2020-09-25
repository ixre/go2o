package impl

import (
	"errors"
	"github.com/ixre/gof/web/ui/tree"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/service/parser"
	"go2o/core/service/proto"
	"golang.org/x/net/context"
	"strconv"
)

var _ proto.ProductServiceServer = new(productService)

// 产品服务
type productService struct {
	pmRepo  promodel.IProModelRepo
	catRepo product.ICategoryRepo
	proRepo product.IProductRepo
	serviceUtil
}


// 获取产品模型
func (p *productService) GetModel(c context.Context, id *proto.ProductModelId) (*proto.SProductModel, error) {
	im :=  p.pmRepo.GetModel(int32(id.Value))
	if im != nil{
		ret := p.parseModelDto(im.Value())
		// 绑定属性
		attrs := im.Attrs()
		ret.Attrs = make([]*proto.SProductAttr,len(attrs))
		for i,v := range attrs{
			attr := p.appendAttrItems(p.parseProductAttrDto(v),v.Items)
			ret.Attrs[i] = attr
		}
		//　绑定规格
		// 绑定品牌
		return ret, nil
	}
	return nil, nil
}

func (p *productService) appendAttrItems(attr *proto.SProductAttr, items []*promodel.AttrItem) *proto.SProductAttr {
	attr.Items = make([]*proto.SProductAttrItem, len(items))
	for i1, v1 := range items {
		attr.Items[i1] = p.parseProductAttrItemDto(v1)
	}
	return attr
}

// 获取产品模型
func (p *productService) GetModels(c context.Context, empty *proto.Empty) (*proto.ProductModelListResponse, error) {
	list := p.pmRepo.SelectProModel("enabled=1")
	arr := make([]*proto.SProductModel,len(list))
	for i,v := range list {
		arr[i] = p.parseModelDto(v)
	}
	return &proto.ProductModelListResponse{
		Value:arr,
	},nil
}

// 获取属性
func (p *productService) GetAttr(c context.Context, id *proto.ProductAttrId) (*proto.SProductAttr, error) {
	v := p.pmRepo.AttrService().GetAttr(int32(id.Value))
	if v != nil{
		attr := p.parseProductAttrDto(v)
		attr = p.appendAttrItems(attr,v.Items)
		return attr,nil
	}
	return nil,nil
}


// 删除产品品牌
func (p *productService) DeleteProBrand_(c context.Context, id *proto.Int64) (*proto.Result, error) {
	err := p.pmRepo.BrandService().DeleteBrand(int32(id.Value))
	return p.result(err), nil
}


// 获取所有产品品牌
func (p *productService) GetBrands(c context.Context, empty *proto.Empty) (*proto.ProductBrandListResponse, error) {
	list := p.pmRepo.BrandService().AllBrands()
	arr := make([]*proto.SProductBrand,len(list))
	for i,v := range list{
		arr[i] = p.parseBrandDto(v)
	}
	return &proto.ProductBrandListResponse{
		Value:arr,
	},nil
}

func (p *productService) GetCategories(c context.Context, empty *proto.Empty) (*proto.ProductCategoriesResponse, error) {
	ic := p.catRepo.GlobCatService()
	cats := ic.GetCategories()
	list := make([]*proto.SProductCategory, len(cats))
	for i, v := range cats {
		cat := p.parseCategoryDto(v.GetValue())
		cat.Icon = format.GetResUrl(cat.Icon)
		p.appendCategoryBrands(ic, v, cat)
		list[i] = cat
	}
	return &proto.ProductCategoriesResponse{
		Value: list,
	}, nil
}

// 为分类绑定品牌
func (p *productService) appendCategoryBrands(ic product.IGlobCatService, v product.ICategory, cat *proto.SProductCategory) {
	brands := ic.RelationBrands(v.GetDomainId())
	cat.Brands = make([]*proto.SProductBrand, len(brands))
	for i1, v1 := range brands {
		cat.Brands[i1] = p.parseBrandDto(v1)
	}
}


// 获取商品分类
func (p *productService) GetCategory(c context.Context, id *proto.Int64) (*proto.SProductCategory, error) {
	ic := p.catRepo.GlobCatService()
	v := ic.GetCategory(int(id.Value))
	if c != nil {
		cat := p.parseCategoryDto(v.GetValue())
		p.appendCategoryBrands(ic, v, cat)
		return cat,nil
	}
	return nil,nil
}

// 删除分类
func (p *productService) DeleteCategory(c context.Context, id *proto.Int64) (*proto.Result, error) {
	err := p.catRepo.GlobCatService().DeleteCategory(int(id.Value))
	return p.error(err),nil
}

// 保存分类
func (p *productService) SaveCategory(c context.Context, category *proto.SProductCategory) (*proto.Result, error) {
	sl := p.catRepo.GlobCatService()
	var ca product.ICategory
	v := p.parseCategory(category)
	if v.Id > 0 {
		ca = sl.GetCategory(v.Id)
	} else {
		ca = sl.CreateCategory(v)
	}
	err := ca.SetValue(v)
	if err == nil {
		_, err= ca.Save()
	}
	return p.error(err),nil
}


// 根据上级编号获取分类列表
func (p *productService) GetChildren(c context.Context, id *proto.CategoryParentId) (*proto.ProductCategoriesResponse, error) {
	ic := p.catRepo.GlobCatService()
	cats := ic.GetCategories()
	var list = make([]*proto.SProductCategory, 0)
	for _, v := range cats {
		if vv := v.GetValue(); vv.ParentId == int(id.Value) && vv.Enabled == 1 {
			cat := p.parseCategoryDto(v.GetValue())
			cat.Icon = format.GetResUrl(cat.Icon)
			p.appendCategoryBrands(ic, v, cat)
			list = append(list, cat)
		}
	}
	return &proto.ProductCategoriesResponse{
		Value: list,
	}, nil
}

// 获取产品值
func (p *productService) GetProductValue(c context.Context, id *proto.ProductId) (*proto.SProduct, error) {
	pro := p.proRepo.GetProduct(id.Value)
	if pro != nil {
		v := p.parseProductDto(pro.GetValue())
		return v,nil
	}
	return nil,nil
}


// 保存产品
func (p *productService) SaveProduct(c context.Context, r *proto.SProduct) (*proto.SaveProductResponse, error) {
	var pro product.IProduct
	v := p.parseProduct(r)
	ret := &proto.SaveProductResponse{}
	if v.Id > 0 {
		pro = p.proRepo.GetProduct(v.Id)
		if pro == nil || pro.GetValue().VendorId != v.VendorId {
			ret.ErrCode = 1
			ret.ErrMsg = product.ErrNoSuchProduct.Error()
			return ret,nil
		}
		// 修改货品时，不会修改详情
		v.Description = pro.GetValue().Description
	} else {
		pro = p.proRepo.CreateProduct(v)
	}
	// 保存
	err := pro.SetValue(v)
	if err == nil {
		// 保存属性
		if v.Attr != nil {
			err = pro.SetAttr(v.Attr)
		}
		if err == nil {
			v.Id, err = pro.Save()
		}
	}
	if err != nil {
		ret.ErrMsg = err.Error()
		ret.ErrCode =2
	}
	ret.ProductId = v.Id
	return ret,nil
}

// 保存货品描述
func (p *productService) SaveProductInfo(c context.Context, r *proto.ProductInfoRequest) (*proto.Result, error) {
	pro := p.proRepo.GetProduct(r.ProductId)
	var err error
	if pro == nil{
		err = product.ErrNoSuchProduct
	}else{
		err = pro.SetDescribe(r.Info)
	}
	return p.error(err),nil
}

// 获取模型属性
func (p *productService) GetModelAttrs_(proModel int32) []*promodel.Attr {
	m := p.pmRepo.CreateModel(&promodel.ProModel{ID: proModel})
	return m.Attrs()
}

func (p *productService) GetAttrItem_(c context.Context, i *proto.Int64) (*proto.SProductAttrItem, error) {
	panic("implement me")
}


// 获取模型属性Html
func (p *productService) GetModelAttrsHtml(c context.Context, id *proto.Int64) (*proto.String, error) {
	m := p.pmRepo.CreateModel(&promodel.ProModel{ID: int32(id.Value)})
	attrs := m.Attrs()
	s := p.pmRepo.AttrService().AttrsHtml(attrs)
	return &proto.String{Value:s},nil
}


// 保存产品模型
func (p *productService) SaveModel(c context.Context, r *proto.SProductModel) (*proto.Result, error) {
	var pm promodel.IModel
	v := p.parseProductModel(r)
	var err error
	if v.ID > 0 {
		pm = p.pmRepo.GetModel(int32(r.Id))
		if pm == nil {
			err = errors.New("模型不存在")
			return p.error(err),nil
		}
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
		v.ID, err = pm.Save()
	}
	return p.result(err), nil
}

// 删除产品模型
func (p *productService) DeleteProModel_(c context.Context, id *proto.Int64) (*proto.Result, error) {
	//err := p.pmRepo.DeleteProModel(id)
	//todo: 暂时不允许删除模型
	return p.result(errors.New("暂时不允许删除模型")), nil
}


// 获取产品品牌
func (p *productService) GetProBrand_(c context.Context, id *proto.Int64) (*proto.SProductBrand, error) {
	brand := p.pmRepo.BrandService().Get(int32(id.Value))
	if brand != nil {
		return p.parseBrandDto(brand),nil
	}
	return nil,nil
}

func (p *productService) SaveProBrand_(c context.Context, brand *proto.SProductBrand) (*proto.Result, error) {
	panic("implement me")
}


func NewProService(pmRepo promodel.IProModelRepo,
	catRepo product.ICategoryRepo,
	proRepo product.IProductRepo) *productService {
	return &productService{
		pmRepo:  pmRepo,
		catRepo: catRepo,
		proRepo: proRepo,
	}
}

// 删除产品
func (p *productService) DeleteProduct(_ context.Context, r *proto.DeleteProductRequest) (*proto.Result, error) {
	var err error
	prod := p.proRepo.GetProduct(r.ProductId)
	if prod == nil || prod.GetValue().VendorId != r.SellerId {
		err = product.ErrNoSuchProduct
	} else {
		err = p.proRepo.DeleteProduct(r.ProductId)
		//todo: 删除商品
	}
	return p.error(err), nil
}


// 获取属性项
func (p *productService) GetAttrItem(id int32) *promodel.AttrItem {
	return p.pmRepo.GetAttrItem(id)
}



// 获取模型规格
func (p *productService) GetModelSpecs(proModel int32) []*promodel.Spec {
	m := p.pmRepo.CreateModel(&promodel.ProModel{ID: proModel})
	return m.Specs()
}



/***** 品牌  *****/

// Save 产品品牌
func (p *productService) SaveProBrand_(v *promodel.ProBrand) (*proto.Result, error) {
	_, err := p.pmRepo.BrandService().SaveBrand(v)
	return p.result(err), nil
}



// 获取模型关联的产品品牌
func (p *productService) GetModelBrands(id int32) []*promodel.ProBrand {
	pm := p.pmRepo.CreateModel(&promodel.ProModel{ID: id})
	return pm.Brands()
}

/***** 分类 *****/


// 获取商品分类和选项
func (p *productService) GetCategoryAndOptions(mchId int64, id int32) (*product.Category,
	domain.IOptionStore) {
	c := p.catRepo.GlobCatService().GetCategory(int(id))
	if c != nil {
		return c.GetValue(), c.GetOption()
	}
	return nil, nil
}



func (p *productService) GetCategoryTreeNode(mchId int64) *tree.TreeNode {
	cats := p.catRepo.GlobCatService().GetCategories()
	rootNode := &tree.TreeNode{
		Title:    "根节点",
		Value:    "",
		Url:      "",
		Icon:     "",
		Expand:   true,
		Children: nil}
	p.walkCategoryTree(rootNode, 0, cats)
	return rootNode
}

// 分类树形
func (p *productService) CategoryTree(parentId int32) *product.Category {
	return p.catRepo.GlobCatService().CategoryTree(int(parentId))
}

// 获取分类关联的品牌
func (p *productService) GetCatBrands(catId int32) []*promodel.ProBrand {
	arr := p.catRepo.GlobCatService().RelationBrands(int(catId))
	for _, v := range arr {
		v.Image = format.GetResUrl(v.Image)
	}
	return arr
}

func (p *productService) walkCategoryTree(node *tree.TreeNode, parentId int, categories []product.ICategory) {
	node.Children = []*tree.TreeNode{}
	for _, v := range categories {
		cate := v.GetValue()
		if cate.ParentId == int(parentId) {
			cNode := &tree.TreeNode{
				Title:    cate.Name,
				Value:    strconv.Itoa(int(cate.Id)),
				Url:      "",
				Icon:     "",
				Expand:   false,
				Children: nil}
			node.Children = append(node.Children, cNode)
			p.walkCategoryTree(cNode, cate.Id, categories)
		}
	}
}

func (p *productService) getCategoryManager(mchId int64) product.IGlobCatService {
	return p.catRepo.GlobCatService()
}

func (p *productService) GetBigCategories(mchId int64) []*proto.SCategory {
	cats := p.catRepo.GlobCatService().GetCategories()
	var list []*proto.SCategory
	for _, v := range cats {
		if v2 := v.GetValue(); v2.ParentId == 0 && v2.Enabled == 1 {
			v2.Icon = format.GetResUrl(v2.Icon)
			list = append(list, parser.CategoryDto(v2))
		}
	}
	return list
}

func (p *productService) GetChildCategories(mchId int64, parentId int64) []*proto.SCategory {
	cats := p.catRepo.GlobCatService().GetCategories()
	var list []*proto.SCategory
	for _, v := range cats {
		if vv := v.GetValue(); vv.ParentId == int(parentId) && vv.Enabled == 1 {
			vv.Icon = format.GetResUrl(vv.Icon)
			p.setChild(cats, vv)
			list = append(list, parser.CategoryDto(vv))
		}
	}
	return list
}

//
//func CopyCategory(src *product.Category, dst *dto.Category) {
//	dst.ID = src.ID
//	dst.Name = src.Name
//	dst.Level = src.Level
//	dst.Icon = src.Icon
//	dst.Url = src.CatUrl
//}

func (p *productService) setChild(list []product.ICategory, dst *product.Category) {
	for _, v := range list {
		if vv := v.GetValue(); vv.ParentId == dst.Id && vv.Enabled == 1 {
			if dst.Children == nil {
				dst.Children = []*product.Category{}
			}
			vv.Icon = format.GetResUrl(vv.Icon)
			dst.Children = append(dst.Children, vv)
		}
	}
}

/***** 产品 *****/





// 删除货品
func (p *productService) DeleteItem(supplierId int64, productId int64) error {
	pro := p.proRepo.GetProduct(productId)
	if pro == nil || pro.GetValue().VendorId != supplierId {
		return product.ErrNoSuchProduct
	}
	return pro.Destroy()
}

// 获取产品属性
func (p *productService) GetAttrArray(productId int64) []*product.Attr {
	pro := p.proRepo.CreateProduct(&product.Product{Id: productId})
	return pro.Attr()
}

// 获取商品的销售标签
func (p *productService) GetItemSaleLabels(mchId int64, itemId int64) []*item.Label {
	var list = make([]*item.Label, 0)
	//todo: refactor

	//sl := s._rep.GetSale(mchId)
	//if goods := sl.ItemManager().GetItem(itemId); goods != nil {
	//	list = goods.GetSaleLabels()
	//}
	return list
}

// 保存商品的销售标签
func (p *productService) SaveItemSaleLabels(mchId, itemId int64, tagIds []int) error {
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

func (p *productService) parseModelDto(v *promodel.ProModel)*proto.SProductModel {

}

func (p *productService) parseProductAttrDto(v *promodel.Attr) *proto.SProductAttr {

}

func (p *productService) parseProductAttrItemDto(v1 *promodel.AttrItem) *proto.SProductAttrItem {

}

func (p *productService) parseBrandDto(v *promodel.ProBrand) *proto.SProductBrand {

}

func (p *productService) parseCategoryDto(value *product.Category) *proto.SProductCategory {

}

func (p *productService) parseCategory(category *proto.SProductCategory)*product.Category{

}

func (p *productService) parseProductDto(value product.Product) *proto.SProduct {

}

func (p *productService) parseProduct(r *proto.SProduct) *product.Product {

}

func (p *productService) parseProductModel(r *proto.SProductModel)*promodel.ProModel {

}
