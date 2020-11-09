package model

// APP产品
type AppProd struct {
	// 产品编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 产品名称
	ProdName string `db:"prod_name"`
	// 产品描述
	ProdDes string `db:"prod_des"`
	// 最新的版本ID
	LatestVid int64 `db:"latest_vid"`
	// 正式版文件hash值
	Md5Hash string `db:"md5_hash"`
	// 发布下载页面地址
	PublishUrl string `db:"publish_url"`
	// 正式版文件地址
	StableFileUrl string `db:"stable_file_url"`
	// 内测版文件地址
	AlphaFileUrl string `db:"alpha_file_url"`
	// 每夜版文件地址
	NightlyFileUrl string `db:"nightly_file_url"`
	// 更新时间
	UpdateTime int64 `db:"update_time"`
}


// APP版本
type AppVersion struct {
	// 编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 产品
	ProductId int64 `db:"product_id"`
	// 更新通道, stable|beta|nightly
	Channel int16 `db:"channel"`
	// 版本号
	Version string `db:"version"`
	// 数字版本
	VersionCode int `db:"version_code"`
	// 是否强制升级
	ForceUpdate int16 `db:"force_update"`
	// 更新内容
	UpdateContent string `db:"update_content"`
	// 发布时间
	CreateTime int64 `db:"create_time"`
}
