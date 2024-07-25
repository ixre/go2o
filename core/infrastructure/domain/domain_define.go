/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:52
 * description :
 * history :
 */

package domain

import "fmt"

// 聚合根
type IAggregateRoot interface {
	// 获取聚合根编号
	GetAggregateRootId() int
}

// 领域对象
type IDomain interface {
	// 获取领域对象编号
	GetDomainId() int
}

// 值对象
type IValueObject interface {
	Equal(interface{}) bool
}

var _ error = &DomainError{}

type DomainError struct {
	Key          string
	DefaultError string
}

func NewError(key string, msg string) *DomainError {
	return &DomainError{
		Key:          key,
		DefaultError: msg,
	}
}

func (d *DomainError) Error() string {
	return d.DefaultError
}

func (d *DomainError) Set(msg string) {
	d.DefaultError = msg
}

// Apply 格式化错误信息
func (d *DomainError) Apply(args ...interface{}) *DomainError {
	return &DomainError{
		Key:          d.Key,
		DefaultError: fmt.Sprintf(d.DefaultError, args...),
	}
}
