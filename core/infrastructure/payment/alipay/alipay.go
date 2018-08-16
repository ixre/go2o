/**
 * Copyright 2015 @ z3q.net.
 * name : alipay.go
 * author : jarryliu
 * date : 2015-07-28 17:21
 * description :
 * history :
 */
package alipay

//作者yeyongchang yycmail@163.com

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	s_sAlipayGatewayWap = "https://wappaygw.alipay.com/service/rest.htm?" //alipay支付网关(WAP)
	s_sAlipayGatewayWeb = "https://mapi.alipay.com/gateway.do?"           //alipay支付网关(WEB)
	s_sAlipayPartner    = "xxxx"                                          //合作者ID
	s_sAlipayKey        = "xxxx"                                          //合作者key(md5)
	s_sWebNotifyUrl     = "http://xxxx/alipay_web_notify"                 //网站异步返回地址
	s_sWebCallbackUrl   = "http://xxxx/alipay_web_callback"               //网站同步返回地址
	s_sWapNotifyUrl     = "http://xxxx/alipay_wap_notify"                 //WAP异步返回地址
	s_sWapCallbackUrl   = "http://xxxx/alipay_wap_callback"               //WAP同步返回地址
	s_sWapMerchantUrl   = "http://xxxx/alipay_wap"                        //WAP商户购物网址
	s_sRefundNotifyUrl  = "http://xxxx/alipay_refund_notify"              //退款通知地址
	s_sSellerEmail      = "xxxx@xxxx.com"                                 //网站卖家邮箱地址
	s_sVersion          = "2.0"
	s_sFormat           = "xml"
	s_sCharset          = "utf-8"
	s_sSigntype         = "MD5"
)

var (
	logger   *log.Logger = nil
	s_nCount int
	s_ch     chan int
)

func getUniqueID() string {
	s_ch <- 0 //多线程保护，随便塞点什么，让管道堵住，其实用sync.Mutex.Lock()更快更好
	s_nCount++
	nID := s_nCount % 100
	if s_nCount == 99 {
		s_nCount = 0
	}
	<-s_ch //清空管道
	return fmt.Sprintf("%s%03d", time.Now().Format("20060102150405"), nID)
}

