package main

import (
	"flag"
	"fmt"
	"github.com/atnet/gof/web"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
)

func main() {
	log.SetOutput(os.Stdout)
	var port *int
	var host *string
	var proxyURL *string
	var proxyedURL string

	port = flag.Int("port", 80, "server port")
	host = flag.String("host", "localhost", "host name by access")
	proxyURL = flag.String("proxy", "", "proxy url")

	flag.Parse()

	if *proxyURL == "" {
		log.Println("[Error]: proxy url example \"http://www.google.com:80\"")
		return
	}

	//host
	if lhost := strings.ToLower(*host); lhost == "localhost" ||
		lhost == "127.0.0.1" {
		log.Println("[Warnning]:only access by localhost,if you want to access other please set --host [yourhost] ")
	} else if strings.Index(lhost, "//") != -1 || lhost[len(lhost)-1:len(lhost)] == "/" {
		log.Println("[Error]: host example \"www.google.com\"")
		return
	}

	//port
	if *port != 80 {
		proxyedURL = fmt.Sprintf("http://%s:%d", *host, *port)
	} else {
		proxyedURL = fmt.Sprintf("http://%s", *host)
	}

	var proxy web.HttpHandleProxy = web.HttpHandleProxy{
		Before: func(w http.ResponseWriter, r *http.Request) bool {
			query := r.URL.Path + "?" + r.URL.RawQuery
			url := *proxyURL + query
			var resp *http.Response

			cookies, err := cookiejar.New(nil)
			if err != nil {
				log.Println("[Error]:", err.Error())
				return false
			}

			client := &http.Client{
				Jar: cookies,
			}

			if r.Method == "GET" {
				resp, _ = client.Get(url)
			} else if r.Method == "POST" {
				r.ParseForm()
				resp, _ = client.PostForm(url, r.Form)
			} else {
				w.Write([]byte("Unsupport request!"))
				w.WriteHeader(500)
				return false
			}

			if resp != nil {
				for i, k := range resp.Header {
					w.Header().Set(i, strings.Join(k, ""))
				}
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {

			} else {
				//fmt.Println(resp.Header.Get("Content-Type"))
				if strings.Index(resp.Header.
					Get("Content-Type"), "text/") == -1 {
					w.Write(data)
				} else {
					source := string(data)
					w.Write([]byte(strings.Replace(source, *proxyURL, proxyedURL, -1)))
				}
			}
			return false
		},
		After: func(http.ResponseWriter, *http.Request) {

		},
		Except: func(w http.ResponseWriter, r *http.Request, err error) {
			w.Write([]byte(fmt.Sprintf("[Error]:%s", err.Error())))
		},
	}

	http.HandleFunc("/", proxy.For(nil))

	log.Println(fmt.Sprintf("Listening on port %d,proxy for url : %s", *port, *proxyURL))

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Println("[Error]:", err.Error())
	}
}
