delete FROM registry where key like 'uams_%'

/* 2023-03-05 23:54 商品快照规格 */

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN product_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN snapshot_key SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN cat_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN vendor_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN brand_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN shop_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN shop_cat_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN express_tid SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN title SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN short_title SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN code SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN image SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN price_range SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN sku_id SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN cost SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN price SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN retail_price SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN weight SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN bulk SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN level_sales SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN shelve_state SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ALTER COLUMN update_time SET NOT NULL;

ALTER TABLE IF EXISTS public.item_snapshot
    ADD COLUMN item_flag integer NOT NULL DEFAULT 0;

COMMENT ON COLUMN public.item_snapshot.item_flag
    IS '商品标志';


UPDATE item_snapshot SET item_flag =  COALESCE((SELECT item_flag FROM item_info where id=item_id),1) WHERE item_flag = 0;


CREATE TABLE public.sys_safeguard
(
    id bigserial NOT NULL,
    bind_type integer NOT NULL,
    flag integer NOT NULL DEFAULT 0,
    name character varying(20) NOT NULL,
    content character varying(120) NOT NULL,
    class_name character varying(20) NOT NULL,
    sort_num integer NOT NULL,
    enabled integer NOT NULL,
    is_internal integer NOT NULL,
    create_time bigint NOT NULL,
    update_time bigint NOT NULL,
    PRIMARY KEY (id)
);


ALTER TABLE IF EXISTS public.sys_safeguard
    OWNER to postgres;

COMMENT ON TABLE public.sys_safeguard
    IS '保障';

COMMENT ON COLUMN public.sys_safeguard.id
    IS '编号';

COMMENT ON COLUMN public.sys_safeguard.bind_type
    IS '绑定类型:1:店铺 2:商品';

COMMENT ON COLUMN public.sys_safeguard.flag
    IS '保障标志';

COMMENT ON COLUMN public.sys_safeguard.name
    IS '保障名称';

COMMENT ON COLUMN public.sys_safeguard.content
    IS '保障内容';

COMMENT ON COLUMN public.sys_safeguard.class_name
    IS '样式表类名';

COMMENT ON COLUMN public.sys_safeguard.sort_num
    IS '序号';

COMMENT ON COLUMN public.sys_safeguard.enabled
    IS '是否启用';

COMMENT ON COLUMN public.sys_safeguard.is_internal
    IS '是否内置';

COMMENT ON COLUMN public.sys_safeguard.create_time
    IS '创建时间';

COMMENT ON COLUMN public.sys_safeguard.update_time
    IS '更新时间';


ALTER TABLE IF EXISTS public.mm_balance_log
    RENAME audit_state TO review_state;


CREATE TABLE public.item_affiliate_rate
(
    id bigserial NOT NULL,
    item_id bigint NOT NULL,
    rate_r1 integer NOT NULL,
    rate_r2 integer NOT NULL,
    rate_c integer NOT NULL,
    origin_rate_r1 integer NOT NULL,
    origin_rate_r2 integer NOT NULL,
    origin_rate_c integer NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.item_affiliate_rate
    OWNER to postgres;

COMMENT ON COLUMN public.item_affiliate_rate.rate_r1
    IS '上一级比例';

COMMENT ON COLUMN public.item_affiliate_rate.rate_r2
    IS '上二级比例';

COMMENT ON COLUMN public.item_affiliate_rate.rate_c
    IS '自定义比例';

COMMENT ON COLUMN public.item_affiliate_rate.origin_rate_r1
    IS '历史上一级比例';

COMMENT ON COLUMN public.item_affiliate_rate.origin_rate_r2
    IS '历史上二级比例';

COMMENT ON COLUMN public.item_affiliate_rate.origin_rate_c
    IS '历史自定义比例';


ALTER TABLE IF EXISTS public.item_info
    ADD COLUMN is_recycle integer NOT NULL DEFAULT 0;

COMMENT ON COLUMN public.item_info.is_recycle
    IS '是否回收';

ALTER TABLE IF EXISTS public.item_info
    RENAME retail_price TO origin_price;

ALTER TABLE IF EXISTS public.item_sku
    RENAME retail_price TO origin_price;

ALTER TABLE IF EXISTS public.item_snapshot
    RENAME retail_price TO origin_price;

/* 2023-03-13 */

ALTER TABLE IF EXISTS public.item_info
    ADD COLUMN safeguard_flag integer NOT NULL DEFAULT 0;
COMMENT ON COLUMN public.item_info.safeguard_flag
    IS '购物保障';

ALTER TABLE IF EXISTS public.item_snapshot
    ADD COLUMN safeguard_flag integer NOT NULL DEFAULT 0;
COMMENT ON COLUMN public.item_snapshot.safeguard_flag
    IS '购物保障';