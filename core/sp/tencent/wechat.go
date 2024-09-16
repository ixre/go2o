package tencent

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/logger"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/storage"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
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
	err := w.o.Get(key, &dst)
	if err != nil {
		return ""
	}
	return dst
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
	_repo      registry.IRegistryRepo
	_wc        *wechat.Wechat
	_mpConfig  *config.Config
	_mpVersion string
	_mux       sync.RWMutex
	_configs   map[string]*config.Config
}

// NewWechat init
func NewWechat(o storage.Interface, repo registry.IRegistryRepo) *Wechat {
	c := &Wechat{
		_repo:      repo,
		_wc:        wechat.NewWechat(),
		_mpVersion: "release",
	}
	c._wc.SetCache(&wechatCache{o})
	c.initConfig()
	return c
}

func (w *Wechat) initConfig() {
	w._repo.CreateUserKey("wechat_mp_app_id", "", "微信小程序AppId")
	w._repo.CreateUserKey("wechat_mp_app_secret", "", "微信小程序AppSecret")
	w._repo.CreateUserKey("wechat_mp_aes", "", "微信小程序AESKey")
	w._repo.CreateUserKey("wechat_mp_version", "release", "微信小程序环境，正式版为:release,体验版为:trial,开发版为:develop")
	cfg := &config.Config{}
	cfg.AppID, _ = w._repo.GetValue("wechat_mp_app_id")
	cfg.AppSecret, _ = w._repo.GetValue("wechat_mp_app_secret")
	cfg.EncodingAESKey, _ = w._repo.GetValue("wechat_mp_aes")
	w._mpVersion, _ = w._repo.GetValue("wechat_mp_version")
	w._mpConfig = cfg
	w.AddConfig(cfg.AppID, cfg)
	w.AddConfig("mp*", cfg)
}

func (w *Wechat) AddConfig(appId string, cfg *config.Config) error {
	w._mux.Lock()
	defer w._mux.Unlock()
	if w._configs == nil {
		w._configs = make(map[string]*config.Config)
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		return errors.New("微信配置信息未初始化")
	}
	w._configs[appId] = cfg
	return nil
}

func (w *Wechat) getMiniProgramConfig(appId string) (*config.Config, error) {
	w._mux.RLock()
	defer w._mux.RUnlock()
	if appId == "" {
		appId = "mp*"
	}
	cfg, ok := w._configs[appId]
	if !ok {
		return nil, errors.New("微信小程序配置信息未初始化")
	}
	if cfg.AppID == "" || cfg.AppSecret == "" {
		return nil, errors.New("微信小程序配置信息未初始化")
	}
	if w._mpVersion != "release" && w._mpVersion != "trial" && w._mpVersion != "develop" {
		return nil, errors.New("微信小程序环境配置错误, 非预期值:" + w._mpVersion)
	}
	return cfg, nil
}

// 根据JsCode获取会话信息
func (w *Wechat) GetMiniProgramOpenId(appId string, jsCode string) (*WxOAuthSession, error) {
	cfg, err := w.getMiniProgramConfig(appId)
	if err != nil {
		return nil, err
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

// GetMiniProgramUnlimitCode 获取小程序无限二维码
func (w *Wechat) GetMiniProgramUnlimitCode(appId, ownerKey string, page string, scene string) ([]byte, error) {
	cfg, err := w.getMiniProgramConfig(appId)
	if err != nil {
		return nil, err
	}
	// 获取二维码文件Key
	key := fmt.Sprintf("mp-%s-%s-page:%s-scene:%s", cfg.AppID, w._mpVersion, page, scene)
	sign := crypto.Md5([]byte(key))[6:24]
	// 生成二维码文件路径
	filePath := fmt.Sprintf("./files/mp/qrcode/%s-%s.png", ownerKey, sign)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		// 不存在文件,调用接口生成二维码
		mp := w._wc.GetMiniProgram(cfg)
		bytes, err := mp.GetQRCode().GetWXACodeUnlimit(qrcode.QRCoder{
			Page:       page,
			Scene:      scene,
			IsHyaline:  true,
			EnvVersion: w._mpVersion,
		})
		if err == nil {
			// 检查目录是否存在并写入文件中
			dir := filepath.Dir(filePath) //检查目录是否存在
			if _, err = os.Stat(dir); os.IsNotExist(err) {
				//创建目录
				if err = os.MkdirAll(dir, os.ModePerm); err != nil {
					return nil, err
				}
			}
			err = os.WriteFile(filePath, bytes, os.ModePerm)
		}
		if err != nil {
			logger.Error("生成小程序二维码失败: %s", err.Error())
		}
		return bytes, err
	}
	if err == nil {
		// 已存在文件,直接读取
		return os.ReadFile(filePath)
	}
	return nil, err
}
