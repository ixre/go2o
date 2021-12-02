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