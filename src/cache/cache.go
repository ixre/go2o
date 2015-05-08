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
    "github.com/atnet/gof"
    "errors"
)

/** this package is manage system cache. **/

var _kvCacheStorage gof.Storage

// Get Key-value storage
func GetKVS()gof.Storage{
    if _kvCacheStorage == nil{
        panic(errors.New("Can't find storage medium."))
    }
    return _kvCacheStorage
}

func Initialize(kvStorage gof.Storage){
    _kvCacheStorage = kvStorage
}