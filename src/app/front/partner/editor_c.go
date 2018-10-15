/**
 * Copyright 2015 @ z3q.net.
 * name : editor_c.go
 * author : jarryliu
 * date : 2015-08-18 17:09
 * description :
 * history :
 */
package partner

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jsix/gof"
	"go2o/src/core/variable"
	"go2o/src/x/echox"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var _ sort.Interface = new(SorterFiles)

type SorterFiles struct {
	files  []os.FileInfo
	sortBy string
}

func (this *SorterFiles) Len() int {
	return len(this.files)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (this *SorterFiles) Less(i, j int) bool {
	switch strings.ToLower(this.sortBy) {
	case "size":
		return this.files[i].Size() < this.files[j].Size()
	case "name":
		return this.files[i].Name() < this.files[j].Name()
	case "type":
		iN, jN := this.files[i].Name(), this.files[j].Name()
		return iN[strings.Index(iN, ".")+1:] < jN[strings.Index(jN, ".")+1:]
	}
	return true
}

// Swap swaps the elements with indexes i and j.
func (this *SorterFiles) Swap(i, j int) {
	tmp := this.files[i]
	this.files[i] = this.files[j]
	this.files[j] = tmp
}

//图片扩展名
var imgFileTypes string = "gif,jpg,jpeg,png,bmp"
var (
	moveUpRegexp *regexp.Regexp = regexp.MustCompile("(.*?)[^\\/]+\\/$")
)

// 文件管理
// @rootDir : 根目录路径，相对路径
// @rootUrl : 根目录URL，可以指定绝对路径，比如 http://www.yoursite.com/attached/
func fileManager(r *http.Request, rootDir, rootUrl string) ([]byte, error) {
	var currentPath = ""
	var currentUrl = ""
	var currentDirPath = ""
	var moveUpDirPath = ""
	var dirPath string = rootDir

	urlQuery := r.URL.Query()
	var dirName string = urlQuery.Get("dir")

	if len(dirName) != 0 {
		if dirName == "image" || dirName == "flash" ||
			dirName == "media" || dirName == "file" {
			dirPath += dirName + "/"
			rootUrl += dirName + "/"
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				os.MkdirAll(dirPath, os.ModePerm)
			}
		} else {
			return nil, errors.New("Invalid Directory name")
		}
	}

	//根据path参数，设置各路径和URL
	var path string = urlQuery.Get("path")
	if len(path) == 0 {
		currentPath = dirPath
		currentUrl = rootUrl
		currentDirPath = ""
		moveUpDirPath = ""
	} else {
		currentPath = dirPath + path
		currentUrl = rootUrl + path
		currentDirPath = path
		moveUpDirPath = moveUpRegexp.ReplaceAllString(currentDirPath, "$1")
	}

	//不允许使用..移动到上一级目录
	if strings.Index(path, "\\.\\.") != -1 {
		return nil, errors.New("Access is not allowed.")
	}

	//最后一个字符不是/
	if path != "" && !strings.HasSuffix(path, "/") {
		return nil, errors.New("Parameter is not valid.")
	}
	//目录不存在或不是目录
	dir, err := os.Stat(currentPath)
	if os.IsNotExist(err) || !dir.IsDir() {
		return nil, errors.New("no such directory or file not directory,path:" + currentPath)
	}

	//排序形式，name or size or type
	var order string = strings.ToLower(urlQuery.Get("order"))

	//遍历目录取得文件信息
	var dirList *SorterFiles = &SorterFiles{
		files:  []os.FileInfo{},
		sortBy: order,
	}
	var fileList *SorterFiles = &SorterFiles{
		files:  []os.FileInfo{},
		sortBy: order,
	}

	// 遍历目录获取子目录和文件
	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		if v.IsDir() {
			dirList.files = append(dirList.files, v)
		} else {
			fileList.files = append(fileList.files, v)
		}
	}

	// 排序
	sort.Sort(dirList)
	sort.Sort(fileList)

	var result = make(map[string]interface{})
	result["moveup_dir_path"] = moveUpDirPath
	result["current_dir_path"] = currentDirPath
	result["current_url"] = currentUrl
	result["total_count"] = dirList.Len() + fileList.Len()
	var dirFileList = []map[string]interface{}{}
	for i := 0; i < dirList.Len(); i++ {
		hash := make(map[string]interface{})
		fs, _ := ioutil.ReadDir(currentPath + "/" + dirList.files[i].Name())
		//fmt.Println("----", currentPath+"/"+dirList.files[i].Name())
		hash["is_dir"] = true
		hash["has_file"] = len(fs) > 0
		hash["is_photo"] = false
		hash["filetype"] = ""
		hash["filename"] = dirList.files[i].Name()
		hash["datetime"] = dirList.files[i].ModTime().Format("2006-01-02 15:04:05")
		dirFileList = append(dirFileList, hash)
	}

	var fN, ext string
	for i := 0; i < fileList.Len(); i++ {
		hash := make(map[string]interface{})
		fN = fileList.files[i].Name()
		ext = fN[strings.Index(fN, ".")+1:]
		hash["is_dir"] = false
		hash["has_file"] = false
		hash["filesize"] = fileList.files[i].Size()
		hash["is_photo"] = strings.Index(imgFileTypes, ext)
		hash["filetype"] = ext
		hash["filename"] = fN
		hash["datetime"] = fileList.files[i].ModTime().Format("2006-01-02 15:04:05")
		dirFileList = append(dirFileList, hash)
	}

	result["file_list"] = dirFileList
	return json.Marshal(result)
}

