/**
 * Copyright 2015 @ S1N1 Team.
 * name : sorter.go
 * author : jarryliu
 * date : 2015-08-17 16:19
 * description :
 * history :
 */
package front

import (
	"bytes"
	"fmt"
	"strings"
)

type SortItem struct {
	Name   string
	Text   string
	Option bool
}

var GoodsListSortItems []*SortItem = []*SortItem{
	{Name: "price", Text: "价格", Option: true},
	{Name: "sale", Text: "销量", Option: true},
	{Name: "rate", Text: "评价", Option: true},
}

func GetSorterHtml(items []*SortItem, selected string, urlPath string) string {
	var buf *bytes.Buffer = bytes.NewBufferString("<ul>")
	var selArr []string = strings.Split(selected, "_")
	var selName = selArr[0]

	if i := strings.Index(urlPath, "sort="); i != -1 {
		s := urlPath[i:]
		if j := strings.Index(s, "&"); j == -1 {
			urlPath = urlPath[:i]
		} else {
			urlPath = urlPath[0:i] + urlPath[i+j+5:]
		}
	}

	if !strings.HasSuffix(urlPath, "?") &&
		!strings.HasSuffix(urlPath, "&") {
		if strings.Index(urlPath, "?") != -1 {
			urlPath = urlPath + "&"
		} else {
			urlPath = urlPath + "?"
		}
	}

	var sortValue string
	var sortUrl string

	for i, v := range items {
		sortValue = ""
		if v.Name == selName {
			buf.WriteString("<li class=\"selected\">")
			if v.Option {
				if selArr[1] == "0" {
					sortValue = "1"
				} else {
					sortValue = "0"
				}
			}
		} else {
			if v.Option {
				sortValue = "0"
			}
			buf.WriteString("<li>")
		}

		if v.Option {
			sortUrl = v.Name + "_" + sortValue
		} else {
			sortUrl = v.Name
		}

		buf.WriteString(fmt.Sprintf("<a href=\"%ssort=%s\" sort-name=\"%s\" sort-val=\"%s\">%s</a>",
			urlPath, sortUrl, v.Name, sortValue, v.Text))

		if v.Option {
			buf.WriteString("<span class=\"d\"></span>")
		}
		if i != len(items)-1 {
			buf.WriteString("<span class=\"split\"></span>")
		}
		buf.WriteString("</li>")
	}
	buf.WriteString("</ul>")
	return buf.String()
}
