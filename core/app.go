/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/ixre/go2o/core/infrastructure"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	_             gof.App = new(AppImpl)
	redisPool     *redis.Pool
	CleanHookFunc func(gof.App) // 当清理数据时候发生
)

const (
	//数据库驱动名称
	DbDriver  = "db_driver"
	DbServer  = "db_server"
	DbPort    = "db_port"
	DbName    = "db_name"
	DbUsr     = "db_user"
	DbPwd     = "db_pwd"
	DbCharset = "db_charset"
)

// application context
type AppImpl struct {
	Loaded       bool
	_confFile    string
	_config      *gof.Config
	_dbConnector db.Connector
	_debugMode   bool
	_logger      log.ILogger
	_storage     storage.Interface
	_registry    *gof.Registry
}

func NewApp(confPath string, cfg *clientv3.Config) *AppImpl {
	fmt.Println(`      
	 ####   ####  #######  ####  
	#    # #    #       # #    # 
	#      #    #  #####  #    # 
	#  ### #    # #       #    # 
	#    # #    # #       #    # 
	 ####   ####  #######  #### 
	`)
	s, err := infrastructure.NewEtcdStorage(*cfg)
	if err != nil {
		panic("[ GO2O][ ERROR]: " + err.Error())
	}
	return &AppImpl{
		_storage:  s,
		_confFile: confPath,
	}
}

func (a *AppImpl) Db() db.Connector {
	if a._dbConnector == nil {
		a._dbConnector = getDb(a.Config(), a._debugMode, a.Log())
	}
	return a._dbConnector
}

func (a *AppImpl) Storage() storage.Interface {
	if a._storage == nil {
		//a._storage = storage.NewRedisStorage(a.Redis())
	}
	return a._storage
}

func (a *AppImpl) Config() *gof.Config {
	if a._config == nil {
		if cfg, err := gof.LoadConfig(a._confFile); err == nil {
			a._config = cfg
		} else {
			log.Fatalln(err)
		}
	}
	return a._config
}
func (a *AppImpl) Registry() *gof.Registry {
	if a._registry == nil {
		conf := a.Config().GetString("conf_path")
		if conf == "" {
			conf = "./conf"
		}
		a._registry, _ = gof.NewRegistry(conf, ".")
	}
	return a._registry
}
func (a *AppImpl) Source() interface{} {
	return a
}

func (a *AppImpl) Debug() bool {
	return a._debugMode
}

func (a *AppImpl) Log() log.ILogger {
	if a._logger == nil {
		var flag = 0
		if a._debugMode {
			flag = log.LOpen | log.LESource | log.LStdFlags
		}
		a._logger = log.NewLogger(nil, " O2O", flag)
	}
	return a._logger
}

func getDb(c *gof.Config, debug bool, l log.ILogger) db.Connector {
	//数据库连接字符串
	//root@tcp(127.0.0.1:3306)/db_name?charset=utf8
	driver := c.GetString(DbDriver)
	dbCharset := c.GetString(DbCharset)
	if dbCharset == "" {
		dbCharset = "utf8"
	}
	//connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=Local",
	//	c.GetString(DbUsr),
	//	c.GetString(DbPwd),
	//	c.GetString(DbServer),
	//	c.GetString(DbPort),
	//	c.GetString(DbName),
	//	dbCharset,
	//)

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.GetString(DbUsr),
		c.GetString(DbPwd),
		c.GetString(DbServer),
		c.GetString(DbPort),
		c.GetString(DbName))
	//todo: charset for connection string?
	conn, err := db.NewConnector(driver, connStr, l, debug)
	if err == nil {
		log.Println("[ GO2O][ INFO]: create database connection..")
		if err := conn.Ping(); err != nil {
			conn.Close()
			//如果异常，则显示并退出
			log.Fatalln("[ GO2O][ ERROR]:" + err.Error())
		}
		conn.SetMaxIdleConns(10000)
		conn.SetMaxIdleConns(5000)
		conn.SetConnMaxLifetime(time.Second * 10)
		return conn
	}
	log.Fatalln("[ GO2O][ ERROR]:" + err.Error())
	return nil
}

// GetRedisPool 获取Redis连接池,will be remove
func GetRedisPool() *redis.Pool {
	if redisPool == nil {
		app := gof.CurrentApp
		if app == nil {
			panic(errors.New("gobal serve not initialize!"))
		}
		var ok bool
		redisPool, ok = app.Storage().Source().(*redis.Pool)
		if !ok {
			panic(errors.New("storage drive not base redis"))
		}
	}
	return redisPool
}

// 获取Redis连接
func GetRedisConn() redis.Conn {
	pool := GetRedisPool()
	if pool != nil {
		return pool.Get()
	}
	return nil
}

// 清除redis缓存
func CleanRedisCache(app gof.App) {
	rs := app.Storage()
	if rs != nil {
		i, err := rs.DeleteWith("go2o:*")
		if err != nil {
			log.Println("[ GO2O][ Redis][ Clean]: Error ", err.Error())
		} else {
			log.Println("[ GO2O][ Redis][ Clean]: clean redis records :", i)
		}
	}
	if CleanHookFunc != nil {
		CleanHookFunc(app)
	}
	log.Println("[ GO2O][ Clean][ Cache]: clean ok !")
}
