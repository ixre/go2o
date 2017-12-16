package kdniao

import "testing"

func TestKdnTraces(t *testing.T) {
	EBusinessID = "1314567"
	AppKey = "27d809c3-51b6-479c-9b77-6b98d7f3d414"
	v, err := KdnTraces("ZTO", "462681586678")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("value=%#v", v)
}
