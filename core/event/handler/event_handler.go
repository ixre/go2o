package handler

import (
	"github.com/ixre/go2o/core/dao"
	"github.com/ixre/go2o/core/domain/interface/content"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/station"
)

type EventHandler struct {
	registryRepo registry.IRegistryRepo
	archiveRepo  content.IArchiveRepo
	stationRepo  station.IStationRepo
	portalDao    dao.IPortalDao
}

func NewEventHandler(repo registry.IRegistryRepo,
	archiveRepo content.IArchiveRepo,
	stationRepo station.IStationRepo,
	portalDao dao.IPortalDao,
) *EventHandler {
	return &EventHandler{
		registryRepo: repo,
		portalDao:    portalDao,
		archiveRepo:  archiveRepo,
		stationRepo:  stationRepo,
	}
}
