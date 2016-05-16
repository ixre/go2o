/**
 * Copyright 2015 @ z3q.net.
 * name : content_rep.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

type IContentRep interface {
	// 获取内容
	GetContent(merchantId int) IContent

	// 根据编号获取页面
	GetPageById(merchantId, id int) *ValuePage

	// 根据标识获取页面
	GetPageByStringIndent(merchantId int, indent string) *ValuePage

	// 删除页面
	DeletePage(merchantId, id int) error

	// 保存页面
	SavePage(merchantId int, v *ValuePage) (int, error)
}
