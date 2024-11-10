
-- 2024-11-09 日志 

CREATE TABLE sys_log_app (
  id        BIGSERIAL NOT NULL, 
  name      varchar(10) NOT NULL, 
  log_level int4 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE sys_log_app IS '日志应用';
COMMENT ON COLUMN sys_log_app.log_level IS '日志级别, 0: 不记录  1: 信息 2:警告  3: 错误 4:全部';
CREATE TABLE sys_log_list (
  id               BIGSERIAL NOT NULL, 
  app_id           int4 NOT NULL, 
  user_id          int4 NOT NULL, 
  username         varchar(20) NOT NULL, 
  log_level        int4 NOT NULL, 
  message          varchar(128) NOT NULL, 
  arguments        varchar(256) NOT NULL, 
  terminal_model   varchar(20) NOT NULL, 
  terminal_name    varchar(20) NOT NULL, 
  terminal_version varchar(20) NOT NULL, 
  extra_info       varchar(256) NOT NULL, 
  create_time      int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE sys_log_list IS '日志记录';
COMMENT ON COLUMN sys_log_list.id IS '编号';
COMMENT ON COLUMN sys_log_list.app_id IS '应用编号';
COMMENT ON COLUMN sys_log_list.user_id IS '用户编号';
COMMENT ON COLUMN sys_log_list.username IS '用户名';
COMMENT ON COLUMN sys_log_list.log_level IS '日志级别, 1:信息  2: 警告  3: 错误 4: 其他';
COMMENT ON COLUMN sys_log_list.arguments IS '参数';
COMMENT ON COLUMN sys_log_list.terminal_model IS '终端设备型号';
COMMENT ON COLUMN sys_log_list.terminal_name IS '终端名称';
COMMENT ON COLUMN sys_log_list.terminal_version IS '终端应用版本';
COMMENT ON COLUMN sys_log_list.extra_info IS '额外信息';
COMMENT ON COLUMN sys_log_list.create_time IS '创建时间';


-- 20241110 系统应用版本
CREATE TABLE IF NOT EXISTS public.sys_app_version
(
    id bigint NOT NULL,
    version character varying(20) COLLATE pg_catalog."default",
    version_code integer,
    terminal_os character varying(10) COLLATE pg_catalog."default",
    is_stable integer,
    start_time bigint,
    update_mode integer,
    update_content character varying(128) COLLATE pg_catalog."default",
    package_url character varying(128) COLLATE pg_catalog."default",
    is_force integer,
    is_notified integer,
    create_time bigint,
    update_time bigint,
    CONSTRAINT sys_app_version_pkey PRIMARY KEY (id)
)

COMMENT ON TABLE public.sys_app_version
    IS '应用版本';

COMMENT ON COLUMN public.sys_app_version.id
    IS '编号';

COMMENT ON COLUMN public.sys_app_version.terminal_os
    IS '终端系统, 如: android / ios';

COMMENT ON COLUMN public.sys_app_version.is_stable
    IS '是否正式版本, 0: 测试版  1: 正式版本';

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
