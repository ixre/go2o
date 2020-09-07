package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/ixre/gof/storage"
	"go.etcd.io/etcd/clientv3"
	"strconv"
	"time"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : etcd_storage.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-07 10:59
 * description :
 * history :
 */

var _ storage.Interface = new(EtcdStorage)
var ctx = context.TODO()
type EtcdStorage struct {
	cli        *clientv3.Client
	timeout time.Duration
}

// 创建Etcd存储
func NewEtcdStorage(config clientv3.Config) (storage.Interface, error) {
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return &EtcdStorage{
		cli: cli,
		timeout: 10*time.Second,
	}, nil
}

func (e EtcdStorage) Driver() string {
	return storage.DriveEtcdStorage
}

func (e EtcdStorage) Source() interface{} {
	return e.cli
}

func (e EtcdStorage) Exists(key string) (exists bool) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	v, err := e.cli.Get(ctx, key)
	cancel()
	return err == nil && v != nil
}

func (e EtcdStorage) Set(key string, v interface{}) error {
	j,err := json.Marshal(v)
	if err == nil{
		ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
		_,err = e.cli.Put(ctx,key,string(j))
		cancel()
	}
	return err
}

func (e EtcdStorage) SetExpire(key string, v interface{}, seconds int64) error {
	panic("implement me")
}

func (e EtcdStorage) Get(key string, dst interface{}) error {
	s,err := e.GetBytes(key)
	if err == nil{
		err = json.Unmarshal(s,&dst)
	}
	return err
}

func (e EtcdStorage) GetRaw(key string) (interface{}, error) {
	return e.GetBytes(key)
}

func (e EtcdStorage) GetBool(key string) (bool, error) {
	panic("implement me")
}

func (e EtcdStorage) GetInt(key string) (int, error) {
	s,err := e.GetString(key)
	if err == nil {
		return strconv.Atoi(s)
	}
	return 0,err
}

func (e EtcdStorage) GetInt64(key string) (int64, error) {
	panic("implement me")
}

func (e EtcdStorage) GetString(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	v, err := e.cli.Get(ctx, key)
	cancel()
	if err == nil && len(v.Kvs) > 0{
		return string(v.Kvs[0].Value), err
	}
	return "",err
}

func (e EtcdStorage) GetFloat64(key string) (float64, error) {
	s,err := e.GetString(key)
	if err == nil {
		return strconv.ParseFloat(s,64)
	}
	return 0,err
}

func (e EtcdStorage) GetBytes(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	v, err := e.cli.Get(ctx, key)
	cancel()
	if err == nil&& len(v.Kvs) > 0 {
		return v.Kvs[0].Value, err
	}
	return []byte(""),err
}

func (e EtcdStorage) Del(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	_, _ = e.cli.Delete(ctx, key)
	cancel()
}

func (e EtcdStorage) RWJson(key string, dst interface{}, src func() interface{}, second int64) error {
	panic("implement me")
}
