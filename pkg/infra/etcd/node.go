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

// 节点
type Node struct {
	// Id 编号
	Id uint32 `json:"id"`
	// Addr 地址
	Addr string `json:"addr"`
}
