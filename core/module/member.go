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
    "github.com/jsix/gof"
    "github.com/jsix/gof/storage"
    "fmt"
    "github.com/jsix/gof/crypto"
)

var _ Module = new(SSOModule)

type MemberModule struct {
    app         gof.App
    storage     storage.Interface
    tokenHours  int64
    tokenOffset string
}

// 模块数据
func (m *MemberModule) SetApp(app gof.App) {
    app = app
    m.storage = app.Storage()
}

// 初始化模块
func (m *MemberModule) Init() {
    m.tokenHours = 24 * 60 //默认保存2个月
    m.tokenOffset = "%$^&@at3.net"
    // m.tokenOffset = "%$^&@#"
}

// 获取会员Token-Key
func (m *MemberModule) getMemberTokenKey(memberId int32) string {
    return fmt.Sprintf("go2o:module:member:token:%d", memberId)
    //return fmt.Sprintf("go2o:api:member:token:%d", memberId)
}

// 获取会员的Token
func (m *MemberModule) getMemberToken(memberId int32) (string, string) {
    var key = m.getMemberTokenKey(memberId)
    var pubToken, tokenBase string
    pubToken, _ = m.storage.GetString(key)
    tokenBase, _ = m.storage.GetString(key + "base")
    return pubToken, tokenBase
}

// 获取会员Token
func (m *MemberModule) GetToken(memberId int32) string {
    pubToken, _ := m.getMemberToken(memberId)
    return pubToken
}

// 检查会员的会话Token是否正确，
func (m *MemberModule) CheckToken(memberId int32, token string) bool {
    if len(token) == 0 {
        return false
    }
    pubToken, tokenBase := m.getMemberToken(memberId)
    if len(pubToken) == 0 || len(tokenBase) == 0 {
        return false
    }
    //fmt.Println("-----",srcToken," IN:",token , " equals:",srcToken==token, len(srcToken),len(token))
    return pubToken == token
}


// 移除会员Token
func (m *MemberModule) RemoveToken(memberId int32) {
    key := m.getMemberTokenKey(memberId)
    m.storage.Del(key)
    m.storage.Del(key + "base")
}

// 重设并返回会员的会员Token，token有效时间默认为60天
func (m *MemberModule)ResetToken(memberId int32, pwd string) string {
    cyp := crypto.NewUnixCrypto(pwd + m.tokenOffset, m.tokenOffset)
    var token string = string(cyp.Encode())
    var key string = m.getMemberTokenKey(memberId)
    // 存储令牌
    m.storage.SetExpire(key, token, m.tokenHours * 3600)
    // 存储令牌凭据
    m.storage.SetExpire(key + "base", pwd, m.tokenHours * 3600)
    return token
}