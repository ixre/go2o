ALTER TABLE "public".order_list 
  ADD COLUMN is_paid int4 DEFAULT 0 NOT NULL;
COMMENT ON COLUMN "public".order_list.id IS '编号';
COMMENT ON COLUMN "public".order_list.order_no IS '订单号';
COMMENT ON COLUMN "public".order_list.order_type IS '订单类型1:普通 2:批发 3:线下';
COMMENT ON COLUMN "public".order_list.subject IS '订单主题';
COMMENT ON COLUMN "public".order_list.buyer_id IS '买家';
COMMENT ON COLUMN "public".order_list.buyer_user IS '买家用户名';
COMMENT ON COLUMN "public".order_list.item_count IS '商品件数';
COMMENT ON COLUMN "public".order_list.item_amount IS '商品金额';
COMMENT ON COLUMN "public".order_list.discount_amount IS '抵扣金额';
COMMENT ON COLUMN "public".order_list.express_fee IS '物流费';
COMMENT ON COLUMN "public".order_list.package_fee IS '包装费';
COMMENT ON COLUMN "public".order_list.final_amount IS '订单最终金额';
COMMENT ON COLUMN "public".order_list.consignee_name IS '收货人姓名';
COMMENT ON COLUMN "public".order_list.consignee_phone IS '收货人电话';
COMMENT ON COLUMN "public".order_list.shipping_address IS '收货人地址';
COMMENT ON COLUMN "public".order_list.is_break IS '是否拆分';
COMMENT ON COLUMN "public".order_list.is_paid IS '是否已支付';
COMMENT ON COLUMN "public".order_list.state IS '订单状态';
COMMENT ON COLUMN "public".order_list.create_time IS '创建时间';
COMMENT ON COLUMN "public".order_list.update_time IS '更新时间';
CREATE INDEX order_list_is_paid 
  ON "public".order_list (is_paid);

ALTER TABLE "public".order_wholesale_order 
  DROP COLUMN is_paid;
COMMENT ON TABLE "public".order_wholesale_order IS '批发订单';
COMMENT ON COLUMN "public".order_wholesale_order.id IS '编号';
COMMENT ON COLUMN "public".order_wholesale_order.order_no IS '订单号';
COMMENT ON COLUMN "public".order_wholesale_order.order_id IS '订单编号';
COMMENT ON COLUMN "public".order_wholesale_order.buyer_id IS '买家';
COMMENT ON COLUMN "public".order_wholesale_order.vendor_id IS '供货商';
COMMENT ON COLUMN "public".order_wholesale_order.shop_id IS '店铺编号';
COMMENT ON COLUMN "public".order_wholesale_order.shop_name IS '店铺名称';
COMMENT ON COLUMN "public".order_wholesale_order.buyer_comment IS '买家留言';
COMMENT ON COLUMN "public".order_wholesale_order.remark IS '备注';
COMMENT ON COLUMN "public".order_wholesale_order.state IS '订单状态';
COMMENT ON COLUMN "public".order_wholesale_order.create_time IS '创建时间';
COMMENT ON COLUMN "public".order_wholesale_order.update_time IS '更新时间';

ALTER TABLE "public".sale_sub_order 
  DROP COLUMN is_paid;
COMMENT ON COLUMN "public".sale_sub_order.shop_name IS '店铺名称';



ALTER TABLE "public".order_trade_order 
 RENAME state TO status;


ALTER TABLE "public".order_wholesale_order 
 RENAME state TO status;

 ALTER TABLE "public".sale_sub_order 
 RENAME state TO status;


 ALTER TABLE "public".order_list 
 RENAME state TO status;


COMMENT ON COLUMN "public".order_trade_order.id IS '编号';
COMMENT ON COLUMN "public".order_trade_order.order_id IS '订单编号';
COMMENT ON COLUMN "public".order_trade_order.vendor_id IS '商家编号';
COMMENT ON COLUMN "public".order_trade_order.shop_id IS '店铺编号';
COMMENT ON COLUMN "public".order_trade_order.subject IS '订单标题';
COMMENT ON COLUMN "public".order_trade_order.order_amount IS '订单金额';
COMMENT ON COLUMN "public".order_trade_order.discount_amount IS '抵扣金额';
COMMENT ON COLUMN "public".order_trade_order.final_amount IS '订单最终金额';
COMMENT ON COLUMN "public".order_trade_order.trade_rate IS '交易结算比例（商户)';
COMMENT ON COLUMN "public".order_trade_order.cash_pay IS '是否现金支付';
COMMENT ON COLUMN "public".order_trade_order.ticket_image IS '发票图片';
COMMENT ON COLUMN "public".order_trade_order.remark IS '订单备注';
COMMENT ON COLUMN "public".order_trade_order.status IS '订单状态';
COMMENT ON COLUMN "public".order_trade_order.create_time IS '订单创建时间';
COMMENT ON COLUMN "public".order_trade_order.update_time IS '订单更新时间';

/** 06-07 */
 ALTER TABLE "public".sale_order_item 
  alter column order_id set not null;
ALTER TABLE "public".sale_order_item 
  ADD COLUMN seller_order_id int8 default 0 NOT NULL;
ALTER TABLE "public".sale_order_item 
  alter column item_id set not null;
ALTER TABLE "public".sale_order_item 
  alter column sku_id set not null;
ALTER TABLE "public".sale_order_item 
  alter column snap_id set not null;
