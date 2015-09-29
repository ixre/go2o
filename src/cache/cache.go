/**
 * Copyright 2015 @ z3q.net.
 * name : cache.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package cache

import (
	"errors"
	"github.com/jsix/gof"
	"github.com/jsix/gof/storage"
)

/** this package is manage system cache. **/

var _kvCacheStorage gof.Storage

// Get Key-value storage
func GetKVS() gof.Storage {
	if _kvCacheStorage == nil {
		panic(errors.New("Can't find storage medium."))
	}
	return _kvCacheStorage
}

func Initialize(kvStorage gof.Storage) {
	if kvStorage.DriverName() == storage.DriveRedisStorage {
		_kvCacheStorage = kvStorage
	} else {
		panic(errors.New("only support redis storage now."))
	}
}
