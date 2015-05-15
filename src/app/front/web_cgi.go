/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-01-22 21:08
 * description :
 * history :
 */
package front

import (
	"bufio"
	"fmt"
	"github.com/atnet/gof/web"
	"go2o/src/core/infrastructure/tool"
	"os"
	"strings"
	"time"
)

const (
	UPLOAD_DIR = "static/uploads/"
)

// Web同一网关接口
type WebCgi struct{}

func (this *WebCgi) Upload(key string, ctx *web.Context, savedir string) []byte {
	r, _ := ctx.Request, ctx.ResponseWriter
	var ext string
	var filePath string
	var err error
	fi, h, err := r.FormFile(key)
	if err != nil {
		return []byte("<html><head><title>" + err.Error() +
			"</title></head></html>")
	}
	ext = h.Filename[strings.LastIndex(h.Filename, "."):]
	filePath = strings.Join([]string{savedir, time.Now().Format("20060102150304"), ext}, "")
	os.MkdirAll(UPLOAD_DIR+savedir, os.ModePerm)

	f, err := os.OpenFile(UPLOAD_DIR+filePath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		os.ModePerm)

	if err != nil {
		return []byte("<html><head><title>" + err.Error() +
			"</title></head></html>")
	}

	buf := bufio.NewWriter(f)
	bufSize := 100
	buffer := make([]byte, bufSize)

	for {
		n, err := fi.Read(buffer)
		if err != nil {
			break
		}
		if n != bufSize {
			buf.Write(buffer[:n])
		} else {
			buf.Write(buffer)
		}
	}

	return []byte("{url:'" + filePath + "'}")

}

//获取位置
func (this *WebCgi) GeoLocation(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	ip := r.RemoteAddr[:strings.Index(r.RemoteAddr, ":")]
	addr := tool.GetLocation(ip)
	w.Write([]byte(fmt.Sprintf(`{"ip":"%s","addr":"%s"}`, ip, addr)))
}