ALTER TABLE "public".sale_order_item 
  alter column quantity set not null;
ALTER TABLE "public".sale_order_item 
  alter column return_quantity set not null;
ALTER TABLE "public".sale_order_item 
  alter column amount set not null;
ALTER TABLE "public".sale_order_item 
  alter column final_amount set not null;
ALTER TABLE "public".sale_order_item 
  alter column is_shipped set not null;
ALTER TABLE "public".sale_order_item 
  alter column update_time set not null;
COMMENT ON COLUMN "public".sale_order_item.id IS '编号';
COMMENT ON COLUMN "public".sale_order_item.order_id IS '订单编号,未支付时并非卖家订单的编号';
COMMENT ON COLUMN "public".sale_order_item.seller_order_id IS '卖家订单编号';
COMMENT ON COLUMN "public".sale_order_item.item_id IS '商品编号';
COMMENT ON COLUMN "public".sale_order_item.sku_id IS 'SKU编号';
COMMENT ON COLUMN "public".sale_order_item.snap_id IS '销售快照编号';
COMMENT ON COLUMN "public".sale_order_item.quantity IS '数量';
COMMENT ON COLUMN "public".sale_order_item.return_quantity IS '退货数量';
COMMENT ON COLUMN "public".sale_order_item.amount IS '金额';
COMMENT ON COLUMN "public".sale_order_item.final_amount IS '最终金额';
COMMENT ON COLUMN "public".sale_order_item.is_shipped IS '是否发货';
COMMENT ON COLUMN "public".sale_order_item.update_time IS '更新时间';

ALTER TABLE "public".sale_sub_order 
  ADD COLUMN break_status int2 DEFAULT 0 NOT NULL;
ALTER TABLE "public".sale_sub_order 
  alter column update_time set default 0;
COMMENT ON COLUMN "public".sale_sub_order.shop_id IS '店铺编号';
COMMENT ON COLUMN "public".sale_sub_order.shop_name IS '店铺名称';
COMMENT ON COLUMN "public".sale_sub_order.break_status IS '拆分状态: 0.默认 1:待拆分 2:无需拆分 3:已拆分';


/* 20220623 */
CREATE TABLE IF NOT EXISTS public.sale_after_order
(
    id bigserial NOT NULL,
    order_no character varying(20),
    order_id bigint NOT NULL,
    vendor_id bigint NOT NULL,
    buyer_id bigint NOT NULL,
    type smallint NOT NULL,
    snapshot_id bigint NOT NULL,
    quantity integer NOT NULL,
    reason character varying(255) COLLATE pg_catalog."default" NOT NULL,
    image_url character varying(255) COLLATE pg_catalog."default" NOT NULL,
    person_name character varying(10) COLLATE pg_catalog."default" NOT NULL,
    person_phone character varying(20) COLLATE pg_catalog."default" NOT NULL,
    shipment_express character varying(10) COLLATE pg_catalog."default" NOT NULL,
    shipment_order_no character varying(20) COLLATE pg_catalog."default" NOT NULL,
    shipment_image character varying(120) COLLATE pg_catalog."default" NOT NULL,
    remark character varying(45) COLLATE pg_catalog."default" NOT NULL,
    vendor_remark character varying(45) COLLATE pg_catalog."default" NOT NULL,
    status int2 NOT NULL,
    create_time bigint NOT NULL,
    update_time bigint NOT NULL,
    CONSTRAINT sale_after_order_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.sale_after_order
    OWNER to postgres;

COMMENT ON COLUMN public.sale_after_order.order_id
    IS '关联销售订单号';

COMMENT ON COLUMN public.sale_after_order.vendor_id
    IS '卖家';

COMMENT ON COLUMN public.sale_after_order.buyer_id
    IS '买家';

COMMENT ON COLUMN public.sale_after_order.type
    IS '售后类型 ';

COMMENT ON COLUMN public.sale_after_order.snapshot_id
    IS '商品快照编号';

COMMENT ON COLUMN public.sale_after_order.quantity
    IS '数量';

COMMENT ON COLUMN public.sale_after_order.reason
    IS '申请售后原因';

COMMENT ON COLUMN public.sale_after_order.image_url
    IS '商品售后图片凭证';

COMMENT ON COLUMN public.sale_after_order.person_name
    IS '联系人';

COMMENT ON COLUMN public.sale_after_order.person_phone
    IS '联系电话';

COMMENT ON COLUMN public.sale_after_order.shipment_express
    IS '退货快递名称';

COMMENT ON COLUMN public.sale_after_order.shipment_order_no
    IS '退货快递单号';

COMMENT ON COLUMN public.sale_after_order.shipment_image
    IS '退货凭证';

COMMENT ON COLUMN public.sale_after_order.remark
    IS '备注';

COMMENT ON COLUMN public.sale_after_order.vendor_remark
    IS '供应商备注';

COMMENT ON COLUMN public.sale_after_order.status
    IS '状态';

COMMENT ON COLUMN public.sale_after_order.create_time
    IS '创建时间';

COMMENT ON COLUMN public.sale_after_order.update_time
    IS '更新时间';

COMMENT ON COLUMN public.sale_after_order.order_no
    IS '售后单单号';


    ALTER TABLE "public".sale_after_order 
  ADD COLUMN shipment_time bigint DEFAULT 0 NOT NULL;


COMMENT ON COLUMN public.sale_after_order.shipment_time
    IS '发货时间';