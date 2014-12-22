package cache

import (
	"bytes"
	"com/ording/dao"
	"com/ording/entity"
	"fmt"
)

type CategoryFormatFunc func(buf *bytes.Buffer, c *entity.Category, level int)

func readToCategoryDropList(partnerId int) []byte {
	categorys := dao.Category().GetCategoriesOfPartner(partnerId)
	buf := bytes.NewBuffer([]byte{})
	var f CategoryFormatFunc = func(buf *bytes.Buffer, c *entity.Category, level int) {
		buf.WriteString(fmt.Sprintf(
			`<option class="opt%d" value="%d">%s</option>`,
			level,
			c.Id,
			c.Name,
		))
	}
	itrCategory(buf, categorys, &entity.Category{Id: 0}, f, 0)

	return buf.Bytes()
}

func itrCategory(buf *bytes.Buffer, categorys []entity.Category, c *entity.Category, f CategoryFormatFunc, level int) {
	if c.Id != 0 {
		f(buf, c, level)
	}

	for _, k := range categorys {
		if k.Pid == c.Id {
			itrCategory(buf, categorys, &k, f, level+1)
		}
	}
}

// 获取分类下拉选项
func GetDropOptionsOfCategory(partnerId int) []byte {
	return readToCategoryDropList(partnerId)
}
