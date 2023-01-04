package events

import (
	"testing"
	"time"

	"github.com/ixre/go2o/tests/ti"
)

func TestWalletLogUpdate(t *testing.T) {
	id := 158113
	repo := ti.Factory.GetWalletRepo()
	l := repo.GetWalletLog_(id)
	l.Subject = l.Subject + "_1"
	_, err := repo.SaveWalletLog_(l)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	time.Sleep(time.Second * 2)
}
