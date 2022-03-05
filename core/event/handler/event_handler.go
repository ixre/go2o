package handler

import (
	"context"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/repos/clickhouse"
	"log"
)

type EventHandler struct {
}

func (h EventHandler) HandleWalletLogWriteEvent(data interface{}) {
	ld := data.(*events.WalletLogClickhouseWriteEvent)
	conn := clickhouse.ConnInstance
	batch, err := conn.PrepareBatch(context.TODO(),
		`INSERT INTO go2o_wal_wallet_log (  
id,wallet_id,wallet_user,kind,title,outer_chan,
outer_no,value,balance,procedure_fee,
opr_uid,opr_name,account_no,
account_name,bank_name,review_state,
review_remark,review_time,remark,create_time,
update_time)`)
	if err != nil {
		log.Println("[ event]: handle wallet log write error", err.Error())
		return
	}
	l := ld.Data
	if err := batch.Append(
		l.Id, l.WalletId, l.WalletUser,
		int32(l.Kind),
		l.Title,
		l.OuterChan,
		l.OuterNo,
		l.Value,
		l.Balance,
		int32(l.ProcedureFee),
		int64(l.OperatorUid),
		l.OperatorName,
		l.AccountNo,
		l.AccountName,
		l.BankName,
		int32(l.ReviewState),
		l.ReviewRemark,
		l.ReviewTime,
		l.Remark,
		l.CreateTime,
		l.UpdateTime,
	); err != nil {
		log.Println("[ event]: handle wallet log write error", err.Error())
		return
	}
	if err := batch.Send(); err != nil {
		log.Println("[ event]: handle wallet log write error", err.Error())
		return
	}
}
