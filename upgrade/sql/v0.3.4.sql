/***** ----------- 此SQL涉及到表结构 ------------------*/

/** 2022-11-25 */
ALTER TABLE IF EXISTS public.perm_dept
    ADD COLUMN code character varying(50) NOT NULL DEFAULT '';

COMMENT ON COLUMN public.perm_dept.code
    IS '部门代码';

ALTER TABLE public.order_list ALTER COLUMN id TYPE bigint;

ALTER TABLE public.order_list
    ALTER COLUMN subject TYPE character varying(40) COLLATE pg_catalog."default";

/** 2022-12-03 */
ALTER TABLE IF EXISTS public.order_list
    ADD COLUMN consignee_modified integer NOT NULL DEFAULT 0;

COMMENT ON COLUMN public.order_list.consignee_modified
    IS '收货地址是否已修改';

ALTER TABLE IF EXISTS public.sale_sub_order
    RENAME is_suspend TO is_forbidden;

COMMENT ON COLUMN public.sale_sub_order.is_forbidden
    IS '是否被用户删除';

/** 2022-12-21 */
ALTER TABLE IF EXISTS public.product_model_attr
    RENAME pro_model TO prod_model;

ALTER TABLE IF EXISTS public.product_model_attr
    RENAME multi_chk TO multi_check;
ALTER TABLE IF EXISTS public.product_model_attr_item
    RENAME pro_model TO prod_model;

ALTER TABLE IF EXISTS public.product_model_brand
    RENAME pro_model TO prod_model;
ALTER TABLE IF EXISTS public.product_model_spec
    RENAME pro_model TO prod_model;
ALTER TABLE IF EXISTS public.product_model_spec_item
    RENAME pro_model TO prod_model;

/** 2022-12-29 15:39 */
ALTER TABLE IF EXISTS public.mm_member
    RENAME "name" TO nick_name;

ALTER TABLE "public".article_list 
  alter column update_time set not null;

CREATE INDEX mm_account_member_id_wallet_code 
  ON "public".mm_account (wallet_code, member_id);


/** 2022-12-29 12:31 */
  CREATE TABLE mm_flag_request (
  id               BIGSERIAL NOT NULL, 
  member_id       int8 NOT NULL, 
  request_flag    int4 NOT NULL, 
  contact_name    varchar(20) NOT NULL, 
  contact_phone   varchar(20) NOT NULL, 
  contact_address varchar(120) NOT NULL, 
  audit_state     int4 NOT NULL, 
  audit_time      int8 NOT NULL, 
  audit_uid       int8 NOT NULL, 
  audit_remark    varchar(120) NOT NULL, 
  create_time     int8 NOT NULL, 
  PRIMARY KEY (id));
COMMENT ON TABLE mm_flag_request IS '会员申请标志请求';
COMMENT ON COLUMN mm_flag_request.id IS '编号';
COMMENT ON COLUMN mm_flag_request.member_id IS '会员编号';
COMMENT ON COLUMN mm_flag_request.request_flag IS '申请标志';
COMMENT ON COLUMN mm_flag_request.contact_name IS '联系人';
COMMENT ON COLUMN mm_flag_request.contact_phone IS '联系电话';
COMMENT ON COLUMN mm_flag_request.contact_address IS '联系地址';
COMMENT ON COLUMN mm_flag_request.audit_state IS '审核状态';
COMMENT ON COLUMN mm_flag_request.audit_time IS '审核时间';
COMMENT ON COLUMN mm_flag_request.audit_uid IS '审核人';
COMMENT ON COLUMN mm_flag_request.audit_remark IS '审核意见';
COMMENT ON COLUMN mm_flag_request.create_time IS '创建时间';

ALTER TABLE IF EXISTS public.mm_balance_log
    RENAME title TO subject;
ALTER TABLE IF EXISTS public.mm_balance_log
    RENAME amount TO change_value;
	ALTER TABLE IF EXISTS public.mm_balance_log
    RENAME csn_fee TO procedure_fee;
		ALTER TABLE IF EXISTS public.mm_balance_log
    RENAME review_state TO audit_state;

ALTER TABLE IF EXISTS public.mm_flow_log
    RENAME title TO subject;
ALTER TABLE IF EXISTS public.mm_flow_log
    RENAME amount TO change_value;
	ALTER TABLE IF EXISTS public.mm_flow_log
    RENAME csn_fee TO procedure_fee;
		ALTER TABLE IF EXISTS public.mm_flow_log
    RENAME review_state TO audit_state;

  
ALTER TABLE IF EXISTS public.mm_integral_log
    RENAME title TO subject;
ALTER TABLE IF EXISTS public.mm_integral_log
    RENAME value TO change_value;
ALTER TABLE IF EXISTS public.mm_integral_log
    RENAME review_state TO audit_state;

