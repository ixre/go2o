/**
 * Copyright 2015 @ z3q.net.
 * name : value_gallery
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

import (
	"sort"
)

var _ sort.Interface = ValueGallery{}

// 轮播广告图片集合
type ValueGallery []*ValueImage

func (this ValueGallery) Len() int {
	return len(this)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (this ValueGallery) Less(i, j int) bool {
	return this[i].SortNumber < this[j].SortNumber || (this[i].SortNumber == this[j].SortNumber &&
		this[i].Id < this[j].Id)
}

// Swap swaps the elements with indexes i and j.
func (this ValueGallery) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
