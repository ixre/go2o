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


ALTER TABLE `pro_product`
CHANGE COLUMN `goods_no` `code` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL COMMENT '供货商编码' ;

ALTER TABLE `pro_product`
ADD COLUMN `sort_num` INT(11) NULL COMMENT '排序序号' AFTER `update_time`;

ALTER TABLE `gs_goods` 
CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '商品编号' ,
CHANGE COLUMN `item_id` `product_id` INT(11) NULL DEFAULT NULL COMMENT '产品编号' ,
CHANGE COLUMN `sku_id` `sku_id` INT(11) NULL DEFAULT NULL COMMENT '默认SKU编号' ,
ADD COLUMN `cat_id` INT(11) NULL COMMENT '分类编号' AFTER `sale_num`,
ADD COLUMN `vendor_id` INT(11) NULL COMMENT '供货商编号' AFTER `cat_id`,
ADD COLUMN `brand_id` INT(11) NULL COMMENT '品牌编号(冗余)\n' AFTER `vendor_id`,
ADD COLUMN `shop_id` INT(11) NULL COMMENT '商铺编号' AFTER `brand_id`,
ADD COLUMN `shop_cat_id` INT(11) NULL COMMENT '商铺分类编号' AFTER `shop_id`,
ADD COLUMN `express_tid` INT(11) NULL COMMENT '快递模板编号' AFTER `shop_cat_id`,
ADD COLUMN `title` VARCHAR(120) NULL COMMENT '商品标题' AFTER `express_tid`,
ADD COLUMN `code` VARCHAR(45) NULL COMMENT '供货商编码' AFTER `title`,
ADD COLUMN `image` VARCHAR(120) NULL COMMENT '主图' AFTER `code`,
ADD COLUMN `cost` INT(11) NULL COMMENT '成本价' AFTER `image`,
ADD COLUMN `retail_price` INT(11) NULL COMMENT '零售价' AFTER `cost`,
ADD COLUMN `price` INT(11) NULL COMMENT '销售价' AFTER `retail_price`,
ADD COLUMN `price_range` VARCHAR(120) NULL COMMENT '销售价格区间' AFTER `price`,
ADD COLUMN `sku_num` INT(2) NULL COMMENT 'SKU数量' AFTER `price_range`,
ADD COLUMN `weight` INT(6) NULL COMMENT '重量:克(g)' AFTER `sku_num`,
ADD COLUMN `bulk` INT(6) NULL COMMENT '体积:毫升(ml)' AFTER `weight`,
ADD COLUMN `shelve_state` INT(1) NULL COMMENT '是否上架' AFTER `bulk`,
ADD COLUMN `review_state` INT(1) NULL COMMENT '审核状态' AFTER `shelve_state`,
ADD COLUMN `review_remark` VARCHAR(120) NULL COMMENT '审核备注' AFTER `review_state`,
ADD COLUMN `sort_num` INT(11) NULL COMMENT '排序序号' AFTER `review_remark`,
ADD COLUMN `create_time` INT(11) NULL COMMENT '创建时间' AFTER `sort_num`,
ADD COLUMN `update_time` INT(11) NULL COMMENT '更新时间' AFTER `create_time`;

ALTER TABLE `gs_goods`
COMMENT = '商品' , RENAME TO  `item_info` ;

ALTER TABLE `item_info`
CHANGE COLUMN `is_present` `is_present` INT(1) NULL DEFAULT NULL COMMENT '是否为赠品\n' AFTER `image`,
CHANGE COLUMN `price_range` `price_range` VARCHAR(120) NULL DEFAULT NULL COMMENT '销售价格区间' AFTER `is_present`,
CHANGE COLUMN `sku_id` `sku_id` INT(11) NULL DEFAULT NULL COMMENT '默认SKU编号' AFTER `sku_num`,
CHANGE COLUMN `stock_num` `stock_num` INT(11) NULL DEFAULT NULL AFTER `bulk`,
CHANGE COLUMN `sale_num` `sale_num` INT(11) NULL DEFAULT NULL AFTER `stock_num`;

