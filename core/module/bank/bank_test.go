package bank

import (
	"encoding/json"
	"testing"
)

func TestSimpleJson(t *testing.T) {

	result := `{
  "code": "10000",
  "charge": false,
  "msg": "查询成功",
  "result": {
    "error_code": 0,
    "reason": "成功",
    "result": {
      "accountNo": "6228480500861233451",
      "name": "张三",
      "idCardCore": "130321198804010180",
      "bankPreMobile": "18600174444",
      "result": "T",
      "message": "认证信息匹配",
      "messagetype": 0
    }
  }
}`

	mp := make(map[string]interface{})
	json.Unmarshal([]byte(result), &mp)
	r1 := mp["result"].(map[string]interface{})
	r2 := r1["result"].(map[string]interface{})
	t.Log(r2["message"])
	t.Log(r2["result"])
	t.Log(r2["messagetype"].(float64))
}

func TestGetNameByAccountNo(t *testing.T) {
	accountNo := "6226220284294245"
	accountNo1 := "6222021001042791910"
	accountNo2 := "6229332000010155164"
	accountNo3 := "5229640795589453"
	t.Log(GetNameByAccountNo(accountNo))
	t.Log(GetNameByAccountNo(accountNo1))
	t.Log(GetNameByAccountNo(accountNo2))
	t.Log(GetNameByAccountNo(accountNo3))
}
