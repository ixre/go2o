package hfb

import (
	"bytes"
	"crypto/rsa"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ixre/go2o/core/infrastructure/qpay"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types"
	http2 "github.com/ixre/gof/util/http"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 快捷（银行侧)
// http://dev.heepay.com/docs/#/KJJK?id=%e5%bf%ab%e6%8d%b7%ef%bc%88%e9%93%b6%e8%a1%8c%e4%be%a7%ef%bc%89

// 快捷支付
//

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : hfb.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-05 09:27
 * description :
 * history :
 */

const cardBinQueryURL = "https://Pay.heepay.com/API/PayTransit/QueryBankCardInfo.aspx"
const applyBankAuthFormURL = "https://Pay.Heepay.com/WithholdAuthPay/ApplyBankAuthForm.aspx"
const bankAuthQueryURL = "https://Pay.Heepay.com/WithholdAuthPay/BankAuthQuery.aspx"
const redirectPaymentURL = "https://Pay.Heepay.com/WithholdAuthPay/SubmitPay.aspx"
const paymentQueryURL = "https://query.heepay.com/Payment/Params.aspx"
const smallAmountTransferURL = "https://Pay.heepay.com/API/PayTransit/PayTransferWithSmallAll.aspx"
const largeAmountTransferURL = "https://Pay.heepay.com/API/PayTransit/PayTransferWithLargeWork.aspx"

type cardBinRsp struct {
	// 查询状态码
	RetCode string `xml:"ret_code"`
	// 查询返回信息，成功为空
	RetMsg string `xml:"ret_msg"`
	// 商户编号，（汇付宝商户内码：七位整数数字）
	AgentId string `xml:"agent_id,omitempty"`
	// 	所查询银行卡号，查询成功时返回)
	BankCardNo string `xml:"bank_card_no,omitempty"`
	// 所属银行名称，查询成功时返回
	BankName string `xml:"bank_name,omitempty"`
	// 银行我方对应编号，查询成功时返回
	BankType int `xml:"bank_type,omitempty"`
	// 银行卡类型（0=储蓄卡,1=信用卡），查询成功时返回
	BankCardType string `xml:"bank_card_type,omitempty"`
}

// 接口返回的加密内容
type EncryptResponse struct {
	// 查询状态码
	RetCode string `xml:"ret_code"`
	// 查询返回信息，成功为空
	RetMsg string `xml:"ret_msg"`
	// 加密后的内容
	EncryptData string `xml:"encrypt_data"`
}

