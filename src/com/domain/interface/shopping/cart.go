/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-08 10:23
 * description :
 * history :
 */

package shopping

type ICart interface {
	GetDomainId() int
	GetValue() ValueCart
	GetSummary() string
	// 获取金额
	GetFee() (totalFee float32, orderFee float32)
}
