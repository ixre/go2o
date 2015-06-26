/**
 * Copyright 2015 @ S1N1 Team.
 * name : content_rep.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

type IContentRep interface {
	// 获取内容
	GetContent(partnerId int)IContent
	
	// 根据编号获取页面
	GetPageById(partnerId,id int)*ValuePage

	// 根据标识获取页面
	GetPageByStringIndent(partnerId int,indent string)*ValuePage

	// 删除页面
	DeletePage(partnerId,id int)error

	// 保存页面
	SavePage(partnerId int,v *ValuePage)(int,error)
}