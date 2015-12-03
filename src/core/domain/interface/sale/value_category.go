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
	Id int `db:"id" auto:"yes" pk:"yes"`
	//父分类
	ParentId int `db:"parent_id"`
	//供应商编号
	PartnerId int `db:"partner_id"`
	//名称
	Name        string `db:"name"`
	OrderIndex  int    `db:"order_index"`
	Url         string `db:"url"`
	CreateTime  int64  `db:"create_time"`
	Enabled     int    `db:"enabled`
	Description string `db:"description"`
}
