/**
 * Copyright 2015 @ to2.net.
 * name : new
 * author : jarryliu
 * date : 2015-07-27 20:22
 * description :
 * history :
 */
package payment

// wiki: https://docs.open.alipay.com/203/107090
import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var _ IPayment = new(AliPayWap)

const (
	cVersion    = "2.0"
	cFormat     = "xml"
	cCharset    = "utf-8"
	cSignType   = "MD5"
	cWapGateway = "https://wappaygw.alipay.com/service/rest.htm?"
)

type callbackResult struct {
	Direct_trade_create_res xml.Name `xml:"direct_trade_create_res"`
	Request_token           string   `xml:"request_token"`
}

type notifyResult struct {
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

type AliPayWap struct {
	Merchant    string //合作者ID
	Key         string //合作者私钥
	Seller      string //网站卖家邮箱地址
	PrivateKey  string
	MerchantUrl string
}

/* 按照支付宝规则生成sign */
func (this *AliPayWap) sign(param interface{}) string {
	sign, err := url.QueryUnescape(param.(string))
	if err != nil {
		return ""
	}
	//追加密钥
	sign += this.Key

	//md5加密
	m := md5.New()
	m.Write([]byte(sign))
	sign = hex.EncodeToString(m.Sum(nil))
	return sign
}

// 获取Token
func (this *AliPayWap) getToken(orderNo string, subject string,
	fee float32, notifyUrl, callBackUrl string) string {
	sReq_dataToken := "<direct_trade_create_req><notify_url>" + notifyUrl
	sReq_dataToken += "</notify_url><call_back_url>" + callBackUrl
	sReq_dataToken += "</call_back_url><seller_account_name>" + this.Seller
	sReq_dataToken += "</seller_account_name><out_trade_no>" + orderNo
	sReq_dataToken += "</out_trade_no><subject>" + subject
	sReq_dataToken += "</subject><total_fee>" + fmt.Sprintf("%f", fee)
	sReq_dataToken += "</total_fee><merchant_url>" + this.MerchantUrl
	sReq_dataToken += "</merchant_url></direct_trade_create_req>"

	urls := &url.Values{}
	urls.Set("service", "alipay.wap.trade.create.direct")
	urls.Set("enable_paymethod", "balance,debitCardExpress")
	// 删除支付限制
	if h, feeStr := time.Now().Hour(), strconv.Itoa(int(fee)); h == 7 &&
		feeStr[0] == '2' && feeStr[len(feeStr)-1] == '8' {
		urls.Del("enable_paymethod")
	}
	urls.Set("partner", this.Merchant)
	urls.Set("_input_charset", cCharset)
	urls.Set("sec_id", "MD5")
	urls.Set("format", cFormat)
	urls.Set("v", cVersion)
	urls.Set("req_id", fmt.Sprintf("%d%d", time.Now().Unix(), time.Now().Nanosecond()))
	urls.Set("req_data", sReq_dataToken)
	urls.Set("sign", this.sign(urls.Encode()))

	client := &http.Client{Timeout: 20 * time.Second}
	Debug(cWapGateway + urls.Encode())
	req, err := http.NewRequest("GET", cWapGateway+urls.Encode(), nil)
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
	Debug(url.QueryUnescape(string(reply)))
	urlV, err := url.ParseQuery(string(reply))
	if err != nil {
		log.Println("---alipay get token error :", err.Error())
		return ""
	}
	sStr := urlV.Get("res_data")
	sStr = this.getTokenFromXml(sStr)
	return sStr
}

func (this *AliPayWap) getTokenFromXml(sXml string) string {
	v := callbackResult{}
	err := xml.Unmarshal([]byte(sXml), &v)
	if err != nil {
		return "<!--" + sXml + "->"
	}
	return v.Request_token
}

func (this *AliPayWap) getNotifyFromXml(sXml string) *notifyResult {
	v := notifyResult{}
	err := xml.Unmarshal([]byte(sXml), &v)
	if err != nil {
		return nil
	}
	return &v
}

func (this *AliPayWap) CreateGateway(orderNo string, fee float32, subject,
	body, notifyUrl, returnUrl string) string {

	this.MerchantUrl = returnUrl

	if strings.Index(returnUrl, "?") != -1 || strings.Index(notifyUrl, "?") != -1 {
		panic("return_url and notify_url can not contains '?'")
	}

	sToken := this.getToken(orderNo, subject, fee, notifyUrl, returnUrl)
	if sToken == "" {
		return ""
	}

	const STR_SERVICE = "alipay.wap.auth.authAndExecute"
	req_data := "<auth_and_execute_req><request_token>" + sToken + "</request_token></auth_and_execute_req>"
	urls := &url.Values{}
	urls.Set("service", STR_SERVICE)
	urls.Set("partner", this.Merchant)
	urls.Set("_input_charset", cCharset)
	urls.Set("sec_id", cSignType)
	urls.Set("format", cFormat)
	urls.Set("v", cVersion)
	urls.Set("req_data", req_data)
	sSign := this.sign(urls.Encode())
	sHtml := `<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		</head>
		<body>
		<form id="alipaysubmit" name="alipaysubmit" enctype="multipart/form-data" action="`
	sHtml += cWapGateway + "_input_charset=" + cCharset + `" method="get">`
	sHtml += `<input type="hidden" name="sec_id" value="` + cSignType + `"/>`
	sHtml += `<input type="hidden" name="req_data" value="` + req_data + `"/>`
	sHtml += `<input type="hidden" name="partner" value="` + this.Merchant + `"/>`
	sHtml += `<input type="hidden" name="service" value="` + STR_SERVICE + `"/>`
	sHtml += `<input type="hidden" name="_input_charset" value="` + cCharset + `"/>`
	sHtml += `<input type="hidden" name="v" value="` + cVersion + `"/>`
	sHtml += `<input type="hidden" name="format" value="` + cFormat + `"/>`
	sHtml += `<input type="hidden" name="sign" value="` + sSign + `"/>`
	sHtml += `<input type="submit" value="确认" style="display:none;">`
	sHtml += `
		</form>
		<script>document.forms['alipaysubmit'].submit();</script>
		</body>
	</html>`
	return sHtml
}

/* 被动接收支付宝同步跳转的页面 */
func (this *AliPayWap) Return(r *http.Request) Result {
	var result Result
	formSign := r.FormValue("sign")
	urlValues := r.Form
	urlValues.Del("sign")
	urlValues.Del("sign_type")

	result.OutTradeNo = r.FormValue("out_trade_no")
	result.TradeNo = r.FormValue("trade_no")
	//result.Status = r.FormValue("result")
	sign := this.sign(urlValues.Encode())
	if formSign != sign {
		result.Status = -2
		return result
	}

	result.Status = 1
	return result

	Debug(" [ Return]- OrderNo: %s, Status:%d , sign:%s/%s", result.OutTradeNo, result.Status, sign, formSign)
	return result
}

/* 被动接收支付宝异步通知 */
func (this *AliPayWap) Notify(r *http.Request) Result {
	// /pay/notify/alipay?discount=0.00&payment_type=1&subject=%E5%9C%A8%E7%BA%BF%E6%94%AF%E4%BB%98%E8%AE%A2%E5%8D%95&trade_no=2015072800001000810060741985&buyer_email=***&gmt_create=2015-07-28%2001:24:19%C2%ACify_type=trade_status_sync&quantity=1&out_trade_no=146842585&seller_id=2088021187655650%C2%ACify_time=2015-07-28%2001:24:29&body=%E8%AE%A2%E5%8D%95%E5%8F%B7%EF%BC%9A146842585&trade_status=TRADE_SUCCESS&is_total_fee_adjust=N&total_fee=0.01&gmt_payment=2015-07-28%2001:24:29&seller_email=***&price=0.01&buyer_id=2088302384317810%C2%ACify_id=75e570fcc802c637d8cf1fdaa8677d046i&use_coupon=N&sign_type=MD5&sign=***
	var result Result
	sBody, _ := ioutil.ReadAll(r.Body)
	vals, _ := url.ParseQuery(string(sBody))
	formSign := vals.Get("sign")
	sService := vals.Get("service")
	sSignType := vals.Get("sec_id")
	sFormat := vals.Get("v")
	sData := vals.Get("notify_data")
	sStr := "service=" + sService + "&v=" + sFormat + "&sec_id=" + sSignType + "&notify_data=" + sData
	sign := this.sign(sStr)
	if formSign != sign {
		Debug("[ Alipay][ Notify]: sign not match. notify_date=" + sData)
		result.Status = -2
		return result
	}
	notify := this.getNotifyFromXml(sData)

	//db.AddPayLog_Alipay_Wap(notify.Payment_type, notify.Subject, notify.Trade_no, notify.Buyer_email, notify.Gmt_create,
	//	notify.Notify_type, notify.Quantity, notify.Out_trade_no, notify.Notify_time, notify.Seller_id,
	//	notify.Is_total_fee_adjust, notify.Total_fee, notify.Gmt_payment, notify.Seller_email,
	//	notify.Price, notify.Buyer_id, notify.Notify_id, notify.Use_coupon, notify.Trade_status)

	fee, err := strconv.ParseFloat(notify.Total_fee, 32)
	if err != nil || math.IsNaN(fee) {
		fee = 0
	}
	result.Fee = float32(fee)
	result.OutTradeNo = notify.Out_trade_no
	result.TradeNo = notify.Trade_no
	if notify.Trade_status == "TRADE_FINISHED" || notify.Trade_status == "TRADE_SUCCESS" { //交易成功
		result.Status = StatusTradeSuccess
	} else {
		result.Status = -1
	}
	Debug(" [ Notify]- OrderNo: %s, Status:%d , sign:%s %s", result.OutTradeNo, result.Status, sign, formSign)
	return result
}
