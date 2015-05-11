/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package domain

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"github.com/atnet/gof/util"
	"github.com/atnet/gof/crypto"
)

const (
	minRand int = 1000000
	maxRand int = 9999999
)

//新订单号
func NewOrderNo(partnerId int) string {
	//PartnerId的首位和末尾再加7位随机数
	rand.Seed(time.Now().Unix())
	rd := minRand + rand.Intn(maxRand-minRand) //minRand - maxRand中间的随机数
	ptstr := strconv.Itoa(partnerId)
	return fmt.Sprintf("%s%s%d", ptstr[:1], ptstr[len(ptstr)-1:], rd)
}

// 创建邀请码(6位)
func GenerateInvitationCode() string{
	var seed string = fmt.Sprintf("%d%s",time.Now().Unix(),util.RandString(6))
	var md5 = crypto.Md5([]byte(seed))
	return md5[8:14]
}