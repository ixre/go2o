/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 15:18
 * description :
 * history :
 */
package delivery

//中国省市行政规划
type AreaValue struct {
	Id   int64  `db:"id" pk:"yes" auto:"no"`
	Pid  int64  `db:"pid"`
	Name string `db:"name"`
}
