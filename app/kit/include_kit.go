package kit

import (
	"bytes"
	"fmt"
	"github.com/jsix/gof"
	"github.com/jsix/gof/storage"
	"github.com/jsix/gof/util"
	"go2o/app/cache"
	"go2o/core/dao/model"
	"go2o/core/domain/interface/content"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/infrastructure/format"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"go2o/core/variable"
	"html/template"
	ht "html/template"
	"strings"
	"sync"
)

var (
	titleSuffix  string
	cacheSeconds int64 = 300
	hashSet            = storage.NewHashStorage()
)

type templateIncludeKitWrapper struct {
	FuncMap    ht.FuncMap
	Middleware func(path string) bool
}

type templateIncludeToolkit struct {
	// 入口链接字点
	entryUrlMap map[string]string
	mutex       sync.Mutex
	rwMut       sync.RWMutex
}

// 返回模板函数
func (t *templateIncludeToolkit) getFuncMap() ht.FuncMap {
	fm := make(map[string]interface{})
	fm["alias"] = t.alias
	fm["script"] = t.scriptTag
	fm["css"] = t.cssTag
	fm["cssXY"] = t.cssXY
	fm["entry"] = t.entryUrl
	fm["catTree"] = t.CatTree
	fm["catParent"] = t.catParent
	fm["catChild"] = t.catChild
	fm["catBrand"] = t.CatBrand
	fm["catItems"] = t.catItems
	fm["modelBrands"] = t.modelBrands
	fm["portalNav"] = t.portalNav
	fm["pageTitle"] = t.pageTitle
	fm["floorAd"] = t.floorAd
	fm["boolInt"] = t.boolInt
	fm["isEmpty"] = t.isEmpty
	fm["img"] = t.imgLabel
	fm["propQuery"] = t.propQuery
	fm["rawHtml"] = t.rawHtml
	fm["resUrl"] = t.resUrl
	fm["hotSaleItems"] = t.hotSaleItems
	fm["randItems"] = t.randItems
	fm["productAttrs"] = t.productAttrs
	fm["add"] = t.add
	fm["multi"] = t.multi
	fm["articles"] = t.articles
	fm["mathRemain"] = t.mathRemain
	fm["kv"] = t.kv
	fm["registry"] = t.registry
	fm["priceStr"] = t.priceStr

	return fm
}

func (t *templateIncludeToolkit) getRds() storage.Interface {
	return gof.CurrentApp.Storage()
}

// 缓存子模板
func (t *templateIncludeToolkit) includeMiddle(path string) bool {
	key := "go2o:front:portal:inc:" + path
	sto := gof.CurrentApp.Storage()
	_, err := sto.GetInt(key)
	if err == nil {
		return true
	}
	sto.SetExpire(key, 1, cache.DefaultMaxSeconds)
	return false
}

// 别名
func (t *templateIncludeToolkit) alias(s string) string {
	switch s {
	case "WalletAccount":
		return variable.AliasWalletAccount
	case "GrowAccount":
		return variable.AliasGrowthAccount
	case "BalanceAccount":
		return variable.AliasBalanceAccount
	case "TradeOrder":
		return variable.AliasTradeOrder
	}
	return s
}

// CSS标签
func (t *templateIncludeToolkit) cssTag(s string) template.HTML {
	registry := RPC.RegistryMap(variable.DStaticServer, variable.DUrlHash)
	staticServe := registry[variable.DStaticServer]
	urlHash := registry[variable.DUrlHash]
	buf := bytes.NewBufferString("")
	arr := strings.Split(s, ",")
	for i, v := range arr {
		if i != 0 {
			buf.WriteString("\n")
		}
		buf.WriteString("<link rel=\"StyleSheet\" type=\"text/css\" href=\"")
		buf.WriteString(staticServe)
		buf.WriteString(v)
		buf.WriteString("?hash=")
		buf.WriteString(urlHash)
		buf.WriteString("\"/>")

	}
	return template.HTML(buf.String())
}

// 脚本标签
func (t *templateIncludeToolkit) scriptTag(s string) template.HTML {
	registry := RPC.RegistryMap(variable.DStaticServer, variable.DUrlHash)
	staticServe := registry[variable.DStaticServer]
	urlHash := registry[variable.DUrlHash]
	buf := bytes.NewBufferString("")
	arr := strings.Split(s, ",")
	for i, v := range arr {
		if i != 0 {
			buf.WriteString("\n")
		}
		buf.WriteString("<script type=\"text/javascript\" src=\"")
		buf.WriteString(staticServe)
		buf.WriteString(v)
		buf.WriteString("?hash=")
		buf.WriteString(urlHash)
		buf.WriteString("\"></script>")
	}
	return template.HTML(buf.String())
}

// 将坐标(x,y)转换为CSS背景坐标
func (t *templateIncludeToolkit) cssXY(xy string) string {
	arr := strings.Split(xy, ",")
	if len(arr) == 2 {
		sa := []string{arr[0], "px ", arr[1], "px"}
		return strings.Join(sa, "")
	}
	return "0px 0px"
}

