package kdniao

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//物流状态：2-在途中,3-签收,4-问题件

var (
	EBusinessID = ""
	AppKey      = ""
	//ReqUrl      = "http://api.kdniao.cc/Ebusiness/EbusinessOrderHandle.aspx"
	ReqUrl = "http://api.kdniao.com/api/dist"
	SHIP_CODE_YUNDA    = "YD"
	SHIP_CODE_SHUNFENG = "SF"
	// You can add more, get code from https://view.officeapps.live.com/op/view.aspx?src=http://www.kdniao.com/file/ExpressCode.xls
)

type RequestData struct {
	ShipperCode  string
	LogisticCode string
}

type PostParams struct {
	RequestData string
	EBusinessID string
	RequestType string
	DataSign    string
	DataType    string
}

type TraceItem struct {
	AcceptTime    string
	AcceptStation string
}

type TraceResult struct {
	Traces  []TraceItem
	Success bool
	State   string
}

//KdnTraces(SHIP_CODE_YUNDA, "xxxx")

func KdnTraces(shipperCode string, logisticCode string) (traceResult *TraceResult, err error) {
	if AppKey == "" || EBusinessID == "" {
		fmt.Println("Please fill AppKey & EBusinessID")
		return nil, nil
	}
	if requestDataJson, err := json.Marshal(&RequestData{
		ShipperCode:  shipperCode,
		LogisticCode: logisticCode,
	}); err != nil {
		fmt.Printf("KdnTraces request to json error:%v\n", err)
		return nil, err
	} else {
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(string(requestDataJson) + AppKey))
		b64 := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(md5Ctx.Sum(nil))))
		form := url.Values{
			"RequestData": {url.QueryEscape(string(requestDataJson))},
			"EBusinessID": {EBusinessID},
			"RequestType": {"1002"},
			"DataSign":    {url.QueryEscape(b64)},
			"DataType":    {"2"},
		}
		resp, err := http.PostForm(ReqUrl, form)
		if err != nil {
			//fmt.Printf("KdnTraces post error:%v\n", err)
			return nil, err
		} else {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil && err != io.EOF && strings.Index(err.Error(), "EOF") == -1 {
				fmt.Printf("KdnTraces read body error:%v\n", err)
				return nil, err
			}
			// Parser body
			traceResult := TraceResult{}
			json.Unmarshal(body, &traceResult)
			//fmt.Printf("Trace result:%v\n", traceResult)
			return &traceResult, nil
		}
	}
}
