/**
 * Copyright 2014 @ to2.net.
 * name : geo.go
 * author : jarryliu
 * date : 2013-12-02 21:34
 * description :
 * history :
 */

package tool

//todo: 新浪ip接口，http://int.dpool.sina.com.cn/iplookup/iplookup.php?format=js&ip=110.110.110.110

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	ip138Regex = regexp.MustCompile("<li>本站主数据：\\s*(.+?)\\s*</li>")
)

func GetLocation(ip string) string {
	rsp, err := http.Get("http://www.ip138.com/ips1388.asp?ip=" + ip + "&action=2")
	if err != nil {
		return "未知地区"
	}
	data, _ := ioutil.ReadAll(rsp.Body)
	//out := make([]byte, len(data)*2)
	//trans := simplifiedchinese.GB18030.NewDecoder()
	//n, _, _ := trans.Transform(out, data, true)
	//m := ip138Regex.FindAllSubmatch(out[:n], 1)
	m := ip138Regex.FindAllSubmatch(data, 1)
	if len(m) != 0 {
		addr := string(m[0][1])
		if addr != "保留地址" {
			return addr
		}
	}
	return "本地网络"
}
