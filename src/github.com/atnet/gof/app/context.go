package app

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/db"
	"github.com/atnet/gof/log"
	"github.com/atnet/gof/web"
)

type Context interface {
	// Provided db access
	Db() db.Connector

	// Return a Wrapper for golang template.

	Template() *web.TemplateWrapper

	// Return application configs.
	Config() *gof.Config

	// Return a logger
	Log() log.ILogger

	// Get a reference of AppContext
	Source() interface{}

	// Application is running debug mode
	Debug() bool
}
