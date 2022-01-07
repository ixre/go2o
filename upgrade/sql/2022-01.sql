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
