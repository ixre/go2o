/** 2023-01-02 订单返利 */
CREATE TABLE order_rebate_list (
  id              BIGSERIAL NOT NULL, 
  plan_id        int4 NOT NULL, 
  trader_id      int8 NOT NULL, 
  affiliate_code varchar(20) NOT NULL, 
  order_no       varchar(20) NOT NULL, 
  order_subject  varchar(40) NOT NULL, 
  order_amount   int8 NOT NULL, 
  rebase_amount  int8 NOT NULL, 
  status         int4 NOT NULL, 
  create_time    int8 NOT NULL, 
  update_time    int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE order_rebate_list IS '订单返利';
COMMENT ON COLUMN order_rebate_list.id IS '编号';
COMMENT ON COLUMN order_rebate_list.plan_id IS '返利方案Id';
COMMENT ON COLUMN order_rebate_list.trader_id IS '成交人Id';
COMMENT ON COLUMN order_rebate_list.affiliate_code IS '分享码';
COMMENT ON COLUMN order_rebate_list.order_no IS '订单号';
COMMENT ON COLUMN order_rebate_list.order_subject IS '订单标题';
COMMENT ON COLUMN order_rebate_list.order_amount IS '订单金额';
COMMENT ON COLUMN order_rebate_list.rebase_amount IS '返利金额';
COMMENT ON COLUMN order_rebate_list.status IS '返利状态';
COMMENT ON COLUMN order_rebate_list.create_time IS '创建时间';
COMMENT ON COLUMN order_rebate_list.update_time IS '更新时间';


CREATE TABLE order_rebate_item (
  id             SERIAL NOT NULL, 
  debate_id     int8 NOT NULL, 
  item_id       int8 NOT NULL, 
  item_name     varchar(20) NOT NULL, 
  item_image    varchar(120) NOT NULL, 
  item_amount   int8 NOT NULL, 
  rebate_amount int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE order_rebate_item IS '订单返利详情';
COMMENT ON COLUMN order_rebate_item.id IS '编号';
COMMENT ON COLUMN order_rebate_item.debate_id IS '返利单Id';
COMMENT ON COLUMN order_rebate_item.item_id IS '商品编号';
COMMENT ON COLUMN order_rebate_item.item_name IS '商品名称';
COMMENT ON COLUMN order_rebate_item.item_image IS '商品图片';
COMMENT ON COLUMN order_rebate_item.item_amount IS '商品金额';
COMMENT ON COLUMN order_rebate_item.rebate_amount IS '返利金额';

/* 2023-01-03 19:01 */
ALTER TABLE IF EXISTS public.wal_wallet_log
    RENAME value TO change_value;

ALTER TABLE IF EXISTS public.wal_wallet_log
    RENAME review_state TO audit_state;

ALTER TABLE IF EXISTS public.wal_wallet_log
    RENAME review_remark TO audit_remark;

ALTER TABLE IF EXISTS public.wal_wallet_log
    RENAME review_time TO audit_time;

ALTER TABLE IF EXISTS public.item_info
    RENAME prom_flag TO item_flag;

ALTER TABLE IF EXISTS public.item_info
    RENAME review_state TO audit_state;
ALTER TABLE IF EXISTS public.item_info
    RENAME review_remark TO audit_remark;

ALTER TABLE IF EXISTS public.mm_member
   DROP invite_code;

ALTER TABLE IF EXISTS public.mm_member
    RENAME "user" TO username;

ALTER TABLE IF EXISTS public.mm_member
    RENAME pwd TO password;
ALTER TABLE IF EXISTS public.mm_member
    RENAME avatar TO portrait;
ALTER TABLE IF EXISTS public.mm_member
    RENAME code TO user_code;
ALTER TABLE IF EXISTS public.mm_member
    RENAME flag TO user_flag;
ALTER TABLE IF EXISTS public.mm_member
    RENAME nick_name TO nickname;
-- 更新默认头像地址
update mm_member set portrait = portrait = 'static/init/avatar.png' where portrait like 'init/%'; 

ALTER TABLE IF EXISTS public.wal_wallet_log
    RENAME "title" TO subject;
 -- 删除店铺表
DROP TABLE IF EXISTS public.mch_shop;
-- 重置商品标志
update public.item_info set item_flag = 0 WHERE item_flag < 0;

update item_info set item_flag  = 199 
WHERE id IN (select id from item_info 
			 where item_flag = 0  order by id desc limit 30);

-- 以下VPP需要更新 --
-- 去掉赠品字段改用flag
ALTER TABLE IF EXISTS public.item_info
   DROP is_present;
   ALTER TABLE IF EXISTS public.item_snapshot
   DROP is_present;

/* 2023-01-10 */
update registry set key='sms_api_2' where key='sms_api_http';
delete FROM registry where key in (
'order_push_affiliate_enabled',
'order_push_affiliate_event',
'sms_push_send_event',
'order_push_sub_order_enabled',
'order_enable_affiliate_rebate',
'orde_affiliater_push_enabled',
'orde_affiliate_push_enabled'
);

/* 2023-01-23 */
ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN product_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN item_flag SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN cat_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN vendor_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN brand_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN shop_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN shop_cat_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN express_tid SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN title SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN short_title SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN code SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN image SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN price_range SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN stock_num SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN sale_num SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN sku_num SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN sku_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN cost SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN price SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN retail_price SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN weight SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN bulk SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN shelve_state SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN audit_state SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN audit_remark SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN sort_num SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN create_time SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ALTER COLUMN update_time SET NOT NULL;

ALTER TABLE IF EXISTS public.item_info
    ADD COLUMN intro_video character varying(120) NOT NULL DEFAULT '' ;

COMMENT ON COLUMN public.item_info.intro_video
    IS '视频介绍';



-- Table: public.order_distribution

DROP TABLE IF EXISTS public.order_rebate_list; 
DROP TABLE IF EXISTS public.order_rebate_item; 

-- DROP TABLE IF EXISTS public.order_distribution;

CREATE TABLE IF NOT EXISTS public.order_distribution
(
    id BIGSERIAL NOT NULL,
    plan_id integer NOT NULL,
    buyer_id bigint NOT NULL,
    owner_id bigint NOT NULL,
    flag int4range NOT NULL,
    is_read int4range NOT NULL,
    affiliate_code character varying(20) COLLATE pg_catalog."default" NOT NULL,
    order_no character varying(20) COLLATE pg_catalog."default" NOT NULL,
    order_subject character varying(40) COLLATE pg_catalog."default" NOT NULL,
    order_amount bigint NOT NULL,
    distribute_amount bigint NOT NULL,
    status integer NOT NULL,
    create_time bigint NOT NULL,
    update_time bigint NOT NULL,
    CONSTRAINT order_distribution_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.order_distribution
    OWNER to postgres;

COMMENT ON TABLE public.order_distribution
    IS '订单分销';

COMMENT ON COLUMN public.order_distribution.id
    IS '编号';

COMMENT ON COLUMN public.order_distribution.plan_id
    IS '返利方案Id';

COMMENT ON COLUMN public.order_distribution.buyer_id
    IS '买家';

COMMENT ON COLUMN public.order_distribution.affiliate_code
    IS '分享码';

COMMENT ON COLUMN public.order_distribution.order_no
    IS '订单号';

COMMENT ON COLUMN public.order_distribution.order_subject
    IS '订单标题';

COMMENT ON COLUMN public.order_distribution.order_amount
    IS '订单金额';

COMMENT ON COLUMN public.order_distribution.distribute_amount
    IS '分销奖励金额';

COMMENT ON COLUMN public.order_distribution.status
    IS '返利状态';

COMMENT ON COLUMN public.order_distribution.create_time
    IS '创建时间';

COMMENT ON COLUMN public.order_distribution.update_time
    IS '更新时间';

COMMENT ON COLUMN public.order_distribution.owner_id
    IS '返利所有人编号';

COMMENT ON COLUMN public.order_distribution.is_read
    IS '是否已读';

COMMENT ON COLUMN public.order_distribution.flag
    IS '标志';

-- Table: public.order_distribution_item

-- DROP TABLE IF EXISTS public.order_distribution_item;

CREATE TABLE IF NOT EXISTS public.order_distribution_item
(
    id BIGSERIAL NOT NULL,
    distribute_id bigint NOT NULL,
    item_id bigint NOT NULL,
    item_name character varying(20) COLLATE pg_catalog."default" NOT NULL,
    item_image character varying(120) COLLATE pg_catalog."default" NOT NULL,
    item_amount bigint NOT NULL,
    distribute_amount bigint NOT NULL,
    CONSTRAINT order_distribution_item_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.order_distribution_item
    OWNER to postgres;

COMMENT ON TABLE public.order_distribution_item
    IS '订单分销详情';

COMMENT ON COLUMN public.order_distribution_item.id
    IS '编号';

COMMENT ON COLUMN public.order_distribution_item.distribute_id
    IS '分销单Id';

COMMENT ON COLUMN public.order_distribution_item.item_id
    IS '商品编号';

COMMENT ON COLUMN public.order_distribution_item.item_name
    IS '商品名称';

COMMENT ON COLUMN public.order_distribution_item.item_image
    IS '商品图片';

COMMENT ON COLUMN public.order_distribution_item.item_amount
    IS '商品金额';

COMMENT ON COLUMN public.order_distribution_item.distribute_amount
    IS '分销金额';


ALTER TABLE public.mch_express_template
    ALTER COLUMN first_fee TYPE integer;
COMMENT ON COLUMN public.mch_express_template.first_fee
    IS '首次计价单价(元),如续重1kg';

ALTER TABLE public.mch_express_template
    ALTER COLUMN add_fee TYPE integer;
COMMENT ON COLUMN public.mch_express_template.add_fee
    IS '超过首次计价单价(元)，如续重1kg';