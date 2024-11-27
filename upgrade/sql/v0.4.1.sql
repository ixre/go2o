
-- 2024-11-09 日志 
DROP TABLE IF EXISTS "public".sys_log CASCADE;
CREATE TABLE sys_log (
  id               BIGSERIAL NOT NULL, 
  user_id          int4 NOT NULL, 
  username         varchar(20) NOT NULL, 
  log_level        int4 NOT NULL, 
  message          varchar(256) NOT NULL, 
  arguments        varchar(256) NOT NULL, 
  terminal_name    varchar(40) NOT NULL, 
  terminal_entry   varchar(20) NOT NULL, 
  terminal_model   varchar(40) NOT NULL, 
  terminal_version varchar(20) NOT NULL, 
  extra_info       varchar(256) NOT NULL, 
  create_time      int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE sys_log IS '日志记录';
COMMENT ON COLUMN sys_log.id IS '编号';
COMMENT ON COLUMN sys_log.terminal_entry IS '终端入口';
COMMENT ON COLUMN sys_log.user_id IS '用户编号';
COMMENT ON COLUMN sys_log.username IS '用户名';
COMMENT ON COLUMN sys_log.log_level IS '日志级别, 1:信息  2: 警告  3: 错误 4: 其他';
COMMENT ON COLUMN sys_log.arguments IS '参数';
COMMENT ON COLUMN sys_log.terminal_model IS '终端设备型号';
COMMENT ON COLUMN sys_log.terminal_name IS '终端名称';
COMMENT ON COLUMN sys_log.terminal_version IS '终端应用版本';
COMMENT ON COLUMN sys_log.extra_info IS '额外信息';
COMMENT ON COLUMN sys_log.create_time IS '创建时间';



			
-- 20241110 系统应用版本

				
DROP TABLE IF EXISTS "public".sys_app_distribution CASCADE;

CREATE TABLE "public".sys_app_distribution (
  id              bigserial NOT NULL, 
  app_name        varchar(20) NOT NULL, 
  app_icon        varchar(128) NOT NULL, 
  app_desc        varchar(128) NOT NULL, 
  update_mode     int4 NOT NULL, 
  distribute_url  varchar(256) NOT NULL, 
  distribute_name varchar(10) NOT NULL, 
  stable_version  varchar(20) NOT NULL, 
  stable_down_url varchar(256) NOT NULL, 
  beta_version    varchar(20) NOT NULL, 
  beta_down_url   varchar(256) NOT NULL, 
  url_scheme      varchar(40) DEFAULT '' NOT NULL, 
  create_time     int8 NOT NULL, 
  update_time     int8 NOT NULL, 
  CONSTRAINT app_prod_pkey 
    PRIMARY KEY (id));
COMMENT ON TABLE "public".sys_app_distribution IS 'APP产品';
COMMENT ON COLUMN "public".sys_app_distribution.id IS '产品编号';
COMMENT ON COLUMN "public".sys_app_distribution.app_name IS '英文应用名称,如:mall';
COMMENT ON COLUMN "public".sys_app_distribution.app_icon IS 'APP图标';
COMMENT ON COLUMN "public".sys_app_distribution.app_desc IS '产品描述';
COMMENT ON COLUMN "public".sys_app_distribution.update_mode IS '更新模式, 1:包更新  2: 更新通知';
COMMENT ON COLUMN "public".sys_app_distribution.distribute_url IS '分发下载页面地址';
COMMENT ON COLUMN "public".sys_app_distribution.distribute_name IS '分发名称, 如:商城';
COMMENT ON COLUMN "public".sys_app_distribution.stable_version IS '最新的版本';
COMMENT ON COLUMN "public".sys_app_distribution.stable_down_url IS '正式版文件地址';
COMMENT ON COLUMN "public".sys_app_distribution.beta_version IS '测试版下载地址';
COMMENT ON COLUMN "public".sys_app_distribution.beta_down_url IS '内测版文件地址';
COMMENT ON COLUMN "public".sys_app_distribution.url_scheme IS '应用URL协议';
COMMENT ON COLUMN "public".sys_app_distribution.create_time IS '创建时间';
COMMENT ON COLUMN "public".sys_app_distribution.update_time IS '更新时间';




CREATE TABLE IF NOT EXISTS public.sys_app_version
(
    id bigserial NOT NULL,
    distribution_id bigint,
    version character varying(20) COLLATE pg_catalog."default",
    version_code integer,
    terminal_os character varying(10) COLLATE pg_catalog."default",
    terminal_channel character varying(10),
    start_time bigint,
    update_mode integer,
    update_content character varying(128) COLLATE pg_catalog."default",
    package_url character varying(128) COLLATE pg_catalog."default",
    is_force integer,
    is_notified integer,
    create_time bigint,
    update_time bigint,
    CONSTRAINT sys_app_version_pkey PRIMARY KEY (id)
);

COMMENT ON TABLE public.sys_app_version
    IS '应用版本';

COMMENT ON COLUMN public.sys_app_version.id
    IS '编号';
COMMENT ON COLUMN public.sys_app_version.distribution_id
    IS '分发应用编号';

COMMENT ON COLUMN public.sys_app_version.terminal_os
    IS '终端系统, 如: android / ios';

COMMENT ON COLUMN public.sys_app_version.terminal_channel
    IS '更新通道, beta: 测试版 nightly:每夜版 stable: 正式版本';

COMMENT ON COLUMN public.sys_app_version.version
    IS '版本号';

COMMENT ON COLUMN public.sys_app_version.version_code
    IS '版本数字代号';

COMMENT ON COLUMN public.sys_app_version.start_time
    IS '开始时间';

COMMENT ON COLUMN public.sys_app_version.update_content
    IS '更新内容';

COMMENT ON COLUMN public.sys_app_version.package_url
    IS '下载包地址';

COMMENT ON COLUMN public.sys_app_version.is_force
    IS '是否强制更新';

COMMENT ON COLUMN public.sys_app_version.is_notified
    IS '是否已完成通知,完成后结束更新';

COMMENT ON COLUMN public.sys_app_version.update_mode
    IS '更新模式, 1:包更新  2: 更新通知';

COMMENT ON COLUMN public.sys_app_version.create_time
    IS '创建时间';

COMMENT ON COLUMN public.sys_app_version.update_time
    IS '更新时间';

-- 20241120 最后在线时间
ALTER TABLE mch_staff 
  ADD COLUMN last_online_time int8 DEFAULT 0 NOT NULL;
COMMENT ON COLUMN mch_staff.last_online_time IS '最后在线时间';

-- 20241122 修改应用私钥key
update registry set key='sys_private_key' WHERE key='sys_jwt_secret';

ALTER TABLE "public".mm_member 
  ALTER COLUMN password SET DATA TYPE varchar(64);
ALTER TABLE "public".mm_member 
  ALTER COLUMN trade_pwd SET DATA TYPE varchar(64);
ALTER TABLE "public".mm_member 
  ALTER COLUMN profile_photo SET DATA TYPE varchar(180);

ALTER TABLE "public".mch_merchant 
  ALTER COLUMN password SET DATA TYPE varchar(64);

ALTER TABLE "public".rbac_user 
  ALTER COLUMN password SET DATA TYPE varchar(64);

-- 20241123 商户认证
ALTER TABLE "public".mch_authenticate 
  ADD COLUMN contact_name varchar(10) DEFAULT '' NOT NULL;
ALTER TABLE "public".mch_authenticate 
  ADD COLUMN contact_phone varchar(11) DEFAULT '' NOT NULL;
COMMENT ON COLUMN "public".mch_authenticate.contact_name IS '联系人姓名';
COMMENT ON COLUMN "public".mch_authenticate.contact_phone IS '联系人电话';

 -- 20241126 商户全称
ALTER TABLE "public".mch_merchant 
  ADD COLUMN full_name varchar(128) DEFAULT '' NOT NULL;
COMMENT ON COLUMN "public".mch_merchant.full_name IS '商户全称';

update "public".mch_merchant set full_name=(select org_name from "public".mch_authenticate where mch_id=id);
update "public".mch_merchant set full_name=mch_name where full_name='';

ALTER TABLE "public".mm_member 
  ADD COLUMN country_code varchar(20) DEFAULT 'CN' NOT NULL;

ALTER TABLE "public".mm_member 
  ADD COLUMN create_time int8 DEFAULT 0 NOT NULL;

COMMENT ON COLUMN "public".mm_member.country_code IS '国家代码';
COMMENT ON COLUMN "public".mm_member.create_time IS '注册时间';

-- 更新注册时间
update "public".mm_member set create_time=reg_time;


DROP TABLE IF EXISTS mm_extra_field CASCADE;
CREATE TABLE mm_extra_field (
  id                   BIGSERIAL NOT NULL, 
  member_id            int8 NOT NULL, 
  region_code          int8 DEFAULT 0 NOT NULL, 
  exp                  int8 DEFAULT '0'::bigint NOT NULL, 
  reg_ip               varchar(20) NOT NULL, 
  reg_from             varchar(20) NOT NULL, 
  reg_time             int8 NOT NULL, 
  check_code           varchar(8) NOT NULL, 
  check_expires        int4 NOT NULL, 
  personal_service_uid int4 NOT NULL, 
  login_time           int4 NOT NULL, 
  last_login_time      int4 DEFAULT 0 NOT NULL, 
  update_time          int8 DEFAULT 0 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE mm_extra_field IS '会员额外属性';
COMMENT ON COLUMN mm_extra_field.region_code IS '城市编码';
COMMENT ON COLUMN mm_extra_field.exp IS '经验值';
COMMENT ON COLUMN mm_extra_field.reg_ip IS '注册IP';
COMMENT ON COLUMN mm_extra_field.reg_from IS '注册来源';
COMMENT ON COLUMN mm_extra_field.reg_time IS '注册时间';
COMMENT ON COLUMN mm_extra_field.check_code IS '校验码';
COMMENT ON COLUMN mm_extra_field.check_expires IS '校验码过期时间';
COMMENT ON COLUMN mm_extra_field.personal_service_uid IS '私人客服人员编号';
COMMENT ON COLUMN mm_extra_field.login_time IS '登录时间';
COMMENT ON COLUMN mm_extra_field.last_login_time IS '最后登录时间';
COMMENT ON COLUMN mm_extra_field.update_time IS '更新时间';


-- 插入会员扩展信息
insert into mm_extra_field (member_id, exp, reg_ip,region_code, reg_from, 
reg_time, check_code, check_expires, personal_service_uid,
 login_time, last_login_time, update_time)
select id, exp, reg_ip,0, reg_from, reg_time, check_code, 
check_expires, 0, login_time, 
last_login_time, update_time from mm_member
WHERE id < 10;


ALTER TABLE "public".mm_member 
  alter column user_code set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column salt set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column profile_photo set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column phone set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column email set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column nickname set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column real_name set default ''::character varying;

ALTER TABLE "public".mm_member 
  DROP COLUMN reg_ip;
ALTER TABLE "public".mm_member 
  DROP COLUMN reg_from;
ALTER TABLE "public".mm_member 
  DROP COLUMN reg_time;
ALTER TABLE "public".mm_member 
  DROP COLUMN check_code;
ALTER TABLE "public".mm_member 
  DROP COLUMN check_expires;
ALTER TABLE "public".mm_member 
  DROP COLUMN login_time;
ALTER TABLE "public".mm_member 
  DROP COLUMN last_login_time;
ALTER TABLE "public".mm_member 
  DROP COLUMN exp;
