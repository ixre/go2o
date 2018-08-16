ALTER TABLE `pt_merchant`
RENAME TO `pt_merchant`;
ALTER TABLE `pt_page`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL;
ALTER TABLE `dlv_partner_bind`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL,
RENAME TO `dlv_merchant_bind`;
ALTER TABLE `gs_category`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL
COMMENT '商户ID(pattern ID);如果为空，则表示模式分类';
ALTER TABLE `gs_sale_label`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL;
ALTER TABLE `mm_relation`
  CHANGE COLUMN `reg_partner_id` `reg_merchant_id` INT(11) NULL DEFAULT NULL
COMMENT '注册商户编号';
ALTER TABLE `mm_relation`
  ADD COLUMN `refer_str` VARCHAR(250) NULL
  AFTER `invi_member_id`;

ALTER TABLE `pm_info`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL;
ALTER TABLE `ad_list`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL;
ALTER TABLE `pt_api`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NOT NULL;
ALTER TABLE `pt_kvset`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL;
ALTER TABLE `pt_kvset_member`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL;
ALTER TABLE `pt_mail_queue`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL;
ALTER TABLE `pt_mail_template`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL;
ALTER TABLE `pt_member_level`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL;
ALTER TABLE `pt_order`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL
COMMENT '商户ID';
ALTER TABLE `pt_saleconf`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NOT NULL;
ALTER TABLE `pt_shop`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL;
ALTER TABLE `pt_siteconf`
  CHANGE COLUMN `partner_id` `merchant_id` INT(11) NOT NULL;
-- ---------------


ALTER TABLE `pt_ad`
RENAME TO `ad_list`;

ALTER TABLE `pt_ad_image`
RENAME TO `ad_image`;

