/**
 * Copyright 2015 @ S1N1 Team.
 * name : member.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package util
import (
    "fmt"
    "github.com/atnet/gof"
    "github.com/atnet/gof/crypto"
    "errors"
)

const offset string = "%$^&@#"

func chkStorage(sto gof.Storage) {
    if sto == nil {
        panic(errors.New("[ Api] - api token storage is null !"))
    }
}


// 获取会员API调用密钥Key
func GetMemberApiTokenKey(memberId int) string {
    return fmt.Sprintf("api:member:token:%d", memberId)
}


// 设置令牌，并返回
func SetMemberApiToken(sto gof.Storage, memberId int, pwd string) string {
    chkStorage(sto)
    cyp := crypto.NewUnixCrypto(pwd+offset, offset)
    var token string = string(cyp.Encode())
    var key string = GetMemberApiTokenKey(memberId)

    sto.Set(key, token)      // 存储令牌
    sto.Set(key+"base", pwd) // 存储令牌凭据

    return token
}

// 获取会员的API令牌
func GetMemberApiToken(sto gof.Storage,memberId int)(string,string){
    chkStorage(sto)

    var key = GetMemberApiTokenKey(memberId)
    var srcToken, tokenBase string

    srcToken,_ = sto.GetString(key)
    tokenBase, _ = sto.GetString(key+"base")
    return srcToken,tokenBase
}

// 校验令牌
func CompareMemberApiToken(sto gof.Storage, memberId int, token string) bool {

    if len(token) == 0 {
        return false
    }

    srcToken, tokenBase := GetMemberApiToken(sto, memberId)
    if len(srcToken) == 0 || len(tokenBase) == 0 {
        return false
    }

    return srcToken == token

    //    cyp := crypto.NewUnixCrypto(tokenBase+offset, offset)
    //    b, _, _ := cyp.Compare(token)
    //    return b
}