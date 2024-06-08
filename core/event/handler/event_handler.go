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
	archiveRepo  content.IArchiveRepo
	stationRepo  station.IStationRepo
	messageRepo  mss.IMessageRepo
	portalDao    dao.IPortalDao
}

func NewEventHandler(repo registry.IRegistryRepo,
	archiveRepo content.IArchiveRepo,
	messageRepo mss.IMessageRepo,
	stationRepo station.IStationRepo,
	portalDao dao.IPortalDao,
) *EventHandler {
	return &EventHandler{
		registryRepo: repo,
		portalDao:    portalDao,
		messageRepo:  messageRepo,
		archiveRepo:  archiveRepo,
		stationRepo:  stationRepo,
	}
}
