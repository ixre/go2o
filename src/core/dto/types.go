/**
 * Copyright 2015 @ z3q.net.
 * name : types
 * author : jarryliu
 * date : 2015-10-29 15:33
 * description :
 * history :
 */
package dto

type (
	TextObject struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
		Title string `json:"title"`
	}
)
