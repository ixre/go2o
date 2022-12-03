/** 2022-11-25 */
ALTER TABLE IF EXISTS public.perm_dept
    ADD COLUMN code character varying(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN public.perm_dept.code
    IS '部门代码';

ALTER TABLE public.order_list ALTER COLUMN id TYPE bigint;

ALTER TABLE public.order_list
    ALTER COLUMN subject TYPE character varying(40) COLLATE pg_catalog."default";


/** 2022-12-03 */
ALTER TABLE IF EXISTS public.order_list
    ADD COLUMN consignee_modified integer NOT NULL DEFAULT 0;

COMMENT ON COLUMN public.order_list.consignee_modified
    IS '收货地址是否已修改';

ALTER TABLE IF EXISTS public.sale_sub_order
    RENAME is_suspend TO is_forbidden;

COMMENT ON COLUMN public.sale_sub_order.is_forbidden
    IS '是否被用户删除';