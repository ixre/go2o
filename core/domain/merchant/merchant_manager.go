package merchant

import (
	"errors"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain/util"
)

var _ merchant.IMerchantManager = new(merchantManagerImpl)

type merchantManagerImpl struct {
	rep     merchant.IMerchantRepo
	valRepo valueobject.IValueRepo
}

func NewMerchantManager(rep merchant.IMerchantRepo,
	valRepo valueobject.IValueRepo) merchant.IMerchantManager {
	return &merchantManagerImpl{
		rep:     rep,
		valRepo: valRepo,
	}
}

// CreateSignUpToken 创建会员申请商户密钥
func (m *merchantManagerImpl) CreateSignUpToken(memberId int64) string {
	return m.rep.CreateSignUpToken(memberId)
}

// GetMemberFromSignUpToken 根据商户申请密钥获取会员编号
func (m *merchantManagerImpl) GetMemberFromSignUpToken(token string) int64 {
	return m.rep.GetMemberFromSignUpToken(token)
}

// RemoveSignUp 删除会员的商户申请资料
func (m *merchantManagerImpl) RemoveSignUp(memberId int) error {
	return m.rep.DeleteMerchantSignUpByMemberId(memberId)
}

// 检查商户注册信息是否正确
func (m *merchantManagerImpl) checkSignUpInfo(v *merchant.MchSignUp) error {
	if v.MemberId <= 0 {
		return errors.New("会员不存在")
	}
	//todo: validate and check merchant name exists
	if v.MchName == "" {
		return merchant.ErrMissingMerchantName
	}
	if v.CompanyName == "" {
		return merchant.ErrMissingCompanyName
	}
	if v.CompanyNo == "" {
		return merchant.ErrMissingCompanyNo
	}
	if v.Address == "" {
		return merchant.ErrMissingAddress
	}
	if v.PersonName == "" {
		return merchant.ErrMissingPersonName
	}
	if v.PersonId == "" {
		return merchant.ErrMissingPersonId
	}
	if util.CheckChineseCardID(v.PersonId) != nil {
		return merchant.ErrPersonCardId
	}
	if v.Phone == "" {
		return merchant.ErrMissingPhone
	}
	if v.CompanyImage == "" {
		return merchant.ErrMissingCompanyImage
	}
	if v.PersonImage == "" {
		return merchant.ErrMissingPersonImage
	}
	return nil
}

// CommitSignUpInfo 提交商户注册信息
func (m *merchantManagerImpl) CommitSignUpInfo(v *merchant.MchSignUp) (int, error) {
	err := m.checkSignUpInfo(v)
	if err != nil {
		return 0, err
	}
	v.Reviewed = enum.ReviewAwaiting
	v.SubmitTime = time.Now().Unix()
	v.UpdateTime = time.Now().Unix()
	return m.rep.SaveSignUpInfo(v)

}

// ReviewMchSignUp 审核商户注册信息
func (m *merchantManagerImpl) ReviewMchSignUp(id int, pass bool, remark string) error {
	var err error
	v := m.GetSignUpInfo(id)
	if v == nil {
		return merchant.ErrNoSuchSignUpInfo
	}
	if pass {
		v.Reviewed = enum.ReviewPass
		v.Remark = ""
		if err = m.createNewMerchant(v); err != nil {
			return err
		}
	} else {
		v.Reviewed = enum.ReviewReject
		v.Remark = remark
		if strings.TrimSpace(v.Remark) == "" {
			return merchant.ErrRequireRejectRemark
		}
	}
	v.UpdateTime = time.Now().Unix()
	_, err = m.rep.SaveSignUpInfo(v)
	return err
}

// 创建新商户
func (m *merchantManagerImpl) createNewMerchant(v *merchant.MchSignUp) error {
	unix := time.Now().Unix()
	panic("implement me")
	mchVal := &merchant.Merchant{
		//MemberId: v.MemberId,
		// // 商户名称
		// Name: v.MchName,
		// // 是否自营
		// SelfSales: 0,
		// // 商户等级
		// Level: 1,
		// // 标志
		// Logo: "",
		// // 公司名称
		// CompanyName: "",
		// // 省
		// Province: int(v.Province),
		// // 市
		// City: int(v.City),
		// // 区
		// District: int(v.District),
		// // 是否启用
		// Enabled: 1,
		// Flag:    1,
		// // 过期时间
		// ExpiresTime: time.Now().Add(time.Hour * time.Duration(24*365)).Unix(),
		// // 注册时间
		// CreateTime: unix,
		// // 更新时间
		// UpdateTime: unix,
		// // 登录时间
		// LoginTime: 0,
		// // 最后登录时间
		// LastLoginTime: 0,
	}
	mch := m.rep.CreateMerchant(mchVal)
	err := mch.SetValue(mchVal)
	if err != nil {
		return err
	}
	mchId, err := mch.Save()
	if err == nil {
		names := m.valRepo.GetAreaNames([]int32{v.Province, v.City, v.District})
		location := strings.Join(names, "")
		ev := &merchant.EnterpriseInfo{
			MchId:        mchId,
			CompanyName:  v.CompanyName,
			CompanyNo:    v.CompanyNo,
			PersonName:   v.PersonName,
			PersonIdNo:   v.PersonId,
			PersonImage:  v.PersonImage,
			Tel:          v.Phone,
			Province:     v.Province,
			City:         v.City,
			District:     v.District,
			Location:     location,
			Address:      v.Address,
			CompanyImage: v.CompanyImage,
			AuthDoc:      v.AuthDoc,
			Reviewed:     v.Reviewed,
			ReviewTime:   unix,
			ReviewRemark: "",
			UpdateTime:   unix,
		}
		_, err = mch.ProfileManager().SaveEnterpriseInfo(ev)
		if err == nil {
			mch.ProfileManager().ReviewEnterpriseInfo(true, "")
		}
	}
	return err
}

// GetSignUpInfo 获取商户申请信息
func (m *merchantManagerImpl) GetSignUpInfo(id int) *merchant.MchSignUp {
	return m.rep.GetMerchantSignUpInfo(id)
}

// GetSignUpInfoByMemberId 获取会员申请的商户信息
func (m *merchantManagerImpl) GetSignUpInfoByMemberId(memberId int) *merchant.MchSignUp {
	return m.rep.GetMerchantSignUpByMemberId(memberId)
}

// GetMerchantByMemberId 获取会员关联的商户
func (m *merchantManagerImpl) GetMerchantByMemberId(memberId int) merchant.IMerchant {
	return m.rep.GetMerchantByMemberId(memberId)
}
