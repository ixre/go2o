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
  `cat_id` INT NULL,
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

ALTER TABLE `flm`.`sale_order_item`
  DROP COLUMN `sku`;



CREATE TABLE `flm`.`express_template` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `user_id` INT NULL,
  `name` VARCHAR(45) NULL,
  `is_free` TINYINT(1) NULL,
  `basis` TINYINT(1) NULL,
  `first_unit` INT(5) NULL,
  `first_fee` DECIMAL(6,2) NULL,
  `add_unit` INT(5) NULL,
  `add_fee` DECIMAL(6,2) NULL,
  `enabled` TINYINT(1) NULL,
  PRIMARY KEY (`id`));


CREATE TABLE `flm`.`express_area_set` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `template_id` INT NULL,
  `code_list` VARCHAR(500) NULL,
  `name_list` VARCHAR(120) NULL,
  `first_unit` INT(5) NULL,
  `first_fee` DECIMAL(6,2) NULL,
  `add_unit` INT(5) NULL,
  `add_fee` DECIMAL(6,2) NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `flm`.`express_provider` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NULL,
  `letter` VARCHAR(1) NULL,
  `code` VARCHAR(10) NULL,
  `api_code` VARCHAR(10) NULL,
  `enabled` TINYINT(1) NULL,
  PRIMARY KEY (`id`));

ALTER TABLE `flm`.`gs_item`
  ADD COLUMN `weight` INT NULL COMMENT '重量,单位:克(g)' AFTER `cost`;

ALTER TABLE `flm`.`gs_snapshot`
  ADD COLUMN `weight` INT NULL COMMENT '单件重量,单位:克(g)' AFTER `img`;

ALTER TABLE `flm`.`sale_order`
  ADD COLUMN `express_fee` DECIMAL(8,2) NULL COMMENT '物流费' AFTER `discount_fee`;


ALTER TABLE `flm`.`mm_member`
  ADD COLUMN `check_code` VARCHAR(8) NULL AFTER `reg_time`,
  ADD COLUMN `check_expires` INT NULL AFTER `check_code`;


CREATE TABLE `flm`.`gs_sales_snapshot` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `snap_key` VARCHAR(45) NULL,
  `sku_id` INT NULL,
  `seller_id` INT NULL,
  `item_id` INT NULL,
  `cat_id` INT NULL,
  `goods_title` VARCHAR(120) NULL,
  `goods_no` VARCHAR(45) NULL,
  `sku` VARCHAR(120) NULL,
  `img` VARCHAR(120) NULL,
  `price` DECIMAL(8,2) NULL,
  `create_time` INT NULL,
  PRIMARY KEY (`id`));

ALTER TABLE `flm`.`sale_order`
  CHANGE COLUMN `total_fee` `goods_fee` DECIMAL(8,2) NULL DEFAULT NULL COMMENT '商品金额' ;


ALTER TABLE `flm`.`sale_sub_order`
  CHANGE COLUMN `total_fee` `goods_fee` DECIMAL(8,2) NULL DEFAULT NULL ,
  ADD COLUMN `express_fee` DECIMAL(4,2) NULL AFTER `discount_fee`;

ALTER TABLE `flm`.`sale_order`
  ADD COLUMN `package_fee` DECIMAL(4,2) NULL AFTER `express_fee`;

ALTER TABLE `flm`.`sale_sub_order`
  ADD COLUMN `package_fee` DECIMAL(4,2) NULL AFTER `express_fee`;


ALTER TABLE `flm`.`sale_sub_order`
  ADD COLUMN `buyer_id` INT NULL AFTER `parent_order`;

ALTER TABLE `zxdb`.`pt_order_log`
  ADD COLUMN `order_state` TINYINT(2) NULL AFTER `type`, RENAME TO  `zxdb`.`sale_order_log` ;

ALTER TABLE `flm`.`sale_sub_order`
  ADD COLUMN `is_paid` TINYINT(1) NULL AFTER `final_fee`;

ALTER TABLE `flm`.`sale_sub_order`
  CHANGE COLUMN `status` `state` TINYINT(1) NULL DEFAULT NULL ;

