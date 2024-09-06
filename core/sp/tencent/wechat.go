package tencent

import (
	"time"

	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/gof/storage"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	"github.com/silenceper/wechat/v2/miniprogram/config"
)

// 全局微信实例
var WECHAT *Wechat

// 初始化默认微信实例
func Configure(o storage.Interface, repo registry.IRegistryRepo) {
	WECHAT = NewWechat(o, repo)
}

var _ cache.Cache = new(wechatCache)

type wechatCache struct {
	o storage.Interface
}

// Delete implements cache.Cache.
func (w *wechatCache) Delete(key string) error {
	w.o.Delete(key)
	return nil
}

// Get implements cache.Cache.
func (w *wechatCache) Get(key string) interface{} {
	var dst interface{}
	return w.o.Get(key, &dst)
}

// IsExist implements cache.Cache.
func (w *wechatCache) IsExist(key string) bool {
	return w.o.Exists(key)
}

// Set implements cache.Cache.
func (w *wechatCache) Set(key string, val interface{}, timeout time.Duration) error {
	return w.o.SetExpire(key, val, int64(time.Second))
}

type Wechat struct {
	_repo registry.IRegistryRepo
	_wc   *wechat.Wechat
}

// NewWechat init
func NewWechat(o storage.Interface, repo registry.IRegistryRepo) *Wechat {
	c := &Wechat{
		_repo: repo,
		_wc:   wechat.NewWechat(),
	}
	c._wc.SetCache(&wechatCache{o})
	return c
}

// 根据JsCode获取会话信息
func (w *Wechat) GetOpenId(jsCode string, cfg *config.Config) (auth.ResCode2Session, error) {
	return w._wc.GetMiniProgram(cfg).GetAuth().Code2Session(jsCode)
}
