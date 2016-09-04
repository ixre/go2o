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
    "github.com/garyburd/redigo/redis"
    "github.com/jsix/gof"
    "log"
    "time"
    "github.com/jsix/gof/storage"
)

var (
    hasGet bool = false
    globPool *redis.Pool
    CleanHookFunc  func(gof.App) // 当清理数据时候发生
)

func CreateRedisPool(c *gof.Config) *redis.Pool {
    redisHost := c.GetString("redis_host")
    redisDb := c.GetString("redis_db")
    redisPort := c.GetInt("redis_port")
    redisAuth := c.GetString("redis_auth")

    if redisPort <= 0 {
        redisPort = 6379
    }
    redisMaxIdle := c.GetInt("redis_maxIdle")
    if redisMaxIdle <= 0 {
        redisMaxIdle = 10000
    }
    redisIdleTimeout := c.GetInt("redis_idleTimeout")
    if redisIdleTimeout <= 0 {
        redisIdleTimeout = 20000
    }

    return &redis.Pool{
        MaxIdle:     redisMaxIdle,
        IdleTimeout: time.Duration(redisIdleTimeout) * time.Second,
        Dial: func() (redis.Conn, error) {
            dial:
            c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", redisHost, redisPort))
            if err != nil {
                for {
                    log.Printf("[ Redis] - redis(%s:%d) dial failed - %s , Redial after 5 seconds\n",
                        redisHost, redisPort, err.Error())
                    time.Sleep(time.Second * 5)
                    goto dial
                }
            }

            if len(redisAuth) != 0 {
                if _, err := c.Do("AUTH", redisAuth); err != nil {
                    c.Close()
                    log.Fatalf("[ Redis][ AUTH] - %s\n", err.Error())
                }
            }
            if _, err = c.Do("SELECT", redisDb); err != nil {
                c.Close()
                log.Fatalf("[ Redis][ SELECT] - redis(%s:%d) select db failed - %s",
                    redisHost, redisPort, err.Error())
            }

            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }
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
            app.Log().Println("[ Redis][ Clean]: happend error ", err.Error())
        } else {
            app.Log().Println("[ Redis][ Clean]: clean redis records :", i)
        }
    }
    if CleanHookFunc != nil {
        CleanHookFunc(app)
    }
    app.Log().Println("[ Clean][ Cache]: clean ok !")
}
