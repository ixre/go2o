/**
 * Copyright 2014 @ Ops.
 * name :
 * author : newmin
 * date : 2013-11-05 17:37
 * description :
 * history :
 */

package apicache

import (
	"bytes"
	"com/service/goclient"
	"fmt"
	"ops/cf/app"
)

func GetCategories(c app.Context, partnerId int, secret string) []byte {
	var buf *bytes.Buffer = bytes.NewBufferString("")
	categories, err := goclient.Partner.Category(partnerId, secret)

	buf.WriteString(`<ul class="categories">
		<li class="s0 current" val="0">
			<div class="name"><span><strong>全部</strong></div>
		</li>
	`)
	if err == nil {
		for i, v := range categories {
			buf.WriteString(fmt.Sprintf(`<li class="s%d" val="%d">
			<div class="name"><span><strong>%s</strong></div>
			</li>`, i+1, v.Id, v.Name))
		}
	}
	buf.WriteString("</ul>")
	return buf.Bytes()
}
