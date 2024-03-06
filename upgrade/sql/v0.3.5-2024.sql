/** 2024-03-03 12:31 */
ALTER TABLE IF EXISTS public.ws_item
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.mm_integral_log
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.mm_balance_log
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.wal_wallet_log
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.mm_flow_log
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.mm_trusted_info
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.mm_levelup
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.mch_sign_up
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.mch_enterprise_info
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.ws_wholesaler
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.product_brand
    RENAME review_state TO review_status;
ALTER TABLE IF EXISTS public.item_info
    RENAME review_state TO review_status;