/**
 * Copyright 2015 @ to2.net.
 * name : cache.go
 * author : jarryliu
 * date : -- :
 * description : 为应用提供缓存
 * history :
 */
package cache

import (
	"errors"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
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

// 删除指定前缀的缓存
func PrefixDel(prefix string) {
	sto := GetKVS().(storage.IRedisStorage)
	_, err := sto.DelWith(prefix)
	if err != nil {
		log.Println("[ Cache][ Clean]: clean by prefix ", prefix, " error:", err)
	}
}

func Initialize(kvStorage storage.Interface) {
	if kvStorage.Driver() == storage.DriveRedisStorage {
		kvCacheStorage = kvStorage
	} else {
		panic(errors.New("only support redis storage now."))
	}
}
