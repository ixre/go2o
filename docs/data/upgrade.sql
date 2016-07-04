ALTER TABLE `zxdb`.`pt_merchant`
RENAME TO  `zxdb`.`pt_merchant` ;
ALTER TABLE `zxdb`.`pt_page`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`dlv_partner_bind`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL , RENAME TO  `zxdb`.`dlv_merchant_bind` ;
ALTER TABLE `zxdb`.`gs_category`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL COMMENT '商户ID(pattern ID);如果为空，则表示模式分类' ;
ALTER TABLE `zxdb`.`gs_sale_label`
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

ALTER TABLE `zxdb`.`mch_merchant`
  ADD COLUMN `level` INT NULL COMMENT '商户等级' AFTER `name`;


ALTER TABLE `flm`.`gs_category`
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NULL
DEFAULT NULL COMMENT '商户ID(merhantId ID);如果为空，则表示系统的f分类 ';

ALTER TABLE `flm`.`gs_category`
  ADD COLUMN `level` TINYINT(1) NULL AFTER `sort_number`;

ALTER TABLE `flm`.`mch_merchant`
  ADD COLUMN `self_sales` TINYINT(1) NULL AFTER `name`;


ALTER TABLE `flm`.`gs_sale_label`
  DROP COLUMN `is_internal`,
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NULL DEFAULT NULL ,
  CHANGE COLUMN `goods_image` `label_image` VARCHAR(100) NULL DEFAULT NULL , RENAME TO  `flm`.`gs_sale_label` ;


ALTER TABLE `flm`.`pt_page`
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NULL DEFAULT NULL , RENAME TO  `flm`.`mch_page` ;

ALTER TABLE `flm`.`gs_item`
  ADD COLUMN `supplier_id` INT NULL AFTER `category_id`;

ALTER TABLE `flm`.`pm_info`
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NULL DEFAULT NULL ;

ALTER TABLE `flm`.`ad_position`
  CHANGE COLUMN `description` `key` VARCHAR(45) NULL DEFAULT NULL AFTER `group_id`,
  CHANGE COLUMN `name` `name` VARCHAR(45) NULL DEFAULT NULL ;

CREATE TABLE `flm`.`msg_list` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `msg_type` TINYINT(1) NULL,
  `use_for` TINYINT(1) NULL,
  `sender_id` INT NULL,
  `sender_role` TINYINT(2) NULL,
  `to_role` TINYINT(2) NULL,
  `all_user` TINYINT(1) NULL,
  `read_only` TINYINT(1) NULL,
  `create_time` INT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `flm`.`msg_content` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `msg_id` INT NULL,
  `msg_data` TEXT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `flm`.`msg_to` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `to_id` INT NULL,
  `to_role` TINYINT(2) NULL,
  `content_id` INT NULL,
  `has_read` TINYINT(1) NULL,
  `read_time` INT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `flm`.`msg_replay` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `refer_id` INT NULL,
  `sender_id` INT NULL,
  `sender_role` TINYINT(2) NULL,
  `content` TEXT NULL,
  PRIMARY KEY (`id`));


ALTER TABLE `flm`.`sale_cart_item`
  ADD COLUMN `mch_id` INT NULL AFTER `cart_id`,
  ADD COLUMN `shop_id` INT NULL AFTER `mch_id`;

ALTER TABLE `flm`.`mm_level`
  ADD COLUMN `is_official` TINYINT(1) NULL AFTER `program_signal`;

CREATE TABLE `flm`.`mm_trusted_info` (
  `member_id` INT NOT NULL,
  `real_name` VARCHAR(10) NULL,
  `body_number` VARCHAR(20) NULL,
  `trust_image` VARCHAR(120) NULL,
  `is_handle` TINYINT(1) NULL,
  `reviewed` TINYINT(1) NULL,
  `review_time` INT NULL,
  `remark` VARCHAR(120) NULL,
  `update_time` INT NULL,
  PRIMARY KEY (`member_id`));


