package handler

import (
	"log"
	"strings"

	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/util"
)

// 子订单推送
func (h EventHandler) HandleAppInitialEvent(data interface{}) {
	v := data.(*events.AppInitialEvent)
	if v == nil {
		return
	}
	initJWTSecret(h.registryRepo)
	initSuperLoginToken(h.registryRepo)

}

func initSuperLoginToken(repo registry.IRegistryRepo) {
	value, _ := repo.GetValue(registry.SysSuperLoginToken)
	if strings.TrimSpace(value) == "" {
		pwd := util.RandString(8)
		log.Printf(`[ GO2O][ INFO]: the initial super pwd is '%s', it only show first time. plese save it.\n`, pwd)
		token := domain.Sha1("master" + crypto.Md5([]byte(pwd)))
		_ = repo.UpdateValue(registry.SysSuperLoginToken, token)
	}

}

// 初始化jwt密钥
func initJWTSecret(repo registry.IRegistryRepo) {
	value, _ := repo.GetValue(registry.SysJWTSecret)
	if strings.TrimSpace(value) == "" {
		_, privateKey, _ := crypto.GenRsaKeys(2048)
		_ = repo.UpdateValue(registry.SysJWTSecret, privateKey)
	}
}