func NewHfb(s storage.Interface) qpay.QuickPayProvider {
	agentId, _ := s.GetString("registry/key/qp_hfb_agent_id")
	md5Key, _ := s.GetString("registry/key/qp_hfb_md5_key")
	privateKey, _ := s.GetString("registry/key/qp_hfb_private_key")
	publicKey, _ := s.GetString("registry/key/qp_hfb_public_key")
	queryMd5Key, _ := s.GetString("registry/key/qp_hfb_query_md5_key")
	batch3DesKey, _ := s.GetString("registry/key/qp_hfb_batch_3des_key")
	//agentId = "1664502"
	//md5Key = "CC08C5E3E69F4E6B85F1DC0B"
	//privateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDF/hqHZZb7r0S5KuuQ1zE4v6BT+irjybOR0mIBbRqUnlBlIK8eayxs7eTazEn7FIFjepvGMxgH/2tC6R7s45KaoQo5Yq9l/rvziyYI50U4SZor1mV24nlCNLbx5BqDBFcGwxOJqwZGVTelBVjDtOsper10rUjhtwDFcLSe82VoPQUt8k9H4zw8+0lC4DsK0JlNtRJNAi380Fmz5JV19+12D2N8Tn9+pqFXzjyvp2EyJ/hS8uHUXZGy3lh7cbeEkFu5sFcKB2RDSs++8Y5vyeXQ6RLqMlEbJIRcRRAeMaCZ2Vn5OATYQKCvTPmITTzKB7NoOvEOC9FO4V6HMjidZzBTAgMBAAECggEADufwi10EnvI1FFO85GyvEfyrT2c4L2oSENpr8nuKUsIQf2yUgo/DCnhmkGps73A9xYWHkMZr+r4qDyGJ6H/Bm86f/G4HkoA5Gj7RoD35IiG4b7B2dxrZ0jgxxchMjqyW+LVbFTRBBq6Hv+7FHgbS5Y6OEOiy4ftrHXI8xvLAIbbEa9k1EVmH2ZvA5iVTBuZGWsEAQMRrIBNpmyB3Lnmo7iK28vpEPLvxADtlr3/1vpwfIPMb2fUYkuMXsCPuxjGxtkiCNhahUyzzwGG8rvszx/JcP/vWwRC7IQQff+YONdGKrJT5VqchJV1oaKbLg9CbU1/xsuLOn2RZP1A3/ssdsQKBgQDrlYhZ8BYSa2l5euKX7r4NFGETD8UGnyJmCGPy22VstJ77vAvffVLkKSzWrZgOlmW8MdRfFUsLfPaolLx56rCtdgS6mwSh4kqz9nKMuQjQbpECJAJtZL4FuMjVKSL/71Kew3/Bc/MNo6uKGxiK54KjxFu4TXWplKHFAI1MPuhdvQKBgQDXJpkFta6XwWbtrBCrgN5+eROA9qP+xC0WF/Ar8jbNJAntoUYXFLkIMt1HJFKAPND71x54G0ZHHpL7LJCP/NiGhY19/4S1oBP79d67HPku9Kbrm1NXKUzafOv2rPXSK7uGR+XSgnnKbs5GicipcqZP3+OGOajb9xxjer0IpU//TwKBgQCfHy8r4FhoNJjXbsMicCV6XCt9XodsA4yOclhgLwSAujcwPUGfwNx+M7mPf01XfQpWZSnW12EK72sDTwNHLdgMMczb5dzpIxnmGC4jEs/7SNM1KPFixkr7PmaYY+K6EAI0LkRafGDM86Hn9IlNOTYqO3TgNaGl2zixAcBuoYb92QKBgFsS7aerFrMKnWVydsQCkyx6WDU5MoZ/yI4XqAUSTPxdiw5aPG88yG6eCWk6COpb1CMnFrDE6uTkHlfQr4kkAQxAsHprlWPE1XDMzXHre9fSnG4TnB3DT9MVGlWbNZu4A3N+L90CekekzBCz9os0Cw64uXlyIvaqDgxWQnrMb6alAoGBAJ44E3SOo9DD5UOk+6swf/YplhqG2sayJruVib+1D2dlWu/+LxJqQZJGI/jtLVO24q7XGdnlA1YXA85DRI9/VUPPOEaLpUI91KWHUaN0Cgcin/O02UR+UWWvtbNEhI8Huk4BDGOPrBxz1tI2Bw1IvkD6u/mKmiExhzCUX/oAAesT"
	//publicKey2 = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsVR6LGVO7kbIBKKuAljjPS+V46Ij8+GVCIhIdx5Nj4kJsByM+wo7Nu8QiZczZsR/Yl9n0hYdb1blAO+O0sA4Dg2ALMJeYamxDe5acC+N5W1aVSiOsqiMmKIX7nOSYL2bPLx6uMG/VZjogZBoqHY5qTQH5AX4nQeqW3rAQACKljuqFTl580+TSZqv+QHcCKQqNDmmFW31a1icELoPWhZF7f+Ry1wr7Q4W1ScpLX3uZZadqsZtH7rvvk+SjxV3y5iCD8ZKFqRdxbuuXXcw+GEth6t0kp5EALkdmJFtIq4uI3lgyqCB+PJq4tyBDZOsU4tY/PqZJ+EbbrPRacRf7ecX0wIDAQAB"
	h := &hfbImpl{
		agentId:              agentId,
		md5Key:               md5Key,
		queryMd5Key:          queryMd5Key,
		batch3DesKey:         batch3DesKey,
		cardBinVersion:       "3",
		batchTransferVersion: "3",
		version:              "1",
		Cache:                qpay.NewCache(s),
	}
	// 转换私钥
	if len(privateKey) > 0 && len(publicKey) > 0 {
		if pk, err := crypto.ParsePrivateKey(strings.TrimSpace(privateKey)); err == nil {
			h.rsaPrivateKey = pk
		}
		if pbk, err := crypto.ParsePublicKey(strings.TrimSpace(publicKey)); err == nil {
			h.rsaPublicKey = pbk
		}
	}
	return h
}

