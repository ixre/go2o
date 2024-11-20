
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
