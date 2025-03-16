package employee

import (
	"errors"
	"fmt"
	"time"

	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/infrastructure/logger"
	"github.com/ixre/gof/domain/eventbus"
	"github.com/ixre/gof/storage"
)

var _ staff.IStaffManager = new(staffManagerImpl)

type staffManagerImpl struct {
	_mch          merchant.IMerchantAggregateRoot
	_repo         staff.IStaffRepo
	_memberRepo   member.IMemberRepo
	_sysRepo      sys.ISystemRepo
	_mchRepo      merchant.IMerchantRepo
	_approvalRepo approval.IApprovalRepository
	_storage      storage.Interface
}

func NewStaffManager(mch merchant.IMerchantAggregateRoot,
	staffRepo staff.IStaffRepo,
	memberRepo member.IMemberRepo,
	sysRepo sys.ISystemRepo,
	mchRepo merchant.IMerchantRepo,
	approvalRepo approval.IApprovalRepository,
	storage storage.Interface,
) staff.IStaffManager {
	return &staffManagerImpl{
		_mch:          mch,
		_repo:         staffRepo,
		_memberRepo:   memberRepo,
		_sysRepo:      sysRepo,
		_mchRepo:      mchRepo,
		_approvalRepo: approvalRepo,
		_storage:      storage,
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
	st := e._sysRepo.GetSystemAggregateRoot().Stations().FindStationByCity(cityCode)
	if st != nil {
		stationId = int(st.GetDomainId())
	}

	mv := mem.GetValue()
	profile := mem.Profile().GetProfile()
	v := &staff.Staff{
		Id:             0,
		MemberId:       memberId,
		StationId:      stationId,
		MchId:          e._mch.GetAggregateRootId(),
		Flag:           0,
		Gender:         int(profile.Gender),
		Nickname:       mv.Nickname,
		WorkStatus:     staff.WorkStatusOffline,
		Grade:          0,
		Status:         1,
		IsCertified:    0,
		CertifiedName:  "",
		PremiumLevel:   0,
		CreateTime:     int(time.Now().Unix()),
		LastOnlineTime: int(time.Now().Unix()),
	}
	_, err := e._repo.Save(v)
	return err
}

func (e *staffManagerImpl) GetStaff(staffId int) *staff.Staff {
	return e._repo.Get(staffId)
}

// RequestTransfer implements staff.IStaffManager.
func (e *staffManagerImpl) RequestTransfer(staffId int, mchId int) (int, error) {
	transRepo := e._repo.TransferRepo()
	count, _ := transRepo.Count("staff_id=? and review_status = ?", staffId, enum.ReviewPending)
	if count > 0 {
		return 0, errors.New("员工存在未审核的转移请求")
	}

	st := e.GetStaff(staffId)
	if st == nil {
		return 0, errors.New("员工不存在")
	}
	if st.MchId != e._mch.GetAggregateRootId() {
		return 0, errors.New("员工不属于当前商户")
	}
	if st.MchId == mchId {
		return 0, errors.New("员工已存在于目标商户")
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
	ret, err := transRepo.Save(transferRequest)
	if err == nil {
		transferRequest.Id = ret.Id
		// 创建审批单
		ia := e._approvalRepo.Create(approval.FlowStaffTransfer, transferRequest.Id)
		err = ia.Save()
		if err == nil {
			// 设置审批单ID
			transferRequest.ApprovalId = ia.GetAggregateRootId()
			transRepo.Save(transferRequest)
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

// TransferApproval implements staff.IStaffManager.
func (e *staffManagerImpl) TransferApproval(trans *staff.StaffTransfer, event *approval.ApprovalProcessEvent) error {
	ap := event.Approval
	if ap.IsFinal() {
		// 审批结束
		isPass := ap.GetApproval().FinalStatus == approval.FinalPassStatus
		if isPass {
			// 通过审核
			trans.ReviewStatus = int(enum.ReviewApproved)
		}
		if ap.GetApproval().FinalStatus == approval.FinalRejectStatus {
			// 未通过审核
			trans.ReviewStatus = int(enum.ReviewRejected)
			trans.ReviewRemark = event.Tx.ApprovalRemark
		}
		trans.UpdateTime = int(time.Now().Unix())
		_, err := e._repo.TransferRepo().Save(trans)
		if err == nil && isPass {
			// 通过审核,更新职员信息
			im := e._mchRepo.GetMerchant(trans.TransferMchId)
			if im == nil {
				return errors.New("商户不存在")
			}
			isn := e._sysRepo.GetSystemAggregateRoot().Stations().FindStationByCity(im.GetValue().City)
			st := e.GetStaff(trans.StaffId)
			st.MchId = trans.TransferMchId
			if isn != nil {
				st.StationId = isn.GetDomainId()
			} else {
				st.StationId = 0
				logger.Error("station not found, staffId: %d, mchId: %d", st.Id, st.MchId)
			}
			_, err = e._repo.Save(st)
			if err == nil {
				// 发布员工转移审批通过事件
				go eventbus.Dispatch(&staff.StaffTransferApprovedEvent{
					Staff:         *st,
					TransferMchId: trans.TransferMchId,
					OriginMchId:   trans.OriginMchId,
				})
			}
		}
		return err
	}
	if event.NodeKey == "aggree" {
		// 指派审批对象为新的商户
		transMch := e._mchRepo.GetMerchant(trans.TransferMchId)
		err := ap.Assign(trans.TransferMchId, transMch.GetValue().MchName)
		if err != nil {
			logger.Error("approval assign error: %v, transferId: %d", err, ap.GetApproval().BizId)
			panic(err)
		}
	}
	return nil
}

// IsKeepOnline implements staff.IStaffManager.
func (e *staffManagerImpl) IsKeepOnline(staffId int) bool {
	// 默认保持在线
	key := fmt.Sprintf("go2o:staff:keep_online:%d", staffId)
	v, _ := e._storage.GetInt(key)
	return v != -1
}

// UpdateWorkStatus implements staff.IStaffManager.
func (e *staffManagerImpl) UpdateWorkStatus(staffId int, workStatus int, isKeepOnline bool) error {
	key := fmt.Sprintf("go2o:staff:keep_online:%d", staffId)
	if isKeepOnline {
		// 继续保持在线状态
		e._storage.Delete(key)
	} else {
		// 下线状态
		err := e._storage.Set(key, "-1")
		if err != nil {
			return err
		}
	}
	st := e.GetStaff(staffId)
	st.WorkStatus = workStatus
	st.LastOnlineTime = int(time.Now().Unix())
	_, err := e._repo.Save(st)
	return err
}
