package model

type (
	// 导航类型
	PortalNavType struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 名称
		Name string `db:"name"`
	}
	// 门户导航
	PortalNav struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 文本
		Text string `db:"text"`
		// 地址
		Url string `db:"url"`
		// 打开目标
		Target string `db:"target"`
		Image  string `db:"image"`
		// 导航类型: 1为电脑，2为手机端
		NavType int32 `db:"nav_type"`
	}
	// 二维码模板
	CommQrTemplate struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
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
)