// 文件上传
func fileUpload(r *http.Request, savePath, rootPath string) (fileUrl string, err error) {

	//定义允许上传的文件扩展名
	var extTable map[string]string = map[string]string{
		"image": "gif,jpg,jpeg,png,bmp",
		"flash": "swf,flv",
		"media": "swf,flv,mp3,wav,wma,wmv,mid,avi,mpg,asf,rm,rmvb",
		"file":  "doc,docx,xls,xlsx,ppt,htm,html,txt,zip,rar,gz,bz2,7z,pdf",
	}

	//最大文件大小
	const maxSize int64 = 1000000

	// 取得上传文件
	r.ParseMultipartForm(maxSize)
	f, header, err := r.FormFile("imgFile")
	if f == nil {
		return "", errors.New("no such upload file")
	}
	if err != nil {
		return "", err
	}

	var fileName string = header.Filename
	var extIdx = strings.LastIndex(fileName, ".")
	if extIdx == -1 {
		return "", errors.New("Unkown file type")
	}
	var fileExt string = strings.ToLower(fileName[extIdx+1:])

	// 检查上传目录
	var dirPath string = rootPath
	var dirName string = r.URL.Query().Get("dir")
	if len(dirName) == 0 {
		dirName = "image"
	}
	if _, ok := extTable[dirName]; !ok {
		return "", errors.New("incorrent file type")
	}

	// 检查扩展名
	if strings.Index(extTable[dirName], fileExt) == -1 &&
		!strings.HasSuffix(extTable[dirName], fileExt) {
		return "", errors.New("上传文件扩展名是不允许的扩展名。\n只允许" + extTable[dirName] + "格式。")
	}

	// 检查上传超出文件大小
	if i, _ := strconv.Atoi(header.Header.Get("Content-Length")); int64(i) > maxSize {
		return "", errors.New("上传文件大小超过限制。")
	}

	//创建文件夹
	dirPath += dirName + "/"
	savePath += dirName + "/"

	var now = time.Now()
	var ymd string = now.Format("200601")
	dirPath += ymd + "/"
	savePath += ymd + "/"

	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		os.MkdirAll(savePath, os.ModePerm)
	}

	var newFileName string = fmt.Sprintf("%d_%d.%s", now.Unix(),
		100+rand.Intn(899), fileExt)
	var filePath string = savePath + newFileName

	fi, err := os.OpenFile(filePath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		os.ModePerm)

	if err == nil {
		defer fi.Close()
		buf := bufio.NewWriter(fi)
		bufSize := 100
		buffer := make([]byte, bufSize)
		var n int
		var leng int
		for {
			if n, err = f.Read(buffer); err == io.EOF {
				break
			}

			if n != bufSize {
				buf.Write(buffer[:n])
			} else {
				buf.Write(buffer)
			}

			leng += n
		}
		buf.Flush()
	}

	return dirPath + newFileName, nil
}

type editorC struct {
}

func (this *editorC) File_manager(ctx *echox.Context) error {
	partnerId := getPartnerId(ctx)
	upDir := ctx.App.Config().GetString(variable.UploadSaveDir)
	d, err := fileManager(ctx.HttpRequest(),
		fmt.Sprintf("%s/%d/upload/", upDir, partnerId),
		fmt.Sprintf("%s/%d/upload/", ctx.App.Config().GetString(variable.ImageServer), partnerId),
	)
	if err != nil {
		return ctx.JSON(http.StatusOK, gof.Result{ErrMsg: err.Error()})
	}
	return ctx.JSON(http.StatusOK, d)
}

func (this *editorC) File_upload(ctx *echox.Context) error {
	if ctx.Request().Method != "POST" {
		return errors.New("error request method")
	}
	partnerId := getPartnerId(ctx)
	upDir := ctx.App.Config().GetString(variable.UploadSaveDir)
	fileUrl, err := fileUpload(ctx.HttpRequest(),
		fmt.Sprintf("%s/%d/upload/", upDir, partnerId),
		fmt.Sprintf("%s/%d/upload/",
			ctx.App.Config().GetString(variable.ImageServer), partnerId),
	)
	var hash map[string]interface{} = make(map[string]interface{})
	if err == nil {
		hash["error"] = 0
		hash["url"] = fileUrl
	} else {
		hash["error"] = 1
		hash["message"] = err.Error()
	}
	return ctx.JSON(http.StatusOK, hash)
}
