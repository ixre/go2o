package uams

import "testing"

func init() {
	API_SERVER = "http://localhost:1419/uams_api_v1"
	API_USER = "3722442566017024"
	API_TOKEN = "3135303037303032383326d878ff2442"
	API_APP = "4297e4cf-2a68-46cb-b88a-de2e62d0f06c"
	API_SIGN_TYPE = "md5"
}

func TestGetAppInfo(t *testing.T) {
	rsp, err := Post("app.info", nil)
	t.Log("Response:", string(rsp))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
