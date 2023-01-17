package handler

import (
	"context"
	"log"

	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/repos/clickhouse"
)

func (h EventHandler) HandleWalletLogWriteEvent(data interface{}) {
	ld := data.(*events.AccountLogPushEvent)
	conn := clickhouse.GetClickhouseConn()
	if ld == nil || conn == nil {
		return
	}
	if err := conn.Exec(context.TODO(),
		`ALTER TABLE go2o_wal_wallet_log UPDATE 
		title= $1, 
		outer_no=$2,
		
		review_state = $8,
		review_remark = $9,
		review_time = $10,
		remark = $11,
		update_time = $12 WHERE 
		wallet_id = $13 AND id= $14
	`,
		ld.Subject,
		ld.OuterNo,
		int32(ld.AuditState),
		ld.CreateTime,
		ld.MemberId,
		ld.LogId,
	); err != nil {
		log.Println("[ event]: update wallet log error", err.Error())
		return
	}
}
