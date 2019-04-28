/**
 * Copyright 2015 @ z3q.net.
 * name : category_manager.go
 * author : jarryliu
 * date : 2016-06-04 13:40
 * description :
 * history :
 */
package product

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ixre/gof/algorithm/iterator"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"sort"
	"strconv"
	"strings"
	"time"
)

var _ product.ICategory = new(categoryImpl)

// 分类实现
type categoryImpl struct {
	value           *product.Category
	rep             product.ICategoryRepo
	parentIdChanged bool
	childIdArr      []int32
	opt             domain.IOptionStore
}

func newCategory(rep product.ICategoryRepo,
	v *product.Category) product.ICategory {
	return &categoryImpl{
		value: v,
		rep:   rep,
	}
}

func (c *categoryImpl) GetDomainId() int32 {
	return c.value.ID
}

func (c *categoryImpl) GetValue() *product.Category {
	return c.value
}

func (c *categoryImpl) GetOption() domain.IOptionStore {
	if c.opt == nil {
		opt := newCategoryOption(c)
		if err := opt.Stat(); err != nil {
			opt.Set(product.C_OptionViewName, &domain.Option{
				Key:   product.C_OptionViewName,
				Type:  domain.OptionTypeString,
				Must:  false,
				Title: "显示页面",
				Value: "goods_list.html",
			})
			opt.Set(product.C_OptionDescribe, &domain.Option{
				Key:   product.C_OptionDescribe,
				Type:  domain.OptionTypeString,
				Must:  false,
				Title: "描述",
				Value: "",
			})
			opt.Flush()
		}
		c.opt = opt
	}
	return c.opt
}

// 检查上级分类是否正确
func (c *categoryImpl) checkParent(parentId int32) error {
	if id := c.GetDomainId(); id > 0 && parentId > 0 {
		//检查上级栏目是否存在
		p := c.rep.GlobCatService().GetCategory(parentId)
		if p == nil {
			return product.ErrNoSuchCategory
		}
		// 检查上级分类
		if p.GetValue().ParentId == id {
			return product.ErrCategoryCycleReference
		}
	}
	return nil
}

// 设置值
func (c *categoryImpl) SetValue(v *product.Category) error {
	val := c.value
	if val.ID == v.ID {
		val.Enabled = v.Enabled
		val.Priority = v.Priority
		val.Name = v.Name
		val.SortNum = v.SortNum
		val.Icon = v.Icon
		val.IconXY = v.IconXY
		// 设置产品模型
		val.ProModel = v.ProModel
		// 设置链接类型
		if c.GetDomainId() > 0 && val.VirtualCat != v.VirtualCat {
			return product.ErrIncorrectCategoryType
		}
		// 检测虚拟分类是否设置URL
		if val.VirtualCat == 1 {
			u := strings.TrimSpace(v.CatUrl)
			if u == "" {
				return product.ErrVirtualCatNoUrl
			} else {
				val.CatUrl = u
			}
		}
		val.VirtualCat = v.VirtualCat
		val.FloorShow = v.FloorShow
		if v.FloorShow == 1 && v.ParentId != 0 {
			return product.ErrCategoryFloorShow
		}
		if val.ParentId != v.ParentId {
			c.parentIdChanged = true
		} else {
			c.parentIdChanged = false
		}

		//todo: 默认都更改了，每次保存都重新设置URL. 这是一个bug需要修复
		if val.VirtualCat == 0 {
			c.parentIdChanged = true
		}

		if c.parentIdChanged {
			err := c.checkParent(v.ParentId)
			if err != nil {
				return err
			}
			val.ParentId = v.ParentId
		}
	}
	return nil
}

// 获取子栏目的编号
func (c *categoryImpl) GetChildes() []int32 {
	if c.childIdArr == nil {
		childCats := c.getChildCategories(c.GetDomainId())
		c.childIdArr = make([]int32, len(childCats))
		for i, v := range childCats {
			c.childIdArr[i] = v.ID
		}
	}
	return c.childIdArr
}
func (c *categoryImpl) setCategoryLevel() {
	var mchId int32 = 0
	list := c.rep.GetCategories(mchId)
	c.parentWalk(list, mchId, &c.value.Level)
	//todo: 未实现
}