//支付宝请求-------------------
//创建一个交易请求，得到交易token，输入订单好，订单标题，总金额
func getToken(sTrade_no string, sSubject string, sTotalAmount string) string {
	sReq_id := getUniqueID()
	sReq_dataToken := "<direct_trade_create_req><notify_url>" + s_sWapNotifyUrl
	sReq_dataToken += "</notify_url><call_back_url>" + s_sWapCallbackUrl
	sReq_dataToken += "</call_back_url><seller_account_name>" + s_sSellerEmail
	sReq_dataToken += "</seller_account_name><out_trade_no>" + sTrade_no
	sReq_dataToken += "</out_trade_no><subject>" + sSubject
	sReq_dataToken += "</subject><total_fee>" + sTotalAmount
	sReq_dataToken += "</total_fee><merchant_url>" + s_sWapMerchantUrl
	sReq_dataToken += "</merchant_url></direct_trade_create_req>"

	urls := &url.Values{}
	urls.Set("service", "alipay.wap.trade.create.direct")
	urls.Set("partner", s_sAlipayPartner)
	urls.Set("_input_charset", s_sCharset)
	urls.Set("sec_id", s_sSigntype)
	urls.Set("format", s_sFormat)
	urls.Set("v", s_sVersion)
	urls.Set("req_id", sReq_id)
	urls.Set("req_data", sReq_dataToken)
	urls.Set("sign", alipayMd5Sign(urls.Encode()))

	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("GET", s_sAlipayGatewayWap+urls.Encode(), nil)
	if err != nil {
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	urlv, err := url.ParseQuery(string(reply))
	if err != nil {
		return ""
	}
	sStr := urlv.Get("res_data")
	sStr = getTokenFromXml(sStr)
	return sStr
}

type CallbackResult struct {
	Direct_trade_create_res xml.Name `xml:"direct_trade_create_res"`
	Request_token           string   `xml:"request_token"`
}

func getTokenFromXml(sXml string) string {
	v := CallbackResult{}
	err := xml.Unmarshal([]byte(sXml), &v)
	if err != nil {
		return ""
	}
	return v.Request_token
}

//发起Wap支付请求，需要两步
func alipayRequest_Wap(sTrade_no string, sSubject string, sTotalAmount string) string {
	sToken := getToken(sTrade_no, sSubject, sTotalAmount)
	if sToken == "" {
		return ""
	}
	const STR_SERVICE = "alipay.wap.auth.authAndExecute"
	req_data := "<auth_and_execute_req><request_token>" + sToken + "</request_token></auth_and_execute_req>"
	urls := &url.Values{}
	urls.Set("service", STR_SERVICE)
	urls.Set("partner", s_sAlipayPartner)
	urls.Set("_input_charset", s_sCharset)
	urls.Set("sec_id", s_sSigntype)
	urls.Set("format", s_sFormat)
	urls.Set("v", s_sVersion)
	urls.Set("req_data", req_data)
	sSign := alipayMd5Sign(urls.Encode())
	sHtml := `<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		</head>
		<body>
		<form id="alipaysubmit" name="alipaysubmit" enctype="multipart/form-data" action="`
	sHtml += s_sAlipayGatewayWap + "_input_charset=" + s_sCharset + `" method="get">`
	sHtml += `<input type="hidden" name="sec_id" value="` + s_sSigntype + `"/>`
	sHtml += `<input type="hidden" name="req_data" value="` + req_data + `"/>`
	sHtml += `<input type="hidden" name="partner" value="` + s_sAlipayPartner + `"/>`
	sHtml += `<input type="hidden" name="service" value="` + STR_SERVICE + `"/>`
	sHtml += `<input type="hidden" name="_input_charset" value="` + s_sCharset + `"/>`
	sHtml += `<input type="hidden" name="v" value="` + s_sVersion + `"/>`
	sHtml += `<input type="hidden" name="format" value="` + s_sFormat + `"/>`
	sHtml += `<input type="hidden" name="sign" value="` + sSign + `"/>`
	sHtml += `<input type="submit" value="确认" style="display:none;">`
	sHtml += `
		</form>
		<script>document.forms['alipaysubmit'].submit();</script>
		</body>
	</html>`
	return sHtml
}

func AlipayHandler_Wap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	sSubject := r.FormValue("subject")
	sFee := r.FormValue("fee")
	if sSubject == "" {
		sSubject = "测试订单"
	}
	if sFee == "" {
		sFee = "0.01"
	}
	html := alipayRequest_Wap(getUniqueID(), sSubject, sFee)
	w.Write([]byte(html))
}

//发起退款请求
func alipayRefundRequest(sTrade_no string, sDesc string, sFee string) string {
	const STR_SERVICE = "refund_fastpay_by_platform_pwd"
	sDate := time.Now().Format("2006-01-02 15:04:05")
	sBatch_no := getUniqueID()
	sDetail := sTrade_no + "^" + sFee + "^" + sDesc
	urls := &url.Values{}
	urls.Set("service", STR_SERVICE)
	urls.Set("partner", s_sAlipayPartner)
	urls.Set("_input_charset", s_sCharset)
	urls.Set("notify_url", s_sRefundNotifyUrl)
	urls.Set("seller_email", s_sSellerEmail)
	urls.Set("refund_date", sDate)
	urls.Set("batch_no", sBatch_no)
	urls.Set("batch_num", "1")
	urls.Set("detail_data", sDetail)
	sSign := alipayMd5Sign(urls.Encode())
	sHtml := `<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		</head>
		<body>
		<form id="alipaysubmit" name="alipaysubmit" enctype="multipart/form-data" action="`
	sHtml += s_sAlipayGatewayWeb + "_input_charset=" + s_sCharset + `" method="get">`
	sHtml += `<input type="hidden" name="service" value="` + STR_SERVICE + `"/>`
	sHtml += `<input type="hidden" name="partner" value="` + s_sAlipayPartner + `"/>`
	sHtml += `<input type="hidden" name="_input_charset" value="` + s_sCharset + `"/>`
	sHtml += `<input type="hidden" name="notify_url" value="` + s_sRefundNotifyUrl + `"/>`
	sHtml += `<input type="hidden" name="seller_email" value="` + s_sSellerEmail + `"/>`
	sHtml += `<input type="hidden" name="refund_date" value="` + sDate + `"/>`
	sHtml += `<input type="hidden" name="batch_no" value="` + sBatch_no + `"/>`
	sHtml += `<input type="hidden" name="batch_num" value="1"/>`
	sHtml += `<input type="hidden" name="detail_data" value="` + sDetail + `"/>`
	sHtml += `<input type="hidden" name="sign" value="` + sSign + `"/>`
	sHtml += `<input type="hidden" name="sign_type" value="` + s_sSigntype + `"/>`
	sHtml += `<input type="submit" value="确认" style="display:none;">`
	sHtml += `
		</form>
		<script>document.forms['alipaysubmit'].submit();
		</script>
		</body>
	</html>`
	return sHtml
}

