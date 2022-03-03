```
CREATE DATABASE [IF NOT EXISTS] go2o;
use go2o;

CREATE TABLE IF NOT EXISTS wal_wallet_log
(
    `id` Int64 COMMENT '编号',
    `wallet_id` Int64 COMMENT '钱包编号',
    `kind` Int32 COMMENT '业务类型',
    `title` String COMMENT '标题',
    `outer_chan` String COMMENT '外部通道',
    `outer_no` String COMMENT '外部订单号',
    `value` Int64 COMMENT '变动金额',
    `balance` Int64 COMMENT '余额',
    `trade_fee` Int32 COMMENT '交易手续费',
    `opr_uid` Int32 COMMENT '操作人员用户编号',
    `opr_name` String COMMENT '操作人员名称',
    `account_no` String COMMENT '提现账号',
    `account_name` String COMMENT '提现银行账户名称',
    `bank_name` String COMMENT '提现银行',
    `review_state` Int32 COMMENT '审核状态',
    `review_remark` String COMMENT '审核备注',
    `review_time` Int64 COMMENT '审核时间',
    `remark` String COMMENT '备注',
    `create_time` Int64 COMMENT '创建时间',
    `update_time` Int64 COMMENT '更新时间'
) ENGINE = MergeTree
ORDER BY id
SETTINGS index_granularity= 8192 ;
```

