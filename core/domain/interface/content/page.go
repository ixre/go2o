package content

import (
	"github.com/ixre/go2o/core/infrastructure/domain"
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
		GetPageById(tenantId, id int) IPage
		// GetPageByCode 根据标识获取页面
		GetPageByCode(tenantId int, code string) IPage
		// DeletePage 删除页面
		DeletePage(tenantId, id int) error
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
		Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
		// 用户编号,系统为0
		UserId int `json:"userId" db:"user_id" gorm:"column:user_id" bson:"userId"`
		// 标题
		Title string `json:"title" db:"title" gorm:"column:title" bson:"title"`
		// 标志
		Flag int `json:"flag" db:"flag" gorm:"column:flag" bson:"flag"`
		// 访问钥匙
		AccessKey string `json:"accessKey" db:"access_key" gorm:"column:access_key" bson:"accessKey"`
		// 页面代码
		Code string `json:"code" db:"code" gorm:"column:code" bson:"code"`
		// 关键词
		Keyword string `json:"keyword" db:"keyword" gorm:"column:keyword" bson:"keyword"`
		// 描述
		Description string `json:"description" db:"description" gorm:"column:description" bson:"description"`
		// 样式表路径
		CssPath string `json:"cssPath" db:"css_path" gorm:"column:css_path" bson:"cssPath"`
		// 是否启用
		Enabled int `json:"enabled" db:"enabled" gorm:"column:enabled" bson:"enabled"`
		// 内容
		Content string `json:"content" db:"content" gorm:"column:content" bson:"content"`
		// 更新时间
		UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
	}
)

func (p Page) TableName() string {
	return "arc_page"
}
