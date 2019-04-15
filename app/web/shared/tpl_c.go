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
	"github.com/ixre/goex/echox"
	"github.com/ixre/gof"
	"strings"
)

type TmlC struct {
	gof.App
}

//通过URL参数登录
//@member_id : 会员编号
//@token  :  密钥/令牌
//@device : 设备类型
func (t *TmlC) Blank(c *echox.Context) error {
	path := c.Request().URL.Path
	tpl := path[strings.Index(path[1:], "/")+1:]
	return c.RenderOK(tpl, c.NewData())
}
