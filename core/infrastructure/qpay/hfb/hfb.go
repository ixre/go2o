package hfb

import "github.com/ixre/gof/storage"

// 快捷（银行侧)
// http://dev.heepay.com/docs/#/KJJK?id=%e5%bf%ab%e6%8d%b7%ef%bc%88%e9%93%b6%e8%a1%8c%e4%be%a7%ef%bc%89


// 快捷支付
//

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : hfb.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-05 09:27
 * description :
 * history :
 */

var agentId = "0000000"
var md5Key = "CC08C5E3E69F4E6B85F1DC0B"
var sto storage.Interface
// 初始化
func Init(s storage.Interface){
	sto = s
	agentId,_ = s.GetString("qp_hfb_agent_id")
	md5Key,_ = s.GetString("qp_hfb_md5_key")
}

