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
	"github.com/jrsix/gof/web/mvc"
	"net/http"
	"strings"
	"errors"
	"os"
	"io/ioutil"
	"sort"
	"encoding/json"
	"github.com/jrsix/gof/web"
	"fmt"
	"gobx/share/variable"
)

var _ sort.Interface = new(SorterFiles)
type SorterFiles struct{
	files []os.FileInfo
	sortBy string
}

func (this *SorterFiles) Len() int{
	return len(this.files)
}
// Less reports whether the element with
// index i should sort before the element with index j.
func (this *SorterFiles) Less(i, j int) bool{
	switch this.sortBy {
	case "size":
		return this.files[i].Size() < this.files[j].Size()
	case "name":
		return this.files[i].Name() < this.files[j].Name()
	case "type":
		iN,jN := this.files[i].Name(),this.files[j].Name()
		return iN[strings.Index(iN,".")+1:] < jN[strings.Index(jN,".")+1:]
	}
	return true
}
// Swap swaps the elements with indexes i and j.
func (this *SorterFiles) Swap(i, j int) {
	tmp := this.files[i]
	this.files[i] = this.files[j]
	this.files[j] = tmp
}



/*
 public class EditorFileManager : IHttpHandler, System.Web.SessionState.IRequiresSessionState
    {
        public void ProcessRequest(HttpContext context)
        {
            String aspxUrl = context.Request.Path.Substring(0, context.Request.Path.LastIndexOf("/") + 1);

            string siteID = Logic.CurrentSite.SiteId.ToString();

            //根目录路径，相对路径
            String rootPath = String.Format("{0}s{1}/", CmsVariables.RESOURCE_PATH,siteID);
            //根目录URL，可以指定绝对路径，比如 http://www.yoursite.com/attached/
            string appPath = AtNet.Cms.Cms.Context.ApplicationPath;
            String rootUrl = String.Format("{0}/{1}s{2}/", appPath == "/" ? "" : appPath,
                CmsVariables.RESOURCE_PATH, siteID);

            //图片扩展名
            String fileTypes = "gif,jpg,jpeg,png,bmp";

            String currentPath = "";
            String currentUrl = "";
            String currentDirPath = "";
            String moveupDirPath = "";

            String dirPath = AppDomain.CurrentDomain.BaseDirectory + rootPath;
            String dirName = context.Request.QueryString["dir"];
            if (!String.IsNullOrEmpty(dirName))
            {
                if (Array.IndexOf("image,flash,media,file".Split(','), dirName) == -1)
                {
                    context.Response.Write("Invalid Directory name.");
                    context.Response.End();
                }
                dirPath += dirName + "/";
                rootUrl += dirName + "/";
                if (!Directory.Exists(dirPath))
                {
                    Directory.CreateDirectory(dirPath).Create();
                }
            }

            //根据path参数，设置各路径和URL
            String path = context.Request.QueryString["path"];
            path = String.IsNullOrEmpty(path) ? "" : path;
            if (path == "")
            {
                currentPath = dirPath;
                currentUrl = rootUrl;
                currentDirPath = "";
                moveupDirPath = "";
            }
            else
            {
                currentPath = dirPath + path;
                currentUrl = rootUrl + path;
                currentDirPath = path;
                moveupDirPath = Regex.Replace(currentDirPath, @"(.*?)[^\/]+\/$", "$1");
            }

            //排序形式，name or size or type
            String order = context.Request.QueryString["order"];
            order = String.IsNullOrEmpty(order) ? "" : order.ToLower();

            //不允许使用..移动到上一级目录
            if (Regex.IsMatch(path, @"\.\."))
            {
                context.Response.Write("Access is not allowed.");
                context.Response.End();
            }
            //最后一个字符不是/
            if (path != "" && !path.EndsWith("/"))
            {
                context.Response.Write("Parameter is not valid.");
                context.Response.End();
            }
            //目录不存在或不是目录
            if (!Directory.Exists(currentPath))
            {
                context.Response.Write("Directory does not exist.");
                context.Response.End();
            }

            //遍历目录取得文件信息
            string[] dirList = Directory.GetDirectories(currentPath);
            string[] fileList = Directory.GetFiles(currentPath);

            switch (order)
            {
                case "size":
                    Array.Sort(dirList, new NameSorter());
                    Array.Sort(fileList, new SizeSorter());
                    break;
                case "type":
                    Array.Sort(dirList, new NameSorter());
                    Array.Sort(fileList, new TypeSorter());
                    break;
                case "name":
                default:
                    Array.Sort(dirList, new NameSorter());
                    Array.Sort(fileList, new NameSorter());
                    break;
            }

            Hashtable result = new Hashtable();
            result["moveup_dir_path"] = moveupDirPath;
            result["current_dir_path"] = currentDirPath;
            result["current_url"] = currentUrl;
            result["total_count"] = dirList.Length + fileList.Length;
            List<Hashtable> dirFileList = new List<Hashtable>();
            for (int i = 0; i < dirList.Length; i++)
            {
                DirectoryInfo dir = new DirectoryInfo(dirList[i]);
                Hashtable hash = new Hashtable();
                hash["is_dir"] = true;
                hash["has_file"] = (dir.GetFileSystemInfos().Length > 0);
                hash["filesize"] = 0;
                hash["is_photo"] = false;
                hash["filetype"] = "";
                hash["filename"] = dir.Name;
                hash["datetime"] = dir.LastWriteTime.ToString("yyyy-MM-dd HH:mm:ss");
                dirFileList.Add(hash);
            }
            for (int i = 0; i < fileList.Length; i++)
            {
                FileInfo file = new FileInfo(fileList[i]);
                Hashtable hash = new Hashtable();
                hash["is_dir"] = false;
                hash["has_file"] = false;
                hash["filesize"] = file.Length;
                hash["is_photo"] = (Array.IndexOf(fileTypes.Split(','), file.Extension.Substring(1).ToLower()) >= 0);
                hash["filetype"] = file.Extension.Substring(1);
                hash["filename"] = file.Name;
                hash["datetime"] = file.LastWriteTime.ToString("yyyy-MM-dd HH:mm:ss");
                dirFileList.Add(hash);
            }

            string files = String.Empty;
            int j = 0;
            foreach (Hashtable h in dirFileList)
            {
                files += JsonAnalyzer.ToJson(h);
                if (++j < dirFileList.Count)
                {
                    files += ",";
                }
            }
            result["file_list"] = "[" + files + "]";
            context.Response.AddHeader("Content-Type", "application/json; charset=UTF-8");
            context.Response.Write(JsonAnalyzer.ToJson(result));
            context.Response.End();
        }

        public class NameSorter : IComparer
        {
            public int Compare(object x, object y)
            {
                if (x == null && y == null)
                {
                    return 0;
                }
                if (x == null)
                {
                    return -1;
                }
                if (y == null)
                {
                    return 1;
                }
                FileInfo xInfo = new FileInfo(x.ToString());
                FileInfo yInfo = new FileInfo(y.ToString());

                return xInfo.FullName.CompareTo(yInfo.FullName);
            }
        }

        public class SizeSorter : IComparer
        {
            public int Compare(object x, object y)
            {
                if (x == null && y == null)
                {
                    return 0;
                }
                if (x == null)
                {
                    return -1;
                }
                if (y == null)
                {
                    return 1;
                }
                FileInfo xInfo = new FileInfo(x.ToString());
                FileInfo yInfo = new FileInfo(y.ToString());

                return xInfo.Length.CompareTo(yInfo.Length);
            }
        }

        public class TypeSorter : IComparer
        {
            public int Compare(object x, object y)
            {
                if (x == null && y == null)
                {
                    return 0;
                }
                if (x == null)
                {
                    return -1;
                }
                if (y == null)
                {
                    return 1;
                }
                FileInfo xInfo = new FileInfo(x.ToString());
                FileInfo yInfo = new FileInfo(y.ToString());

                return xInfo.Extension.CompareTo(yInfo.Extension);
            }
        }

        public bool IsReusable
        {
            get { return true; }
        }
    }
 */

