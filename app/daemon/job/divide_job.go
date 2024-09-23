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

// SubmitPaymentProviderEndpointDivide 提交支付渠道端分账请求
// 如果实时分账则订阅事件，对分账进行处理后，定时任务中将不会再检测待提交分账的任务
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
		divides, err := pq.QueryAwaitSubmitDivides(unix, size)
		if err != nil {
			logger.Error("查询待提交分账记录失败，错误信息:%s", err.Error())
		}
		for _, dv := range divides {
			divideNo, err := f(dv)
			remark := "成功"
			if err != nil {
				logger.Error("提交分账记录失败，分账ID:%d, 错误信息:%s", dv.Id, err.Error())
				remark = err.Error()
			}
			ret, _ := ps.UpdateDivideStatus(context.TODO(), &proto.UpdateDivideStatusRequest{
				PayId:    int64(dv.PayId),
				DivideId: int64(dv.Id),
				Success:  err == nil,
				Remark:   remark,
				DivideNo: divideNo,
			})
			if ret.Code > 0 {
				logger.Error("更新分账记录提交状态失败，id:%d,错误信息:%s", dv.Id, ret.Message)
				time.Sleep(time.Second * 5)
			}
		}
		if len(divides) < size {
			break
		}
	}
}
