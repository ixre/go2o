// create for src 29/11/2017 ( jarrysix@gmail.com )
package tests

import (
	"testing"

	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/inject"
)

const walletId int64 = 1

// 测试创建钱包
func TestCreateWallet(t *testing.T) {
	repo := inject.GetWalletRepo()
	wlt := repo.GetWallet(int(walletId))
	if wlt == nil {
		wlt = repo.CreateWallet(1, "", wallet.TMerchant, "商户钱包", wallet.FlagCharge|wallet.FlagDiscount)
	}
	id, err := wlt.Save()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("id=", id)
}

// 测试充值钱包
func TestChargeWallet(t *testing.T) {
	repo := inject.GetWalletRepo()
	wlt := repo.GetWallet(int(walletId))
	totalCharge := wlt.Get().TotalCharge
	err := wlt.Charge(100000, wallet.CServiceAgentCharge,
		"客服充值", "", "1234", 1, "洛洛")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = wlt.Charge(50000, wallet.CUserCharge,
		"用户充值", "-", "", 0, "洛洛")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if wlt.Get().TotalCharge-totalCharge != 150000 {
		t.Log("原先累计充值:", totalCharge, "现在累计重置:", wlt.Get().TotalCharge)
		t.FailNow()
	}
	t.Log("余额=", wlt.Get().Balance)
}

// 测试钱包支付和退款
func TestDiscountRefundWallet(t *testing.T) {
	repo := inject.GetWalletRepo()
	wlt := repo.GetWallet(int(walletId))
	var value = 10000
	var tradeNo = "02af1208xa209sl2"
	var balance = wlt.Get().Balance
	_, err := wlt.Discount(value, "支付订单"+tradeNo, tradeNo, true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = wlt.Refund(-value, wallet.KPaymentOrderRefund, "订单退款", tradeNo, 0, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if v := wlt.Get().Balance; v != balance {
		t.Error("退款不成功", balance, v)
		t.FailNow()
	}

}

// 测试冻结钱包
func TestFreezeWallet(t *testing.T) {
	repo := inject.GetWalletRepo()
	wlt := repo.GetWallet(int(walletId))
	var value = 10000
	var freeze = wlt.Get().FreezeAmount
	var balance = wlt.Get().Balance
	_, err := wlt.Freeze(wallet.TransactionData{
		TransactionTitle:  "冻结金额",
		Amount:            value,
		OuterTxNo:         "",
		TransactionRemark: "",
	}, wallet.Operator{OperatorUid: 0})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("冻结金额=", wlt.Get().FreezeAmount)
	err = wlt.Unfreeze(value, "解冻金额", "", false, 0, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if v := wlt.Get().FreezeAmount; v != freeze {
		t.Error("解冻不正确", freeze, v)
		t.FailNow()
	}
	if v := wlt.Get().Balance; v != balance {
		t.Error("解冻后余额不正确", balance, v)
		t.FailNow()
	}
}

// 测试调整钱包金额
func TestAdjustWallet(t *testing.T) {
	repo := inject.GetWalletRepo()
	wlt := repo.GetWallet(int(walletId))
	adjust := wlt.Get().AdjustAmount
	err := wlt.Adjust(1000, "客服调整",
		"-", "", 2, "TOM")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = wlt.Adjust(-1000, "客服取消调整",
		"-", "", 2, "TOM")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if v := wlt.Get().AdjustAmount; v != adjust {
		t.Log("调整后金额不正确", adjust, v)
		t.FailNow()
	}
}

// 测试提现失败
func TestTakeOutWalletFail(t *testing.T) {
	repo := inject.GetWalletRepo()
	wlt := repo.GetWallet(int(walletId))
	var amount = 10000
	balance := wlt.Get().Balance
	id, _, err := wlt.RequestWithdrawal(wallet.WithdrawTransaction{
		Amount:           -amount,
		TransactionFee:   200,
		Kind:             wallet.KWithdrawToBankCard,
		TransactionTitle: "提现到银行卡",
		AccountNo:        "6226220285888888",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if v := wlt.Get().Balance; v != balance-int(amount) {
		t.Error("提现扣款不正确", balance, v)
		t.FailNow()
	}
	err = wlt.ReviewWithdrawal(id, false, "银行卡号不正确", 1, "管理员")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if v := wlt.Get().Balance; v != balance {
		t.Error("提现退回后余额不正确", balance, v)
		t.FailNow()
	}
}

// 测试提现失败
func TestTakeOutWalletSuccess(t *testing.T) {
	repo := inject.GetWalletRepo()
	wlt := repo.GetWallet(int(walletId))
	var amount = 10000
	balance := wlt.Get().Balance
	id, _, err := wlt.RequestWithdrawal(wallet.WithdrawTransaction{
		Amount:           -amount,
		TransactionFee:   200,
		Kind:             wallet.KWithdrawToBankCard,
		TransactionTitle: "提现到银行卡",
		AccountNo:        "6226220285888888",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if v := wlt.Get().Balance; v != balance-int(amount) {
		t.Error("提现扣款不正确", balance, v)
		t.FailNow()
	}
	err = wlt.ReviewWithdrawal(id, true, "", 1, "管理员")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = wlt.CompleteTransaction(id, "96699999999")
	if v := wlt.Get().Balance; v != balance-int(amount) {
		t.Error("提现退回后余额不正确", balance, v)
		t.FailNow()
	}
}

// 测试转账
func TestTransferWallet(t *testing.T) {
	repo := inject.GetWalletRepo()
	wlt := repo.GetWallet(int(walletId))
	var amount = 10000
	var transactionFee = 1000
	var toWalletId int64 = 2
	var balance2 int64 = 0
	wlt2 := repo.GetWallet(int(toWalletId))
	if wlt2 == nil {
		t.Error("目标账户不存在")
		t.FailNow()
	} else {
		balance2 = int64(wlt2.Get().Balance)
	}
	balance := int64(wlt.Get().Balance)
	err := wlt.Transfer(toWalletId, amount, -transactionFee, "转账给2", "收款1", "给你发个红包")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if v := wlt.Get().Balance; v != int(balance)-int(amount)-int(transactionFee) {
		t.Error("转账扣款不正确", balance-int64(amount)-int64(transactionFee), v)
		t.FailNow()
	}
	wlt2 = repo.GetWallet(int(toWalletId))
	if v := wlt2.Get().Balance; v != int(balance2)+int(amount) {
		t.Error("转账收款不正确", balance2, int(v)-int(amount))
	}

}
