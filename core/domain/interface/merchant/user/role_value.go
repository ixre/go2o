/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 16:59
 * description :
 * history :
 */
package user

type RoleValue struct {
	Id   int32  `db:"id" pk:"yes" auto:"no"`
	Name string `db:"name"`
	// 表示角色位值
	Flag    int `db:"flag"`
	Enabled int `db:"enabled"`
}
