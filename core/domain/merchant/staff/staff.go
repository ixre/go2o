package employee

import (
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/domain/interface/station"
)

var _ staff.IStaffManager = new(staffManagerImpl)

type staffManagerImpl struct {
	_mch         merchant.IMerchant
	_repo        staff.IStaffRepo
	_memberRepo  member.IMemberRepo
	_stationRepo station.IStationRepo
}

func NewStaffManager(mch merchant.IMerchant,
	staffRepo staff.IStaffRepo,
	memberRepo member.IMemberRepo,
	stationRepo station.IStationRepo) staff.IStaffManager {
	return &staffManagerImpl{
		_mch:         mch,
		_repo:        staffRepo,
		_memberRepo:  memberRepo,
		_stationRepo: stationRepo,
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
	id, err := e._repo.SaveStaff(v)
	if err == nil {
		v.Id = id
	}
	return err
}
