package uams

import "testing"

func init() {
	API_SERVER = "http://localhost:1419/uams_api_v1"
	API_USER = "3722442566017024"
	API_TOKEN = "3135303037303032383326d878ff2442"
	API_SIGN_TYPE = "md5"
}

func TestGetAppInfo(t *testing.T) {
	rsp, err := Post("app.info", nil)
	t.Log("Response:", rsp)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
