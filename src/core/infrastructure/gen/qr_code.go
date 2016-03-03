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
<<<<<<< HEAD
func BuildQrCodeForUrl(url string, scale int) []byte {
	if code, err := qr.Encode(url, qr.L); err == nil {
		code.Scale = scale
=======
func BuildQrCodeForUrl(url string) []byte {
	if code, err := qr.Encode(url, qr.M); err == nil {
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
		return code.PNG()
	}
	return []byte{}
}
