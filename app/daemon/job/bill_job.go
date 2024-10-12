/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: bill_job.go
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2024-10-12 19:30:58
 * description: 账单生成任务
 * history:
 */

package job

import (
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/infrastructure/logger"
	"github.com/ixre/go2o/core/inject"
)

// GenerateMerchantDailyBill 生成商户日度账单
func GenerateMerchantDailyBill() {
	qs := inject.GetMerchantQueryService()
	size := 5
	lastId := 0
	for {
		list := qs.QueryWaitGenerateBills(size, lastId)
		for _, bill := range list {
			generateMerchantBill(bill)
			lastId = bill.Id
		}
		if len(list) < size {
			break
		}
	}
}

// generateMerchantBill 生成商户账单
func generateMerchantBill(bill *merchant.MerchantBill) {
	repo := inject.GetMerchantRepo()
	mch := repo.GetMerchant(bill.MchId)
	if mch == nil {
		logger.Error("生成账单失败: 商户不存在,商户ID:%d, 账单编号:%d", bill.MchId, bill.Id)
		return
	}
	err := mch.SaleManager().GenerateBill(bill.Id)
	if err != nil {
		logger.Error("生成账单失败: %v, 商户ID:%d, 账单编号:%d", err, bill.MchId, bill.Id)
	}
}
