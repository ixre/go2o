CREATE TABLE "public".perm_res (
   id             bigserial NOT NULL,
   name           varchar(20) NOT NULL,
   res_type       int2 NOT NULL,
   pid            int8 NOT NULL,
   depth          int2 NOT NULL,
   "key"          varchar(120) NOT NULL,
   path           varchar(256) NOT NULL,
   icon           varchar(120) NOT NULL,
   permission     varchar(120) NOT NULL,
   sort_num       int4 NOT NULL,
   is_external    int2 NOT NULL,
   is_hidden      int2 DEFAULT 0 NOT NULL,
   create_time    int8 NOT NULL,
   component_path varchar(120) NOT NULL,
   cache_         varchar(20) DEFAULT ''::character varying NOT NULL,
   CONSTRAINT perm_res_pkey
       PRIMARY KEY (id));
COMMENT ON COLUMN "public".perm_res.id IS '资源ID';
COMMENT ON COLUMN "public".perm_res.name IS '资源名称';
COMMENT ON COLUMN "public".perm_res.res_type IS '资源类型, 0: 资源　 1: 菜单  2:　 按钮';
COMMENT ON COLUMN "public".perm_res.pid IS '上级菜单ID';
COMMENT ON COLUMN "public".perm_res.depth IS '深度/层级';
COMMENT ON COLUMN "public".perm_res."key" IS '资源键';
COMMENT ON COLUMN "public".perm_res.path IS '资源路径';
COMMENT ON COLUMN "public".perm_res.icon IS '图标';
COMMENT ON COLUMN "public".perm_res.permission IS '权限,多个值用|分隔';
COMMENT ON COLUMN "public".perm_res.sort_num IS '排序';
COMMENT ON COLUMN "public".perm_res.is_external IS '是否外部';
COMMENT ON COLUMN "public".perm_res.is_hidden IS '是否隐藏';
COMMENT ON COLUMN "public".perm_res.create_time IS '创建日期';
COMMENT ON COLUMN "public".perm_res.component_path IS '组件路径';
COMMENT ON COLUMN "public".perm_res.cache_ IS '缓存';


CREATE TABLE "public".perm_dict (
                                    id          bigserial NOT NULL,
                                    name        varchar(255) NOT NULL,
                                    remark      varchar(255) DEFAULT 'NULL::character varying' NOT NULL,
                                    create_time int8 NOT NULL,
                                    CONSTRAINT perm_dict_pkey
                                        PRIMARY KEY (id));
COMMENT ON TABLE "public".perm_dict IS '数据字典';
COMMENT ON COLUMN "public".perm_dict.name IS '字典名称';
COMMENT ON COLUMN "public".perm_dict.remark IS '描述';
COMMENT ON COLUMN "public".perm_dict.create_time IS '创建日期';
CREATE TABLE "public".perm_dict_detail (
                                           id          bigserial NOT NULL,
                                           label       varchar(255) NOT NULL,
                                           value       varchar(255) NOT NULL,
                                           sort        varchar(255) DEFAULT 'NULL::character varying' NOT NULL,
                                           dict_id     int8 NOT NULL,
                                           create_time int8 NOT NULL,
                                           CONSTRAINT perm_dict_detail_pkey
                                               PRIMARY KEY (id));
COMMENT ON TABLE "public".perm_dict_detail IS '数据字典详情';
COMMENT ON COLUMN "public".perm_dict_detail.label IS '字典标签';
COMMENT ON COLUMN "public".perm_dict_detail.value IS '字典值';
COMMENT ON COLUMN "public".perm_dict_detail.sort IS '排序';
COMMENT ON COLUMN "public".perm_dict_detail.dict_id IS '字典id';
COMMENT ON COLUMN "public".perm_dict_detail.create_time IS '创建日期';
CREATE TABLE "public".perm_dept (
                                    id          bigserial NOT NULL,
                                    name        varchar(40) NOT NULL,
                                    pid         int8 NOT NULL,
                                    enabled     int2 NOT NULL,
                                    create_time int8 NOT NULL,
                                    CONSTRAINT perm_dept_pkey
                                        PRIMARY KEY (id));
COMMENT ON TABLE "public".perm_dept IS '部门';
COMMENT ON COLUMN "public".perm_dept.id IS 'ID';
COMMENT ON COLUMN "public".perm_dept.name IS '名称';
COMMENT ON COLUMN "public".perm_dept.pid IS '上级部门';
COMMENT ON COLUMN "public".perm_dept.enabled IS '状态';
COMMENT ON COLUMN "public".perm_dept.create_time IS '创建日期';
CREATE TABLE "public".perm_job (
                                   id          bigserial NOT NULL,
                                   name        varchar(40) NOT NULL,
                                   enabled     int2 NOT NULL,
                                   sort        int4 NOT NULL,
                                   dept_id     int8 NOT NULL,
                                   create_time int8 NOT NULL,
                                   CONSTRAINT perm_job_pkey
                                       PRIMARY KEY (id));
