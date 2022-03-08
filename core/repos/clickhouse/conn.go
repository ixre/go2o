package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"log"
	"time"
)

var connInstance driver.Conn

// IsCluster 是否为集群模式
var IsCluster bool

// GetClickhouseConn 获取Clickhouse写入连接
func GetClickhouseConn() driver.Conn {
	return connInstance
}

// Configure 配置clickhouse写入连接
func Configure(servers []string, database string, password string) {
	if len(servers) == 0 || servers[0] == ""{
		return
	}
	log.Println("[ Go2o][ Info]: configure clickhouse connection..")
	IsCluster = len(servers) > 1
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: servers,
		Auth: clickhouse.Auth{
			Database: database,
			Username: "default",
			Password: password,
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
	if err != nil {
		log.Fatal(err)
	}
	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 10,
	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
		fmt.Println("progress: ", p)
	}))
	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		log.Fatal(err)
	}
	connInstance = conn
}

var connDb *sql.DB

func GetClickhouseDB() *sql.DB {
	return connDb
}

// InitializeDB 初始化clickhouse查询连接
func InitializeDB(servers []string, database string, password string) {
	if len(servers) == 0 || servers[0] == ""{
		return
	}
	log.Println("[ Go2o][ Info]: configure clickhouse sql connection..")
	IsCluster = len(servers) > 1
	// 初始化连接
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: servers,
		Auth: clickhouse.Auth{
			Database: database,
			Username: "default",
			Password: password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug: false,
	})
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(50)
	conn.SetConnMaxLifetime(time.Hour)
	if err := conn.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		log.Fatal(err)
	}
	connDb = conn
}
