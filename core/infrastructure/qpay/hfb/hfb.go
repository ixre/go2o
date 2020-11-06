package hfb

import (
	"encoding/xml"
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


type(
	cardBinRsp struct{
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

func NewHfb(s storage.Interface)qpay.QuickPayProvider{
	agentId,_ := s.GetString("registry/key/qp_hfb_agent_id")
	md5Key,_ := s.GetString("registry/key/qp_hfb_md5_key")
	//agentId = "1664502"
	//md5Key = "CC08C5E3E69F4E6B85F1DC0B"
	return &hfbImpl{
		agentId:agentId,
		md5Key: md5Key,
		version: "3",
	}
}

var _ qpay.QuickPayProvider = new(hfbImpl)
type hfbImpl struct{
	 agentId string
	 md5Key string
	 version string
}

func (h *hfbImpl) CheckSign(params map[string]string, signType string, sign string) bool {
	panic("implement me")
}

func (h *hfbImpl) RequestBankSideAuth(nonce string, bankCardNo string, accountName string, idCardNo string, mobile string) (*qpay.BankAuthResult, error) {
	panic("implement me")
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
func (h *hfbImpl) signParams(mp url.Values)string{
	query := string(api.ParamsToBytes(mp,h.md5Key,false))
	query = strings.ToLower(query)
	sign := crypto.Md5([]byte(query))
	println("---",query)
	return sign
}


// 查询银行卡信息
func (h *hfbImpl) QueryCardBin(bankCardNo string)*qpay.CardBinQueryResult{
	apiUrl :=  "https://Pay.heepay.com/API/PayTransit/QueryBankCardInfo.aspx"
	mp := url.Values{
		"agent_id":[]string{h.agentId},
		"bank_card_no":[]string{bankCardNo},
		"key":[]string{h.md5Key},
		"version":[]string{h.version},
	}
	sign := h.signParams(mp)
	mp["sign"] = []string{sign}
	cli := http.Client{}
	rsp,err := cli.PostForm(apiUrl,mp)
	if err != nil{
		return &qpay.CardBinQueryResult{ErrMsg: err.Error()}
	}
	reader := transform.NewReader(rsp.Body,simplifiedchinese.GBK.NewDecoder())
	body,err := ioutil.ReadAll(reader)
	if err != nil{
		return &qpay.CardBinQueryResult{ErrMsg: err.Error()}
	}
	var ret cardBinRsp
	err = xml.Unmarshal(body,&ret)
	if err != nil{
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
		CardType:            int32(types.IntCond(ret.BankCardType == "1", 1, 0)),
		RequireBankSideAuth: true,
	}
}