COMMENT ON TABLE "public".perm_job IS '岗位';
COMMENT ON COLUMN "public".perm_job.id IS 'ID';
COMMENT ON COLUMN "public".perm_job.name IS '岗位名称';
COMMENT ON COLUMN "public".perm_job.enabled IS '岗位状态';
COMMENT ON COLUMN "public".perm_job.sort IS '岗位排序';
COMMENT ON COLUMN "public".perm_job.dept_id IS '部门ID';
COMMENT ON COLUMN "public".perm_job.create_time IS '创建日期';
CREATE TABLE "public".perm_role (
                                    id          bigserial NOT NULL,
                                    name        varchar(40) NOT NULL,
                                    level       int4 NOT NULL,
                                    data_scope  varchar(255) DEFAULT 'NULL::character varying' NOT NULL,
                                    permission  varchar(255) DEFAULT 'NULL::character varying' NOT NULL,
                                    remark      varchar(120) DEFAULT 'NULL::character varying' NOT NULL,
                                    create_time int8 NOT NULL,
                                    CONSTRAINT perm_role_pkey
                                        PRIMARY KEY (id));
COMMENT ON TABLE "public".perm_role IS '角色';
COMMENT ON COLUMN "public".perm_role.id IS 'ID';
COMMENT ON COLUMN "public".perm_role.name IS '名称';
COMMENT ON COLUMN "public".perm_role.level IS '角色级别';
COMMENT ON COLUMN "public".perm_role.data_scope IS '数据权限';
COMMENT ON COLUMN "public".perm_role.permission IS '功能权限';
COMMENT ON COLUMN "public".perm_role.remark IS '备注';
COMMENT ON COLUMN "public".perm_role.create_time IS '创建日期';
CREATE TABLE "public".perm_role_dept (
                                         id       BIGSERIAL NOT NULL,
                                         role_id int8 NOT NULL,
                                         dept_id int8 NOT NULL,
                                         CONSTRAINT perm_role_dept_pkey
                                             PRIMARY KEY (id));
COMMENT ON TABLE "public".perm_role_dept IS '角色部门关联';
COMMENT ON COLUMN "public".perm_role_dept.id IS '编号';
COMMENT ON COLUMN "public".perm_role_dept.role_id IS '角色编号';
COMMENT ON COLUMN "public".perm_role_dept.dept_id IS '部门编号';
CREATE TABLE "public".perm_role_res (
                                        id       BIGSERIAL NOT NULL,
                                        res_id  int8 NOT NULL,
                                        role_id int8 NOT NULL,
                                        CONSTRAINT perm_role_menu_pkey
                                            PRIMARY KEY (id));
COMMENT ON TABLE "public".perm_role_res IS '角色菜单关联';
COMMENT ON COLUMN "public".perm_role_res.id IS '编号';
COMMENT ON COLUMN "public".perm_role_res.res_id IS '菜单ID';
COMMENT ON COLUMN "public".perm_role_res.role_id IS '角色ID';
CREATE TABLE "public".perm_user (
                                    id          bigserial NOT NULL,
                                    usr      varchar(20) DEFAULT 'NULL::character varying' NOT NULL,
                                    pwd         varchar(40) DEFAULT 'NULL::character varying' NOT NULL,
                                    salt        varchar(10) NOT NULL,
                                    flag        int4 NOT NULL,
                                    avatar      varchar(256) NOT NULL,
                                    nick_name   varchar(20) DEFAULT 'NULL::character varying' NOT NULL,
                                    gender         varchar(20) DEFAULT 'NULL::character varying' NOT NULL,
                                    email       varchar(64) DEFAULT 'NULL::character varying' NOT NULL,
                                    phone       varchar(11) DEFAULT 'NULL::character varying' NOT NULL,
                                    dept_id     int8 NOT NULL,
                                    job_id      int8 NOT NULL,
                                    enabled     int2 NOT NULL,
                                    last_login  int8 NOT NULL,
                                    create_time int8 NOT NULL,
                                    CONSTRAINT perm_user_pkey
                                        PRIMARY KEY (id));
COMMENT ON TABLE "public".perm_user IS '系统用户';
COMMENT ON COLUMN "public".perm_user.id IS 'ID';
COMMENT ON COLUMN "public".perm_user.usr IS '用户名';
COMMENT ON COLUMN "public".perm_user.pwd IS '密码';
COMMENT ON COLUMN "public".perm_user.salt IS '加密盐';
COMMENT ON COLUMN "public".perm_user.flag IS '标志';
COMMENT ON COLUMN "public".perm_user.avatar IS '头像';
COMMENT ON COLUMN "public".perm_user.nick_name IS '姓名';
COMMENT ON COLUMN "public".perm_user.gender IS '性别';
COMMENT ON COLUMN "public".perm_user.email IS '邮箱';
COMMENT ON COLUMN "public".perm_user.phone IS '手机号码';
COMMENT ON COLUMN "public".perm_user.dept_id IS '部门编号';
COMMENT ON COLUMN "public".perm_user.job_id IS '岗位编号';
COMMENT ON COLUMN "public".perm_user.enabled IS '状态：1启用、0禁用';
COMMENT ON COLUMN "public".perm_user.last_login IS '最后登录的日期';
COMMENT ON COLUMN "public".perm_user.create_time IS '创建日期';

CREATE TABLE "public".perm_user_role (
                                         id       BIGSERIAL NOT NULL,
                                         user_id int8 NOT NULL,
                                         role_id int8 NOT NULL,
                                         CONSTRAINT perm_user_role_pkey
                                             PRIMARY KEY (id));
COMMENT ON TABLE "public".perm_user_role IS '用户角色关联';
COMMENT ON COLUMN "public".perm_user_role.id IS '编号';
COMMENT ON COLUMN "public".perm_user_role.user_id IS '用户ID';
COMMENT ON COLUMN "public".perm_user_role.role_id IS '角色ID';
CREATE INDEX perm_dict_detail_dict_id_idx
    ON "public".perm_dict_detail (dict_id);

