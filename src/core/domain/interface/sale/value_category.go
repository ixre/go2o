/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:47
 * description :
 * history :
 */

package sale

import "sort"

//分类
type ValueCategory struct {
	Id int `db:"id" auto:"yes" pk:"yes"`
	//父分类
	ParentId int `db:"parent_id"`
	//供应商编号
	PartnerId int `db:"partner_id"`
	//名称
	Name        string           `db:"name"`
	SortNumber  int              `db:"sort_number"`
	Icon        string           `db:"icon"`
	Url         string           `db:"url"`
	CreateTime  int64            `db:"create_time"`
	Enabled     int              `db:"enabled"`
	Description string           `db:"description"`
	Child       []*ValueCategory `db:"-"`
}

var _ sort.Interface = new(CategoryList)

type CategoryList []*ValueCategory

func (c CategoryList) Len() int {
	return len(c)
}

func (c CategoryList) Less(i, j int) bool {
	return c[i].SortNumber < c[j].SortNumber || (c[i].SortNumber == c[j].SortNumber && c[i].Id < c[j].Id)
}

func (c CategoryList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
