package handler

import "github.com/ixre/go2o/core/domain/interface/registry"

type EventHandler struct {
	registryRepo registry.IRegistryRepo
}

func NewEventHandler(repo registry.IRegistryRepo) *EventHandler {
	return &EventHandler{
		registryRepo: repo,
	}
}
