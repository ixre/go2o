package glob

import (
	"com/share/variable"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/atnet/gof"
	"github.com/atnet/gof/log"
	"strconv"
	"time"
)

func chkerr(err error) {
	if err != nil {
		log.Fatalln("[Redis][Init Error]:", err.Error())
	}
}
func createRedisPool(c *gof.Config) *redis.Pool {
	redisHost := c.GetString(variable.RedisHost)
	redisDb := c.GetString(variable.RedisDb)
	redisPort, err := strconv.Atoi(c.GetString(variable.RedisPort))
	chkerr(err)
	redisMaxIdle, err := strconv.Atoi(c.GetString(variable.RedisMaxIdle))
	chkerr(err)
	redisIdleTout, err := strconv.Atoi(c.GetString(variable.RedisIdleTout))
	chkerr(err)

	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: time.Duration(redisIdleTout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", redisHost, redisPort))
			if err != nil {
				log.Fatalf("FATAL: redis(%s:%d) initialize failed - %s",
					redisHost, redisPort, err.Error())
			}

			if _, err := c.Do("select", redisDb); err != nil {
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
