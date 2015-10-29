/**
 * Copyright 2015 @ z3q.net.
 * name : member
 * author : jarryliu
 * date : 2015-10-29 15:06
 * description :
 * history :
 */
package dto

type(
	SimpleMember struct{
		Id int `db:"id"`
		Name string `db:"name"`
		User string `db:"user"`
		Phone string `db:"phone"`
	}
)