package tencent

import (
	"errors"
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

type WxOAuthSession struct {
	auth.ResCode2Session
	// 小程序或公众号应用ID
	AppId string
}

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
	_repo     registry.IRegistryRepo
	_wc       *wechat.Wechat
	_mpConfig *config.Config
}

// NewWechat init
func NewWechat(o storage.Interface, repo registry.IRegistryRepo) *Wechat {
	c := &Wechat{
		_repo: repo,
		_wc:   wechat.NewWechat(),
	}
	c._wc.SetCache(&wechatCache{o})
	c.initConfig()
	return c
}

func (w *Wechat) initConfig() {
	w._repo.CreateUserKey("wechat_mp_app_id", "", "微信小程序AppId")
	w._repo.CreateUserKey("wechat_mp_app_secret", "", "微信小程序AppSecret")
	w._repo.CreateUserKey("wechat_mp_aes", "", "微信小程序AESKey")
	cfg := &config.Config{}
	cfg.AppID, _ = w._repo.GetValue("wechat_mp_app_id")
	cfg.AppSecret, _ = w._repo.GetValue("wechat_mp_app_secret")
	cfg.EncodingAESKey, _ = w._repo.GetValue("wechat_mp_aes")
	w._mpConfig = cfg
}

// 根据JsCode获取会话信息
func (w *Wechat) GetOpenId(jsCode string, cfg *config.Config) (*WxOAuthSession, error) {
	if cfg == nil {
		cfg = w._mpConfig
		if cfg.AppID == "" || cfg.AppSecret == "" {
			return nil, errors.New("微信小程序配置信息未初始化")
		}
	}
	ret, err := w._wc.GetMiniProgram(cfg).GetAuth().Code2Session(jsCode)
	if err != nil {
		return nil, err
	}
	return &WxOAuthSession{
		ResCode2Session: ret,
		AppId:           cfg.AppID,
	}, err
}
