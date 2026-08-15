/**
 * Copyright 2015 @ 56x.net.
 * name : download.go
 * author : jarryliu
 * date : 2015-12-31 12:23
 * description :
 * history :
 */
package domain

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func getFileName(disposition string, guessExt string) string {
	ext := guessExt
	if len(disposition) > 0 && strings.Index(disposition, "filename=") != -1 {
		ext = disposition[strings.Index(disposition, ".")+1:]
	}
	return fmt.Sprintf("%d_%d.%s", time.Now().Unix(),
		100+rand.Intn(899), ext)
}

// 下载远程资源并返回本地地址
func DownloadToLocal(url string, savePath string, ext string) string {
	var req *http.Request
	var err error
	req, err = http.NewRequest("GET", url, nil)
	if err == nil {
		req.Header.Set("User_Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:133.0) Gecko/20100101 Firefox/133.0")
		rsp, err := http.DefaultClient.Do(req)
		if err == nil {
			fileName := getFileName(rsp.Header.Get("Content-Disposition"), ext)
			var filePath = savePath + fileName
			dir := path.Dir(filePath)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.MkdirAll(dir, os.ModePerm)
			}
			src := rsp.Body
			defer src.Close()

			fi, err := os.OpenFile(filePath,
				os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
				os.ModePerm)

			if err == nil {
				defer fi.Close()
				buf := bufio.NewWriter(fi)
				bufSize := 100
				buffer := make([]byte, bufSize)
				var n int
				var l int
				for {
					if n, err = src.Read(buffer); err == io.EOF {
						break
					}

					if n != bufSize {
						buf.Write(buffer[:n])
					} else {
						buf.Write(buffer)
					}

					l += n
				}
				buf.Flush()

				return filePath[2:] // 去掉"./"
			}
		}
	}

	if err != nil {
		log.Println("[ Download]- ", url, "\n  ", err.Error())
	}

	return url
}
