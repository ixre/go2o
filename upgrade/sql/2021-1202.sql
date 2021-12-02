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