package impl

import (
	"go2o/core/dao"
	"go2o/core/dao/impl"
	"go2o/core/dao/model"
	"go2o/core/service/proto"
	 "golang.org/x/net/context"
)

var _ proto.PortalServiceServer = new(portalService)

type portalService struct {
	repo *impl.CommonDao
	dao  dao.IPortalDao
	serviceUtil
}

func NewPortalService(d *impl.CommonDao, portalDao dao.IPortalDao) *portalService {
	return &portalService{
		repo: d,
		dao:  portalDao,
	}
}

func (p *portalService) SaveNav(context context.Context, r *proto.SaveNavRequest) (*proto.SaveNavResponse, error) {
	var dst *model.PortalNav
	if r.Id > 0 {
		if dst = p.dao.GetNav(r.Id); dst == nil {
			return &proto.SaveNavResponse{
				ErrCode: 2,
				ErrMsg:  "no such record",
			}, nil
		}
	} else {
		dst = &model.PortalNav{}

	}
	dst.Text = r.Text
	dst.Url = r.Url
	dst.Target = r.Target
	dst.Image = r.Image
	dst.NavType = r.NavType
	id, err := p.dao.SaveNav(dst)
	ret := &proto.SaveNavResponse{
		Id: int64(id),
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

func (p *portalService) GetNav(context context.Context, id *proto.PortalNavId) (*proto.SNav, error) {
	v := p.dao.GetNav(id.Value)
	if v == nil {
		return nil, nil
	}
	return p.parseNav(v), nil
}

func (p *portalService) QueryNavList(context context.Context, r *proto.QueryNavRequest) (*proto.QueryNavResponse, error) {
	arr := p.dao.SelectNav("nav_type= $1", r.NavType)
	ret := &proto.QueryNavResponse{
		List: make([]*proto.SNav, len(arr)),
	}
	for i, v := range arr {
		ret.List[i] = p.parseNav(v)
	}
	return ret, nil
}

func (p *portalService) DeleteNav(context context.Context, id *proto.PortalNavId) (*proto.Result, error) {
	err := p.dao.DeleteNav(id.Value)
	return p.error(err), nil
}

// 获取门户导航
func (p *portalService) GetPortalNav_(id int32) *model.PortalNav {
	return p.repo.GetPortalNav(id)
}

// 保存门户导航
func (p *portalService) SavePortalNav_(v *model.PortalNav) (*proto.Result, error) {
	_, err := p.repo.SavePortalNav(v)
	return p.result(err), nil
}

// 删除门户导航
func (p *portalService) DeletePortalNav_(id int32) (*proto.Result, error) {
	err := p.repo.DeletePortalNav(id)
	return p.result(err), nil
}

// 获取门户导航
func (p *portalService) SelectPortalNav(navType int32) []*model.PortalNav {
	return p.repo.SelectPortalNav("nav_type= $1", navType)
}

// 获取导航类型
func (p *portalService) GetPortalNavType_(id int32) *model.PortalNavType {
	return p.repo.GetPortalNavType(id)
}

// 保存导航类型
func (p *portalService) SavePortalNavType_(v *model.PortalNavType) (*proto.Result, error) {
	_, err := p.repo.SavePortalNavType(v)
	return p.result(err), nil
}

// 删除导航类型
func (p *portalService) DeletePortalNavType_(id int32) (*proto.Result, error) {
	err := p.repo.DeletePortalNavType(id)
	return p.result(err), nil
}

func (p *portalService) parseNav(v *model.PortalNav) *proto.SNav {
	return &proto.SNav{
		Id:      v.Id,
		Text:    v.Text,
		Url:     v.Url,
		Target:  v.Target,
		Image:   v.Image,
		NavType: v.NavType,
	}
}
