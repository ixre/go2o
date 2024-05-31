package msq

import (
	"github.com/ixre/go2o/core/event/msq"
	_ "github.com/ixre/go2o/tests"
)

const id = 1

func init() {
	// 初始化producer
	msq.Configure(msq.NATS, []string{"127.0.0.1:4222"})
	//msq.Configure(msq.NATS, []string{"www.dev1.super4bit:4222"})
}

// func TestMemberTrustPassedMQ(t *testing.T) {
// 	defer msq.Close()
// 	err := msq.Push(msq.MemberTrustInfoPassed,
// 		fmt.Sprintf("%d|%d|%s|%s",
// 			id, 1, "513701981888455487", "刘铭"))
// 	if err != nil {
// 		t.Log("--", err)
// 	}
// 	time.Sleep(5 * time.Second)
// }
