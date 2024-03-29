/**
 * This file is auto generated by tto v0.4.5 !
 * If you want to modify this code, please read the guide
 * to modify code template.
 *
 * Get started: https://github.com/ixre/tto
 *
 * Copyright (C) 2021 <no value>, All rights reserved.
 *
 * name : comm_qr_template_dao.go
 * author : jarrysix
 * date : 2021/12/02 10:37:45
 * description :
 * history :
 */
package dao

import (
	"github.com/ixre/go2o/core/dao/model"
)

type ICommQrTemplateDao interface {
	GetQrTemplates() []*model.QrTemplate
	GetQrTemplate(id int64) *model.QrTemplate
	SaveQrTemplate(q *model.QrTemplate) (int64, error)
	DeleteQrTemplate(id int64) error
}
