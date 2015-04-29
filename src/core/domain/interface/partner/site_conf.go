/**
 * Copyright 2014 @ ops.
 * name :
 * author : jarryliu
 * date : 2013-11-11 21:46
 * description :
 * history :
 */

package partner

//合作商网站配置
type SiteConf struct {
	//合作商编号
	PartnerId int `db:"pt_id" auto:"no" pk:"yes"`

	//主机
	Host string `db:"host"`

	//前台Logo
	Logo string `db:"logo"`

	//首页标题
	IndexTitle string `db:"index_title"`

	//子页面标题
	SubTitle string `db:"sub_title"`

	//状态: 0:暂停  1：正常
	State     int    `db:"state"`
	StateHtml string `db:"state_html"`
}
