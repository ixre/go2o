package hfb

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/ixre/gof/api"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types"
	"go2o/core/infrastructure/qpay"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

const ApplyBankAuthFormURL = "https://Pay.Heepay.com/WithholdAuthPay/ApplyBankAuthForm.aspx"

type (
	cardBinRsp struct {
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
)

func NewHfb(s storage.Interface) qpay.QuickPayProvider {
	agentId, _ := s.GetString("registry/key/qp_hfb_agent_id")
	md5Key, _ := s.GetString("registry/key/qp_hfb_md5_key")
	privateKey := "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDF/hqHZZb7r0S5KuuQ1zE4v6BT+irjybOR0mIBbRqUnlBlIK8eayxs7eTazEn7FIFjepvGMxgH/2tC6R7s45KaoQo5Yq9l/rvziyYI50U4SZor1mV24nlCNLbx5BqDBFcGwxOJqwZGVTelBVjDtOsper10rUjhtwDFcLSe82VoPQUt8k9H4zw8+0lC4DsK0JlNtRJNAi380Fmz5JV19+12D2N8Tn9+pqFXzjyvp2EyJ/hS8uHUXZGy3lh7cbeEkFu5sFcKB2RDSs++8Y5vyeXQ6RLqMlEbJIRcRRAeMaCZ2Vn5OATYQKCvTPmITTzKB7NoOvEOC9FO4V6HMjidZzBTAgMBAAECggEADufwi10EnvI1FFO85GyvEfyrT2c4L2oSENpr8nuKUsIQf2yUgo/DCnhmkGps73A9xYWHkMZr+r4qDyGJ6H/Bm86f/G4HkoA5Gj7RoD35IiG4b7B2dxrZ0jgxxchMjqyW+LVbFTRBBq6Hv+7FHgbS5Y6OEOiy4ftrHXI8xvLAIbbEa9k1EVmH2ZvA5iVTBuZGWsEAQMRrIBNpmyB3Lnmo7iK28vpEPLvxADtlr3/1vpwfIPMb2fUYkuMXsCPuxjGxtkiCNhahUyzzwGG8rvszx/JcP/vWwRC7IQQff+YONdGKrJT5VqchJV1oaKbLg9CbU1/xsuLOn2RZP1A3/ssdsQKBgQDrlYhZ8BYSa2l5euKX7r4NFGETD8UGnyJmCGPy22VstJ77vAvffVLkKSzWrZgOlmW8MdRfFUsLfPaolLx56rCtdgS6mwSh4kqz9nKMuQjQbpECJAJtZL4FuMjVKSL/71Kew3/Bc/MNo6uKGxiK54KjxFu4TXWplKHFAI1MPuhdvQKBgQDXJpkFta6XwWbtrBCrgN5+eROA9qP+xC0WF/Ar8jbNJAntoUYXFLkIMt1HJFKAPND71x54G0ZHHpL7LJCP/NiGhY19/4S1oBP79d67HPku9Kbrm1NXKUzafOv2rPXSK7uGR+XSgnnKbs5GicipcqZP3+OGOajb9xxjer0IpU//TwKBgQCfHy8r4FhoNJjXbsMicCV6XCt9XodsA4yOclhgLwSAujcwPUGfwNx+M7mPf01XfQpWZSnW12EK72sDTwNHLdgMMczb5dzpIxnmGC4jEs/7SNM1KPFixkr7PmaYY+K6EAI0LkRafGDM86Hn9IlNOTYqO3TgNaGl2zixAcBuoYb92QKBgFsS7aerFrMKnWVydsQCkyx6WDU5MoZ/yI4XqAUSTPxdiw5aPG88yG6eCWk6COpb1CMnFrDE6uTkHlfQr4kkAQxAsHprlWPE1XDMzXHre9fSnG4TnB3DT9MVGlWbNZu4A3N+L90CekekzBCz9os0Cw64uXlyIvaqDgxWQnrMb6alAoGBAJ44E3SOo9DD5UOk+6swf/YplhqG2sayJruVib+1D2dlWu/+LxJqQZJGI/jtLVO24q7XGdnlA1YXA85DRI9/VUPPOEaLpUI91KWHUaN0Cgcin/O02UR+UWWvtbNEhI8Huk4BDGOPrBxz1tI2Bw1IvkD6u/mKmiExhzCUX/oAAesT"

	pk, err := qpay.ParseRSAPrivateKey(privateKey)
	if err != nil {
		println("[ Go2o][ Warning]: quick pay error ", err.Error())
	}
	//agentId = "1664502"
	//md5Key = "CC08C5E3E69F4E6B85F1DC0B"
	return &hfbImpl{
		agentId:       agentId,
		md5Key:        md5Key,
		version:       "3",
		rsaPrivateKey: pk,
		Cache:         qpay.NewCache(s),
	}
}

var _ qpay.QuickPayProvider = new(hfbImpl)

type hfbImpl struct {
	agentId       string
	md5Key        string
	version       string
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
	*qpay.Cache
}

func (h *hfbImpl) CheckSign(params map[string]string, signType string, sign string) bool {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return h.signParams(values) == sign
}

func (h *hfbImpl) RequestBankSideAuth(nonce string, bankCardNo string, accountName string, idCardNo string, mobile string) (*qpay.BankAuthResult, error) {
	if h.rsaPrivateKey == nil {
		return nil, errors.New("未设置汇付宝私钥或私钥不正确")
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
	// 签名
	sign, err := crypto.Sha1WithRSA(signParams, h.rsaPrivateKey)
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
	encryptData, _ := rsa.EncryptPKCS1v15(rand.Reader,
		&h.rsaPrivateKey.PublicKey, signParams)
	// 拼装Form表单
	formData := map[string]string{
		"agent_id":     h.agentId,
		"sign":         sign,
		"encrypt_data": string(encryptData),
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
	sb.WriteString(ApplyBankAuthFormURL)
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

func (h *hfbImpl) QueryBankAuthByNonceId(id string) (*qpay.BankAuthQueryResponse, error) {
	panic("implement me")
}

func (h *hfbImpl) QueryBankAuth(bankCardNo string) (*qpay.BankAuthQueryResponse, error) {
	panic("implement me")
}

func (h *hfbImpl) DirectPayment(orderNo string, fee int32, subject string, bankToken string, tradeIp string, notifyUrl string, returnUrl string) (*qpay.QPaymentResponse, error) {
	panic("implement me")
}

func (h *hfbImpl) BatchTransfer(batchTradeNo string, batchTradeFee int32, list []*qpay.CardTransferReq, nonce string, tradeIp string, notifyUrl string) (*qpay.BatchTransferResponse, error) {
	panic("implement me")
}

// 签名
func (h *hfbImpl) signParams(mp url.Values) string {
	query := string(api.ParamsToBytes(mp, h.md5Key, false))
	query = strings.ToLower(query)
	return crypto.Md5([]byte(query))
}

// 查询银行卡信息
func (h *hfbImpl) QueryCardBin(bankCardNo string) *qpay.CardBinQueryResult {
	apiUrl := "https://Pay.heepay.com/API/PayTransit/QueryBankCardInfo.aspx"
	mp := url.Values{
		"agent_id":     []string{h.agentId},
		"bank_card_no": []string{bankCardNo},
		"key":          []string{h.md5Key},
		"version":      []string{h.version},
	}
	sign := h.signParams(mp)
	mp["sign"] = []string{sign}
	cli := http.Client{}
	rsp, err := cli.PostForm(apiUrl, mp)
	if err != nil {
		return &qpay.CardBinQueryResult{ErrMsg: err.Error()}
	}
	reader := transform.NewReader(rsp.Body, simplifiedchinese.GBK.NewDecoder())
	body, err := ioutil.ReadAll(reader)
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
		CardType:            types.IntCond(ret.BankCardType == "1", 1, 0),
		RequireBankSideAuth: true,
	}
}
