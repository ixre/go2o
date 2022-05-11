
-- DROP TABLE IF EXISTS public.pay_integrate_app;

CREATE TABLE IF NOT EXISTS public.pay_integrate_app
(
    id serial NOT NULL,
    app_name character varying(20) COLLATE pg_catalog."default" NOT NULL,
    app_url character varying(120) COLLATE pg_catalog."default" NOT NULL,
    integrate_type integer NOT NULL,
    sort_number integer NOT NULL,
    enabled integer NOT NULL,
    CONSTRAINT pay_integrate_app_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.pay_integrate_app
    OWNER to postgres;

COMMENT ON TABLE public.pay_integrate_app
    IS '集成支付应用';

COMMENT ON COLUMN public.pay_integrate_app.id
    IS '编号';

COMMENT ON COLUMN public.pay_integrate_app.app_name
    IS '支付应用名称';

COMMENT ON COLUMN public.pay_integrate_app.app_url
    IS '支付应用接口';

COMMENT ON COLUMN public.pay_integrate_app.enabled
    IS '是否启用';

COMMENT ON COLUMN public.pay_integrate_app.integrate_type
    IS '集成方式: 1:API调用 2: 跳转';

COMMENT ON COLUMN public.pay_integrate_app.sort_number
    IS '显示顺序';