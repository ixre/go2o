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
<<<<<<< HEAD
	"go2o/src/core/infrastructure/tool"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
=======
	"github.com/jsix/gof/web"
	"go2o/src/core/infrastructure/tool"
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	"io"
	"os"
	"strings"
	"time"
)

<<<<<<< HEAD
var (
	uploadSaveDir string
)

// Web同一网关接口
type WebCgi struct {
}

func (this *WebCgi) Upload(key string, ctx *echox.Context, saveDir string) []byte {
	r := ctx.Request()
=======
const (
	UPLOAD_DIR = "static/uploads/"
)

// Web同一网关接口
type WebCgi struct{}

func (this *WebCgi) Upload(key string, ctx *web.Context, savedir string) []byte {
	r, _ := ctx.Request, ctx.Response
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	var ext string
	var filePath string
	var err error
	fi, h, err := r.FormFile(key)
	if err != nil {
		return []byte("<html><head><title>" + err.Error() +
			"</title></head></html>")
	}
<<<<<<< HEAD

	if len(uploadSaveDir) == 0 {
		uploadSaveDir = ctx.App.Config().GetString(variable.UploadSaveDir) // 上传存放目录
	}

	ext = h.Filename[strings.LastIndex(h.Filename, "."):]
	filePath = strings.Join([]string{saveDir, time.Now().Format("20060102150304"), ext}, "")
	os.MkdirAll(uploadSaveDir+saveDir, os.ModePerm)

	f, err := os.OpenFile(uploadSaveDir+filePath,
=======
	ext = h.Filename[strings.LastIndex(h.Filename, "."):]
	filePath = strings.Join([]string{savedir, time.Now().Format("20060102150304"), ext}, "")
	os.MkdirAll(UPLOAD_DIR+savedir, os.ModePerm)

	f, err := os.OpenFile(UPLOAD_DIR+filePath,
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
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
<<<<<<< HEAD
func (this *WebCgi) GeoLocation(ctx *echox.Context) {
	r, w := ctx.Request(), ctx.Response()
=======
func (this *WebCgi) GeoLocation(ctx *web.Context) {
	r, w := ctx.Request, ctx.Response
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	ip := r.RemoteAddr[:strings.Index(r.RemoteAddr, ":")]
	add := tool.GetLocation(ip)
	w.Write([]byte(fmt.Sprintf(`{"ip":"%s","addr":"%s"}`, ip, add)))
}
