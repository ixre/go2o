/**
 * Copyright 2015 @ z3q.net.
 * name : cache.go
 * author : jarryliu
 * date : -- :
 * description : 为应用提供缓存
 * history :
 */
package cache

import (
	"errors"
	"github.com/jsix/gof/storage"
)

/** this package is manage system cache. **/

var (
	DefaultMaxSeconds int64 = 300 //默认存储300秒
	kvCacheStorage    storage.Interface
)

// Get Key-value storage
func GetKVS() storage.Interface {
	if kvCacheStorage == nil {
		panic(errors.New("Can't find storage medium."))
	}
	return kvCacheStorage
}

func Initialize(kvStorage storage.Interface) {
	if kvStorage.DriverName() == storage.DriveRedisStorage {
		kvCacheStorage = kvStorage
	} else {
		panic(errors.New("only support redis storage now."))
	}
}
