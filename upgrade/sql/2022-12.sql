/** 添加数据库字典 */
ALTER TABLE "public".registry 
  alter column group_name set default ''::character varying;
COMMENT ON TABLE "public".registry IS '注册表';
COMMENT ON COLUMN "public".registry."key" IS '键';
COMMENT ON COLUMN "public".registry.flag IS '是否用户定义,0:否,1:是';
COMMENT ON COLUMN "public".registry.value IS '值';
COMMENT ON COLUMN "public".registry.description IS '描述';
COMMENT ON COLUMN "public".registry.default_value IS '默认值';
COMMENT ON COLUMN "public".registry.options IS '可选值';
COMMENT ON COLUMN "public".registry.group_name IS '分组名称';


COMMENT ON TABLE "public".app_version IS 'APP版本';
COMMENT ON COLUMN "public".app_version.id IS '编号';
COMMENT ON COLUMN "public".app_version.product_id IS '产品';
COMMENT ON COLUMN "public".app_version.channel IS '更新通道, stable:0|alpha:1|nightly:2';
COMMENT ON COLUMN "public".app_version.version IS '版本号';
COMMENT ON COLUMN "public".app_version.version_code IS '数字版本';
COMMENT ON COLUMN "public".app_version.force_update IS '是否强制升级';
COMMENT ON COLUMN "public".app_version.update_content IS '更新内容';
COMMENT ON COLUMN "public".app_version.create_time IS '发布时间';


COMMENT ON TABLE "public".app_prod IS 'APP产品';
COMMENT ON COLUMN "public".app_prod.id IS '产品编号';
COMMENT ON COLUMN "public".app_prod.prod_name IS '产品名称';
COMMENT ON COLUMN "public".app_prod.prod_des IS '产品描述';
COMMENT ON COLUMN "public".app_prod.latest_vid IS '最新的版本ID';
COMMENT ON COLUMN "public".app_prod.md5_hash IS '正式版文件hash值';
COMMENT ON COLUMN "public".app_prod.publish_url IS '发布下载页面地址';
COMMENT ON COLUMN "public".app_prod.stable_file_url IS '正式版文件地址';
COMMENT ON COLUMN "public".app_prod.alpha_file_url IS '内测版文件地址';
COMMENT ON COLUMN "public".app_prod.nightly_file_url IS '每夜版文件地址';
COMMENT ON COLUMN "public".app_prod.update_type IS '更新方式,比如APK, EXE等';
COMMENT ON COLUMN "public".app_prod.update_time IS '更新时间';

COMMENT ON TABLE "public".mm_lock_info IS '会员锁定记录';
COMMENT ON COLUMN "public".mm_lock_info.id IS '编号';
COMMENT ON COLUMN "public".mm_lock_info.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_lock_info.lock_time IS '锁定时间';
COMMENT ON COLUMN "public".mm_lock_info.unlock_time IS '解锁时间';
COMMENT ON COLUMN "public".mm_lock_info.remark IS '备注';

COMMENT ON TABLE "public".mm_lock_history IS '会员锁定历史';
COMMENT ON COLUMN "public".mm_lock_history.id IS '编号';
COMMENT ON COLUMN "public".mm_lock_history.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_lock_history.lock_time IS '锁定时间';
COMMENT ON COLUMN "public".mm_lock_history.duration IS '锁定持续分钟数';
COMMENT ON COLUMN "public".mm_lock_history.remark IS '备注';

COMMENT ON TABLE "public".mm_bank_card IS '银行卡';
COMMENT ON COLUMN "public".mm_bank_card.id IS '编号';
COMMENT ON COLUMN "public".mm_bank_card.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_bank_card.bank_account IS '银行账号';
COMMENT ON COLUMN "public".mm_bank_card.account_name IS '户名';
COMMENT ON COLUMN "public".mm_bank_card.bank_id IS '银行编号';
COMMENT ON COLUMN "public".mm_bank_card.bank_name IS '银行名称';
COMMENT ON COLUMN "public".mm_bank_card.bank_code IS '银行卡代码';
COMMENT ON COLUMN "public".mm_bank_card.network IS '网点';
COMMENT ON COLUMN "public".mm_bank_card.auth_code IS '快捷支付授权码';
COMMENT ON COLUMN "public".mm_bank_card.state IS '状态';
COMMENT ON COLUMN "public".mm_bank_card.create_time IS '添加时间';

/** 2022-12 */
-- Table: public.sys_member_flag

-- DROP TABLE IF EXISTS public.sys_member_flag;

CREATE TABLE IF NOT EXISTS public.sys_member_flag
(
    id bigint NOT NULL,
    flag_name character varying(20) COLLATE pg_catalog."default" NOT NULL,
    flag_value integer NOT NULL,
    update_time bigint NOT NULL,
    editable integer NOT NULL,
    is_show integer NOT NULL,
    CONSTRAINT sys_member_flag_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.sys_member_flag
    OWNER to postgres;

COMMENT ON TABLE public.sys_member_flag
    IS '会员标志名称';

COMMENT ON COLUMN public.sys_member_flag.flag_name
    IS '标志名称';

COMMENT ON COLUMN public.sys_member_flag.flag_value
    IS '标志值';

COMMENT ON COLUMN public.sys_member_flag.update_time
    IS '更新时间';

COMMENT ON COLUMN public.sys_member_flag.editable
    IS '是否能编辑【0：否，1：是】';

COMMENT ON COLUMN public.sys_member_flag.is_show
    IS '是否展示【0：否，1：是】';

ALTER TABLE "public".mm_member 
  alter column code set default ''::character varying;

DROP TABLE mm_balance_info;

CREATE TABLE portal_nav_type (
  id   serial NOT NULL, 
  name varchar(20) NOT NULL, 
  CONSTRAINT portal_nav_type_pkey1
    PRIMARY KEY (id));
COMMENT ON COLUMN "public".portal_nav_type.id IS '编号';
COMMENT ON COLUMN "public".portal_nav_type.name IS '名称';
INSERT INTO portal_nav_type (id, name)
VALUES (1, 'PC商城');
INSERT INTO portal_nav_type (id, name)
VALUES (2, '移动商城');
