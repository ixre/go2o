package tests

import (
	"errors"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/merchant/wholesaler"
	"go2o/core/infrastructure/domain"
	"go2o/tests/ti"
	"testing"
)

func TestMerchantPwd2(t *testing.T) {
	s := domain.Md5("123456")
	println(domain.Sha1Pwd(s))
}

// 测试创建商户
func TestCreateMerchant(t *testing.T) {
	repo := ti.Factory.GetMerchantRepo()
	v := &merchant.Merchant{
		LoginUser:   "zy",
		LoginPwd:    domain.Md5("123456"),
		Name:        "天猫",
		SelfSales:   1,
		Level:       0,
		Logo:        "",
		CompanyName: "天猫",
		Province:    0,
		City:        0,
		District:    0,
	}
	v.LoginPwd = domain.Sha1Pwd(v.LoginPwd)
	im := repo.CreateMerchant(v)
	err := im.SetValue(v)
	if err == nil {
		_, err = im.Save()
		if err == nil {
			o := shop.OnlineShop{
				ShopName:   v.Name,
				Logo:       "https://raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif",
				Host:       "",
				Alias:      "zy",
				Tel:        "",
				Addr:       "",
				ShopTitle:  "",
				ShopNotice: "",
			}
			_, err = im.ShopManager().CreateOnlineShop(&o)
		}
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 测试商家分组设置
func TestMchBuyerGroupSet(t *testing.T) {
	repo := ti.Factory.GetMerchantRepo()
	conf := repo.GetMerchant(1).ConfManager()
	g := conf.GetGroupByGroupId(2)
	g.Alias = "VIP1"
	g.EnableRetail = 1
	g.EnableWholesale = 0
	_, err := conf.SaveMchBuyerGroup(g)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	groups := conf.SelectBuyerGroup()
	for _, v := range groups {
		if v.Name == g.Alias {
			return
		}
	}
	t.Error("没有获取到名称号为:VIP1的分组")
	t.Fail()
}

// 测试成为批发商
func TestStartMchWholesale(t *testing.T) {
	mch := ti.Factory.GetMerchantRepo().GetMerchant(1)
	err := mch.EnableWholesale()
	if err == nil {
		ws := mch.Wholesaler()
		if ws == nil {
			err = errors.New("become wholesaler failed.")
		} else {
			err = ws.Review(true, "")
			return
		}
	}
	t.Error(err)
	t.Fail()
}

// 测试设置返点比例
func TestGroupRebateRate(t *testing.T) {
	mmRepo := ti.Factory.GetMemberRepo()
	groups := mmRepo.GetManager().GetAllBuyerGroups()
	mch := ti.Factory.GetMerchantRepo().GetMerchant(1)
	ws := mch.Wholesaler()
	if ws == nil {
		t.Error("merchant is not a wholesaler!")
		t.Fail()
	}
	for _, v := range groups {
		r := []*wholesaler.WsRebateRate{
			{
				RequireAmount: 2,
				RebateRate:    0.1,
			},
			{
				RequireAmount: 30,
				RebateRate:    0.15,
			},
		}
		err := ws.SaveGroupRebateRate(v.ID, r)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}
	groupId := groups[0].ID
	//计算商品订单的折扣率
	rate1 := ws.GetRebateRate(groupId, 1)
	if rate1 != 0 {
		t.Error("价格应无折扣,实际:", rate1)
		t.Fail()
	}
	rate2 := ws.GetRebateRate(groupId, 2)
	if rate2 != 0.1 {
		t.Error("折扣不正确,实际:", rate2)
		t.Fail()
	}
	rate3 := ws.GetRebateRate(groupId, 31)
	if rate3 != 0.15 {
		t.Error("折扣不正确,实际:", rate3)
		t.Fail()
	}
}

// 测试结算订单到账户中
func TestMchSettleOrder(t *testing.T) {
	repo := ti.Factory.GetMerchantRepo()
	mch := repo.GetMerchant(111)
	err := mch.Account().SettleOrder("123", 1000, 20, 0, "零售订单结算")
	if err != nil {
		t.Log("结算订单出错：", err)
		t.FailNow()
	}
	t.Log("结算成功")
}
