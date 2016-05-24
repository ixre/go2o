/**
 * Copyright 2015 @ z3q.net.
 * name : merchant.go
 * author : jarryliu
 * date : 2016-05-21 13:04
 * description :
 * history :
 */
package impl

type MasterCredentialRpcImpl struct {
}

func (this *MasterCredentialRpcImpl) Login(args *Args,
	relay *map[string]interface{}) (err error) {
	checkArgs(args)
}
