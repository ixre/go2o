package job

// IJobAggregate 任务聚合
type IJobAggregate interface {
	// GetAggregateRootId 获取编号
	GetAggregateRootId()int64
	// GetValue 获取值
	GetValue()ExecData
	// SetValue 设置值
	SetValue(data ExecData)error
	// AddFail 添加失败计数
	AddFail(recordId int)error
	// UpdateExecCursor 更新执行游标位置
	UpdateExecCursor(id int)error
	// Save 保存
	Save()error
}


type IJobRepo interface{
	CreateJob(*ExecData) IJobAggregate
	// GetExecData Get JobExecData
	GetExecData(primary interface{})*ExecData
	// GetJobByName GetBy JobExecData
	GetJobByName(name string)IJobAggregate
	// SaveExecData Save JobExecData
	SaveExecData(v *ExecData)(int,error)
	// GetExecFailBy GetBy 任务执行失败
	GetExecFailBy(where string,v ...interface{})*ExecFail
	// SaveExecFail Save 任务执行失败
	SaveExecFail(v *ExecFail)(int,error)
}


// ExecData 任务执行数据
type ExecData struct{
	// 编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 任务名称
	JobName string `db:"job_name"`
	// 上次执行位置索引
	LastExecIndex int64 `db:"last_exec_index"`
	// 最后执行时间
	LastExecTime int64 `db:"last_exec_time"`
}


// ExecFail 任务执行失败
type ExecFail struct{
	// 编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 任务编号
	JobId int64 `db:"job_id"`
	// 任务数据编号
	JobDataId int64 `db:"job_data_id"`
	// 重试次数
	RetryCount int `db:"retry_count"`
	// 创建时间
	CreateTime int64 `db:"create_time"`
	// 重试时间
	RetryTime int64 `db:"retry_time"`
}

