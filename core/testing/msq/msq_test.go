package msq

import (
	"fmt"
	"go2o/core/msq"
	"strconv"
	"testing"
	"time"
)

const id = 22149

func init() {
	// 初始化producer
	msq.Configure(msq.KAFKA, []string{"127.0.0.1:9092"})
}

func TestMemberUpdate(t *testing.T) {
	defer msq.Close()
	msq.Push(msq.MemberUpdated, strconv.Itoa(id), "update")
	time.Sleep(time.Second)
	msq.Push(msq.MemberAccountUpdated, strconv.Itoa(id), "")
	time.Sleep(time.Second)
	msq.Push(msq.MemberProfileUpdated, strconv.Itoa(id), "")
	time.Sleep(time.Second)
	msq.Push(msq.MemberRelationUpdated, strconv.Itoa(id), "")
	time.Sleep(5 * time.Second)
}

func TestMemberTrustPassedMQ(t *testing.T) {
	defer msq.Close()
	msq.Push(msq.MemberTrustInfoPassed, strconv.Itoa(id),
		fmt.Sprintf("%d|%s|%s",
			1, "513701981888455487", "刘铭"))
	time.Sleep(5 * time.Second)
}