// 入口URL
func (t *templateIncludeToolkit) entryUrl(k string) string {
	key := k
	switch strings.TrimSpace(k) {
	case "retail", "retail_portal", "retailPortal":
		key = variable.DRetailPortal
	case "retail_m", "retail_portal_m":
		key = variable.DRetailMobilePortal
	case "wholesale", "wholesale_portal", "wholesalePortal":
		key = variable.DWholesalePortal
	case "image_serve", "img_serve", "img":
		key = variable.DImageServer
	case "static_serve", "static":
		key = variable.DStaticServer
	}
	t.rwMut.RLock()
	if t.entryUrlMap != nil {
		if v, ok := t.entryUrlMap[k]; ok {
			t.rwMut.RUnlock()
			return v
		}
	}
	t.rwMut.RUnlock()
	t.rwMut.Lock()
	if t.entryUrlMap == nil {
		t.entryUrlMap = make(map[string]string)
	}
	v := RPC.RegistryMap(key)[key]
	t.entryUrlMap[key] = v
	t.rwMut.Unlock()
	return v
}

// 去掉虚拟未启用的分类
func (t *templateIncludeToolkit) fixCatTree(cat *product.Category) {
	catArr := []*product.Category{}
	for _, v := range cat.Children {
		if v.Enabled == 1 {
			catArr = append(catArr, v)
			t.fixCatTree(v)
		}
	}
	cat.Children = catArr
}

// 分类树形
func (t *templateIncludeToolkit) CatTree(parentId int32) product.Category {
	c := rsi.ProductService.CategoryTree(parentId)
	if c != nil {
		t.fixCatTree(c)
		return *c
	}
	return product.Category{}
}

// 获取分类的品牌
func (t *templateIncludeToolkit) CatBrand(catId int32, num int32) []*promodel.ProBrand {
	key := fmt.Sprintf("go2o:portal:cache:cat-brands-%d-%d", catId, num)
	_, err := t.getRds().GetInt(key)
	if err == nil {
		r, err := hashSet.GetRaw(key)
		if err == nil {
			return r.([]*promodel.ProBrand)
		}
	}
	arr := rsi.ProductService.GetCatBrands(catId)
	if num > 0 && int(num) < len(arr) {
		arr = arr[:num]
	}
	hashSet.Set(key, arr)
	t.getRds().SetExpire(key, 1, cacheSeconds)
	return arr
}

// 获取模型的品牌
func (t *templateIncludeToolkit) modelBrands(proModel int32) []*promodel.ProBrand {
	return rsi.ProductService.GetModelBrands(proModel)
}

// 栏目上级栏目
func (t *templateIncludeToolkit) catParent(catId int32) []*product.Category {
	s := rsi.ProductService
	arr := []*product.Category{}
	for catId > 0 {
		cat := s.GetCategory(0, catId)
		if cat != nil {
			arr = append([]*product.Category{cat}, arr...)
			catId = cat.ParentId
		} else {
			break
		}
	}
	return arr
}

// 获取栏目下级分类
func (t *templateIncludeToolkit) catChild(catId int32) []*ttype.SCategory {
	return rsi.ProductService.GetChildCategories(0, catId)
}

// 门户导航链接
func (t *templateIncludeToolkit) portalNav(navType int32) []*model.PortalNav {
	return rsi.PortalService.SelectPortalNav(navType)
}

// 页面标题
func (t *templateIncludeToolkit) pageTitle(tit string) string {
	if titleSuffix == "" {
		trans, cli, err := thrift.FoundationServeClient()
		if err == nil {
			defer trans.Close()
			r, _ := cli.GetRegistryMapV1(thrift.Context, []string{"PlatformName"})
			titleSuffix = r["PlatformName"]
		}
	}
	if tit == "" {
		return titleSuffix
	}
	return tit + "-" + titleSuffix
}

// 拼接属性URL-Query
func (t *templateIncludeToolkit) propQuery(query string, k interface{}, v interface{}) string {
	key := util.Str(k)
	val := util.Str(v)
	// 如没有属性
	if query == "" {
		s := []string{key, "E", val}
		return strings.Join(s, "")
	}
	keyI := strings.Index(query, key+"E")
	if keyI == -1 {
		s := []string{query, "_", key, "E", val}
		return strings.Join(s, "")
	}
	s := []string{query[:keyI]}
	s = append(s, key)
	s = append(s, "E")
	s = append(s, val)
	//查找下一个"_"的位置
	afterStr := query[keyI:]
	if eI := strings.Index(afterStr, "_"); eI != -1 {
		s = append(s, afterStr[eI:])
		//log.Println("---",key,keyI,query[keyI:],eI)
	}
	return strings.Join(s, "")
}

