package provide

import (
	"github.com/ixre/go2o/core/initial/wrap"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
)

// 应用当前的上下文
var (
	_app gof.App

	_db      db.Connector
	_orm     *wrap.ORM
	_storage storage.Interface
)

func Configure(a gof.App) {
	_app = a
	_db = a.Db()
	_orm = wrap.NewORM(_db)
	_storage = a.Storage()
}

// 获取应用
func GetApp() gof.App {
	return _app
}

// 返回orm实例
func GetDb() db.Connector {
	return _db
}

// 返回orm实例
func GetOrm() *wrap.ORM {
	return _orm
}

// 返回旧orm实例(不包含gorm)
func GetOrmInstance() orm.Orm {
	return _orm.Orm
}

func GetStorageInstance() storage.Interface {
	return _storage
}
