package clickhouse

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ixre/gof"
	"log"
	"time"
)

var ConnInstance driver.Conn

// Initialize 初始化clickhouse
func Initialize(app gof.App) {
	cfg := app.Config()
	server := []string{cfg.GetString("clickhouse_server")}
	database := cfg.GetString("clickhouse_database")
	password := cfg.GetString("clickhouse_password")
	configure(server, database, password)
}

func configure(server []string, database string, password string) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: server,
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
	ConnInstance = conn
}
