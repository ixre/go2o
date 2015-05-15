/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-24 19:56
 * description :
 * history :
 */

package ui

import (
	"fmt"
	"github.com/atnet/gof/web/pager"
	"strconv"
)

var _ pager.PagerGetter = new(jsPagerGetter)

const (
	format = "javascript:gp(%d)"
)

type jsPagerGetter struct {
}

func (this *jsPagerGetter) Get(page, total, nowPage, flag int) (url, text string) {
	if flag&CONTROL != 0 {
		if flag&PREVIOUS != 0 {
			if page == 1 {
				return "javascript:;", FirstPageText
			}
			return fmt.Sprintf(format, nowPage), PreviousPageText
		}

		if flag&NEXT != 0 {
			if page == total {
				return "javascript:;", LastPageText
			}
			return fmt.Sprintf(format, nowPage), NextPageText
		}
	}

	return fmt.Sprintf(format, nowPage), strconv.Itoa(nowPage)
}
