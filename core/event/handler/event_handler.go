package handler

import (
	"context"
	"log"

	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/repos/clickhouse"
)

type EventHandler struct {
}

func (h EventHandler) HandleWalletLogWriteEvent(data interface{}) {
	ld := data.(*events.WalletLogClickhouseUpdateEvent)
	conn := clickhouse.GetClickhouseConn()
	if ld == nil || conn == nil {
		return
	}
	l := ld.Data
	if err := conn.Exec(context.TODO(),
		`ALTER TABLE go2o_wal_wallet_log UPDATE 
		title= $1, 
		outer_no=$2,
		opr_uid = $3,
		opr_name = $4,
		account_no = $5,
		account_name = $6,
		bank_name = $7,
		review_state = $8,
		review_remark = $9,
		review_time = $10,
		remark = $11,
		update_time = $12 WHERE 
		wallet_id = $13 AND id= $14
	`,
		l.Title,
		l.OuterNo,
		int64(l.OperatorUid),
		l.OperatorName,
		l.AccountNo,
		l.AccountName,
		l.BankName,
		int32(l.AuditState),
		l.AuditRemark,
		l.AuditTime,
		l.Remark,
		l.UpdateTime,
		l.WalletId,
		l.Id,
	); err != nil {
		log.Println("[ event]: update wallet log error", err.Error())
		return
	}
}
