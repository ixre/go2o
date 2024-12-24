/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-10 21:16
 * description :
 * history :
 */

package domain

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"

	"github.com/ixre/gof/crypto"
)

var (
	// 私钥
	privateKey = ""
)

func ConfigureHmacPrivateKey(key string) {
	// 解析私钥
	block, _ := pem.Decode([]byte(key))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		panic("privateKey is not a valid RSA private key")
	}

	_, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic("privateKey is not a valid RSA private key:" + err.Error())
	}
	privateKey = strings.TrimSpace(key)
}

// MD5加密
func Md5(pwd string) string {
	if strings.TrimSpace(pwd) == "" {
		return ""
	}
	return crypto.Md5([]byte(pwd))
}

func ChkPwdRight(pwd string) (bool, error) {
	if len(pwd) < 6 {
		return false, errors.New("密码必须大于6位")
	}
	return true, nil
}

// HmacSha256
func HmacSha256(s string) string {
	if len(privateKey) == 0 {
		panic("privateKey is empty, please call ConfigPrivateKey to set it")
	}
	return crypto.HmacSha256([]byte(s), []byte(privateKey))
}

// 加密会员密码,因为可能会使用手机号码登录，
// 所以密码不能依据用户名作为生成凭据
func MemberSha256Pwd(pwd string, salt string) string {
	if strings.TrimSpace(pwd) == "" {
		return ""
	}
	return HmacSha256(pwd + salt)
}

// 交易密码
func TradePassword(pwd string, salt string) string {
	if strings.TrimSpace(pwd) == "" {
		return ""
	}
	return HmacSha256(pwd + salt)
}

// 加密合作商密码
func MerchantSha265Pwd(pwd string, salt string) string {
	if strings.TrimSpace(pwd) == "" {
		return ""
	}
	return HmacSha256(pwd + salt)
}

// 超级管理员密码
func SuperPassword(username, pwd string) string {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(pwd) == "" {
		return ""
	}
	return HmacSha256(username + pwd)
}

// 系统用户密码
func RbacPassword(pwd string, salt string) string {
	if strings.TrimSpace(pwd) == "" {
		return ""
	}
	return HmacSha256(pwd + salt)
}

func AesEncrypt(data string, key string, iv string) string {
	if strings.TrimSpace(data) == "" {
		return ""
	}
	return crypto.AesEncrypt([]byte(data), []byte(key), []byte(iv))
}

func AesDecrypt(data string, key string, iv string) string {
	if strings.TrimSpace(data) == "" {
		return ""
	}
	return crypto.AesDecrypt([]byte(data), []byte(key), []byte(iv))
}
