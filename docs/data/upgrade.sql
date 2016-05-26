ALTER TABLE `zxdb`.`pt_merchant`
RENAME TO  `zxdb`.`pt_merchant` ;
ALTER TABLE `zxdb`.`pt_page`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`dlv_partner_bind`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL , RENAME TO  `zxdb`.`dlv_merchant_bind` ;
ALTER TABLE `zxdb`.`gs_category`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL COMMENT '商家ID(pattern ID);如果为空，则表示模式分类' ;
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
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL COMMENT '商家ID' ;
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
  `default` int(11) DEFAULT NULL,
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



