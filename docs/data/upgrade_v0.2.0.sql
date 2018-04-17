

/*ALTER TABLE mm_trusted_info
  DROP COLUMN reviewed;*/
ALTER TABLE mm_trusted_info
  modify column real_name varchar(10) NOT NULL;
ALTER TABLE mm_trusted_info
  ADD COLUMN card_type int(1) NOT NULL comment '证件类型'
ALTER TABLE mm_trusted_info
  ADD COLUMN card_area varchar(5) NOT NULL comment '证件区域';
ALTER TABLE mm_trusted_info
  modify column card_id varchar(20) NOT NULL;
ALTER TABLE mm_trusted_info
  ADD COLUMN card_image varchar(120) NOT NULL comment '证件图像';
ALTER TABLE mm_trusted_info
  modify column trust_image varchar(120) NOT NULL;
ALTER TABLE mm_trusted_info
  ADD COLUMN review_state int(1) NOT NULL comment '审核状态';
ALTER TABLE mm_trusted_info
  modify column review_time int(11) NOT NULL;
ALTER TABLE mm_trusted_info
  ADD COLUMN manual_review int(1) NOT NULL comment '是否人工认证';
ALTER TABLE mm_trusted_info
  modify column remark varchar(120) NOT NULL;
ALTER TABLE mm_trusted_info
  modify column update_time int(11) NOT NULL;

/* 更新会员等级 */

ALTER TABLE mm_bank 
  modify column name varchar(45) NOT NULL;
ALTER TABLE mm_bank
  modify column account varchar(45) NOT NULL;
ALTER TABLE mm_bank 
  modify column account_name varchar(45) NOT NULL;
ALTER TABLE mm_bank
  modify column network varchar(45) NOT NULL;
ALTER TABLE mm_bank
  modify column is_locked int(1) NOT NULL;
ALTER TABLE mm_bank
  modify column state int(11) NOT NULL;
ALTER TABLE mm_bank 
  modify column update_time int(11) NOT NULL;
ALTER TABLE mm_favorite 
  modify column member_id int(11) NOT NULL;
ALTER TABLE mm_favorite 
  modify column fav_type int(1) NOT NULL;
ALTER TABLE mm_favorite 
  modify column refer_id int(11) NOT NULL;
ALTER TABLE mm_favorite 
  modify column update_time int(11) NOT NULL;
ALTER TABLE mm_income_log 
  modify column order_id int(11) NOT NULL;
ALTER TABLE mm_income_log 
  modify column member_id int(11) NOT NULL;
ALTER TABLE mm_income_log 
  modify column type varchar(10) NOT NULL;
ALTER TABLE mm_income_log 
  modify column fee float NOT NULL;
ALTER TABLE mm_income_log 
  modify column log varchar(100) NOT NULL;
ALTER TABLE mm_income_log 
  modify column record_time int(11) NOT NULL;
ALTER TABLE mm_income_log 
  modify column state int(11) NOT NULL;
ALTER TABLE mm_integral_log 
  modify column outer_no varchar(45) NOT NULL;
ALTER TABLE mm_integral_log 
  modify column remark varchar(100) NOT NULL;
ALTER TABLE mm_level 
  modify column name varchar(45) NOT NULL;
ALTER TABLE mm_level 
  modify column require_exp int(11) NOT NULL;
ALTER TABLE mm_level 
  modify column program_signal varchar(45) NOT NULL;
ALTER TABLE mm_level 
  modify column is_official int(1) NOT NULL;
ALTER TABLE mm_level 
  ADD COLUMN allow_upgrade int(1) NOT NULL comment '允许自动升级';
ALTER TABLE mm_level 
  modify column enabled int(1) NOT NULL;
ALTER TABLE mm_levelup 
  ADD COLUMN upgrade_type int(1) NOT NULL comment '升级方式1:自动升级 2:客服更改 3:系统升级';
ALTER TABLE mm_relation 
  modify column card_no varchar(20) NOT NULL;
ALTER TABLE mm_relation 
  modify column inviter_str varchar(250) NOT NULL;
ALTER TABLE mm_relation 
  modify column reg_mchid int(11) NOT NULL;
ALTER TABLE wal_wallet 
  alter column balance set default 0.00;
ALTER TABLE wal_wallet 
  alter column total_charge set default 0.00;
ALTER TABLE wal_wallet 
  alter column total_pay set default 0.00;


ALTER TABLE `pay_order`
CHANGE COLUMN `total_fee` `total_amount` DECIMAL(8,2) NULL DEFAULT NULL ,
CHANGE COLUMN `payment_opt` `payment_opt` INT(2) NULL DEFAULT NULL ,
CHANGE COLUMN `payment_sign` `payment_sign` INT(2) NULL DEFAULT NULL ,
CHANGE COLUMN `state` `state` INT(1) NULL DEFAULT NULL ,
ADD COLUMN `trade_type` VARCHAR(20) NULL AFTER `order_type`;

ALTER TABLE `mch_online_shop`
CHANGE COLUMN `host` `host` VARCHAR(40) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL COMMENT '主机头' ;
