/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:47
 * description :
 * history :
 */

package sale

//分类
type ValueCategory struct {
	Id int `db:"id" auto:"yes" pk:"yes" json:"id"`
	//父分类
	ParentId int `db:"parent_id" json:"parentId"`
	//供应商编号
	PartnerId int `db:"partner_id" json:"PartnerId"`
	//名称
	Name        string           `db:"name" json:"name"`
	SortNumber  int              `db:"sort_number" json:"sortNumber"`
	Icon        string           `db:"icon" json:"icon"`
	Url         string           `db:"url" json:"url"`
	CreateTime  int64            `db:"create_time" json:"createTime"`
	Enabled     int              `db:"enabled" json:"enabled"`
	Description string           `db:"description" json:"description"`
	Child       []*ValueCategory `db:"-" json:"child"`
}
