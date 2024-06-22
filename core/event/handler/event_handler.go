package handler

import (
	"github.com/ixre/go2o/core/dao"
	"github.com/ixre/go2o/core/domain/interface/content"
	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/station"
)

type EventHandler struct {
	registryRepo registry.IRegistryRepo
	archiveRepo  content.IArticleRepo
	stationRepo  station.IStationRepo
	messageRepo  mss.IMessageRepo
	pageRepo     content.IPageRepo
	portalDao    dao.IPortalDao
}

func NewEventHandler(repo registry.IRegistryRepo,
	archiveRepo content.IArticleRepo,
	messageRepo mss.IMessageRepo,
	stationRepo station.IStationRepo,
	pageRepo content.IPageRepo,
	portalDao dao.IPortalDao,
) *EventHandler {
	return &EventHandler{
		registryRepo: repo,
		portalDao:    portalDao,
		messageRepo:  messageRepo,
		pageRepo:     pageRepo,
		archiveRepo:  archiveRepo,
		stationRepo:  stationRepo,
	}
}
