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
	"image/jpeg"

	"github.com/disintegration/imaging"
)

// 生成缩略图
func MakeThumbnail(filename string, width, height int) ([]byte, error) {
	img, err := imaging.Open(filename, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}
	rate := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
	if height <= 0 {
		height = int(float64(width) / rate)
	}
	if width <= 0 {
		width = int(float64(height) * rate)
	}
	dstImg := imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
	buf := bytes.NewBuffer(nil)
	err = jpeg.Encode(buf, dstImg, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
