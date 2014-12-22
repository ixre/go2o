package entity

import (
	"time"
)

//分类
type Category struct {
	Id int
	//父分类
	Pid int
	//供应商编号
	Ptid int
	//名称
	Name       string
	Idx        int
	CreateTime time.Time
	IsEnabled  bool
	Descript   string
}
