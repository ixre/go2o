/**
 * Copyright 2015 @ z3q.net.
 * name : qr_code
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package gen

import (
	"go2o/src/core/infrastructure/gen/rsc/qr"
)

// 生成网址对应的二维码
func BuildQrCodeForUrl(url string, scale int) []byte {
	if code, err := qr.Encode(url, qr.M); err == nil {
		code.Scale = scale
		return code.PNG()
	}
	return []byte{}
}
