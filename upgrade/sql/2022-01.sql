-- Table: public.sys_search_word

-- DROP TABLE IF EXISTS public.sys_search_word;

CREATE TABLE IF NOT EXISTS public.sys_search_word
(
    id bigint NOT NULL,
    word character varying(20) COLLATE pg_catalog."default" NOT NULL,
    search_count integer NOT NULL,
    flag integer NOT NULL,
    CONSTRAINT sys_search_word_pkey PRIMARY KEY (id)
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
