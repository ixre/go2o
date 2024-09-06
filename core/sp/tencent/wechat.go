package tencent

import (
	wechat "github.com/silenceper/wechat/v2"
)

type Wechat struct {
}

func NewWechat() *Wechat {
	wechat.NewWechat()
	return &Wechat{}
}
