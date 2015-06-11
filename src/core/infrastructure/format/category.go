/**
 * Copyright 2015 @ S1N1 Team.
 * name : category
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package format
import (
    "strconv"
    "strings"
)

// 获取栏目编号字符串
func GetCategoryIdStr(ids []int)string{
    var strIds []string = make([]string, len(ids))
    for i, v := range ids {
        strIds[i] = strconv.Itoa(v)
    }
    return strings.Join(strIds, ",")
}