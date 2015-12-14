/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-01-22 21:08
 * description :
 * history :
 */
package front

type Pager struct {
	Total int                      `json:"total"`
	Rows  []map[string]interface{} `json:"rows"`
	Text  string                   `json:"text"`
}
