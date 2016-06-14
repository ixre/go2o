/**
 * Copyright 2015 @ z3q.net.
 * name : error
 * author : jarryliu
 * date : 2015-07-27 09:22
 * description :
 * history :
 */
package mss

import "go2o/core/infrastructure/domain"

var (
    ErrNotSupportMessageType *domain.DomainError = domain.NewDomainError(
        "err_not_support_message_type", "不支持的消息类型")

    ErrNotEnabled *domain.DomainError = domain.NewDomainError(
        "err_template_not_enabled", "模板未启用")

    ErrTemplateUsed *domain.DomainError = domain.NewDomainError(
        "err_template_used", "模板被使用，无法删除")

    ErrNoSuchNotifyItem *domain.DomainError = domain.NewDomainError(
        "err_no_such_notify_item", "通知项不存在")
)
