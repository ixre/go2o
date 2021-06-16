ALTER TABLE public.mch_enterprise_info DROP COLUMN review_state;

ALTER TABLE public.mch_enterprise_info
    ADD COLUMN review_state integer;

COMMENT ON COLUMN public.mch_enterprise_info.review_state
    IS '审核状态';


ALTER TABLE public.mch_shop DROP COLUMN  shop_type;

ALTER TABLE public.mch_shop DROP COLUMN opening_state;


ALTER TABLE public.mch_shop
    ADD COLUMN shop_type integer;

COMMENT ON COLUMN public.mch_shop.shop_type
    IS '店铺类型';

ALTER TABLE public.mch_shop
    ADD COLUMN opening_state integer;

COMMENT ON COLUMN public.mch_shop.opening_state
    IS '营业状态';


ALTER TABLE public.mch_shop
    DROP COLUMN state;

ALTER TABLE public.mch_shop
    ADD COLUMN state int2;

COMMENT ON COLUMN public.mch_shop.state
    IS '状态 1:表示正常,2:表示关闭 ';

ALTER TABLE public.mch_saleconf
    RENAME TO mch_sale_conf;



DROP TABLE mch_online_shop;
DROP TABLE mch_merchant;
TRUNCATE TABLE mch_sale_conf;
TRUNCATE TABLE mch_api_info;
TRUNCATE TABLE mch_account;



CREATE TABLE "public".mch_online_shop (
  id           SERIAL NOT NULL,
  vendor_id   int4 NOT NULL,
  shop_name   varchar(20) NOT NULL,
  logo        varchar(120) NOT NULL,
  host        varchar(40) NOT NULL,
  alias       varchar(20) NOT NULL,
  tel         varchar(45) NOT NULL,
  addr        varchar(120) NOT NULL,
  shop_title  varchar(120) NOT NULL,
  shop_notice varchar(255) NOT NULL,
  flag        int4 NOT NULL,
  state       int2 NOT NULL,
  create_time int8 NOT NULL,
  CONSTRAINT mch_online_shop_pkey
    PRIMARY KEY (id));
COMMENT ON COLUMN "public".mch_online_shop.id IS '店铺编号';
COMMENT ON COLUMN "public".mch_online_shop.vendor_id IS '商户编号';
COMMENT ON COLUMN "public".mch_online_shop.shop_name IS '店铺名称';
COMMENT ON COLUMN "public".mch_online_shop.logo IS '店铺标志';
COMMENT ON COLUMN "public".mch_online_shop.host IS '自定义 域名';
COMMENT ON COLUMN "public".mch_online_shop.alias IS '个性化域名';
COMMENT ON COLUMN "public".mch_online_shop.tel IS '电话';
COMMENT ON COLUMN "public".mch_online_shop.addr IS '地址';
COMMENT ON COLUMN "public".mch_online_shop.shop_title IS '店铺标题';
COMMENT ON COLUMN "public".mch_online_shop.shop_notice IS '店铺公告';
COMMENT ON COLUMN "public".mch_online_shop.flag IS '标志';
COMMENT ON COLUMN "public".mch_online_shop.state IS '状态';
COMMENT ON COLUMN "public".mch_online_shop.create_time IS '创建时间';

CREATE TABLE "public".mch_merchant (
  id              serial NOT NULL,
  member_id       int8 NOT NULL,
  login_user      varchar(20) NOT NULL,
  login_pwd       varchar(45) NOT NULL,
  name            varchar(20) NOT NULL,
  company_name    varchar(45) NOT NULL,
  self_sales      int2 NOT NULL,
  level           int4 NOT NULL,
  logo            varchar(120) NOT NULL,
  province        int4 NOT NULL,
  city            int4 NOT NULL,
  district        int4 NOT NULL,
  create_time     int4 NOT NULL,
  flag            int4 NOT NULL,
  enabled         int2 NOT NULL,
  expires_time    int4 NOT NULL,
  update_time     int4 NOT NULL,
  login_time      int4 NOT NULL,
  last_login_time int4 NOT NULL,
  CONSTRAINT mch_merchant_pkey
    PRIMARY KEY (id));
COMMENT ON TABLE "public".mch_merchant IS '商户';
COMMENT ON COLUMN "public".mch_merchant.member_id IS '会员编号';
COMMENT ON COLUMN "public".mch_merchant.login_user IS '登录用户';
COMMENT ON COLUMN "public".mch_merchant.login_pwd IS '登录密码';
COMMENT ON COLUMN "public".mch_merchant.name IS '名称';
COMMENT ON COLUMN "public".mch_merchant.company_name IS '公司名称';
COMMENT ON COLUMN "public".mch_merchant.self_sales IS '是否自营';
COMMENT ON COLUMN "public".mch_merchant.level IS '商户等级';
COMMENT ON COLUMN "public".mch_merchant.logo IS '标志';
COMMENT ON COLUMN "public".mch_merchant.province IS '省';
COMMENT ON COLUMN "public".mch_merchant.city IS '市';
COMMENT ON COLUMN "public".mch_merchant.district IS '区';
COMMENT ON COLUMN "public".mch_merchant.create_time IS '创建时间';
COMMENT ON COLUMN "public".mch_merchant.flag IS '标志';
COMMENT ON COLUMN "public".mch_merchant.enabled IS '是否启用';
COMMENT ON COLUMN "public".mch_merchant.expires_time IS '过期时间';
COMMENT ON COLUMN "public".mch_merchant.update_time IS '更新时间';
COMMENT ON COLUMN "public".mch_merchant.login_time IS '登录时间';
COMMENT ON COLUMN "public".mch_merchant.last_login_time IS '最后登录时间';





/** 删除无用表 */
DROP TABLE gs_item_tag;
DROP TABLE gs_category;
DROP TABLE gs_sale_snapshot;
DROP TABLE gs_sale_tag;
DROP TABLE gs_snapshot;
DROP TABLE gs_item;
DROP TABLE gs_goods;
DROP TABLE gc_order_confirm;
DROP TABLE gc_member;
DROP TABLE pt_page;
DROP TABLE pt_positions;
DROP TABLE pt_shop;

DROP TABLE pt_saleconf;
DROP TABLE pt_order_log;
DROP TABLE pt_order_item;
DROP TABLE pt_order;
DROP TABLE pt_kvset_member;
DROP TABLE pt_kvset;
DROP TABLE pt_api;
DROP TABLE pt_siteconf;