ALTER TABLE `item_info`
CHANGE COLUMN `stock_num` `stock_num` INT(11) NULL DEFAULT NULL COMMENT '总库存' AFTER `price_range`,
CHANGE COLUMN `sale_num` `sale_num` INT(11) NULL DEFAULT NULL COMMENT '销售数量' AFTER `stock_num`,
CHANGE COLUMN `sku_num` `sku_num` INT(2) NULL DEFAULT NULL COMMENT 'SKU数量' AFTER `sale_num`,
CHANGE COLUMN `sku_id` `sku_id` INT(11) NULL DEFAULT NULL COMMENT '默认SKU编号' AFTER `sku_num`,
CHANGE COLUMN `price` `price` INT(11) NULL DEFAULT NULL COMMENT '销售价' AFTER `cost`;

ALTER TABLE `item_info`
CHANGE COLUMN `cost` `cost` DECIMAL(8,2) NULL DEFAULT NULL COMMENT '成本价' ,
CHANGE COLUMN `price` `price` DECIMAL(8,2) NULL DEFAULT NULL COMMENT '销售价' ,
CHANGE COLUMN `retail_price` `retail_price` DECIMAL(8,2) NULL DEFAULT NULL COMMENT '零售价' ;

ALTER TABLE `gs_sale_snapshot`
CHANGE COLUMN `price` `retail_price` DECIMAL(8,2) NULL DEFAULT '0.00' COMMENT '售价(市场价)' ;

ALTER TABLE `pro_product`
DROP COLUMN `bulk`,
DROP COLUMN `weight`,
DROP COLUMN `price`,
DROP COLUMN `cost`,
DROP COLUMN `small_title`,
DROP COLUMN `express_tid`,
DROP COLUMN `shop_id`;


ALTER TABLE `gs_snapshot`
CHANGE COLUMN `sku_id` `sku_id` INT(11) NOT NULL COMMENT '商品快照' , RENAME TO  `item_snapshot` ;

ALTER TABLE `item_snapshot`
DROP COLUMN `shelve_state`,
DROP COLUMN `stock_num`,
DROP COLUMN `sale_num`,
DROP COLUMN `level_sales`,
CHANGE COLUMN `vendor_id` `vendor_id` INT(11) NULL DEFAULT NULL COMMENT '供货商编号' ,
ADD COLUMN `brand_id` INT(11) NULL COMMENT '编号' AFTER `vendor_id`,
ADD COLUMN `shop_id` INT(11) NULL COMMENT '商铺编号' AFTER `brand_id`,
ADD COLUMN `shop_cat_id` INT(11) NULL COMMENT '编号分类编号' AFTER `shop_id`,
ADD COLUMN `is_present` INT(1) NULL COMMENT '是否为赠品' AFTER `image`,
ADD COLUMN `price_range` VARCHAR(20) NULL COMMENT '价格区间' AFTER `is_present`,
CHANGE COLUMN `item_id` `item_id` INT(11) NOT NULL COMMENT '商品编号' FIRST,
CHANGE COLUMN `cat_id` `cat_id` INT(11) NULL DEFAULT NULL COMMENT '分类编号' AFTER `snapshot_key`,
CHANGE COLUMN `express_tid` `express_tid` INT(11) NULL COMMENT '运费模板' AFTER `shop_cat_id`,
CHANGE COLUMN `sku_id` `sku_id` INT(11) NULL COMMENT '默认SKU' AFTER `price_range`,
CHANGE COLUMN `weight` `weight` INT(11) NULL DEFAULT NULL COMMENT '重量(g)' AFTER `retail_price`,
CHANGE COLUMN `snapshot_key` `snapshot_key` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL COMMENT '快照编码' ,
CHANGE COLUMN `goods_title` `title` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL COMMENT '商品标题' ,
CHANGE COLUMN `small_title` `short_title` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL COMMENT '短标题' ,
CHANGE COLUMN `goods_no` `code` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL COMMENT '商户编码' ,
CHANGE COLUMN `img` `image` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL COMMENT '商品图片' ,
CHANGE COLUMN `cost` `cost` DECIMAL(8,2) NULL COMMENT '成本' ,
CHANGE COLUMN `price` `price` DECIMAL(8,2) NULL DEFAULT '0.00' COMMENT '售价' ,
CHANGE COLUMN `sale_price` `retail_price` DECIMAL(8,2) NULL DEFAULT NULL COMMENT '零售价' ,
CHANGE COLUMN `update_time` `update_time` INT(11) NULL DEFAULT NULL COMMENT '更新时间' ,

