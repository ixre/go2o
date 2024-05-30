package provide

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
)

// 应用当前的上下文
var currentApp gof.App
var (
	_db      db.Connector
	_orm     orm.Orm
	_storage storage.Interface
)

func Configure(a gof.App) {
	_db = a.Db()
	_orm = orm.NewOrm(_db.Driver(), _db.Raw())
	_storage = a.Storage()
}

// 获取应用
func GetApp() gof.App {
	return currentApp
}

// 返回orm实例
func GetDb() db.Connector {
	return _db
}

// 返回orm实例
func GetOrmInstance() orm.Orm {
	return _orm
}

func GetStorageInstance() storage.Interface {
	return _storage
}
