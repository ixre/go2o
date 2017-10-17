package module

import (
	"encoding/json"
	"errors"
	"github.com/jsix/gof"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/factory"
	"go2o/core/infrastructure/format"
	"go2o/core/module/bank"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var _ Module = new(Bank4E)

type Bank4E struct {
	memberRepo member.IMemberRepo
	valueRepo  valueobject.IValueRepo
}

func (b *Bank4E) SetApp(app gof.App) {
}

func (b *Bank4E) Init() {
	b.memberRepo = factory.Repo.GetMemberRepo()
}

// 获取基础信息
func (b *Bank4E) GetBasicInfo(memberId int64) map[string]string {
	data := map[string]string{}
	m := b.memberRepo.GetMember(memberId)
	if m == nil {
		data["Error"] = "会员不存在"
		return data
	}
	pr := m.Profile().GetProfile()
	info := m.Profile().GetTrustedInfo()
	bank := m.Profile().GetBank()
	data["RealName"] = info.RealName
	data["IDCard"] = info.CardId
	data["Phone"] = pr.Phone
	data["BankAccount"] = bank.Account
	data["Remark"] = info.Remark
	if info.Reviewed == enum.ReviewPass {
		data["Reviewed"] = "true"
	} else {
		data["Reviewed"] = "false"
	}
	return data
}

// 判断四要素是否一致
func (b *Bank4E) Check(realName, idCard, phone, bankAccount string) map[string]string {
	data := map[string]string{}
	err := b.b4eApi(realName, idCard, phone, bankAccount)
	if err == nil {
		bankName := bank.GetNameByAccountNo(bankAccount)
		bankArr := strings.Split(bankName, ".")
		data["Result"] = "true"
		data["BankName"] = bankArr[0]
		data["Message"] = "PASS"
	} else {
		data["Result"] = "false"
		data["Message"] = err.Error()
	}
	return data
}

// 更新信息
func (b *Bank4E) UpdateInfo(memberId int64, realName, idCard, phone, bankAccount string) error {
	m := b.memberRepo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	mv := m.Profile().GetProfile()
	if mv.Phone != "" && mv.Phone != phone {
		return errors.New("手机号码非法`")
	}
	info := m.Profile().GetTrustedInfo()
	if info.Reviewed == enum.ReviewPass {
		return errors.New("您已通过实名认证")
	}

	// 验证四要素
	result := b.Check(realName, idCard, phone, bankAccount)
	// 验证不通过，则返回错误
	if result["Result"] == "false" {
		return errors.New(result["Message"])
	}

	// 保存手机号码
	if mv.Phone == "" {
		mv.Phone = phone
		err := m.Profile().SaveProfile(&mv)
		if err != nil {
			return err
		}
	}

	// 保存实名信息
	if err := m.Profile().SaveTrustedInfo(&member.TrustedInfo{
		RealName:   realName,
		CardId:     idCard,
		TrustImage: format.GetResUrl(""),
	}); err != nil {
		return err
	}

	// 审核通过实名信息
	if err := m.Profile().ReviewTrustedInfo(true, ""); err != nil {
		return err
	}

	// 保存银行信息
	m.Profile().UnlockBank()
	if err := m.Profile().SaveBank(&member.BankInfo{
		BankName:    result["BankName"],
		AccountName: realName,
		Account:     bankAccount,
	}); err != nil {
		return err
	}

	return nil
}

// 调用验证接口
func (b *Bank4E) b4eApi(realName, idCard, phone, bankAccount string) error {
	appKey := b.valueRepo.GetsRegistry([]string{"go2o:bank-e4:config:app_key"})[0]
	apiServer := "https://way.jd.com/youhuoBeijing/QryBankCardBy4Element"
	if appKey == "" {
		return errors.New("验证接口未配置")
	}
	cli := http.Client{}
	data := url.Values{
		"appkey":        []string{appKey},
		"accountNo":     []string{bankAccount},
		"name":          []string{realName},
		"idCardCode":    []string{idCard},
		"bankPreMobile": []string{phone},
	}
	rsp, err := cli.PostForm(apiServer, data)
	if err == nil {
		// 响应状态正确
		if rsp.StatusCode == 200 {
			body, _ := ioutil.ReadAll(rsp.Body)
			mp := make(map[string]interface{})
			json.Unmarshal(body, &mp)
			r1 := mp["result"].(map[string]interface{})
			r2 := r1["result"].(map[string]interface{})
			//authResult := r2["result"].(string)
			authMsg := r2["message"].(string)
			authMsgType := r2["messagetype"].(float64)
			if authMsgType == 1 {
				return errors.New("认证信息不正确：" + authMsg)
			}
			if authMsgType == 2 {
				return errors.New("银行卡号异常")
			}
			return nil
		}
		return errors.New("请求失败：HTTP " + strconv.Itoa(rsp.StatusCode))
	}
	return err
}
