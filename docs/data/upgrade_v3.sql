CREATE TABLE `mm_levelup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) NOT NULL COMMENT '会员编号',
  `origin_level` tinyint(2) NOT NULL COMMENT '原来等级',
  `target_level` tinyint(2) NOT NULL COMMENT '现在等级',
  `is_free` int(1) NOT NULL COMMENT '是否为免费升级的会员',
  `payment_id` int(11) NOT NULL COMMENT '支付单编号',
  `reviewed` int(1) NOT NULL COMMENT '是否审核及处理',
  `create_time` int(11) NOT NULL COMMENT '升级时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='升级日志表';

CREATE TABLE `comm_qr_template` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `title` varchar(45) COLLATE utf8_unicode_ci NOT NULL COMMENT '模板标题',
  `bg_image` varchar(120) COLLATE utf8_unicode_ci NOT NULL COMMENT '背景图片',
  `offset_x` int(11) NOT NULL COMMENT '垂直偏离量',
  `offset_y` int(11) NOT NULL COMMENT '垂直偏移量',
  `comment` varchar(120) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '二维码模板文本',
  `enabled` int(1) DEFAULT NULL COMMENT '是否启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='二维码模板';

ALTER TABLE `go2o`.`comm_qr_template`
ADD COLUMN `callback_url` VARCHAR(120) NULL COMMENT '回调地址' AFTER `comment`;


ALTER TABLE `mm_member`
CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '编号' ,
CHANGE COLUMN `trade_pwd` `trade_pwd` VARCHAR(45) NULL DEFAULT NULL COMMENT '交易密码' ,
CHANGE COLUMN `exp` `exp` INT(11) UNSIGNED NULL DEFAULT '0' COMMENT '经验值' ,
CHANGE COLUMN `level` `level` INT(11) NULL DEFAULT '1' COMMENT '等级' ,
CHANGE COLUMN `reg_ip` `reg_ip` VARCHAR(20) NULL DEFAULT NULL COMMENT '注册IP' ,
CHANGE COLUMN `reg_from` `reg_from` VARCHAR(20) NULL DEFAULT NULL COMMENT '注册来源' ,
CHANGE COLUMN `reg_time` `reg_time` INT(11) NULL DEFAULT NULL COMMENT '注册时间' ,
CHANGE COLUMN `check_code` `check_code` VARCHAR(8) NULL DEFAULT NULL COMMENT '校验码' ,
CHANGE COLUMN `check_expires` `check_expires` INT(11) NULL DEFAULT NULL COMMENT '校验码过期时间' ,
CHANGE COLUMN `login_time` `login_time` INT(11) NULL DEFAULT NULL COMMENT '登陆时间' ,
CHANGE COLUMN `state` `state` INT(1) NULL DEFAULT '1' COMMENT '状态' ,
CHANGE COLUMN `update_time` `update_time` INT(11) NULL DEFAULT NULL COMMENT '更新时间' ;

ALTER TABLE `mm_account`
CHANGE COLUMN `integral` `integral` INT(11) NOT NULL DEFAULT '0' COMMENT '积分' ,
CHANGE COLUMN `freeze_integral` `freeze_integral` INT(11) NOT NULL COMMENT '冻结积分' ,
CHANGE COLUMN `freeze_balance` `freeze_balance` DECIMAL(10,2) NOT NULL COMMENT '冻结账户余额' ,
CHANGE COLUMN `expired_balance` `expired_balance` DECIMAL(10,2) NOT NULL COMMENT '失效账户余额' ,
CHANGE COLUMN `present_balance` `present_balance` DECIMAL(10,2) NOT NULL DEFAULT '0.00' COMMENT '赠送余额' ,
CHANGE COLUMN `total_present_fee` `total_present_amount` DECIMAL(10,2) NOT NULL DEFAULT '0.00' COMMENT '累计赠送金额' ,
CHANGE COLUMN `flow_balance` `flow_balance` DECIMAL(10,2) NOT NULL DEFAULT '0.00' COMMENT '浮动账户余额' ,
CHANGE COLUMN `grow_balance` `grow_balance` DECIMAL(10,2) NOT NULL COMMENT '增长账户余额' ,
CHANGE COLUMN `grow_amount` `grow_amount` DECIMAL(10,2) NOT NULL DEFAULT '0.00' COMMENT '浮动账户余额总投资金额,不含收益' ,
CHANGE COLUMN `grow_earnings` `grow_earnings` DECIMAL(10,2) NOT NULL COMMENT '当前收益金额' ,
CHANGE COLUMN `grow_total_earnings` `grow_total_earnings` DECIMAL(10,2) NOT NULL COMMENT '累积收益金额' ,
CHANGE COLUMN `total_charge` `total_charge` DECIMAL(10,2) NOT NULL DEFAULT '0.00' COMMENT '总充值金额' ,
CHANGE COLUMN `total_pay` `total_pay` DECIMAL(10,2) NOT NULL DEFAULT '0.00' COMMENT '总支付额' ,
CHANGE COLUMN `total_consumption` `total_consumption` DECIMAL(10,2) NOT NULL COMMENT '总消费金额' ,
CHANGE COLUMN `update_time` `update_time` INT(11) NOT NULL DEFAULT '0' COMMENT '更新时间' ;

ALTER TABLE `mm_profile`
DROP COLUMN `qq`,
CHANGE COLUMN `email` `email` VARCHAR(50) NULL DEFAULT NULL COMMENT '电子邮件' AFTER `im`,
CHANGE COLUMN `province` `province` INT(8) NULL DEFAULT NULL COMMENT '省' AFTER `email`,
CHANGE COLUMN `city` `city` INT(8) NULL DEFAULT NULL COMMENT '市' AFTER `province`,
CHANGE COLUMN `district` `district` INT(8) NULL DEFAULT NULL COMMENT '区' AFTER `city`,
CHANGE COLUMN `avatar` `avatar` VARCHAR(80) NULL DEFAULT NULL COMMENT '头像' ,
CHANGE COLUMN `birthday` `birthday` VARCHAR(20) NULL DEFAULT NULL COMMENT '生日' ,
CHANGE COLUMN `phone` `phone` VARCHAR(15) NULL DEFAULT NULL COMMENT '电话' ,
CHANGE COLUMN `address` `address` VARCHAR(100) NULL DEFAULT NULL COMMENT '地址' ,
CHANGE COLUMN `im` `im` VARCHAR(45) NULL DEFAULT NULL COMMENT '即时通讯' ,
CHANGE COLUMN `ext_1` `ext_1` VARCHAR(45) NULL DEFAULT NULL COMMENT '扩展1' ,
CHANGE COLUMN `ext_2` `ext_2` VARCHAR(45) NULL DEFAULT NULL COMMENT '扩展2' ,
CHANGE COLUMN `ext_3` `ext_3` VARCHAR(45) NULL DEFAULT NULL COMMENT '扩展3' ,
CHANGE COLUMN `ext_4` `ext_4` VARCHAR(45) NULL DEFAULT NULL COMMENT '扩展4' ,
CHANGE COLUMN `ext_5` `ext_5` VARCHAR(45) NULL DEFAULT NULL COMMENT '扩展5' ,
CHANGE COLUMN `ext_6` `ext_6` VARCHAR(45) NULL DEFAULT NULL COMMENT '扩展6' ,
CHANGE COLUMN `remark` `remark` VARCHAR(100) NULL DEFAULT NULL COMMENT '备注' ,
CHANGE COLUMN `update_time` `update_time` INT(11) NULL DEFAULT NULL COMMENT '更新时间' ,
COMMENT = '会员资料' ;


ALTER TABLE `mm_relation`
CHANGE COLUMN `member_id` `member_id` INT(11) NOT NULL COMMENT '会员编号' ,
CHANGE COLUMN `card_id` `card_no` VARCHAR(20) NULL DEFAULT NULL ,
CHANGE COLUMN `invi_member_id` `inviter_id` INT(11) NOT NULL COMMENT '邀请人(会员)编号' ,
CHANGE COLUMN `refer_str` `inviter_str` VARCHAR(250) NULL DEFAULT NULL COMMENT '邀请人(会员）字符表示' ,
CHANGE COLUMN `reg_merchant_id` `reg_mchid` INT(11) NULL DEFAULT NULL COMMENT '注册关联的商户编号' ,
COMMENT = '会员关系' ;


ALTER TABLE `gs_category`
DROP COLUMN `description`,
DROP COLUMN `mch_id`,
CHANGE COLUMN `icon` `icon` VARCHAR(150) NULL COMMENT '分类图片' AFTER `name`,
CHANGE COLUMN `level` `level` INT(3) NOT NULL COMMENT '分类层级' AFTER `url`,
CHANGE COLUMN `sort_number` `sort_number` INT(11) NOT NULL COMMENT '排序序号' AFTER `level`,
CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '编号' ,
CHANGE COLUMN `parent_id` `parent_id` INT(11) NOT NULL COMMENT '父分类' ,
CHANGE COLUMN `name` `name` VARCHAR(20) NOT NULL COMMENT '分类名称' ,
CHANGE COLUMN `url` `url` VARCHAR(120) NULL COMMENT '品牌链接地址' ,
CHANGE COLUMN `enabled` `enabled` INT(1) NOT NULL COMMENT '是否启用' ,
CHANGE COLUMN `create_time` `create_time` INT(11) NOT NULL COMMENT '创建时间' ,
ADD COLUMN `spec_model` INT(11) NULL COMMENT '商品规格模型' AFTER `parent_id`,
COMMENT = '商品分类' , RENAME TO  `cat_category` ;


ALTER TABLE `cat_category`
CHANGE COLUMN `spec_model` `pro_model` INT(11) NOT NULL DEFAULT 0 COMMENT '产品模型' AFTER `parent_id`,
CHANGE COLUMN `sort_number` `sort_num` INT(11) NOT NULL COMMENT '排序序号' ;


/** update **/

ALTER TABLE `gs_item`
CHANGE COLUMN `express_tid` `express_tid` INT(11) NULL COMMENT '快递模板编号' AFTER `supplier_id`,
CHANGE COLUMN `goods_no` `goods_no` VARCHAR(45) NULL DEFAULT NULL AFTER `small_title`,
CHANGE COLUMN `weight` `weight` FLOAT(6,2) NULL DEFAULT NULL COMMENT '重量:克(g)' AFTER `sale_price`,
CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '产品编号' ,
CHANGE COLUMN `category_id` `cat_id` INT(11) NULL DEFAULT NULL COMMENT '分类编号' ,
CHANGE COLUMN `supplier_id` `supplier_id` INT(11) NULL DEFAULT NULL COMMENT '供货商编号' ,
ADD COLUMN `brand_id` INT(11) NULL COMMENT '品牌编号' AFTER `supplier_id`,
ADD COLUMN `shop_id` INT(11) NULL COMMENT '商铺编号' AFTER `supplier_id`,
ADD COLUMN `bulk` INT(11) NULL COMMENT '体积:毫升(ml)' AFTER `weight`,
COMMENT = '产品' , RENAME TO  `pro_product` ;



/** ======== new table **/

CREATE TABLE spec_model (
  id      int(10) NOT NULL AUTO_INCREMENT comment '编号',
  name    varchar(10) NOT NULL comment '名称',
  enabled int(1) NOT NULL comment '是否启用',
  PRIMARY KEY (id)) comment='规格模型';
CREATE TABLE o_seo (
  id          int(10) NOT NULL AUTO_INCREMENT comment '编号',
  use_id      int(10) NOT NULL comment '使用者编号',
  use_type    int(10) NOT NULL comment '使用者类型',
  title       varchar(120) comment '标题',
  keywords    varchar(120) comment '关键词',
  description varchar(200) comment '描述',
  PRIMARY KEY (id)) comment='SEO信息表';
CREATE TABLE cat_brand (
  Id       int(10) NOT NULL AUTO_INCREMENT,
  brand_id int(10) NOT NULL,
  cat_id   int(10) NOT NULL,
  PRIMARY KEY (Id)) comment='分类品牌关联';
CREATE TABLE brand (
  id          int(11) NOT NULL AUTO_INCREMENT comment '编号',
  name        varchar(45) NOT NULL comment '品牌名称',
  image       varchar(200) NOT NULL comment '品牌图片',
  site_url    varchar(120) comment '品牌网址',
  intro       varchar(255) comment '介绍',
  review      bit(1) NOT NULL comment '是否审核',
  create_time int(11) comment '加入时间',
  PRIMARY KEY (id)) comment='品牌';


