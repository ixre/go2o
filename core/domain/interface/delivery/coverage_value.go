/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-12 17:08
 * description :
 * history :
 */
package delivery

// 覆盖区域
type CoverageValue struct {
	Id      int32   `db:"id" auto:"yes" pk:"true"`
	Name    string  `db:"name"`
	Lng     float64 `db:"lng"`
	Lat     float64 `db:"lat"`
	Radius  int     `db:"radius"`
	Address string  `db:"address"`
	AreaId  int32   `db:"area_id"`
}
