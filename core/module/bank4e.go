package module

import (
	"encoding/json"
	"errors"
	"github.com/ixre/gof"
	"github.com/ixre/gof/storage"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain/util"
	"go2o/core/infrastructure/format"
	"go2o/core/module/bank"
	"go2o/core/repos"
	"go2o/core/service/thrift"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	_            Module = new(Bank4E)
	zhNameRegexp        = regexp.MustCompile("^[\u4e00-\u9fa5]{2,6}$")
)
var keys = []string{"bank4e_trust_on", "bank4e_jd_app_key"}

type Bank4E struct {
	memberRepo member.IMemberRepo
	valueRepo  valueobject.IValueRepo
	storage    storage.Interface
	appKey     string
	open       bool
}

func (b *Bank4E) SetApp(app gof.App) {
	b.storage = app.Storage()
}

func (b *Bank4E) Init() {
	b.memberRepo = repos.Repo.GetMemberRepo()
	b.valueRepo = repos.Repo.GetValueRepo()
	trans, cli, err := thrift.FoundationServeClient()
	if err == nil {
		defer trans.Close()
		cli.CreateUserRegistry(thrift.Context, keys[0], "false", "是否开启四要素实名认证")
		cli.CreateUserRegistry(thrift.Context, keys[1], "", "京东银行四要素接口KEY")
		data, _ := cli.GetRegistries(thrift.Context, keys)
		b.open, _ = strconv.ParseBool(data[keys[0]])
		b.appKey = data[keys[1]]
	}
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
	if info.ReviewState == int(enum.ReviewPass) {
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
		bankName, _, _ := bank.GetNameByAccountNo(bankAccount)
		bankArr := strings.Split(bankName, ".")
		data["ErrCode"] = "true"
		data["BankName"] = bankArr[0]
		data["ErrMsg"] = "PASS"
	} else {
		data["ErrCode"] = "false"
		data["ErrMsg"] = err.Error()
	}
	return data
}

// 更新信息
func (b *Bank4E) UpdateInfo(memberId int64, realName, idCard, phone, bankAccount string) error {
	m := b.memberRepo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	// 校验姓名
	if !zhNameRegexp.MatchString(realName) {
		return errors.New("真实姓名不正确")
	}
	// 校验身份证
	idCard = strings.ToUpper(idCard)
	err := util.CheckChineseCardID(idCard)
	if err != nil {
		return err
	}
	bankName, cardType, err := bank.GetNameByAccountNo(bankAccount)
	if err != nil {
		return err
	}
	if bankName == "" || cardType == 0 {
		return errors.New("银行卡号无法识别")
	}
	if cardType != 1 {
		return errors.New("您填写的卡号不是有效的储蓄卡")
	}

	mv := m.Profile().GetProfile()
	if mv.Phone != "" && mv.Phone != phone {
		return errors.New("手机号码非法`")
	}
	info := m.Profile().GetTrustedInfo()
	if info.ReviewState == int(enum.ReviewPass) {
		return errors.New("您已通过实名认证")
	}

	if !b.checkLatestInfo(memberId, realName, idCard, phone, bankAccount) {
		return errors.New("信息不正确，认证未通过")
	}

	//log.Println("---请求服务器")
	//return errors.New("未通过")

	// 验证四要素
	result := b.Check(realName, idCard, phone, bankAccount)
	// 验证不通过，则返回错误
	if result["ErrCode"] == "false" {
		return errors.New(result["ErrMsg"])
	}

	// 移除四要素验证记录
	b.removeLatestInfo(memberId)

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

// 检查是否为上次提交，如果不是，则更新
func (b *Bank4E) checkLatestInfo(memberId int64, realName, idCard, phone, bankAccount string) bool {
	// 获取是否关闭检查
	keyStat := "sys:go2o:b4e:turn_stat"
	turnStat, _ := b.storage.GetInt(keyStat)
	// 获取之前提交信息
	key := "sys:go2o:b4e:last-post:" + strconv.Itoa(int(memberId))
	result := strings.Join([]string{realName, idCard, phone, bankAccount}, "|")
	src, _ := b.storage.GetString(key)
	if src != result || turnStat == 0 {
		b.storage.SetExpire(key, result, int64(time.Hour)*24*100)
		return true
	}
	return false
}

func (b *Bank4E) removeLatestInfo(memberId int64) {
	key := "sys:go2o:b4e:last-post:" + strconv.Itoa(int(memberId))
	b.storage.Del(key)
}

// 打开/关闭重复信息验证
func (b *Bank4E) turnCheckInfo(r bool) {
	key := "sys:go2o:b4e:turn_stat"
	if r {
		b.storage.Set(key, 1)
	} else {
		b.storage.Set(key, 0)
	}
}

// 调用验证接口
func (b *Bank4E) b4eApi(realName, idCard, phone, bankAccount string) error {
	apiServer := "https://way.jd.com/youhuoBeijing/QryBankCardBy4Element"
	if !b.open {
		return errors.New("未开启四要素实名认证")
	}
	if b.appKey == "" {
		return errors.New("验证接口未配置")
	}
	cli := http.Client{}
	data := url.Values{
		"appkey":        []string{b.appKey},
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
			rc := mp["code"].(string)
			if rc != "10000" {
				// 关闭重复信息验证
				b.turnCheckInfo(false)
				return errors.New("接口返回异常:" + mp["msg"].(string))
			}
			// 打开重复信息验证
			b.turnCheckInfo(true)
			r1 := mp["result"]
			if r1 == nil {
				return errors.New("认证服务异常")
			}
			rm1 := r1.(map[string]interface{})
			r2 := rm1["result"]
			if r2 == nil {
				return errors.New("认证服务异常")
			}
			rm2 := r2.(map[string]interface{})
			//authResult := r2["result"].(string)
			//authMsg := r2["message"].(string)
			authResult := rm2["result"].(string)
			if authResult == "F" {
				return errors.New("信息不正确，认证未通过")
			}
			if authResult == "N" {
				return errors.New("银行卡号无法完成认证")
			}
			return nil
		}
		return errors.New("请求失败：HTTP " + strconv.Itoa(rsp.StatusCode))
	}
	return err
}
