package initial

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
)

// InitTestService 初始化测试服务
func InitTestService(ctx gof.App, db db.Connector, orm orm.Orm, sto storage.Interface) {
	//initService(ctx, db, orm, sto)
	// 初始化clickhouse
	//initializeClickhouse(ctx)
	// 初始化事件
	//event.InitEvent()
}
