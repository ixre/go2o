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
	msq.Configure(msq.NATS, []string{"127.0.0.1:4222"})
	//msq.Configure(msq.NATS, []string{"www.dev1.super4bit:4222"})
}

func TestMemberUpdate(t *testing.T) {
	defer msq.Close()
	msq.Push(msq.MemberUpdated, "update!"+ strconv.Itoa(id))
	msq.PushDelay(msq.MemberAccountUpdated, strconv.Itoa(id), "", 1000)
	msq.PushDelay(msq.MemberProfileUpdated, strconv.Itoa(id), "", 1000)
	msq.PushDelay(msq.MemberRelationUpdated, strconv.Itoa(id), "", 1000)
	time.Sleep(5 * time.Second)
}

func TestMemberTrustPassedMQ(t *testing.T) {
	defer msq.Close()
	msq.Push(msq.MemberTrustInfoPassed,
		fmt.Sprintf("%d|%d|%s|%s",
			id,1, "513701981888455487", "刘铭"))
	time.Sleep(5 * time.Second)
}
