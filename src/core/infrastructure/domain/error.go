/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package domain

type DomainError struct {
	Key          string
	DefaultError string
}

func NewDomainError(key string, msg string) *DomainError {
	return &DomainError{
		Key:          key,
		DefaultError: msg,
	}
}

func (this *DomainError) Error() string {
	return this.DefaultError
}
