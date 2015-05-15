/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:47
 * description :
 * history :
 */

package sale

//分类
type ValueCategory struct {
	Id int `db:"id" auto:"yes" pk:"yes"`
	//父分类
	ParentId int `db:"parent_id"`
	//供应商编号
	PartnerId int `db:"partner_id"`
	//名称
	Name        string
	OrderIndex  int    `db:"order_index"`
	CreateTime  int64  `db:"create_time"`
	Enabled     int    `db:"enabled`
	Description string `db:"description"`
}
