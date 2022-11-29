/** 2022-11-25 */
ALTER TABLE IF EXISTS public.perm_dept
    ADD COLUMN code character varying(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN public.perm_dept.code
    IS '部门代码';

ALTER TABLE public.order_list ALTER COLUMN id TYPE bigint;

ALTER TABLE public.order_list
    ALTER COLUMN subject TYPE character varying(40) COLLATE pg_catalog."default";