ALTER TABLE `flm`.`sale_order`
  CHANGE COLUMN `status` `state` TINYINT(1) NULL DEFAULT NULL ;

ALTER TABLE `flm`.`mch_enterprise_info`
  ADD COLUMN `person_id` VARCHAR(20) NULL COMMENT '法人身份证号' AFTER `person_name`;

ALTER TABLE `flm`.`mch_enterprise_info`
  ADD COLUMN `is_handled` TINYINT(1) NULL AFTER `company_imageurl`;

CREATE TABLE `flm`.`ship_order` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `order_id` INT NULL,
  `sp_id` INT NULL COMMENT '快递SP编号',
  `sp_order` VARCHAR(20) NULL COMMENT '快递SP单号',
  `exporess_log` VARCHAR(512) NULL,
  `amount` DECIMAL(8,2) NULL,
  `final_amount` DECIMAL(8,2) NULL,
  `ship_time` INT NULL COMMENT '发货时间',
  `state` TINYINT(1) NULL COMMENT '是否已收货',
  `update_time` INT NULL,
  PRIMARY KEY (`id`));


CREATE TABLE `flm`.`ship_item` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `ship_order` INT NULL,
  `snap_id` INT NULL,
  `quantity` INT NULL,
  `amount` DECIMAL(8,2) NULL,
  `final_amount` DECIMAL(8,2) NULL,
  PRIMARY KEY (`id`));
ALTER TABLE `flm`.`sale_order`
  CHANGE COLUMN `goods_fee` `goods_amount` DECIMAL(8,2) NULL DEFAULT NULL COMMENT '商品金额' ,
  CHANGE COLUMN `discount_fee` `discount_amount` DECIMAL(8,2) NULL DEFAULT NULL ,
  CHANGE COLUMN `final_fee` `final_amount` DECIMAL(8,2) NULL DEFAULT NULL ;

ALTER TABLE `flm`.`sale_order_item`
  CHANGE COLUMN `fee` `amount` DECIMAL(8,2) NULL DEFAULT NULL ,
  CHANGE COLUMN `final_fee` `final_amount` DECIMAL(8,2) NULL DEFAULT NULL ;

ALTER TABLE `flm`.`sale_sub_order`
  CHANGE COLUMN `goods_fee` `goods_amount` DECIMAL(8,2) NULL DEFAULT NULL ,
  CHANGE COLUMN `discount_fee` `discount_amount` DECIMAL(8,2) NULL DEFAULT NULL ,
  CHANGE COLUMN `final_fee` `final_amount` DECIMAL(8,2) NULL DEFAULT NULL ;

ALTER TABLE `flm`.`sale_order_item`
  ADD COLUMN `is_shipped` TINYINT(1) NULL AFTER `final_amount`;

ALTER TABLE `flm`.`mm_integral_log`
  CHANGE COLUMN `partner_id` `mch_id` INT(11) NULL DEFAULT NULL ;

ALTER TABLE `flm`.`pay_order`
  CHANGE COLUMN `sub_fee` `sub_amount` DECIMAL(8,2) NULL DEFAULT NULL COMMENT '立减金额' ,
  ADD COLUMN `adjustment_amount` DECIMAL(8,2) NULL COMMENT '调整金额' AFTER `sub_amount`;


