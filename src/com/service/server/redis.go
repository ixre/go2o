package server

import (
	"com/share/glob"
	"github.com/garyburd/redigo/redis"
	"github.com/newmin/gof/net/jsv"
)

var (
	_redis *redis.Pool
)

func Redis() *redis.Pool {
	if _redis == nil {
		gc := jsv.Context.(*glob.AppContext)
		_redis = gc.Redis
	}
	return _redis
}
