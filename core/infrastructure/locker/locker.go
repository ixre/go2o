package locker

import (
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util/concurrent"
)

var lock *concurrent.DistributedLock

func Configure(s storage.Interface) {
	lock = concurrent.NewDistributedLock(s)
}

func Lock(key string, expires int64) bool {
	return lock.Lock(key, expires)
}

func Unlock(key string) {
	lock.Unlock(key)
}
