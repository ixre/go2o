
CREATE TABLE mch_trade_conf (
  id           int(10) NOT NULL AUTO_INCREMENT comment '编号',
  mch_id       int(10) NOT NULL comment '商户编号',
  trade_type   int(1) NOT NULL comment '交易类型',
  plan_id      int(10) NOT NULL comment '交易方案，根据方案来自动调整比例',
  flag         int(4) NOT NULL comment '交易标志',
  amount_basis int(1) NOT NULL comment '交易手续费依据,1:未设置 2:按金额 3:按比例',
  trade_fee    int(4) NOT NULL comment '交易费，按单笔收取',
  trade_rate   int(4) NOT NULL comment '交易手续费比例',
  update_time  int(11) NOT NULL comment '更新时间',
  PRIMARY KEY (id)) comment='商户交易设置';