func (c *categoryImpl) parentWalk(list []*product.Category,
	parentId int32, level *int32) {
	*level += 1
	if parentId <= 0 {
		return
	}
	for _, v := range list {
		if v.ID == v.ParentId {
			panic(errors.New(fmt.Sprintf(
				"Bad category , id is same of parent id , id:%s",
				v.ID)))
		} else if v.ID == parentId {
			c.parentWalk(list, v.ParentId, level)
			break
		}
	}
}

func (c *categoryImpl) Save() (int32, error) {
	//if c._manager.ReadOnly() {
	//    return c.GetDomainId(), product.ErrReadonlyCategory
	//}
	c.setCategoryLevel()
	id, err := c.rep.SaveCategory(c.value)
	if err == nil {
		c.value.ID = id
		// 非虚拟分类，自动设置链接地址
		if c.value.VirtualCat == 0 {
			//todo: ??? await refactor
			if c.parentIdChanged {
				c.value.CatUrl = c.getAutomaticUrl(id)
				c.parentIdChanged = false
				return c.Save()
			}
		}
	}
	return id, err
}

// 获取子栏目
func (c *categoryImpl) getChildCategories(catId int32) []*product.Category {
	var all = c.rep.GetCategories(0)
	var newArr []*product.Category

	var cdt iterator.Condition = func(v, v1 interface{}) bool {
		return v1.(*product.Category).ParentId == v.(*product.Category).ID
	}
	var start iterator.WalkFunc = func(v interface{}, level int) {
		c := v.(*product.Category)
		if c.ID != catId {
			newArr = append(newArr, c)
		}
	}

	var arr = make([]interface{}, len(all))
	for i := range arr {
		arr[i] = all[i]
	}

	iterator.Walk(arr, &product.Category{ID: catId}, cdt, start, nil, 1)

	return newArr
}

// 获取与栏目相关的栏目
func (c *categoryImpl) getRelationCategories(catId int32) []*product.Category {
	var all = c.rep.GetCategories(0)
	var newArr []*product.Category
	var isMatch bool
	var pid int32
	var l = len(all)

	for i := 0; i < l; i++ {
		if !isMatch && all[i].ID == catId {
			isMatch = true
			pid = all[i].ParentId
			newArr = append(newArr, all[i])
			i = -1
		} else {
			if all[i].ID == pid {
				newArr = append(newArr, all[i])
				pid = all[i].ParentId
				i = -1
				if pid == 0 {
					break
				}
			}
		}
	}
	return newArr
}

func (c *categoryImpl) getAutomaticUrl(id int32) string {
	relCats := c.getRelationCategories(id)
	var buf = bytes.NewBufferString("/list")
	var l = len(relCats)
	for i := l; i > 0; i-- {
		buf.WriteString("-")
		buf.WriteString(strconv.Itoa(int(relCats[i-1].ID)))
	}
	buf.WriteString(".html")
	return buf.String()
}

var _ domain.IOptionStore = new(categoryOption)

// 分类数据选项
type categoryOption struct {
	domain.IOptionStore
	cat *categoryImpl
}

func newCategoryOption(c *categoryImpl) domain.IOptionStore {
	file := fmt.Sprintf("conf/core/sale/cate_opt_%d", c.GetDomainId())
	return &categoryOption{
		cat:          c,
		IOptionStore: domain.NewOptionStoreWrapper(file),
	}
}

var _ product.IGlobCatService = new(categoryManagerImpl)

//当商户共享系统的分类时,没有修改的权限,既只读!
type categoryManagerImpl struct {
	readonly       bool
	repo           product.ICategoryRepo
	valRepo        valueobject.IValueRepo
	vendorId       int32
	lastUpdateTime int64
	categories     []product.ICategory
}

