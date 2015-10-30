/**
 * Copyright 2015 @ z3q.net.
 * name : member
 * author : jarryliu
 * date : 2015-10-29 15:06
 * description :
 * history :
 */
package dto

type (
	SimpleMember struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		User  string `json:"user"`
		Phone string `json:"phone"`
	}
)