//trade_no 交易号 , fee 退还金额, desc 退还说明
func Alipay_Refund_Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	sTradeNo := r.FormValue("trade_no")
	sFee := r.FormValue("fee")
	sDesc := r.FormValue("desc")
	if sTradeNo == "" {
		w.Write([]byte("错误，没有交易单号"))
		return
	}
	if sDesc == "" {
		sDesc = "退款测试"
	}
	if sFee == "" {
		sFee = "0.01"
	}
	html := alipayRefundRequest(sTradeNo, sDesc, sFee)
	w.Write([]byte(html))
}

//支付宝同步回调---------------------
//页面回调，返回信息包括is_success,sign_type,sign,out_trade_no,subject,payment_type,exterface,trade_no,trade_status,
//notify_id,notify_time,notify_type,seller_email,buyer_email,seller_id,buyer_id,total_fee,body,extra_common_param,agent_user_id
func Alipay_Web_CallbackHandler(w http.ResponseWriter, r *http.Request) {
	sSign := r.FormValue("sign")
	vals := r.Form
	vals.Del("sign")
	vals.Del("sign_type")
	w.Header().Set("Content-Type", "text/plain")
	if sSign != alipayMd5Sign(vals.Encode()) {
		w.Write([]byte("通知信息错误，校验失败"))
		return
	}
	sResult := fmt.Sprintf("调用结果:%s 交易结果:%s 订单号:%s 交易号:%s 总金额:%s 信用卡代理:%s", r.FormValue("is_success"), r.FormValue("trade_status"),
		r.FormValue("out_trade_no"), r.FormValue("trade_no"), r.FormValue("total_fee"), r.FormValue("agent_user_id"))
	w.Write([]byte(sResult))
}

//WAP支付回调，sign，result，out_trade_no，trade_no，request_token
func Alipay_Wap_CallbackHandler(w http.ResponseWriter, r *http.Request) {
	sSign := r.FormValue("sign")
	vals := r.Form
	vals.Del("sign")
	vals.Del("sign_type")
	w.Header().Set("Content-Type", "text/plain")
	if sSign != alipayMd5Sign(vals.Encode()) {
		w.Write([]byte("通知信息错误，校验失败"))
		return
	}
	sResult := fmt.Sprintf("交易结果:%s 订单号:%s 交易号:%s", r.FormValue("result"), r.FormValue("out_trade_no"), r.FormValue("trade_no"))
	w.Write([]byte(sResult))
}

