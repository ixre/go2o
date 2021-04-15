package impl

import (
	"errors"
	"github.com/ixre/gof/types"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/infrastructure/format"
	"go2o/core/service/proto"
	"golang.org/x/net/context"
)

var _ proto.ProductServiceServer = new(productService)

// 产品服务
type productService struct {
	pmRepo      promodel.IProductModelRepo
	catRepo     product.ICategoryRepo
	productRepo product.IProductRepo
	serviceUtil
}

func NewProductService(pmRepo promodel.IProductModelRepo,
	catRepo product.ICategoryRepo,
	proRepo product.IProductRepo) *productService {
	return &productService{
		pmRepo:      pmRepo,
		catRepo:     catRepo,
		productRepo: proRepo,
	}
}

// GetModel 获取产品模型
func (p *productService) GetModel(_ context.Context, id *proto.ProductModelId) (*proto.SProductModel, error) {
	im := p.pmRepo.GetModel(int32(id.Value))
	if im != nil {
		ret := p.parseModelDto(im.Value())
		// 绑定属性
		attrs := im.Attrs()
		ret.Attrs = make([]*proto.SProductAttr, len(attrs))
		for i, v := range attrs {
			attr := p.appendAttrItems(p.parseProductAttrDto(v), v.Items)
			ret.Attrs[i] = attr
		}
		// 绑定规格
		specList := im.Specs()
		ret.Specs = make([]*proto.SProductSpec, len(specList))
		for i, v := range specList {
			spec := p.appendSpecItems(p.parseSpecDto(v), v.Items)
			ret.Specs[i] = spec
		}
		// 绑定品牌
		brands := im.Brands()
		ret.Brands = make([]*proto.SProductBrand, len(brands))
		for i, v := range brands {
			ret.Brands[i] = p.parseBrandDto(v)
		}
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

func (p *productService) appendSpecItems(spec *proto.SProductSpec, items promodel.SpecItemList) *proto.SProductSpec {
	spec.Items = make([]*proto.SProductSpecItem, len(items))
	for i1, v1 := range items {
		spec.Items[i1] = p.parseProductSpecItemDto(v1)
	}
	return spec
}

// GetModels 获取产品模型
func (p *productService) GetModels(_ context.Context, _ *proto.Empty) (*proto.ProductModelListResponse, error) {
	list := p.pmRepo.SelectProModel("")
	arr := make([]*proto.SProductModel, len(list))
	for i, v := range list {
		arr[i] = p.parseModelDto(v)
	}
	return &proto.ProductModelListResponse{
		Value: arr,
	}, nil
}

// GetAttr 获取属性
func (p *productService) GetAttr(_ context.Context, id *proto.ProductAttrId) (*proto.SProductAttr, error) {
	v := p.pmRepo.AttrService().GetAttr(int32(id.Value))
	if v != nil {
		attr := p.parseProductAttrDto(v)
		attr = p.appendAttrItems(attr, v.Items)
		return attr, nil
	}
	return nil, nil
}

// 获取属性项
func (p *productService) GetAttrItem(_ context.Context, id *proto.ProductAttrItemId) (*proto.SProductAttrItem, error) {
	it := p.pmRepo.GetAttrItem(id.Value)
	if it != nil {
		return p.parseProductAttrItemDto(it), nil
	}
	return nil, nil
}

// 获取所有产品品牌
func (p *productService) GetBrands(_ context.Context, _ *proto.Empty) (*proto.ProductBrandListResponse, error) {
	list := p.pmRepo.BrandService().AllBrands()
	arr := make([]*proto.SProductBrand, len(list))
	for i, v := range list {
		arr[i] = p.parseBrandDto(v)
	}
	return &proto.ProductBrandListResponse{
		Value: arr,
	}, nil
}

func (p *productService) GetCategories(_ context.Context, _ *proto.Empty) (*proto.ProductCategoriesResponse, error) {
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
func (p *productService) GetCategory(_ context.Context, id *proto.Int64) (*proto.SProductCategory, error) {
	ic := p.catRepo.GlobCatService()
	v := ic.GetCategory(int(id.Value))
	if v != nil {
		cat := p.parseCategoryDto(v.GetValue())
		p.appendCategoryBrands(ic, v, cat)
		return cat, nil
	}
	return nil, nil
}

// 删除分类
func (p *productService) DeleteCategory(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	err := p.catRepo.GlobCatService().DeleteCategory(int(id.Value))
	return p.error(err), nil
}

// 保存分类
func (p *productService) SaveCategory(_ context.Context, category *proto.SProductCategory) (*proto.Result, error) {
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
		_, err = ca.Save()
	}
	return p.error(err), nil
}

// 根据上级编号获取分类列表
func (p *productService) GetChildren(_ context.Context, id *proto.CategoryParentId) (*proto.ProductCategoriesResponse, error) {
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
func (p *productService) GetProduct(_ context.Context, id *proto.ProductId) (*proto.SProduct, error) {
	pro := p.productRepo.GetProduct(id.Value)
	if pro != nil {
		ret := p.parseProductDto(pro.GetValue())
		attrs := pro.Attr()
		ret.Attrs = make([]*proto.SProductAttrValue, len(attrs))
		for i, v := range attrs {
			ret.Attrs[i] = p.parseProductAttrValueDto(v)
		}
		return ret, nil
	}
	return nil, nil
}

// 保存产品
func (p *productService) SaveProduct(_ context.Context, r *proto.SProduct) (*proto.SaveProductResponse, error) {
	var pro product.IProduct
	v := p.parseProduct(r)
	ret := &proto.SaveProductResponse{}
	if v.Id > 0 {
		pro = p.productRepo.GetProduct(v.Id)
		if pro == nil || pro.GetValue().VendorId != v.VendorId {
			ret.ErrCode = 1
			ret.ErrMsg = product.ErrNoSuchProduct.Error()
			return ret, nil
		}
		// 修改货品时，不会修改详情
		v.Description = pro.GetValue().Description
	} else {
		pro = p.productRepo.CreateProduct(v)
	}
	// 保存
	err := pro.SetValue(v)
	if err == nil {
		// 保存属性
		if v.Attrs != nil {
			err = pro.SetAttr(v.Attrs)
		}
		if err == nil {
			v.Id, err = pro.Save()
		}
	}
	if err != nil {
		ret.ErrMsg = err.Error()
		ret.ErrCode = 2
	}
	ret.ProductId = v.Id
	return ret, nil
}

// 保存货品描述
func (p *productService) SaveProductInfo(_ context.Context, r *proto.ProductInfoRequest) (*proto.Result, error) {
	pro := p.productRepo.GetProduct(r.ProductId)
	var err error
	if pro == nil {
		err = product.ErrNoSuchProduct
	} else {
		err = pro.SetDescribe(r.Info)
	}
	return p.error(err), nil
}

// 获取模型属性
func (p *productService) GetModelAttrs_(proModel int32) []*promodel.Attr {
	m := p.pmRepo.CreateModel(&promodel.ProductModel{ID: proModel})
	return m.Attrs()
}

// 获取模型属性Html
func (p *productService) GetModelAttrsHtml(_ context.Context, id *proto.ProductModelId) (*proto.String, error) {
	m := p.pmRepo.CreateModel(&promodel.ProductModel{ID: int32(id.Value)})
	attrs := m.Attrs()
	s := p.pmRepo.AttrService().AttrsHtml(attrs)
	return &proto.String{Value: s}, nil
}

// 保存产品模型
func (p *productService) SaveModel(_ context.Context, r *proto.SProductModel) (*proto.Result, error) {
	var pm promodel.IProductModel
	v := p.parseProductModel(r)
	if v.ID > 0 {
		pm = p.pmRepo.GetModel(int32(r.Id))
		if pm == nil {
			return p.error(errors.New("模型不存在")), nil
		}
	} else {
		pm = p.pmRepo.CreateModel(v)
	}
	err := pm.SetValue(v)
	if err == nil {
		// 保存属性
		if v.Attrs != nil {
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
	}
	// 保存模型
	if err == nil {
		v.ID, err = pm.Save()
	}
	return p.result(err), nil
}

// 删除产品模型
func (p *productService) DeleteModel_(_ context.Context, id *proto.ProductModelId) (*proto.Result, error) {
	//err := p.pmRepo.DeleteProModel(id)
	//todo: 暂时不允许删除模型
	return p.result(errors.New("暂时不允许删除模型")), nil
}

// 获取产品品牌
func (p *productService) GetBrand(_ context.Context, id *proto.Int64) (*proto.SProductBrand, error) {
	brand := p.pmRepo.BrandService().Get(int32(id.Value))
	if brand != nil {
		return p.parseBrandDto(brand), nil
	}
	return nil, nil
}

// Save 产品品牌
func (p *productService) SaveBrand(_ context.Context, brand *proto.SProductBrand) (*proto.Result, error) {
	v := p.parseBrand(brand)
	_, err := p.pmRepo.BrandService().SaveBrand(v)
	return p.result(err), nil
}

// 删除产品品牌
func (p *productService) DeleteBrand(_ context.Context, id *proto.Int64) (*proto.Result, error) {
	err := p.pmRepo.BrandService().DeleteBrand(int32(id.Value))
	return p.result(err), nil
}

// 删除产品
func (p *productService) DeleteProduct(_ context.Context, r *proto.DeleteProductRequest) (*proto.Result, error) {
	var err error
	prod := p.productRepo.GetProduct(r.ProductId)
	if prod == nil || prod.GetValue().VendorId != r.SellerId {
		err = product.ErrNoSuchProduct
	} else {
		err = p.productRepo.DeleteProduct(r.ProductId)
		//todo: 删除商品
	}
	return p.error(err), nil
}

// 获取模型规格
func (p *productService) GetModelSpecs(proModel int32) []*promodel.Spec {
	m := p.pmRepo.CreateModel(&promodel.ProductModel{ID: proModel})
	return m.Specs()
}

// GetCategoryTreeNode 分类
func (p *productService) GetCategoryTreeNode(_ context.Context, req *proto.CategoryTreeRequest) (*proto.STreeNode, error) {
	cats := p.catRepo.GlobCatService().GetCategories()
	rootNode := &proto.STreeNode{
		Title:    "根节点",
		Value:    "",
		Icon:     "",
		Expand:   true,
		Children: nil}
	p.walkCategoryTree(rootNode, 0, cats, 0, req)
	return rootNode, nil
}

// 获取分类关联的品牌
func (p *productService) GetCatBrands(catId int32) []*promodel.ProductBrand {
	arr := p.catRepo.GlobCatService().RelationBrands(int(catId))
	for _, v := range arr {
		v.Image = format.GetResUrl(v.Image)
	}
	return arr
}

// 排除分类
func (p *productService) testWalkCondition(req *proto.CategoryTreeRequest, cat *product.Category, depth int) bool {
	if req.Depth > 0 && int(req.Depth) < depth+1 {
		return false
	}
	if req.ExcludeIdList == nil {
		return true
	}
	for _, v := range req.ExcludeIdList {
		if v == int64(cat.Id) {
			return false
		}
	}
	return true
}

func (p *productService) walkCategoryTree(node *proto.STreeNode, parentId int,
	categories []product.ICategory, depth int,
	req *proto.CategoryTreeRequest) {
	node.Children = []*proto.STreeNode{}
	// 遍历子分类
	for _, v := range categories {
		cat := v.GetValue()
		if cat.ParentId == parentId &&
			p.testWalkCondition(req, cat, depth) {
			cNode := &proto.STreeNode{
				Id:       int64(cat.Id),
				Title:    cat.Name,
				Icon:     "",
				Expand:   false,
				Children: nil}
			node.Children = append(node.Children, cNode)
			p.walkCategoryTree(cNode, cat.Id, categories, depth+1, req)
		}
	}
}

func (p *productService) GetBigCategories(mchId int64) []*proto.SProductCategory {
	cats := p.catRepo.GlobCatService().GetCategories()
	var list []*proto.SProductCategory
	for _, v := range cats {
		if v2 := v.GetValue(); v2.ParentId == 0 && v2.Enabled == 1 {
			v2.Icon = format.GetResUrl(v2.Icon)
			list = append(list, p.parseCategoryDto(v2))
		}
	}
	return list
}

func (p *productService) GetChildCategories(mchId int64, parentId int64) []*proto.SProductCategory {
	cats := p.catRepo.GlobCatService().GetCategories()
	var list []*proto.SProductCategory
	for _, v := range cats {
		if vv := v.GetValue(); vv.ParentId == int(parentId) && vv.Enabled == 1 {
			vv.Icon = format.GetResUrl(vv.Icon)
			p.setChild(cats, vv)
			list = append(list, p.parseCategoryDto(vv))
		}
	}
	return list
}

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
	pro := p.productRepo.GetProduct(productId)
	if pro == nil || pro.GetValue().VendorId != supplierId {
		return product.ErrNoSuchProduct
	}
	return pro.Destroy()
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

func (p *productService) parseModelDto(v *promodel.ProductModel) *proto.SProductModel {
	ret := &proto.SProductModel{
		Id:      int64(v.ID),
		Name:    v.Name,
		AttrStr: v.AttrStr,
		SpecStr: v.SpecStr,
		Attrs:   nil,
		Specs:   nil,
		Brands:  nil,
		Enabled: int32(v.Enabled),
	}
	return ret
}

func (p *productService) parseProductAttrDto(v *promodel.Attr) *proto.SProductAttr {
	return &proto.SProductAttr{
		Id:         int64(v.Id),
		Name:       v.Name,
		IsFilter:   v.IsFilter,
		MultiCheck: v.MultiChk,
		SortNum:    v.SortNum,
		ItemValues: v.ItemValues,
		Items:      nil,
	}
}

func (p *productService) parseProductAttrItemDto(v *promodel.AttrItem) *proto.SProductAttrItem {
	return &proto.SProductAttrItem{
		Id:      int64(v.Id),
		Value:   v.Value,
		SortNum: v.SortNum,
	}
}

func (p *productService) parseBrandDto(v *promodel.ProductBrand) *proto.SProductBrand {
	return &proto.SProductBrand{
		Id:           int64(v.ID),
		Name:         v.Name,
		Image:        v.Image,
		SiteUrl:      v.SiteUrl,
		Introduce:    v.Introduce,
		ReviewState:  v.ReviewState,
		ReviewRemark: v.ReviewRemark,
		Enabled:      int32(v.Enabled),
		CreateTime:   v.CreateTime,
	}
}

func (p *productService) parseCategoryDto(v *product.Category) *proto.SProductCategory {
	return &proto.SProductCategory{
		Id:          int64(v.Id),
		ParentId:    int64(v.ParentId),
		ModelId:     int64(v.ModelId),
		Priority:    int32(v.Priority),
		Name:        v.Name,
		IsVirtual:   v.VirtualCat == 1,
		CategoryUrl: v.CatUrl,
		Icon:        v.Icon,
		IconPoint:   v.IconPoint,
		Level:       int32(v.Level),
		SortNum:     int32(v.SortNum),
		FloorShow:   v.FloorShow == 1,
		Enabled:     v.Enabled == 1,
		CreateTime:  v.CreateTime,
		Options:     map[string]string{},
		Brands:      nil,
		Children:    nil,
	}
}

func (p *productService) parseCategory(v *proto.SProductCategory) *product.Category {
	return &product.Category{
		Id:         int(v.Id),
		ParentId:   int(v.ParentId),
		ModelId:    int(v.ModelId),
		Priority:   int(v.Priority),
		Name:       v.Name,
		VirtualCat: types.IntCond(v.IsVirtual, 1, 0),
		CatUrl:     v.CategoryUrl,
		Icon:       v.Icon,
		IconPoint:  v.IconPoint,
		Level:      int(v.Level),
		SortNum:    int(v.SortNum),
		FloorShow:  types.IntCond(v.FloorShow, 1, 0),
		Enabled:    types.IntCond(v.Enabled, 1, 0),
	}
}

func (p *productService) parseProductDto(v product.Product) *proto.SProduct {
	return &proto.SProduct{
		Id:          v.Id,
		CategoryId:  int64(v.CatId),
		Name:        v.Name,
		VendorId:    v.VendorId,
		BrandId:     int64(v.BrandId),
		Code:        v.Code,
		Image:       v.Image,
		Description: v.Description,
		Remark:      v.Remark,
		State:       v.State,
		SortNum:     v.SortNum,
		CreateTime:  v.CreateTime,
		UpdateTime:  v.UpdateTime,
	}
}

func (p *productService) parseProduct(v *proto.SProduct) *product.Product {
	ret := &product.Product{
		Id:          v.Id,
		CatId:       int32(v.CategoryId),
		Name:        v.Name,
		VendorId:    v.VendorId,
		BrandId:     int32(v.BrandId),
		Code:        v.Code,
		Image:       v.Image,
		Description: v.Description,
		Remark:      v.Remark,
		State:       v.State,
		SortNum:     v.SortNum,
	}
	if v.Attrs != nil {
		ret.Attrs = make([]*product.AttrValue, len(v.Attrs))
		for i, v := range v.Attrs {
			ret.Attrs[i] = p.parseProductAttrValue(v)
		}
	}
	return ret
}

func (p *productService) parseProductModel(v *proto.SProductModel) *promodel.ProductModel {
	ret := &promodel.ProductModel{
		ID:      int32(v.Id),
		Name:    v.Name,
		Enabled: int(v.Enabled),
	}
	if v.Attrs != nil {
		ret.Attrs = make([]*promodel.Attr, len(v.Attrs))
		for i, v := range v.Attrs {
			ret.Attrs[i] = p.parseProductAttr(v)
		}
	}
	if v.Specs != nil {
		ret.Specs = make([]*promodel.Spec, len(v.Specs))
		for i, v := range v.Specs {
			ret.Specs[i] = p.parseProductSpec(v)
		}
	}
	if v.Brands != nil {
		ret.BrandArray = make([]int32, len(v.Brands))
		for i, v := range v.Brands {
			ret.BrandArray[i] = int32(v.Id)
		}
	}
	return ret
}

func (p *productService) parseBrand(v *proto.SProductBrand) *promodel.ProductBrand {
	return &promodel.ProductBrand{
		ID:           int32(v.Id),
		Name:         v.Name,
		Image:        v.Image,
		SiteUrl:      v.SiteUrl,
		Introduce:    v.Introduce,
		ReviewState:  v.ReviewState,
		ReviewRemark: v.ReviewRemark,
		Enabled:      int(v.Enabled),
		CreateTime:   v.CreateTime,
	}
}

func (p *productService) parseProductAttr(v *proto.SProductAttr) *promodel.Attr {
	ret := &promodel.Attr{
		Id:         int32(v.Id),
		Name:       v.Name,
		IsFilter:   v.IsFilter,
		MultiChk:   v.MultiCheck,
		ItemValues: "",
		SortNum:    v.SortNum,
	}
	if v.Items != nil {
		ret.Items = make([]*promodel.AttrItem, len(v.Items))
		for i, v := range v.Items {
			ret.Items[i] = p.parseProductAttrItem(v)
		}
	}
	return ret
}

func (p *productService) parseProductSpec(v *proto.SProductSpec) *promodel.Spec {
	ret := &promodel.Spec{
		Id:         int32(v.Id),
		Name:       v.Name,
		ItemValues: v.ItemValues,
		SortNum:    v.SortNum,
	}
	if v.Items != nil {
		ret.Items = make([]*promodel.SpecItem, len(v.Items))
		for i, v := range v.Items {
			ret.Items[i] = p.parseProductSpecItem(v)
		}
	}
	return ret
}

func (p *productService) parseProductAttrItem(v *proto.SProductAttrItem) *promodel.AttrItem {
	return &promodel.AttrItem{
		Id:      int32(v.Id),
		Value:   v.Value,
		SortNum: v.SortNum,
	}
}

func (p *productService) parseProductSpecItem(v *proto.SProductSpecItem) *promodel.SpecItem {
	return &promodel.SpecItem{
		Id:      int32(v.Id),
		Value:   v.Value,
		Color:   v.Color,
		SortNum: v.SortNum,
	}
}

func (p *productService) parseProductAttrValueDto(v *product.AttrValue) *proto.SProductAttrValue {
	return &proto.SProductAttrValue{
		Id:       v.ID,
		AttrId:   v.AttrId,
		AttrName: v.AttrName,
		AttrData: v.AttrData,
		AttrWord: v.AttrWord,
	}
}

func (p *productService) parseSpecDto(v *promodel.Spec) *proto.SProductSpec {
	return &proto.SProductSpec{
		Id:         int64(v.Id),
		Name:       v.Name,
		SortNum:    v.SortNum,
		ItemValues: v.ItemValues,
		Items:      nil,
	}
}

func (p *productService) parseProductSpecItemDto(v *promodel.SpecItem) *proto.SProductSpecItem {
	return &proto.SProductSpecItem{
		Id:      int64(v.Id),
		Value:   v.Value,
		Color:   v.Color,
		SortNum: v.SortNum,
	}
}

func (p *productService) parseProductAttrValue(v *proto.SProductAttrValue) *product.AttrValue {
	return &product.AttrValue{
		ID:       v.Id,
		AttrName: v.AttrName,
		AttrId:   v.AttrId,
		AttrData: v.AttrData,
		AttrWord: v.AttrWord,
	}
}
