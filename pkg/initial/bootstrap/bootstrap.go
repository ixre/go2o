/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package bootstrap

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	infrastructure "github.com/ixre/go2o/pkg/infra"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	_             gof.App = new(AppConfigLoader)
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

var _AppConfigLoaderInstance = new(AppConfigLoader)

// application context
type AppConfigLoader struct {
	Loaded       bool
	_confFile    string
	_config      *gof.Config
	_dbConnector db.Connector
	_debugMode   bool
	_logger      log.ILogger
	_storage     storage.Interface
	_registry    *gof.Registry
}

func NewApp(confPath string, cfg *clientv3.Config) *AppConfigLoader {
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
	_AppConfigLoaderInstance = &AppConfigLoader{
		_storage:  s,
		_confFile: confPath,
	}
	return _AppConfigLoaderInstance
}

func (a *AppConfigLoader) Db() db.Connector {
	if a._dbConnector == nil {
		a._dbConnector = getDb(a.Config(), a._debugMode, a.Log())
	}
	return a._dbConnector
}

func (a *AppConfigLoader) Storage() storage.Interface {
	// if a._storage == nil {
	// 	//a._storage = storage.NewRedisStorage(a.Redis())
	// }
	return a._storage
}

func (a *AppConfigLoader) Config() *gof.Config {
	if a._config == nil {
		// 优先加载本地开发环境配置
		for _, v := range []string{
			"./local.conf",
			"../local.conf",
			"../../local.conf",
		} {
			cfg, err := gof.LoadConfig(v)
			if err == nil {
				a._config = cfg
				return a._config
			}
		}

		for _, v := range []string{
			a._confFile,
			"../" + a._confFile,
			"../../" + a._confFile,
		} {
			cfg, err := gof.LoadConfig(v)
			if err == nil {
				a._config = cfg
				return a._config
			}
		}
	}
	return a._config
}
func (a *AppConfigLoader) Registry() *gof.Registry {
	if a._registry == nil {
		conf := a.Config().GetString("conf_path")
		if conf == "" {
			conf = "./conf"
		}
		a._registry, _ = gof.NewRegistry(conf, ".")
	}
	return a._registry
}
func (a *AppConfigLoader) Source() interface{} {
	return a
}

func (a *AppConfigLoader) Debug() bool {
	return a._debugMode
}

func (a *AppConfigLoader) Log() log.ILogger {
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
			log.Fatalln("[ GO2O][ ERROR]: database connect failed; error:",
				err.Error(), "; connection string:", connStr)
		} else {
			log.Println("[ GO2O][ INFO]: database connection success")
			conn.SetMaxIdleConns(10000)
			conn.SetMaxIdleConns(5000)
			conn.SetConnMaxLifetime(time.Second * 10)
		}
	}
	return conn
}

// GetRedisPool 获取Redis连接池,will be remove
func GetRedisPool() *redis.Pool {
	if redisPool == nil {
		app := _AppConfigLoaderInstance
		if app == nil {
			panic("gobal serve not initialize!")
		}
		var ok bool
		redisPool, ok = app.Storage().Source().(*redis.Pool)
		if !ok {
			panic("storage drive not base redis")
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
