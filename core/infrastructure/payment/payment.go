/**
 * Copyright 2015 @ z3q.net.
 * name : payment
 * author : jarryliu
 * date : 2015-07-27 21:51
 * description :
 * history :
 */
package payment

import (
	"github.com/ixre/gof/log"
	"net/http"
)

// 交易成功
const StatusTradeSuccess = 1

// 交易失败
const StatusTradeFail = 2

var (
	logF log.ILogger
)

func init() {
	//fi, _ := os.OpenFile("pay.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	//logF = log.NewLogger(fi, "payment", log.LOpen)
}

func Debug(format string, data ...interface{}) {
	//fmt.Printf("[ Pay][ Gateway]:"+format+"\n",data...)
	//logF.Printf("[ Pay][ Gateway]:"+format+"\n\n", data...)
}

type IPayment interface {
	// 创建网关
	CreateGateway(orderNo string, fee float32, subject, body, notifyUrl, returnUrl string) string
	// 返回
	Return(r *http.Request) Result
	// 通知
	Notify(r *http.Request) Result
}
