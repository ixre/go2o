package item

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/infrastructure/format"
	"log"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

var _ item.ISkuService = new(skuServiceImpl)

type skuServiceImpl struct {
	repo     item.IGoodsItemRepo
	proMRepo promodel.IProModelRepo
	su       *skuServiceUtil
}

func NewSkuServiceImpl(repo item.IGoodsItemRepo,
	proMRepo promodel.IProModelRepo) item.ISkuService {
	s := &skuServiceImpl{
		repo:     repo,
		proMRepo: proMRepo,
	}
	s.su = newSkuUtil(s)
	return s
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
func (s *skuServiceImpl) GetSpecItemArray(sku []*item.Sku) (
	sa []int, ia []int) {
	sa, ia = []int{}, []int{} //规格与规格项编号的数组
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

// 根据SKU更新商品的信息
func (s *skuServiceImpl) UpgradeBySku(it *item.GoodsItem,
	arr []*item.Sku) error {
	//更新SKU数量
	it.SkuNum = int32(len(arr))
	//如果包含SKU，则更新库存和价格区间
	if it.SkuNum > 0 {
		//更新库存
		it.StockNum = 0
		//更新销售数量
		it.SaleNum = 0
		var pl, ph float32
		for i := 0; i < len(arr); i++ {
			it.StockNum += arr[i].Stock
			it.SaleNum += arr[i].SaleNum
			price := arr[i].Price
			if price < pl || pl == 0 {
				pl = price
			}
			if price > ph || ph == 0 {
				ph = price
			}
		}
		//更新价格区间
		it.Price = pl
		if pl == ph {
			it.PriceRange = format.FormatFloat(pl)
		} else {
			it.PriceRange = format.FormatFloat(pl) +
				"~" + format.FormatFloat(ph)
		}
	}
	return nil
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
			//log.Println("SKU MERGE > dst: ",v.ID,"; src:",fs.ID)
			if v.ID == 0 {
				v.ID = fs.ID
			}
		}
	}
}

// 根据SKU获取规格名称字典，多个SKU的规格名称是相同的
func (s *skuServiceImpl) GetNameMap(skuArr []*item.Sku) (
	specMap map[int32]string, itemMap map[int32]string) {
	// 获取传入的规格信息,按传入规格名称SKU
	specMap = make(map[int32]string)
	itemMap = make(map[int32]string)
	for _, sku := range skuArr {
		sku.SpecData = strings.TrimSpace(sku.SpecData)
		sku.SpecWord = strings.TrimSpace(sku.SpecWord)
		if sku.SpecData == "" || sku.SpecWord == "" {
			continue
		}
		arr := strings.Split(sku.SpecData, ";")
		wordArr := strings.Split(sku.SpecWord, " ")
		if len(arr) != len(wordArr) {
			continue
		}
		for j := 0; j < len(arr); j++ {
			ki := strings.Index(arr[j], ":")
			wi := strings.Index(wordArr[j], ":")
			if ki == -1 || wi == -1 {
				continue
			}
			//绑定规格
			specId, _ := util.I32Err(strconv.Atoi(arr[j][:ki]))
			if _, ok := specMap[specId]; !ok {
				specMap[specId] = wordArr[j][:wi]
			}
			//绑定项
			itemId, _ := util.I32Err(strconv.Atoi(arr[j][ki+1:]))
			if _, ok := itemMap[itemId]; !ok {
				itemMap[itemId] = wordArr[j][wi+1:]
			}
		}
	}
	return specMap, itemMap
}

// 重建SKU数组，将信息附加
func (s *skuServiceImpl) RebuildSkuArray(sku *[]*item.Sku,
	it *item.GoodsItem) (err error) {
	skuArr := *sku
	// 获取传入的规格信息,按传入规格名称SKU直到找到结果为止
	sName, iName := s.GetNameMap(skuArr)
	// 获取当前规格及规格项存储于map中
	skMap := map[int32]*promodel.Spec{}
	siMap := map[int32]*promodel.SpecItem{}
	sa, ia := s.GetSpecItemArray(skuArr)
	for _, v := range sa {
		skMap[int32(v)] = s.proMRepo.GetSpec(v)
	}
	for _, v := range ia {
		siMap[int32(v)] = s.proMRepo.GetSpecItem(v)
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
		// 设置规格字符
		items := make([]string, l)
		for i, v := range arr {
			idx := strings.Index(v, ":")
			specId, _ := util.I32Err(strconv.Atoi(v[:idx]))
			if n, ok := sName[specId]; ok && n != "" {
				items[i] = n
			} else {
				items[i] = skMap[specId].Name
			}
			items[i] += ":"
			itemId, _ := util.I32Err(strconv.Atoi(v[idx+1:]))
			if n, ok := iName[itemId]; ok && n != "" {
				items[i] += n
			} else {
				items[i] += siMap[itemId].Value
			}
		}
		v.SpecWord = strings.TrimSpace(strings.Join(items, " "))

		// 标题为空，则自动设置
		if strings.TrimSpace(v.Title) == "" {
			titArr := make([]string, l+1)
			titArr[0] = it.Title
			for i, v := range arr {
				ii := strings.Index(v, ":")
				iid, _ := util.I32Err(strconv.Atoi(v[ii+1:]))
				//先从自定义规格值获取文本，再从规格项获取
				if n, ok := iName[iid]; ok && n != "" {
					titArr[i+1] = n
				} else if im := siMap[iid]; im != nil {
					titArr[i+1] = im.Value
				}
			}
			v.Title = strings.TrimSpace(strings.Join(titArr, " "))
		}
	}
	return nil
}

