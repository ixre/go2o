ALTER TABLE `pay_order`
CHANGE COLUMN `pay_flag` `pay_flag` INT(6) NOT NULL COMMENT '可⽤支付方式' ,
ADD COLUMN `paid_fee` INT(10) NOT NULL COMMENT '实付金额' AFTER `final_fee`,
ADD COLUMN `final_flag` INT(6) NOT NULL COMMENT '最终使用支付方式' AFTER `pay_flag`;
