package hfb

import (
	"bytes"
	"encoding/hex"
	"github.com/ixre/go2o/core/infrastructure"
	"github.com/ixre/go2o/core/infrastructure/qpay"
	"github.com/ixre/gof/crypto"
	"go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : hfb_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-05 20:01
 * description :
 * history :
 */

var h qpay.QuickPayProvider

func init() {
	// 默认的ETCD端点
	etcdEndPoints := []string{"http://127.0.0.1:2379"}
	cfg := clientv3.Config{
		Endpoints:   etcdEndPoints,
		DialTimeout: 5 * time.Second,
	}
	s, _ := infrastructure.NewEtcdStorage(cfg)
	h = NewHfb(s)
}

func TestCardBin(t *testing.T) {
	bankCardNo := "6227000010990006191"
	r := h.QueryCardBin(bankCardNo)
	t.Logf("%#v", r)
}

func TestHfbImpl_RequestBankSideAuth(t *testing.T) {
	bankCardNo := "6227000010990006191"
	nonce := strconv.Itoa(int(time.Now().Unix()))
	r, err := h.RequestBankSideAuth(nonce, bankCardNo, "闫雪龙",
		"22011219850823101X", "13810512111")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(r.AuthForm)
}
func TestHfbImpl_DirectPayment(t *testing.T) {
	bankCardNo := "6227000010990006191"
	rsp, err := h.QueryBankAuth(bankCardNo)
	if err == nil {
		t.Logf("%#v", rsp)
	}
	if rsp.Code != 0 {
		t.Log("卡片未授权成功")
		t.FailNow()
	}
	orderNo := "BZ" + time.Now().Format("20060102150405")
	ret, err := h.DirectPayment(orderNo, 1, "补差价链接", rsp.BankAuthToken,
		"127.0.0.1", "http://www.go2o-dev.56x.net/qpay/notify_url",
		"")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("支付成功,订单号：", orderNo, "第三方订单号：", ret.BillNo)
}

func TestHfbImpl_QueryPaymentStatus(t *testing.T) {
	orderNo := "bz20201107224102"
	billTime := "20201107224102"
	ret, err := h.QueryPaymentStatus(orderNo, map[string]string{"agent_bill_time": billTime})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("%#v", ret)
}

func TestHfbImpl_QueryBankAuth(t *testing.T) {
	//bankCardNo := "6227000010990006191"
	bankCardNo := "9559981014352796313"

	r, err := h.QueryBankAuth(bankCardNo)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(r)
}

func TestHfbImpl_BatchTransfer(t *testing.T) {
	orderNo := "BZ" + time.Now().Format("20060102150405")
	nonce := strconv.Itoa(int(time.Now().Unix()))
	batch := []*qpay.CardTransferReq{
		{
			OrderNo:         orderNo,
			PersonTransfer:  true,
			TradeFee:        1, //0.01
			BankCardNo:      "6227000010990006191",
			BankCode:        "2",
			BankAccountName: "闫雪龙",
			Subject:         "",
			Province:        "-",
			City:            "-",
			StoreName:       "-",
		},
	}
	r, err := h.BatchTransfer(orderNo, batch, nonce, "http://www.go2o-dev.56x.net/qpay/callback")
	if err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		t.Logf("%#v", r)
	}
}

func TestHfbImpl_Encrypt3DES(t *testing.T) {
	detail_data := "agent_id=2126129&batch_amt=0.01&batch_no=bz20201108080915&batch_num=1&detail_data=bz20201108080915^2^0^6227000010990006191^闫雪龙^0.01^上游结算款^^^&ext_param1=1604794155&key=0e05aac3be0746269f114bd7&notify_url=http://www.go2o-dev.56x.net/qpay/callback&version=3"
	key := []byte("4865534446254C0F8837DFB3")
	bytes, _ := crypto.EncryptECB3DES([]byte(detail_data), key)
	s := hex.EncodeToString(bytes)
	t.Log(strings.ToUpper(s))
}