var _ qpay.QuickPayProvider = new(hfbImpl)

type hfbImpl struct {
	agentId        string
	md5Key         string // (2)查询CardBin使用
	version        string
	cardBinVersion string          // (2)查询CardBin接口版本号
	rsaPrivateKey  *rsa.PrivateKey // (1)银行侧快捷支付使用
	rsaPublicKey   *rsa.PublicKey  // (1)银行侧快捷支付使用
	queryMd5Key    string          // (3)查询支付用KEY
	batch3DesKey   string          // (4)批付3DES密钥,加密和解密都使用这个KEY
	*qpay.Cache
	batchTransferVersion string // (3)批付接口版本号
}

func (h *hfbImpl) CheckSign(params map[string]string, _ string, sign string) bool {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return h.signParams(values) == sign
}

func (h *hfbImpl) checkPrivateKey() error {
	if h.rsaPrivateKey == nil || h.rsaPublicKey == nil {
		return errors.New("未设置或设置了不正确的汇付宝公钥或私钥")
	}
	return nil
}

// 申请银行侧认证授权(某些银行需跳转到银行页面进行授权)
func (h *hfbImpl) RequestBankSideAuth(nonce string, bankCardNo string, accountName string, idCardNo string, mobile string) (*qpay.BankAuthResult, error) {
	if err := h.checkPrivateKey(); err != nil {
		return nil, err
	}
	data := h.Cache.GetBankAuthData(nonce)
	if data != nil {
		return nil, errors.New("请勿重复提交")
	}
	// 查询CardBin信息
	bin := h.QueryCardBin(bankCardNo)
	if bin.ErrMsg != "" {
		return nil, errors.New(bin.ErrMsg)
	}
	// 快捷请求只有3个参数agent_id、encrypt_data、和sign。
	// 除去agent_id商户id和RSA 签名所得sign，其他请求参数进行
	// RSA加密得出的值就是encrypt_data。
	signParams := []byte(fmt.Sprintf("bank_card_no=%s&bank_user=%s&cert_no=%s"+
		"&mobile=%s&version=%s", bankCardNo, accountName, idCardNo, mobile, h.version))
	sign, err := crypto.Sha1WithRSA(h.rsaPrivateKey, signParams)
	if err != nil {
		return nil, errors.New("签名失败:" + err.Error())
	}
	// 存储认证请求数据
	data = &qpay.BankAuthSwapData{
		BankCardNo:  bankCardNo,
		AccountName: accountName,
		IdCardNo:    idCardNo,
		Mobile:      mobile,
		BankName:    bin.BankName,
		BankCode:    bin.BankCode,
		CardType:    bin.CardType,
	}
	h.Cache.SaveBankAuthData(nonce, data, 3600)
	// 加密请求数据
	encryptData, _ := crypto.EncryptRSAToBase64(
		h.rsaPublicKey, signParams)

	// 拼装Form表单
	formData := map[string]string{
		"agent_id":     h.agentId,
		"sign":         sign,
		"encrypt_data": encryptData,
	}
	return &qpay.BankAuthResult{
		NonceId:  nonce,
		AuthForm: h.buildAuthForm(formData),
		AuthData: formData,
	}, nil
}

