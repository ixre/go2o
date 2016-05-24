/**
 * Copyright 2015 @ z3q.net.
 * name : category
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package util

import (
	"github.com/jsix/gof/algorithm/iterator"
	"go2o/core/domain/interface/sale"
)

type CategoryFormatFunc func(c *sale.ValueCategory, level int)

// 迭代栏目
func IterateCategory(categories []*sale.ValueCategory, c *sale.ValueCategory,
	iterateFunc CategoryFormatFunc, finishFunc CategoryFormatFunc, level int) {
	if c.Id != 0 {
		iterateFunc(c, level)
	}

	for _, k := range categories {
		if k.ParentId == c.Id {
			IterateCategory(categories, k, iterateFunc, finishFunc, level+1)
		}
	}

	if finishFunc != nil {
		finishFunc(c, level)
	}

}

// 迭代栏目
func WalkCategory(cs []*sale.ValueCategory, v *sale.ValueCategory, start iterator.WalkFunc, over iterator.WalkFunc) {
	var condition iterator.Condition = func(v, v1 interface{}) bool {
		return v1.(*sale.ValueCategory).ParentId == v.(*sale.ValueCategory).Id
	}
	var arr []interface{} = make([]interface{}, len(cs))
	for i, v := range cs {
		arr[i] = v
	}
	iterator.Walk(arr, v, condition, start, over, 0)
}
