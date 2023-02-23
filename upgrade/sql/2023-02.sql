delete FROM registry where key in ('order_disallow_user_cancel',"domain_file_server_prefix");

/** 2023-02-17 09:50 */
ALTER TABLE IF EXISTS public.pay_order
    RENAME final_fee TO final_amount;

ALTER TABLE IF EXISTS public.pay_order
    RENAME paid_fee TO paid_amount;

ALTER TABLE IF EXISTS public.pay_order
    RENAME pay_uid TO payer_id;

ALTER TABLE IF EXISTS public.sale_sub_order
    ADD COLUMN item_count integer NOT NULL DEFAULT 0;

COMMENT ON COLUMN public.sale_sub_order.item_count
    IS '商品数量';

update sale_sub_order set item_count = 
(SELECT  coalesce(SUM(quantity),0) FROM sale_order_item 
 WHERE order_id = sale_sub_order.id)
WHERE item_count = 0;

/** 2023-02-20 更改收货地址的is_default类型 */
ALTER TABLE public.mm_deliver_addr
    ALTER COLUMN is_default TYPE int USING is_default::integer; 

/** 2023-02-23 菜单添加is_forbidden */
ALTER TABLE "public".perm_res 
  alter column cache_ set default ''::character varying;
ALTER TABLE "public".perm_res 
  ADD COLUMN is_forbidden int2 DEFAULT 0 NOT NULL;
COMMENT ON COLUMN "public".perm_res.id IS '资源ID';
COMMENT ON COLUMN "public".perm_res.pid IS '上级菜单ID';
COMMENT ON COLUMN "public".perm_res.name IS '资源名称';
COMMENT ON COLUMN "public".perm_res.res_type IS '资源类型, 0: 目录 1: 资源　2: 菜单  3:　 按钮';
COMMENT ON COLUMN "public".perm_res."key" IS '资源键';
COMMENT ON COLUMN "public".perm_res.path IS '资源路径';
COMMENT ON COLUMN "public".perm_res.icon IS '图标';
COMMENT ON COLUMN "public".perm_res.permission IS '权限,多个值用|分隔';
COMMENT ON COLUMN "public".perm_res.sort_num IS '排序';
COMMENT ON COLUMN "public".perm_res.is_external IS '是否外部';
COMMENT ON COLUMN "public".perm_res.is_hidden IS '是否隐藏';
COMMENT ON COLUMN "public".perm_res.create_time IS '创建日期';
COMMENT ON COLUMN "public".perm_res.component_name IS '组件路径';
COMMENT ON COLUMN "public".perm_res.cache_ IS '缓存';
COMMENT ON COLUMN "public".perm_res.depth IS '深度/层级';
COMMENT ON COLUMN "public".perm_res.is_forbidden IS '是否禁止';
