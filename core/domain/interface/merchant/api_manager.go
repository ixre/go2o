/**
 * Copyright 2015 @ to2.net.
 * name : api_manager
 * author : jarryliu
 * date : 2016-05-27 13:23
 * description :
 * history :
 */
package merchant

type (
	// 商户接口信息
	ApiInfo struct {
		// 商户编号
		MerchantId int32 `db:"mch_id" pk:"yes" auto:"no"`
		// 商户接口编号(10位数字)
		ApiId string `db:"api_id"`
		// 密钥
		ApiSecret string `db:"api_secret"`
		// IP白名单
		WhiteList string `db:"white_list"`
		// 是否启用,0:停用,1启用
		Enabled int `db:"enabled"`
	}

	// Api接口管理器
	IApiManager interface {
		// 获取API信息,管理员可停用。
		GetApiInfo() ApiInfo

		// 保存API信息, 一般情况只有内部保存,其他为查看权限
		SaveApiInfo(*ApiInfo) error

		// 启用API权限
		EnableApiPerm() error

		// 禁用API权限
		DisableApiPerm() error
	}
)