CREATE TABLE `flm`.`sale_after_order` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `order_id` INT NULL,
  `vendor_id` INT NULL,
  `buyer_id` INT NULL,
  `type` TINYINT(1) NULL,
  `snap_id` INT NULL,
  `quantity` INT NULL,
  `reason` VARCHAR(255) NULL,
  `person_name` VARCHAR(10) NULL,
  `person_phone` VARCHAR(20) NULL,
  `rsp_name` VARCHAR(10) NULL COMMENT '退货快递名称',
  `rsp_order` VARCHAR(20) NULL COMMENT '退货快递单号',
  `rsp_image` VARCHAR(120) NULL,
  `remark` VARCHAR(45) NULL,
  `vendor_remark` VARCHAR(45) NULL,
  `state` TINYINT(1) NULL,
  `create_time` INT NULL,
  `update_time` INT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `flm`.`sale_return` (
  `id` INT NOT NULL,
  `amount` DECIMAL(8,2) NULL,
  `is_refund` TINYINT(1) NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `flm`.`sale_exchange` (
  `id` INT NOT NULL,
  `is_shipped` TINYINT(1) NULL,
  `sp_name` VARCHAR(20) NULL,
  `sp_order` VARCHAR(20) NULL,
  `ship_time` INT NULL,
  `is_received` TINYINT(1) NULL,
  `receive_time` INT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `flm`.`sale_refund` (
  `id` INT NOT NULL,
  `amount` DECIMAL(8,2) NULL,
  `is_refund` TINYINT(1) NULL,
  PRIMARY KEY (`id`));

ALTER TABLE `flm`.`sale_order_item`
  ADD COLUMN `return_quantity` INT NULL AFTER `quantity`;



ALTER TABLE `flm`.`sale_sub_order`
  ADD COLUMN `create_time` INT NULL AFTER `remark`;

ALTER TABLE `flm`.`mm_trusted_info`
  CHANGE COLUMN `body_number` `card_id` VARCHAR(20) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL ;


ALTER TABLE `flm`.`pay_order`
  ADD COLUMN `order_type` INT NULL COMMENT '支付单的类型，如购物或其他' AFTER `vendor_id`;

ALTER TABLE `flm`.`pay_order`
  ADD COLUMN `subject` VARCHAR(45) NULL COMMENT '支付单标题' AFTER `order_id`;


ALTER TABLE `flm`.`mm_integral_log`
  DROP COLUMN `mch_id`,
  CHANGE COLUMN `member_id` `member_id` INT(11) NOT NULL ,
  CHANGE COLUMN `type` `type` INT(11) NOT NULL ,
  CHANGE COLUMN `integral` `value` INT(11) NOT NULL ,
  CHANGE COLUMN `log` `remark` VARCHAR(100) NULL DEFAULT NULL ,
  CHANGE COLUMN `record_time` `create_time` INT(11) NOT NULL ;

ALTER TABLE `flm`.`mm_integral_log`
  ADD COLUMN `outer_no` VARCHAR(45) NULL AFTER `type`;


ALTER TABLE `flm`.`mm_account`
  CHANGE COLUMN `freezes_fee` `freezes_balance` FLOAT(10,2) NOT NULL AFTER `balance`,
  CHANGE COLUMN `freezes_present` `freezes_present` FLOAT(10,2) NOT NULL AFTER `present_balance`,
  CHANGE COLUMN `total_fee` `total_consumption` FLOAT(10,2) NOT NULL COMMENT '总消费' AFTER `total_pay`,
  CHANGE COLUMN `integral` `integral` INT(11) NOT NULL ,
  CHANGE COLUMN `balance` `balance` FLOAT(10,2) NOT NULL ,
  CHANGE COLUMN `present_balance` `present_balance` FLOAT(10,2) NOT NULL ,
  CHANGE COLUMN `total_present_fee` `total_present_fee` FLOAT(10,2) NOT NULL ,
  CHANGE COLUMN `flow_balance` `flow_balance` FLOAT(10,2) NOT NULL ,
  CHANGE COLUMN `grow_balance` `grow_balance` FLOAT(10,2) NOT NULL ,
  CHANGE COLUMN `grow_amount` `grow_amount` FLOAT(10,2) NOT NULL ,
  CHANGE COLUMN `grow_earnings` `grow_earnings` FLOAT(10,2) NOT NULL ,
  CHANGE COLUMN `grow_total_earnings` `grow_total_earnings` FLOAT(10,2) NOT NULL ,
  CHANGE COLUMN `total_charge` `total_charge` FLOAT(10,2) NOT NULL ,
  CHANGE COLUMN `total_pay` `total_pay` FLOAT(10,2) NOT NULL ,
  CHANGE COLUMN `update_time` `update_time` INT(11) NOT NULL COMMENT '积分' ,
  ADD COLUMN `freezes_integral` INT NOT NULL COMMENT '不可用积分' AFTER `integral`;


ALTER TABLE `flm`.`con_page`
  CHANGE COLUMN `mch_id` `user_id` INT(11) NULL DEFAULT NULL ;

DROP TABLE `flm`.`gs_member_price`;

CREATE TABLE `gs_member_price` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `goods_id` int(11) NOT NULL,
  `level` int(11) NOT NULL,
  `price` decimal(8,2) NOT NULL,
  `max_quota` int(11) NOT NULL,
  `enabled` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


ALTER TABLE `flm`.`express_provider`
  CHANGE COLUMN `letter` `group_flag` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL ;


ALTER TABLE `flm`.`gs_item`
  ADD COLUMN `express_tid` INT NULL COMMENT '快递模板编号' AFTER `sale_price`;

ALTER TABLE `flm`.`gs_snapshot`
  ADD COLUMN `express_tid` INT NULL AFTER `sale_price`;

ALTER TABLE `flm`.`gs_item`
  CHANGE COLUMN `img` `img` VARCHAR(120) NULL DEFAULT NULL ;

ALTER TABLE `flm`.`gs_item`
  CHANGE COLUMN `weight` `weight` FLOAT(6,2) NULL DEFAULT NULL COMMENT '重量,单位:克(g)' ;

ALTER TABLE `flm`.`gs_snapshot`
  CHANGE COLUMN `weight` `weight` FLOAT(6,2) NULL DEFAULT NULL COMMENT '单件重量,单位:克(g)' ;

ALTER TABLE `flm`.`gs_snapshot`
  ADD COLUMN `cost` DECIMAL(8,2) NULL AFTER `weight`;


ALTER TABLE `flm`.`gs_sales_snapshot`
  ADD COLUMN `cost` DECIMAL(8,2) NULL COMMENT '供货价' AFTER `img`;


CREATE TABLE mch_account (mch_id int(10) NOT NULL AUTO_INCREMENT comment '商户编号',
  balance decimal(10, 2) NOT NULL comment '余额',
  freeze_amount decimal(10, 2) NOT NULL comment '冻结金额',
  await_amount decimal(10, 2) NOT NULL comment '待入账金额',
  present_amount decimal(10, 2) NOT NULL comment '平台赠送金额',
  sales_amount decimal(10, 2) NOT NULL comment '累计销售总额',
  refund_amount decimal(10, 2) NOT NULL comment '累计退款金额',
  take_amount decimal(10, 2) NOT NULL comment '已提取金额',
  offline_sales decimal(10, 2) NOT NULL comment '线下销售金额',
  update_time int(11) NOT NULL comment '更新时间',
  PRIMARY KEY (mch_id)) comment='商户账户表';


CREATE TABLE mch_balance_log (id int(10) NOT NULL AUTO_INCREMENT,
                              mch_id int(10) NOT NULL comment '商户编号',
                              kind int(10) NOT NULL comment '日志类型',
                              title varchar(45) NOT NULL comment '标题',
                              outer_no varchar(45) NOT NULL comment '外部订单号',
                              amount float NOT NULL comment '金额',
                              csn_amount float DEFAULT 0.00 NOT NULL comment '手续费',
                              state tinyint(1) NOT NULL comment '状态',
                              create_time int(10) NOT NULL comment '创建时间',
                              update_time int(10) NOT NULL comment '更新时间',
  PRIMARY KEY (id)) comment='商户余额日志';

CREATE TABLE mch_day_chart (id int(11) NOT NULL AUTO_INCREMENT comment '编号',
                            mch_id int(11) NOT NULL comment '商户编号',
                            order_number int(11) NOT NULL comment '新增订单数量',
                            order_amount decimal(10, 2) NOT NULL comment '订单额',
                            buyer_number int(11) NOT NULL comment '购物会员数',
                            paid_number int(11) NOT NULL comment '支付单数量',
                            paid_amount decimal(10, 2) NOT NULL comment '支付总金额',
                            complete_orders int(11) NOT NULL comment '完成订单数',
                            in_amount decimal(10, 2) NOT NULL comment '入帐金额',
                            offline_orders int(11) NOT NULL comment '线下订单数量',
                            offline_amount decimal(10, 2) NOT NULL comment '线下订单金额',
                            `date` int(11) NOT NULL comment '日期',
                            date_str varchar(10) NOT NULL comment '日期字符串',
                            update_time int(11) NOT NULL comment '更新时间',
  PRIMARY KEY (id)) comment='商户每日报表';


ALTER TABLE `flm`.`mch_enterprise_info`
  CHANGE COLUMN `person_imageurl` `person_image` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL ,
  CHANGE COLUMN `company_imageurl` `company_image` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL ,
  ADD COLUMN `auth_doc` VARCHAR(120) NULL COMMENT '授权书' AFTER `company_image`;


CREATE TABLE mch_sign_up (id int(11) NOT NULL AUTO_INCREMENT,
  member_id int(11) NOT NULL,
  sign_no varchar(20) NOT NULL, usr varchar(45) NOT NULL,
  pwd varchar(45) NOT NULL, mch_name varchar(20) NOT NULL,
  province int(10) NOT NULL, city int(10) NOT NULL,
  district int(10) NOT NULL, shop_name varchar(20) NOT NULL,
  company_name varchar(20) NOT NULL, company_no varchar(20) NOT NULL,
  person_name varchar(10) NOT NULL, person_id varchar(20) NOT NULL,
  phone varchar(20) NOT NULL, address varchar(120) NOT NULL,
  person_image varchar(120) NOT NULL, company_image varchar(120) NOT NULL,
  auth_doc varchar(120) NOT NULL, reviewed int(1) NOT NULL,
  remark varchar(120) NOT NULL, submit_time int(11) NOT NULL,
  update_time int(11) NOT NULL, PRIMARY KEY (id));


ALTER TABLE `zxdb`.`msg_to`
  CHANGE COLUMN `to_id` `to_id` INT(11) NOT NULL ,
  CHANGE COLUMN `to_role` `to_role` TINYINT(2) NOT NULL ,
  CHANGE COLUMN `content_id` `content_id` INT(11) NOT NULL ,
  CHANGE COLUMN `has_read` `has_read` TINYINT(1) NOT NULL ,
  CHANGE COLUMN `read_time` `read_time` INT(11) NOT NULL ,
  ADD COLUMN `msg_id` INT(11) NOT NULL AFTER `to_role`;



ALTER TABLE `zxdb`.`mm_account`
  ADD COLUMN `priority_pay` TINYINT(1) NULL COMMENT '优先（默认）支付账户' AFTER `total_consumption`;

CREATE TABLE `mm_balance_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) NOT NULL,
  `kind` int(11) NOT NULL COMMENT '业务类型',
  `title` varchar(45) NOT NULL COMMENT '标题',
  `outer_no` varchar(45) NOT NULL COMMENT '外部订单号',
  `amount` float(8,2) NOT NULL COMMENT '金额',
  `csn_fee` float(8,2) NOT NULL COMMENT '手续费',
  `state` tinyint(1) NOT NULL COMMENT '状态，比如提现需要确认等',
  `rel_user` int(11) NOT NULL COMMENT '关联操作人员编号',
  `create_time` int(11) NOT NULL,
  `update_time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


CREATE TABLE `mm_present_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) NOT NULL,
  `kind` int(11) NOT NULL COMMENT '业务类型',
  `title` varchar(45) NOT NULL COMMENT '标题',
  `outer_no` varchar(45) NOT NULL COMMENT '外部订单号',
  `amount` float(8,2) NOT NULL COMMENT '金额',
  `csn_fee` float(8,2) NOT NULL COMMENT '手续费',
  `state` tinyint(1) NOT NULL COMMENT '状态，比如提现需要确认等',
  `rel_user` int(11) NOT NULL COMMENT '关联操作人员编号',
  `create_time` int(11) NOT NULL,
  `update_time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
















