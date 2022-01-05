package model

type (
	// PortalNavType 导航类型
	PortalNavType struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 名称
		Name string `db:"name"`
	}
	// PortalNav 门户导航
	PortalNav struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 文本
		Text string `db:"text"`
		// 地址
		Url string `db:"url"`
		// 打开目标
		Target string `db:"target"`
		// 图片
		Image string `db:"image"`
		// 导航类型: 1为电脑，2为手机端
		NavType int32 `db:"nav_type"`
		// 分组
		NavGroup string `db:"nav_group"`
	}

	// NavGroup 导航分组
	NavGroup struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 名称
		Name string `db:"name"`
	}
	// 二维码模板
	QrTemplate struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 模板标题
		Title string `db:"title"`
		// 背景图片
		BgImage string `db:"bg_image"`
		// 垂直偏离量
		OffsetX int `db:"offset_x"`
		// 垂直偏移量
		OffsetY int `db:"offset_y"`
		// 二维码模板文本
		Comment string `db:"comment"`
		// 回调地址
		CallbackUrl string `db:"callback_url"`
		// 是否启用
		Enabled int `db:"enabled"`
	}
	// 楼层广告设置
	PortalFloorAd struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 分类编号
		CatId int32 `db:"cat_id"`
		// 广告位编号
		PosId int32 `db:"pos_id"`
		// 广告顺序
		AdIndex int32 `db:"ad_index"`
	}

	// 楼层链接
	PortalFloorLink struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 分类编号
		CatId int32 `db:"cat_id"`
		// 文本
		Text string `db:"text"`
		// 链接地址
		LinkUrl string `db:"link_url"`
		// 打开方式
		Target string `db:"target"`
		// 序号
		SortNum int32 `db:"sort_num"`
	}
)

// SearchWord 热搜词
type SearchWord struct {
	// 编号
	Id int64 `db:"id" pk:"yes"`
	// Word
	Word string `db:"word"`
	// 搜索次数
	SearchCount int `db:"search_count"`
	// 1:启用　2:特殊显示 4: 手动创建
	Flag int `db:"flag"`
}
