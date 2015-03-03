/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package server

import (
	"github.com/atnet/gof/net/jsv"
	"github.com/garyburd/redigo/redis"
	"go2o/share/glob"
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
