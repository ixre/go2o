/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-11-22 20:01
 * description :
 * history :
 */

package merchant

type (
    IShop interface {
        GetDomainId() int

        GetValue() ValueShop

        SetValue(*ValueShop) error

        //	// 获取经维度
        //	GetLngLat() (float64, float64)
        //
        //	// 是否可以配送
        //	// 返回是否可以配送，以及距离(米)
        //	CanDeliver(lng, lat float64) (bool, int)
        //
        //	// 是否可以配送
        //	// 返回是否可以配送，以及距离(米)
        //	CanDeliverTo(address string) (bool, int)

        Save() (int, error)
    }
    //门店
    ValueShop struct {
        Id         int    `db:"id" pk:"yes" auto:"yes"`
        MerchantId  int    `db:"merchant_id"`
        Name       string `db:"name"`
        Address    string `db:"address"`
        Phone      string `db:"phone"`
        SortNumber int    `db:"sort_number"`
        State      int    `db:"state"`
        CreateTime int64  `db:"create_time"`

        //    // 位置(经度+"/"+纬度)
        //    Location string `db:"location"`
        //
        //    // 配送最大半径(公里)
        //    DeliverRadius int `db:"deliver_radius"`
    }
    //合作商网站配置
    SiteConf struct {
        //合作商编号
        MerchantId int `db:"merchant_id" auto:"no" pk:"yes"`

        //主机
        Host       string `db:"host"`

        //前台Logo
        Logo       string `db:"logo"`

        //首页标题
        IndexTitle string `db:"index_title"`

        //子页面标题
        SubTitle   string `db:"sub_title"`

        //状态: 0:暂停  1：正常
        State      int    `db:"state"`
        StateHtml  string `db:"state_html"`
    }
)