CREATE TABLE `mm_profile` (
  `member_id` int(11) NOT NULL,
  `name` varchar(20) DEFAULT NULL COMMENT '名字',
  `sex` int(1) DEFAULT NULL COMMENT '性别(0: 未知,1:男,2：女)',
  `avatar` varchar(80) DEFAULT NULL,
  `birthday` varchar(20) DEFAULT NULL,
  `phone` varchar(15) DEFAULT NULL,
  `address` varchar(100) DEFAULT NULL COMMENT '送餐地址',
  `qq` varchar(15) DEFAULT NULL,
  `im` varchar(45) DEFAULT NULL,
  `ext_1` varchar(45) DEFAULT NULL,
  `ext_2` varchar(45) DEFAULT NULL,
  `ext_3` varchar(45) DEFAULT NULL,
  `ext_4` varchar(45) DEFAULT NULL,
  `ext_5` varchar(45) DEFAULT NULL,
  `ext_6` varchar(45) DEFAULT NULL,
  `email` varchar(50) DEFAULT NULL,
  `remark` varchar(100) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`member_id`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8;

# copy profile info to mm_profile

INSERT INTO mm_profile SELECT `id`,`name`,`sex`,`avatar`,`birthday`,`phone`,
`address`,`qq`,`im`,`ext_1`, `ext_2`,`ext_3`,`ext_4`,`ext_5`,`ext_6`,`email`,
`remark`,`update_time` FROM mm_member;

ALTER TABLE `flm`.`mm_profile`
  ADD COLUMN `province` INT NULL AFTER `email`,
  ADD COLUMN `city` INT NULL AFTER `province`,
  ADD COLUMN `district` INT NULL AFTER `city`;

ALTER TABLE `flm`.`mm_member`
  DROP COLUMN `remark`,
  DROP COLUMN `email`,
  DROP COLUMN `ext_6`,
  DROP COLUMN `ext_5`,
  DROP COLUMN `ext_4`,
  DROP COLUMN `ext_3`,
  DROP COLUMN `ext_2`,
  DROP COLUMN `ext_1`,
  DROP COLUMN `im`,
  DROP COLUMN `qq`,
  DROP COLUMN `address`,
  DROP COLUMN `phone`,
  DROP COLUMN `birthday`,
  DROP COLUMN `avatar`,
  DROP COLUMN `sex`,
  DROP COLUMN `name`;


CREATE TABLE `flm`.`mm_favorite` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '会员收藏表',
  `member_id` INT NULL,
  `fav_type` TINYINT(1) NULL,
  `refer_id` INT NULL,
  `update_time` INT NULL,
  PRIMARY KEY (`id`));


ALTER TABLE `flm`.`mm_deliver_addr`
  CHANGE COLUMN `address` `address` VARCHAR(80) NULL DEFAULT NULL COMMENT '详细地址' ,
  ADD COLUMN `province` INT NULL AFTER `phone`,
  ADD COLUMN `city` INT NULL AFTER `province`,
  ADD COLUMN `district` INT NULL AFTER `city`,
  ADD COLUMN `area` VARCHAR(50) NULL COMMENT '省市区' AFTER `district`;

ALTER TABLE `flm`.`mch_page`
  RENAME TO  `flm`.`con_page` ;

CREATE TABLE `flm`.`con_article_category` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `parent_id` INT NULL,
  `name` VARCHAR(45) NULL,
  `alias` VARCHAR(45) NULL,
  `title` VARCHAR(120) NULL,
  `keywords` VARCHAR(120) NULL,
  `describe` VARCHAR(250) NULL,
  `sort_number` INT NULL,
  `location` VARCHAR(120) NULL,
  `update_time` INT NULL,
  PRIMARY KEY (`id`));


CREATE TABLE `flm`.`con_article` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `category_id` INT NULL,
  `title` VARCHAR(120) NULL,
  `small_title` VARCHAR(45) NULL,
  `thumbnail` VARCHAR(120) NULL,
  `location` VARCHAR(120) NULL,
  `publisher_id` INT NULL,
  `content` TEXT NULL,
  `tags` VARCHAR(120) NULL,
  `view_count` INT NULL,
  `sort_number` INT NULL,
  `create_time` INT NULL,
  `update_time` INT NULL,
  PRIMARY KEY (`id`));


ALTER TABLE `zxdb`.`gs_snapshot`
  RENAME TO  `zxdb`.`gs_sale_snapshot` ;


CREATE TABLE `zxdb`.`gs_snapshot` (
  `sku_id` INT NOT NULL,
  `vendor_id` INT NULL,
  `snapshot_key` VARCHAR(45) NULL,
  `goods_title` VARCHAR(80) NULL,
  `small_title` VARCHAR(45) NULL,
  `goods_no` VARCHAR(45) NULL,
  `item_id` INT NULL,
  `category_id` INT NULL,
  `img` VARCHAR(120) NULL,
  `price` DECIMAL(8,2) NULL,
  `sale_price` DECIMAL(8,2) NULL,
  `update_time` INT NULL,
  PRIMARY KEY (`id`));

ALTER TABLE `flm`.`gs_item`
  ADD COLUMN `has_review` TINYINT(1) NULL AFTER `state`,
  ADD COLUMN `review_pass` TINYINT(1) NULL AFTER `has_review`;

ALTER TABLE `flm`.`gs_snapshot`
  CHANGE COLUMN `category_id` `cat_id` INT(11) NULL DEFAULT NULL ,
  ADD COLUMN `on_shelves` TINYINT(1) NULL DEFAULT 1 AFTER `cat_id`;


ALTER TABLE `flm`.`gs_snapshot`
  ADD COLUMN `level_sales` TINYINT(1) NULL COMMENT '是否有会员价' AFTER `sale_price`;

