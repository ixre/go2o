/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: divide_job.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-09-09 21:14:59
 * description: 分账任务，该任务应在具体的启动项目中注册为定时任务，并传递处理分账的具体实现
 * history:
 */
package job

import (
	"context"
	"time"

	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/infrastructure/locker"
	"github.com/ixre/go2o/core/infrastructure/logger"
	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
)

// SubmitPaymentProviderEndpointDivide提交支付渠道端分账请求
func SubmitPaymentProviderEndpointDivide(f func(o *payment.PayDivide) (string, error)) {
	if f == nil {
		panic("分账处理函数不能为空")
	}
	jobName := "/SubmitDivideJob"
	if !locker.Lock(jobName, 600) {
		return
	}
	defer locker.Unlock(jobName)
	pq := inject.GetPaymentQueryService()
	ps := inject.GetPaymentService()
	size := 50
	unix := time.Now().Unix()
	for {
		orders, err := pq.QueryAwaitSubmitDivides(unix, size)
		if err != nil {
			logger.Error("查询待提交分账记录失败，错误信息:%s", err.Error())
		}
		for _, order := range orders {
			divideNo, err := f(order)
			remark := ""
			if err != nil {
				remark = err.Error()
			}
			ret, _ := ps.UpdateDivideStatus(context.TODO(), &proto.UpdateDivideStatusRequest{
				PayId:    int64(order.PayId),
				DivideId: int64(order.Id),
				Success:  err == nil,
				Remark:   remark,
				DivideNo: divideNo,
			})
			if ret.Code > 0 {
				logger.Error("更新分账记录提交状态失败，id:%d,错误信息:%s", order.Id, ret.Message)
				time.Sleep(time.Second * 5)
			}
		}
		if len(orders) < size {
			break
		}
	}
}
