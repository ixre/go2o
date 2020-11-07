package qpay

import "github.com/ixre/gof/storage"

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : cache
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-07 10:38
 * description :
 * history :
 */

type Cache struct{
	s storage.Interface
}

func NewCache(s storage.Interface)*Cache{
	return &Cache{s:s}
}
func (c *Cache)  BankAuthNonceKey(nonce string) string {
	return "pay/quick/bank/auth/"+nonce
}

func (c *Cache) GetBankAuthData(nonce string)*BankAuthSwapData {
	key := c.BankAuthNonceKey(nonce)
	var ret *BankAuthSwapData
	c.s.Get(key,&ret)
	return ret
}

func (c *Cache) SaveBankAuthData(nonce string, data *BankAuthSwapData, expires int) {
	key := c.BankAuthNonceKey(nonce)
	c.s.SetExpire(key,data,int64(expires))
}