/**
 * Copyright 2015 @ at3.net.
 * name : sso.go
 * author : jarryliu
 * date : 2016-11-25 13:02
 * description :
 * history :
 */
package module

import (
	"fmt"
	"github.com/ixre/gof"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/storage"
	"strings"
	"time"
)

var _ Module = new(MemberModule)

type MemberModule struct {
	app         gof.App
	storage     storage.Interface
	tokenHours  int64
	tokenOffset string
}

// 模块数据
func (m *MemberModule) SetApp(app gof.App) {
	m.app = app
	m.storage = app.Storage()
}

// 初始化模块
func (m *MemberModule) Init() {
	m.tokenHours = 24 * 7 //默认保存7天
	m.tokenOffset = ""
}

// 获取会员Token-Key
func (m *MemberModule) getMemberTokenKey(memberId int64) string {
	return fmt.Sprintf("go2o:module:member:token:%d", memberId)
}

// 获取会员的Token
func (m *MemberModule) getMemberToken(memberId int64) (string, string) {
	var key = m.getMemberTokenKey(memberId)
	var pubToken, tokenBase string
	pubToken, _ = m.storage.GetString(key)
	tokenBase, _ = m.storage.GetString(key + "base")
	return pubToken, tokenBase
}

// 获取会员Token
func (m *MemberModule) GetToken(memberId int64) string {
	pubToken, _ := m.getMemberToken(memberId)
	return pubToken
}

// 检查会员的会话Token是否正确，
func (m *MemberModule) CheckToken(memberId int64, token string) bool {
	token = strings.TrimSpace(token)
	if len(token) == 0 {
		return false
	}
	pubToken, tokenBase := m.getMemberToken(memberId)
	// 清除token
	if pubToken == "" || tokenBase == "" {
		m.RemoveToken(memberId)
		return false
	}
	// return pubToken == token
	//if m.serve.Debug() {
	//    m.serve.Log().Println("[ Module][ Member]: check token for ",
	//        memberId, "; ", pubToken, " IN:", token, " equals:",
	//        pubToken == token, "; Len:", len(pubToken), len(token))
	//}
	// 验证Token
	cyp := crypto.NewUnixCrypto(tokenBase+m.tokenOffset, m.tokenOffset)
	b, decBytes, unix := cyp.Compare([]byte(token))
	// token不匹配
	if !b {
		if m.app.Debug() {
			m.app.Log().Println("dec:", string(decBytes), "; real:"+tokenBase+m.tokenOffset)
		}
		return false
	}
	// token已过期
	if unix < time.Now().Add(time.Hour*time.Duration(-m.tokenHours)).Unix() {
		m.RemoveToken(memberId)
		return false
	}
	return b
}

// 移除会员Token
func (m *MemberModule) RemoveToken(memberId int64) {
	key := m.getMemberTokenKey(memberId)
	m.storage.Del(key)
	m.storage.Del(key + "base")
}

// 重设并返回会员的会员Token，token有效时间默认为60天
func (m *MemberModule) ResetToken(memberId int64, pwd string) string {
	cyp := crypto.NewUnixCrypto(pwd+m.tokenOffset, m.tokenOffset)
	println("--", pwd, "|", m.tokenOffset)
	var token = string(cyp.Encode())
	var key = m.getMemberTokenKey(memberId)
	// 存储令牌
	m.storage.SetExpire(key, token, m.tokenHours*3600)
	// 存储令牌凭据
	m.storage.SetExpire(key+"base", pwd, m.tokenHours*3600)
	return token
}
