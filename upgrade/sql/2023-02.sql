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