//-------------支付宝异步通知------------------
//WEB支付(退款也会产生通知refund_statuss)异步通知， 返回信息包括notify_time,notify_type,notify_id,sign_type,sign
//out_trade_no,subject,payment_type,trade_no,trade_status,gmt_create,gmt_payment,gmt_close,refund_status,gmt_refund,
//seller_email,buyer_email,seller_id,buyer_id,price,total_fee,quantity,body,discount,is_total_fee_adjust,use_coupon,extra_common_param,
//out_channel_type,out_channel_amount,out_channel_inst,business_scene
func Alipay_Web_NotifyHandler(w http.ResponseWriter, r *http.Request) {
	sBody, _ := ioutil.ReadAll(r.Body)
	logger.Printf("即时到账WEB支付异步通知:\n")
	logger.Printf("url:\n%s", r.RequestURI)
	logger.Printf("body:\n%s", string(sBody))
	vals, _ := url.ParseQuery(string(sBody))
	sSign := vals.Get("sign")
	vals.Del("sign")
	vals.Del("sign_type")
	w.Header().Set("Content-Type", "text/plain")
	if sSign != alipayMd5Sign(vals.Encode()) {
		w.Write([]byte("fail"))
		return
	}
	//db.AddPayLog_Alipay_Web(vals.Get("payment_type"), vals.Get("subject"), vals.Get("trade_no"), vals.Get("buyer_email"), vals.Get("gmt_create"),
	//	vals.Get("notify_type"), vals.Get("quantity"), vals.Get("out_trade_no"), vals.Get("notify_time"), vals.Get("seller_id"),
	//	vals.Get("is_total_fee_adjust"), vals.Get("total_fee"), vals.Get("gmt_payment"), vals.Get("seller_email"), vals.Get("price"),
	//	vals.Get("buyer_id"), vals.Get("notify_id"), vals.Get("use_coupon"), vals.Get("out_channel_type"), vals.Get("out_channel_amount"),
	//	vals.Get("out_channel_inst"), vals.Get("business_scene"), vals.Get("trade_status"), vals.Get("refund_status"))
	w.Write([]byte("success"))
}

//手机WAP支付异步通知，签名不排序
//包含service,v,sec_id,sign,notify_data
//notify_data包含:payment_type,subject,trade_no,buyer_email,gmt_create,notify_type,quantity,out_trade_no,notify_time
//seller_id,trade_status,is_total_fee_adjust,total_fee,gmt_payment,seller_email,gmt_close,price,buyer_id,notify_id,use_coupon
func Alipay_Wap_NotifyHandler(w http.ResponseWriter, r *http.Request) {
	sBody, _ := ioutil.ReadAll(r.Body)
	logger.Printf("即时到账WAP支付异步通知:\n")
	logger.Printf("url:\n%s", r.RequestURI)
	logger.Printf("body:\n%s", string(sBody))
	vals, _ := url.ParseQuery(string(sBody))
	sSign := vals.Get("sign")
	sService := vals.Get("service")
	sSignType := vals.Get("sec_id")
	sFormat := vals.Get("v")
	sData := vals.Get("notify_data")
	sStr := "service=" + sService + "&v=" + sFormat + "&sec_id=" + sSignType + "&notify_data=" + sData
	w.Header().Set("Content-Type", "text/plain")
	if sSign != alipayMd5Sign(sStr) {
		w.Write([]byte("fail"))
		return
	}
	notify := getNotifyFromXml(sData)
	logger.Printf("payment_type:%s ...", notify.Payment_type)
	//db.AddPayLog_Alipay_Wap(notify.Payment_type, notify.Subject, notify.Trade_no, notify.Buyer_email, notify.Gmt_create,
	//	notify.Notify_type, notify.Quantity, notify.Out_trade_no, notify.Notify_time, notify.Seller_id,
	//	notify.Is_total_fee_adjust, notify.Total_fee, notify.Gmt_payment, notify.Seller_email,
	//	notify.Price, notify.Buyer_id, notify.Notify_id, notify.Use_coupon, notify.Trade_status)
	w.Write([]byte("success"))
}

type NotifyResult struct {
	Notify              xml.Name `xml:"notify"`
	Payment_type        string   `xml:"payment_type"`
	Subject             string   `xml:"subject"`
	Trade_no            string   `xml:"trade_no"`
	Buyer_email         string   `xml:"buyer_email"`
	Gmt_create          string   `xml:"gmt_create"`
	Notify_type         string   `xml:"notify_type"`
	Quantity            string   `xml:"quantity"`
	Out_trade_no        string   `xml:"out_trade_no"`
	Notify_time         string   `xml:"notify_time"`
	Seller_id           string   `xml:"seller_id"`
	Trade_status        string   `xml:"trade_status"`
	Is_total_fee_adjust string   `xml:"is_total_fee_adjust"`
	Total_fee           string   `xml:"total_fee"`
	Gmt_payment         string   `xml:"gmt_payment"`
	Seller_email        string   `xml:"seller_email"`
	Price               string   `xml:"price"`
	Buyer_id            string   `xml:"buyer_id"`
	Notify_id           string   `xml:"notify_id"`
	Use_coupon          string   `xml:"use_coupon"`
}

