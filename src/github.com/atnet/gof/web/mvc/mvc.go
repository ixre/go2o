package mvc

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var (
	actionRegexp = regexp.MustCompile("/([^/]+)$")
)

//控制器处理
//@controller ： 包含多种动作，URL中的文件名自动映射到控制器的函数
//				 注意，是区分大小写的,默认映射到index函数
//				 如果是POST请求将映射到控制器“函数名+_post”的函数执行
// @re_post : 是否为post请求额外加上_post来区分Post和Get请求
func HandleRequest(controller interface{}, w http.ResponseWriter, r *http.Request, re_post bool, args ...interface{}) {
	//	defer func(){
	//		if err := recover(); err!= nil{
	//			if _,ok := err.(error);ok{
	//				w.Write([]byte("Inject error, plese check parameter type or index!"))
	//			}
	//		}
	//	}()
	var do string
	groups := actionRegexp.FindAllStringSubmatch(r.URL.Path, 1)
	if len(groups) == 0 || len(groups[0]) == 0 {
		do = "Index"
	} else {
		do = groups[0][1]

		//去扩展名
		extIndex := strings.Index(do, ".")
		if extIndex != -1 {
			do = do[0:extIndex]
		}

		//将第一个字符转为大写,这样才可以
		upperFirstLetter := strings.ToUpper(do[0:1])
		if upperFirstLetter != do[0:1] {
			do = upperFirstLetter + do[1:]
		}
	}

	if re_post && r.Method == "POST" {
		do += "_post"
	}

	t := reflect.ValueOf(controller)
	method := t.MethodByName(do)

	if !method.IsValid() {
		w.Write([]byte("No action named \"" + strings.Replace(do, "_post", "", 1) +
			"\" in " + reflect.TypeOf(controller).String() + "."))
	} else {
		//包含基础的ResponseWriter和http.Request 2个参数
		argsLen := len(args)
		numIn := method.Type().NumIn()

		if argsLen == 0 || numIn == 2 {
			params := []reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)}
			method.Call(params)
		} else {
			params := make([]reflect.Value, numIn)
			params[0] = reflect.ValueOf(w)
			params[1] = reflect.ValueOf(r)

			min := numIn - 2
			if min > argsLen {
				min = argsLen
			}
			for i := 0; i < min; i++ {
				params[i+2] = reflect.ValueOf(args[i])
			}
			method.Call(params)
		}
	}
}
