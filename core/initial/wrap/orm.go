package wrap

import (
	"log"

	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ORM 数据库操作接口
type ORM struct {
	// 引用orm包
	orm.Orm
	// 引用gorm包
	DB *gorm.DB
}

func NewORM(db db.Connector) *ORM {
	o := orm.NewOrm(db.Driver(), db.Raw())
	n, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db.Raw(),
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("gorm.Open: %v", err)
	}
	return &ORM{
		Orm: o,
		DB:  n,
	}
}
