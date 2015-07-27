/**
 * Copyright 2015 @ S1N1 Team.
 * name : error
 * author : jarryliu
 * date : 2015-07-27 09:22
 * description :
 * history :
 */
package mss

import "go2o/src/core/infrastructure/domain"

var (
	ErrNotSupportMessageType *domain.DomainError = domain.NewDomainError(
		"err_not_support_message_type", "不支持的消息类型")
	ErrNotEnabled *domain.DomainError = domain.NewDomainError(
		"err_template_not_enabled", "模板未启用")
)
