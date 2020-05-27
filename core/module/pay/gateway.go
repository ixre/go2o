package pay

import (
	"errors"
	"fmt"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/storage"
	"go2o/core/domain/interface/member"
	"go2o/core/repos"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type wrapperData struct {
	// 支付结果
	State string
	// 支付金额
	Amount int
	// 手续费
	ProcedureFee int
	// 主题
	Subject string
	// 交易号
	TradeNo string
	// 支付标志
	PayFlag int
	// 商品地址
	ItemUrl string
	// 通知URL
	NotifyUrl string
	// 返回URL
	ReturnUrl string
	// 数据
	Data map[string]string
}

const (
	// 仅验证密码
	FlagOnlyCheck = 1 << iota
	// 余额抵扣
	FlagBalanceDiscount
	// 积分抵扣
	FlagIntegralDiscount
	// 钱包支付
	FlagWalletPayment
)

// 支付网关
type Gateway struct {
	s          storage.Interface
	memberRepo member.IMemberRepo
}

func NewGateway(s storage.Interface) *Gateway {
	return &Gateway{
		s:          s,
		memberRepo: repos.Repo.GetMemberRepo(),
	}
}

// 生成支付网关提交令牌,5分钟内有效
func (g *Gateway) CreatePostToken(userId int64) string {
	unix := time.Now().UnixNano()
	str := fmt.Sprintf("%d-%d", unix, userId)
	token := crypto.Md5([]byte(str))
	rdsKey := "go2o:pay:gateway:token:user-" + strconv.Itoa(int(userId))
	g.s.SetExpire(rdsKey, token, 300) //5分钟过期
	return token
}

// 对比支付网关提交的令牌
func (g *Gateway) verifyPostToken(userId int64, token string) bool {
	rdsKey := "go2o:pay:gateway:token:user-" + strconv.Itoa(int(userId))
	src, _ := g.s.GetString(rdsKey)
	return token == src
}

// 连接URL
func (g *Gateway) urlJoin(url string, query string) string {
	i := strings.Index(url, "?")
	if i == -1 {
		return url + "?" + query
	}
	return url + "&" + query
}

func (g *Gateway) getTradeKey(tradeNo string) string {
	return "go2o:pay:gateway:trade-" + tradeNo
}

// 提交到网关
func (g *Gateway) Submit(userId int64, data map[string]string) error {
	amount, err := strconv.Atoi(data["amount"])
	if err != nil {
		amount = 0
	}
	prFee, err1 := strconv.Atoi(data["procedure_fee"])
	if err1 != nil {
		prFee = 0
	}
	flag, err2 := strconv.Atoi(data["pay_flag"])
	if err2 != nil {
		flag = FlagOnlyCheck
	}

	d := &wrapperData{
		TradeNo:      data["trade_no"],
		Amount:       amount,
		ProcedureFee: prFee,
		PayFlag:      flag,
		Subject:      data["subject"],
		ItemUrl:      data["item_url"],
		NotifyUrl:    data["notify_url"],
		ReturnUrl:    data["return_url"],
		Data: map[string]string{
			"token": data["token"],
		},
	}
	delete(data, "trade_no")
	delete(data, "amount")
	delete(data, "procedure_fee")
	delete(data, "pay_flag")
	delete(data, "subject")
	delete(data, "item_url")
	delete(data, "notify_url")
	delete(data, "return_url")
	delete(data, "token")
	for k, v := range data {
		d.Data[k] = v
	}
	return g.realSubmit(userId, d)
}

// 存储提交数据
func (g *Gateway) realSubmit(userId int64, data *wrapperData) error {
	token := data.Data["token"]
	if token == "" {
		return errors.New("提交支付网关错误:NO_TOKEN")
	}
	if !g.verifyPostToken(userId, token) {
		return errors.New("提交支付网关错误:TOKEN_NOT_MATCH")
	}
	if data.TradeNo == "" {
		return errors.New("参数不完整:trade_no")
	}
	if data.Amount+data.ProcedureFee <= 0 {
		return errors.New("支付金额错误")
	}
	if data.NotifyUrl == "" {
		return errors.New("参数不完整:notify_url")
	}
	data.State = "prepare"
	rdsKey := g.getTradeKey(data.TradeNo)
	err := g.s.SetExpire(rdsKey, data, 3600*72)
	return err
}

// 模拟支付
func (g *Gateway) CheckAndPayment(userId int64, tradeNo string, tradePwd string) error {
	m := g.memberRepo.GetMember(userId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	mv := m.GetValue()
	if mv.TradePwd == "" {
		return errors.New("您还未设置交易密码")
	}
	if mv.TradePwd != tradePwd {
		return errors.New("交易密码不正确")
	}
	rk := g.getTradeKey(tradeNo)
	data := wrapperData{}
	g.s.Get(rk, &data)
	// 处理付款
	err := g.handlePayment(userId, tradeNo, data)
	if err == nil {
		// 处理通知
		err = g.notify(userId, &data)
		if err == nil {
			data.State = "success"
			g.s.SetExpire(rk, &data, 3600*12)
		}
	}
	return err
}
func (g *Gateway) handlePayment(userId int64, tradeNo string, data wrapperData) error {
	// 仅验证交易密码
	if data.PayFlag&FlagOnlyCheck != 0 {
		return nil
	}
	// 余额抵扣
	if data.PayFlag&FlagBalanceDiscount != 0 {
		//todo:
	}
	return nil
}

//通知支付结果,响应端返回success表示处理完成
func (g *Gateway) notify(userId int64, data *wrapperData) error {
	cli := http.Client{}
	values := url.Values{
		"user_id":       []string{strconv.Itoa(int(userId))},
		"trade_no":      []string{data.TradeNo},
		"state":         []string{"success"},
		"amount":        []string{strconv.Itoa(data.Amount)},
		"procedure_fee": []string{strconv.Itoa(data.ProcedureFee)},
		"flag":          []string{strconv.Itoa(data.PayFlag)},
		"subject":       []string{data.Subject},
	}
	for k, v := range data.Data {
		if k != "token" {
			values[k] = []string{v}
		}
	}
	rsp, err := cli.PostForm(data.NotifyUrl, values)
	// 未通知成功
	if err != nil {
		log.Println("[ Go2o][ Pay][ Gateway]: notify failed :",
			err.Error(), " [URL]:", data.NotifyUrl)
		return errors.New("通知支付结果失败")
	}
	body, _ := ioutil.ReadAll(rsp.Body)
	rspTxt := string(body)
	// 响应状态不正确
	if rsp.StatusCode != 200 {
		log.Println("[ Go2o][ Pay][ Gateway]: notify failed :",
			rspTxt, " [URL]:", data.NotifyUrl)
		return errors.New("通知支付结果失败")
	}
	// 判断响应内容
	if rspTxt != "success" {
		log.Println("[ Go2o][ Pay][ Gateway]: notify response :", rspTxt)
		return errors.New("通知支付结果异常：" + rspTxt)
	}
	return nil
}

func (g *Gateway) GetData() *wrapperData {
	return nil
}
