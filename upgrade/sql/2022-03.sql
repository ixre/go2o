/* 03-04 */
CREATE INDEX wallet_log_wallet_id_title ON public.wal_wallet_log (wallet_id,title);
CREATE INDEX wallet_log_wallet_id ON public.wal_wallet_log (wallet_id);

ALTER TABLE IF EXISTS public.wal_wallet_log
    ADD COLUMN wallet_user character varying(40) NOT NULL DEFAULT '';

COMMENT ON COLUMN public.wal_wallet_log.wallet_user
    IS '钱包用户';

ALTER TABLE public.wal_wallet_log RENAME COLUMN trade_fee TO procedure_fee;


-- Table: public.job_exec_data

-- DROP TABLE IF EXISTS public.job_exec_data;

CREATE TABLE IF NOT EXISTS public.job_exec_data
(
    id bigserial NOT NULL,
    job_name character varying(40) NOT NULL,
    last_exec_index bigint NOT NULL,
    last_exec_time bigint,
    CONSTRAINT job_data_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.job_exec_data
    OWNER to postgres;

COMMENT ON COLUMN public.job_exec_data.id
    IS '编号';

COMMENT ON COLUMN public.job_exec_data.job_name
    IS '任务名称';

COMMENT ON COLUMN public.job_exec_data.last_exec_index
    IS '上次执行位置索引';

COMMENT ON COLUMN public.job_exec_data.last_exec_time
    IS '最后执行时间';


-- Table: public.job_exec_fail

-- DROP TABLE IF EXISTS public.job_exec_fail;

CREATE TABLE IF NOT EXISTS public.job_exec_fail
(
    id bigserial NOT NULL,
    job_id bigint NOT NULL,
    job_data_id bigint NOT NULL,
    retry_count integer NOT NULL,
    create_time bigint NOT NULL,
    retry_time bigint NOT NULL,
    CONSTRAINT job_exec_fail_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.job_exec_fail
    OWNER to postgres;

COMMENT ON TABLE public.job_exec_fail
    IS '任务执行失败';

COMMENT ON COLUMN public.job_exec_fail.id
    IS '编号';

COMMENT ON COLUMN public.job_exec_fail.job_id
    IS '任务编号';

COMMENT ON COLUMN public.job_exec_fail.job_data_id
    IS '任务数据编号';

COMMENT ON COLUMN public.job_exec_fail.retry_count
    IS '重试次数';

COMMENT ON COLUMN public.job_exec_fail.create_time
    IS '创建时间';

COMMENT ON COLUMN public.job_exec_fail.retry_time
    IS '重试时间';

/* 03-17 */
-- Table: public.exec_re_queue

-- DROP TABLE IF EXISTS public.exec_re_queue;

CREATE TABLE IF NOT EXISTS public.exec_re_queue
(
    id bigserial NOT NULL,
    queue_name character varying(40) COLLATE pg_catalog."default" NOT NULL,
    relate_id bigint NOT NULL,
    relate_data text COLLATE pg_catalog."default" NOT NULL,
    create_time bigint NOT NULL,
    CONSTRAINT exec_re_queue_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.exec_re_queue
    OWNER to postgres;

COMMENT ON TABLE public.exec_re_queue
    IS '重新加入队列';

COMMENT ON COLUMN public.exec_re_queue.queue_name
    IS '队列名称';

COMMENT ON COLUMN public.exec_re_queue.relate_id
    IS '关联数据编号';

COMMENT ON COLUMN public.exec_re_queue.relate_data
    IS '数据';

COMMENT ON COLUMN public.exec_re_queue.create_time
    IS '创建时间';

CREATE INDEX exec_re_queue_queue_name ON public.exec_re_queue (queue_name);
