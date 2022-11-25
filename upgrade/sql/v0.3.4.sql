/** 2022-11-25 */
ALTER TABLE IF EXISTS public.perm_dept
    ADD COLUMN code character varying(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN public.perm_dept.code
    IS '部门代码';