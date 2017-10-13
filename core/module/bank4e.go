package module

import (
	"errors"
	"github.com/jsix/gof"
	"go2o/core/domain/interface/member"
	"go2o/core/factory"
)

var _ Module = new(Bank4E)

type Bank4E struct {
	memberRepo member.IMemberRepo
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
	return data
}

// 判断四要素是否一致
func (b *Bank4E) Check(realName, idCard, phone, bankAccount string) (
	data map[string]string, err error) {
	data = map[string]string{}
	data["Result"] = "true"
	data["BankName"] = "中国民生银行"
	data["Message"] = "PASS"
	return data, nil
}

// 更新信息
func (b *Bank4E) UpdateInfo(memberId int64, realName, idCard, phone, bankAccount string) error {
	m := b.memberRepo.GetMember(memberId)
	if m == nil {
		return member.ErrNoSuchMember
	}
	// 验证四要素
	result, err := b.Check(realName, idCard, phone, bankAccount)
	if err != nil {
		return err
	}
	// 验证不通过，则返回错误
	if result["Result"] == "false" {
		return errors.New(result["Message"])
	}

	// 保存银行信息
	m.Profile().UnlockBank()
	if err = m.Profile().SaveBank(&member.BankInfo{
		BankName:    result["BankName"],
		AccountName: realName,
		Account:     bankAccount,
	}); err != nil {
		return err
	}

	// 保存实名信息
	if err = m.Profile().SaveTrustedInfo(&member.TrustedInfo{
		RealName:   realName,
		CardId:     idCard,
		TrustImage: "",
	}); err != nil {
		return err
	}
	// 审核通过实名信息
	if err = m.Profile().ReviewTrustedInfo(true, ""); err != nil {
		return err
	}
	return nil
}
