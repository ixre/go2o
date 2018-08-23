ALTER TABLE `pay_order`
CHANGE COLUMN `pay_flag` `pay_flag` INT(6) NOT NULL COMMENT '可⽤支付方式' ,
ADD COLUMN `paid_fee` INT(10) NOT NULL COMMENT '实付金额' AFTER `final_fee`,
ADD COLUMN `final_flag` INT(6) NOT NULL COMMENT '最终使用支付方式' AFTER `pay_flag`;

ALTER TABLE `pay_trade_chan`
CHANGE COLUMN `chan_data` `pay_code` VARCHAR(40) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL COMMENT '支付代码' AFTER `pay_method`,
CHANGE COLUMN `trade_no` `trade_no` VARCHAR(40) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL COMMENT '支付单号' ,
CHANGE COLUMN `pay_chan` `pay_method` INT(6) NOT NULL COMMENT '支付方式' ,
CHANGE COLUMN `internal_chan` `internal` INT(1) NOT NULL COMMENT '是否为内置支付方式' ,
ADD COLUMN `out_trade_no` VARCHAR(45) NOT NULL COMMENT '外部交易单号' AFTER `pay_amount`, RENAME TO  `pay_trade_data` ;

ALTER TABLE `pay_trade_data`
ADD COLUMN `pay_time` INT(11) NOT NULL COMMENT '支付时间' AFTER `out_trade_no`;
