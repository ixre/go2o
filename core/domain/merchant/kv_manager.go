/**
 * Copyright 2015 @ 56x.net.
 * name : kv_manager
 * author : jarryliu
 * date : 2015-07-26 22:44
 * description :
 * history :
 */
package merchant

import (
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"strconv"
	"time"
)

var _ merchant.IKvManager = new(KvManager)

type KvManager struct {
	mch   *merchantImpl
	mchId int
	// 标识
	indent string
}

func newKvManager(p *merchantImpl, indent string) merchant.IKvManager {
	return &KvManager{
		mch:    p,
		mchId:  int(p.GetAggregateRootId()),
		indent: indent,
	}
}

// 获取键值
func (k *KvManager) Get(key string) string {
	return k.mch._repo.GetKeyValue(k.mchId, k.indent, key)
}

// 获取int类型的键值
func (k *KvManager) GetInt(key string) int {
	i, _ := strconv.Atoi(k.Get(key))
	return i
}

// 设置
func (k *KvManager) Set(key, v string) {
	k.mch._repo.SaveKeyValue(k.mchId, k.indent, key, v, time.Now().Unix())
}

// 获取多项
func (k *KvManager) Gets(key []string) map[string]string {
	return k.mch._repo.GetKeyMap(k.mchId, k.indent, key)
}

// 设置多项
func (k *KvManager) Sets(v map[string]string) error {
	for key, val := range v {
		k.Set(key, val)
	}
	return nil
}

// 根据关键字获取字典
func (k *KvManager) GetsByChar(keyword string) map[string]string {
	return k.mch._repo.GetKeyMapByChar(k.mchId, k.indent, keyword)
}