ALTER TABLE `flm`.`sale_cart_item`
  CHANGE COLUMN `mch_id` `vendor_id` INT(11) NULL DEFAULT NULL ,
  CHANGE COLUMN `quantity` `quantity` INT(8) NULL DEFAULT NULL ,
  ADD COLUMN `checked` TINYINT(1) NULL AFTER `quantity`;

ALTER TABLE `flm`.`gs_snapshot`
  ADD COLUMN `sale_num` INT NULL AFTER `level_sales`,
  ADD COLUMN `stock_num` INT NULL AFTER `sale_num`;


ALTER TABLE `flm`.`pt_order`
  CHANGE COLUMN `member_id` `buyner_id` INT(11) NULL DEFAULT NULL COMMENT '-1代表游客订餐' ,
  CHANGE COLUMN `merchant_id` `vendor_id` INT(11) NULL DEFAULT NULL COMMENT '商家ID' ;

CREATE TABLE `flm`.`pay_order` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `trade_no` VARCHAR(45) NULL,
  `vendor_id` INT NULL,
  `order_id` INT NULL,
  `buy_user` INT NULL,
  `payment_user` INT NULL,
  `total_fee` DECIMAL(8,2) NULL,
  `balance_discount` DECIMAL(8,2) NULL,
  `integral_discount` DECIMAL(8,2) NULL,
  `system_discount` DECIMAL(8,2) NULL,
  `coupon_discount` DECIMAL(8,2) NULL,
  `sub_fee` DECIMAL(8,2) NULL,
  `final_fee` DECIMAL(8,2) NULL,
  `payment_opt` TINYINT(2) NULL,
  `payment_sign` TINYINT(1) NULL,
  `outer_no` VARCHAR(45) NULL COMMENT '外部订单号',
  `create_time` INT NULL,
  `paid_time` INT NULL,
  `state` VARCHAR(45) NULL,
  PRIMARY KEY (`id`));


ALTER TABLE `zxdb`.`pt_order`
  CHANGE COLUMN `buyner_id` `buyer_id` INT(11) NULL DEFAULT NULL COMMENT '-1代表游客订餐' ;

ALTER TABLE `zxdb`.`pt_order_item`
  CHANGE COLUMN `snapshot_id` `snap_id` INT(11) NULL DEFAULT NULL ,
  CHANGE COLUMN `quantity` `quantity` INT(6) NULL DEFAULT NULL ,
  CHANGE COLUMN `update_time` `update_time` INT NULL DEFAULT NULL ,
  ADD COLUMN `vendor_id` INT NULL AFTER `order_id`,
  ADD COLUMN `shop_id` INT NULL AFTER `vendor_id`,
  ADD COLUMN `sku_id` INT NULL AFTER `shop_id`,
  ADD COLUMN `final_fee` DECIMAL(8,2) NULL AFTER `fee`;

CREATE TABLE `flm`.`sale_order` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `order_no` VARCHAR(20) NULL,
  `buyer_id` INT NULL,
  `items_info` VARCHAR(255) NULL,
  `total_fee` DECIMAL(8,2) NULL,
  `discount_fee` DECIMAL(8,2) NULL,
  `final_fee` DECIMAL(8,2) NULL,
  `is_paid` TINYINT(1) NULL,
  `paid_time` INT NULL,
  `consignee_person` VARCHAR(45) NULL,
  `consignee_phone` VARCHAR(45) NULL,
  `shipping_address` VARCHAR(120) NULL,
  `shipping_time` VARCHAR(45) NULL,
  `create_time` INT NULL,
  `update_time` INT NULL,
  `status` TINYINT(1) NULL,
  PRIMARY KEY (`id`));

ALTER TABLE `flm`.`pt_order_pb`
  CHANGE COLUMN `order_no` `order_id` INT NULL DEFAULT NULL AFTER `id`;


CREATE TABLE `flm`.`sale_sub_order` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `order_no` VARCHAR(20) NULL,
  `parent_order` INT NULL,
  `vendor_id` INT NULL,
  `shop_id` INT NULL,
  `subject` VARCHAR(45) NULL,
  `items_info` VARCHAR(255) NULL,
  `total_fee` DECIMAL(8,2) NULL,
  `discount_fee` DECIMAL(8,2) NULL,
  `final_fee` DECIMAL(8,2) NULL,
  `is_suspend` TINYINT(1) NULL,
  `note` VARCHAR(120) NULL,
  `remark` VARCHAR(120) NULL,
  `update_time` INT NULL,
  `status` TINYINT(1) NULL,
  PRIMARY KEY (`id`));


ALTER TABLE `flm`.`pt_order_item`
  CHANGE COLUMN `update_time` `update_time` INT NULL DEFAULT NULL ,
  ADD COLUMN `sku_id` INT NULL AFTER `order_id`,
  ADD COLUMN `final_fee` DECIMAL(8,2) NULL AFTER `fee`, RENAME TO  `flm`.`sale_order_item` ;

ALTER TABLE `flm`.`sale_order_item`
  CHANGE COLUMN `snapshot_id` `snap_id` INT(11) NULL DEFAULT NULL ;

























