package web

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var (
	DefaultHttpExceptHandle func(http.ResponseWriter, *http.Request, error)
	HttpBeforePrintHandle   func(http.ResponseWriter, *http.Request) bool
	HttpAfterPrintHandle    func(http.ResponseWriter, *http.Request)
)

//Http请求处理代理
type HttpHandleProxy struct {
	//请求之前发生
	//返回false,则终止运行
	Before func(http.ResponseWriter, *http.Request) bool
	After  func(http.ResponseWriter, *http.Request)
	Except func(http.ResponseWriter, *http.Request, error)
}

type ResponseProxyWriter struct {
	writer http.ResponseWriter
	Output []byte
}

func (this *ResponseProxyWriter) Header() http.Header {
	return this.writer.Header()
}
func (this *ResponseProxyWriter) Write(bytes []byte) (int, error) {
	this.Output = append(this.Output, bytes[0:len(bytes)]...)
	return this.writer.Write(bytes)
}
func (this *ResponseProxyWriter) WriteHeader(i int) {
	this.writer.WriteHeader(i)
}

//创建一个新的HttpWriter
func NewRespProxyWriter(w http.ResponseWriter) *ResponseProxyWriter {
	return &ResponseProxyWriter{
		writer: w,
		Output: []byte{},
	}
}

//应用代理
func (this *HttpHandleProxy) For(handle func(http.ResponseWriter,
	*http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		//todo: panic可以抛出任意对象，所以recover()返回一个interface{}
		if this.Except != nil {
			defer func() {
				if err := recover(); err != nil {
					this.Except(w, r, errors.New(fmt.Sprintf("%s", err)))
				}
			}()
		}

		if this.Before != nil {
			if this.Before(w, r) {
				return
			}
		}

		proxy := NewRespProxyWriter(w)

		if handle != nil {
			handle(proxy, r)
		}

		if this.After != nil {
			this.After(proxy, r)
		}
	}
}

func init() {
	DefaultHttpExceptHandle = func(w http.ResponseWriter, r *http.Request, err error) {
		_, f, line, _ := runtime.Caller(1)
		var header http.Header = w.Header()
		header.Add("Content-Type", "text/html")
		w.WriteHeader(500)
		stack := strings.Replace(string(debug.Stack()), "\n", "<br />", -1)
		w.Write([]byte(fmt.Sprintf(`<h1 style="color:red;font-size:20px">ERROR :%s</h1>
				Source:%s line:%d<br />
				</strong><br /><br /><b>Statck:</b><br />%s`,
			err.Error(), f, line, stack)))
		fmt.Fprint(w, err)
	}

	HttpBeforePrintHandle = func(w http.ResponseWriter, r *http.Request) bool {
		fmt.Println("[Request] ", time.Now().Format("2006-01-02 15:04:05"), ": URL:", r.RequestURI)
		for k, v := range r.Header {
			fmt.Println(k, ":", v)
		}
		if r.Method == "POST" {
			r.ParseForm()
		}
		for k, v := range r.Form {
			fmt.Println("form", k, ":", v)
		}
		return true
	}

	HttpAfterPrintHandle = func(w http.ResponseWriter, r *http.Request) {
		proxy, ok := w.(*ResponseProxyWriter)
		if !ok {
			fmt.Println("[Response] convert error")
			return
		}
		fmt.Println("[Respose]\n" + string(proxy.Output))
	}
}