func NewCategoryManager(mchId int32, rep product.ICategoryRepo,
	valRepo valueobject.IValueRepo) product.IGlobCatService {
	c := &categoryManagerImpl{
		repo:     rep,
		vendorId: mchId,
		valRepo:  valRepo,
	}
	return c.init()
}

func (c *categoryManagerImpl) init() product.IGlobCatService {
	mchConf := c.valRepo.GetPlatformConf()
	if !mchConf.MchGoodsCategory && c.vendorId > 0 {
		c.readonly = true
		c.vendorId = 0
	}
	return c
}

// 获取栏目关联的编号,系统用0表示
func (c *categoryManagerImpl) getRelationId() int32 {
	return c.vendorId
}

// 清理缓存
func (c *categoryManagerImpl) clean() {
	c.categories = nil
}

// 是否只读,当商户共享系统的分类时,
// 没有修改的权限,即只读!
func (c *categoryManagerImpl) ReadOnly() bool {
	return c.readonly
}

// 创建分类
func (c *categoryManagerImpl) CreateCategory(v *product.Category) product.ICategory {
	if v.CreateTime == 0 {
		v.CreateTime = time.Now().Unix()
	}
	return newCategory(c.repo, v)
}

// 获取分类
func (c *categoryManagerImpl) GetCategory(id int32) product.ICategory {
	v := c.repo.GetCategory(c.getRelationId(), id)
	if v != nil {
		return c.CreateCategory(v)
	}
	return nil
}

// 获取所有分类
func (c *categoryManagerImpl) GetCategories() []product.ICategory {
	var list product.CategoryList = c.repo.GetCategories(c.getRelationId())
	sort.Sort(list)
	slice := make([]product.ICategory, len(list))
	for i, v := range list {
		slice[i] = c.CreateCategory(v)
	}
	return slice
}

// 删除分类
func (c *categoryManagerImpl) DeleteCategory(id int32) error {
	cat := c.GetCategory(id)
	if cat == nil {
		return product.ErrNoSuchCategory
	}
	if len(cat.GetChildes()) > 0 {
		return product.ErrHasChildCategories
	}
	if c.repo.CheckContainGoods(c.getRelationId(), id) {
		return product.ErrCategoryContainGoods
	}
	err := c.repo.DeleteCategory(c.getRelationId(), id)
	if err == nil {
		err = cat.GetOption().Destroy()
		cat = nil
	}
	return err
}

// 递归获取下级分类
func (c *categoryManagerImpl) CategoryTree(parentId int32) *product.Category {
	list := c.repo.GetCategories(0)
	var cat *product.Category
	if parentId == 0 {
		cat = &product.Category{ID: parentId}
	} else {
		for _, v := range list {
			if v.ID == parentId {
				cat = v
				break
			}
		}
		if cat == nil {
			return nil
		}
	}
	c.walkCategoryTree(cat, list)
	return cat
}

func (c *categoryManagerImpl) walkCategoryTree(node *product.Category,
	categories []*product.Category) {
	node.Children = []*product.Category{}
	for _, v := range categories {
		if v.ID != 0 && v.ParentId == node.ID {
			node.Children = append(node.Children, v)
			c.walkCategoryTree(v, categories)
		}
	}
}

// 获取分类关联的品牌
func (c *categoryManagerImpl) RelationBrands(catId int32) []*promodel.ProBrand {
	p := c.GetCategory(catId)
	if p != nil {
		var idArr []int32
		c.childWalk(p, &idArr)
		return c.repo.GetRelationBrands(idArr)
	}
	return []*promodel.ProBrand{}
}

func (c *categoryManagerImpl) childWalk(p product.ICategory, idArr *[]int32) {
	childes := p.GetChildes()
	if len(childes) > 0 {
		*idArr = append(*idArr, childes...)
		for _, v := range childes {
			c.childWalk(c.GetCategory(v), idArr)
		}
	}
}