// 移除属性查询
func (t *templateIncludeToolkit) RemovePropQuery(query string, k interface{}) string {
	key := util.Str(k)
	keyI := strings.Index(query, key+"E")
	if keyI != -1 {
		s := []string{}
		//查找下一个"_"的位置,如果为-1,既末尾
		afterStr := query[keyI:]
		si := strings.Index(afterStr, "&")
		if si != -1 {
			afterStr = afterStr[:si]
		}
		eI := strings.Index(afterStr, "_")
		//最后的属性
		if eI == -1 {
			s = append(s, query[:keyI-1])
		} else {
			s = append(s, query[:keyI])
			s = append(s, afterStr[eI+1:])
		}
		if si != -1 {
			s = append(s, query[keyI+si:])
		}
		return strings.Join(s, "")
	}
	return query
}

// 判断是否为true
func (t *templateIncludeToolkit) boolInt(i int32) bool {
	return i > 0
}

// 加法
func (t *templateIncludeToolkit) add(x, y int) int {
	return x + y
}

// 乘法
func (t *templateIncludeToolkit) multi(x, y interface{}) interface{} {
	fx, ok := x.(float64)
	if ok {
		switch y.(type) {
		case float32:
			return fx * float64(y.(float32))
		case float64:
			return fx * y.(float64)
		case int:
			return fx * float64(y.(int))
		case int32:
			return fx * float64(y.(int32))
		case int64:
			return fx * float64(y.(int64))
		}
	}
	panic("not support")
}

// I32转为字符
func (t *templateIncludeToolkit) str(i interface{}) string {
	return util.Str(i)
}

// 是否为空
func (t *templateIncludeToolkit) isEmpty(s string) bool {
	if s == "" {
		return true
	}
	return strings.TrimSpace(s) == ""
}

// 图片
func (t *templateIncludeToolkit) imgLabel(img string) ht.HTML {
	htm := ""
	if img != "" {
		htm = fmt.Sprintf("<img src=\"%s\"/>",
			format.GetResUrl(img))
	}
	return ht.HTML(htm)
}

// 资源地址
func (t *templateIncludeToolkit) resUrl(u string) string {
	return format.GetResUrl(u)
}

// 转换为HTML
func (t *templateIncludeToolkit) rawHtml(v interface{}) ht.HTML {
	return ht.HTML(util.Str(v))
}

// 获取销售排行商品
func (t *templateIncludeToolkit) hotSaleItems(catId int32, quantity int32) []*ttype.SOldItem {
	_, arr := rsi.ItemService.GetPagedOnShelvesItem(item.ItemNormal,
		catId, 0, quantity, "", "item_info.sale_num DESC")
	return arr
}

// 获取随机商品
func (t *templateIncludeToolkit) randItems(catId int32, quantity int32) []*ttype.SOldItem {
	if catId <= 0 {
		catId = 0
	}
	return rsi.ItemService.GetRandomItem(catId, quantity, "")
}

// 获取大分类商品的
func (t *templateIncludeToolkit) catItems(catId int32, quantity int32) []*ttype.SOldItem {
	key := fmt.Sprintf("go2o:portal:cache:cat-items-%d-%d", catId, quantity)
	_, err := t.getRds().GetInt(key)
	if err == nil {
		r, err := hashSet.GetRaw(key)
		if err == nil {
			return r.([]*ttype.SOldItem)
		}
	}
	arr := rsi.ItemService.GetBigCatItems(catId, quantity, "")
	hashSet.Set(key, arr)
	t.getRds().SetExpire(key, 1, cacheSeconds)
	return arr
}

// 获取产品属性
func (t *templateIncludeToolkit) productAttrs(productId int64) []ttype.Pair {
	var arr []ttype.Pair
	attrs := rsi.ProductService.GetAttrArray(productId)
	for _, v := range attrs {
		attr := rsi.ProductService.GetAttr(v.AttrId)
		arr = append(arr, ttype.Pair{Key: attr.Name, Value: v.AttrWord})
	}
	return arr
}

// 获取文章列表
func (t *templateIncludeToolkit) articles(cat string, quantity int32) []*content.Article {
	_, arr := rsi.ContentService.PagedArticleList(cat, 0, int(quantity), "")
	return arr
}

//求余
func (t *templateIncludeToolkit) mathRemain(i int, j int) int {
	return i % j
}

// 根据键获取值
func (t *templateIncludeToolkit) kv(key string) string {
	r, _ := rsi.FoundationService.GetValue(thrift.Context, key)
	return r
}

func (t *templateIncludeToolkit) registry(keys ...string) map[string]string {
	trans, cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer trans.Close()
		r, _ := cli.GetRegistryMapV1(thrift.Context, keys)
		return r
	}
	return map[string]string{}
}

// 楼层广告设置
func (t *templateIncludeToolkit) floorAd(catId int32) string {
	return rsi.CommonDao.GetFloorAdPos(catId)
}

// 价格字符串
func (t *templateIncludeToolkit) priceStr(v float64) string {
	return format.DecimalToString(v)
}
