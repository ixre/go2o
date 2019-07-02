/**
 * Copyright 2015 @ to2.net.
 * name : 253com
 * author : jarryliu
 * date : 2016-07-06 20:50
 * description :
 * history :
 */
package cl253

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//https://www.253.com/api-docs-10.html
// 签名要在后台添加，并在短信中包含。
const url = "http://sms.253.com/msg/send"

func SendMsgToMobile(account, pwd, phone, content string) error {
	strUrl := fmt.Sprintf("%s?un=%s&pw=%s&phone=%s&msg=%s&rd=1",
		url, account, pwd, phone, content)
	rsp, err := http.Get(strUrl)
	if err == nil {
		defer rsp.Body.Close()
		if rsp.StatusCode != http.StatusOK {
			err = errors.New("error : " + strconv.Itoa(rsp.StatusCode))
		}
		var data []byte
		data, err = ioutil.ReadAll(rsp.Body)
		if err == nil {
			arr := strings.Split(string(data), ",")
			if arr[1][0] != '0' {
				err = errors.New("status code : " + arr[1] +
					" ; response : " + string(data))
			}
		}
	}
	return err
}
