
-- DROP TABLE IF EXISTS public.pay_integrate_app;

CREATE TABLE IF NOT EXISTS public.pay_integrate_app
(
    id serial NOT NULL,
    app_name character varying(20) COLLATE pg_catalog."default" NOT NULL,
    app_url character varying(120) COLLATE pg_catalog."default" NOT NULL,
    integrate_type integer NOT NULL,
    sort_number integer NOT NULL,
    enabled integer NOT NULL,
    CONSTRAINT pay_integrate_app_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.pay_integrate_app
    OWNER to postgres;

COMMENT ON TABLE public.pay_integrate_app
    IS '集成支付应用';

COMMENT ON COLUMN public.pay_integrate_app.id
    IS '编号';

COMMENT ON COLUMN public.pay_integrate_app.app_name
    IS '支付应用名称';

COMMENT ON COLUMN public.pay_integrate_app.app_url
    IS '支付应用接口';

COMMENT ON COLUMN public.pay_integrate_app.enabled
    IS '是否启用';

COMMENT ON COLUMN public.pay_integrate_app.integrate_type
    IS '集成方式: 1:API调用 2: 跳转';

COMMENT ON COLUMN public.pay_integrate_app.sort_number
    IS '显示顺序';

ALTER TABLE IF EXISTS public.pay_integrate_app
    ADD COLUMN hint character varying(30) NOT NULL DEFAULT '';

COMMENT ON COLUMN public.pay_integrate_app.hint
    IS '支付提示信息';

ALTER TABLE IF EXISTS public.pay_integrate_app
    ADD COLUMN highlight integer NOT NULL DEFAULT 0;

COMMENT ON COLUMN public.pay_integrate_app.highlight
    IS '是否高亮显示';


/** 05-25 订单数据表调整　*/

CREATE TABLE "public".order_list (
     id               serial NOT NULL,
     order_no         varchar(40) NOT NULL,
     order_type       int4 NOT NULL,
     subject          varchar(20) NOT NULL,
     buyer_id         int4 NOT NULL,
     buyer_user       varchar(20) NOT NULL,
     item_amount      int8 NOT NULL,
     discount_amount  int8 NOT NULL,
     express_fee      int8 NOT NULL,
     package_fee      int8 NOT NULL,
     final_amount     int8 NOT NULL,
     consignee_name varchar(45) NOT NULL,
     consignee_phone  varchar(45) NOT NULL,
     shipping_address varchar(120) NOT NULL,
     is_break         int4 NOT NULL,
     state            int4 NOT NULL,
     create_time      int8 NOT NULL,
     update_time      int8 NOT NULL,
     CONSTRAINT order_list_pkey
         PRIMARY KEY (id));
COMMENT ON COLUMN "public".order_list.id IS '编号';
COMMENT ON COLUMN "public".order_list.order_no IS '订单号';
COMMENT ON COLUMN "public".order_list.order_type IS '订单类型1:普通 2:批发 3:线下';
COMMENT ON COLUMN "public".order_list.subject IS '订单主题';
COMMENT ON COLUMN "public".order_list.buyer_id IS '买家';
COMMENT ON COLUMN "public".order_list.buyer_user IS '买家用户名';
COMMENT ON COLUMN "public".order_list.item_amount IS '商品金额';
COMMENT ON COLUMN "public".order_list.discount_amount IS '抵扣金额';
COMMENT ON COLUMN "public".order_list.express_fee IS '物流费';
COMMENT ON COLUMN "public".order_list.package_fee IS '包装费';
COMMENT ON COLUMN "public".order_list.final_amount IS '订单最终金额';
COMMENT ON COLUMN "public".order_list.consignee_name IS '收货人姓名';
COMMENT ON COLUMN "public".order_list.consignee_phone IS '收货人电话';
COMMENT ON COLUMN "public".order_list.shipping_address IS '收货人地址';
COMMENT ON COLUMN "public".order_list.is_break IS '是否拆分';
COMMENT ON COLUMN "public".order_list.state IS '订单状态';
COMMENT ON COLUMN "public".order_list.create_time IS '创建时间';
COMMENT ON COLUMN "public".order_list.update_time IS '更新时间';


DROP TABLE IF EXISTS "public".order_wholesale_order CASCADE;

CREATE TABLE "public".order_wholesale_order (
    id            bigserial NOT NULL,
    order_no      varchar(20) NOT NULL,
    order_id      int8 NOT NULL,
    buyer_id      int8 NOT NULL,
    vendor_id     int8 NOT NULL,
    shop_id       int8 NOT NULL,
    is_paid       int4 NOT NULL,
    buyer_comment varchar(120) NOT NULL,
    remark        varchar(120) NOT NULL,
    state         int4 NOT NULL,
    create_time   int8 NOT NULL,
    update_time   int8 NOT NULL,
    CONSTRAINT order_wholesale_order_pkey
        PRIMARY KEY (id));
COMMENT ON TABLE "public".order_wholesale_order IS '批发订单';
COMMENT ON COLUMN "public".order_wholesale_order.id IS '编号';
COMMENT ON COLUMN "public".order_wholesale_order.order_no IS '订单号';
COMMENT ON COLUMN "public".order_wholesale_order.order_id IS '订单编号';
COMMENT ON COLUMN "public".order_wholesale_order.buyer_id IS '买家';
COMMENT ON COLUMN "public".order_wholesale_order.vendor_id IS '供货商';
COMMENT ON COLUMN "public".order_wholesale_order.shop_id IS '店铺编号';
COMMENT ON COLUMN "public".order_wholesale_order.is_paid IS '是否支付';
COMMENT ON COLUMN "public".order_wholesale_order.buyer_comment IS '买家留言';
COMMENT ON COLUMN "public".order_wholesale_order.remark IS '备注';
COMMENT ON COLUMN "public".order_wholesale_order.state IS '订单状态';
COMMENT ON COLUMN "public".order_wholesale_order.create_time IS '创建时间';
COMMENT ON COLUMN "public".order_wholesale_order.update_time IS '更新时间';

-- Table: public.article_list

-- DROP TABLE IF EXISTS public.article_list;

CREATE TABLE IF NOT EXISTS public.article_list
(
    id bigserial NOT NULL ,
    cat_id bigint NOT NULL,
    title character varying(120) COLLATE pg_catalog."default" NOT NULL,
    small_title character varying(45) COLLATE pg_catalog."default" NOT NULL,
    thumbnail character varying(120) COLLATE pg_catalog."default" NOT NULL,
    publisher_id integer NOT NULL,
    location character varying(120) COLLATE pg_catalog."default" NOT NULL,
    priority integer NOT NULL,
    access_key character varying(45) COLLATE pg_catalog."default" NOT NULL,
    content text COLLATE pg_catalog."default" NOT NULL,
    tags character varying(120) COLLATE pg_catalog."default" NOT NULL,
    view_count integer NOT NULL,
    sort_num integer NOT NULL,
    create_time bigint NOT NULL,
    update_time bigint,
    CONSTRAINT article_list_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.article_list
    OWNER to postgres;

COMMENT ON COLUMN public.article_list.id
    IS '编号';

COMMENT ON COLUMN public.article_list.cat_id
    IS '分类编号';

COMMENT ON COLUMN public.article_list.title
    IS '标题';

COMMENT ON COLUMN public.article_list.small_title
    IS '小标题';

COMMENT ON COLUMN public.article_list.thumbnail
    IS '缩略图';

COMMENT ON COLUMN public.article_list.publisher_id
    IS '作者';

COMMENT ON COLUMN public.article_list.location
    IS '地址';

COMMENT ON COLUMN public.article_list.priority
    IS '优先级';

COMMENT ON COLUMN public.article_list.access_key
    IS '访问密钥';

COMMENT ON COLUMN public.article_list.content
    IS '内容';

COMMENT ON COLUMN public.article_list.tags
    IS '标签';

COMMENT ON COLUMN public.article_list.view_count
    IS '浏览次数';

COMMENT ON COLUMN public.article_list.sort_num
    IS '排列序号';

COMMENT ON COLUMN public.article_list.create_time
    IS '创建时间';

COMMENT ON COLUMN public.article_list.update_time
    IS '更新时间';


ALTER TABLE IF EXISTS public.article_category
    OWNER to postgres;

COMMENT ON COLUMN public.article_category.id
    IS '编号';

COMMENT ON COLUMN public.article_category.parent_id
    IS '上级编号';

COMMENT ON COLUMN public.article_category.perm_flag
    IS '权限标志';

COMMENT ON COLUMN public.article_category.name
    IS '分类编号';

COMMENT ON COLUMN public.article_category.cat_alias
    IS '分类别名';

COMMENT ON COLUMN public.article_category.title
    IS '标题';

COMMENT ON COLUMN public.article_category.keywords
    IS '关键词';

COMMENT ON COLUMN public.article_category.describe
    IS '描述';

COMMENT ON COLUMN public.article_category.sort_num
    IS '排序编号';

COMMENT ON COLUMN public.article_category.location
    IS '地址';

COMMENT ON COLUMN public.article_category.update_time
    IS '更新时间';