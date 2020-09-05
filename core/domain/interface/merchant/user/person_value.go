/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 16:19
 * description :
 * history :
 */
package user

// 人员资料
type PersonValue struct {
	Id       int32  `db:"id" pk:"yes" auto:"yes"`
	Name     string `db:"name"`
	RealName string `db:"real_name"`
	Phone    string `db:"phone"`
	Sex      int    `db:"sex"`
	BirthDay int    `db:"birth_day"`
	Enabled  int    `db:"enabled`
}
