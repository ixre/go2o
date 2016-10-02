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
    "github.com/garyburd/redigo/redis"
    "github.com/jsix/gof"
    "github.com/jsix/gof/storage"
    "log"
)

var (
    hasGet bool = false
    globPool      *redis.Pool
    CleanHookFunc func(gof.App) // 当清理数据时候发生
)

func CreateRedisPool(c *gof.Config) *redis.Pool {
    host := c.GetString("redis_host")
    db := c.GetInt("redis_db")
    port := c.GetInt("redis_port")
    auth := c.GetString("redis_auth")

    if port <= 0 {
        port = 6379
    }
    maxIdle := c.GetInt("redis_maxIdle")
    if maxIdle <= 0 {
        maxIdle = 10000
    }
    idleTimeout := c.GetInt("redis_idleTimeout")
    if idleTimeout <= 0 {
        idleTimeout = 20000
    }
    return storage.NewRedisPool(host,port,db,auth,maxIdle,idleTimeout)
}

// 获取Redis连接池
func GetRedisPool() *redis.Pool {
    if !hasGet {
        app := gof.CurrentApp
        if app == nil {
            panic(errors.New("gobal app not initialize!"))
        }
        var ok bool
        globPool, ok = app.Storage().Driver().(*redis.Pool)
        if !ok {
            panic(errors.New("storage drive not base redis"))
        }
        hasGet = true
    }
    return globPool
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
        i, err := rs.PrefixDel("go2o:*")
        if err != nil {
            log.Println("[ Go2o][ Redis][ Clean]: happend error ", err.Error())
        } else {
            log.Println("[ Go2o][ Redis][ Clean]: clean redis records :", i)
        }
    }
    if CleanHookFunc != nil {
        CleanHookFunc(app)
    }
    log.Println("[ Go2o][ Clean][ Cache]: clean ok !")
}

// 删除指定前缀的缓存
func RemovePrefixKeys(sto storage.Interface, prefix string) {
    rds := sto.(storage.IRedisStorage)
    _, err := rds.PrefixDel(prefix)
    if err != nil {
        log.Println("[ Cache][ Clean]: clean by prefix ", prefix, " error:", err)
    }
}
