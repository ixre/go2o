/**
 * Copyright 2015 @ 56x.net.
 * name : 2.express_service.go
 * author : jarryliu
 * date : 2016-07-05 18:57
 * description :
 * history :
 */
package impl

import (
	"context"
	"sort"
	"strings"

	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types"
)

var _ proto.ExpressServiceServer = new(expressServiceImpl)

type expressServiceImpl struct {
	_repo express.IExpressRepo
	serviceUtil
	proto.UnimplementedExpressServiceServer
}

// 获取快递服务
func NewExpressService(rep express.IExpressRepo) *expressServiceImpl {
	return &expressServiceImpl{
		_repo: rep,
	}
}

// 获取快递公司
func (e *expressServiceImpl) GetExpressProvider(_ context.Context, name *proto.IdOrName) (*proto.SExpressProvider, error) {
	var v *express.Provider
	if name.Id > 0 {
		v = e._repo.GetExpressProvider(int32(name.Id))
	} else {
		//v = e._rep.GetExpressProviderByName(name.Name)
	}
	if v != nil {
		return e.parseProviderDto(v), nil
	}
	return nil, express.ErrNotSupportProvider
}

// SaveExpressProvider 保存快递公司
func (e *expressServiceImpl) SaveExpressProvider(_ context.Context, r *proto.SExpressProvider) (*proto.Result, error) {
	v := e.parseProvider(r)
	_, err := e._repo.SaveExpressProvider(v)
	return e.error(err), nil
}

// GetProviders 获取卖家的快递公司
func (e *expressServiceImpl) GetProviders(_ context.Context, _ *proto.Empty) (*proto.ExpressProviderListResponse, error) {
	var arr []*proto.SExpressProvider
	list := e._repo.GetExpressProviders()
	for _, v := range list {
		if v.Enabled == 1 {
			arr = append(arr, e.parseProviderDto(v))
		}
	}
	return &proto.ExpressProviderListResponse{
		Value: arr,
	}, nil
}

// 获取卖家的快递公司分组
func (e *expressServiceImpl) GetProviderGroup(_ context.Context, _ *proto.Empty) (*proto.ExpressProviderGroupResponse, error) {
	list := e._repo.GetExpressProviders()
	for i, v := range list {
		if v.Enabled == 0 {
			list = append(list[:i], list[i+1:]...)
		}
	}
	iMap := make(map[string][]*proto.SMinifiyExpressProvider, 0)
	var letArr []string
	for _, v := range list {
		for _, g := range strings.Split(v.GroupFlag, ",") {
			if g == "" {
				continue
			}
			arr, ok := iMap[g]
			if !ok {
				arr = []*proto.SMinifiyExpressProvider{}
				letArr = append(letArr, g)
			}
			arr = append(arr, &proto.SMinifiyExpressProvider{
				Id:     int64(v.Id),
				Name:   v.Name,
				Letter: v.FirstLetter,
			})
			iMap[g] = arr
		}
	}
	sort.Strings(letArr)
	// 将常用移动到数组开始位置
	l := len(letArr)
	if letArr[l-1] == "常用" {
		letArr = append(letArr[l-1:], letArr[:l-1]...)
	}
	dst := &proto.ExpressProviderGroupResponse{
		List: []*proto.SExpressProviderGroup{},
	}
	for _, k := range letArr {
		dst.List = append(dst.List, &proto.SExpressProviderGroup{
			Group: k,
			List:  iMap[k],
		})
	}
	return dst, nil
}

// 保存快递模板
func (e *expressServiceImpl) SaveExpressTemplate(_ context.Context, r *proto.SExpressTemplate) (*proto.SaveTemplateResponse, error) {
	u := e._repo.GetUserExpress(int(r.SellerId))
	v := e.parseExpressTemplate(r)
	var ie express.IExpressTemplate
	if r.Id > 0 {
		ie = u.GetTemplate(int(r.Id))
	} else {
		ie = u.CreateTemplate(&express.ExpressTemplate{
			VendorId: int(r.SellerId),
		})
	}
	var id int
	err := ie.Set(v)
	if err == nil {
		ie.SetRegionExpress(e.parseRegionsTemplate(r.Regions))
		id, err = ie.Save()
	}
	ret := &proto.SaveTemplateResponse{
		TemplateId: int64(id),
	}
	if err != nil {
		ret.ErrCode = 1
		ret.ErrMsg = err.Error()
	}
	return ret, nil
}

// 获取快递模板
func (e *expressServiceImpl) GetTemplate(_ context.Context, id *proto.ExpressTemplateId) (*proto.SExpressTemplate, error) {
	u := e._repo.GetUserExpress(int(id.SellerId))
	t := u.GetTemplate(int(id.TemplateId))
	if t != nil {
		v := t.Value()
		v2 := t.RegionExpress()
		ret := e.parseExpressTemplateDto(&v)
		ret.Regions = e.parseExpressRegions(&v2)
		return ret, nil
	}
	return nil, express.ErrNoSuchTemplate
}

