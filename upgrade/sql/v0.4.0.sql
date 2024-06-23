-- 子站,聊天,投诉等
DROP TABLE IF EXISTS chat_conversation;
CREATE TABLE chat_conversation (
  id             BIGSERIAL NOT NULL, 
  "key"          varchar(20) NOT NULL, 
  start_user_id  int8 NOT NULL, 
  join_user_id   int8 NOT NULL, 
  flag           int4 NOT NULL, 
  greet_word     varchar(20) NOT NULL, 
  create_time    int8 NOT NULL, 
  update_time    int8 NOT NULL, 
  last_chat_time int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE chat_conversation IS '聊天会话';
COMMENT ON COLUMN chat_conversation.id IS '编号';
COMMENT ON COLUMN chat_conversation."key" IS '编码';
COMMENT ON COLUMN chat_conversation.start_user_id IS '会话创建人';
COMMENT ON COLUMN chat_conversation.join_user_id IS '会话参与人';
COMMENT ON COLUMN chat_conversation.flag IS '预留标志';
COMMENT ON COLUMN chat_conversation.greet_word IS '打招呼内容';
COMMENT ON COLUMN chat_conversation.create_time IS '创建时间';
COMMENT ON COLUMN chat_conversation.update_time IS '更新时间';
COMMENT ON COLUMN chat_conversation.last_chat_time IS '最后聊天时间';
DROP TABLE IF EXISTS chat_msg;
CREATE TABLE chat_msg (
  id              BIGSERIAL NOT NULL, 
  conversation_id int8 NOT NULL, 
  msg_type        int4 NOT NULL, 
  msg_data        varchar(120) NOT NULL, 
  msg_content     varchar(255) NOT NULL, 
  starter_msg     int4 NOT NULL, 
  is_revert       int4 NOT NULL, 
  expires_time    int8 NOT NULL, 
  create_time     int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE chat_msg IS '消息消息';
COMMENT ON COLUMN chat_msg.id IS '编号';
COMMENT ON COLUMN chat_msg.conversation_id IS '会话编号';
COMMENT ON COLUMN chat_msg.msg_type IS '消息类型, 1: 文本  2: 图片  3: 表情  4: 文件  5:语音  6:位置  7:语音  8:红包  9:名片  11: 委托申请';
COMMENT ON COLUMN chat_msg.msg_data IS '消息数据';
COMMENT ON COLUMN chat_msg.msg_content IS '消息内容';
COMMENT ON COLUMN chat_msg.starter_msg IS '是否为发起人的消息, 0:否 1:是';
COMMENT ON COLUMN chat_msg.is_revert IS '是否撤回 0:否 1:是';
COMMENT ON COLUMN chat_msg.expires_time IS '过期时间';
COMMENT ON COLUMN chat_msg.create_time IS '创建时间';
DROP TABLE IF EXISTS m_block_list;
CREATE TABLE m_block_list (
  id              BIGSERIAL NOT NULL, 
  member_id       int8 NOT NULL, 
  block_member_id int8 NOT NULL, 
  block_flag      int4 NOT NULL, 
  remark          varchar(20) NOT NULL, 
  create_time     int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE m_block_list IS '会员拉黑列表';
COMMENT ON COLUMN m_block_list.id IS '编号';
COMMENT ON COLUMN m_block_list.member_id IS '会员编号';
COMMENT ON COLUMN m_block_list.block_member_id IS '拉黑会员编号';
COMMENT ON COLUMN m_block_list.block_flag IS '拉黑标志，1: 屏蔽  2: 拉黑';
COMMENT ON COLUMN m_block_list.remark IS '备注';
COMMENT ON COLUMN m_block_list.create_time IS '拉黑时间';
DROP TABLE IF EXISTS m_complain_case;
CREATE TABLE m_complain_case (
  id               BIGSERIAL NOT NULL, 
  member_id        int8 NOT NULL, 
  complain_type    int4 NOT NULL, 
  order_id         int8 NOT NULL, 
  mch_id           int8 NOT NULL, 
  target_member_id int8 NOT NULL, 
  complain_desc    varchar(255) NOT NULL, 
  hope_desc        varchar(120) NOT NULL, 
  first_pic        varchar(80) NOT NULL, 
  pic_list         varchar(350) NOT NULL, 
  is_resolved      int4 NOT NULL, 
  is_closed        int4 NOT NULL, 
  status           int4 NOT NULL, 
  service_agent_id int8 NOT NULL, 
  service_rank     int4 NOT NULL, 
  service_apprise  varchar(120) NOT NULL, 
  create_time      int8 NOT NULL, 
  update_time      int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON COLUMN m_complain_case.id IS '编号';
COMMENT ON COLUMN m_complain_case.member_id IS '会员编号';
COMMENT ON COLUMN m_complain_case.complain_type IS '投诉类型:  1:常规   11: 咨询服务';
COMMENT ON COLUMN m_complain_case.order_id IS '订单号';
COMMENT ON COLUMN m_complain_case.mch_id IS '商户编号';
COMMENT ON COLUMN m_complain_case.target_member_id IS '投诉目标会员';
COMMENT ON COLUMN m_complain_case.complain_desc IS '投诉内容';
COMMENT ON COLUMN m_complain_case.hope_desc IS '诉求描述';
COMMENT ON COLUMN m_complain_case.first_pic IS '图片';
COMMENT ON COLUMN m_complain_case.pic_list IS '图片列表';
COMMENT ON COLUMN m_complain_case.is_resolved IS '是否已解决 0:否 1:是';
COMMENT ON COLUMN m_complain_case.is_closed IS '是否用户关闭 0:否 1:是';
COMMENT ON COLUMN m_complain_case.status IS '状态,1:待处理 2:处理中 3:已完结';
COMMENT ON COLUMN m_complain_case.service_agent_id IS '客服编号';
COMMENT ON COLUMN m_complain_case.service_rank IS '服务评分';
COMMENT ON COLUMN m_complain_case.service_apprise IS '服务评价';
COMMENT ON COLUMN m_complain_case.create_time IS '创建时间';
COMMENT ON COLUMN m_complain_case.update_time IS '更新时间';
DROP TABLE IF EXISTS m_complain_details;
CREATE TABLE m_complain_details (
  id          BIGSERIAL NOT NULL, 
  case_id     int8 NOT NULL, 
  sender_type int4 NOT NULL, 
  content     varchar(255) NOT NULL, 
  is_revert   int4 NOT NULL, 
  create_time int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE m_complain_details IS '投诉详情';
COMMENT ON COLUMN m_complain_details.id IS '编号';
COMMENT ON COLUMN m_complain_details.case_id IS '案件编号';
COMMENT ON COLUMN m_complain_details.sender_type IS '发送类型: 1:发起人  2: 投诉对象  3: 平台客服';
COMMENT ON COLUMN m_complain_details.is_revert IS '是否撤回 0:否 1:是';



DROP TABLE IF EXISTS sys_general_option;
CREATE TABLE sys_general_option (
  id          BIGSERIAL NOT NULL, 
  type        varchar(20) NOT NULL, 
  pid         int8 NOT NULL, 
  name        varchar(20) NOT NULL, 
  value       varchar(20) NOT NULL, 
  sort_num    int4 NOT NULL, 
  enabled     int4 NOT NULL, 
  create_time int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE sys_general_option IS '系统通用选项(用于存储分类和选项等数据)';
COMMENT ON COLUMN sys_general_option.id IS '编号';
COMMENT ON COLUMN sys_general_option.type IS '类型';
COMMENT ON COLUMN sys_general_option.pid IS '上级编号';
COMMENT ON COLUMN sys_general_option.name IS '名称';
COMMENT ON COLUMN sys_general_option.value IS '值';
COMMENT ON COLUMN sys_general_option.sort_num IS '排列序号';
COMMENT ON COLUMN sys_general_option.enabled IS '是否启用';
COMMENT ON COLUMN sys_general_option.create_time IS '创建时间';
DROP TABLE IF EXISTS sys_sub_station;
CREATE TABLE sys_sub_station (
  id          SERIAL NOT NULL, 
  city_code   int4 NOT NULL, 
  status      int4 NOT NULL, 
  create_time int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE sys_sub_station IS '地区子站';
COMMENT ON COLUMN sys_sub_station.id IS '编号';
COMMENT ON COLUMN sys_sub_station.city_code IS '城市代码';
COMMENT ON COLUMN sys_sub_station.status IS '状态: 0: 待开通  1: 已开通  2: 已关闭';
COMMENT ON COLUMN sys_sub_station.create_time IS '创建时间';


-- 商户
DROP TABLE IF EXISTS mch_merchant;
CREATE TABLE "public".mch_merchant (
  id              BIGSERIAL NOT NULL, 
  member_id       int8 NOT NULL, 
  username        varchar(30) NOT NULL, 
  password        varchar(45) NOT NULL, 
  mail_addr       varchar(30) NOT NULL, 
  salt            varchar(10) DEFAULT ''::character varying NOT NULL, 
  mch_name        varchar(20) NOT NULL, 
  is_self         int2 NOT NULL, 
  flag            int4 NOT NULL, 
  level           int4 NOT NULL, 
  province        int4 NOT NULL, 
  city            int4 NOT NULL, 
  district        int4 NOT NULL, 
  address         varchar(120) NOT NULL, 
  logo            varchar(120) NOT NULL, 
  tel             varchar(20) NOT NULL, 
  status          int2 NOT NULL, 
  expires_time    int4 NOT NULL, 
  last_login_time int4 NOT NULL, 
  create_time     int4 NOT NULL, 
  CONSTRAINT mch_merchant_pkey 
    PRIMARY KEY (id));
COMMENT ON TABLE "public".mch_merchant IS '商户';
COMMENT ON COLUMN "public".mch_merchant.id IS '编号';
COMMENT ON COLUMN "public".mch_merchant.member_id IS '会员编号';
COMMENT ON COLUMN "public".mch_merchant.username IS '登录用户';
COMMENT ON COLUMN "public".mch_merchant.password IS '登录密码';
COMMENT ON COLUMN "public".mch_merchant.mail_addr IS '邮箱地址';
COMMENT ON COLUMN "public".mch_merchant.salt IS '加密盐';
COMMENT ON COLUMN "public".mch_merchant.mch_name IS '名称';
COMMENT ON COLUMN "public".mch_merchant.is_self IS '是否自营';
COMMENT ON COLUMN "public".mch_merchant.flag IS '标志';
COMMENT ON COLUMN "public".mch_merchant.level IS '商户等级';
COMMENT ON COLUMN "public".mch_merchant.province IS '所在省';
COMMENT ON COLUMN "public".mch_merchant.city IS '所在市';
COMMENT ON COLUMN "public".mch_merchant.district IS '所在区';
COMMENT ON COLUMN "public".mch_merchant.address IS '公司地址';
COMMENT ON COLUMN "public".mch_merchant.logo IS '标志';
COMMENT ON COLUMN "public".mch_merchant.tel IS '公司电话';
COMMENT ON COLUMN "public".mch_merchant.status IS '状态: 0:待审核 1:已开通  2:停用  3: 关闭';
COMMENT ON COLUMN "public".mch_merchant.expires_time IS '过期时间';
COMMENT ON COLUMN "public".mch_merchant.last_login_time IS '最后登录时间';
COMMENT ON COLUMN "public".mch_merchant.create_time IS '创建时间';



DROP TABLE IF EXISTS mch_authenticate;
CREATE TABLE "public".mch_authenticate (
  id                BIGSERIAL NOT NULL, 
  mch_id            int4 NOT NULL, 
  org_name          varchar(45) NOT NULL, 
  org_no            varchar(45) NOT NULL, 
  org_pic           varchar(120) NOT NULL, 
  work_city         int4 NOT NULL, 
  qualification_pic varchar(120) NOT NULL, 
  person_id         varchar(20) NOT NULL, 
  person_name       varchar(10) NOT NULL, 
  person_pic        varchar(120) NOT NULL, 
  authority_pic     varchar(120) NOT NULL, 
  bank_name         varchar(20) NOT NULL, 
  bank_account      varchar(20) NOT NULL, 
  bank_no           varchar(20) NOT NULL, 
  extra_data        varchar(512) NOT NULL, 
  review_time       int4 NOT NULL, 
  review_status     int4 NOT NULL, 
  review_remark     varchar(45) NOT NULL, 
  update_time       int8 NOT NULL, 
  CONSTRAINT mch_authenticate_pkey 
    PRIMARY KEY (id));
COMMENT ON TABLE "public".mch_authenticate IS '商户认证信息';
COMMENT ON COLUMN "public".mch_authenticate.mch_id IS '商户编号';
COMMENT ON COLUMN "public".mch_authenticate.org_name IS '公司名称';
COMMENT ON COLUMN "public".mch_authenticate.org_no IS '营业执照编号';
COMMENT ON COLUMN "public".mch_authenticate.org_pic IS '营业执照照片';
COMMENT ON COLUMN "public".mch_authenticate.work_city IS '办公地';
COMMENT ON COLUMN "public".mch_authenticate.qualification_pic IS '资质图片';
COMMENT ON COLUMN "public".mch_authenticate.person_id IS '法人身份证号';
COMMENT ON COLUMN "public".mch_authenticate.person_name IS '法人姓名';
COMMENT ON COLUMN "public".mch_authenticate.person_pic IS '法人身份证照片';
COMMENT ON COLUMN "public".mch_authenticate.authority_pic IS '授权书';
COMMENT ON COLUMN "public".mch_authenticate.bank_name IS '开户银行';
COMMENT ON COLUMN "public".mch_authenticate.bank_account IS '开户户名';
COMMENT ON COLUMN "public".mch_authenticate.bank_no IS '开户账号';
COMMENT ON COLUMN "public".mch_authenticate.extra_data IS '扩展数据';
COMMENT ON COLUMN "public".mch_authenticate.review_time IS '审核时间';
COMMENT ON COLUMN "public".mch_authenticate.review_status IS '审核状态';
COMMENT ON COLUMN "public".mch_authenticate.review_remark IS '审核备注';
COMMENT ON COLUMN "public".mch_authenticate.update_time IS '更新时间';


ALTER TABLE "public".mch_authenticate 
  ADD COLUMN person_phone varchar(11) NOT NULL;
COMMENT ON COLUMN "public".mch_authenticate.person_phone IS '联系人手机';

ALTER TABLE "public".mch_authenticate 
  ADD COLUMN version int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN "public".mch_authenticate.version IS '版本号: 0: 待审核 1: 已审核';

ALTER TABLE "public".mm_member 
  ADD COLUMN role_flag int4 DEFAULT 0 NOT NULL;
ALTER TABLE "public".mm_member 
  alter column user_code set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column portrait set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column phone set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column email set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column nickname set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column real_name set default ''::character varying;
ALTER TABLE "public".mm_member 
  alter column salt set default ''::character varying;
COMMENT ON COLUMN "public".mm_member.role_flag IS '角色标志';


ALTER TABLE sys_sub_station 
  ADD COLUMN letter varchar(2) DEFAULT '' NOT NULL;
ALTER TABLE sys_sub_station 
  ADD COLUMN is_hot int4 DEFAULT 0 NOT NULL;
COMMENT ON COLUMN sys_sub_station.letter IS '首字母';
COMMENT ON COLUMN sys_sub_station.is_hot IS '是否热门';


-- 系统通知模板
CREATE TABLE "public".sys_notify_template (
  id          bigserial NOT NULL, 
  code        varchar(20) NOT NULL, 
  temp_type   int4 NOT NULL, 
  temp_name   varchar(20) NOT NULL, 
  content     varchar(512) NOT NULL, 
  labels      varchar(40) NOT NULL, 
  sp_code     varchar(20) NOT NULL, 
  sp_tid      varchar(56) NOT NULL, 
  create_time int8 NOT NULL, 
  update_time int8 NOT NULL, 
  is_deleted  int4 NOT NULL, 
  CONSTRAINT sys_notify_template_key 
    PRIMARY KEY (id));
COMMENT ON TABLE "public".sys_notify_template IS '系统通知模板';
COMMENT ON COLUMN "public".sys_notify_template.id IS '编号';
COMMENT ON COLUMN "public".sys_notify_template.code IS '模板编号';
COMMENT ON COLUMN "public".sys_notify_template.temp_type IS '模板类型,1:站内信 2:短信 3:邮件';
COMMENT ON COLUMN "public".sys_notify_template.temp_name IS '模板名称';
COMMENT ON COLUMN "public".sys_notify_template.content IS '模板内容';
COMMENT ON COLUMN "public".sys_notify_template.labels IS '模板标签, 多个用,隔开';
COMMENT ON COLUMN "public".sys_notify_template.sp_code IS '短信服务商代码';
COMMENT ON COLUMN "public".sys_notify_template.sp_tid IS '短信服务商模板编号';
COMMENT ON COLUMN "public".sys_notify_template.create_time IS '创建时间';
COMMENT ON COLUMN "public".sys_notify_template.is_deleted IS '是否删除,0:否 1:是';


DROP TABLE IF EXISTS mch_agent;
DROP TABLE IF EXISTS mch_agent_extent;
DROP TABLE IF EXISTS mch_agent_revenue;
DROP TABLE IF EXISTS mch_service_order;

CREATE TABLE mch_service_order (
  id                SERIAL NOT NULL, 
  order_no          int4 NOT NULL, 
  mch_id            int4 NOT NULL, 
  station_id        int4 NOT NULL, 
  staff_id          int4 NOT NULL, 
  member_id         int4 NOT NULL, 
  charge_amount     int4 NOT NULL, 
  is_transformed    int4 NOT NULL, 
  transform_time    int8 NOT NULL, 
  transform_deposit int4 NOT NULL, 
  service_time      int4 NOT NULL, 
  service_rank      int4 NOT NULL, 
  service_apprise   varchar(120) NOT NULL, 
  status            int4 NOT NULL, 
  create_time       int8 NOT NULL, 
  update_time       int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE mch_service_order IS '商户服务单';
COMMENT ON COLUMN mch_service_order.id IS '编号';
COMMENT ON COLUMN mch_service_order.order_no IS '订单号';
COMMENT ON COLUMN mch_service_order.mch_id IS '商户编号';
COMMENT ON COLUMN mch_service_order.station_id IS '站点编号';
COMMENT ON COLUMN mch_service_order.staff_id IS '代理人编号';
COMMENT ON COLUMN mch_service_order.member_id IS '会员编号';
COMMENT ON COLUMN mch_service_order.charge_amount IS '充值金额(服务单)';
COMMENT ON COLUMN mch_service_order.is_transformed IS '是否转化';
COMMENT ON COLUMN mch_service_order.transform_time IS '转化时间';
COMMENT ON COLUMN mch_service_order.transform_deposit IS '定金';
COMMENT ON COLUMN mch_service_order.service_time IS '服务计时(分钟)';
COMMENT ON COLUMN mch_service_order.service_rank IS '服务评分';
COMMENT ON COLUMN mch_service_order.service_apprise IS '服务评价';
COMMENT ON COLUMN mch_service_order.status IS '状态: 1: 待服务  2: 服务中   3: 已结束  4: 已关闭';
COMMENT ON COLUMN mch_service_order.create_time IS '创建时间';
COMMENT ON COLUMN mch_service_order.update_time IS '更新时间';
CREATE TABLE mch_staff (
  id             BIGSERIAL NOT NULL, 
  member_id      int8 NOT NULL, 
  station_id     int4 NOT NULL, 
  mch_id         int8 NOT NULL, 
  flag           int4 NOT NULL, 
  gender         int4 NOT NULL, 
  nickname       varchar(20) NOT NULL, 
  work_status    int4 NOT NULL, 
  grade          int4 NOT NULL, 
  status         int4 NOT NULL, 
  is_certified   int4 NOT NULL, 
  certified_name varchar(10) NOT NULL, 
  premium_level  int4 NOT NULL, 
  create_time    int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE mch_staff IS '商户代理人坐席(员工)';
COMMENT ON COLUMN mch_staff.id IS '编号';
COMMENT ON COLUMN mch_staff.member_id IS '会员编号';
COMMENT ON COLUMN mch_staff.station_id IS '站点编号';
COMMENT ON COLUMN mch_staff.mch_id IS '商户编号';
COMMENT ON COLUMN mch_staff.flag IS '坐席标志';
COMMENT ON COLUMN mch_staff.gender IS '性别: 0: 未知 1:男 2:女';
COMMENT ON COLUMN mch_staff.nickname IS '昵称';
COMMENT ON COLUMN mch_staff.work_status IS '工作状态: 1: 离线 2:在线空闲 3: 工作中';
COMMENT ON COLUMN mch_staff.grade IS '评分';
COMMENT ON COLUMN mch_staff.status IS '状态: 1: 正常  2: 锁定';
COMMENT ON COLUMN mch_staff.is_certified IS '是否认证 0:否 1:是';
COMMENT ON COLUMN mch_staff.certified_name IS '认证姓名';
COMMENT ON COLUMN mch_staff.premium_level IS '高级用户等级';
COMMENT ON COLUMN mch_staff.create_time IS '创建时间';
CREATE TABLE mch_staff_extent (
  id              BIGSERIAL NOT NULL, 
  staff_id        int8 NOT NULL, 
  certified_time  int8 NOT NULL, 
  focus_fields    varchar(20) NOT NULL, 
  unit_price      numeric(6, 2) NOT NULL, 
  work_begin      int8 NOT NULL, 
  work_years      int4 NOT NULL, 
  birthday        int8 NOT NULL, 
  age             int4 NOT NULL, 
  city_code       int4 NOT NULL, 
  introduce       varchar(80) NOT NULL, 
  commission_rate numeric(4, 2) NOT NULL, 
  id_no           varchar(20) NOT NULL, 
  license_pic     varchar(120) NOT NULL, 
  license_no      varchar(20) NOT NULL, 
  update_time     int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE mch_staff_extent IS '商户坐席(员工)扩展表';
COMMENT ON COLUMN mch_staff_extent.id IS '编号';
COMMENT ON COLUMN mch_staff_extent.certified_time IS '认证时间';
COMMENT ON COLUMN mch_staff_extent.focus_fields IS '聚焦领域';
COMMENT ON COLUMN mch_staff_extent.unit_price IS '每小时单价';
COMMENT ON COLUMN mch_staff_extent.work_begin IS '工作起始时间';
COMMENT ON COLUMN mch_staff_extent.work_years IS '工龄';
COMMENT ON COLUMN mch_staff_extent.birthday IS '生日';
COMMENT ON COLUMN mch_staff_extent.age IS '年龄';
COMMENT ON COLUMN mch_staff_extent.city_code IS '所在城市';
COMMENT ON COLUMN mch_staff_extent.introduce IS '个人介绍';
COMMENT ON COLUMN mch_staff_extent.commission_rate IS '提成比例';
COMMENT ON COLUMN mch_staff_extent.id_no IS '身份证号码';
COMMENT ON COLUMN mch_staff_extent.license_pic IS '执业资格图片';
COMMENT ON COLUMN mch_staff_extent.license_no IS '执业资格证编号';
COMMENT ON COLUMN mch_staff_extent.update_time IS '更新时间';
CREATE TABLE mch_staff_revenue (
  id              SERIAL NOT NULL, 
  revenue_type    int4 NOT NULL, 
  order_id        int8 NOT NULL, 
  order_no        varchar(30) NOT NULL, 
  consumer_name   varchar(20) NOT NULL, 
  procedure_rate  numeric(4, 2) NOT NULL, 
  commission_rate numeric(4, 2) NOT NULL, 
  commission_fee  numeric(6, 2) NOT NULL, 
  amount          int4 NOT NULL, 
  review_status   int4 NOT NULL, 
  review_remark   varchar(40) NOT NULL, 
  grant_time      int8 NOT NULL, 
  is_granted      int4 NOT NULL, 
  create_time     int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON COLUMN mch_staff_revenue.id IS '编号';
COMMENT ON COLUMN mch_staff_revenue.revenue_type IS '收入类型,预留默认传1';
COMMENT ON COLUMN mch_staff_revenue.order_id IS '订单编号';
COMMENT ON COLUMN mch_staff_revenue.order_no IS '订单号';
COMMENT ON COLUMN mch_staff_revenue.consumer_name IS '消费者名称';
COMMENT ON COLUMN mch_staff_revenue.review_status IS '1: 待审核  2: 已通过  3: 未通过';
COMMENT ON COLUMN mch_staff_revenue.grant_time IS '佣金发放时间';
COMMENT ON COLUMN mch_staff_revenue.is_granted IS '是否已发放';
COMMENT ON COLUMN mch_staff_revenue.create_time IS ' 创建时间';


ALTER TABLE "public".mm_trusted_info 
  ADD COLUMN cert_image varchar(120) DEFAULT '' NOT NULL;
ALTER TABLE "public".mm_trusted_info 
  ADD COLUMN cert_reverse_image varchar(120) DEFAULT '' NOT NULL;
ALTER TABLE "public".mm_trusted_info 
  ADD COLUMN extra_cert_file varchar(120) DEFAULT '' NOT NULL;
ALTER TABLE "public".mm_trusted_info 
  ADD COLUMN extra_cert_ext1 varchar(120) DEFAULT '' NOT NULL;
ALTER TABLE "public".mm_trusted_info 
  ADD COLUMN extra_cert_ext2 varchar(120) DEFAULT '' NOT NULL;
COMMENT ON COLUMN "public".mm_trusted_info.cert_image IS '证件图片';
COMMENT ON COLUMN "public".mm_trusted_info.cert_reverse_image IS '证件反面图片';
COMMENT ON COLUMN "public".mm_trusted_info.extra_cert_file IS '其他认证资料';
COMMENT ON COLUMN "public".mm_trusted_info.extra_cert_ext1 IS '扩展认证资料1';
COMMENT ON COLUMN "public".mm_trusted_info.extra_cert_ext2 IS '扩展认证资料2';


ALTER TABLE mm_trusted_info
    RENAME TO mm_cert_info;

ALTER TABLE "public".mm_cert_info 
  ADD COLUMN "version" int4 DEFAULT 0 NOT NULL;
COMMENT ON COLUMN "public".mm_cert_info.version IS '版本';

DROP TABLE IF EXISTS "public".mm_cert_info;
CREATE TABLE "public".mm_cert_info (
  id                 serial NOT NULL, 
  member_id          int8 NOT NULL, 
  country_code       varchar(10) NOT NULL, 
  real_name          varchar(10) NOT NULL, 
  card_type          int4 NOT NULL, 
  card_id            varchar(20) NOT NULL, 
  cert_image         varchar(120) DEFAULT ''::character varying NOT NULL, 
  cert_reverse_image varchar(120) DEFAULT ''::character varying NOT NULL, 
  trust_image        varchar(120) NOT NULL, 
  extra_cert_file    varchar(120) DEFAULT ''::character varying NOT NULL, 
  extra_cert_ext1    varchar(120) DEFAULT ''::character varying NOT NULL, 
  extra_cert_ext2    varchar(120) DEFAULT ''::character varying NOT NULL, 
  version            int4 DEFAULT 0 NOT NULL, 
  manual_review      int4 NOT NULL, 
  review_status      int2 DEFAULT 0 NOT NULL, 
  review_time        int4 NOT NULL, 
  remark             varchar(120) NOT NULL, 
  update_time        int4 NOT NULL, 
  CONSTRAINT mm_cert_info_pkey 
    PRIMARY KEY (id));
COMMENT ON TABLE "public".mm_cert_info IS '会员认证信息';
COMMENT ON COLUMN "public".mm_cert_info.id IS '编号';
COMMENT ON COLUMN "public".mm_cert_info.country_code IS '国家代码';
COMMENT ON COLUMN "public".mm_cert_info.real_name IS '真实姓名';
COMMENT ON COLUMN "public".mm_cert_info.card_type IS '证件类型';
COMMENT ON COLUMN "public".mm_cert_info.card_id IS '证件编号';
COMMENT ON COLUMN "public".mm_cert_info.cert_image IS '证件图片';
COMMENT ON COLUMN "public".mm_cert_info.cert_reverse_image IS '证件反面图片';
COMMENT ON COLUMN "public".mm_cert_info.trust_image IS '认证图片,人与身份证的图像等';
COMMENT ON COLUMN "public".mm_cert_info.extra_cert_file IS '其他认证资料';
COMMENT ON COLUMN "public".mm_cert_info.extra_cert_ext1 IS '扩展认证资料1';
COMMENT ON COLUMN "public".mm_cert_info.extra_cert_ext2 IS '扩展认证资料2';
COMMENT ON COLUMN "public".mm_cert_info.version IS '版本';
COMMENT ON COLUMN "public".mm_cert_info.manual_review IS '人工审核';
COMMENT ON COLUMN "public".mm_cert_info.review_status IS '审核状态';
COMMENT ON COLUMN "public".mm_cert_info.review_time IS '审核时间';
COMMENT ON COLUMN "public".mm_cert_info.remark IS '备注';
COMMENT ON COLUMN "public".mm_cert_info.update_time IS '更新时间';

-- 删除域名相关配置
delete FROM registry WHERE key like 'domain_%';

-- 更改地区表名
ALTER TABLE china_area RENAME TO sys_district;

-- 文章表添加mch_id字段
ALTER TABLE article_list 
  ADD COLUMN "mch_id" int4 DEFAULT 0 NOT NULL;
  COMMENT ON COLUMN article_list.mch_id IS '商户编号';

ALTER TABLE article_list 
  ADD COLUMN "short_title" varchar(20) NOT NULL ;
  COMMENT ON COLUMN article_list.mch_id IS '短标题';

ALTER TABLE article_list 
  DROP COLUMN IF EXISTS "small_title";

ALTER TABLE "public".article_list 
  ADD COLUMN like_count int4 NOT NULL;
ALTER TABLE "public".article_list 
  ADD COLUMN dislike_count int4 NOT NULL;
COMMENT ON COLUMN "public".article_list.like_count IS '喜欢的数量';
COMMENT ON COLUMN "public".article_list.dislike_count IS '不喜欢的数量';



-- 发票抬头
CREATE TABLE invoice_headers (
  id           int8 NOT NULL UNIQUE, 
  user_type    int4 NOT NULL, 
  user_id      int8 NOT NULL, 
  invoice_type int4 NOT NULL, 
  issue_type   int4 NOT NULL, 
  header_name  varchar(20) NOT NULL, 
  tax_code   varchar(40) NOT NULL, 
  sign_address varchar(60) NOT NULL, 
  sign_tel     varchar(20) NOT NULL, 
  bank_name    varchar(20) NOT NULL, 
  bank_account varchar(20) NOT NULL, 
  remarks      varchar(20) NOT NULL, 
  is_default   int4 NOT NULL, 
  create_time  int8 NOT NULL);
COMMENT ON TABLE invoice_headers IS '发票抬头';
COMMENT ON COLUMN invoice_headers.id IS '编号';
COMMENT ON COLUMN invoice_headers.user_type IS '用户类型,1:会员  2:商户';
COMMENT ON COLUMN invoice_headers.user_id IS '会员/商户编号';
COMMENT ON COLUMN invoice_headers.invoice_type IS '发票类型: 1:增值税普通发票 2:增值税专用发票 3:形式发票';
COMMENT ON COLUMN invoice_headers.issue_type IS '开具类型, 1: 个人 2:企业';
COMMENT ON COLUMN invoice_headers.tax_code IS '纳税人识别号';
COMMENT ON COLUMN invoice_headers.sign_address IS '注册场所地址';
COMMENT ON COLUMN invoice_headers.sign_tel IS '注册固定电话';
COMMENT ON COLUMN invoice_headers.bank_name IS '基本户开户银行名';
COMMENT ON COLUMN invoice_headers.bank_account IS '基本户开户账号';
COMMENT ON COLUMN invoice_headers.remarks IS '备注';
COMMENT ON COLUMN invoice_headers.is_default IS '是否默认';

-- 发票内容/ 发票备注 / 邮箱 / 留言
CREATE TABLE invoice (
  id           int8 NOT NULL UNIQUE, 
  invoice_code varchar(32) NOT NULL,
  invoice_no varchar(32) NOT NULL,
  user_type    int4 NOT NULL, 
  user_id      int8 NOT NULL, 
  invoice_type int4 NOT NULL, 
  issue_type   int4 NOT NULL,
  seller_name varchar(20) NOT NULL,
  seller_tax_code varchar(64) NOT NULL,
  purchaser_name  varchar(20) NOT NULL, 
  purchaser_tax_code  varchar(40) NOT NULL, 
  invoice_amount decimal(10,2) NOT NULL,
  tax_amount decimal(10,2) NOT NULL,
  remark varchar(64) NOT NULL,
  issue_remark varchar(64) NOT NULL,
  receive_email varchar(64) NOT NULL,
  invoice_status int4 NOT NULL,
  invoice_time int8 NOT NULL,
  create_time int8 NOT NULL,
  update_time int8 NOT NULL);

COMMENT ON TABLE invoice IS '发票';
COMMENT ON COLUMN invoice.id IS '编号';
COMMENT ON COLUMN invoice.invoice_code IS '发票代码';
COMMENT ON COLUMN invoice.invoice_no IS '发票号码';
COMMENT ON COLUMN invoice.user_type IS '用户类型,1:会员  2:商户';
COMMENT ON COLUMN invoice.user_id IS '会员/商户编号';
COMMENT ON COLUMN invoice.invoice_type IS '发票类型: 1:增值税普通发票 2:增值税专用发票 3:形式发票';
COMMENT ON COLUMN invoice.issue_type IS '开具类型, 1: 个人 2:企业';
COMMENT ON COLUMN invoice.seller_name IS '销售方名称';
COMMENT ON COLUMN invoice.seller_tax_code IS  '销售方纳税人识别号';
COMMENT ON COLUMN invoice.purchaser_name IS '买方名称';
COMMENT ON COLUMN invoice.purchaser_tax_code IS '买方纳税人识别号';
COMMENT ON COLUMN invoice.invoice_amount IS '合计金额';
COMMENT ON COLUMN invoice.tax_amount IS '合计税额';
COMMENT ON COLUMN invoice.remark IS '备注';
COMMENT ON COLUMN invoice.issue_remark IS '开具备注/开票失败备注';
COMMENT ON COLUMN invoice.receive_email IS '发票接收邮箱地址';
COMMENT ON COLUMN invoice.invoice_status IS '发票状态,1:待开票 2:开票完成 3:未通过';
COMMENT ON COLUMN invoice.invoice_time IS '开票时间';
COMMENT ON COLUMN invoice.create_time IS '创建时间';
COMMENT ON COLUMN invoice.update_time IS '更新时间';

CREATE TABLE invoice_item (
  id int8 NOT NULL UNIQUE, 
  invoice_id  int8 NOT NULL, 
  item_name varchar(50) NOT NULL,
  item_spec varchar(50) NOT NULL,
  price decimal(10,2) NOT NULL,
  quantity int4 NOT NULL,
  tax_rate decimal(4,2) NOT NULL,
  unit varchar(4) NOT NULL,
  amount decimal(10,2) NOT NULL,
  tax_amount decimal(6,2) NOT NULL,
  create_time int8 NOT NULL,
  update_time int8 NOT NULL
);

COMMENT ON TABLE invoice_item IS '发票项目';
COMMENT ON COLUMN invoice_item.id IS '编号';
COMMENT ON COLUMN invoice_item.invoice_id IS '发票编号';
COMMENT ON COLUMN invoice_item.item_name IS '项目名称';
COMMENT ON COLUMN invoice_item.item_spec IS '项目规格';
COMMENT ON COLUMN invoice_item.price IS '价格';
COMMENT ON COLUMN invoice_item.quantity IS  '数量';
COMMENT ON COLUMN invoice_item.tax_rate IS '税率';
COMMENT ON COLUMN invoice_item.unit IS  '计量单位';
COMMENT ON COLUMN invoice_item.amount IS '总金额';
COMMENT ON COLUMN invoice_item.tax_amount IS '税额';
COMMENT ON COLUMN invoice_item.create_time IS '创建时间';
COMMENT ON COLUMN invoice_item.update_time IS '更新时间';