package locker

import (
	"github.com/ixre/go2o/core/initial/provide"
	"github.com/ixre/gof/util/concurrent"
)

var lock *concurrent.DistributedLock

func getLocker() *concurrent.DistributedLock {
	if lock == nil {
		lock = concurrent.NewDistributedLock(provide.GetStorageInstance())
	}
	return lock
}

func Lock(key string, expires int64) bool {
	return getLocker().Lock(key, expires)
}

func Unlock(key string) {
	getLocker().Unlock(key)
}
