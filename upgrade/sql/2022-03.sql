/* 03-04 */
CREATE INDEX wallet_log_wallet_id_title ON public.wal_wallet_log (wallet_id,title);
CREATE INDEX wallet_log_wallet_id ON public.wal_wallet_log (wallet_id);

ALTER TABLE IF EXISTS public.wal_wallet_log
    ADD COLUMN wallet_user character varying(40) NOT NULL DEFAULT '';

COMMENT ON COLUMN public.wal_wallet_log.wallet_user
    IS '钱包用户';
