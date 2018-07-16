Connecting Database...
Generating DDL...
Exporting to database...
Executing SQL: ALTER TABLE pro_product
  DROP COLUMN sale_price
Executing SQL: ALTER TABLE ws_rebate_rate
  DROP COLUMN wss_id
Executing SQL: DROP TABLE _express_provider
Executing SQL: DROP TABLE pay_order_old
Executing SQL: ALTER TABLE ad_userset
  modify column pos_id int(11) NOT NULL
Executing SQL: ALTER TABLE ad_userset
  modify column user_id int(11) NOT NULL
Executing SQL: ALTER TABLE ad_userset
  modify column ad_id int(11) NOT NULL
Executing SQL: ALTER TABLE mch_sign_up
  modify column review_state int(1) NOT NULL
Executing SQL: ALTER TABLE mm_trusted_info
  modify column real_name varchar(10) NOT NULL
Executing SQL: ALTER TABLE mm_trusted_info
  modify column card_id varchar(20) NOT NULL
Executing SQL: ALTER TABLE mm_trusted_info
  modify column trust_image varchar(120) NOT NULL
Executing SQL: ALTER TABLE mm_trusted_info
  modify column review_state int(1) NOT NULL
Executing SQL: ALTER TABLE mm_trusted_info
  modify column review_time int(11) NOT NULL
Executing SQL: ALTER TABLE mm_trusted_info
  modify column remark varchar(120) NOT NULL
Executing SQL: ALTER TABLE mm_trusted_info
  modify column update_time int(11) NOT NULL
Executing SQL: ALTER TABLE pro_category
  modify column priority int(2) DEFAULT 0 NOT NULL
Executing SQL: ALTER TABLE pro_category
  modify column enabled int(2) NOT NULL
Executing SQL: ALTER TABLE mm_account
  alter column balance set default 0.00
Executing SQL: ALTER TABLE ws_rebate_rate
  ADD COLUMN ws_id int(10) NOT NULL comment '批发商编号'
Executing SQL: ALTER TABLE article_list
  modify column update_time int(11) NOT NULL
Generate database finish...