//图片扩展名
var fileTypes string = "gif,jpg,jpeg,png,bmp"

//
// @rootDir : 根目录路径，相对路径
// @rootUrl : 根目录URL，可以指定绝对路径，比如 http://www.yoursite.com/attached/
func fileManager(r *http.Request,rootDir,rootUrl string)([]byte,error) {
	var currentPath = ""
	var currentUrl = ""
	var currentDirPath = ""
	var moveupDirPath = ""
	var dirPath string = rootDir

	urlQuery := r.URL.Query()
	var dirName string = urlQuery.Get("dir");

	if len(dirName)!= 0 {
		if dirName == "image" || dirName == "flash" ||
		dirName == "media" || dirName == "file" {
			dirPath += dirName + "/"
			rootUrl += dirName + "/"
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				os.MkdirAll(dirPath, os.ModePerm)
			}
		}else {
			return nil, errors.New("Invalid Directory name")
		}
	}



	//根据path参数，设置各路径和URL
	var path string = urlQuery.Get("path")
	if len(path) == 0 {
		currentPath = dirPath
		currentUrl = rootUrl
		currentDirPath = ""
		moveupDirPath = ""
	}else {
		currentPath = dirPath + path
		currentUrl = rootUrl + path
		currentDirPath = path
		//reg := regexp.MustCompile("(.*?)[^\\/]+\\/$")
		moveupDirPath = currentDirPath[:strings.LastIndex(currentDirPath, "\\")]
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
	dir, err := os.Stat(currentPath);
	if os.IsNotExist(err) || !dir.IsDir() {
		return nil, errors.New("no such directory or file not directory")
	}

	//排序形式，name or size or type
	var order string = strings.ToLower(urlQuery.Get("order"))

	//遍历目录取得文件信息

	var dirList SorterFiles = &SorterFiles{
		files:[]os.FileInfo{},
		sortBy:order,
	}
	var fileList SorterFiles = &SorterFiles{
		files:[]os.FileInfo{},
		sortBy:order,
	}

	files, err := ioutil.ReadDir(currentDirPath)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		if v.IsDir() {
			dirList.files = append(dirList.files, v)
		}else {
			fileList.files = append(fileList.files, v)
		}
	}

	var result = make(map[string]interface{})
	result["moveup_dir_path"] = moveupDirPath
	result["current_dir_path"] = currentDirPath
	result["current_url"] = currentUrl;
	result["total_count"] = dirList.Len() + fileList.Len()
	var dirFileList = []map[string]interface{}{}
	for i := 0; i < dirList.Len(); i++ {
		hash := make(map[string]interface{})
		fs, _ := ioutil.ReadDir(currentDirPath+"/"+dirList.files[i].Name())
		hash["is_dir"] = true
		hash["has_file"] = len(fs) > 0
		hash["is_photo"] = false
		hash["filetype"] = ""
		hash["filename"] = dirList.files[i].Name()
		hash["datetime"] = dirList.files[i].ModTime().Format("2006-01-02 15:04:05")
		dirFileList = append(dirFileList, hash)
	}

	var fN,ext string
	for i := 0; i < fileList.Len(); i++ {
		hash := make(map[string]interface{})
		fN = fileList.files[i].Name()
		ext = fN[strings.Index(fN,".")+1:]
		hash["is_dir"] = false;
		hash["has_file"] = false;
		hash["filesize"] = fileList.files[i].Size();
		hash["is_photo"] = strings.Index(fileTypes,ext)
		hash["filetype"] = ext
		hash["filename"] = fN
		hash["datetime"] = fileList.files[i].ModTime().Format("2006-01-02 15:04:05")
		dirFileList = append(dirFileList, hash)
	}

	result["file_list"] =  dirFileList
	return json.Marshal(result)
}


var _ mvc.Filter = new(editorC)
type editorC struct{
	*baseC
}

func (this *editorC) File_manager(ctx *web.Context) {
	partnerId := this.GetPartnerId(ctx)
	d, err := fileManager(ctx.Request,
		fmt.Sprintf("./static/uploads/%d/", partnerId),
		fmt.Sprintf("%s/%d/", ctx.App.Storage().GetString(variable.StaticServer), partnerId),
	)
	ctx.Response.Header().Add("Content-Type","application/json")
	if err != nil {
		ctx.Response.Write([]byte("{error:'"+strings.Replace(err.Error(), "'", "\\'", -1)+"'}"))
	}else {
		ctx.Response.Write(d)
	}
}



