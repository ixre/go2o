package domain

import (
	"strconv"
	"testing"
	"time"

	"github.com/ixre/gof/util"
)

func TestGenerateAppId(t *testing.T) {
	for {
		s := strconv.Itoa(util.RandInt(8))
		t.Log(s)
		time.Sleep(1000)
	}
}
