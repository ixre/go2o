/**
 * Copyright 2015 @ 56x.net.
 * name : img.go
 * author : jarryliu
 * date : 2016-08-23 12:53
 * description :
 * history :
 */
package util

import (
	"bytes"
	"github.com/smartwalle/resize"
	"image"
	"image/jpeg"
	"io"
)

// 生成缩略图
func MakeThumbnail(r io.Reader, width, height uint) ([]byte, error) {
	img, _, err := image.Decode(r)
	if err == nil {
		aResize := resize.Resize(width, height, img, resize.Lanczos3)
		w := bytes.NewBuffer(nil)
		jpeg.Encode(w, aResize, &jpeg.Options{100})
		return w.Bytes(), err
	}
	return nil, err
}
