package etcd

import (
	"context"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var Locker *DistributeLock

// Etcd 分布式锁
type DistributeLock struct {
	cli *clientv3.Client
}

// 配置分布式锁
func InitializeLocker(config *clientv3.Config) error {
	l, err := NewDistributeLock(*config)
	Locker = l
	if err != nil {
		log.Println("[ GO2O][ ERROR]: Failed to initialize ETCD Locker: ", err)
	}
	return err
}

// NewDistributeLock creates a new DistributeLock
func NewDistributeLock(config clientv3.Config) (*DistributeLock, error) {
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return &DistributeLock{
		cli: cli,
	}, nil
}

// Lock 尝试获得锁(阻塞)
func (d *DistributeLock) Lock(prefix string) (cancel func(), err1 error) {
	m, err := concurrency.NewSession(d.cli, concurrency.WithTTL(5))
	if err == nil {
		locker := concurrency.NewLocker(m, prefix)
		locker.Lock()
		return func() {
			locker.Unlock()
		}, err
	}
	return nil, err
}

// TryLock 尝试获得锁(非阻塞)
func (d *DistributeLock) TryLock(prefix string) (cancel func(), err1 error) {
	m, err := concurrency.NewSession(d.cli, concurrency.WithTTL(5))
	if err == nil {
		locker := concurrency.NewMutex(m, prefix)
		err = locker.TryLock(context.Background())
		return func() {
			locker.Unlock(context.Background())
		}, err
	}
	return nil, err
}
