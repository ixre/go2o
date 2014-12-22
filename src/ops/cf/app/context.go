package app

import (
	"ops/cf"
	"ops/cf/db"
	"ops/cf/log"
	"ops/cf/web"
)

type Context interface {
	// Provided db access
	Db() db.Connector

	// Return a Wrapper for golang template.

	Template() *web.TemplateWrapper

	// Return application configs.
	Config() *cf.Config

	// Return a logger
	Log() log.ILogger

	// Get a reference of AppContext
	Source() interface{}

	// Application is running debug mode
	Debug() bool
}
