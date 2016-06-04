/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package cache

import (
	"bytes"
	"fmt"
	"github.com/jsix/gof/algorithm/iterator"
	"go2o/core/domain/interface/sale"
	"go2o/core/infrastructure/domain/util"
	"go2o/core/service/dps"
)

func readToCategoryDropList(merchantId int) []byte {
	categories := dps.SaleService.GetCategories(merchantId)
	buf := bytes.NewBuffer([]byte{})
	var f iterator.WalkFunc = func(v1 interface{}, level int) {
		c := v1.(*sale.Category)
		if c.Id != 0 {
			buf.WriteString(fmt.Sprintf(
				`<option class="opt%d" value="%d">%s</option>`,
				level,
				c.Id,
				c.Name,
			))
		}
	}
	util.WalkCategory(categories, &sale.Category{Id: 0}, f, nil)
	return buf.Bytes()
}

// 获取分类下拉选项
func GetDropOptionsOfCategory(merchantId int) []byte {
	return readToCategoryDropList(merchantId)
}
