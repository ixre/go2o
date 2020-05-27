/**
 * Copyright 2015 @ to2.net.
 * name : img.go
 * author : jarryliu
 * date : 2016-08-23 12:53
 * description :
 * history :
 */
package tool

import (
	"bytes"
	"github.com/nfnt/resize"
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
