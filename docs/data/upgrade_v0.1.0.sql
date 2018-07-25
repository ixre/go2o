CREATE TABLE `mm_levelup` (
  `id`           int(11)    NOT NULL AUTO_INCREMENT,
  `member_id`    int(11)    NOT NULL
  COMMENT '会员编号',
  `origin_level` tinyint(2) NOT NULL
  COMMENT '原来等级',
  `target_level` tinyint(2) NOT NULL
  COMMENT '现在等级',
  `is_free`      int(1)     NOT NULL
  COMMENT '是否为免费升级的会员',
  `payment_id`   int(11)    NOT NULL
  COMMENT '支付单编号',
  `reviewed`     int(1)     NOT NULL
  COMMENT '是否审核及处理',
  `create_time`  int(11)    NOT NULL
  COMMENT '升级时间',
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci
  COMMENT ='升级日志表';

CREATE TABLE `comm_qr_template` (
  `id`       int(11)                              NOT NULL AUTO_INCREMENT
  COMMENT '编号',
  `title`    varchar(45) COLLATE utf8_unicode_ci  NOT NULL
  COMMENT '模板标题',
  `bg_image` varchar(120) COLLATE utf8_unicode_ci NOT NULL
  COMMENT '背景图片',
  `offset_x` int(11)                              NOT NULL
  COMMENT '垂直偏离量',
  `offset_y` int(11)                              NOT NULL
  COMMENT '垂直偏移量',
  `comment`  varchar(120) COLLATE utf8_unicode_ci          DEFAULT NULL
  COMMENT '二维码模板文本',
  `enabled`  int(1)                                        DEFAULT NULL
  COMMENT '是否启用',
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci
  COMMENT ='二维码模板';

ALTER TABLE `go2o`.`comm_qr_template`
  ADD COLUMN `callback_url` VARCHAR(120) NULL
COMMENT '回调地址'
  AFTER `comment`;


ALTER TABLE `mm_member`
  CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT
COMMENT '编号',
  CHANGE COLUMN `trade_pwd` `trade_pwd` VARCHAR(45) NULL DEFAULT NULL
COMMENT '交易密码',
  CHANGE COLUMN `exp` `exp` INT(11) UNSIGNED NULL DEFAULT '0'
COMMENT '经验值',
  CHANGE COLUMN `level` `level` INT(11) NULL DEFAULT '1'
COMMENT '等级',
  CHANGE COLUMN `reg_ip` `reg_ip` VARCHAR(20) NULL DEFAULT NULL
COMMENT '注册IP',
  CHANGE COLUMN `reg_from` `reg_from` VARCHAR(20) NULL DEFAULT NULL
COMMENT '注册来源',
  CHANGE COLUMN `reg_time` `reg_time` INT(11) NULL DEFAULT NULL
COMMENT '注册时间',
  CHANGE COLUMN `check_code` `check_code` VARCHAR(8) NULL DEFAULT NULL
COMMENT '校验码',
  CHANGE COLUMN `check_expires` `check_expires` INT(11) NULL DEFAULT NULL
COMMENT '校验码过期时间',
  CHANGE COLUMN `login_time` `login_time` INT(11) NULL DEFAULT NULL
COMMENT '登陆时间',
  CHANGE COLUMN `state` `state` INT(1) NULL DEFAULT '1'
COMMENT '状态',
  CHANGE COLUMN `update_time` `update_time` INT(11) NULL DEFAULT NULL
COMMENT '更新时间';

ALTER TABLE `mm_account`
  CHANGE COLUMN `integral` `integral` INT(11) NOT NULL DEFAULT '0'
COMMENT '积分',
  CHANGE COLUMN `freeze_integral` `freeze_integral` INT(11) NOT NULL
COMMENT '冻结积分',
  CHANGE COLUMN `freeze_balance` `freeze_balance` DECIMAL(10, 2) NOT NULL
COMMENT '冻结账户余额',
  CHANGE COLUMN `expired_balance` `expired_balance` DECIMAL(10, 2) NOT NULL
COMMENT '失效账户余额',
  CHANGE COLUMN `present_balance` `present_balance` DECIMAL(10, 2) NOT NULL DEFAULT '0.00'
COMMENT '赠送余额',
  CHANGE COLUMN `total_present_fee` `total_present_amount` DECIMAL(10, 2) NOT NULL DEFAULT '0.00'
COMMENT '累计赠送金额',
  CHANGE COLUMN `flow_balance` `flow_balance` DECIMAL(10, 2) NOT NULL DEFAULT '0.00'
COMMENT '浮动账户余额',
  CHANGE COLUMN `grow_balance` `grow_balance` DECIMAL(10, 2) NOT NULL
COMMENT '增长账户余额',
  CHANGE COLUMN `grow_amount` `grow_amount` DECIMAL(10, 2) NOT NULL DEFAULT '0.00'
COMMENT '浮动账户余额总投资金额,不含收益',
  CHANGE COLUMN `grow_earnings` `grow_earnings` DECIMAL(10, 2) NOT NULL
COMMENT '当前收益金额',
  CHANGE COLUMN `grow_total_earnings` `grow_total_earnings` DECIMAL(10, 2) NOT NULL
COMMENT '累积收益金额',
  CHANGE COLUMN `total_charge` `total_charge` DECIMAL(10, 2) NOT NULL DEFAULT '0.00'
COMMENT '总充值金额',
  CHANGE COLUMN `total_pay` `total_pay` DECIMAL(10, 2) NOT NULL DEFAULT '0.00'
COMMENT '总支付额',
  CHANGE COLUMN `total_consumption` `total_consumption` DECIMAL(10, 2) NOT NULL
COMMENT '总消费金额',
  CHANGE COLUMN `update_time` `update_time` INT(11) NOT NULL DEFAULT '0'
COMMENT '更新时间';

ALTER TABLE `mm_profile`
  DROP COLUMN `qq`,
  CHANGE COLUMN `email` `email` VARCHAR(50) NULL DEFAULT NULL
COMMENT '电子邮件'
  AFTER `im`,
  CHANGE COLUMN `province` `province` INT(8) NULL DEFAULT NULL
COMMENT '省'
  AFTER `email`,
  CHANGE COLUMN `city` `city` INT(8) NULL DEFAULT NULL
COMMENT '市'
  AFTER `province`,
  CHANGE COLUMN `district` `district` INT(8) NULL DEFAULT NULL
COMMENT '区'
  AFTER `city`,
  CHANGE COLUMN `avatar` `avatar` VARCHAR(80) NULL DEFAULT NULL
COMMENT '头像',
  CHANGE COLUMN `birthday` `birthday` VARCHAR(20) NULL DEFAULT NULL
COMMENT '生日',
  CHANGE COLUMN `phone` `phone` VARCHAR(15) NULL DEFAULT NULL
COMMENT '电话',
  CHANGE COLUMN `address` `address` VARCHAR(100) NULL DEFAULT NULL
COMMENT '地址',
  CHANGE COLUMN `im` `im` VARCHAR(45) NULL DEFAULT NULL
COMMENT '即时通讯',
  CHANGE COLUMN `ext_1` `ext_1` VARCHAR(45) NULL DEFAULT NULL
COMMENT '扩展1',
  CHANGE COLUMN `ext_2` `ext_2` VARCHAR(45) NULL DEFAULT NULL
COMMENT '扩展2',
  CHANGE COLUMN `ext_3` `ext_3` VARCHAR(45) NULL DEFAULT NULL
COMMENT '扩展3',
  CHANGE COLUMN `ext_4` `ext_4` VARCHAR(45) NULL DEFAULT NULL
COMMENT '扩展4',
  CHANGE COLUMN `ext_5` `ext_5` VARCHAR(45) NULL DEFAULT NULL
COMMENT '扩展5',
  CHANGE COLUMN `ext_6` `ext_6` VARCHAR(45) NULL DEFAULT NULL
COMMENT '扩展6',
  CHANGE COLUMN `remark` `remark` VARCHAR(100) NULL DEFAULT NULL
COMMENT '备注',
  CHANGE COLUMN `update_time` `update_time` INT(11) NULL DEFAULT NULL
COMMENT '更新时间',
  COMMENT = '会员资料';


ALTER TABLE `mm_relation`
  CHANGE COLUMN `member_id` `member_id` INT(11) NOT NULL
COMMENT '会员编号',
  CHANGE COLUMN `card_id` `card_no` VARCHAR(20) NULL DEFAULT NULL,
  CHANGE COLUMN `invi_member_id` `inviter_id` INT(11) NOT NULL
COMMENT '邀请人(会员)编号',
  CHANGE COLUMN `refer_str` `inviter_str` VARCHAR(250) NULL DEFAULT NULL
COMMENT '邀请人(会员）字符表示',
  CHANGE COLUMN `reg_merchant_id` `reg_mchid` INT(11) NULL DEFAULT NULL
COMMENT '注册关联的商户编号',
  COMMENT = '会员关系';


ALTER TABLE `gs_category`
  DROP COLUMN `description`,
  DROP COLUMN `mch_id`,
  CHANGE COLUMN `icon` `icon` VARCHAR(150) NULL
COMMENT '分类图片'
  AFTER `name`,
  CHANGE COLUMN `level` `level` INT(3) NOT NULL
COMMENT '分类层级'
  AFTER `url`,
  CHANGE COLUMN `sort_number` `sort_number` INT(11) NOT NULL
COMMENT '排序序号'
  AFTER `level`,
  CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT
COMMENT '编号',
  CHANGE COLUMN `parent_id` `parent_id` INT(11) NOT NULL
COMMENT '父分类',
  CHANGE COLUMN `name` `name` VARCHAR(20) NOT NULL
COMMENT '分类名称',
  CHANGE COLUMN `url` `url` VARCHAR(120) NULL
COMMENT '品牌链接地址',
  CHANGE COLUMN `enabled` `enabled` INT(1) NOT NULL
COMMENT '是否启用',
  CHANGE COLUMN `create_time` `create_time` INT(11) NOT NULL
COMMENT '创建时间',
  ADD COLUMN `spec_model` INT(11) NULL
COMMENT '商品规格模型'
  AFTER `parent_id`,
  COMMENT = '商品分类', RENAME TO `cat_category`;


ALTER TABLE `cat_category`
  CHANGE COLUMN `spec_model` `pro_model` INT(11) NOT NULL DEFAULT 0
COMMENT '产品模型'
  AFTER `parent_id`,
  CHANGE COLUMN `sort_number` `sort_num` INT(11) NOT NULL
COMMENT '排序序号';


/** update **/

ALTER TABLE `gs_item`
  CHANGE COLUMN `express_tid` `express_tid` INT(11) NULL
COMMENT '快递模板编号'
  AFTER `supplier_id`,
  CHANGE COLUMN `goods_no` `goods_no` VARCHAR(45) NULL DEFAULT NULL
  AFTER `small_title`,
  CHANGE COLUMN `weight` `weight` FLOAT(6, 2) NULL DEFAULT NULL
COMMENT '重量:克(g)'
  AFTER `sale_price`,
  CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT
COMMENT '产品编号',
  CHANGE COLUMN `category_id` `cat_id` INT(11) NULL DEFAULT NULL
COMMENT '分类编号',
  CHANGE COLUMN `supplier_id` `supplier_id` INT(11) NULL DEFAULT NULL
COMMENT '供货商编号',
  ADD COLUMN `brand_id` INT(11) NULL
COMMENT '品牌编号'
  AFTER `supplier_id`,
  ADD COLUMN `shop_id` INT(11) NULL
COMMENT '商铺编号'
  AFTER `supplier_id`,
  ADD COLUMN `bulk` INT(11) NULL
COMMENT '体积:毫升(ml)'
  AFTER `weight`,
  COMMENT = '产品', RENAME TO `pro_product`;


ALTER TABLE `pro_product`
  CHANGE COLUMN `goods_no` `code` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '供货商编码';

ALTER TABLE `pro_product`
  ADD COLUMN `sort_num` INT(11) NULL
COMMENT '排序序号'
  AFTER `update_time`;

ALTER TABLE `gs_goods`
  CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT
COMMENT '商品编号',
  CHANGE COLUMN `item_id` `product_id` INT(11) NULL DEFAULT NULL
COMMENT '产品编号',
  CHANGE COLUMN `sku_id` `sku_id` INT(11) NULL DEFAULT NULL
COMMENT '默认SKU编号',
  ADD COLUMN `cat_id` INT(11) NULL
COMMENT '分类编号'
  AFTER `sale_num`,
  ADD COLUMN `vendor_id` INT(11) NULL
COMMENT '供货商编号'
  AFTER `cat_id`,
  ADD COLUMN `brand_id` INT(11) NULL
COMMENT '品牌编号(冗余)\n'
  AFTER `vendor_id`,
  ADD COLUMN `shop_id` INT(11) NULL
COMMENT '商铺编号'
  AFTER `brand_id`,
  ADD COLUMN `shop_cat_id` INT(11) NULL
COMMENT '商铺分类编号'
  AFTER `shop_id`,
  ADD COLUMN `express_tid` INT(11) NULL
COMMENT '快递模板编号'
  AFTER `shop_cat_id`,
  ADD COLUMN `title` VARCHAR(120) NULL
COMMENT '商品标题'
  AFTER `express_tid`,
  ADD COLUMN `code` VARCHAR(45) NULL
COMMENT '供货商编码'
  AFTER `title`,
  ADD COLUMN `image` VARCHAR(120) NULL
COMMENT '主图'
  AFTER `code`,
  ADD COLUMN `cost` INT(11) NULL
COMMENT '成本价'
  AFTER `image`,
  ADD COLUMN `retail_price` INT(11) NULL
COMMENT '零售价'
  AFTER `cost`,
  ADD COLUMN `price` INT(11) NULL
COMMENT '销售价'
  AFTER `retail_price`,
  ADD COLUMN `price_range` VARCHAR(120) NULL
COMMENT '销售价格区间'
  AFTER `price`,
  ADD COLUMN `sku_num` INT(2) NULL
COMMENT 'SKU数量'
  AFTER `price_range`,
  ADD COLUMN `weight` INT(6) NULL
COMMENT '重量:克(g)'
  AFTER `sku_num`,
  ADD COLUMN `bulk` INT(6) NULL
COMMENT '体积:毫升(ml)'
  AFTER `weight`,
  ADD COLUMN `shelve_state` INT(1) NULL
COMMENT '是否上架'
  AFTER `bulk`,
  ADD COLUMN `review_state` INT(1) NULL
COMMENT '审核状态'
  AFTER `shelve_state`,
  ADD COLUMN `review_remark` VARCHAR(120) NULL
COMMENT '审核备注'
  AFTER `review_state`,
  ADD COLUMN `sort_num` INT(11) NULL
COMMENT '排序序号'
  AFTER `review_remark`,
  ADD COLUMN `create_time` INT(11) NULL
COMMENT '创建时间'
  AFTER `sort_num`,
  ADD COLUMN `update_time` INT(11) NULL
COMMENT '更新时间'
  AFTER `create_time`;

ALTER TABLE `gs_goods`
  COMMENT = '商品', RENAME TO `item_info`;

ALTER TABLE `item_info`
  CHANGE COLUMN `is_present` `is_present` INT(1) NULL DEFAULT NULL
COMMENT '是否为赠品\n'
  AFTER `image`,
  CHANGE COLUMN `price_range` `price_range` VARCHAR(120) NULL DEFAULT NULL
COMMENT '销售价格区间'
  AFTER `is_present`,
  CHANGE COLUMN `sku_id` `sku_id` INT(11) NULL DEFAULT NULL
COMMENT '默认SKU编号'
  AFTER `sku_num`,
  CHANGE COLUMN `stock_num` `stock_num` INT(11) NULL DEFAULT NULL
  AFTER `bulk`,
  CHANGE COLUMN `sale_num` `sale_num` INT(11) NULL DEFAULT NULL
  AFTER `stock_num`;

ALTER TABLE `item_info`
  CHANGE COLUMN `stock_num` `stock_num` INT(11) NULL DEFAULT NULL
COMMENT '总库存'
  AFTER `price_range`,
  CHANGE COLUMN `sale_num` `sale_num` INT(11) NULL DEFAULT NULL
COMMENT '销售数量'
  AFTER `stock_num`,
  CHANGE COLUMN `sku_num` `sku_num` INT(2) NULL DEFAULT NULL
COMMENT 'SKU数量'
  AFTER `sale_num`,
  CHANGE COLUMN `sku_id` `sku_id` INT(11) NULL DEFAULT NULL
COMMENT '默认SKU编号'
  AFTER `sku_num`,
  CHANGE COLUMN `price` `price` INT(11) NULL DEFAULT NULL
COMMENT '销售价'
  AFTER `cost`;

ALTER TABLE `item_info`
  CHANGE COLUMN `cost` `cost` DECIMAL(8, 2) NULL DEFAULT NULL
COMMENT '成本价',
  CHANGE COLUMN `price` `price` DECIMAL(8, 2) NULL DEFAULT NULL
COMMENT '销售价',
  CHANGE COLUMN `retail_price` `retail_price` DECIMAL(8, 2) NULL DEFAULT NULL
COMMENT '零售价';

ALTER TABLE `gs_sale_snapshot`
  CHANGE COLUMN `price` `retail_price` DECIMAL(8, 2) NULL DEFAULT '0.00'
COMMENT '售价(市场价)';

ALTER TABLE `pro_product`
  DROP COLUMN `bulk`,
  DROP COLUMN `weight`,
  DROP COLUMN `price`,
  DROP COLUMN `cost`,
  DROP COLUMN `small_title`,
  DROP COLUMN `express_tid`,
  DROP COLUMN `shop_id`;


ALTER TABLE `gs_snapshot`
  CHANGE COLUMN `sku_id` `sku_id` INT(11) NOT NULL
COMMENT '商品快照', RENAME TO `item_snapshot`;

ALTER TABLE `item_snapshot`
  DROP COLUMN `shelve_state`,
  DROP COLUMN `stock_num`,
  DROP COLUMN `sale_num`,
  DROP COLUMN `level_sales`,
  CHANGE COLUMN `vendor_id` `vendor_id` INT(11) NULL DEFAULT NULL
COMMENT '供货商编号',
  ADD COLUMN `brand_id` INT(11) NULL
COMMENT '编号'
  AFTER `vendor_id`,
  ADD COLUMN `shop_id` INT(11) NULL
COMMENT '商铺编号'
  AFTER `brand_id`,
  ADD COLUMN `shop_cat_id` INT(11) NULL
COMMENT '编号分类编号'
  AFTER `shop_id`,
  ADD COLUMN `is_present` INT(1) NULL
COMMENT '是否为赠品'
  AFTER `image`,
  ADD COLUMN `price_range` VARCHAR(20) NULL
COMMENT '价格区间'
  AFTER `is_present`,
  CHANGE COLUMN `item_id` `item_id` INT(11) NOT NULL
COMMENT '商品编号'
  FIRST,
  CHANGE COLUMN `cat_id` `cat_id` INT(11) NULL DEFAULT NULL
COMMENT '分类编号'
  AFTER `snapshot_key`,
  CHANGE COLUMN `express_tid` `express_tid` INT(11) NULL
COMMENT '运费模板'
  AFTER `shop_cat_id`,
  CHANGE COLUMN `sku_id` `sku_id` INT(11) NULL
COMMENT '默认SKU'
  AFTER `price_range`,
  CHANGE COLUMN `weight` `weight` INT(11) NULL DEFAULT NULL
COMMENT '重量(g)'
  AFTER `retail_price`,
  CHANGE COLUMN `snapshot_key` `snapshot_key` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '快照编码',
  CHANGE COLUMN `goods_title` `title` VARCHAR(120) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '商品标题',
  CHANGE COLUMN `small_title` `short_title` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '短标题',
  CHANGE COLUMN `goods_no` `code` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '商户编码',
  CHANGE COLUMN `img` `image` VARCHAR(120) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '商品图片',
  CHANGE COLUMN `cost` `cost` DECIMAL(8, 2) NULL
COMMENT '成本',
  CHANGE COLUMN `price` `price` DECIMAL(8, 2) NULL DEFAULT '0.00'
COMMENT '售价',
  CHANGE COLUMN `sale_price` `retail_price` DECIMAL(8, 2) NULL DEFAULT NULL
COMMENT '零售价',
  CHANGE COLUMN `update_time` `update_time` INT(11) NULL DEFAULT NULL
COMMENT '更新时间',

  DROP PRIMARY KEY,
  ADD PRIMARY KEY (`item_id`);

ALTER TABLE `item_snapshot`
  ADD COLUMN `bulk` INT(11) NULL
COMMENT '体积(ml)'
  AFTER `weight`;

ALTER TABLE `item_snapshot`
  ADD COLUMN `shelve_state` INT(1) NULL
COMMENT '上架状态'
  AFTER `bulk`;

ALTER TABLE `item_snapshot`
  ADD COLUMN `product_id` INT(11) NULL
COMMENT '产品编号'
  AFTER `item_id`;

ALTER TABLE `item_snapshot`
  ADD COLUMN `level_sales` INT(1) NULL
COMMENT '会员价'
  AFTER `bulk`;

CREATE TABLE item_sku (
  id           int(10)        NOT NULL AUTO_INCREMENT
  comment '编号',
  product_id   int(10)        NOT NULL
  comment '产品编号',
  item_id      int(10)        NOT NULL
  comment '商品编号',
  title        varchar(120) comment '标题',
  image        varchar(200) comment '图片',
  spec_data    varchar(200)   NOT NULL
  comment '规格数据',
  spec_word    varchar(200)   NOT NULL
  comment '规格字符',
  code         varchar(45)    NOT NULL
  comment '产品编码',
  retail_price decimal(10, 2) NOT NULL
  comment '参考价',
  price        decimal(10, 2) NOT NULL
  comment '价格（分)',
  cost         decimal(10, 2) NOT NULL
  comment '成本（分)',
  weight       int(11)        NOT NULL
  comment '重量(克)',
  `bulk`       int(11)        NOT NULL
  comment '体积（毫升)',
  stock        int(11)        NOT NULL
  comment '库存',
  sale_num     int(11)        NOT NULL
  comment '已销售数量',
  PRIMARY KEY (id)
);


ALTER TABLE `sale_cart_item`
  CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT
COMMENT '编号',
  CHANGE COLUMN `cart_id` `cart_id` INT(11) NULL DEFAULT NULL
COMMENT '购物车编号',
  CHANGE COLUMN `vendor_id` `vendor_id` INT(11) NULL DEFAULT NULL
COMMENT '运营商编号',
  CHANGE COLUMN `shop_id` `shop_id` INT(11) NULL DEFAULT NULL
COMMENT '店铺编号',
  CHANGE COLUMN `goods_id` `item_id` INT(11) NULL DEFAULT NULL
COMMENT '商品编号',
  CHANGE COLUMN `snap_id` `sku_id` INT(11) NULL DEFAULT NULL
COMMENT 'SKU编号',
  CHANGE COLUMN `quantity` `quantity` INT(8) NULL DEFAULT NULL
COMMENT '数量',
  CHANGE COLUMN `checked` `checked` TINYINT(1) NULL DEFAULT NULL
COMMENT '是否勾选结算',
  COMMENT = '购物车商品项';

ALTER TABLE `sale_order_item`
  ADD COLUMN `item_id` INT(11) NULL
COMMENT '商品编号'
  AFTER `order_id`;


/** ======== new table **/

CREATE TABLE pro_model (
  id       int(10)     NOT NULL AUTO_INCREMENT
  comment '编号',
  name     varchar(10) NOT NULL
  comment '名称',
  attr_str varchar(200) comment '属性字符',
  spec_str varchar(200) comment '规格字符',
  enabled  int(1)      NOT NULL
  comment '是否启用',
  PRIMARY KEY (id)
)
  comment ='产品模型';
CREATE TABLE pro_attr_info (
  id         int(10)      NOT NULL AUTO_INCREMENT
  comment '编号',
  product_id int(10)      NOT NULL
  comment '产品编号',
  attr_id    int(10)      NOT NULL
  comment '属性编号',
  attr_data  varchar(100) NOT NULL
  comment '属性值',
  PRIMARY KEY (id)
)
  comment ='产品属性';
CREATE TABLE pro_attr_item (
  id        int(10)     NOT NULL AUTO_INCREMENT
  comment '编号',
  attr_id   int(10)     NOT NULL
  comment '属性编号',
  pro_model int(10)     NOT NULL
  comment '产品模型',
  value     varchar(20) NOT NULL
  comment '属性值',
  sort_num  int(2)      NOT NULL
  comment '排列序号',
  PRIMARY KEY (id)
)
  comment ='产品属性项';
CREATE TABLE pro_attr (
  id          int(10)     NOT NULL AUTO_INCREMENT
  comment '编号',
  pro_model   int(10)     NOT NULL
  comment '产品模型',
  name        varchar(20) NOT NULL
  comment '属性名称',
  is_filter   int(1)      NOT NULL
  comment '是否作为筛选条件',
  multi_chk   int(1)      NOT NULL
  comment '是否多选',
  item_values varchar(200) comment '属性项值',
  sort_num    int(2)      NOT NULL
  comment '排列序号',
  PRIMARY KEY (id)
)
  comment ='属性';
CREATE TABLE pro_spec (
  id          int(10)     NOT NULL AUTO_INCREMENT
  comment '编号',
  pro_model   int(10)     NOT NULL
  comment '产品模型',
  name        varchar(20) NOT NULL
  comment '规格名称',
  item_values varchar(200) comment '规格项值',
  sort_num    int(2)      NOT NULL
  comment '排列序号',
  PRIMARY KEY (id)
)
  comment ='规格';
CREATE TABLE pro_spec_item (
  id        int(10)     NOT NULL AUTO_INCREMENT
  comment '编号',
  spec_id   int(10)     NOT NULL
  comment '规格编号',
  pro_model int(10)     NOT NULL
  comment '产品模型（冗余)',
  value     varchar(20) NOT NULL
  comment '规格项值',
  color     varchar(20) NOT NULL
  comment '规格项颜色',
  sort_num  int(2)      NOT NULL
  comment '排列序号',
  PRIMARY KEY (id)
)
  comment ='规格项';

CREATE TABLE pro_model_brand (
  id        int(10) NOT NULL AUTO_INCREMENT,
  brand_id  int(10) NOT NULL
  comment '品牌编号',
  pro_model int(10) NOT NULL
  comment '产品模型',
  PRIMARY KEY (id)
)
  comment ='产品模型与品牌关联';
CREATE TABLE pro_brand (
  id          int(11)      NOT NULL AUTO_INCREMENT
  comment '编号',
  name        varchar(45)  NOT NULL
  comment '品牌名称',
  image       varchar(200) NOT NULL
  comment '品牌图片',
  site_url    varchar(120) comment '品牌网址',
  intro       varchar(255) comment '介绍',
  review      int(1)       NOT NULL
  comment '是否审核',
  create_time int(11) comment '加入时间',
  PRIMARY KEY (id)
)
  comment ='产品品牌';


/** 2016-12-30 **/


DROP TABLE `gc_member`, `gc_order_confirm`;
DROP TABLE `gs_category`;
DROP TABLE `sg_bonus`, `sg_bonus_log`, `sg_day_total`, `sg_member`;
DROP TABLE `pt_order`, `pt_order_item`;
DROP TABLE `t_ips`, `t_members`, `t_usrcount`;

CREATE TABLE portal_nav_type (
  id   int(10)     NOT NULL AUTO_INCREMENT
  comment '编号',
  name varchar(20) NOT NULL
  comment '名称',
  PRIMARY KEY (id)
)
  comment ='导航类型';
CREATE TABLE portal_nav (
  id       int(10)      NOT NULL AUTO_INCREMENT
  comment '编号',
  text     varchar(20)  NOT NULL
  comment '文本',
  url      varchar(120) NOT NULL
  comment '地址',
  target   varchar(10)  NOT NULL
  comment '打开目标',
  image    varchar(120) NOT NULL,
  nav_type int(2)       NOT NULL
  comment '导航类型: 1为电脑，2为手机端',
  PRIMARY KEY (id)
)
  comment ='门户导航';

INSERT INTO `portal_nav_type` (`id`, `name`)
VALUES ('1', 'PC商城');
INSERT INTO `portal_nav_type` (`id`, `name`)
VALUES ('2', '移动商城');

ALTER TABLE `pro_attr_info`
  CHANGE COLUMN `product_id` `product_id` INT(10) NOT NULL
COMMENT '产品编号'
  AFTER `id`,
  ADD COLUMN `attr_word` VARCHAR(200) NOT NULL
COMMENT '属性文本'
  AFTER `attr_data`;

/** 2017-01-08 **/

ALTER TABLE `item_info`
  ADD COLUMN `short_title` VARCHAR(120) NULL
COMMENT '短标题'
  AFTER `title`;


/** 2017-01-10 **/
DROP TABLE `pt_merchant`;
DROP TABLE `pt_shop`;
DROP TABLE `pt_api`;

ALTER TABLE `cat_category`
  CHANGE COLUMN `pro_model` `pro_model` INT(11) NOT NULL
COMMENT '商品模型'
  AFTER `parent_id`,
  ADD COLUMN `floor_show` INT(1) NOT NULL
COMMENT '楼层显示'
  AFTER `sort_num`;

CREATE TABLE sys_kv (
  id          int(10)      NOT NULL AUTO_INCREMENT
  comment '编号',
  `key`       varchar(100) NOT NULL
  comment '键',
  value       text         NOT NULL
  comment '值',
  update_time int(10)      NOT NULL
  comment '更新时间',
  PRIMARY KEY (id)
)
  comment ='系统键值';


CREATE TABLE portal_floor_ad (
  id       int(10) NOT NULL AUTO_INCREMENT
  comment '编号',
  cat_id   int(10) NOT NULL
  comment '分类编号',
  pos_id   int(10) NOT NULL
  comment '广告位编号',
  ad_index int(10) NOT NULL
  comment '广告顺序',
  PRIMARY KEY (id)
)
  comment ='楼层广告设置';

CREATE TABLE portal_floor_link (
  id       int(10)      NOT NULL AUTO_INCREMENT
  comment '编号',
  cat_id   int(10)      NOT NULL
  comment '分类编号',
  text     varchar(50)  NOT NULL
  comment '文本',
  link_url varchar(150) NOT NULL
  comment '链接地址',
  target   varchar(10)  NOT NULL
  comment '打开方式',
  sort_num int(10)      NOT NULL
  comment '序号',
  PRIMARY KEY (id)
)
  comment ='楼层链接';


/** 2017-02-17 **/
ALTER TABLE `mm_account`
  CHANGE COLUMN `present_balance` `wallet_balance` DECIMAL(10, 2) NOT NULL DEFAULT '0.00'
COMMENT '钱包余额',
  CHANGE COLUMN `freeze_present` `freeze_wallet` DECIMAL(10, 2) NOT NULL
COMMENT '冻结的钱包金额',
  CHANGE COLUMN `expired_present` `expired_wallet` DECIMAL(10, 2) NOT NULL
COMMENT '过期的钱包金额',
  CHANGE COLUMN `total_present_amount` `total_wallet_amount` DECIMAL(10, 2) NOT NULL DEFAULT '0.00'
COMMENT '累计钱包金额';


ALTER TABLE `mm_present_log`
RENAME TO `mm_wallet_log`;

ALTER TABLE `gs_sales_snapshot`
RENAME TO `item_trade_snapshot`;

DROP TABLE `gs_sale_snapshot`;

ALTER TABLE `sale_after_order`
  ADD COLUMN `image_url` VARCHAR(255) NULL
COMMENT '商品售后图片凭证'
  AFTER `reason`;

/* 2017-02-23 */
DROP TABLE `pt_order_log`;
DROP TABLE `sale_cart`;
DROP TABLE `sale_cart_item`;
CREATE TABLE sale_cart (
  id          int(11)     NOT NULL AUTO_INCREMENT
  comment '编号',
  code        varchar(32) NOT NULL
  comment '购物车编码',
  buyer_id    int(11)     NOT NULL
  comment '买家编号',
  deliver_id  int(11)     NOT NULL
  comment '送货地址',
  payment_opt int(11)     NOT NULL
  comment '支付选项',
  create_time int(11)     NOT NULL
  comment '创建时间',
  update_time int(11)     NOT NULL
  comment '修改时间',
  PRIMARY KEY (id)
)
  comment ='购物车';

CREATE TABLE sale_cart_item (
  id        int(11) NOT NULL AUTO_INCREMENT
  comment '编号',
  cart_id   int(11) NOT NULL
  comment '购物车编号',
  vendor_id int(11) NOT NULL
  comment '运营商编号',
  shop_id   int(11) NOT NULL
  comment '店铺编号',
  item_id   int(11) NOT NULL
  comment '商品编号',
  sku_id    int(11) NOT NULL
  comment 'SKU编号',
  quantity  int(8)  NOT NULL
  comment '数量',
  checked   int(2)  NOT NULL
  comment '是否勾选结算',
  PRIMARY KEY (id)
)
  comment ='购物车商品项';

/* 2017-02-28 */

CREATE TABLE order_list (
  id          int(10)     NOT NULL AUTO_INCREMENT,
  order_no    varchar(45) NOT NULL
  comment '订单号',
  buyer_id    int(10)     NOT NULL
  comment '买家编号',
  order_type  int(2)      NOT NULL
  comment '订单类型',
  create_time int(10)     NOT NULL
  comment '下单时间',
  state       int(2)      NOT NULL
  comment '订单状态',
  PRIMARY KEY (id)
)
  comment ='订单';


/* 订单表更名或者删除后，再重新创建订单 */
CREATE TABLE sale_order (
  id               int(11)       NOT NULL AUTO_INCREMENT
  comment '编号',
  order_id         int(11)       NOT NULL
  comment '订单编号',
  item_amount      decimal(8, 2) NOT NULL
  comment '商品金额',
  discount_amount  decimal(8, 2) NOT NULL
  comment '抵扣金额',
  express_fee      decimal(8, 2) NOT NULL
  comment '物流费',
  package_fee      decimal(4, 2) NOT NULL
  comment '包装费',
  final_amount     decimal(8, 2) NOT NULL
  comment '订单最终金额',
  consignee_person varchar(45)   NOT NULL
  comment '收货人姓名',
  consignee_phone  varchar(45)   NOT NULL
  comment '收货人电话',
  shipping_address varchar(120)  NOT NULL
  comment '收货人地址',
  is_break         int(2)        NOT NULL
  comment '是否拆分',
  update_time      int(11)       NOT NULL
  comment '更新时间',
  PRIMARY KEY (id)
)
  comment ='普通订单';

CREATE TABLE sale_sub_order (
  id              int(11)       NOT NULL AUTO_INCREMENT
  comment '编号',
  order_no        varchar(20)   NOT NULL
  comment '订单号',
  order_id        int(11)       NOT NULL
  comment '订单编号',
  order_pid       int(11)       NOT NULL
  comment '父订单编号',
  buyer_id        int(11)       NOT NULL
  comment '买家编号',
  vendor_id       int(11)       NOT NULL
  comment '商家编号',
  shop_id         int(11)       NOT NULL
  comment '店铺编号',
  subject         varchar(45)   NOT NULL
  comment '主题',
  item_amount     decimal(8, 2) NOT NULL
  comment '商品总价',
  discount_amount decimal(8, 2) NOT NULL
  comment '抵扣金额',
  express_fee     decimal(4, 2) NOT NULL
  comment '运费',
  package_fee     decimal(4, 2) NOT NULL
  comment '包装费',
  final_amount    decimal(8, 2) NOT NULL
  comment '订单最终金额',
  is_paid         int(2)        NOT NULL
  comment '是否支付',
  is_suspend      int(2)        NOT NULL
  comment '是否挂起',
  remark          varchar(120)  NOT NULL
  comment '订单备注',
  buyer_remark    varchar(120)  NOT NULL
  comment '订单买家备注',
  state           int(2)        NOT NULL
  comment '订单状态',
  create_time     int(11)       NOT NULL
  comment '订单创建时间',
  update_time     int(11)       NOT NULL
  comment '订单更新时间',
  PRIMARY KEY (id)
)
  comment ='子订单';

CREATE TABLE sale_order_item (
  id              int(11) NOT NULL AUTO_INCREMENT
  comment '编号',
  order_id        int(11) comment '子订单编号',
  item_id         int(11) comment '商品编号',
  sku_id          int(11) comment 'SKU编号',
  snap_id         int(11) comment '商品快照编号',
  quantity        int(11) comment '销售数量',
  return_quantity int(11) comment '退货数量',
  amount          decimal(8, 2) comment '商品总金额',
  final_amount    decimal(8, 2) comment '商品实际金额',
  is_shipped      bit(1) comment '是否已发货',
  update_time     int(11) comment '更新时间',
  PRIMARY KEY (id)
)
  comment ='普通订单商品项';


/* 2017-03-03 */

ALTER TABLE `sale_sub_order`
  DROP COLUMN `order_pid`;


CREATE TABLE order_wholesale_order (
  id               int(11)       NOT NULL AUTO_INCREMENT
  comment '编号',
  order_no         varchar(20)   NOT NULL
  comment '订单号',
  order_id         int(11)       NOT NULL
  comment '订单编号',
  buyer_id         int(11)       NOT NULL
  comment '买家编号',
  vendor_id        int(11)       NOT NULL
  comment '商家编号',
  shop_id          int(11)       NOT NULL
  comment '店铺编号',
  item_amount      decimal(8, 2) NOT NULL
  comment '商品总价',
  discount_amount  decimal(8, 2) NOT NULL
  comment '抵扣金额',
  express_fee      decimal(4, 2) NOT NULL
  comment '运费',
  package_fee      decimal(4, 2) NOT NULL
  comment '包装费',
  final_amount     decimal(8, 2) NOT NULL
  comment '订单最终金额',
  consignee_person varchar(45)   NOT NULL
  comment '收货人姓名',
  consignee_phone  varchar(45)   NOT NULL
  comment '收货人电话',
  shipping_address varchar(120)  NOT NULL
  comment '收货人地址',
  is_paid          int(2)        NOT NULL
  comment '是否支付',
  remark           varchar(120)  NOT NULL
  comment '订单备注',
  buyer_remark     varchar(120)  NOT NULL
  comment '订单买家备注',
  state            int(2)        NOT NULL
  comment '订单状态',
  create_time      int(11)       NOT NULL
  comment '订单创建时间',
  update_time      int(11)       NOT NULL
  comment '订单更新时间',
  PRIMARY KEY (id)
)
  comment ='批发订单';

CREATE TABLE order_wholesale_item (
  id              int(11)       NOT NULL AUTO_INCREMENT
  comment '编号',
  order_id        int(11)       NOT NULL
  comment '订单编号',
  item_id         int(11)       NOT NULL
  comment '商品编号',
  sku_id          int(11)       NOT NULL
  comment 'SKU编号',
  snapshot_id     int(11)       NOT NULL
  comment '商品快照编号',
  quantity        int(11)       NOT NULL
  comment '销售数量',
  return_quantity int(11)       NOT NULL
  comment '退货数量',
  amount          decimal(8, 2) NOT NULL
  comment '商品总金额',
  final_amount    decimal(8, 2) NOT NULL
  comment '商品实际金额',
  is_shipped      int(2)        NOT NULL
  comment '是否已发货',
  update_time     int(11)       NOT NULL
  comment '更新时间',
  PRIMARY KEY (id)
)
  comment ='批发订单商品';


ALTER TABLE `ship_order`
  CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT
COMMENT '编号',
  CHANGE COLUMN `order_id` `order_id` INT(11) NOT NULL
COMMENT '订单编号',
  CHANGE COLUMN `sp_id` `sp_id` INT(11) NOT NULL
COMMENT '快递SP编号',
  CHANGE COLUMN `sp_order` `sp_order` VARCHAR(20) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NOT NULL
COMMENT '快递SP单号',
  CHANGE COLUMN `exporess_log` `exporess_log` VARCHAR(512) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NOT NULL
COMMENT '物流日志',
  CHANGE COLUMN `amount` `amount` DECIMAL(8, 2) NOT NULL
COMMENT '运费',
  CHANGE COLUMN `final_amount` `final_amount` DECIMAL(8, 2) NOT NULL
COMMENT '实际运费',
  CHANGE COLUMN `ship_time` `ship_time` INT(11) NOT NULL
COMMENT '发货时间',
  CHANGE COLUMN `state` `state` INT(2) NOT NULL
COMMENT '状态',
  CHANGE COLUMN `update_time` `update_time` INT(11) NOT NULL
COMMENT '更新时间',
  ADD COLUMN `sub_orderid` INT(11) NOT NULL
COMMENT '子订单编号'
  AFTER `order_id`,
  COMMENT = '发货单';

ALTER TABLE `ship_order`
  CHANGE COLUMN `exporess_log` `shipment_log` VARCHAR(512) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NOT NULL
COMMENT '物流日志';


ALTER TABLE `ship_item`
  CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT
COMMENT '编号',
  CHANGE COLUMN `ship_order` `ship_order` INT(11) NOT NULL
COMMENT '发货单编号',
  CHANGE COLUMN `snap_id` `snapshot_id` INT(11) NOT NULL
COMMENT '商品交易快照编号',
  CHANGE COLUMN `quantity` `quantity` INT(11) NOT NULL
COMMENT '商品数量',
  CHANGE COLUMN `amount` `amount` DECIMAL(8, 2) NOT NULL
COMMENT '运费',
  CHANGE COLUMN `final_amount` `final_amount` DECIMAL(8, 2) NOT NULL
COMMENT '实际运费',
  COMMENT = '发货单详情';


/* 2017-03-05 */

CREATE TABLE order_trade_order (
  id              int(11)       NOT NULL AUTO_INCREMENT
  comment '编号',
  order_id        int(11)       NOT NULL
  comment '订单编号',
  vendor_id       int(11)       NOT NULL
  comment '商家编号',
  shop_id         int(11)       NOT NULL
  comment '店铺编号',
  subject         varchar(45)   NOT NULL
  comment '订单标题',
  order_amount    decimal(8, 2) NOT NULL
  comment '订单金额',
  discount_amount decimal(8, 2) NOT NULL
  comment '抵扣金额',
  final_amount    decimal(8, 2) NOT NULL
  comment '订单最终金额',
  trade_rate      decimal(8, 2) NOT NULL
  comment '交易结算比例（商户)',
  cash_pay        int(2)        NOT NULL
  comment '是否现金支付',
  remark          varchar(120)  NOT NULL
  comment '订单备注',
  state           int(2)        NOT NULL
  comment '订单状态',
  create_time     int(11)       NOT NULL
  comment '订单创建时间',
  update_time     int(11)       NOT NULL
  comment '订单更新时间',
  PRIMARY KEY (id)
)
  comment ='交易类订单';

ALTER TABLE `mm_account`
  CHANGE COLUMN `total_consumption` `total_expense` DECIMAL(10, 2) NOT NULL
COMMENT '总消费金额';

ALTER TABLE `mm_member`
  CHANGE COLUMN `usr` `usr` VARCHAR(20) NOT NULL
COMMENT '用户名',
  CHANGE COLUMN `pwd` `pwd` VARCHAR(45) NOT NULL
COMMENT '密码',
  CHANGE COLUMN `trade_pwd` `trade_pwd` VARCHAR(45) NOT NULL
COMMENT '交易密码',
  CHANGE COLUMN `exp` `exp` INT(11) UNSIGNED NOT NULL DEFAULT '0'
COMMENT '经验值',
  CHANGE COLUMN `level` `level` INT(11) NOT NULL DEFAULT '1'
COMMENT '等级',
  CHANGE COLUMN `invitation_code` `invitation_code` VARCHAR(10) NOT NULL
COMMENT '邀请码',
  CHANGE COLUMN `reg_ip` `reg_ip` VARCHAR(20) NOT NULL
COMMENT '注册IP',
  CHANGE COLUMN `reg_from` `reg_from` VARCHAR(20) NOT NULL
COMMENT '注册来源',
  CHANGE COLUMN `reg_time` `reg_time` INT(11) NOT NULL
COMMENT '注册时间',
  CHANGE COLUMN `check_code` `check_code` VARCHAR(8) NOT NULL
COMMENT '校验码',
  CHANGE COLUMN `check_expires` `check_expires` INT(11) NOT NULL
COMMENT '校验码过期时间',
  CHANGE COLUMN `login_time` `login_time` INT(11) NOT NULL
COMMENT '登陆时间',
  CHANGE COLUMN `last_login_time` `last_login_time` INT(11) NOT NULL
COMMENT '最后登录时间',
  CHANGE COLUMN `state` `state` INT(1) NOT NULL DEFAULT '1'
COMMENT '状态',
  CHANGE COLUMN `update_time` `update_time` INT(11) NOT NULL
COMMENT '更新时间',
  ADD COLUMN `premium_user` INT(2) NOT NULL
COMMENT '高级用户,0表示非高级'
  AFTER `level`,
  ADD COLUMN `premium_expires` INT(11) NOT NULL
COMMENT '高级会员过期时间'
  AFTER `premium_user`;


/* 2017-03-09 */

ALTER TABLE `order_trade_order`
  ADD COLUMN `ticket_image` VARCHAR(150) NOT NULL
COMMENT '发票图片'
  AFTER `cash_pay`;


CREATE TABLE mm_buyer_group (
  id         int(10)     NOT NULL AUTO_INCREMENT
  comment '编号',
  name       varchar(45) NOT NULL
  comment '名称',
  is_default int(2)      NOT NULL
  comment '是否为默认分组,未设置分组的客户作为该分组。',
  PRIMARY KEY (id)
)
  comment ='买家（客户）分组';


CREATE TABLE mch_buyer_group (
  id               int(10) NOT NULL AUTO_INCREMENT,
  mch_id           int(10) NOT NULL,
  group_id         int(10) NOT NULL,
  alias            int(10) NOT NULL,
  enable_retail    int(2)  NOT NULL
  comment '是否启用零售',
  enable_wholesale int(2)  NOT NULL
  comment '是否启用批发',
  rebate_period    int(2)  NOT NULL
  comment '批发返点周期',
  PRIMARY KEY (id)
);

CREATE TABLE ws_wholesaler (
  mch_id       int(10) NOT NULL AUTO_INCREMENT
  comment '供货商编号等于供货商（等同与商户编号)',
  rate         int(2)  NOT NULL
  comment '批发商评级',
  review_state int(2)  NOT NULL
  comment '批发商审核状态',
  PRIMARY KEY (mch_id)
);

CREATE TABLE ws_rebate_rate (
  id             int(10)       NOT NULL AUTO_INCREMENT,
  wss_id         int(10)       NOT NULL
  comment '批发商编号',
  buyer_gid      int(10)       NOT NULL
  comment '客户分组编号',
  require_amount int(10)       NOT NULL
  comment '下限金额',
  rebate_rate    decimal(6, 4) NOT NULL
  comment '返点率',
  PRIMARY KEY (id)
)
  comment ='批发客户分组返点比例设置';

CREATE TABLE ws_item (
  id            int(10)     NOT NULL AUTO_INCREMENT
  comment '编号',
  vendor_id     int(10)     NOT NULL
  comment '运营商编号',
  item_id       int(10)     NOT NULL
  comment '商品编号',
  shelve_state  int(2)      NOT NULL
  comment '上架状态',
  review_state  int(2)      NOT NULL
  comment '是否审核通过',
  review_remark varchar(45) NOT NULL
  comment '审核备注',
  PRIMARY KEY (id)
)
  comment ='批发商品';


CREATE TABLE ws_sku_price (
  id               int(10)        NOT NULL AUTO_INCREMENT,
  item_id          int(10)        NOT NULL
  comment '商品编号',
  sku_id           int(10)        NOT NULL
  comment 'SKU编号',
  require_quantity int(10)        NOT NULL
  comment '需要数量以上',
  wholesale_price  decimal(10, 2) NOT NULL
  comment '批发价',
  PRIMARY KEY (id)
)
  comment ='商品批发价';

CREATE TABLE ws_item_discount (
  id             int(10)       NOT NULL AUTO_INCREMENT,
  item_id        int(10)       NOT NULL
  comment '商品编号',
  buyer_gid      int(10)       NOT NULL
  comment '客户分组',
  require_amount int(10)       NOT NULL
  comment '要求金额，默认为0',
  discount_rate  decimal(4, 2) NOT NULL
  comment '折扣率',
  PRIMARY KEY (id)
)
  comment ='批发商品折扣';

/* 2017-04-24 */

ALTER TABLE `mch_merchant`
  ADD COLUMN `company_name` VARCHAR(45) NULL
COMMENT '公司名称'
  AFTER `name`;

ALTER TABLE `mch_enterprise_info`
  DROP COLUMN `is_handled`,
  CHANGE COLUMN `address` `address` VARCHAR(120) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '公司地址'
  AFTER `location`,
  CHANGE COLUMN `mch_id` `mch_id` INT(11) NULL DEFAULT NULL
COMMENT '商户编号',
  CHANGE COLUMN `name` `company_name` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '公司名称',
  CHANGE COLUMN `company_no` `company_no` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '营业执照编号',
  CHANGE COLUMN `person_name` `person_name` VARCHAR(10) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '法人姓名',
  CHANGE COLUMN `tel` `tel` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicm.A<QS,KKKKKKKKJUUFASDTHJJ MMGBV
  MMode_ci' NULL DEFAULT NULL
COMMENT '公司电话',
  CHANGE COLUMN `province` `province` INT(11) NOT NULL
COMMENT '所在省',
  CHANGE COLUMN `city` `city` INT(11) NOT NULL
COMMENT '所在市',
  CHANGE COLUMN `district` `district` INT(11) NOT NULL
COMMENT '所在区',
  CHANGE COLUMN `location` `location` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '位置',
  CHANGE COLUMN `person_image` `person_image` VARCHAR(120) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '法人身份证照片',
  CHANGE COLUMN `company_image` `company_image` VARCHAR(120) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '营业执照照片',
  CHANGE COLUMN `review_time` `review_time` INT(11) NULL DEFAULT NULL
COMMENT '审核时间',
  CHANGE COLUMN `remark` `review_remark` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL
COMMENT '审核备注';

ALTER TABLE `ws_item`
  CHANGE COLUMN `vendor_id` `vendor_id` INT(10) NOT NULL
COMMENT '运营商编号'
  AFTER `item_id`,
  ADD COLUMN `price` DECIMAL(10, 2) NULL
COMMENT '价格'
  AFTER `vendor_id`,
  ADD COLUMN `price_range` VARCHAR(45) NULL
COMMENT '价格区间'
  AFTER `price`;

CREATE TABLE ws_cart (
  id          int(11) NOT NULL AUTO_INCREMENT
  comment '编号',
  code        varchar(32) comment '购物车编码',
  buyer_id    int(11) comment '买家编号',
  deliver_id  int(11) comment '送货地址',
  payment_opt int(11) comment '支付选项',
  create_time int(11) comment '创建时间',
  update_time int(11) comment '修改时间',
  PRIMARY KEY (id)
)
  comment ='商品批发购物车';

CREATE TABLE ws_cart_item (
  id        int(11) NOT NULL AUTO_INCREMENT
  comment '编号',
  cart_id   int(11) comment '购物车编号',
  vendor_id int(11) comment '运营商编号',
  shop_id   int(11) comment '店铺编号',
  item_id   int(11) comment '商品编号',
  sku_id    int(11) comment 'SKU编号',
  quantity  int(8) comment '数量',
  checked   int(2) comment '是否勾选结算',
  PRIMARY KEY (id)
)
  comment ='批发购物车商品项';


/* 2017-06-09 */

ALTER TABLE `con_article_category`
  CHANGE COLUMN `alias` `cat_alias` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT NULL,
RENAME TO `article_category`;

ALTER TABLE `con_article`
RENAME TO `article_list`;

ALTER TABLE `con_page`
RENAME TO `content_page`;

/* 2017-07-04 */


ALTER TABLE `mch_shop`
  CHANGE COLUMN `id` `id` INT(11) NOT NULL AUTO_INCREMENT
COMMENT '商店编号',
  CHANGE COLUMN `mch_id` `vendor_id` INT(11) NULL DEFAULT NULL
COMMENT '商户编号',
  CHANGE COLUMN `shop_type` `shop_type` TINYINT(1) NULL DEFAULT NULL
COMMENT '商店类型',
  CHANGE COLUMN `name` `name` VARCHAR(50) NULL DEFAULT NULL
COMMENT '商店名称',
  CHANGE COLUMN `sort_number` `sort_number` INT(11) NULL DEFAULT '0'
COMMENT '排序序号',
  CHANGE COLUMN `state` `state` INT(2) NULL DEFAULT NULL
COMMENT '状态 1:表示正常,2:表示关闭 ',
  ADD COLUMN `opening_state` INT(2) NULL
COMMENT '商店营业状态,1:正常,2:暂停营业'
  AFTER `create_time`;


ALTER TABLE `article_category`
  ADD COLUMN `perm_flag` INT(2) NULL
COMMENT '访问权限'
  AFTER `parent_id`;


ALTER TABLE `article_list`
  CHANGE COLUMN `publisher_id` `publisher_id` INT(11) NULL DEFAULT NULL
  AFTER `thumbnail`,
  ADD COLUMN `priority` INT(2) NULL
COMMENT '优先级'
  AFTER `location`,
  ADD COLUMN `access_key` VARCHAR(45) NULL
COMMENT '访问钥匙'
  AFTER `priority`;

ALTER TABLE `content_page`
  CHANGE COLUMN `title` `title` VARCHAR(100) NULL DEFAULT NULL
COMMENT '标题',
  ADD COLUMN `perm_flag` INT(2) NULL
COMMENT '访问权限'
  AFTER `title`,
  ADD COLUMN `access_key` VARCHAR(45) NULL
COMMENT '访问钥匙'
  AFTER `perm_flag`, RENAME TO `ex_page`;

ALTER TABLE `cat_category`
  ADD COLUMN `priority` INT(2) NULL
COMMENT '优先级'
  AFTER `pro_model`;

ALTER TABLE `cat_category`
RENAME TO `pro_category`;

UPDATE pro_category
SET priority = 0
WHERE id > 0
  AND priority IS NULL;

ALTER TABLE `ws_item`
  CHANGE COLUMN `price` `price` DECIMAL(10, 2) NULL DEFAULT 0
COMMENT '价格',
  CHANGE COLUMN `price_range` `price_range` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL DEFAULT 0
COMMENT '价格区间';

ALTER TABLE `ws_item`
  CHANGE COLUMN `price` `price` DECIMAL(10, 2) NOT NULL DEFAULT '0.00'
COMMENT '价格',
  CHANGE COLUMN `review_remark` `review_remark` VARCHAR(45) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NULL
COMMENT '审核备注';

/* 2017-07-07 */

update pro_category
set icon = ''
WHERE id > 0
  AND icon IS NULL;
ALTER TABLE `pro_category`
  ADD COLUMN `virtual_cat` INT(2) NOT NULL DEFAULT 0
  AFTER `name`,
  CHANGE COLUMN `url` `cat_url` VARCHAR(120) NOT NULL
COMMENT '品牌链接地址'
  AFTER `virtual_cat`,
  CHANGE COLUMN `icon` `icon` VARCHAR(150) NOT NULL
COMMENT '分类图片';

ALTER TABLE `mch_shop`
  CHANGE COLUMN `shop_type` `shop_type` TINYINT(1) NOT NULL,
  CHANGE COLUMN `state` `state` TINYINT(1) NOT NULL,
  CHANGE COLUMN `opening_state` `opening_state` TINYINT(1) NOT NULL;

ALTER TABLE `pro_category`
  ADD COLUMN `icon_xy` VARCHAR(45) NOT NULL
  AFTER `icon`;

update pro_category
set icon_xy = '0,0'
WHERE id > 0 && icon_xy IS NULL;


/* 2017-07-15 */

ALTER TABLE `mch_buyer_group`
  CHANGE COLUMN `alias` `alias` VARCHAR(45) NOT NULL,
  CHANGE COLUMN `enable_retail` `enable_retail` TINYINT(1) NOT NULL
COMMENT '是否启用零售',
  CHANGE COLUMN `enable_wholesale` `enable_wholesale` TINYINT(1) NOT NULL
COMMENT '是否启用批发';

ALTER TABLE `mch_online_shop`
  DROP COLUMN `sub_tit`,
  CHANGE COLUMN `index_tit` `shop_title` VARCHAR(120) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NOT NULL,
  CHANGE COLUMN `notice_html` `shop_notice` VARCHAR(255) CHARACTER SET 'utf8'
COLLATE 'utf8_unicode_ci' NOT NULL;

ALTER TABLE `order_wholesale_order`
  CHANGE COLUMN `buyer_remark` `buyer_comment` VARCHAR(120) NOT NULL
COMMENT '订单买家备注'
  AFTER `is_paid`;

ALTER TABLE `sale_sub_order`
  CHANGE COLUMN `buyer_remark` `buyer_comment` VARCHAR(120) NOT NULL
COMMENT '订单买家备注'
  AFTER `is_suspend`;

CREATE TABLE wal_wallet (
  id              int(11)              NOT NULL AUTO_INCREMENT
  comment '钱包编号',
  hash_code       varchar(40)          NOT NULL
  comment '哈希值',
  node_id         int(2)               NOT NULL
  comment '节点编号',
  user_id         int(11)              NOT NULL
  comment '用户编号',
  wallet_type     int(1)               NOT NULL
  comment '钱包类型',
  wallet_flag     int(4)               NOT NULL
  comment '钱包标志',
  balance         int(11) DEFAULT 0.00 NOT NULL
  comment '余额',
  present_balance int(11)              NOT NULL
  comment '赠送余额',
  adjust_amount   int(11)              NOT NULL
  comment '调整金额',
  freeze_amount   int(11)              NOT NULL
  comment '冻结余额',
  latest_amount   int(11)              NOT NULL
  comment '结余金额',
  expired_amount  int(11)              NOT NULL
  comment '失效账户余额',
  total_charge    int(11) DEFAULT 0.00 NOT NULL
  comment '总充值金额',
  total_present   int(11)              NOT NULL
  comment '累计赠送金额',
  total_pay       int(11) DEFAULT 0.00 NOT NULL
  comment '总支付额',
  remark          varchar(40)          NOT NULL
  comment '备注',
  state           int(1)               NOT NULL
  comment '状态,1:正常 2:锁定 3:关停',
  create_time     int(11)              NOT NULL
  comment '创建时间',
  update_time     int(11) DEFAULT 0    NOT NULL
  comment '更新时间',
  PRIMARY KEY (id)
)
  comment ='钱包';

CREATE TABLE wal_wallet_log (
  id            int(11)      NOT NULL AUTO_INCREMENT
  comment '编号',
  wallet_id     int(11)      NOT NULL
  comment '钱包编号',
  kind          int(1)       NOT NULL
  comment '业务类型',
  title         varchar(45)  NOT NULL
  comment '标题',
  outer_chan    varchar(20)  NOT NULL
  comment '外部通道',
  outer_no      varchar(45)  NOT NULL
  comment '外部订单号',
  value         int(11)      NOT NULL
  comment '变动金额',
  balance       int(11)      NOT NULL
  comment '余额',
  trade_fee     int(8)       NOT NULL
  comment '交易手续费',
  op_uid        int(10)      NOT NULL
  comment '操作人员用户编号',
  op_name       varchar(20)  NOT NULL
  comment '操作人员名称',
  remark        varchar(45)  NOT NULL
  comment '备注',
  review_state  int(1)       NOT NULL
  comment '审核状态',
  review_remark varchar(120) NOT NULL
  comment '审核备注',
  review_time   int(11)      NOT NULL
  comment '审核时间',
  create_time   int(11)      NOT NULL
  comment '创建时间',
  update_time   int(11)      NOT NULL
  comment '更新时间',
  PRIMARY KEY (id)
)
  comment ='钱包日志';
