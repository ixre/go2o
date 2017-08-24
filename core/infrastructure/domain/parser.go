package domain

import (
	"encoding/json"
	"strconv"
)

// 转换勾选字典,数据如：{"1":["10","11"],"2":["20","21"]}
func ParseCartCheckedMap(data string) (m map[int64][]int64) {
	if data != "" || data != "{}" {
		src := map[string][]string{}
		err := json.Unmarshal([]byte(data), &src)
		if err == nil {
			m = map[int64][]int64{}
			for k, v := range src {
				itemId, _ := strconv.Atoi(k)
				skuList := []int64{}
				for _, v2 := range v {
					skuId, _ := strconv.Atoi(v2)
					skuList = append(skuList, int64(skuId))
				}
				m[int64(itemId)] = skuList
			}
			return m
		}
	}
	return nil
}
