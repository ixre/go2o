package rsi

import (
	"github.com/jsix/gof/util"
	"go2o/core/dao"
	"go2o/core/dao/model"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/thrift/parser"
)

type portalService struct {
	repo *dao.CommonDao
}

func NewPortalService(d *dao.CommonDao) *portalService {
	return &portalService{
		repo: d,
	}
}

// 获取门户导航
func (p *portalService) GetPortalNav_(id int32) *model.PortalNav {
	return p.repo.GetPortalNav(id)
}

// 保存门户导航
func (p *portalService) SavePortalNav_(v *model.PortalNav) (*ttype.Result_, error) {
	return parser.Result_(util.I32Err(p.repo.SavePortalNav(v))), nil
}

// 删除门户导航
func (p *portalService) DeletePortalNav_(id int32) (*ttype.Result_, error) {
	err := p.repo.DeletePortalNav(id)
	return parser.Result_(0, err), nil
}

// 获取门户导航
func (p *portalService) SelectPortalNav(navType int32) []*model.PortalNav {
	return p.repo.SelectPortalNav("nav_type=?", navType)
}

// 获取导航类型
func (p *portalService) GetPortalNavType_(id int32) *model.PortalNavType {
	return p.repo.GetPortalNavType(id)
}

// 保存导航类型
func (p *portalService) SavePortalNavType_(v *model.PortalNavType) (*ttype.Result_, error) {
	return parser.Result_(util.I32Err(p.repo.SavePortalNavType(v))), nil
}

// 删除导航类型
func (p *portalService) DeletePortalNavType_(id int32) (*ttype.Result_, error) {
	err := p.repo.DeletePortalNavType(id)
	return parser.Result_(0, err), nil
}
