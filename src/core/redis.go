/**
 * Copyright 2014 @ s1n1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package core

import (
	"fmt"
	"github.com/atnet/gof"
	"log"
	"github.com/garyburd/redigo/redis"
	"time"
)

func createRedisPool(c *gof.Config) *redis.Pool {
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
				if _,err := c.Do("AUTH",redisAuth);err != nil{
					c.Close()
					log.Fatalf("[ Redis][ AUTH] - %s\n",err.Error())
				}
			}
			if _, err = c.Do("SELECT", redisDb); err != nil {
				c.Close()
				log.Fatalf("FATAL: redis(%s:%d) initialize failed - %s",
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
