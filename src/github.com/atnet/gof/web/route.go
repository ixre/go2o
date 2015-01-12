package web

import (
	"net/http"
	"regexp"
)

//路由映射
type RouteMap struct {
	//地址模式
	UrlPatterns []string
	//路由集合
	RouteCollection map[string]func(http.ResponseWriter, *http.Request)
}

//添加路由
func (this *RouteMap) Add(
	urlPattern string,
	requestFunc func(http.ResponseWriter, *http.Request)) {
	if this.RouteCollection == nil {
		this.RouteCollection =
			make(map[string]func(http.ResponseWriter, *http.Request))
	}
	_, exists := this.RouteCollection[urlPattern]

	if !exists {
		this.RouteCollection[urlPattern] = requestFunc
		this.UrlPatterns = append(this.UrlPatterns, urlPattern)
	}
}

//处理请求
func (this *RouteMap) HandleRequest(w http.ResponseWriter, r *http.Request) {
	routes := this.RouteCollection
	path := r.URL.Path
	var isHandled bool = false

	//range 顺序是随机的，参见：http://golanghome.com/post/155
	for _, k := range this.UrlPatterns {
		v, exist := routes[k]
		if exist {
			matched, err := regexp.Match(k, []byte(path))
			//fmt.Println("\n",k,path)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			if matched && v != nil {
				isHandled = true
				v(w, r)
				break
			}
		}
	}

	if !isHandled {
		w.Write([]byte("404 Not found!"))
	}
}

//处理路由请求
func handleMapRoute(
	w http.ResponseWriter,
	r *http.Request,
	routes map[string]func(http.ResponseWriter, *http.Request)) {
	path := r.URL.Path
	var isHandled bool = false

	for k, v := range routes {
		if !isHandled {
			matched, err := regexp.Match(k, []byte(path))

			//			if path == "/a/" {
			//				fmt.Println(k + "==>" + path)
			//				fmt.Println(matched)
			//				fmt.Println(v)
			//			}

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			if matched && v != nil {
				isHandled = true
				v(w, r)
				//fmt.Println("----")
				break
			}
		}
	}

	if !isHandled {
		w.Write([]byte("404 Not found!"))
	}
}
