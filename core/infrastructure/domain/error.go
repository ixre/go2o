/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package domain

import (
	"fmt"

	"github.com/ixre/go2o/core/initial/provide"
)

// todo: 可以做通过后台设置错误信息
// 处理错误
func HandleError(err error, src string) error {
	debug := provide.GetApp().Debug()

	if err != nil && debug {
		logger := provide.GetApp().Log()
		logger.Println("[ GO2O][ ERROR] - ", err.Error())
	}
	return err
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
