package etcd

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : nodes.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-02 17:19
 * description :
 * history :
 */

type Node struct {
	Id   uint32 `json:"id"`
	Addr string `json:"addr"`
}