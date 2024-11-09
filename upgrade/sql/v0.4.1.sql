
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


