/**
 * Copyright 2014 @ z3q.net.
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
	"github.com/gomodule/redigo/redis"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
	"time"
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
	DbUsr     = "db_usr"
	DbPwd     = "db_pwd"
	DbCharset = "db_charset"
	//REDIS配置
	RedisHost        = "redis_host"
	RedisDb          = "redis_db"
	RedisPort        = "redis_port"
	RedisAuth        = "redis_auth"
	RedisMaxIdle     = "redis_maxIdle"
	RedisIdleTimeOut = "redis_idleTimeout"
)

// application context
type AppImpl struct {
	Loaded       bool
	_confFile    string
	_config      *gof.Config
	_redis       *redis.Pool
	_dbConnector db.Connector
	_debugMode   bool
	_logger      log.ILogger
	_storage     storage.Interface
	_registry    *gof.Registry
}

func NewApp(confPath string) *AppImpl {
	return &AppImpl{
		_confFile: confPath,
	}
}

func (a *AppImpl) Db() db.Connector {
	if a._dbConnector == nil {
		a._dbConnector = getDb(a.Config(), a._debugMode, a.Log())
		orm.CacheProxy(a._dbConnector.GetOrm(), a.Storage())
	}
	return a._dbConnector
}

func (a *AppImpl) Storage() storage.Interface {
	if a._storage == nil {
		a._storage = storage.NewRedisStorage(a.Redis())
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
		var flag int = 0
		if a._debugMode {
			flag = log.LOpen | log.LESource | log.LStdFlags
		}
		a._logger = log.NewLogger(nil, " O2O", flag)
	}
	return a._logger
}

func (a *AppImpl) Redis() *redis.Pool {
	if a._redis == nil {
		a._redis = CreateRedisPool(a.Config())
	}
	return a._redis
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
	conn := db.NewConnector(driver, connStr, l, debug)
	conn.SetMaxIdleConns(10000)
	conn.SetMaxIdleConns(5000)
	conn.SetConnMaxLifetime(time.Second * 10)
	return conn
}

func CreateRedisPool(c *gof.Config) *redis.Pool {
	host := c.GetString(RedisHost)
	db := c.GetInt(RedisDb)
	port := c.GetInt(RedisPort)
	auth := c.GetString(RedisAuth)
	maxIdle := c.GetInt(RedisMaxIdle)
	idleTimeout := c.GetInt(RedisIdleTimeOut)
	return storage.NewRedisPool(host, port, db, auth, maxIdle, idleTimeout)
}

// 获取Redis连接池
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
	rs := app.Storage().(storage.IRedisStorage)
	if rs != nil {
		i, err := rs.DelWith("go2o:*")
		if err != nil {
			log.Println("[ Go2o][ Redis][ Clean]: Error ", err.Error())
		} else {
			log.Println("[ Go2o][ Redis][ Clean]: clean redis records :", i)
		}
	}
	if CleanHookFunc != nil {
		CleanHookFunc(app)
	}
	log.Println("[ Go2o][ Clean][ Cache]: clean ok !")
}
