package item

import (
	"errors"
	"github.com/jsix/gof/log"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
	"sort"
	"strconv"
	"strings"
)

var _ item.ISkuService = new(skuServiceImpl)

type skuServiceImpl struct {
	repo     item.IGoodsItemRepo
	proMRepo promodel.IProModelRepo
}

func NewSkuServiceImpl(repo item.IGoodsItemRepo,
	proMRepo promodel.IProModelRepo) item.ISkuService {
	return &skuServiceImpl{
		repo:     repo,
		proMRepo: proMRepo,
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

// 合并SKU数组；主要是SKU编号的复制
func (s *skuServiceImpl) Merge(from []*item.Sku, to *[]*item.Sku) {
	if to == nil || from == nil || len(from) == 0 || len(*to) == 0 {
		return
	}
	dst := *to
	fromMap := make(map[string]*item.Sku, len(from))
	for _, v := range from {
		fromMap[v.SpecData] = v
	}
	for _, v := range dst {
		if fs, ok := fromMap[v.SpecData]; ok {
			//log.Println("SKU MERGE > dst: ",v.Id,"; src:",fs.Id)
			if v.Id == 0 {
				v.Id = fs.Id
			}
		}
	}
}

// 根据SKU获取规格名称字典，多个SKU的规格名称是相同的
func (s *skuServiceImpl) GetSpecNameMap(sku *item.Sku) (
	specMap map[int]string, err error) {
	specMap = make(map[int]string)
	sku.SpecData = strings.TrimSpace(sku.SpecData)
	sku.SpecWord = strings.TrimSpace(sku.SpecWord)
	if sku.SpecData == "" || sku.SpecWord == "" {
		return nil, errors.New("sku pair is null")
	}
	arr := strings.Split(sku.SpecData, ";")
	wordArr := strings.Split(sku.SpecWord, " ")
	if len(arr) != len(wordArr) {
		err = errors.New("sku pair not match")
	} else {
		for j := 0; j < len(arr); j++ {
			ki := strings.Index(arr[j], ":")
			wi := strings.Index(wordArr[j], ":")
			if ki == -1 || wi == -1 {
				continue
			}
			specId, err := strconv.Atoi(arr[j][:ki])
			if err != nil {
				return specMap, err
			}
			specMap[specId] = wordArr[j][:wi]
		}
	}
	return specMap, err
}

// 根据SKU获取规格项名称字典,多个SKU的规格项名称可能不同
func (s *skuServiceImpl) GetItemNameMap(sku *item.Sku) (
	itemMap map[int]string, err error) {
	itemMap = make(map[int]string)
	sku.SpecData = strings.TrimSpace(sku.SpecData)
	sku.SpecWord = strings.TrimSpace(sku.SpecWord)
	if sku.SpecData == "" || sku.SpecWord == "" {
		return nil, errors.New("sku pair is null")
	}
	arr := strings.Split(sku.SpecData, ";")
	wordArr := strings.Split(sku.SpecWord, " ")
	if len(arr) != len(wordArr) {
		err = errors.New("sku pair not match")
	} else {
		for j := 0; j < len(arr); j++ {
			ki := strings.Index(arr[j], ":")
			wi := strings.Index(wordArr[j], ":")
			if ki == -1 || wi == -1 {
				continue
			}
			itemId, err := strconv.Atoi(arr[j][ki+1:])
			if err != nil {
				return nil, err
			}
			itemMap[itemId] = wordArr[j][wi+1:]
		}
	}
	return itemMap, err
}

// 重建SKU数组，将信息附加
func (s *skuServiceImpl) RebuildSkuArray(sku *[]*item.Sku, it *item.GoodsItem) (err error) {
	skuArr := *sku
	// 获取传入的规格信息,按传入规格名称SKU直到找到结果为止
	var inSpecNameMap map[int]string
	for _, v := range skuArr {
		inSpecNameMap, err = s.GetSpecNameMap(v)
		if err == nil {
			break
		}
	}
	// 获取当前规格及规格项存储于map中
	skMap := map[int]*promodel.Spec{}
	siMap := map[int]*promodel.SpecItem{}
	sa, ia := s.GetSpecItemArray(skuArr)
	for _, v := range sa {
		spec := s.proMRepo.GetSpec(v)
		skMap[v] = spec
	}
	for _, v := range ia {
		s := s.proMRepo.GetSpecItem(v)
		siMap[v] = s
	}
	// 赋值SpecWord
	for _, v := range skuArr {
		// 图片
		if strings.TrimSpace(v.Image) == "" {
			v.Image = it.Image
		}
		arr := strings.Split(v.SpecData, ";")
		l := len(arr)
		if l == 0 {
			continue
		}
		itemNameMap, _ := s.GetItemNameMap(v)

		//log.Println("--",itemNameMap,er)

		// 设置规格字符
		items := make([]string, l)
		for i, v := range arr {
			idx := strings.Index(v, ":")
			specId, _ := strconv.Atoi(v[:idx])
			if inSpecNameMap != nil {
				if n, ok := inSpecNameMap[specId]; ok && n != "" {
					items[i] = n
				} else {
					if spec := skMap[specId]; spec != nil {
						items[i] = spec.Name
					}
				}
			} else {
				if spec := skMap[specId]; spec != nil {
					items[i] = spec.Name
				}
			}
			items[i] += ":"
			itemId, _ := strconv.Atoi(v[idx+1:])
			if itemNameMap != nil {
				if n, ok := itemNameMap[itemId]; ok && n != "" {
					items[i] += n
				} else {
					if im := siMap[itemId]; im != nil {
						items[i] += im.Value
					}
				}
			} else {
				if im := siMap[itemId]; im != nil {
					items[i] += im.Value
				}
			}

		}
		v.SpecWord = strings.TrimSpace(strings.Join(items, " "))

		// 标题为空，则自动设置
		if strings.TrimSpace(v.Title) == "" {
			titArr := make([]string, l+1)
			titArr[0] = it.Title
			for i, v := range arr {
				ii := strings.Index(v, ":")
				iid, _ := strconv.Atoi(v[ii+1:])
				if im := siMap[iid]; im != nil {
					titArr[i+1] = im.Value
				}
			}
			v.Title = strings.TrimSpace(strings.Join(titArr, " "))
		}
	}
	return nil
}
