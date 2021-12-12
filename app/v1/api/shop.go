package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ixre/gof"
	"github.com/ixre/gof/api"
	"go2o/core/service"
	"go2o/core/service/proto"
)

/**
 * Copyright 2009-2019 @ 56x.net
 * name : store.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-12-09 11:27
 * description :
 * history :
 */

var _ api.Handler = new(shopApi)

type shopApi struct {
	utils
}

func NewShopApi() *shopApi {
	return &shopApi{}
}

func (s shopApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"category":  s.shopCat,
		"favorite":  s.Favorite,
		"my_addrss": s.addressList,
	})
}

/**
 * @api {post} /shop/category 获取商城的商品分类
 * @apiGroup shop
 * @apiParam {int} parent_id 上级分类编号
 * @apiParam {int} shop_id 店铺编号
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (s shopApi) shopCat(ctx api.Context) interface{} {
	parentId := ctx.Form().GetInt("parent_id")
	shopId := ctx.Form().GetInt("shop_id")
	var list []*proto.SProductCategory
	key := fmt.Sprintf("go2o:repo:cat:%d:json:%d", shopId, parentId)
	sto := gof.CurrentApp.Storage()
	trans, cli, _ := service.ProductServiceClient()
	defer trans.Close()
	if err := sto.Get(key, &list); err != nil {
		if parentId == 0 {
			ret, _ := cli.GetChildren(context.TODO(), &proto.CategoryParentId{})
			list = ret.Value
		} else {
			ret, _ := cli.GetCategory(context.TODO(),
				&proto.Int64{Value: int64(parentId)})
			list = ret.Children
		}
		var d []byte
		d, err = json.Marshal(list)
		sto.SetExpire(key, string(d), 3600*24)
	}
	return list
}

/**
 * @api {post} /shop/favorite 收藏店铺
 * @apiGroup shop
 * @apiParam {int} item_id 店铺编号
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (s shopApi) Favorite(ctx api.Context) interface{} {
	memberId := getMemberId(ctx)
	id := ctx.Form().GetInt("shop_id")
	trans, cli, _ := service.MemberServiceClient()
	r, err := cli.Favorite(context.TODO(), &proto.FavoriteRequest{
		MemberId:     int64(memberId),
		FavoriteType: proto.FavoriteType_SHOP,
		ReferId:      int64(id),
	})
	trans.Close()
	if r.ErrCode > 0 {
		err = errors.New(r.ErrMsg)
	}
	if err != nil {
		return api.ResponseWithCode(1, err.Error())
	}
	return api.NewResponse(nil)
}

/*

// 登陆状态
func (s *serviceC) LoginState(c *echox.Context) error {
	mp := make(map[string]interface{})
	mobileReq := ut.MobileRequest(c.Request())
	mPrefix := types.StringCond(mobileReq, consts.DOMAIN_PREFIX_M_PASSPORT,
		consts.DOMAIN_PREFIX_PASSPORT)
	pstUrl := fmt.Sprintf("//%s%s", mPrefix, variable.Domain)
	memberId := getMemberId(c)
	if memberId <= 0 {
		registry, _ := impl.RegistryService.GetValues(context.TODO(),
			[]string{"PlatformName"})
		mp["PtName"] = registry["PlatformName"]
		mp["LoginUrl"] = pstUrl + "/auth/login"
		mp["RegisterUrl"] = pstUrl + "/register"
		mp["login"] = 0
	} else {
		mmUrl := fmt.Sprintf("//%s%s",
			consts.DOMAIN_PREFIX_MEMBER, variable.Domain)
		m, _ := impl.MemberService.GetProfile(context.TODO(),&proto.Int64{Value:int64(memberId)})
		mp["MMName"] = m.Name
		mp["LogoutUrl"] = pstUrl + "/auth/logout"
		mp["MMUrl"] = mmUrl
		mp["login"] = 1
	}
	return c.JSONP(http.StatusOK, c.QueryParam("callback"), mp)
}
*/

//

/**
 * @api {post} /shop/my_address 收货地址列表
 * @apiGroup shop
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (s shopApi) addressList(ctx api.Context) interface{} {
	memberId := getMemberId(ctx)
	trans, cli, _ := service.MemberServiceClient()
	defer trans.Close()
	address, _ := cli.GetAddressList(context.TODO(), &proto.Int64{Value: int64(memberId)})
	return address
}
