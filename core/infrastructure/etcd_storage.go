package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types/typeconv"
	clientv3 "go.etcd.io/etcd/client/v3"
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
	cli     *clientv3.Client
	timeout time.Duration
}

// 创建Etcd存储
func NewEtcdStorage(config clientv3.Config) (storage.Interface, error) {
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return &EtcdStorage{
		cli:     cli,
		timeout: 5 * time.Second,
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
	return err == nil && v.Kvs != nil
}

// Set
func (e EtcdStorage) Set(key string, v interface{}) error {
	//j, err := e.serialize(v)
	if v == nil {
		return errors.New("value is nil")
	}
	bytes, err := e.marshal(v)
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
		_, err = e.cli.Put(ctx, key, bytes)
		cancel()
	}
	return err
}

func (e EtcdStorage) SetExpire(key string, v interface{}, seconds int64) error {
	bytes, err := e.marshal(v)
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
		var rsp *clientv3.LeaseGrantResponse
		rsp, err = e.cli.Grant(ctx, seconds)
		if err == nil {
			_, err = e.cli.Put(ctx, key, bytes, clientv3.WithLease(rsp.ID))
		}
		cancel()
	}
	return err
}

func (e EtcdStorage) Get(key string, dst interface{}) error {
	s, err := e.GetBytes(key)
	if err == nil {
		err = e.unmarshal(s, &dst)
	}
	return err
}

func (e EtcdStorage) GetRaw(key string) (interface{}, error) {
	return e.GetBytes(key)
}

func (e EtcdStorage) GetBool(key string) (bool, error) {
	s, err := e.GetString(key)
	if err == nil {
		return strconv.ParseBool(s)
	}
	return false, err
}

func (e EtcdStorage) GetInt(key string) (int, error) {
	s, err := e.GetString(key)
	if err == nil {
		return strconv.Atoi(s)
	}
	return 0, err
}

func (e EtcdStorage) GetInt64(key string) (int64, error) {
	s, err := e.GetString(key)
	if err == nil {
		i, err := strconv.Atoi(s)
		return int64(i), err
	}
	return 0, err
}

func (e EtcdStorage) GetString(key string) (string, error) {
	s, err := e.GetBytes(key)
	if err == nil {
		return string(s), nil
	}
	return "", err
}

func (e EtcdStorage) GetFloat64(key string) (float64, error) {
	s, err := e.GetString(key)
	if err == nil {
		return strconv.ParseFloat(s, 64)
	}
	return 0, err
}

func (e EtcdStorage) GetBytes(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	v, err := e.cli.Get(ctx, key)
	cancel()
	if err == nil {
		if v.Kvs == nil {
			return nil, errors.New("no such key")
		}
		return v.Kvs[0].Value, err
	}
	return []byte(nil), err
}

func (e EtcdStorage) Delete(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	_, _ = e.cli.Delete(ctx, key)
	cancel()
}

func (e EtcdStorage) DeleteWith(prefix string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	ret, _ := e.cli.Get(ctx, prefix, clientv3.WithPrefix())
	var err error
	if ret.Count > 0 {
		for _, v := range ret.Kvs {
			_, err = e.cli.Delete(ctx, string(v.Key))
			if err != nil {
				break
			}
		}
	}
	cancel()
	return int(ret.Count), err
}

func (e EtcdStorage) RWJson(key string, dst interface{}, src func() interface{}, second int64) error {
	jsonBytes, err := e.GetBytes(key)
	if err == nil {
		err = json.Unmarshal(jsonBytes, &dst)
	}
	if err != nil {
		if src == nil {
			panic(errors.New("src is null pointer"))
		}
		dst = src()
		if dst != nil {
			jsonBytes, err = json.Marshal(dst)
			if err == nil {
				if second > 0 {
					_ = e.SetExpire(key, jsonBytes, second)
				} else {
					_ = e.Set(key, jsonBytes)
				}
			}
		}
	}
	return err
}

func (e EtcdStorage) marshal(v interface{}) (string, error) {
	s, b := typeconv.String(v)
	if !b {
		if j, err := json.Marshal(v); err != nil {
			return "", err
		} else {
			s = string(j)
		}
	}
	return s, nil
}

func (e EtcdStorage) unmarshal(s []byte, dst *interface{}) error {
	//err = storage.DecodeBytes(s,&dst)
	return json.Unmarshal(s, &dst)
}
