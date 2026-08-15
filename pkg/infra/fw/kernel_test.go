package fw

import (
	"testing"

	"github.com/sirupsen/logrus"
)

type Staff struct {
	Id int
}

type IStaffRepo interface {
	Repository[Staff]
}

type IStaffService interface {
	Service[Staff]
}

// 员工扩展表
type StaffRepo struct {
	BaseRepository[Staff]
}

// 创建仓储
func NewStaffRepo(o ORM) *StaffRepo {
	s := &StaffRepo{}
	s.ORM = o
	return s
}

type StaffService struct {
	BaseService[Staff]
}

// 创建服务
func NewStaffService(repo IStaffRepo) *StaffService {
	s := &StaffService{}
	s.Repo = repo
	return s
}

func TestGet(t *testing.T) {
	var o ORM
	s := NewStaffRepo(o)
	ss := NewStaffService(s)
	staff := Staff{}
	_, err := ss.Save(&staff)
	if err != nil {
		t.Fatal(err)
	}
	id := 1
	s2 := ss.Get(id)
	if s2 == nil || s2.Id != id {
		t.Error("get staff error")
	}
	logrus.Debugf
}
