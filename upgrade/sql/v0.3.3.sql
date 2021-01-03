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
