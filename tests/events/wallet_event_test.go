package events

import (
	"github.com/ixre/go2o/tests/ti"
	"testing"
	"time"
)

func TestWalletLogUpdate(t *testing.T) {
	id := 158113
   repo :=	ti.Factory.GetWalletRepo()
	l := repo.GetWalletLog_(id)
	l.Title = l.Title+"_1"
	_,err := repo.SaveWalletLog_(l)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	time.Sleep(time.Second*2)
}