// 拼装Form表单
func (h *hfbImpl) buildAuthForm(formData map[string]string) string {
	sb := bytes.NewBufferString("<form action=\"")
	sb.WriteString(applyBankAuthFormURL)
	sb.WriteString("\" method=\"POST\">")
	for k, v := range formData {
		sb.WriteString("<input type=\"hidden\" name=\"")
		sb.WriteString(k)
		sb.WriteString("\" value=\"")
		sb.WriteString(v)
		sb.WriteString("\"/>")
	}
	sb.WriteString("</form><small>正在提交..</small>")
	sb.WriteString("<script type=\"text/javascript\">document.forms[0].submit()</script>")
	authForm := sb.String()
	return authForm
}

func (h *hfbImpl) QueryBankAuthByNonceId(nonce string) (*qpay.BankAuthQueryResponse, error) {
	if err := h.checkPrivateKey(); err != nil {
		return nil, err
	}
	data := h.Cache.GetBankAuthData(nonce)
	if data == nil {
		return nil, errors.New("授权操作超时")
	}
	return h.QueryBankAuth(data.BankCardNo)
}

func (h *hfbImpl) QueryBankAuth(bankCardNo string) (*qpay.BankAuthQueryResponse, error) {
	if err := h.checkPrivateKey(); err != nil {
		return nil, err
	}
	signParams := []byte(fmt.Sprintf("bank_card_no=%s&version=%s", bankCardNo, h.version))
	sign, err := crypto.Sha1WithRSA(h.rsaPrivateKey, signParams)
	if err != nil {
		return nil, errors.New("签名失败:" + err.Error())
	}
	// 加密请求数据
	encryptData, _ := crypto.EncryptRSAToBase64(h.rsaPublicKey, signParams)
	// 拼装Form表单
	formData := map[string]string{
		"agent_id":     h.agentId,
		"sign":         sign,
		"encrypt_data": encryptData,
	}
	body, err := h.request(bankAuthQueryURL, http2.ParseUrlValues(formData), false)
	if err != nil {
		return nil, err
	}
	ret, mp, err := h.readEncryptResponse(body)
	if err != nil {
		return nil, err
	}
	// 未成功
	if ret.RetCode != "0000" {
		return nil, errors.New(ret.RetMsg)
	}
	// 处理中
	if ret.RetMsg == "处理中" {
		return &qpay.BankAuthQueryResponse{
			Code:          1,
			Message:       "正在授权",
			BankAuthToken: "",
		}, nil
	}
	// 已授权,获取授权码
	return &qpay.BankAuthQueryResponse{
		Code:          0,
		Message:       "已授权",
		BankAuthToken: mp["hy_auth_uid"],
	}, nil
}

func (h *hfbImpl) DirectPayment(orderNo string, fee int32, subject string, bankToken string, tradeIp string, notifyUrl string, returnUrl string) (*qpay.QPaymentResponse, error) {
	if err := h.checkPrivateKey(); err != nil {
		return nil, err
	}
	params := map[string]string{
		"agent_bill_id":   orderNo,
		"agent_bill_time": time.Now().Format("20060102150405"),
		"goods_name":      subject,
		"hy_auth_uid":     bankToken,
		"notify_url":      notifyUrl,
		"pay_amt":         types.FixedMoney(float64(fee) / 100),
		"return_url":      returnUrl,
		"user_ip":         tradeIp,
		"version":         h.version,
	}
	signParams := []byte(http2.SortedQuery(http2.ParseUrlValues(params)))
	sign, err := crypto.Sha1WithRSA(h.rsaPrivateKey, signParams)
	if err != nil {
		return nil, errors.New("签名失败:" + err.Error())
	}
	// 加密请求数据
	encryptData, _ := crypto.EncryptRSAToBase64(h.rsaPublicKey, signParams)
	// 拼装Form表单
	formData := map[string]string{
		"agent_id":     h.agentId,
		"sign":         sign,
		"encrypt_data": encryptData,
	}
	body, err := h.request(redirectPaymentURL, http2.ParseUrlValues(formData), false)
	if err != nil {
		return nil, err
	}
	ret, mp, err := h.readEncryptResponse(body)
	if err != nil {
		return nil, err
	}
	qp := &qpay.QPaymentResponse{
		Code:   ret.RetCode,
		BillNo: "",
	}
	// 未成功
	if ret.RetCode != "0000" {
		return qp, errors.New(ret.RetMsg)
	}
	// 支付成功
	qp.Code = "0"
	qp.BillNo = mp["hy_bill_no"]
	return qp, nil
}

