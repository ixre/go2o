package query

import (
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

type WalletQuery struct {
	walletRepo    fw.BaseRepository[wallet.Wallet]
	walletLogRepo fw.BaseRepository[wallet.WalletLog]
}

func NewWalletQuery(o fw.ORM) *WalletQuery {
	w := &WalletQuery{}
	w.walletRepo.ORM = o
	w.walletLogRepo.ORM = o
	return w
}

func (m *WalletQuery) GetWalletId(userId int, walletType int) int {
	v := m.walletRepo.FindBy("user_id = ? and wallet_type = ?", userId, walletType)
	if v != nil {
		return int(v.Id)
	}
	return 0
}

// QueryPagingAccountLog 查询商户账户钱包明细
func (m *WalletQuery) QueryMerchantPagingAccountLog(walletId int, p *fw.PagingParams) (*fw.PagingResult, error) {
	p.Equal("wallet_id", walletId)
	p.OrderBy("create_time desc")
	return m.walletLogRepo.QueryPaging(p)
}

// 查询总收入金额
func (m *WalletQuery) QueryTotalCarryAmount(walletId int) int {
	count, _ := m.walletLogRepo.Count("wallet_id = ? and kind = ?",
		walletId, wallet.KCarry)
	return count
}

// 查询月度总收入金额
func (m *WalletQuery) QueryMonthCarryAmount(walletId int, unix int) int {
	count, _ := m.walletLogRepo.Count(`wallet_id = ? and kind = ?
	AND DATE_TRUNC('month',to_timestamp(create_time)) = DATE_TRUNC('month',?)`,
		walletId, wallet.KCarry, unix)
	return count
}
