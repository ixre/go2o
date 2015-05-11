/**
 * Copyright 2015 @ S1N1 Team.
 * name : cache.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package cache

import (
	"errors"
	"github.com/atnet/gof"
	"github.com/atnet/gof/storage"
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
	if kvStorage.Driver() == storage.DriveRedisStorage {
		_kvCacheStorage = kvStorage
	} else {
		panic(errors.New("only support redis storage now."))
	}
}