func (h *hfbImpl) QueryPaymentStatus(orderNo string, options map[string]string) (*qpay.QPaymentQueryResponse, error) {
	billTime := options["agent_bill_time"]
	returnMode := "1"
	// 注：此接口的签名参数没有按照字符排序
	query := fmt.Sprintf("version=%s&agent_id=%s&agent_bill_id=%s"+
		"&agent_bill_time=%s&return_mode=%s&key=%s",
		h.version, h.agentId, orderNo, billTime, returnMode, h.queryMd5Key)
	sign := crypto.Md5([]byte(query))
	mp := url.Values{
		"agent_id":        {h.agentId},
		"agent_bill_id":   {orderNo},
		"agent_bill_time": {billTime},
		"key":             {h.queryMd5Key},
		"return_mode":     {returnMode},
		"version":         {h.version},
		"remark":          {"query"},
		"sign":            {sign},
	}
	body, err := h.request(paymentQueryURL, mp, false)
	if err != nil {
		return nil, err
	}
	// 返回以"|"分割的数据
	qp := &qpay.QPaymentQueryResponse{}
	re := regexp.MustCompile("(\\w+)=([^|]*)")
	arr := re.FindAllStringSubmatch(string(body), -1)
	for _, v := range arr {
		switch v[1] {
		case "result":
			qp.Code = types.ElseInt(v[2] == "1", 1, 0)
		case "third_bill_no":
			qp.BillNo = v[2]
		case "detail_error_message":
			qp.ErrMsg = v[2]
		}
	}
	return qp, nil
}

// 请求接口
func (h *hfbImpl) request(url string, mp url.Values, gbk bool) ([]byte, error) {
	cli := http.Client{}
	rsp, err := cli.PostForm(url, mp)
	if err == nil {
		var reader io.Reader
		if gbk {
			reader = transform.NewReader(rsp.Body, simplifiedchinese.GBK.NewDecoder())
		} else {
			reader = rsp.Body
		}
		body, err := io.ReadAll(reader)
		if err == nil {
			return body, nil
		}
	}
	return []byte(""), nil
}

func (h *hfbImpl) BatchTransfer(batchTradeNo string, list []*qpay.CardTransferReq, nonce string, notifyUrl string) (*qpay.BatchTransferResponse, error) {
	if len(h.batch3DesKey) == 0 {
		return nil, errors.New("未设置汇付宝3DES密钥")
	}
	detailData, totalFee := h.batchDetailData(batchTradeNo, list)
	values := url.Values{
		"agent_id":    {h.agentId},
		"batch_amt":   {types.Money(totalFee)},
		"batch_no":    {batchTradeNo},
		"batch_num":   {strconv.Itoa(len(list))},
		"detail_data": {detailData},
		"ext_param1":  {url.QueryEscape(nonce)},
		"key":         {h.md5Key},
		"notify_url":  {notifyUrl},
		"version":     {h.batchTransferVersion},
	}
	sign := h.signParams(values)
	values["sign"] = []string{sign}
	// 批付信息先转换为GBK编码,再用3DES加密
	gbkBytes, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(detailData))
	bytes, err := crypto.EncryptECB3DES(gbkBytes, []byte(h.batch3DesKey))
	if err != nil {
		return nil, err
	}
	detailData = strings.ToUpper(hex.EncodeToString(bytes))
	delete(values, "key")
	values["detail_data"] = []string{detailData}
	// 大于5W使用大额交易接口
	reqURL := types.StringCond(totalFee > 50000, largeAmountTransferURL, smallAmountTransferURL)
	body, err := h.request(reqURL, values, true)
	if err != nil {
		return nil, err
	}
	println(string(body))
	var ret EncryptResponse
	err = xml.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	if ret.RetCode != "0000" {
		return &qpay.BatchTransferResponse{
			Code:    "0",
			NonceId: nonce,
		}, errors.New(ret.RetMsg)
	}
	return &qpay.BatchTransferResponse{
		Code:    "0",
		NonceId: nonce,
	}, nil
}

