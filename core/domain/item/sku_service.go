package item

import (
	"errors"
	"go2o/core/domain/interface/item"
	"sort"
	"strconv"
	"strings"
)

var _ item.ISkuService = new(skuServiceImpl)

type skuServiceImpl struct {
	repo item.IGoodsItemRepo
}

func NewSkuServiceImpl(repo item.IGoodsItemRepo) item.ISkuService {
	return &skuServiceImpl{
		repo: repo,
	}
}

// 将SKU字符串转为字典,如: 1:2;2:3
func (s *skuServiceImpl) SpecDataToMap(specData string) map[int]int {
	arr := strings.Split(specData, ";")
	l := len(arr)
	if l == 0 {
		panic("incorrent spec arr string")
	}
	mp := make(map[int]int, l)
	for _, s := range arr {
		i := strings.Index(s, ":")
		k, err := strconv.Atoi(s[:i])
		v, err1 := strconv.Atoi(s[i+1:])
		if err != nil || err1 != nil {
			panic(errors.New("spec arr key or value" +
				"not as int type! " + specData))
		}
		mp[k] = v
	}
	return mp
}

// 数组中是否存在元素
func (s *skuServiceImpl) arrExists(arr []int, v int) bool {
	for _, j := range arr {
		if j == v {
			return true
		}
	}
	return false
}

// 获取规格和项的数组
func (s *skuServiceImpl) GetSpecItemArray(sku []*item.Sku) ([]int, []int) {
	sa, ia := []int{}, []int{} //规格与规格项编号的数组
	for _, v := range sku {
		for k2, v2 := range s.SpecDataToMap(v.SpecData) {
			if !s.arrExists(sa, k2) {
				sa = append(sa, k2)
			}
			if !s.arrExists(ia, v2) {
				ia = append(ia, v2)
			}
		}
	}
	sort.Ints(sa)
	sort.Ints(ia)
	return sa, ia
}