// 获取商品的规格(从SKU中读取)
func (s *skuServiceImpl) GetSpecArray(skuArr []*item.Sku) (
	specArr promodel.SpecList) {
	sa, ia := s.GetSpecItemArray(skuArr) //规格与规格项编号的数组
	if l := len(sa); l > 0 {
		// 获取传入的规格信息,按传入规格名称SKU直到找到结果为止
		sName, iName := s.GetNameMap(skuArr)
		// 绑定规格
		specArr = make([]*promodel.Spec, l)
		imp := make(map[int32]int, l) //记录规格对应数组的索引
		for i, v := range sa {
			spec := s.proMRepo.GetSpec(v)
			if spec == nil {
				continue
			}
			spec.Items = []*promodel.SpecItem{}
			//重新绑定规格名字
			if n, ok := sName[spec.ID]; ok && n != "" {
				spec.Name = n
			}
			specArr[i] = spec
			imp[spec.ID] = i
		}
		// 绑定规格项
		for _, v := range ia {
			item := s.proMRepo.GetSpecItem(v)
			if item == nil {
				log.Println("no such spec product, id:", v)
				continue
			}
			if n, ok := iName[item.ID]; ok && n != "" {
				item.Value = n
			}
			i2 := imp[item.SpecId]
			specArr[i2].Items = append(specArr[i2].Items, item)
		}
		// 排序
		s.sortSpecArray(specArr)
		return specArr
	}
	return []*promodel.Spec{}
}

func (s *skuServiceImpl) sortSpecArray(arr promodel.SpecList) {
	for _, v := range arr {
		sort.Sort(v.Items)
	}
	sort.Sort(arr)
}

// 获取规格选择HTML
func (s *skuServiceImpl) GetSpecHtml(spec promodel.SpecList) string {
	return s.su.GetSpecHtml(spec)
}

// 获取规格JSON数据
func (s *skuServiceImpl) GetSpecJson(spec promodel.SpecList) []byte {
	arr := iJsonUtil.getSpecJdo(spec)
	data, _ := json.Marshal(arr)
	return data
}

// 获取SKU的JSON字符串
func (s *skuServiceImpl) GetSkuJson(skuArr []*item.Sku) []byte {
	arr := iJsonUtil.getSkuJdo(skuArr)
	b, _ := json.Marshal(arr)
	return b
}

type skuTempStruct struct {
	Spec []*promodel.Spec
}

type skuServiceUtil struct {
	service *skuServiceImpl
	skuTemp *template.Template
}

func newSkuUtil(s *skuServiceImpl) *skuServiceUtil {
	return (&skuServiceUtil{
		service: s,
	}).init()
}

// 初始化模板,需使用text/template
func (s *skuServiceUtil) init() *skuServiceUtil {
	var err error
	s.skuTemp = &template.Template{}
	htm := `{{range $i,$v := .Spec}}
        <div class="mod-sku-spec">
            <div class="spec-label">{{$v.Name}}</div>
            <div class="spec-option">
                {{range $i2,$vi := $v.Items}}
                <a class="spec-option-check spec-option-item" href="javascript:void(0)" sid="{{$v.ID}}:{{$vi.ID}}">
                    {{$vi.Value}}
                </a>
                {{end}}
            </div>
        </div>
       {{end}}
    `
	s.skuTemp, err = s.skuTemp.Parse(htm)
	if err != nil {
		log.Println("convert sku template error:", err.Error())
	}
	return s
}

// 获取规格选择HTML
func (s *skuServiceUtil) GetSpecHtml(spec promodel.SpecList) string {
	if len(spec) == 0 {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	err := s.skuTemp.Execute(buf, &skuTempStruct{spec})
	if err != nil {
		log.Println("execute sku template error:", err.Error())
	}
	return buf.String()
}
