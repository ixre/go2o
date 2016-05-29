ALTER TABLE `zxdb`.`pt_merchant`
RENAME TO  `zxdb`.`pt_merchant` ;
ALTER TABLE `zxdb`.`pt_page`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`dlv_partner_bind`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL , RENAME TO  `zxdb`.`dlv_merchant_bind` ;
ALTER TABLE `zxdb`.`gs_category`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL COMMENT '商户ID(pattern ID);如果为空，则表示模式分类' ;
ALTER TABLE `zxdb`.`gs_sale_tag`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`mm_relation`
CHANGE COLUMN `reg_partner_id` `reg_merchant_id` INT(11) NULL DEFAULT NULL COMMENT '注册商户编号' ;
ALTER TABLE `zxdb`.`pm_info`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`ad_list`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_api`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NOT NULL ;
ALTER TABLE `zxdb`.`pt_kvset`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_kvset_member`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_mail_queue`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_mail_template`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_member_level`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_order`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL COMMENT '商户ID' ;
ALTER TABLE `zxdb`.`pt_saleconf`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NOT NULL ;
ALTER TABLE `zxdb`.`pt_shop`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_siteconf`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NOT NULL ;

-- ---------------

ALTER TABLE `zxdb`.`pt_ad`
RENAME TO  `zxdb`.`ad_list` ;

ALTER TABLE `zxdb`.`pt_ad_image`
RENAME TO  `zxdb`.`ad_image` ;

CREATE TABLE `ad_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(10) DEFAULT NULL,
  `opened` tinyint(1) DEFAULT NULL,
  `enabled` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `ad_position` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_id` int(11) DEFAULT NULL,
  `name` varchar(20) DEFAULT NULL,
  `description` varchar(100) DEFAULT NULL,
  `default_id` int(11) DEFAULT NULL,
  `opened` tinyint(1) DEFAULT NULL,
  `enabled` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_idx` (`group_id`),
  CONSTRAINT `id` FOREIGN KEY (`group_id`) REFERENCES `ad_group` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `ad_userset` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pos_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `ad_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `zxdb`.`ad_list`
DROP COLUMN `enabled`,
DROP COLUMN `is_internal`,
CHANGE COLUMN `merchant_id` `user_id` INT(11) NULL DEFAULT NULL ,
ADD COLUMN `show_times` INT NULL COMMENT '展现数量' AFTER `type_id`,
ADD COLUMN `click_times` INT NULL COMMENT '点击次数' AFTER `show_time`,
ADD COLUMN `show_days` INT NULL COMMENT '投放天数' AFTER `click_count`;

CREATE TABLE `zxdb`.`ad_hyperlink` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `ad_id` INT NULL,
  `title` VARCHAR(50) NULL,
  `link_url` VARCHAR(120) NULL,
  PRIMARY KEY (`id`));

--------------------------------


CREATE TABLE `zxdb`.`mm_level` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NULL,
  `require_exp` INT NULL,
  `program_signal` VARCHAR(45) NULL,
  `enabled` TINYINT(1) NULL,
  PRIMARY KEY (`id`));


ALTER TABLE `zxdb`.`pt_api`
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NOT NULL ,
  RENAME TO  `zxdb`.`mch_api` ;


  CREATE TABLE `zxdb`.`mch_enterpriseinfo` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `mch_id` INT NULL,
  `name` VARCHAR(45) NULL,
  `company_no` VARCHAR(45) NULL,
  `person_name` VARCHAR(10) NULL,
  `tel` VARCHAR(45) NULL,
  `address` VARCHAR(120) NULL,
  `province` INT NOT NULL,
  `city`  INT NOT NULL,
  `district` INT NOT NULL,
  `location` VARCHAR(45) NULL,
  `person_imageurl` VARCHAR(120) NULL,
  `company_imageurl` VARCHAR(120) NULL,
  `reviewed` TINYINT(1) NULL COMMENT '是否审核通过',
  `review_time` INT NULL,
  `remark` VARCHAR(45) NULL,
  `update_time` INT NULL,
  PRIMARY KEY (`id`));


ALTER TABLE `zxdb`.`pt_merchant`
DROP COLUMN `address`,
DROP COLUMN `phone`,
DROP COLUMN `tel`,
ADD COLUMN `province` INT NULL AFTER `logo`,
ADD COLUMN `city` INT NULL AFTER `province`,
ADD COLUMN `district` INT NULL AFTER `city`,
ADD COLUMN `enabled` TINYINT(1) NULL AFTER `join_time`,
ADD COLUMN `member_id` INT UNSIGNED NULL AFTER `id`,
RENAME TO  `zxdb`.`mch_merchant` ;


ALTER TABLE `zxdb`.`pt_saleconf`
  DROP COLUMN `present_convert_csn`,
  DROP COLUMN `flow_convert_csn`,
  DROP COLUMN `apply_csn`,
  DROP COLUMN `trans_csn`,
  DROP COLUMN `register_mode`,
  DROP COLUMN `ib_extra`,
  DROP COLUMN `ib_num`,
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NOT NULL ,
  ADD COLUMN `fx_sales` TINYINT(1) NULL COMMENT '是否启用分销' AFTER `mch_id`,
  DROP PRIMARY KEY,
  ADD PRIMARY KEY (`mch_id`), RENAME TO `zxdb`.`mch_saleconf` ;


CREATE TABLE `zxdb`.`mch_offline_shop` (
  `shop_id` INT NOT NULL,
  `tel` VARCHAR(45) NULL,
  `addr` VARCHAR(45) NULL,
  `lng` FLOAT(5,2) NULL,
  `lat` FLOAT(5,2) NULL,
  `deliver_radius` INT NULL COMMENT '配送范围',
  `province` INT NULL,
  `city` INT NULL,
  `district` INT NULL,
  PRIMARY KEY (`shop_id`));


ALTER TABLE `zxdb`.`pt_shop`
  DROP COLUMN `deliver_radius`,
  DROP COLUMN `location`,
  DROP COLUMN `phone`,
  DROP COLUMN `address`,
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NULL DEFAULT NULL ,
  ADD COLUMN `shop_type` TINYINT(1) NULL AFTER `mch_id`, RENAME TO  `zxdb`.`mch_shop` ;


CREATE TABLE `zxdb`.`mch_online_shop` (
  `shop_id` INT NOT NULL,
  `alias` VARCHAR(20) NULL,
  `tel` VARCHAR(45) NULL,
  `addr` VARCHAR(120) NULL,
  `host` VARCHAR(20) NULL,
  `logo` VARCHAR(120) NULL,
  `index_tit` VARCHAR(120) NULL,
  `sub_tit` VARCHAR(120) NULL,
  `notice_html` TEXT NULL,
  PRIMARY KEY (`shop_id`));



