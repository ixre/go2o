package events

import (
	"testing"
	"time"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/inject"
	_ "github.com/ixre/go2o/tests"
	"github.com/ixre/gof/domain/eventbus"
)

func TestWalletLogUpdate(t *testing.T) {
	id := 158113
	repo := inject.GetWalletRepo()
	l := repo.GetWalletLog_(id)
	l.Subject = l.Subject + "_1"
	_, err := repo.SaveWalletLog_(l)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	time.Sleep(time.Second * 2)
}

func TestPushMemberRegisterEvent(t *testing.T) {
	eventbus.Dispatch(&events.MemberPushEvent{
		IsCreate:  false,
		Member:    &member.Member{},
		InviterId: 0,
	})
	time.Sleep(100 * time.Second)
	t.Log("test finished...")
}
