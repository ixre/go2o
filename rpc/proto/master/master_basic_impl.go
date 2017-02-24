/**
 * Copyright 2015 @ z3q.net.
 * name : master_basic_impl.go
 * author : jarryliu
 * date : 2016-05-21 17:09
 * description :
 * history :
 */
package admin

import (
	"golang.org/x/net/context"
)

type MasterBasicImpl struct {
}

func (this *MasterBasicImpl) Login(c context.Context, u *Credential) (*Message, error) {
	return nil, nil
}

func (this *MasterBasicImpl) Passwd(c context.Context, u *NewCredential) (*Message, error) {
	return nil, nil
}