// 获取所有的快递模板
func (e *expressServiceImpl) GetAllTemplate(userId int32) []*express.ExpressTemplate {
	u := e._repo.GetUserExpress(int(userId))
	list := u.GetAllTemplate()
	arr := make([]*express.ExpressTemplate, len(list))
	for i, v := range list {
		v2 := v.Value()
		arr[i] = &v2
	}
	return arr
}

// 获取可有的快递模板
func (e *expressServiceImpl) GetTemplates(_ context.Context, r *proto.GetTemplatesRequest) (*proto.ExpressTemplateListResponse, error) {
	u := e._repo.GetUserExpress(int(r.SellerId))
	list := u.GetAllTemplate()
	var arr []*proto.SExpressTemplate
	for _, v := range list {
		v2 := v.Value()
		if v2.Enabled == 1 {
			arr = append(arr, e.parseExpressTemplateDto(&v2))
		}
	}
	return &proto.ExpressTemplateListResponse{
		Value: arr,
	}, nil
}

// 删除模板
func (e *expressServiceImpl) DeleteTemplate(_ context.Context, id *proto.ExpressTemplateId) (*proto.Result, error) {
	u := e._repo.GetUserExpress(int(id.SellerId))
	err := u.DeleteTemplate(int(id.TemplateId))
	return e.error(err), nil
}

func (e *expressServiceImpl) parseProviderDto(v *express.Provider) *proto.SExpressProvider {
	return &proto.SExpressProvider{
		Id:        int64(v.Id),
		Name:      v.Name,
		Letter:    v.FirstLetter,
		GroupFlag: v.GroupFlag,
		Code:      v.Code,
		ApiCode:   v.ApiCode,
		Enabled:   v.Enabled == 1,
	}
}

func (e *expressServiceImpl) parseProvider(r *proto.SExpressProvider) *express.Provider {
	return &express.Provider{
		Id:          int32(r.Id),
		Name:        r.Name,
		FirstLetter: r.Letter,
		GroupFlag:   r.GroupFlag,
		Code:        r.Code,
		ApiCode:     r.ApiCode,
		Enabled:     types.ElseInt(r.Enabled, 1, 0),
	}
}

func (e *expressServiceImpl) parseExpressTemplate(r *proto.SExpressTemplate) *express.ExpressTemplate {
	return &express.ExpressTemplate{
		Id:        int(r.Id),
		VendorId:  int(r.SellerId),
		Name:      r.Name,
		IsFree:    types.ElseInt(r.IsFree, 1, 0),
		Basis:     int(r.Basis),
		FirstUnit: int(r.FirstUnit),
		FirstFee:  r.FirstFee,
		AddUnit:   int(r.AddUnit),
		AddFee:    r.AddFee,
		Enabled:   types.ElseInt(r.Enabled, 1, 0),
	}
}

func (e *expressServiceImpl) parseExpressTemplateDto(v *express.ExpressTemplate) *proto.SExpressTemplate {
	return &proto.SExpressTemplate{
		Id:        int64(v.Id),
		SellerId:  int64(v.VendorId),
		Name:      v.Name,
		IsFree:    v.IsFree == 1,
		Basis:     int32(v.Basis),
		FirstUnit: int32(v.FirstUnit),
		FirstFee:  v.FirstFee,
		AddUnit:   int32(v.AddUnit),
		AddFee:    v.AddFee,
		Enabled:   v.Enabled == 1,
	}
}

func (e *expressServiceImpl) parseExpressRegions(regions *[]express.RegionExpressTemplate) []*proto.SRegionExpressTemplate {
	arr := make([]*proto.SRegionExpressTemplate, 0)
	for _, v := range *regions {
		arr = append(arr, &proto.SRegionExpressTemplate{
			Id:        int64(v.Id),
			CodeList:  v.CodeList,
			NameList:  v.NameList,
			FirstUnit: v.FirstUnit,
			FirstFee:  v.FirstFee,
			AddUnit:   v.AddUnit,
			AddFee:    v.AddFee,
		})
	}
	return arr
}

func (e *expressServiceImpl) parseRegionsTemplate(regions []*proto.SRegionExpressTemplate) *[]express.RegionExpressTemplate {
	arr := make([]express.RegionExpressTemplate, 0)
	for _, v := range regions {
		arr = append(arr, express.RegionExpressTemplate{
			Id:        int(v.Id),
			CodeList:  v.CodeList,
			NameList:  v.NameList,
			FirstUnit: v.FirstUnit,
			FirstFee:  v.FirstFee,
			AddUnit:   v.AddUnit,
			AddFee:    v.AddFee,
		})
	}
	return &arr
}