DROP PRIMARY KEY,
ADD PRIMARY KEY (`item_id`);

ALTER TABLE `item_snapshot`
ADD COLUMN `bulk` INT(11) NULL COMMENT '体积(ml)' AFTER `weight`;

ALTER TABLE `item_snapshot`
ADD COLUMN `shelve_state` INT(1) NULL COMMENT '上架状态' AFTER `bulk`;

ALTER TABLE `item_snapshot`
ADD COLUMN `product_id` INT(11) NULL COMMENT '产品编号' AFTER `item_id`;

ALTER TABLE `item_snapshot`
ADD COLUMN `level_sales` INT(1) NULL COMMENT '会员价' AFTER `bulk`;

CREATE TABLE item_sku (
  id           int(10) NOT NULL AUTO_INCREMENT comment '编号',
  product_id   int(10) NOT NULL comment '产品编号',
  item_id      int(10) NOT NULL comment '商品编号',
  title        varchar(120) comment '标题',
  image        varchar(200) comment '图片',
  spec_data    varchar(200) NOT NULL comment '规格数据',
  spec_word    varchar(200) NOT NULL comment '规格字符',
  code         varchar(45) NOT NULL comment '产品编码',
  retail_price decimal(10, 2) NOT NULL comment '参考价',
  price        decimal(10, 2) NOT NULL comment '价格（分)',
  cost         decimal(10, 2) NOT NULL comment '成本（分)',
  weight       int(11) NOT NULL comment '重量(克)',
  `bulk`       int(11) NOT NULL comment '体积（毫升)',
  stock        int(11) NOT NULL comment '库存',
  sale_num     int(11) NOT NULL comment '已销售数量',
  PRIMARY KEY (id));


ALTER TABLE `sale_cart_item`
CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '编号' ,
CHANGE COLUMN `cart_id` `cart_id` INT(11) NULL DEFAULT NULL COMMENT '购物车编号' ,
CHANGE COLUMN `vendor_id` `vendor_id` INT(11) NULL DEFAULT NULL COMMENT '运营商编号' ,
CHANGE COLUMN `shop_id` `shop_id` INT(11) NULL DEFAULT NULL COMMENT '店铺编号' ,
CHANGE COLUMN `goods_id` `item_id` INT(11) NULL DEFAULT NULL COMMENT '商品编号' ,
CHANGE COLUMN `snap_id` `sku_id` INT(11) NULL DEFAULT NULL COMMENT 'SKU编号' ,
CHANGE COLUMN `quantity` `quantity` INT(8) NULL DEFAULT NULL COMMENT '数量' ,
CHANGE COLUMN `checked` `checked` TINYINT(1) NULL DEFAULT NULL COMMENT '是否勾选结算' ,
COMMENT = '购物车商品项' ;

ALTER TABLE `sale_order_item`
ADD COLUMN `item_id` INT(11) NULL COMMENT '商品编号' AFTER `order_id`;


/** ======== new table **/

CREATE TABLE pro_model (
  id       int(10) NOT NULL AUTO_INCREMENT comment '编号',
  name     varchar(10) NOT NULL comment '名称',
  attr_str varchar(200) comment '属性字符',
  spec_str varchar(200) comment '规格字符',
  enabled  int(1) NOT NULL comment '是否启用',
  PRIMARY KEY (id)) comment='产品模型';
CREATE TABLE pro_attr_info (
  id         int(10) NOT NULL AUTO_INCREMENT comment '编号',
  product_id int(10) NOT NULL comment '产品编号',
  attr_id    int(10) NOT NULL comment '属性编号',
  attr_data  varchar(100) NOT NULL comment '属性值',
  PRIMARY KEY (id)) comment='产品属性';
CREATE TABLE pro_attr_item (
  id        int(10) NOT NULL AUTO_INCREMENT comment '编号',
  attr_id   int(10) NOT NULL comment '属性编号',
  pro_model int(10) NOT NULL comment '产品模型',
  value     varchar(20) NOT NULL comment '属性值',
  sort_num  int(2) NOT NULL comment '排列序号',
  PRIMARY KEY (id)) comment='产品属性项';
