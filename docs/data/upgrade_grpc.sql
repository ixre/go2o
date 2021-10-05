ALTER TABLE public.product_brand DROP COLUMN review;

ALTER TABLE public.product ALTER COLUMN img TYPE character varying(200) COLLATE pg_catalog."default";

ALTER TABLE public.item_info
ALTER
COLUMN image TYPE character varying(200) COLLATE pg_catalog."default";

ALTER TABLE public.item_snapshot
ALTER
COLUMN image TYPE character varying(200) COLLATE pg_catalog."default";

ALTER TABLE public.item_trade_snapshot
ALTER
COLUMN img TYPE character varying(200) COLLATE pg_catalog."default";
