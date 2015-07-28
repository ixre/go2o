/**
 * Copyright 2015 @ S1N1 Team.
 * name : payment
 * author : jarryliu
 * date : 2015-07-27 21:51
 * description :
 * history :
 */
package payment

import (
	"github.com/atnet/gof/log"
	"os"
)

// 交易成功
const StatusTradeSuccess = 1

var (
	logF log.ILogger
)

func init() {
	//fi, _ := os.OpenFile("pay.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	//logF = log.NewLogger(fi, "payment", log.LOpen)
}

func Debug(format string, data ...interface{}) {
	//logF.Printf(format+"\n\n", data...)
}
