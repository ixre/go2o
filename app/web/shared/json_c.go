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
	"github.com/jsix/gof"
	"github.com/jsix/gof/crypto"
	"github.com/jsix/gof/storage"
	"go2o/core/domain/interface/ad"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/dto"
	"go2o/core/service/dps"
	"go2o/x/echox"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	//todo: ??? 设置为可配置
	maxSeconds int64 = 10
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

// 广告
func (j *JsonC) Ad(ctx *echox.Context) error {
	namesParams := strings.TrimSpace(ctx.Query("keys"))
	names := strings.Split(namesParams, "|")
	userId, _ := strconv.Atoi(ctx.Query("ad_user"))
	as := dps.AdService
	result := make(map[string]*ad.AdDto, len(names))
	key := fmt.Sprintf("go2o:rep:ad:%d:front:%s", userId, getMd5(namesParams))
	sto := ctx.App.Storage()
	if err := sto.Get(key, &result); err != nil {
		//从缓存中读取
		for _, n := range names {
			//分别绑定广告
			dto := as.GetAdAndDataByKey(userId, n)
			if dto == nil {
				result[n] = nil
				continue
			}
			result[n] = dto
		}
		seconds := dps.BaseService.GetRegistry().CacheAdMaxAge
		sto.SetExpire(key, result, seconds)
		//log.Println("---- 更新广告缓存 ",err)
	}
	return ctx.JSON(http.StatusOK, result)
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
func (j *JsonC) ShopCat(ctx *echox.Context) error {
	parentId, _ := strconv.Atoi(ctx.Form("parent_id"))
	shopId, _ := strconv.Atoi(ctx.Form("shop_id"))
	list := []dto.Category{}
	key := fmt.Sprintf("go2o:rep:cat:%d:json:%d", shopId, parentId)
	sto := ctx.App.Storage()
	if err := j.unmarshal(sto, key, &list); err != nil {
		//if err := sto.Get(key,*list);err != nil{
		if parentId == 0 {
			list = dps.SaleService.GetBigCategories(shopId)
		} else {
			list = dps.SaleService.GetChildCategories(shopId, parentId)
		}
		//sto.Set(key,list)
		var d []byte
		d, err = json.Marshal(list)
		sto.Set(key, string(d))
		//log.Println("---- 更新分类缓存 ", err)
	}
	return ctx.JSON(http.StatusOK, list)
}

func (j *JsonC) Get_shop(ctx *echox.Context) error {
	typeParams := strings.TrimSpace(ctx.Form("params"))
	types := strings.Split(typeParams, "|")
	result := make(map[string]interface{}, len(types))
	key := fmt.Sprint("go2o:rep:shop:front:glob_%s", typeParams)
	sto := ctx.App.Storage()
	//从缓存中读取
	if err := sto.Get(key, &result); err != nil {
		ss := dps.ShopService
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
	return ctx.Debug(ctx.JSON(http.StatusOK, result))
}

// 商品
func (j *JsonC) Get_goods(ctx *echox.Context) error {
	shopId, _ := strconv.Atoi(ctx.Form("shop_id"))
	typeParams := strings.TrimSpace(ctx.Form("params"))
	types := strings.Split(typeParams, "|")
	result := make(map[string]interface{}, len(types))
	key := fmt.Sprint("go2o:rep:gs:fc3:%d_%s", shopId, typeParams)
	sto := ctx.App.Storage()
	if err := sto.Get(key, &result); err != nil {
		//从缓存中读取
		ss := dps.SaleService
		for _, t := range types {
			p, size, begin := j.getMultiParams(t)
			switch p {
			case "new-goods":
				_, result[p] = ss.GetPagedOnShelvesGoods(shopId,
					-1, begin, begin+size, "gs_goods.id DESC")
			case "hot-sales":
				_, result[p] = ss.GetPagedOnShelvesGoods(shopId,
					-1, begin, begin+size, "gs_goods.sale_num DESC")
			}
		}
		sto.SetExpire(key, result, maxSeconds)
	}
	return ctx.Debug(ctx.JSON(http.StatusOK, result))
}

func (j *JsonC) Mch_goods(ctx *echox.Context) error {
	typeParams := strings.TrimSpace(ctx.Form("params"))
	types := strings.Split(typeParams, "|")
	mchId, _ := strconv.Atoi(ctx.Form("mch_id"))
	result := make(map[string]interface{}, len(types))

	key := fmt.Sprint("go2o:rep:sg:front:%d_%s", mchId, typeParams)
	sto := ctx.App.Storage()
	if err := sto.Get(key, &result); err != nil {
		//从缓存中读取
		ss := dps.SaleService
		for _, t := range types {
			p, size, begin := j.getMultiParams(t)
			switch p {
			case "new-goods":
				_, result[p] = ss.GetShopPagedOnShelvesGoods(mchId,
					-1, begin, begin+size, "gs_goods.id DESC")
			case "hot-sales":
				_, result[p] = ss.GetShopPagedOnShelvesGoods(mchId,
					-1, begin, begin+size, "gs_goods.sale_num DESC")
			}
		}
		sto.SetExpire(key, result, maxSeconds)
		/////////
	}
	return ctx.Debug(ctx.JSON(http.StatusOK, result))
}

// 获取销售标签获取商品
func (j *JsonC) SaleLabelGoods(ctx *echox.Context) error {
	codeParams := strings.TrimSpace(ctx.Form("params"))
	codes := strings.Split(codeParams, "|")
	mchId, _ := strconv.Atoi(ctx.Form("mch_id"))
	result := make(map[string]interface{}, len(codes))

	key := fmt.Sprint("go2o:rep:stg:front:%d--%s", mchId, getMd5(codeParams))
	sto := ctx.App.Storage()
	if err := sto.Get(key, &result); err != nil {
		//从缓存中读取
		log.Println(err)
		for _, param := range codes {
			code, size, begin := j.getMultiParams(param)
			list := dps.SaleService.GetValueGoodsBySaleLabel(
				mchId, code, "", begin, begin+size)
			result[code] = list
		}
		sto.SetExpire(key, result, maxSeconds)
	}
	return ctx.Debug(ctx.JSON(http.StatusOK, result))
}
