package service

import (
	"github.com/ixre/gof/crypto"
	"go2o/core/domain/interface/registry"
	"go2o/core/service/impl"
	"strings"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : init.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-14 11:35
 * description :
 * history :
 */


func sysInit(){
	repo := impl.Repos.GetRegistryRepo()
	initJWTSecret(repo)
}

// 初始化jwt密钥
func initJWTSecret(repo registry.IRegistryRepo) {
	value, _ := repo.GetValue(registry.SysJWTSecret)
	if strings.TrimSpace(value) == "" {
		_, privateKey, _ := crypto.GenRsaKeys(2048)
		_ = repo.UpdateValue(registry.SysJWTSecret, privateKey)
	}
}