CREATE TABLE `ad_group` (
  `id`      INT(11) NOT NULL AUTO_INCREMENT,
  `name`    VARCHAR(10)      DEFAULT NULL,
  `opened`  TINYINT(1)       DEFAULT NULL,
  `enabled` TINYINT(1)       DEFAULT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `ad_position` (
  `id`          INT(11) NOT NULL AUTO_INCREMENT,
  `group_id`    INT(11)          DEFAULT NULL,
  `name`        VARCHAR(20)      DEFAULT NULL,
  `description` VARCHAR(100)     DEFAULT NULL,
  `default_id`  INT(11)          DEFAULT NULL,
  `opened`      TINYINT(1)       DEFAULT NULL,
  `enabled`     TINYINT(1)       DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_idx` (`group_id`),
  CONSTRAINT `id` FOREIGN KEY (`group_id`) REFERENCES `ad_group` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `ad_userset` (
  `id`      INT(11) NOT NULL AUTO_INCREMENT,
  `pos_id`  INT(11)          DEFAULT NULL,
  `user_id` INT(11)          DEFAULT NULL,
  `ad_id`   INT(11)          DEFAULT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

ALTER TABLE `ad_list`
  DROP COLUMN `enabled`,
  DROP COLUMN `is_internal`,
  CHANGE COLUMN `merchant_id` `user_id` INT(11) NULL DEFAULT NULL,
  ADD COLUMN `show_times` INT NULL
COMMENT '展现数量'
  AFTER `type_id`,
  ADD COLUMN `click_times` INT NULL
COMMENT '点击次数'
  AFTER `show_time`,
  ADD COLUMN `show_days` INT NULL
COMMENT '投放天数'
  AFTER `click_count`;

CREATE TABLE `ad_hyperlink` (
  `id`       INT          NOT NULL AUTO_INCREMENT,
  `ad_id`    INT          NULL,
  `title`    VARCHAR(50)  NULL,
  `link_url` VARCHAR(120) NULL,
  PRIMARY KEY (`id`)
);


CREATE TABLE `mm_level` (
  `id`             INT         NOT NULL AUTO_INCREMENT,
  `name`           VARCHAR(45) NULL,
  `require_exp`    INT         NULL,
  `program_signal` VARCHAR(45) NULL,
  `enabled`        TINYINT(1)  NULL,
  PRIMARY KEY (`id`)
);


ALTER TABLE `pt_api`
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NOT NULL,
RENAME TO `mch_api`;


CREATE TABLE `mch_enterpriseinfo` (
  `id`               INT          NOT NULL AUTO_INCREMENT,
  `mch_id`           INT          NULL,
  `name`             VARCHAR(45)  NULL,
  `company_no`       VARCHAR(45)  NULL,
  `person_name`      VARCHAR(10)  NULL,
  `tel`              VARCHAR(45)  NULL,
  `address`          VARCHAR(120) NULL,
  `province`         INT          NOT NULL,
  `city`             INT          NOT NULL,
  `district`         INT          NOT NULL,
  `location`         VARCHAR(45)  NULL,
  `person_imageurl`  VARCHAR(120) NULL,
  `company_imageurl` VARCHAR(120) NULL,
  `reviewed`         TINYINT(1)   NULL
  COMMENT '是否审核通过',
  `review_time`      INT          NULL,
  `remark`           VARCHAR(45)  NULL,
  `update_time`      INT          NULL,
  PRIMARY KEY (`id`)
);


ALTER TABLE `pt_merchant`
  DROP COLUMN `address`,
  DROP COLUMN `phone`,
  DROP COLUMN `tel`,
  ADD COLUMN `province` INT NULL
  AFTER `logo`,
  ADD COLUMN `city` INT NULL
  AFTER `province`,
  ADD COLUMN `district` INT NULL
  AFTER `city`,
  ADD COLUMN `enabled` TINYINT(1) NULL
  AFTER `join_time`,
  ADD COLUMN `member_id` INT UNSIGNED NULL
  AFTER `id`,
RENAME TO `mch_merchant`;


ALTER TABLE `pt_saleconf`
  DROP COLUMN `present_convert_csn`,
  DROP COLUMN `flow_convert_csn`,
  DROP COLUMN `apply_csn`,
  DROP COLUMN `trans_csn`,
  DROP COLUMN `register_mode`,
  DROP COLUMN `ib_extra`,
  DROP COLUMN `ib_num`,
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NOT NULL,
  ADD COLUMN `fx_sales` TINYINT(1) NULL
COMMENT '是否启用分销'
  AFTER `mch_id`,
  DROP PRIMARY KEY,
  ADD PRIMARY KEY (`mch_id`),
RENAME TO `mch_saleconf`;


CREATE TABLE `mch_offline_shop` (
  `shop_id`        INT         NOT NULL,
  `tel`            VARCHAR(45) NULL,
  `addr`           VARCHAR(45) NULL,
  `lng`            FLOAT(5, 2) NULL,
  `lat`            FLOAT(5, 2) NULL,
  `deliver_radius` INT         NULL
  COMMENT '配送范围',
  `province`       INT         NULL,
  `city`           INT         NULL,
  `district`       INT         NULL,
  PRIMARY KEY (`shop_id`)
);


ALTER TABLE `pt_shop`
  DROP COLUMN `deliver_radius`,
  DROP COLUMN `location`,
  DROP COLUMN `phone`,
  DROP COLUMN `address`,
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NULL DEFAULT NULL,
  ADD COLUMN `shop_type` TINYINT(1) NULL
  AFTER `mch_id`,
RENAME TO `mch_shop`;


CREATE TABLE `mch_online_shop` (
  `shop_id`     INT          NOT NULL,
  `alias`       VARCHAR(20)  NULL,
  `tel`         VARCHAR(45)  NULL,
  `addr`        VARCHAR(120) NULL,
  `host`        VARCHAR(20)  NULL,
  `logo`        VARCHAR(120) NULL,
  `index_tit`   VARCHAR(120) NULL,
  `sub_tit`     VARCHAR(120) NULL,
  `notice_html` TEXT         NULL,
  PRIMARY KEY (`shop_id`)
);

ALTER TABLE `mch_merchant`
  ADD COLUMN `level` INT NULL
COMMENT '商户等级'
  AFTER `name`;


ALTER TABLE `gs_category`
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NULL
  DEFAULT NULL
COMMENT '商户ID(merhantId ID);如果为空，则表示系统的f分类 ';

ALTER TABLE `gs_category`
  ADD COLUMN `level` TINYINT(1) NULL
  AFTER `sort_number`;

ALTER TABLE `mch_merchant`
  ADD COLUMN `self_sales` TINYINT(1) NULL
  AFTER `name`;


ALTER TABLE `gs_sale_label`
  DROP COLUMN `is_internal`,
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NULL DEFAULT NULL,
  CHANGE COLUMN `goods_image` `label_image` VARCHAR(100) NULL DEFAULT NULL,
RENAME TO `gs_sale_label`;


ALTER TABLE `pt_page`
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NULL DEFAULT NULL,
RENAME TO `mch_page`;

ALTER TABLE `gs_item`
  ADD COLUMN `supplier_id` INT NULL
  AFTER `category_id`;

ALTER TABLE `pm_info`
  CHANGE COLUMN `merchant_id` `mch_id` INT(11) NULL DEFAULT NULL;

ALTER TABLE `ad_position`
  CHANGE COLUMN `description` `key` VARCHAR(45) NULL DEFAULT NULL
  AFTER `group_id`,
  CHANGE COLUMN `name` `name` VARCHAR(45) NULL DEFAULT NULL;

CREATE TABLE `msg_list` (
  `id`          INT        NOT NULL AUTO_INCREMENT,
  `msg_type`    TINYINT(1) NULL,
  `use_for`     TINYINT(1) NULL,
  `sender_id`   INT        NULL,
  `sender_role` TINYINT(2) NULL,
  `to_role`     TINYINT(2) NULL,
  `all_user`    TINYINT(1) NULL,
  `read_only`   TINYINT(1) NULL,
  `create_time` INT        NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `msg_content` (
  `id`       INT  NOT NULL AUTO_INCREMENT,
  `msg_id`   INT  NULL,
  `msg_data` TEXT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `msg_to` (
  `id`         INT        NOT NULL AUTO_INCREMENT,
  `to_id`      INT        NULL,
  `to_role`    TINYINT(2) NULL,
  `content_id` INT        NULL,
  `has_read`   TINYINT(1) NULL,
  `read_time`  INT        NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `msg_replay` (
  `id`          INT        NOT NULL AUTO_INCREMENT,
  `refer_id`    INT        NULL,
  `sender_id`   INT        NULL,
  `sender_role` TINYINT(2) NULL,
  `content`     TEXT       NULL,
  PRIMARY KEY (`id`)
);


ALTER TABLE `sale_cart_item`
  ADD COLUMN `mch_id` INT NULL
  AFTER `cart_id`,
  ADD COLUMN `shop_id` INT NULL
  AFTER `mch_id`;

ALTER TABLE `mm_level`
  ADD COLUMN `is_official` TINYINT(1) NULL
  AFTER `program_signal`;

CREATE TABLE `mm_trusted_info` (
  `member_id`   INT          NOT NULL,
  `real_name`   VARCHAR(10)  NULL,
  `body_number` VARCHAR(20)  NULL,
  `trust_image` VARCHAR(120) NULL,
  `reviewed`    TINYINT(1)   NULL,
  `review_time` INT          NULL,
  `remark`      VARCHAR(120) NULL,
  `update_time` INT          NULL,
  PRIMARY KEY (`member_id`)
);


CREATE TABLE `mm_profile` (
  `member_id`   INT(11) NOT NULL,
  `name`        VARCHAR(20)  DEFAULT NULL
  COMMENT '名字',
  `sex`         INT(1)       DEFAULT NULL
  COMMENT '性别(0: 未知,1:男,2：女)',
  `avatar`      VARCHAR(80)  DEFAULT NULL,
  `birthday`    VARCHAR(20)  DEFAULT NULL,
  `phone`       VARCHAR(15)  DEFAULT NULL,
  `address`     VARCHAR(100) DEFAULT NULL
  COMMENT '送餐地址',
  `qq`          VARCHAR(15)  DEFAULT NULL,
  `im`          VARCHAR(45)  DEFAULT NULL,
  `ext_1`       VARCHAR(45)  DEFAULT NULL,
  `ext_2`       VARCHAR(45)  DEFAULT NULL,
  `ext_3`       VARCHAR(45)  DEFAULT NULL,
  `ext_4`       VARCHAR(45)  DEFAULT NULL,
  `ext_5`       VARCHAR(45)  DEFAULT NULL,
  `ext_6`       VARCHAR(45)  DEFAULT NULL,
  `email`       VARCHAR(50)  DEFAULT NULL,
  `remark`      VARCHAR(100) DEFAULT NULL,
  `update_time` INT(11)      DEFAULT NULL,
  PRIMARY KEY (`member_id`)
)
  ENGINE = MyISAM
  DEFAULT CHARSET = utf8;

# copy profile info to mm_profile

INSERT INTO mm_profile
SELECT `id`,
       `name`,
       `sex`,
       `avatar`,
       `birthday`,
       `phone`,
       `address`,
       `qq`,
       `im`,
       `ext_1`,
       `ext_2`,
       `ext_3`,
       `ext_4`,
       `ext_5`,
       `ext_6`,
       `email`,
       `remark`,
       `update_time`
FROM mm_member;

ALTER TABLE `mm_profile`
  ADD COLUMN `province` INT NULL
  AFTER `email`,
  ADD COLUMN `city` INT NULL
  AFTER `province`,
  ADD COLUMN `district` INT NULL
  AFTER `city`;

ALTER TABLE `mm_member`
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


CREATE TABLE `mm_favorite` (
  `id`          INT        NOT NULL AUTO_INCREMENT
  COMMENT '会员收藏表',
  `member_id`   INT        NULL,
  `fav_type`    TINYINT(1) NULL,
  `refer_id`    INT        NULL,
  `update_time` INT        NULL,
  PRIMARY KEY (`id`)
);


ALTER TABLE `mm_deliver_addr`
  CHANGE COLUMN `address` `address` VARCHAR(80) NULL DEFAULT NULL
COMMENT '详细地址',
  ADD COLUMN `province` INT NULL
  AFTER `phone`,
  ADD COLUMN `city` INT NULL
  AFTER `province`,
  ADD COLUMN `district` INT NULL
  AFTER `city`,
  ADD COLUMN `area` VARCHAR(50) NULL
COMMENT '省市区'
  AFTER `district`;

ALTER TABLE `mch_page`
RENAME TO `con_page`;

CREATE TABLE `con_article_category` (
  `id`          INT          NOT NULL AUTO_INCREMENT,
  `parent_id`   INT          NULL,
  `name`        VARCHAR(45)  NULL,
  `alias`       VARCHAR(45)  NULL,
  `title`       VARCHAR(120) NULL,
  `keywords`    VARCHAR(120) NULL,
  `describe`    VARCHAR(250) NULL,
  `sort_number` INT          NULL,
  `location`    VARCHAR(120) NULL,
  `update_time` INT          NULL,
  PRIMARY KEY (`id`)
);


CREATE TABLE `con_article` (
  `id`           INT          NOT NULL AUTO_INCREMENT,
  `cat_id`       INT          NULL,
  `title`        VARCHAR(120) NULL,
  `small_title`  VARCHAR(45)  NULL,
  `thumbnail`    VARCHAR(120) NULL,
  `location`     VARCHAR(120) NULL,
  `publisher_id` INT          NULL,
  `content`      TEXT         NULL,
  `tags`         VARCHAR(120) NULL,
  `view_count`   INT          NULL,
  `sort_number`  INT          NULL,
  `create_time`  INT          NULL,
  `update_time`  INT          NULL,
  PRIMARY KEY (`id`)
);


ALTER TABLE `gs_snapshot`
RENAME TO `gs_sale_snapshot`;


CREATE TABLE `gs_snapshot` (
  `sku_id`       INT           NOT NULL,
  `vendor_id`    INT           NULL,
  `snapshot_key` VARCHAR(45)   NULL,
  `goods_title`  VARCHAR(80)   NULL,
  `small_title`  VARCHAR(45)   NULL,
  `goods_no`     VARCHAR(45)   NULL,
  `item_id`      INT           NULL,
  `category_id`  INT           NULL,
  `img`          VARCHAR(120)  NULL,
  `price`        DECIMAL(8, 2) NULL,
  `sale_price`   DECIMAL(8, 2) NULL,
  `update_time`  INT           NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `gs_item`
  ADD COLUMN `has_review` TINYINT(1) NULL
  AFTER `state`,
  ADD COLUMN `review_pass` TINYINT(1) NULL
  AFTER `has_review`;

ALTER TABLE `gs_snapshot`
  CHANGE COLUMN `category_id` `cat_id` INT(11) NULL DEFAULT NULL,
  ADD COLUMN `on_shelves` TINYINT(1) NULL DEFAULT 1
  AFTER `cat_id`;


ALTER TABLE `gs_snapshot`
  ADD COLUMN `level_sales` TINYINT(1) NULL
COMMENT '是否有会员价'
  AFTER `sale_price`;

ALTER TABLE `sale_cart_item`
  CHANGE COLUMN `mch_id` `vendor_id` INT(11) NULL DEFAULT NULL,
  CHANGE COLUMN `quantity` `quantity` INT(8) NULL DEFAULT NULL,
  ADD COLUMN `checked` TINYINT(1) NULL
  AFTER `quantity`;

ALTER TABLE `gs_snapshot`
  ADD COLUMN `sale_num` INT NULL
  AFTER `level_sales`,
  ADD COLUMN `stock_num` INT NULL
  AFTER `sale_num`;


ALTER TABLE `pt_order`
  CHANGE COLUMN `member_id` `buyner_id` INT(11) NULL DEFAULT NULL
COMMENT '-1代表游客订餐',
  CHANGE COLUMN `merchant_id` `vendor_id` INT(11) NULL DEFAULT NULL
COMMENT '商家ID';

CREATE TABLE `pay_order` (
  `id`                INT           NOT NULL AUTO_INCREMENT,
  `trade_no`          VARCHAR(45)   NULL,
  `vendor_id`         INT           NULL,
  `order_id`          INT           NULL,
  `buy_user`          INT           NULL,
  `payment_user`      INT           NULL,
  `total_fee`         DECIMAL(8, 2) NULL,
  `balance_discount`  DECIMAL(8, 2) NULL,
  `integral_discount` DECIMAL(8, 2) NULL,
  `system_discount`   DECIMAL(8, 2) NULL,
  `coupon_discount`   DECIMAL(8, 2) NULL,
  `sub_fee`           DECIMAL(8, 2) NULL,
  `final_fee`         DECIMAL(8, 2) NULL,
  `payment_opt`       TINYINT(2)    NULL,
  `payment_sign`      TINYINT(1)    NULL,
  `outer_no`          VARCHAR(45)   NULL
  COMMENT '外部订单号',
  `create_time`       INT           NULL,
  `paid_time`         INT           NULL,
  `state`             VARCHAR(45)   NULL,
  PRIMARY KEY (`id`)
);


ALTER TABLE `pt_order`
  CHANGE COLUMN `buyner_id` `buyer_id` INT(11) NULL DEFAULT NULL
COMMENT '-1代表游客订餐';

ALTER TABLE `pt_order_item`
  CHANGE COLUMN `snapshot_id` `snap_id` INT(11) NULL DEFAULT NULL,
  CHANGE COLUMN `quantity` `quantity` INT(6) NULL DEFAULT NULL,
  CHANGE COLUMN `update_time` `update_time` INT NULL DEFAULT NULL,
  ADD COLUMN `vendor_id` INT NULL
  AFTER `order_id`,
  ADD COLUMN `shop_id` INT NULL
  AFTER `vendor_id`,
  ADD COLUMN `sku_id` INT NULL
  AFTER `shop_id`,
  ADD COLUMN `final_fee` DECIMAL(8, 2) NULL
  AFTER `fee`;

CREATE TABLE `sale_order` (
  `id`               INT           NOT NULL AUTO_INCREMENT,
  `order_no`         VARCHAR(20)   NULL,
  `buyer_id`         INT           NULL,
  `items_info`       VARCHAR(255)  NULL,
  `total_fee`        DECIMAL(8, 2) NULL,
  `discount_fee`     DECIMAL(8, 2) NULL,
  `final_fee`        DECIMAL(8, 2) NULL,
  `is_paid`          TINYINT(1)    NULL,
  `paid_time`        INT           NULL,
  `consignee_person` VARCHAR(45)   NULL,
  `consignee_phone`  VARCHAR(45)   NULL,
  `shipping_address` VARCHAR(120)  NULL,
  `shipping_time`    VARCHAR(45)   NULL,
  `create_time`      INT           NULL,
  `update_time`      INT           NULL,
  `status`           TINYINT(1)    NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `pt_order_pb`
  CHANGE COLUMN `order_no` `order_id` INT NULL DEFAULT NULL
  AFTER `id`;


CREATE TABLE `sale_sub_order` (
  `id`           INT           NOT NULL AUTO_INCREMENT,
  `order_no`     VARCHAR(20)   NULL,
  `parent_order` INT           NULL,
  `vendor_id`    INT           NULL,
  `shop_id`      INT           NULL,
  `subject`      VARCHAR(45)   NULL,
  `items_info`   VARCHAR(255)  NULL,
  `total_fee`    DECIMAL(8, 2) NULL,
  `discount_fee` DECIMAL(8, 2) NULL,
  `final_fee`    DECIMAL(8, 2) NULL,
  `is_suspend`   TINYINT(1)    NULL,
  `note`         VARCHAR(120)  NULL,
  `remark`       VARCHAR(120)  NULL,
  `update_time`  INT           NULL,
  `status`       TINYINT(1)    NULL,
  PRIMARY KEY (`id`)
);


ALTER TABLE `pt_order_item`
  CHANGE COLUMN `update_time` `update_time` INT NULL DEFAULT NULL,
  ADD COLUMN `sku_id` INT NULL
  AFTER `order_id`,
  ADD COLUMN `final_fee` DECIMAL(8, 2) NULL
  AFTER `fee`,
RENAME TO `sale_order_item`;

ALTER TABLE `sale_order_item`
  CHANGE COLUMN `snapshot_id` `snap_id` INT(11) NULL DEFAULT NULL;

ALTER TABLE `sale_order_item`
  DROP COLUMN `sku`;


CREATE TABLE `express_template` (
  `id`         INT           NOT NULL AUTO_INCREMENT,
  `user_id`    INT           NULL,
  `name`       VARCHAR(45)   NULL,
  `is_free`    TINYINT(1)    NULL,
  `basis`      TINYINT(1)    NULL,
  `first_unit` INT(5)        NULL,
  `first_fee`  DECIMAL(6, 2) NULL,
  `add_unit`   INT(5)        NULL,
  `add_fee`    DECIMAL(6, 2) NULL,
  `enabled`    TINYINT(1)    NULL,
  PRIMARY KEY (`id`)
);


CREATE TABLE `express_area_set` (
  `id`          INT           NOT NULL AUTO_INCREMENT,
  `template_id` INT           NULL,
  `code_list`   VARCHAR(500)  NULL,
  `name_list`   VARCHAR(120)  NULL,
  `first_unit`  INT(5)        NULL,
  `first_fee`   DECIMAL(6, 2) NULL,
  `add_unit`    INT(5)        NULL,
  `add_fee`     DECIMAL(6, 2) NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `express_provider` (
  `id`       INT         NOT NULL AUTO_INCREMENT,
  `name`     VARCHAR(45) NULL,
  `letter`   VARCHAR(1)  NULL,
  `code`     VARCHAR(10) NULL,
  `api_code` VARCHAR(10) NULL,
  `enabled`  TINYINT(1)  NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `gs_item`
  ADD COLUMN `weight` INT NULL
COMMENT '重量,单位:克(g)'
  AFTER `cost`;

ALTER TABLE `gs_snapshot`
  ADD COLUMN `weight` INT NULL
COMMENT '单件重量,单位:克(g)'
  AFTER `img`;

ALTER TABLE `sale_order`
  ADD COLUMN `express_fee` DECIMAL(8, 2) NULL
COMMENT '物流费'
  AFTER `discount_fee`;


ALTER TABLE `mm_member`
  ADD COLUMN `check_code` VARCHAR(8) NULL
  AFTER `reg_time`,
  ADD COLUMN `check_expires` INT NULL
  AFTER `check_code`;


CREATE TABLE `gs_sales_snapshot` (
  `id`          INT           NOT NULL AUTO_INCREMENT,
  `snap_key`    VARCHAR(45)   NULL,
  `sku_id`      INT           NULL,
  `seller_id`   INT           NULL,
  `item_id`     INT           NULL,
  `cat_id`      INT           NULL,
  `goods_title` VARCHAR(120)  NULL,
  `goods_no`    VARCHAR(45)   NULL,
  `sku`         VARCHAR(120)  NULL,
  `img`         VARCHAR(120)  NULL,
  `price`       DECIMAL(8, 2) NULL,
  `create_time` INT           NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `sale_order`
  CHANGE COLUMN `total_fee` `goods_fee` DECIMAL(8, 2) NULL DEFAULT NULL
COMMENT '商品金额';


ALTER TABLE `sale_sub_order`
  CHANGE COLUMN `total_fee` `goods_fee` DECIMAL(8, 2) NULL DEFAULT NULL,
  ADD COLUMN `express_fee` DECIMAL(4, 2) NULL
  AFTER `discount_fee`;

ALTER TABLE `sale_order`
  ADD COLUMN `package_fee` DECIMAL(4, 2) NULL
  AFTER `express_fee`;

ALTER TABLE `sale_sub_order`
  ADD COLUMN `package_fee` DECIMAL(4, 2) NULL
  AFTER `express_fee`;


ALTER TABLE `sale_sub_order`
  ADD COLUMN `buyer_id` INT NULL
  AFTER `parent_order`;

ALTER TABLE `pt_order_log`
  ADD COLUMN `order_state` TINYINT(2) NULL
  AFTER `type`,
RENAME TO `sale_order_log`;

ALTER TABLE `sale_sub_order`
  ADD COLUMN `is_paid` TINYINT(1) NULL
  AFTER `final_fee`;

ALTER TABLE `sale_sub_order`
  CHANGE COLUMN `status` `state` TINYINT(1) NULL DEFAULT NULL;

ALTER TABLE `sale_order`
  CHANGE COLUMN `status` `state` TINYINT(1) NULL DEFAULT NULL;

ALTER TABLE `mch_enterprise_info`
  ADD COLUMN `person_id` VARCHAR(20) NULL
COMMENT '法人身份证号'
  AFTER `person_name`;

ALTER TABLE `mch_enterprise_info`
  ADD COLUMN `is_handled` TINYINT(1) NULL
  AFTER `company_imageurl`;

CREATE TABLE `ship_order` (
  `id`           INT           NOT NULL AUTO_INCREMENT,
  `order_id`     INT           NULL,
  `sp_id`        INT           NULL
  COMMENT '快递SP编号',
  `sp_order`     VARCHAR(20)   NULL
  COMMENT '快递SP单号',
  `exporess_log` VARCHAR(512)  NULL,
  `amount`       DECIMAL(8, 2) NULL,
  `final_amount` DECIMAL(8, 2) NULL,
  `ship_time`    INT           NULL
  COMMENT '发货时间',
  `state`        TINYINT(1)    NULL
  COMMENT '是否已收货',
  `update_time`  INT           NULL,
  PRIMARY KEY (`id`)
);


CREATE TABLE `ship_item` (
  `id`           INT           NOT NULL AUTO_INCREMENT,
  `ship_order`   INT           NULL,
  `snap_id`      INT           NULL,
  `quantity`     INT           NULL,
  `amount`       DECIMAL(8, 2) NULL,
  `final_amount` DECIMAL(8, 2) NULL,
  PRIMARY KEY (`id`)
);
ALTER TABLE `sale_order`
  CHANGE COLUMN `goods_fee` `goods_amount` DECIMAL(8, 2) NULL DEFAULT NULL
COMMENT '商品金额',
  CHANGE COLUMN `discount_fee` `discount_amount` DECIMAL(8, 2) NULL DEFAULT NULL,
  CHANGE COLUMN `final_fee` `final_amount` DECIMAL(8, 2) NULL DEFAULT NULL;

ALTER TABLE `sale_order_item`
  CHANGE COLUMN `fee` `amount` DECIMAL(8, 2) NULL DEFAULT NULL,
  CHANGE COLUMN `final_fee` `final_amount` DECIMAL(8, 2) NULL DEFAULT NULL;

ALTER TABLE `sale_sub_order`
  CHANGE COLUMN `goods_fee` `goods_amount` DECIMAL(8, 2) NULL DEFAULT NULL,
  CHANGE COLUMN `discount_fee` `discount_amount` DECIMAL(8, 2) NULL DEFAULT NULL,
  CHANGE COLUMN `final_fee` `final_amount` DECIMAL(8, 2) NULL DEFAULT NULL;

ALTER TABLE `sale_order_item`
  ADD COLUMN `is_shipped` TINYINT(1) NULL
  AFTER `final_amount`;

ALTER TABLE `mm_integral_log`
  CHANGE COLUMN `partner_id` `mch_id` INT(11) NULL DEFAULT NULL;

ALTER TABLE `pay_order`
  CHANGE COLUMN `sub_fee` `sub_amount` DECIMAL(8, 2) NULL DEFAULT NULL
COMMENT '立减金额',
  ADD COLUMN `adjustment_amount` DECIMAL(8, 2) NULL
COMMENT '调整金额'
  AFTER `sub_amount`;


CREATE TABLE `sale_after_order` (
  `id`            INT          NOT NULL AUTO_INCREMENT,
  `order_id`      INT          NULL,
  `vendor_id`     INT          NULL,
  `buyer_id`      INT          NULL,
  `type`          TINYINT(1)   NULL,
  `snap_id`       INT          NULL,
  `quantity`      INT          NULL,
  `reason`        VARCHAR(255) NULL,
  `person_name`   VARCHAR(10)  NULL,
  `person_phone`  VARCHAR(20)  NULL,
  `rsp_name`      VARCHAR(10)  NULL
  COMMENT '退货快递名称',
  `rsp_order`     VARCHAR(20)  NULL
  COMMENT '退货快递单号',
  `rsp_image`     VARCHAR(120) NULL,
  `remark`        VARCHAR(45)  NULL,
  `vendor_remark` VARCHAR(45)  NULL,
  `state`         TINYINT(1)   NULL,
  `create_time`   INT          NULL,
  `update_time`   INT          NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `sale_return` (
  `id`        INT           NOT NULL,
  `amount`    DECIMAL(8, 2) NULL,
  `is_refund` TINYINT(1)    NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `sale_exchange` (
  `id`           INT         NOT NULL,
  `is_shipped`   TINYINT(1)  NULL,
  `sp_name`      VARCHAR(20) NULL,
  `sp_order`     VARCHAR(20) NULL,
  `ship_time`    INT         NULL,
  `is_received`  TINYINT(1)  NULL,
  `receive_time` INT         NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `sale_refund` (
  `id`        INT           NOT NULL,
  `amount`    DECIMAL(8, 2) NULL,
  `is_refund` TINYINT(1)    NULL,
  PRIMARY KEY (`id`)
);

ALTER TABLE `sale_order_item`
  ADD COLUMN `return_quantity` INT NULL
  AFTER `quantity`;


ALTER TABLE `sale_sub_order`
  ADD COLUMN `create_time` INT NULL
  AFTER `remark`;

ALTER TABLE `mm_trusted_info`
  CHANGE COLUMN `body_number` `card_id` VARCHAR(20)
CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL;


ALTER TABLE `pay_order`
  ADD COLUMN `order_type` INT NULL
COMMENT '支付单的类型，如购物或其他'
  AFTER `vendor_id`;

ALTER TABLE `pay_order`
  ADD COLUMN `subject` VARCHAR(45) NULL
COMMENT '支付单标题'
  AFTER `order_id`;


ALTER TABLE `mm_integral_log`
  DROP COLUMN `mch_id`,
  CHANGE COLUMN `member_id` `member_id` INT(11) NOT NULL,
  CHANGE COLUMN `type` `type` INT(11) NOT NULL,
  CHANGE COLUMN `integral` `value` INT(11) NOT NULL,
  CHANGE COLUMN `log` `remark` VARCHAR(100) NULL DEFAULT NULL,
  CHANGE COLUMN `record_time` `create_time` INT(11) NOT NULL;

ALTER TABLE `mm_integral_log`
  ADD COLUMN `outer_no` VARCHAR(45) NULL
  AFTER `type`;


ALTER TABLE `mm_account`
  CHANGE COLUMN `freezes_fee` `freezes_balance` FLOAT(10, 2) NOT NULL
  AFTER `balance`,
  CHANGE COLUMN `freezes_present` `freezes_present` FLOAT(10, 2) NOT NULL
  AFTER `present_balance`,
  CHANGE COLUMN `total_fee` `total_consumption` FLOAT(10, 2) NOT NULL
COMMENT '总消费'
  AFTER `total_pay`,
  CHANGE COLUMN `integral` `integral` INT(11) NOT NULL,
  CHANGE COLUMN `balance` `balance` FLOAT(10, 2) NOT NULL,
  CHANGE COLUMN `present_balance` `present_balance` FLOAT(10, 2) NOT NULL,
  CHANGE COLUMN `total_present_fee` `total_present_fee` FLOAT(10, 2) NOT NULL,
  CHANGE COLUMN `flow_balance` `flow_balance` FLOAT(10, 2) NOT NULL,
  CHANGE COLUMN `grow_balance` `grow_balance` FLOAT(10, 2) NOT NULL,
  CHANGE COLUMN `grow_amount` `grow_amount` FLOAT(10, 2) NOT NULL,
  CHANGE COLUMN `grow_earnings` `grow_earnings` FLOAT(10, 2) NOT NULL,
  CHANGE COLUMN `grow_total_earnings` `grow_total_earnings` FLOAT(10, 2) NOT NULL,
  CHANGE COLUMN `total_charge` `total_charge` FLOAT(10, 2) NOT NULL,
  CHANGE COLUMN `total_pay` `total_pay` FLOAT(10, 2) NOT NULL,
  CHANGE COLUMN `update_time` `update_time` INT(11) NOT NULL
COMMENT '积分',
  ADD COLUMN `freezes_integral` INT NOT NULL
COMMENT '不可用积分'
  AFTER `integral`;


ALTER TABLE `con_page`
  CHANGE COLUMN `mch_id` `user_id` INT(11) NULL DEFAULT NULL;

DROP TABLE `gs_member_price`;

CREATE TABLE `gs_member_price` (
  `id`        INT(11)       NOT NULL AUTO_INCREMENT,
  `goods_id`  INT(11)       NOT NULL,
  `level`     INT(11)       NOT NULL,
  `price`     DECIMAL(8, 2) NOT NULL,
  `max_quota` INT(11)       NOT NULL,
  `enabled`   TINYINT(1)    NOT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;


ALTER TABLE `express_provider`
  CHANGE COLUMN `letter` `group_flag` VARCHAR(45)
CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL;


ALTER TABLE `gs_item`
  ADD COLUMN `express_tid` INT NULL
COMMENT '快递模板编号'
  AFTER `sale_price`;

ALTER TABLE `gs_snapshot`
  ADD COLUMN `express_tid` INT NULL
  AFTER `sale_price`;

ALTER TABLE `gs_item`
  CHANGE COLUMN `img` `img` VARCHAR(120) NULL DEFAULT NULL;

ALTER TABLE `gs_item`
  CHANGE COLUMN `weight` `weight` FLOAT(6, 2) NULL DEFAULT NULL
COMMENT '重量,单位:克(g)';

ALTER TABLE `gs_snapshot`
  CHANGE COLUMN `weight` `weight` FLOAT(6, 2) NULL DEFAULT NULL
COMMENT '单件重量,单位:克(g)';

ALTER TABLE `gs_snapshot`
  ADD COLUMN `cost` DECIMAL(8, 2) NULL
  AFTER `weight`;


ALTER TABLE `gs_sales_snapshot`
  ADD COLUMN `cost` DECIMAL(8, 2) NULL
COMMENT '供货价'
  AFTER `img`;


CREATE TABLE mch_account (
  mch_id         INT(10)        NOT NULL AUTO_INCREMENT
  COMMENT '商户编号',
  balance        DECIMAL(10, 2) NOT NULL
  COMMENT '余额',
  freeze_amount  DECIMAL(10, 2) NOT NULL
  COMMENT '冻结金额',
  await_amount   DECIMAL(10, 2) NOT NULL
  COMMENT '待入账金额',
  present_amount DECIMAL(10, 2) NOT NULL
  COMMENT '平台赠送金额',
  sales_amount   DECIMAL(10, 2) NOT NULL
  COMMENT '累计销售总额',
  refund_amount  DECIMAL(10, 2) NOT NULL
  COMMENT '累计退款金额',
  take_amount    DECIMAL(10, 2) NOT NULL
  COMMENT '已提取金额',
  offline_sales  DECIMAL(10, 2) NOT NULL
  COMMENT '线下销售金额',
  update_time    INT(11)        NOT NULL
  COMMENT '更新时间',
  PRIMARY KEY (mch_id)
)
  COMMENT = '商户账户表';


CREATE TABLE mch_balance_log (
  id          INT(10)            NOT NULL AUTO_INCREMENT,
  mch_id      INT(10)            NOT NULL
  COMMENT '商户编号',
  kind        INT(10)            NOT NULL
  COMMENT '日志类型',
  title       VARCHAR(45)        NOT NULL
  COMMENT '标题',
  outer_no    VARCHAR(45)        NOT NULL
  COMMENT '外部订单号',
  amount      FLOAT              NOT NULL
  COMMENT '金额',
  csn_amount  FLOAT DEFAULT 0.00 NOT NULL
  COMMENT '手续费',
  state       TINYINT(1)         NOT NULL
  COMMENT '状态',
  create_time INT(10)            NOT NULL
  COMMENT '创建时间',
  update_time INT(10)            NOT NULL
  COMMENT '更新时间',
  PRIMARY KEY (id)
)
  COMMENT = '商户余额日志';

CREATE TABLE mch_day_chart (
  id              INT(11)        NOT NULL AUTO_INCREMENT
  COMMENT '编号',
  mch_id          INT(11)        NOT NULL
  COMMENT '商户编号',
  order_number    INT(11)        NOT NULL
  COMMENT '新增订单数量',
  order_amount    DECIMAL(10, 2) NOT NULL
  COMMENT '订单额',
  buyer_number    INT(11)        NOT NULL
  COMMENT '购物会员数',
  paid_number     INT(11)        NOT NULL
  COMMENT '支付单数量',
  paid_amount     DECIMAL(10, 2) NOT NULL
  COMMENT '支付总金额',
  complete_orders INT(11)        NOT NULL
  COMMENT '完成订单数',
  in_amount       DECIMAL(10, 2) NOT NULL
  COMMENT '入帐金额',
  offline_orders  INT(11)        NOT NULL
  COMMENT '线下订单数量',
  offline_amount  DECIMAL(10, 2) NOT NULL
  COMMENT '线下订单金额',
  `date`          INT(11)        NOT NULL
  COMMENT '日期',
  date_str        VARCHAR(10)    NOT NULL
  COMMENT '日期字符串',
  update_time     INT(11)        NOT NULL
  COMMENT '更新时间',
  PRIMARY KEY (id)
)
  COMMENT = '商户每日报表';


ALTER TABLE `mch_enterprise_info`
  CHANGE COLUMN `person_imageurl` `person_image` VARCHAR(120)
CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL,
  CHANGE COLUMN `company_imageurl` `company_image` VARCHAR(120)
CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL,
  ADD COLUMN `auth_doc` VARCHAR(120) NULL
COMMENT '授权书'
  AFTER `company_image`;


CREATE TABLE mch_sign_up (
  id            INT(11)      NOT NULL AUTO_INCREMENT,
  member_id     INT(11)      NOT NULL,
  sign_no       VARCHAR(20)  NOT NULL,
  usr           VARCHAR(45)  NOT NULL,
  pwd           VARCHAR(45)  NOT NULL,
  mch_name      VARCHAR(20)  NOT NULL,
  province      INT(10)      NOT NULL,
  city          INT(10)      NOT NULL,
  district      INT(10)      NOT NULL,
  shop_name     VARCHAR(20)  NOT NULL,
  company_name  VARCHAR(20)  NOT NULL,
  company_no    VARCHAR(20)  NOT NULL,
  person_name   VARCHAR(10)  NOT NULL,
  person_id     VARCHAR(20)  NOT NULL,
  phone         VARCHAR(20)  NOT NULL,
  address       VARCHAR(120) NOT NULL,
  person_image  VARCHAR(120) NOT NULL,
  company_image VARCHAR(120) NOT NULL,
  auth_doc      VARCHAR(120) NOT NULL,
  reviewed      INT(1)       NOT NULL,
  remark        VARCHAR(120) NOT NULL,
  submit_time   INT(11)      NOT NULL,
  update_time   INT(11)      NOT NULL,
  PRIMARY KEY (id)
);


ALTER TABLE `msg_to`
  CHANGE COLUMN `to_id` `to_id` INT(11) NOT NULL,
  CHANGE COLUMN `to_role` `to_role` TINYINT(2) NOT NULL,
  CHANGE COLUMN `content_id` `content_id` INT(11) NOT NULL,
  CHANGE COLUMN `has_read` `has_read` TINYINT(1) NOT NULL,
  CHANGE COLUMN `read_time` `read_time` INT(11) NOT NULL,
  ADD COLUMN `msg_id` INT(11) NOT NULL
  AFTER `to_role`;


ALTER TABLE `mm_account`
  ADD COLUMN `priority_pay` TINYINT(1) NULL
COMMENT '优先（默认）支付账户'
  AFTER `total_consumption`;

CREATE TABLE `mm_balance_log` (
  `id`          INT(11)     NOT NULL AUTO_INCREMENT,
  `member_id`   INT(11)     NOT NULL,
  `kind`        TINYINT(2)  NOT NULL
  COMMENT '业务类型',
  `title`       VARCHAR(45) NOT NULL
  COMMENT '标题',
  `outer_no`    VARCHAR(45) NOT NULL
  COMMENT '外部订单号',
  `amount`      FLOAT(8, 2) NOT NULL
  COMMENT '金额',
  `csn_fee`     FLOAT(8, 2) NOT NULL
  COMMENT '手续费',
  `state`       TINYINT(1)  NOT NULL
  COMMENT '状态，比如提现需要确认等',
  `rel_user`    INT(11)     NOT NULL
  COMMENT '关联操作人员编号',
  `remark`      VARCHAR(45)          DEFAULT NULL,
  `create_time` INT(11)     NOT NULL,
  `update_time` INT(11)     NOT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;


CREATE TABLE `mm_present_log` (
  `id`          INT(11)     NOT NULL AUTO_INCREMENT,
  `member_id`   INT(11)     NOT NULL,
  `kind`        TINYINT(2)  NOT NULL
  COMMENT '业务类型',
  `title`       VARCHAR(45) NOT NULL
  COMMENT '标题',
  `outer_no`    VARCHAR(45) NOT NULL
  COMMENT '外部订单号',
  `amount`      FLOAT(8, 2) NOT NULL
  COMMENT '金额',
  `csn_fee`     FLOAT(8, 2) NOT NULL
  COMMENT '手续费',
  `state`       TINYINT(1)  NOT NULL
  COMMENT '状态，比如提现需要确认等',
  `rel_user`    INT(11)     NOT NULL
  COMMENT '关联操作人员编号',
  `remark`      VARCHAR(45)          DEFAULT NULL,
  `create_time` INT(11)     NOT NULL,
  `update_time` INT(11)     NOT NULL,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;


ALTER TABLE `mm_account`
  CHANGE COLUMN `freezes_integral` `freeze_integral` INT(11) NOT NULL
COMMENT '不可用积分',
  CHANGE COLUMN `balance` `balance` DECIMAL(10, 2) NOT NULL
COMMENT '余额',
  CHANGE COLUMN `freezes_balance` `freeze_balance` DECIMAL(10, 2) NOT NULL
COMMENT '冻结的账户余额',
  CHANGE COLUMN `present_balance` `present_balance` DECIMAL(10, 2) NOT NULL,
  CHANGE COLUMN `freezes_present` `freeze_present` DECIMAL(10, 2) NOT NULL
COMMENT '冻结的赠送金额',
  CHANGE COLUMN `total_present_fee` `total_present_fee` DECIMAL(10, 2) NOT NULL,
  CHANGE COLUMN `flow_balance` `flow_balance` DECIMAL(10, 2) NOT NULL,
  CHANGE COLUMN `grow_balance` `grow_balance` DECIMAL(10, 2) NOT NULL,
  CHANGE COLUMN `grow_amount` `grow_amount` DECIMAL(10, 2) NOT NULL,
  CHANGE COLUMN `grow_earnings` `grow_earnings` DECIMAL(10, 2) NOT NULL,
  CHANGE COLUMN `grow_total_earnings` `grow_total_earnings` DECIMAL(10, 2) NOT NULL,
  CHANGE COLUMN `total_charge` `total_charge` DECIMAL(10, 2) NOT NULL,
  CHANGE COLUMN `total_pay` `total_pay` DECIMAL(10, 2) NOT NULL,
  CHANGE COLUMN `total_consumption` `total_consumption` DECIMAL(10, 2) NOT NULL
COMMENT '总消费',
  CHANGE COLUMN `priority_pay` `priority_pay` DECIMAL(10, 2) NOT NULL
COMMENT '优先（默认）支付账户',
  ADD COLUMN `expired_balance` DECIMAL(10, 2) NOT NULL
COMMENT '失效的账户余额'
  AFTER `freeze_balance`,
  ADD COLUMN `expired_present` DECIMAL(10, 2) NOT NULL
COMMENT '失效的赠送金额'
  AFTER `freeze_present`;


ALTER TABLE gs_item
  DROP COLUMN `review_pass`,
  DROP COLUMN `apply_subs`,
  CHANGE COLUMN `on_shelves` `shelve_state` TINYINT(1) NULL DEFAULT NULL
COMMENT '是否上架'
  AFTER `state`,
  CHANGE COLUMN `remark` `remark` VARCHAR(255) NULL DEFAULT NULL
COMMENT '备注'
  AFTER `review_state`,
  CHANGE COLUMN `has_review` `review_state` TINYINT(1) NULL DEFAULT NULL;

ALTER TABLE `gs_snapshot`
  CHANGE COLUMN `on_shelves` `shelve_state` TINYINT(1) NULL DEFAULT '1'
  AFTER `stock_num`;

ALTER TABLE `mm_present_log`
  CHANGE COLUMN `amount` `amount` FLOAT(12, 2) NOT NULL
COMMENT '金额';

ALTER TABLE `mm_balance_log`
  CHANGE COLUMN `amount` `amount` FLOAT(12, 2) NOT NULL
COMMENT '金额';


ALTER TABLE mm_member
  ADD COLUMN `login_time` INT(11) NULL
  AFTER `check_expires`;
















