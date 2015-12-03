/**
 * Copyright 2015 @ z3q.net.
 * name : api_info.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package partner

// 商户接口信息
type ApiInfo struct {
	// 商户编号
	PartnerId int `db:"partner_id" pk:"yes" auto:"no"`
	// 商户接口编号(10位数字)
	ApiId string `db:"api_id"`
	// 密钥
	ApiSecret string `db:"api_secret"`
	// IP白名单
	WhiteList string `db:"white_list"`
	// 是否启用,0:停用,1启用
	Enabled int `db:"enabled"`
}
