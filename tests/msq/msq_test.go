package msq

import (
	_ "github.com/ixre/go2o/tests"
)

const id = 1

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
