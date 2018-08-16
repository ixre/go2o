/**
 * Copyright 2015 @ z3q.net.
 * name : json_c.go
 * author : jarryliu
 * date : 2016-04-25 23:09
 * description :
 * history :
 */
package shared

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/jsix/goex/echox"
	"github.com/jsix/gof"
	"github.com/jsix/gof/crypto"
	"github.com/jsix/gof/storage"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/ad"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/rsi"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	//todo: ??? 设置为可配置
	maxSeconds int64 = 120
)

func init() {
	gob.Register(map[string]map[string]interface{}{})
	gob.Register(ad.ValueGallery{})
	gob.Register(ad.Ad{})
	gob.Register([]*valueobject.Goods{})
	gob.Register(valueobject.Goods{})
	gob.Register(ad.HyperLink{})
	gob.Register(ad.Image{})
}

type JsonC struct {
	gof.App
	mux *sync.RWMutex
}

func NewJsonC() *JsonC {
	return &JsonC{
		App: gof.CurrentApp,
		mux: new(sync.RWMutex),
	}
}

func getMd5(s string) string {
	return crypto.Md5([]byte(s))[8:16]
}

func (j *JsonC) getMultiParams(s string) (p string, size, begin int) {
	arr := strings.Split(s, "*")
	l := len(arr)
	if l == 1 {
		p = s
		size = 10 //返回默认10条
	} else {
		p = arr[0]
		size, _ = strconv.Atoi(arr[1])
		if l > 2 {
			begin, _ = strconv.Atoi(arr[2])
		}
	}
	return p, size, begin
}

func (j *JsonC) unmarshal(sto storage.Interface, key string, dst interface{}) error {
	jsStr, err := sto.GetString(key)
	if err == nil {
		err = json.Unmarshal([]byte(jsStr), &dst)
	}
	return err
}

// 商城/商铺分类JSON，shop_id为0，则返回商城的分类
// todo: ??? gob编码提示错误
func (j *JsonC) ShopCat(c *echox.Context) error {
	parentId, _ := util.I32Err(strconv.Atoi(c.FormValue("parent_id")))
	shopId, _ := util.I32Err(strconv.Atoi(c.FormValue("shop_id")))
	var list []*ttype.SCategory
	key := fmt.Sprintf("go2o:repo:cat:%d:json:%d", shopId, parentId)
	sto := c.App.Storage()
	if err := j.unmarshal(sto, key, &list); err != nil {
		//if err := sto.Get(key,*list);err != nil{
		if parentId == 0 {
			list = rsi.ProductService.GetBigCategories(shopId)
		} else {
			list = rsi.ProductService.GetChildCategories(shopId, parentId)
		}
		//sto.Set(key,list)
		var d []byte
		d, err = json.Marshal(list)
		sto.Set(key, string(d))
		//log.Println("---- 更新分类缓存 ", err)
	}
	return c.JSON(http.StatusOK, list)
}

func (j *JsonC) Get_shop(c *echox.Context) error {
	typeParams := strings.TrimSpace(c.FormValue("params"))
	types := strings.Split(typeParams, "|")
	result := make(map[string]interface{}, len(types))
	key := fmt.Sprint("go2o:repo:shop:front:glob_%s", typeParams)
	sto := c.App.Storage()
	//从缓存中读取
	if err := sto.Get(key, &result); err != nil {
		ss := rsi.ShopService
		for _, t := range types {
			p, size, begin := j.getMultiParams(t)
			switch p {
			case "new-shop":
				_, result[p] = ss.PagedOnBusinessOnlineShops(
					begin, begin+size, "", "sp.create_time DESC")
			case "hot-shop":
				_, result[p] = ss.PagedOnBusinessOnlineShops(
					begin, begin+size, "", "")
			}
		}
		sto.SetExpire(key, result, maxSeconds)
	}
	return c.Debug(c.JSON(http.StatusOK, result))
}

