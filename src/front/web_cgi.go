/**
 * Copyright 2014 @ z3q.net.
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
	"go2o/src/core/infrastructure/tool"
	"go2o/src/x/echox"
	"io"
	"os"
	"strings"
	"time"
	"go2o/src/core/variable"
)

// Web同一网关接口
type WebCgi struct{}

func (this *WebCgi) Upload(key string, ctx *echox.Context, savedir string) []byte {
	r := ctx.Request()
	var ext string
	var filePath string
	var err error
	fi, h, err := r.FormFile(key)
	if err != nil {
		return []byte("<html><head><title>" + err.Error() +
			"</title></head></html>")
	}

	var upSaveDir string = ctx.App.Config().Get(variable.UploadSaveDir) // 上传存放目录

	ext = h.Filename[strings.LastIndex(h.Filename, "."):]
	filePath = strings.Join([]string{savedir, time.Now().Format("20060102150304"), ext}, "")
	os.MkdirAll(upSaveDir+savedir, os.ModePerm)

	f, err := os.OpenFile(upSaveDir+filePath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		os.ModePerm)

	if err != nil {
		return []byte("<html><head><title>" + err.Error() +
			"</title></head></html>")
	}

	defer f.Close()

	if err == nil {
		buf := bufio.NewWriter(f)
		bufSize := 100
		buffer := make([]byte, bufSize)
		var n int
		var totalLen int
		for {
			if n, err = fi.Read(buffer); err == io.EOF {
				break
			}

			if n != bufSize {
				buf.Write(buffer[:n])
			} else {
				buf.Write(buffer)
			}

			totalLen += n
		}
		buf.Flush()

		return []byte(fmt.Sprintf("{url:'%s',len:%d}", filePath, totalLen))
	}
	return []byte("{error:'" + err.Error() + "'}")
}

//获取位置
func (this *WebCgi) GeoLocation(ctx *echox.Context) {
	r, w := ctx.Request(), ctx.Response()
	ip := r.RemoteAddr[:strings.Index(r.RemoteAddr, ":")]
	add := tool.GetLocation(ip)
	w.Write([]byte(fmt.Sprintf(`{"ip":"%s","addr":"%s"}`, ip, add)))
}
