package handler

import (
	"github.com/ixre/go2o/core/dao"
	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/sys"
)

type EventHandler struct {
	registryRepo  registry.IRegistryRepo
	archiveRepo   content.IArticleRepo
	messageRepo   mss.IMessageRepo
	pageRepo      content.IPageRepo
	portalDao     dao.IPortalDao
	_mchRepo      merchant.IMerchantRepo
	_approvalRepo approval.IApprovalRepository
	_staffRepo    staff.IStaffRepo
	_sysRepo      sys.ISystemRepo
}

func NewEventHandler(repo registry.IRegistryRepo,
	archiveRepo content.IArticleRepo,
	messageRepo mss.IMessageRepo,
	pageRepo content.IPageRepo,
	portalDao dao.IPortalDao,
	mchRepo merchant.IMerchantRepo,
	approvalRepo approval.IApprovalRepository,
	staffRepo staff.IStaffRepo,
	sysRepo sys.ISystemRepo,
) *EventHandler {
	return &EventHandler{
		registryRepo:  repo,
		portalDao:     portalDao,
		messageRepo:   messageRepo,
		pageRepo:      pageRepo,
		archiveRepo:   archiveRepo,
		_mchRepo:      mchRepo,
		_approvalRepo: approvalRepo,
		_staffRepo:    staffRepo,
		_sysRepo:      sysRepo,
	}
}
