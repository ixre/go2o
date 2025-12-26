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
	"context"
	"time"

	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/infrastructure/locker"
	"github.com/ixre/go2o/core/infrastructure/logger"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/proto"
)

func GenerateMerchantBill() {
	jobName := "merchant.generateMerchantBill"
	if !locker.Lock(jobName, 300) {
		return
	}
	defer locker.Unlock(jobName)
	ms := inject.GetMerchantService()
	qs := inject.GetMerchantQueryService()
	// 生成日度账单
	generateMerchantDailyBill(ms, qs)
	// 生成月度账单
	generateMerchantMonthlyBill(ms, qs)
}

// generateMerchantDailyBill 生成商户日度账单
func generateMerchantDailyBill(ms proto.MerchantServiceServer, qs *query.MerchantQuery) {
	size, lastId := 5, 0
	for {
		list := qs.QueryWaitGenerateDailyBills(size, lastId)
		for _, bill := range list {
			ret, _ := ms.GenerateBill(context.TODO(), &proto.GenerateMerchantBillRequest{
				MchId:    int64(bill.MchId),
				BillId:   int64(bill.Id),
				BillType: int32(merchant.BillTypeDaily),
			})
			if ret.Code != 0 {
				logger.Error("生成商户日度账单失败: %v, 商户ID:%d, 账单编号:%d", ret.Message, bill.MchId, bill.Id)
			}
			lastId = bill.Id
		}
		if len(list) < size {
			break
		}
	}
}

// generateMerchantMonthlyBill 生成商户月度账单
func generateMerchantMonthlyBill(ms proto.MerchantServiceServer, qs *query.MerchantQuery) {
	day := time.Now().Day()
	if day < 3 {
		// 每月3日前不生成上个月账单
		return
	}
	// 生成上个月账单
	dt := time.Now().AddDate(0, -1, 0)
	billTime := time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, time.Local).Unix()
	year, month := dt.Year(), dt.Month()
	begin, size := 0, 100
	mchList := make([]int64, 0)
	for {
		list := qs.QueryWaitGenerateMonthBillMerchants(billTime, begin, size)
		mchList = append(mchList, list...)
		l := len(list)
		if l < size {
			break
		}
		begin += l
	}

	for _, mchId := range mchList {
		ret, _ := ms.GenerateBill(context.TODO(), &proto.GenerateMerchantBillRequest{
			MchId:    mchId,
			BillType: int32(merchant.BillTypeMonthly),
			Unixtime: billTime,
		})
		if ret.Code != 0 {
			logger.Error("生成商户月度账单失败: %v, 商户ID:%d, 账单日期:%d-%d", ret.Message, mchId, year, month)
		}
	}
	logger.Info("[ Job]: 生成商户月度(%d-%d)账单完成, 共%d个商户", year, month, len(mchList))
}
