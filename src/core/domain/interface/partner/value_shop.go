/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-12-22 17:59
 * description :
 * history :
 */

package partner

//门店
type ValueShop struct {
	Id         int    `db:"id" pk:"yes" auto:"yes"`
	PartnerId  int    `db:"partner_id"`
	Name       string `db:"name"`
	Address    string `db:"address"`
	Phone      string `db:"phone"`
	SortNumber int    `db:"sort_number"`
	State      int    `db:"state"`
	CreateTime int64  `db:"create_time"`

	//    // 位置(经度+"/"+纬度)
	//    Location string `db:"location"`
	//
	//    // 配送最大半径(公里)
	//    DeliverRadius int `db:"deliver_radius"`
}
