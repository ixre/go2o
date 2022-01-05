delete  FROM mm_member where id>100

DELETE FROM mm_profile WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM mm_account WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM mm_relation WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM mm_trusted_info WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM mm_receipts_code WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM mm_lock_info WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM mm_lock_history WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM mm_deliver_addr WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM mm_wallet_log WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM mm_balance_log WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM mm_integral_log WHERE member_id NOT IN (select id FROM mm_member);
DELETE FROM wal_wallet WHERE wallet_type=1 AND user_id NOT IN (select id FROM mm_member);
DELETE FROM wal_wallet_log WHERE wallet_id NOT IN (select id FROM wal_wallet);