func getNotifyFromXml(sXml string) *NotifyResult {
	v := NotifyResult{}
	err := xml.Unmarshal([]byte(sXml), &v)
	if err != nil {
		return nil
	}
	return &v
}

//退款异步通知，包含notify_time,notify_type,notify_id,sign_type,sign,batch_no,success_num,result_details
//result_details包含 交易号^退款金额^处理结果($退费账号^退费账户^ID^退费金额^处理结果)
func Alipay_Refund_NotifyHandler(w http.ResponseWriter, r *http.Request) {
	sBody, _ := ioutil.ReadAll(r.Body)
	logger.Printf("即时到账退款异步通知:\n")
	logger.Printf("url:\n%s", r.RequestURI)
	logger.Printf("body:\n%s", string(sBody))
	vals, _ := url.ParseQuery(string(sBody))
	sSign := vals.Get("sign")
	vals.Del("sign")
	vals.Del("sign_type")
	w.Header().Set("Content-Type", "text/plain")
	if sSign != alipayMd5Sign(vals.Encode()) {
		w.Write([]byte("fail"))
		return
	}
	sDetail := vals.Get("result_details")
	var paras []string
	if sDetail != "" {
		paras = strings.Split(sDetail, "^")
	}
	paras = append(paras, "", "", "") //避免分离失败导致后面的异常
	logger.Printf("退款记录:交易号:%s 金额:%s 结果:%s 批次:%s 笔数:%s", paras[0], paras[1], paras[2], vals.Get("batch_no"), vals.Get("success_num"))
	//db.AddPayLog_Refund(vals.Get("batch_no"), vals.Get("success_num"), paras[0], paras[1], paras[2])
	w.Write([]byte("success"))
}

type AlipayLog struct {
	Pay_type     string `json:"pay_type"`
	Trade_no     string `json:"trade_no"`
	Subject      string `json:"subject"`
	Buyer_email  string `json:"buyer_email"`
	Gmt_create   string `json:"gmt_create"`
	Out_trade_no string `json:"out_trade_no"`
	Total_fee    string `json:"total_fee"`
	Gmt_payment  string `json:"gmt_payment"`
}

type AlipayLogResult struct {
	Errcode int         `json:"errcode"`
	ErrMsg  string      `json:"errmsg,omitempty"`
	Logs    []AlipayLog `json:"logs,omitempty"`
}

//func GetAlipayTracklog_Handler(w http.ResponseWriter, r *http.Process){
//	result		:=AlipayLogResult{}
//	sCount		:=r.FormValue("count")
//	sSort		:=r.FormValue("sort")
//	nCount,_	:=strconv.Atoi(sCount)
//	rows, err 	:=db.Get_AliPay_Log(nCount, sSort =="desc")
//	if err == nil {
//    	defer rows.Close()
//	    for rows.Next() {
//			var pay_type,subject,trade_no,buyer_email,gmt_create,out_trade_no,total_fee,gmt_payment sql.NullString
//			rows.Scan(&pay_type, &subject, &trade_no, &buyer_email, &gmt_create, &out_trade_no, &total_fee, &gmt_payment)
//	        if  pay_type.String !="" {
//				log 			:=AlipayLog{}
//				log.Pay_type	=pay_type.String
//				log.Subject		=subject.String
//				log.Trade_no	=trade_no.String
//				log.Buyer_email	=buyer_email.String
//				log.Gmt_create	=gmt_create.String
//				log.Out_trade_no	=out_trade_no.String
//				log.Total_fee	=total_fee.String
//				log.Gmt_payment	=gmt_payment.String
//				result.Logs		=append(result.Logs, log)
//	        }
//	    }
//		result.Errcode	=0
//    }else{
//		result.Errcode	=global.ERRCODE_SERVER
//		result.Message	=err.Error()
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	json.NewEncoder(w).Encode(&result)
//}
