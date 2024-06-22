package content

import (
	"github.com/ixre/go2o/core/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

type (

	// IPageManager 页面管理器
	IPageManager interface {
		// CreatePage 创建页面
		CreatePage(*Page) IPage
		// GetPage 获取页面
		GetPage(id int) IPage
		// GetPageByCode 根据字符串标识获取页面
		GetPageByCode(indent string) IPage
		// DeletePage 删除页面
		DeletePage(id int) error
	}

	// IPageRepo 页面仓储
	IPageRepo interface {
		fw.Repository[Page]
		// GetPageById 根据编号获取页面
		GetPageById(zoneId, id int) IPage
		// GetPageByCode 根据标识获取页面
		GetPageByCode(zoneId int, code string) IPage
		// DeletePage 删除页面
		DeletePage(zoneId, id int) error
		// SavePage 保存页面
		SavePage(zondId int, v *Page) error
	}

	// IPage 页面
	IPage interface {
		domain.IDomain
		// GetValue 获取值
		GetValue() *Page
		// SetValue 设置值
		SetValue(*Page) error
		// Save 保存
		Save() (int, error)
	}
	// Page 页面
	Page struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 商户编号
		UserId int `db:"user_id"`
		// 标题
		Title string `db:"title"`
		// 字符标识
		Code string `db:"code"`
		// 浏览权限
		Flag int `db:"flag"`
		// 浏览钥匙
		AccessKey string `db:"access_key"`
		// 关键词
		KeyWord string `db:"keyword"`
		// 描述
		Description string `db:"description"`
		// 样式表地址
		CssPath string `db:"css_path"`
		// 内容
		Content string `db:"content"`
		// 是否启用
		Enabled int `db:"enabled"`
		// 修改时间
		UpdateTime int64 `db:"update_time"`
	}
)

func (p Page) TableName() string {
	return "arc_page"
}
