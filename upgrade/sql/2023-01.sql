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