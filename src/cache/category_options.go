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
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/infrastructure/domain/util"
	"go2o/src/core/service/dps"
)

func readToCategoryDropList(partnerId int) []byte {
	categories := dps.SaleService.GetCategories(partnerId)
	buf := bytes.NewBuffer([]byte{})
	var f iterator.WalkFunc = func(v1 interface{}, level int) {
		c := v1.(*sale.ValueCategory)
		if c.Id != 0 {
			buf.WriteString(fmt.Sprintf(
				`<option class="opt%d" value="%d">%s</option>`,
				level,
				c.Id,
				c.Name,
			))
		}
	}
	util.WalkCategory(categories, &sale.ValueCategory{Id: 0}, f, nil)
<<<<<<< HEAD
=======

>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	return buf.Bytes()
}

// 获取分类下拉选项
func GetDropOptionsOfCategory(partnerId int) []byte {
	return readToCategoryDropList(partnerId)
}
