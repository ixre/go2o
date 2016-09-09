/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package shared

import (
	"github.com/jsix/gof"
	"go2o/x/echox"
	"strings"
)

type TmlC struct {
	gof.App
}

//通过URL参数登录
//@member_id : 会员编号
//@token  :  密钥/令牌
//@device : 设备类型
func (this *TmlC) Blank(ctx *echox.Context) error {
	path := ctx.Request().URL.Path
	tpl := path[strings.Index(path[1:], "/")+1:]
	return ctx.RenderOK(tpl, ctx.NewData())
}
