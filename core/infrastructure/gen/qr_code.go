/**
 * Copyright 2015 @ 56x.net.
 * name : qr_code
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package gen

import (
	"github.com/ixre/go2o/core/infrastructure/gen/rsc/qr"
)

// 生成网址对应的二维码
func BuildQrCodeForUrl(url string, scale int) []byte {
	if code, err := qr.Encode(url, qr.L); err == nil {
		code.Scale = scale
		return code.PNG()
	}
	return []byte{}
}
