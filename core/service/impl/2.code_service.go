package impl

import (
	"context"
	"github.com/ixre/go2o/core/dao"
	"github.com/ixre/go2o/core/dao/impl"
	"github.com/ixre/go2o/core/dao/model"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
)

var _ proto.CodeServiceServer = new(codeServiceImpl)

type codeServiceImpl struct {
	dao dao.ICommQrTemplateDao
	s   storage.Interface
	serviceUtil
}

func NewCodeService(sto storage.Interface, o orm.Orm) *codeServiceImpl {
	return &codeServiceImpl{
		dao: impl.NewCommQrTemplateDao(o, sto),
		s:   sto,
	}
}

// SaveQrTemplate 保存CommQrTemplate
func (c *codeServiceImpl) SaveQrTemplate(_ context.Context, r *proto.SaveQrTemplateRequest) (*proto.SaveQrTemplateResponse, error) {
	var dst *model.QrTemplate
	if r.Id > 0 {
		if dst = c.dao.GetQrTemplate(r.Id); dst == nil {
			return &proto.SaveQrTemplateResponse{
				ErrCode: 2,
				ErrMsg:  "no such record",
			}, nil
		}
	} else {
		dst = &model.QrTemplate{}

	}
	dst.Title = r.Title
	dst.BgImage = r.BgImage
	dst.OffsetX = int(r.OffsetX)
	dst.OffsetY = int(r.OffsetY)
	dst.Comment = r.Comment
	dst.CallbackUrl = r.CallbackUrl
	dst.Enabled = int(r.Enabled)

	id, err := c.dao.SaveQrTemplate(dst)
	ret := &proto.SaveQrTemplateResponse{
		Id: id,
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

// 获取CommQrTemplate
func (c *codeServiceImpl) GetQrTemplate(_ context.Context, id *proto.CommQrTemplateId) (*proto.SQrTemplate, error) {
	v := c.dao.GetQrTemplate(id.Value)
	if v == nil {
		return nil, nil
	}
	return c.parseQrTemplate(v), nil
}

func (c *codeServiceImpl) DeleteQrTemplate(_ context.Context, id *proto.CommQrTemplateId) (*proto.Result, error) {
	err := c.dao.DeleteQrTemplate(id.Value)
	return c.error(err), nil
}

func (c *codeServiceImpl) parseQrTemplate(v *model.QrTemplate) *proto.SQrTemplate {
	return &proto.SQrTemplate{
		Id:          int64(v.Id),
		Title:       v.Title,
		BgImage:     v.BgImage,
		OffsetX:     int32(v.OffsetX),
		OffsetY:     int32(v.OffsetY),
		Comment:     v.Comment,
		CallbackUrl: v.CallbackUrl,
		Enabled:     int32(v.Enabled),
	}
}
