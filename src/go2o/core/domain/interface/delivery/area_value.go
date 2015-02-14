/**
 * Copyright 2014 @ ops Inc.
 * name : 
 * author : newmin
 * date : 2014-02-14 15:18
 * description :
 * history :
 */
package delivery

//中国省市行政规划
type AreaValue struct {
	Id   int	`db:"id" pk:"yes" auto:"no"`
	Pid  int	`db:"pid"`
	Name string `db:"name"`
}
