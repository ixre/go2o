/**
 * Copyright 2015 @ z3q.net.
 * name : category_manager.go
 * author : jarryliu
 * date : 2016-06-04 13:40
 * description :
 * history :
 */
package sale

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jsix/gof/algorithm/iterator"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"sort"
	"strconv"
	"strings"
	"time"
)

var _ sale.ICategory = new(categoryImpl)

// 分类实现
type categoryImpl struct {
	value           *sale.Category
	rep             sale.ICategoryRep
	parentIdChanged bool
	childIdArr      []int
	opt             domain.IOptionStore
}

func NewCategory(rep sale.ICategoryRep, v *sale.Category) sale.ICategory {
	return &categoryImpl{
		value: v,
		rep:   rep,
	}
}

func (c *categoryImpl) GetDomainId() int {
	return c.value.Id
}

func (c *categoryImpl) GetValue() *sale.Category {
	return c.value
}

func (c *categoryImpl) GetOption() domain.IOptionStore {
	if c.opt == nil {
		opt := newCategoryOption(c)
		if err := opt.Stat(); err != nil {
			opt.Set(sale.C_OptionViewName, &domain.Option{
				Key:   sale.C_OptionViewName,
				Type:  domain.OptionTypeString,
				Must:  false,
				Title: "显示页面",
				Value: "goods_list.html",
			})
			opt.Set(sale.C_OptionDescribe, &domain.Option{
				Key:   sale.C_OptionDescribe,
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
func (c *categoryImpl) checkParent(parentId int) error {
	if id := c.GetDomainId(); id > 0 && parentId > 0 {
		//检查上级栏目是否存在
		p := c.rep.GetGlobManager().GetCategory(parentId)
		if p == nil {
			return sale.ErrNoSuchCategory
		}
		// 检查上级分类
		if p.GetValue().ParentId == id {
			return sale.ErrCategoryCycleReference
		}
	}
	return nil
}

// 设置值
func (c *categoryImpl) SetValue(v *sale.Category) error {
	val := c.value
	if val.Id == v.Id {
		val.Description = v.Description
		val.Enabled = v.Enabled
		val.Name = v.Name
		val.SortNumber = v.SortNumber
		val.Icon = v.Icon
		if val.ParentId != v.ParentId {
			c.parentIdChanged = true
		} else {
			c.parentIdChanged = false
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
func (c *categoryImpl) GetChildes() []int {
	if c.childIdArr == nil {
		childCats := c.GetChildCategories(
			c.value.MerchantId, c.GetDomainId())
		c.childIdArr = make([]int, len(childCats))
		for i, v := range childCats {
			c.childIdArr[i] = v.Id
		}
	}
	return c.childIdArr
}
func (c *categoryImpl) setCategoryLevel() {
	mchId := c.value.MerchantId
	list := c.rep.GetCategories(mchId)
	c.parentWalk(list, mchId, &c.value.Level)
}

func (c *categoryImpl) parentWalk(list []*sale.Category,
	parentId int, level *int) {
	*level += 1
	if parentId <= 0 {
		return
	}
	for _, v := range list {
		if v.Id == v.ParentId {
			panic(errors.New(fmt.Sprintf(
				"Bad category , id is same of parent id , id:%s",
				v.Id)))
		} else if v.Id == parentId {
			c.parentWalk(list, v.ParentId, level)
			break
		}
	}
}

func (c *categoryImpl) Save() (int, error) {
	//if c._manager.ReadOnly() {
	//    return c.GetDomainId(), sale.ErrReadonlyCategory
	//}
	c.setCategoryLevel()
	id, err := c.rep.SaveCategory(c.value)
	if err == nil {
		c.value.Id = id
		if len(c.value.Url) == 0 || (c.parentIdChanged &&
			strings.HasPrefix(c.value.Url, "/c-")) {
			c.value.Url = c.getAutomaticUrl(c.value.MerchantId, id)
			c.parentIdChanged = false
			return c.Save()
		}
	}
	return id, err
}

// 获取子栏目
func (c *categoryImpl) GetChildCategories(mchId, categoryId int) []*sale.Category {
	var all []*sale.Category = c.rep.GetCategories(mchId)
	var newArr []*sale.Category = []*sale.Category{}

	var cdt iterator.Condition = func(v, v1 interface{}) bool {
		return v1.(*sale.Category).ParentId == v.(*sale.Category).Id
	}
	var start iterator.WalkFunc = func(v interface{}, level int) {
		c := v.(*sale.Category)
		if c.Id != categoryId {
			newArr = append(newArr, c)
		}
	}

	var arr []interface{} = make([]interface{}, len(all))
	for i := range arr {
		arr[i] = all[i]
	}

	iterator.Walk(arr, &sale.Category{Id: categoryId}, cdt, start, nil, 1)

	return newArr
}

// 获取与栏目相关的栏目
func (c *categoryImpl) GetRelationCategories(mchId, categoryId int) []*sale.Category {
	var all []*sale.Category = c.rep.GetCategories(mchId)
	var newArr []*sale.Category = []*sale.Category{}
	var isMatch bool
	var pid int
	var l int = len(all)

	for i := 0; i < l; i++ {
		if !isMatch && all[i].Id == categoryId {
			isMatch = true
			pid = all[i].ParentId
			newArr = append(newArr, all[i])
			i = -1
		} else {
			if all[i].Id == pid {
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

func (c *categoryImpl) getAutomaticUrl(merchantId, id int) string {
	relCats := c.GetRelationCategories(merchantId, id)
	var buf *bytes.Buffer = bytes.NewBufferString("/c")
	var l int = len(relCats)
	for i := l; i > 0; i-- {
		buf.WriteString("-" + strconv.Itoa(relCats[i-1].Id))
	}
	buf.WriteString(".htm")
	return buf.String()
}

var _ domain.IOptionStore = new(categoryOption)

// 分类数据选项
type categoryOption struct {
	domain.IOptionStore
	_mchId int
	_c     *categoryImpl
}

func newCategoryOption(c *categoryImpl) domain.IOptionStore {
	var file string
	mchId := c.GetValue().MerchantId
	if mchId > 0 {
		file = fmt.Sprintf("conf/mch/%d/option/c/%d", mchId, c.GetDomainId())
	} else {
		file = fmt.Sprintf("conf/core/sale/cate_opt_%d", c.GetDomainId())
	}
	return &categoryOption{
		_mchId:       c.GetValue().ParentId,
		_c:           c,
		IOptionStore: domain.NewOptionStoreWrapper(file),
	}
}

var _ sale.ICategoryManager = new(categoryManagerImpl)

//当商户共享系统的分类时,没有修改的权限,既只读!
type categoryManagerImpl struct {
	_readonly      bool
	_rep           sale.ICategoryRep
	_valRep        valueobject.IValueRep
	_mchId         int
	lastUpdateTime int64
	_categories    []sale.ICategory
}

func NewCategoryManager(mchId int, rep sale.ICategoryRep,
	valRep valueobject.IValueRep) sale.ICategoryManager {
	c := &categoryManagerImpl{
		_rep:    rep,
		_mchId:  mchId,
		_valRep: valRep,
	}
	return c.init()
}

func (c *categoryManagerImpl) init() sale.ICategoryManager {
	mchConf := c._valRep.GetPlatformConf()
	if !mchConf.MchGoodsCategory && c._mchId > 0 {
		c._readonly = true
		c._mchId = 0
	}
	return c
}

// 获取栏目关联的编号,系统用0表示
func (c *categoryManagerImpl) getRelationId() int {
	return c._mchId
}

// 清理缓存
func (c *categoryManagerImpl) clean() {
	c._categories = nil
}

// 是否只读,当商户共享系统的分类时,
// 没有修改的权限,即只读!
func (c *categoryManagerImpl) ReadOnly() bool {
	return c._readonly
}

// 创建分类
func (c *categoryManagerImpl) CreateCategory(v *sale.Category) sale.ICategory {
	if v.CreateTime == 0 {
		v.CreateTime = time.Now().Unix()
	}
	v.MerchantId = c.getRelationId()
	return NewCategory(c._rep, v)
}

// 获取分类
func (c *categoryManagerImpl) GetCategory(id int) sale.ICategory {
	v := c._rep.GetCategory(c.getRelationId(), id)
	if v != nil {
		return c.CreateCategory(v)
	}
	return nil
}

// 获取所有分类
func (c *categoryManagerImpl) GetCategories() []sale.ICategory {
	var list sale.CategoryList = c._rep.GetCategories(c.getRelationId())
	sort.Sort(list)
	slice := make([]sale.ICategory, len(list))
	for i, v := range list {
		slice[i] = c.CreateCategory(v)
	}
	return slice
}

// 删除分类
func (c *categoryManagerImpl) DeleteCategory(id int) error {
	cat := c.GetCategory(id)
	if cat == nil {
		return sale.ErrNoSuchCategory
	}
	if len(cat.GetChildes()) > 0 {
		return sale.ErrHasChildCategories
	}
	if c._rep.CheckGoodsContain(c.getRelationId(), id) {
		return sale.ErrCategoryContainGoods
	}
	err := c._rep.DeleteCategory(c.getRelationId(), id)
	if err == nil {
		err = cat.GetOption().Destroy()
		cat = nil
	}
	return err
}
