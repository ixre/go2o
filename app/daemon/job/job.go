/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: job.go
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2024-10-12 19:30:58
 * description: 定时任务
 * history:
 */

package job

// JobInfo 任务信息
type JobInfo struct {
	// 任务表达式
	Spec string
	// 任务执行函数
	Cmd func()
}

var jobs []JobInfo

// add 添加任务
func add(spec string, cmd func()) {
	jobs = append(jobs, JobInfo{Spec: spec, Cmd: cmd})
}

func GetJobs() []JobInfo {
	// 检查订单过期,1分钟检测一次
	add("*/1 * * * *", CheckExpiresPaymentOrder)
	// 生成商户日度账单,每天00:00执行
	//add("0 0 0 * *", GenerateMerchantDailyBill)
	add("*/1 * * * *", GenerateMerchantDailyBill)
	return jobs
}
