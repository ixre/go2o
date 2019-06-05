package tests

import (
	"go2o/core/domain/interface/member"
	"go2o/core/msq"
	"go2o/tests/ti"
	"strconv"
	"testing"
	"time"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : batch_member_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-06-05 01:02
 * description :
 * history :
 */
func init() {
	// 初始化producer
	msq.Configure(msq.KAFKA, []string{"127.0.0.1:9092"})
}

func TestBatchPushMember(t *testing.T){
	defer msq.Close()
	orm := ti.GetApp().Db().GetOrm()

	var members []member.Member
	err := orm.SelectByQuery(&members,"select * FROM mm_member where id > 0 LIMIT 100 OFFSET 0")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	for _,v := range members{
		id := int(v.Id)
		msq.Push(msq.MemberUpdated, strconv.Itoa(id), "update")
		msq.PushDelay(msq.MemberAccountUpdated, strconv.Itoa(id), "", 1000)
		msq.PushDelay(msq.MemberProfileUpdated, strconv.Itoa(id), "", 1000)
		msq.PushDelay(msq.MemberRelationUpdated, strconv.Itoa(id), "", 1000)
		t.Log("notify ",id)
	}
	t.Log("finished")
	time.Sleep(5 * time.Minute)
}