func TestHfbImpl_CheckBankCert(t *testing.T) {
	action := "https://gateway.95516.com/gateway/api/frontTransReq.do"
	query := "bizType=000902&tokenPayData=%7BtrId%3D62000004640%26tokenType%3D01%7D&orderId=0000015558661117&txnSubType=00&backUrl=https%3A%2F%2Fpay.heepay.com%2FBank%2FUnionPayAuthAsyncNotify.aspx&signature=Iz8E3Y6hNNM50CIr5DR4t4Cm1UBx8QF2ZlrdCq%2BqzydgOZmvN6oajR6im7KaHdyVRIhJrRQvStFYKisFCHlDFQyXvNBEe8uoTOefep8e4CchpK0PfLX%2FQToPDEcuGMbSK9FinVU%2FUvn3WK5pidRwhnSPDc8WvcbLLMlOgdkODqsBhqhly0NUhfoi89OJ9Hk2LbjGDqPEplSIyNZ35Fjw3IK8o%2FLDYIf%2FR%2FskdtEnae48D%2FiqEjCFCaJf7QXbcSWRVeb5PJ6df4h8jgvQjcYBP8TfWoBQ91AXNyPBGUltcCMOnjMNAXNbv51WIeQov8F%2F%2B4htSZmBkwErkfVEdJa3rg%3D%3D&accNo=a8Fpjs5ylNGHKULsEY1KRYkernhcslfV4LB4X615lkh3EQAqvqjy6a1kaHjjRRalQm89xW0bPhyFEwKjblUGnjD1xMoHdEXOfhsPTmOm6WB5lNNnszpG63uXGftWqtkpUE2VQ3UXddybcLhbh%2BhkiuGEn0d71ruaFBU98R9BDgvc6cO4o0fc5E6XEredAOcvAVAh4LvOfJ3%2FMvs21PsfXdhLsEkP94fjRXa%2BP9gnt5ymaWLDPt%2BoNPDwZKFl2dYa4Rvmq2kxF9uGO9fHuGrBOb2T82GIAwjrvNr68Hd%2B1XtwFyFMIzOR%2FEMu87ynMqf7X3D%2B36ZD59IsGZ%2Fl3XoFzg%3D%3D&customerInfo=e2VuY3J5cHRlZEluZm89M0RCOElSOEZ0K0t0WW1LTVdqTFBzRUl1dzZKQ1Z6R1hJRWk5aFNQdmZUaStxMWE0cWNEVE54QWFiWGVZaFo2elNlOUFiY05ESUQveXhneWtPZVNsSEdlWkZYZ2hWZjVmQjQrb3BWRnRybWo3NXZJL3VDUkdKandJWTAwMzErMnY1Qk1OMHJRSHdVeHZiZ0czVXB6Z0hpalVhbVZQWEFyMEtSVVkwN0J2VDd1NUVRRUh4S01kWXQrcExDSTF6WFhkVHRuUlNBZVNBL3BYYjg4aFlORWg1UWo0d0l1b1ZkOU9wMkNLdEJLZXFyaWpnTWQ0NXdoSytTOWk1REtVZUlRRGE1SGZxTFA0N2o1bmlFUlNjK0RBUGYvMzJWdGtqc25CWW9xSERvWkVUbThDSThvYmhxUG1WZzY4UGpWelpiMDVTMXJYNy9XditFZFpGTWdJbWhjVER3PT0mY2VydGlmVHA9MDEmY2VydGlmSWQ9MjIwMTEyMTk3MjAzMDgxODA3JmN1c3RvbWVyTm09546L5pa5fQ%3D%3D&merName=%E5%B9%BF%E4%B8%9C%E7%A5%A5%E7%9D%BF%E7%A7%91%E6%8A%80%E6%9C%89%E9%99%90%E5%85%AC%E5%8F%B8&txnType=79&channelType=07&frontUrl=https%3A%2F%2Fpay.heepay.com%2FBank%2FUnionPayAuthReceive.aspx&certId=75456652545&encoding=UTF-8&acqInsCode=49473930&version=5.1.0&merAbbr=%E4%BF%9D%E8%B4%A6&accessType=1&encryptCertId=77447321186&txnTime=20201117121248&merId=947393048161131&payTimeout=20201117122749&merCatCode=5211&signMethod=01"
	values, _ := url.ParseQuery(query)
	cli := http.Client{}
	req, _ := http.NewRequest("POST", action, bytes.NewReader([]byte(values.Encode())))
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:82.0) Gecko/20100101 Firefox/82.0")
	rsp, _ := cli.Do(req)
	bytes, _ := ioutil.ReadAll(rsp.Body)
	t.Log(string(bytes))

}