CREATE TABLE pro_attr (
  id          int(10) NOT NULL AUTO_INCREMENT comment '编号',
  pro_model   int(10) NOT NULL comment '产品模型',
  name        varchar(20) NOT NULL comment '属性名称',
  is_filter   int(1) NOT NULL comment '是否作为筛选条件',
  multi_chk   int(1) NOT NULL comment '是否多选',
  item_values varchar(200) comment '属性项值',
  sort_num    int(2) NOT NULL comment '排列序号',
  PRIMARY KEY (id)) comment='属性';
CREATE TABLE pro_spec (
  id          int(10) NOT NULL AUTO_INCREMENT comment '编号',
  pro_model   int(10) NOT NULL comment '产品模型',
  name        varchar(20) NOT NULL comment '规格名称',
  item_values varchar(200) comment '规格项值',
  sort_num    int(2) NOT NULL comment '排列序号',
  PRIMARY KEY (id)) comment='规格';
CREATE TABLE pro_spec_item (
  id        int(10) NOT NULL AUTO_INCREMENT comment '编号',
  spec_id   int(10) NOT NULL comment '规格编号',
  pro_model int(10) NOT NULL comment '产品模型（冗余)',
  value     varchar(20) NOT NULL comment '规格项值',
  color     varchar(20) NOT NULL comment '规格项颜色',
  sort_num  int(2) NOT NULL comment '排列序号',
  PRIMARY KEY (id)) comment='规格项';

CREATE TABLE pro_model_brand (
  id        int(10) NOT NULL AUTO_INCREMENT,
  brand_id  int(10) NOT NULL comment '品牌编号',
  pro_model int(10) NOT NULL comment '产品模型',
  PRIMARY KEY (id)) comment='产品模型与品牌关联';
CREATE TABLE pro_brand (
  id          int(11) NOT NULL AUTO_INCREMENT comment '编号',
  name        varchar(45) NOT NULL comment '品牌名称',
  image       varchar(200) NOT NULL comment '品牌图片',
  site_url    varchar(120) comment '品牌网址',
  intro       varchar(255) comment '介绍',
  review      int(1) NOT NULL comment '是否审核',
  create_time int(11) comment '加入时间',
  PRIMARY KEY (id)) comment='产品品牌';




/** 2016-12-30 **/


DROP TABLE `gc_member`, `gc_order_confirm`;
DROP TABLE `gs_category`;
DROP TABLE `sg_bonus`, `sg_bonus_log`, `sg_day_total`, `sg_member`;
DROP TABLE `pt_order`, `pt_order_item`;
DROP TABLE `t_ips`, `t_members`, `t_usrcount`;

CREATE TABLE portal_nav_type (
  id   int(10) NOT NULL AUTO_INCREMENT comment '编号',
  name varchar(20) NOT NULL comment '名称',
  PRIMARY KEY (id)) comment='导航类型';
CREATE TABLE portal_nav (
  id       int(10) NOT NULL AUTO_INCREMENT comment '编号',
  text     varchar(20) NOT NULL comment '文本',
  url      varchar(120) NOT NULL comment '地址',
  target   varchar(10) NOT NULL comment '打开目标',
  image    varchar(120) NOT NULL,
  nav_type int(2) NOT NULL comment '导航类型: 1为电脑，2为手机端',
  PRIMARY KEY (id)) comment='门户导航';

INSERT INTO `portal_nav_type` (`id`, `name`) VALUES ('1', 'PC商城');
INSERT INTO `portal_nav_type` (`id`, `name`) VALUES ('2', '移动商城');

ALTER TABLE `pro_attr_info`
CHANGE COLUMN `product_id` `product_id` INT(10) NOT NULL COMMENT '产品编号' AFTER `id`,
ADD COLUMN `attr_word` VARCHAR(200) NOT NULL COMMENT '属性文本' AFTER `attr_data`;

/** 2017-01-08 **/

ALTER TABLE `txmall`.`item_info`
ADD COLUMN `short_title` VARCHAR(120) NULL COMMENT '短标题' AFTER `title`;