package employee

import (
	"errors"
	"time"

	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/domain/interface/station"
)

var _ staff.IStaffManager = new(staffManagerImpl)

type staffManagerImpl struct {
	_mch               merchant.IMerchantAggregateRoot
	_repo              staff.IStaffRepo
	_staffTransferRepo staff.IStaffTransferRepo
	_memberRepo        member.IMemberRepo
	_stationRepo       station.IStationRepo
	_mchRepo           merchant.IMerchantRepo
	_approvalRepo      approval.IApprovalRepository
}

func NewStaffManager(mch merchant.IMerchantAggregateRoot,
	staffRepo staff.IStaffRepo,
	staffTransferRepo staff.IStaffTransferRepo,
	memberRepo member.IMemberRepo,
	stationRepo station.IStationRepo,
	mchRepo merchant.IMerchantRepo,
	approvalRepo approval.IApprovalRepository,
) staff.IStaffManager {
	return &staffManagerImpl{
		_mch:               mch,
		_repo:              staffRepo,
		_staffTransferRepo: staffTransferRepo,
		_memberRepo:        memberRepo,
		_stationRepo:       stationRepo,
		_mchRepo:           mchRepo,
		_approvalRepo:      approvalRepo,
	}
}

// Create implements staff.IStaffManager.
func (e *staffManagerImpl) Create(memberId int) error {
	// 获取会员信息,并检查会员是否有效
	mem := e._memberRepo.GetMember(int64(memberId))
	if mem == nil {
		return member.ErrNoSuchMember
	}
	if mem.ContainFlag(member.FlagLocked) {
		return member.ErrMemberLocked
	}
	// role := mem.GetValue().RoleFlag
	// if domain.TestFlag(role, member.RoleEmployee) {
	// }
	// 查询会员是否已存在在职
	exists := e._repo.GetStaffByMemberId(memberId)
	if exists != nil {
		return staff.ErrStaffAlreadyExists
	}
	// 获取站点,站点允许为0
	stationId := 0
	cityCode := e._mch.GetValue().City
	st := e._stationRepo.GetStationByCity(cityCode)
	if st != nil {
		stationId = int(st.GetAggregateRootId())
	}
	// 创建职员
	mv := mem.GetValue()
	profile := mem.Profile().GetProfile()
	v := &staff.Staff{
		Id:            0,
		MemberId:      memberId,
		StationId:     stationId,
		MchId:         e._mch.GetAggregateRootId(),
		Flag:          0,
		Gender:        int(profile.Gender),
		Nickname:      mv.Nickname,
		WorkStatus:    staff.WorkStatusOffline,
		Grade:         0,
		Status:        1,
		IsCertified:   0,
		CertifiedName: "",
		PremiumLevel:  0,
		CreateTime:    int(time.Now().Unix()),
	}
	_, err := e._repo.Save(v)
	return err
}

// RequestTransfer implements staff.IStaffManager.
func (e *staffManagerImpl) RequestTransfer(staffId int, mchId int) (int, error) {
	count, _ := e._staffTransferRepo.Count("staff_id=? and review_status = ?", staffId, enum.ReviewPending)
	if count > 0 {
		return 0, errors.New("员工存在未审核的转移请求")
	}
	mch := e._mchRepo.GetMerchant(mchId)
	if mch == nil {
		return 0, errors.New("商户不存在")
	}
	// 创建转移请求
	transferRequest := &staff.StaffTransfer{
		Id:            0,
		StaffId:       staffId,
		OriginMchId:   e._mch.GetAggregateRootId(),
		TransferMchId: mchId,
		ApprovalId:    0,
		ReviewStatus:  int(enum.ReviewPending),
		ReviewRemark:  "",
		CreateTime:    int(time.Now().Unix()),
		UpdateTime:    int(time.Now().Unix()),
	}
	ret, err := e._staffTransferRepo.Save(transferRequest)
	if err == nil {
		transferRequest.Id = ret.Id
		// 创建审批单
		ia := e._approvalRepo.Create(approval.FlowStaffTransfer, transferRequest.Id)
		err = ia.Save()
		if err == nil {
			// 分配审批人
			err = ia.Assign(e._mch.GetAggregateRootId(), e._mch.GetValue().MchName)
		}
		if err == nil {
			return ret.Id, nil
		}
		return 0, err
	}
	return transferRequest.Id, err
}
