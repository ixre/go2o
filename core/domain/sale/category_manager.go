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
	_value           *sale.Category
	_rep             sale.ICategoryRep
	_parentIdChanged bool
	_childIdArr      []int
	_opt             domain.IOptionStore
}

func NewCategory(rep sale.ICategoryRep, v *sale.Category) sale.ICategory {
	return &categoryImpl{
		_value: v,
		_rep:   rep,
	}
}

func (c *categoryImpl) GetDomainId() int {
	return c._value.Id
}

func (c *categoryImpl) GetValue() *sale.Category {
	return c._value
}

func (c *categoryImpl) GetOption() domain.IOptionStore {
	if c._opt == nil {
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
		c._opt = opt
	}
	return c._opt
}

func (c *categoryImpl) SetValue(v *sale.Category) error {
	val := c._value
	if val.Id == v.Id {
		val.Description = v.Description
		val.Enabled = v.Enabled
		val.Name = v.Name
		val.SortNumber = v.SortNumber
		val.Icon = v.Icon
		if val.ParentId != v.ParentId {
			c._parentIdChanged = true
			val.ParentId = v.ParentId
		} else {
			c._parentIdChanged = false
		}
	}
	return nil
}

// 获取子栏目的编号
func (c *categoryImpl) GetChildes() []int {
	if c._childIdArr == nil {
		childCats := c.GetChildCategories(
			c._value.MerchantId, c.GetDomainId())
		c._childIdArr = make([]int, len(childCats))
		for i, v := range childCats {
			c._childIdArr[i] = v.Id
		}
	}
	return c._childIdArr
}
func (c *categoryImpl) setCategoryLevel() {
	mchId := c._value.MerchantId
	list := c._rep.GetCategories(mchId)
	c.parentWalk(list, mchId, &c._value.Level)
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
	id, err := c._rep.SaveCategory(c._value)
	if err == nil {
		c._value.Id = id
		if len(c._value.Url) == 0 || (c._parentIdChanged &&
			strings.HasPrefix(c._value.Url, "/c-")) {
			c._value.Url = c.getAutomaticUrl(c._value.MerchantId, id)
			c._parentIdChanged = false
			return c.Save()
		}
	}
	return id, err
}

// 获取子栏目
func (c *categoryImpl) GetChildCategories(mchId, categoryId int) []*sale.Category {
	var all []*sale.Category = c._rep.GetCategories(mchId)
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
	var all []*sale.Category = c._rep.GetCategories(mchId)
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
		return sale.ErrCategoryNotExist
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
