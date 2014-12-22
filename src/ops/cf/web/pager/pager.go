/**
 * Copyright 2014 @ Ops.
 * name :
 * author : newmin
 * date : 2013-11-17 07:49
 * description :
 * history :
 */

package pager

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	FirstLinkFormat      = ""
	PagerLinkFormat      = "?page=%d"
	PagerLinkCount       = 10
	FirstPageText        = "首页"
	LastPageText         = "末页"
	NextPageText         = "下一页"
	PreviousPageText     = "上一页"
	CollagePagerLinkText = "..."
)
const (
	CONTROL = 1 << iota
	PREVIOUS
	NEXT
)

var (
	DefaultPagerGetter PagerGetter = new(defaultPagerGetter)
)

type PagerGetter interface {
	Get(page, total, nowPage, flag int) (url, text string)
}

type defaultPagerGetter struct {
}

func (this *defaultPagerGetter) Get(page, total, nowPage, flag int) (url, text string) {
	if flag&CONTROL != 0 {
		if flag&PREVIOUS != 0 {
			if page == 1 {
				return "javascript:;", FirstPageText
			}
			return fmt.Sprintf(PagerLinkFormat, nowPage), PreviousPageText
		}

		if flag&NEXT != 0 {
			if page == total {
				return "javascript:;", LastPageText
			}
			return fmt.Sprintf(PagerLinkFormat, nowPage), NextPageText
		}
	}

	if nowPage == 1 && len(FirstLinkFormat) != 0 {
		return FirstLinkFormat, strconv.Itoa(nowPage)
	}
	return fmt.Sprintf(PagerLinkFormat, nowPage), strconv.Itoa(nowPage)
}

//分页方法
type PagerFuncGetter struct {
	GetFunc func(page, total, nowPage, flag int) (url, text string)
}

func (this *PagerFuncGetter) Get(page, total, nowPage, flag int) (url, text string) {
	return this.GetFunc(page, total, nowPage, flag)
}

type UrlPager struct {
	//当前页面索引（从1开始）
	cpi int

	//连接和页码
	getter PagerGetter

	//页面总数
	pageCount int

	//链接长度,创建多少个跳页链接
	LinkCount int

	//记录条数
	RecordCount int

	//选页框文本
	//SelectPageText string

	//第一页链接格式
	firstPageLink string

	//分页链接格式
	linkFormat string

	//页码文本格式
	pageTextFormat string

	//是否允许输入页码调页
	enableInput bool

	//使用选页
	enableSelect bool

	//分页详细记录,如果为空字符则用默认,为空则不显示
	PagerTotal string
}

func (this *UrlPager) Pager() []byte {
	var bys *bytes.Buffer
	var cls string
	var u, t string

	bys = bytes.NewBufferString("<div class=\"pager\">")

	//输出上一页
	if this.cpi > 1 {
		cls = "previous"
		u, t = this.getter.Get(this.cpi, this.pageCount, this.cpi-1, CONTROL|PREVIOUS)
	} else {
		cls = "disabled"
		u, t = this.getter.Get(this.cpi, this.pageCount, this.cpi, CONTROL|PREVIOUS)
	}
	bys.WriteString(fmt.Sprintf(`<span class="%s"><a href="%s">%s</a></span>`, cls, u, t))

	//起始页:CurrentPageIndex / 10 * 10+1
	//结束页:(CurrentPageIndex%10==0?CurrentPageIndex-1: CurrentPageIndex) / 10 * 10
	//当前页数能整除10的时候需要减去10页，否则不能选中

	//链接页码数量(默认10)
	var c int = this.LinkCount
	var startPage int = (this.cpi-1)/c*c + 1
	var _gotoPrevious bool = false //是否上一栏分页

	for i, j := 1, startPage; i <= c && j < this.pageCount; i++ {
		if this.cpi%c == 0 {
			j = (this.cpi-1)/c*c + i
		} else {
			j = this.cpi/c*c + i
		}

		if j == this.cpi {
			_gotoPrevious = j != 1 && j%c == 1

			//上一栏分页
			if _gotoPrevious {
				u, _ := this.getter.Get(this.cpi, this.pageCount, j-1, 0)
				bys.WriteString(fmt.Sprintf(`<a class="page" href="%s">%s</a>`, u, CollagePagerLinkText))
			}

			//如果为页码为当前页
			bys.WriteString(fmt.Sprintf(`<span class=\"current\">%d</span>`, j))

			//下一栏分页
			if !_gotoPrevious && j%c == 0 && j != this.pageCount {
				u, _ := this.getter.Get(this.cpi, this.pageCount, j+1, 0)
				bys.WriteString(fmt.Sprintf(`<a class="page" href="%s">%s</a>`, u, CollagePagerLinkText))
			}

		} else {
			//页码不为当前页，则输出页码
			u, t := this.getter.Get(this.cpi, this.pageCount, j, 0)
			bys.WriteString(fmt.Sprintf(`<a class="page" href="%s">%s</a>`, u, t))
		}
	}

	//输出下一页链接
	if this.cpi < this.pageCount {
		cls = "next"
		u, t = this.getter.Get(this.cpi, this.pageCount, this.cpi+1, CONTROL|NEXT)
	} else {
		cls = "disabled"
		u, t = this.getter.Get(this.cpi, this.pageCount, this.cpi, CONTROL|NEXT)
	}
	bys.WriteString(fmt.Sprintf(`<span class="%s"><a href="%s">%s</a></span>`, cls, u, t))

	//显示信息
	const pagerStr string = "<span class=\"pageinfo\">&nbsp;第%d/%d页，共%d条。</span>"
	if len(this.PagerTotal) == 0 {
		this.PagerTotal = pagerStr
	}
	bys.WriteString(fmt.Sprintf(this.PagerTotal, this.cpi, this.pageCount, this.RecordCount))

	bys.WriteString("</div>")
	return bys.Bytes()
}

func (this *UrlPager) PagerString() string {
	return string(this.Pager())
}

func NewUrlPager(totalPage, page int, pg PagerGetter) *UrlPager {
	if page < 1 {
		page = 1
	}
	if totalPage < 1 {
		totalPage = 1
	}

	p := &UrlPager{}
	p.pageCount = totalPage
	p.cpi = page
	p.LinkCount = PagerLinkCount

	if pg == nil {
		p.getter = DefaultPagerGetter
	} else {
		p.getter = pg
	}
	return p
}
