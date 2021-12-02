ALTER TABLE "public".mch_merchant
    ADD COLUMN salt varchar(10) DEFAULT '' NOT NULL;
COMMENT
ON COLUMN "public".mch_merchant.salt IS '加密盐';

ALTER TABLE public.product
ALTER
COLUMN name TYPE character varying(120) COLLATE pg_catalog."default";


ALTER TABLE IF EXISTS public.product DROP COLUMN IF EXISTS shelve_state;
ALTER TABLE IF EXISTS public.product DROP COLUMN IF EXISTS review_state;
ALTER TABLE IF EXISTS public.product DROP COLUMN IF EXISTS sale_price;


/* 用bigint存储金额　*/
ALTER TABLE public.item_info
ALTER COLUMN cost TYPE bigint;

ALTER TABLE public.item_info
ALTER COLUMN price TYPE bigint;

ALTER TABLE public.item_info
ALTER COLUMN retail_price TYPE bigint;

ALTER TABLE public.item_sku
ALTER COLUMN retail_price TYPE bigint;

ALTER TABLE public.item_sku
ALTER COLUMN price TYPE bigint;

ALTER TABLE public.item_sku
ALTER COLUMN cost TYPE bigint;


ALTER TABLE public.item_snapshot
ALTER COLUMN cost TYPE bigint;

ALTER TABLE public.item_snapshot
ALTER COLUMN price TYPE bigint;

ALTER TABLE public.item_snapshot
ALTER COLUMN retail_price TYPE bigint;

ALTER TABLE public.item_trade_snapshot
ALTER COLUMN cost TYPE bigint;

ALTER TABLE public.item_trade_snapshot
ALTER COLUMN price TYPE bigint;

ALTER TABLE public.gs_member_price
ALTER COLUMN price TYPE bigint;

ALTER TABLE public.ws_item
ALTER COLUMN price TYPE bigint;


ALTER TABLE public.sale_order
ALTER COLUMN item_amount TYPE bigint;

ALTER TABLE public.sale_order
ALTER COLUMN discount_amount TYPE bigint;

ALTER TABLE public.sale_order
ALTER COLUMN express_fee TYPE bigint;

ALTER TABLE public.sale_order
ALTER COLUMN package_fee TYPE bigint;

ALTER TABLE public.sale_order
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.sale_order_item
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.sale_order_item
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.sale_sub_order
ALTER COLUMN discount_amount TYPE bigint;

ALTER TABLE public.sale_sub_order
ALTER COLUMN express_fee TYPE bigint;

ALTER TABLE public.sale_sub_order
ALTER COLUMN package_fee TYPE bigint;

ALTER TABLE public.sale_sub_order
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN freeze_balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN expired_balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN wallet_balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN freeze_wallet TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN expired_wallet TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN total_wallet_amount TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN flow_balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN grow_balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN grow_amount TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN grow_earnings TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN grow_total_earnings TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN total_charge TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN total_pay TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN total_expense TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN priority_pay TYPE integer;

ALTER TABLE public.express_area_set
ALTER COLUMN first_fee TYPE bigint;

ALTER TABLE public.express_area_set
ALTER COLUMN add_fee TYPE bigint;

ALTER TABLE public.mch_balance_log
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.mch_balance_log
ALTER COLUMN csn_amount TYPE bigint;

ALTER TABLE public.mm_balance_info
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.mm_balance_info
ALTER COLUMN csn_amount TYPE bigint;

ALTER TABLE public.mm_balance_log
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.mm_balance_log
ALTER COLUMN csn_fee TYPE bigint;

ALTER TABLE public.mm_flow_log
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.mm_flow_log
ALTER COLUMN csn_fee TYPE bigint;

ALTER TABLE public.mm_wallet_log
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.mm_wallet_log
ALTER COLUMN csn_fee TYPE bigint;

ALTER TABLE public.order_wholesale_order
ALTER COLUMN item_amount TYPE bigint;

ALTER TABLE public.order_wholesale_order
ALTER COLUMN discount_amount TYPE bigint;

ALTER TABLE public.order_wholesale_order
ALTER COLUMN express_fee TYPE bigint;

ALTER TABLE public.order_wholesale_order
ALTER COLUMN package_fee TYPE bigint;

ALTER TABLE public.order_wholesale_order
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.order_wholesale_item
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.order_wholesale_item
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.order_trade_order
ALTER COLUMN order_amount TYPE bigint;

ALTER TABLE public.order_trade_order
ALTER COLUMN discount_amount TYPE bigint;

ALTER TABLE public.order_trade_order
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.order_trade_order
ALTER COLUMN trade_rate TYPE bigint;

COMMENT ON TABLE public.item_info
  IS '商品信息';

COMMENT ON COLUMN public.item_info.id
    IS '编号';

COMMENT ON COLUMN public.item_info.product_id
    IS '产品编号';

COMMENT ON COLUMN public.item_info.prom_flag
    IS '营销标志';

COMMENT ON COLUMN public.item_info.cat_id
    IS '分类编号';

COMMENT ON COLUMN public.item_info.vendor_id
    IS '供应商编号';

COMMENT ON COLUMN public.item_info.brand_id
    IS '品牌编号';

COMMENT ON COLUMN public.item_info.shop_id
    IS '店铺编号';

COMMENT ON COLUMN public.item_info.shop_cat_id
    IS '店铺分类编号';

COMMENT ON COLUMN public.item_info.express_tid
    IS '快递模板';

COMMENT ON COLUMN public.item_info.title
    IS '商品标题';

COMMENT ON COLUMN public.item_info.short_title
    IS '商品小标题';

COMMENT ON COLUMN public.item_info.code
    IS '商品编码';

COMMENT ON COLUMN public.item_info.image
    IS '商品主图';

COMMENT ON COLUMN public.item_info.is_present
    IS '是否为赠品';

COMMENT ON COLUMN public.item_info.price_range
    IS '价格区间';

COMMENT ON COLUMN public.item_info.stock_num
    IS '库存数量';

COMMENT ON COLUMN public.item_info.sale_num
    IS '销售数量';

COMMENT ON COLUMN public.item_info.sku_num
    IS '规格数量';

COMMENT ON COLUMN public.item_info.sku_id
    IS 'SKU编号';

COMMENT ON COLUMN public.item_info.cost
    IS '成本价';

COMMENT ON COLUMN public.item_info.price
    IS '销售价';

COMMENT ON COLUMN public.item_info.retail_price
    IS '零售价';

COMMENT ON COLUMN public.item_info.weight
    IS '重量';

COMMENT ON COLUMN public.item_info.bulk
    IS '容积';

COMMENT ON COLUMN public.item_info.shelve_state
    IS '上架状态';

COMMENT ON COLUMN public.item_info.review_state
    IS '审核状态';

COMMENT ON COLUMN public.item_info.review_remark
    IS '审核意见';

COMMENT ON COLUMN public.item_info.sort_num
    IS '排列序号';

COMMENT ON COLUMN public.item_info.create_time
    IS '创建时间';

COMMENT ON COLUMN public.item_info.update_time
    IS '更新时间';