// 生成批量付款的数据
func (h *hfbImpl) batchDetailData(batchTradeNo string, list []*qpay.CardTransferReq) (string, float64) {
	detailData := bytes.NewBuffer(nil)
	total := 0.00
	for i, v := range list {
		if i > 0 {
			detailData.WriteString("|")
		}
		isCompany := types.StringCond(v.PersonTransfer, "0", "1")
		fee := float64(v.TransactionFee) / 100
		total += fee
		arr := []string{
			batchTradeNo,
			v.BankCode,
			isCompany,
			v.BankCardNo,
			v.BankAccountName,
			types.Money(fee),
			types.StringCond(v.Subject == "", "上游结算款", v.Subject),
			v.Province,
			v.City,
			v.StoreName,
		}
		detailData.WriteString(strings.Join(arr, "^"))
	}
	return detailData.String(), total
}

// 查询银行卡信息
func (h *hfbImpl) QueryCardBin(bankCardNo string) *qpay.CardBinQueryResult {
	mp := url.Values{
		"agent_id":     {h.agentId},
		"bank_card_no": {bankCardNo},
		"key":          {h.md5Key},
		"version":      {h.cardBinVersion},
	}
	sign := h.signParams(mp)
	mp["sign"] = []string{sign}
	body, err := h.request(cardBinQueryURL, mp, true)
	if err != nil {
		return &qpay.CardBinQueryResult{ErrMsg: err.Error()}
	}
	var ret cardBinRsp
	err = xml.Unmarshal(body, &ret)
	if err != nil {
		return &qpay.CardBinQueryResult{ErrMsg: err.Error()}
	}
	//	String b = "中国工商银行，中国光大银行，中国浦发银行，中国银行，中国浙商银行，中国建设银行，中国中信银行，中国华夏银行，中国平安银行，中国上海银行，中国民生银行";
	//	if (!b.contains(rootElm.element("BankName").getText())) {
	//		map.put("success", false);
	//		//map.put("err_msg","暂时不支持的银行卡")
	//	}
	//	return R.ok(map);
	//}
	return &qpay.CardBinQueryResult{
		ErrMsg:              ret.RetMsg,
		BankName:            ret.BankName,
		BankCardNo:          ret.BankCardNo,
		BankCode:            strconv.Itoa(ret.BankType),
		CardType:            types.ElseInt(ret.BankCardType == "1", 1, 0),
		RequireBankSideAuth: true,
	}
}

// 签名
func (h *hfbImpl) signParams(mp url.Values) string {
	query := http2.SortedQuery(mp)
	println(query)
	query = strings.ToLower(query)
	return crypto.Md5([]byte(query))
}

// 读取响应并解密加密内容
func (h *hfbImpl) readEncryptResponse(rsp []byte) (*EncryptResponse, map[string]string, error) {
	var ret EncryptResponse
	err := xml.Unmarshal(rsp, &ret)
	if err != nil {
		return nil, nil, err
	}
	retBytes, _ := crypto.DecryptRSAFromBase64(h.rsaPrivateKey, ret.EncryptData)
	retMsg, _ := http2.ParseQuery(string(retBytes))
	return &ret, retMsg, nil
}
