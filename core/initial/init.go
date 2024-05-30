package initial

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/db/orm"
)

func Init(ctx gof.App) {
	panic("implement me")
	Context := ctx
	db := Context.Db()
	sto := Context.Storage()
	o := orm.NewOrm(db.Driver(), db.Raw())
	orm.CacheProxy(o, sto)
	// 初始化clickhouse
	//initializeClickhouse(ctx)
	// 初始化服务
	//initService(ctx, db, o, sto)
	// 初始化事件
	//event.InitEvent()
	// 初始化数据
	//InitData(o, nil)
}
