alter table mm_profile rename column sex to gender;
alter table perm_user rename column sex to gender;


ALTER TABLE "public".mm_member
    alter column code set default ' '::character varying;
ALTER TABLE "public".mm_member
    alter column avatar set default ' '::character varying;
ALTER TABLE "public".mm_member
    alter column phone set default ' '::character varying;
ALTER TABLE "public".mm_member
    alter column email set default ' '::character varying;
ALTER TABLE "public".mm_member
    alter column name set default ' '::character varying;
ALTER TABLE "public".mm_member
    alter column real_name set default ''::character varying;
ALTER TABLE "public".mm_member
    ADD COLUMN salt varchar(10) DEFAULT '' NOT NULL;
COMMENT ON COLUMN "public".mm_member."user" IS '用户名';
COMMENT ON COLUMN "public".mm_member.flag IS '会员标志';
COMMENT ON COLUMN "public".mm_member.name IS '昵称';
COMMENT ON COLUMN "public".mm_member.real_name IS '真实姓名';
COMMENT ON COLUMN "public".mm_member.salt IS '加密盐';

COMMENT ON COLUMN public.perm_res.res_type
    IS '资源类型, 0: 目录 1: 资源　2: 菜单  3:　 按钮';


/* 2021-03-11 */
COMMENT ON COLUMN public.mm_level.id
    IS '编号';

COMMENT ON COLUMN public.mm_level.name
    IS '等级名称';

COMMENT ON COLUMN public.mm_level.require_exp
    IS '需要经验值';

COMMENT ON COLUMN public.mm_level.program_signal
    IS '编程符号';

COMMENT ON COLUMN public.mm_level.is_official
    IS '是否正式的会员等级';

COMMENT ON COLUMN public.mm_level.enabled
    IS '是否启用';

COMMENT ON COLUMN public.mm_level.allow_upgrade
    IS '是否允许自动升级';


COMMENT ON COLUMN public.mm_member.id
    IS '编号';

COMMENT ON COLUMN public.mm_member.code
    IS '用户编码';

COMMENT ON COLUMN public.mm_member.pwd
    IS '密码';

COMMENT ON COLUMN public.mm_member.trade_pwd
    IS '交易密码';

COMMENT ON COLUMN public.mm_member.exp
    IS '经验值';

COMMENT ON COLUMN public.mm_member.level
    IS '等级';

COMMENT ON COLUMN public.mm_member.premium_user
    IS '高级用户类型';

COMMENT ON COLUMN public.mm_member.premium_expires
    IS '高级用户过期时间';

COMMENT ON COLUMN public.mm_member.invite_code
    IS '邀请码';

COMMENT ON COLUMN public.mm_member.reg_ip
    IS '注册IP';

COMMENT ON COLUMN public.mm_member.reg_from
    IS '注册来源';

COMMENT ON COLUMN public.mm_member.reg_time
    IS '注册时间';

COMMENT ON COLUMN public.mm_member.check_code
    IS '校验码';

COMMENT ON COLUMN public.mm_member.check_expires
    IS '校验码过期时间';

COMMENT ON COLUMN public.mm_member.login_time
    IS '登录时间';

COMMENT ON COLUMN public.mm_member.last_login_time
    IS '最后登录时间';

COMMENT ON COLUMN public.mm_member.state
    IS '状态';

COMMENT ON COLUMN public.mm_member.update_time
    IS '更新时间';

COMMENT ON COLUMN public.mm_member.avatar
    IS '头像';

COMMENT ON COLUMN public.mm_member.phone
    IS '手机号码';

COMMENT ON COLUMN public.mm_member.email
    IS '电子邮箱';

COMMENT ON TABLE public.mm_level
  IS '会员等级';

ALTER TABLE public.wal_wallet
    ADD COLUMN user_name character varying(20) NOT NULL DEFAULT '';

COMMENT ON COLUMN public.wal_wallet.user_name
    IS '用户名';

COMMENT ON COLUMN public.mch_merchant.id
    IS '编号';

COMMENT ON COLUMN public.mch_sign_up.submit_time
    IS '申请时间';

COMMENT ON COLUMN public.mch_sign_up.update_time
    IS '更新时间';