// 商品
func (j *JsonC) Get_goods(c *echox.Context) error {
	shopId, _ := util.I32Err(strconv.Atoi(c.FormValue("shop_id")))
	typeParams := strings.TrimSpace(c.FormValue("params"))
	types := strings.Split(typeParams, "|")
	result := make(map[string]interface{}, len(types))
	key := fmt.Sprint("go2o:repo:item:fc:%d_%s", shopId, typeParams)
	sto := c.App.Storage()
	if err := sto.Get(key, &result); err != nil {
		//从缓存中读取
		ss := rsi.ItemService
		for _, t := range types {
			p, size, begin := j.getMultiParams(t)
			switch p {
			case "new-goods":
				_, result[p] = ss.GetPagedOnShelvesGoods__(shopId,
					-1, begin, begin+size, "it.id DESC")
			case "hot-sales":
				_, result[p] = ss.GetPagedOnShelvesGoods__(shopId,
					-1, begin, begin+size, "it.sale_num DESC")
			}
		}
		//sto.SetExpire(key, result, maxSeconds)
	}
	return c.Debug(c.JSON(http.StatusOK, result))
}

// 最新店铺
func (j *JsonC) Get_Newgoods(c *echox.Context) error {
	shopId, _ := util.I32Err(strconv.Atoi(c.FormValue("shop_id")))
	begin, _ := strconv.Atoi(c.FormValue("begin"))
	size, _ := strconv.Atoi(c.FormValue("size"))
	ss := rsi.ItemService
	_, result := ss.GetPagedOnShelvesGoods__(shopId,
		-1, begin, begin+size, "it.id DESC")

	return c.Debug(c.JSON(http.StatusOK, result))
}

// 最新商品
func (j *JsonC) Get_Newshop(c *echox.Context) error {
	begin, _ := strconv.Atoi(c.FormValue("begin"))
	size, _ := strconv.Atoi(c.FormValue("size"))
	ss := rsi.ShopService
	_, result := ss.PagedOnBusinessOnlineShops(
		begin, begin+size, "", "sp.create_time DESC")

	return c.Debug(c.JSON(http.StatusOK, result))
}

//最热商品
func (j *JsonC) Get_hotGoods(c *echox.Context) error {
	shopId, _ := util.I32Err(strconv.Atoi(c.FormValue("shop_id")))
	ss := rsi.ItemService
	begin, _ := strconv.Atoi(c.FormValue("begin"))
	size, _ := strconv.Atoi(c.FormValue("size"))
	_, result := ss.GetPagedOnShelvesGoods__(shopId,
		-1, begin, begin+size, "it.sale_num DESC")
	return c.Debug(c.JSON(http.StatusOK, result))
}

func (j *JsonC) Mch_goods(c *echox.Context) error {
	typeParams := strings.TrimSpace(c.FormValue("params"))
	types := strings.Split(typeParams, "|")
	mchId, _ := util.I32Err(strconv.Atoi(c.FormValue("mch_id")))
	result := make(map[string]interface{}, len(types))
	key := fmt.Sprint("go2o:repo:sg:front:%d_%s", mchId, typeParams)
	sto := c.App.Storage()
	if err := sto.Get(key, &result); err != nil {
		//从缓存中读取
		ss := rsi.ItemService
		for _, t := range types {
			p, size, begin := j.getMultiParams(t)
			switch p {
			case "new-goods":
				_, result[p] = ss.GetShopPagedOnShelvesGoods(mchId,
					-1, begin, begin+size, "it.id DESC")
			case "hot-sales":
				_, result[p] = ss.GetShopPagedOnShelvesGoods(mchId,
					-1, begin, begin+size, "it.sale_num DESC")
			}
		}
		sto.SetExpire(key, result, maxSeconds)
	}
	return c.Debug(c.JSON(http.StatusOK, result))
}

// 获取销售标签获取商品
func (j *JsonC) SaleLabelGoods(c *echox.Context) error {
	codeParams := strings.TrimSpace(c.FormValue("params"))
	codes := strings.Split(codeParams, "|")
	mchId, _ := util.I32Err(strconv.Atoi(c.FormValue("mch_id")))
	result := make(map[string]interface{}, len(codes))

	key := fmt.Sprint("go2o:repo:stg:front:%d--%s", mchId, getMd5(codeParams))
	sto := c.App.Storage()
	if err := sto.Get(key, &result); err != nil {
		//从缓存中读取
		for _, param := range codes {
			code, size, begin := j.getMultiParams(param)
			list := rsi.ItemService.GetValueGoodsBySaleLabel(
				code, "", begin, begin+size)
			result[code] = list
		}
		sto.SetExpire(key, result, maxSeconds)
	}
	return c.Debug(c.JSON(http.StatusOK, result))
}
