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


UPDATE item_snapshot SET item_flag = (SELECT item_flag FROM item_info where id=item_id) WHERE item_flag = 0;