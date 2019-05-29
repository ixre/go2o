-- 更换为POSTGRESQL --
ALTER TABLE "public".mm_member ADD COLUMN flag int4 DEFAULT 0 NOT NULL;

-- mm_balance_log  mm_wallet_log  mm_integral_log 表结构更改
CREATE TABLE "public".mm_integral_log (id serial NOT NULL, member_id int4 NOT NULL, kind int4 NOT NULL, title varchar(60) DEFAULT '""'::character varying NOT NULL, outer_no varchar(40) DEFAULT '""'::character varying NOT NULL, value int4 NOT NULL, remark varchar(40) NOT NULL, rel_user int4 DEFAULT 0 NOT NULL, review_state int2 DEFAULT 0 NOT NULL, create_time int8 NOT NULL, update_time int8 DEFAULT 0 NOT NULL, CONSTRAINT mm_integral_log_pkey PRIMARY KEY (id));
COMMENT ON TABLE "public".mm_integral_log IS '积分明细';
COMMENT ON COLUMN "public".mm_integral_log.id IS '编号';
COMMENT ON COLUMN "public".mm_integral_log.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_integral_log.kind IS '类型';
COMMENT ON COLUMN "public".mm_integral_log.title IS '标题';
COMMENT ON COLUMN "public".mm_integral_log.outer_no IS '关联的编号';
COMMENT ON COLUMN "public".mm_integral_log.value IS '积分值';
COMMENT ON COLUMN "public".mm_integral_log.remark IS '备注';
COMMENT ON COLUMN "public".mm_integral_log.rel_user IS '关联用户';
COMMENT ON COLUMN "public".mm_integral_log.review_state IS '审核状态';
COMMENT ON COLUMN "public".mm_integral_log.create_time IS '创建时间';
COMMENT ON COLUMN "public".mm_integral_log.update_time IS '更新时间';

