package domain

import (
	"errors"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/wholesaler"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/initial/provide"
	"github.com/ixre/go2o/core/inject"
	_ "github.com/ixre/go2o/tests"
)

func TestMerchantPwd2(t *testing.T) {
	s := domain.Md5("123456")
	println(domain.MerchantSha265Pwd(s, ""))
}

// 测试创建商户
func TestCreateMerchant(t *testing.T) {
	repo := inject.GetMerchantRepo()
	v := &merchant.Merchant{
		MailAddr: "zy11@fze.net",
		MchName:  "天猫B",
		Salt:     "000",
		IsSelf:   1,
		Level:    0,
		Logo:     "",
		Province: 110000,
		City:     110000,
		District: 110000,
	}
	v.Username = v.MailAddr
	v.Password = domain.MerchantSha265Pwd(domain.Md5("123456"), v.Salt)
	im := repo.CreateMerchant(v)
	err := im.SetValue(v)
	if err == nil {
		_, err = im.Save()
		// 商户要求认证后，才能开通
		// if err == nil {
		// 	o := shop.OnlineShop{
		// 		ShopName:   v.MchName,
		// 		Logo:       "https://raw.githubusercontent.com/jsix/go2o/master/docs/mark.gif",
		// 		Host:       "",
		// 		Alias:      "zy",
		// 		Telephone:  "",
		// 		Addr:       "",
		// 		ShopTitle:  "",
		// 		ShopNotice: "",
		// 	}
		// 	_, err = im.ShopManager().CreateOnlineShop(&o)
		// }
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// 保存商户认证信息
func TestSaveMerchantAuthenticate(t *testing.T) {
	mch := inject.GetMerchantRepo().GetMerchant(21)
	v := &merchant.Authenticate{
		OrgName:          "中国天空有限公司",
		MchName:          "中的国际",
		LicenceNo:        "00000000",
		LicencePic:       "https://so1.360tres.com/dr/220__/t0146eaced4b2c0a82d.jpg",
		WorkCity:         0,
		QualificationPic: "https://so1.360tres.com/dr/220__/t0146eaced4b2c0a82d.jpg",
		PersonId:         "513701980102345678",
		PersonName:       "田猫",
		PersonFrontPic:   "https://so1.360tres.com/dr/220__/t0146eaced4b2c0a82d.jpg",
		PersonBackPic:    "https://so1.360tres.com/dr/220__/t0146eaced4b2c0a82d.jpg",
		PersonPhone:      "13888888888",
		AuthorityPic:     "https://so1.360tres.com/dr/220__/t0146eaced4b2c0a82d.jpg",
		BankName:         "花旗银行",
		BankAccount:      "田猫",
		BankNo:           "622601345897234",
		City:             110000,
		District:         110101,
		Province:         110000,
		ContactName:      "田猫",
		ContactPhone:     "13888888888",
	}
	id, err := mch.ProfileManager().SaveAuthenticate(v)
	if err != nil {
		t.Error(err)
	}
	t.Logf("id:%d", id)
}

// 拒绝审核商户认证信息
func TestRejectMerchantAuthenticateRequest(t *testing.T) {
	mch := inject.GetMerchantRepo().GetMerchant(1)
	err := mch.ProfileManager().ReviewAuthenticate(false, "不通过")
	if err != nil {
		t.Error(err)
	}
	// 再次审核
	err = mch.ProfileManager().ReviewAuthenticate(false, "不通过")
	if err == nil {
		t.Error(errors.New("再次审核未提示错误"))
	}
}

// 测试审核通过商户认证信息
func TestPassMerchantAuthenticateRequest(t *testing.T) {
	/**
	SQL:
	select * FROM mch_merchant where id=1;
	select * FROM mch_authenticate WHERE mch_id=1;
	*/
	mch := inject.GetMerchantRepo().GetMerchant(1)
	err := mch.ProfileManager().ReviewAuthenticate(true, "")
	if err == nil {
		// 审核通过
		err = mch.ProfileManager().ReviewAuthenticate(true, "通过")
		if err != nil {
			t.Error(errors.New("审核失败:" + err.Error()))
		}
	} else {
		t.Error(errors.New("审核失败:" + err.Error()))
	}
}

// 测试更改绑定会员
func TestBindMember(t *testing.T) {
	var mchId = 1
	var memberId = 4
	mch := inject.GetMerchantRepo().GetMerchant(mchId)
	err := mch.BindMember(memberId)
	if err == nil {
		err = mch.BindMember(memberId + 1)
		if err == nil {
			err = mch.BindMember(memberId)
		}
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if mch.GetValue().MemberId != memberId {
		t.Log("now bind member id is ", mch.GetValue().MemberId)
		t.FailNow()
	}
}

// 测试商家分组设置
func TestMchBuyerGroupSet(t *testing.T) {
	repo := inject.GetMerchantRepo()
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
	mch := inject.GetMerchantRepo().GetMerchant(1)
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
	mmRepo := inject.GetMemberRepo()
	groups := mmRepo.GetManager().GetAllBuyerGroups()
	mch := inject.GetMerchantRepo().GetMerchant(1)
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
func TestMchCarry(t *testing.T) {
	repo := inject.GetMerchantRepo()
	mch := repo.GetMerchant(1)
	sd := merchant.CarryParams{
		OuterTxNo:         "TS:202407241000001",
		Amount:            10000,
		TransactionFee:    1000,
		RefundAmount:      0,
		TransactionTitle:  "测试订单结算",
		TransactionRemark: "虚拟订单",
	}
	txId, err := mch.Account().Carry(sd)
	if err != nil {
		t.Log("结算订单出错：", err)
		t.FailNow()
	}
	t.Logf("结算成功,交易流水号:%d", txId)
}

// 测试员工转移
func TestTransferStaff(t *testing.T) {
	staffId := 2
	repo := inject.GetMerchantRepo()
	mch := repo.GetMerchant(1)
	staffManager := mch.EmployeeManager()
	transferId, err := staffManager.RequestTransfer(staffId, 2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("转移ID:%d", transferId)
	staffRepo := inject.GetStaffRepo()
	transfer := staffRepo.TransferRepo().Get(transferId)
	// 进行审批
	approvalRepo := inject.GetApprovalRepo()
	ia := approvalRepo.GetApproval(transfer.ApprovalId)
	err = ia.Approve()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = ia.Approve()
	//err = ia.Reject("测试")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("审批完成")
}

// 测试商户交易账单
func TestMerchantTransactionBill(t *testing.T) {
	repo := inject.GetMerchantRepo()
	mch := repo.GetMerchant(1)
	tx := mch.TransactionManager().GetCurrentDailyBill()
	err := tx.UpdateBillAmount(10000, 1000)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = tx.UpdateAmount()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = tx.Generate()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = tx.Confirm()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = tx.Review(1, "测试")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if tx.Value().BillType == merchant.BillTypeMonthly {
		// 结算月度账单
		err = tx.Settle()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		err = tx.UpdateSettleInfo("SP000001", "TX000001", merchant.BillResultSuccess, "success")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		t.Logf("结算成功")
	}
}

func TestBatchChangeMerchantPassword(t *testing.T) {
	db := provide.GetGOrm()
	var merchantList []*merchant.Merchant
	tx := db.Model(&merchant.Merchant{})
	tx.Scan(&merchantList)
	for _, v := range merchantList {
		v.Password = domain.MemberSha256Pwd(domain.Md5("123456"), v.Salt)
		db.Save(v)
	}
}
