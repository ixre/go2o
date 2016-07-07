/**
 * Copyright 2015 @ z3q.net.
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

const url = "http://222.73.117.169/msg/HttpBatchSendSM"

func SendMsgToMobile(account, pwd, phone, content string) error {
	strUrl := fmt.Sprintf("%s?account=%s&pswd=%s&mobile=%s&msg=%s",
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
			if arr[1] != "0" {
				err = errors.New("status code : " + arr[1])
			}
		}
	}
	return err
}
