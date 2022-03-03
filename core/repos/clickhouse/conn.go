package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"log"
	"time"
)

var Conn driver.Conn

func Configure(server []string,database string,password string){
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: server,
		Auth: clickhouse.Auth{
			Database: database,
			//Username: "default",
			Password:password,
		},
		//Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    50,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil{
		log.Fatal(err)
	}
	Conn = conn
}
