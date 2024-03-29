/**
 * Copyright 2015 @ 56x.net.
 * name : category
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package format

import (
	"github.com/ixre/gof/util"
	"strconv"
	"strings"
)

// 获取栏目编号字符串
func IntArrStrJoin(ids []int) string {
	return util.JoinIntArray(ids, ",")
}

func I32ArrStrJoin(ids []int32) string {
	var strIds = make([]string, len(ids))
	for i, v := range ids {
		strIds[i] = strconv.Itoa(int(v))
	}
	return strings.Join(strIds, ",")
}

func I64ArrStrJoin(ids []int64) string {
	var strIds = make([]string, len(ids))
	for i, v := range ids {
		strIds[i] = strconv.Itoa(int(v))
	}
	return strings.Join(strIds, ",")
}
