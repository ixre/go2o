-- Table: public.sys_search_word

-- DROP TABLE IF EXISTS public.sys_search_word;

CREATE TABLE IF NOT EXISTS public.sys_search_word
(
    id serial NOT NULL,
    word character varying(20) COLLATE pg_catalog."default" NOT NULL,
    search_count integer NOT NULL,
    flag integer NOT NULL
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.sys_search_word
    OWNER to postgres;

COMMENT ON TABLE public.sys_search_word
    IS '热搜词';

COMMENT ON COLUMN public.sys_search_word.id
    IS '编号';

COMMENT ON COLUMN public.sys_search_word.search_count
    IS '搜索次数';

COMMENT ON COLUMN public.sys_search_word.flag
    IS '1:启用　2:特殊显示 4: 手动创建';

/** 2022-01-07 */
ALTER TABLE public.mm_integral_log
ALTER COLUMN value TYPE bigint;

ALTER TABLE IF EXISTS public.mm_integral_log
    ADD COLUMN balance bigint NOT NULL DEFAULT 0;

COMMENT ON COLUMN public.mm_integral_log.balance
    IS '变动后的余额';

ALTER TABLE IF EXISTS public.mm_balance_log
    ADD COLUMN balance bigint NOT NULL DEFAULT 0;

COMMENT ON COLUMN public.mm_balance_log.balance
    IS '变动后的余额';

DROP TABLE public.mm_balance_info;
DROP TABLE public.mm_income_log;

/** 2022-01-10 */
ALTER TABLE public.wal_wallet_log
ALTER COLUMN id TYPE bigint;

ALTER TABLE public.wal_wallet_log
ALTER COLUMN value TYPE bigint;

ALTER TABLE public.wal_wallet_log
ALTER COLUMN balance TYPE bigint;


ALTER TABLE public.wal_wallet
ALTER COLUMN id TYPE bigint;

ALTER TABLE public.wal_wallet
ALTER COLUMN balance TYPE bigint;

ALTER TABLE public.wal_wallet
ALTER COLUMN present_balance TYPE bigint;

ALTER TABLE public.wal_wallet
ALTER COLUMN latest_amount TYPE bigint;

ALTER TABLE public.wal_wallet
ALTER COLUMN total_charge TYPE bigint;

ALTER TABLE public.wal_wallet
ALTER COLUMN total_pay TYPE bigint;

ALTER TABLE public.wal_wallet
ALTER COLUMN adjust_amount TYPE bigint;

ALTER TABLE public.wal_wallet
ALTER COLUMN freeze_amount TYPE bigint;

ALTER TABLE public.wal_wallet
ALTER COLUMN expired_amount TYPE bigint;

ALTER TABLE public.wal_wallet
ALTER COLUMN total_present TYPE bigint;

/** 2022-01-19 */
/** 更新商品价格 */
update item_info set price=CAST(replace(price_range,'.','') AS integer) where price =0


-- Table: public.item_image

-- DROP TABLE IF EXISTS public.item_image;

CREATE TABLE IF NOT EXISTS public.item_image
(
    id bigserial NOT NULL,
    item_id bigint NOT NULL,
    image_url character varying(180) COLLATE pg_catalog."default" NOT NULL,
    sort_num integer NOT NULL,
    create_time bigint NOT NULL,
    CONSTRAINT item_image_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.item_image
    OWNER to postgres;

COMMENT ON TABLE public.item_image
    IS '产品图片';

COMMENT ON COLUMN public.item_image.id
    IS '图片编号';

COMMENT ON COLUMN public.item_image.item_id
    IS '商品编号';

COMMENT ON COLUMN public.item_image.image_url
    IS '图片地址';

COMMENT ON COLUMN public.item_image.sort_num
    IS '排列序号';

COMMENT ON COLUMN public.item_image.create_time
    IS '创建时间';
