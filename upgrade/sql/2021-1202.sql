ALTER TABLE "public".mch_merchant
    ADD COLUMN salt varchar(10) DEFAULT '' NOT NULL;
COMMENT
ON COLUMN "public".mch_merchant.salt IS '加密盐';

ALTER TABLE public.product
ALTER
COLUMN name TYPE character varying(120) COLLATE pg_catalog."default";


ALTER TABLE IF EXISTS public.product DROP COLUMN IF EXISTS shelve_state;
ALTER TABLE IF EXISTS public.product DROP COLUMN IF EXISTS review_state;
ALTER TABLE IF EXISTS public.product DROP COLUMN IF EXISTS sale_price;


/* 用bigint存储金额　*/
ALTER TABLE public.item_info
ALTER COLUMN cost TYPE bigint;

ALTER TABLE public.item_info
ALTER COLUMN price TYPE bigint;

ALTER TABLE public.item_info
ALTER COLUMN retail_price TYPE bigint;

ALTER TABLE public.item_sku
ALTER COLUMN retail_price TYPE bigint;

ALTER TABLE public.item_sku
ALTER COLUMN price TYPE bigint;

ALTER TABLE public.item_sku
ALTER COLUMN cost TYPE bigint;


ALTER TABLE public.item_snapshot
ALTER COLUMN cost TYPE bigint;

ALTER TABLE public.item_snapshot
ALTER COLUMN price TYPE bigint;

ALTER TABLE public.item_snapshot
ALTER COLUMN retail_price TYPE bigint;

ALTER TABLE public.item_trade_snapshot
ALTER COLUMN cost TYPE bigint;

ALTER TABLE public.item_trade_snapshot
ALTER COLUMN price TYPE bigint;

ALTER TABLE public.gs_member_price
ALTER COLUMN price TYPE bigint;

ALTER TABLE public.ws_item
ALTER COLUMN price TYPE bigint;


ALTER TABLE public.sale_order
ALTER COLUMN item_amount TYPE bigint;

ALTER TABLE public.sale_order
ALTER COLUMN discount_amount TYPE bigint;

ALTER TABLE public.sale_order
ALTER COLUMN express_fee TYPE bigint;

ALTER TABLE public.sale_order
ALTER COLUMN package_fee TYPE bigint;

ALTER TABLE public.sale_order
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.sale_order_item
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.sale_order_item
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.sale_sub_order
ALTER COLUMN discount_amount TYPE bigint;

ALTER TABLE public.sale_sub_order
ALTER COLUMN express_fee TYPE bigint;

ALTER TABLE public.sale_sub_order
ALTER COLUMN package_fee TYPE bigint;

ALTER TABLE public.sale_sub_order
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN freeze_balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN expired_balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN wallet_balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN freeze_wallet TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN expired_wallet TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN total_wallet_amount TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN flow_balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN grow_balance TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN grow_amount TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN grow_earnings TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN grow_total_earnings TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN total_charge TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN total_pay TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN total_expense TYPE bigint;

ALTER TABLE public.mm_account
ALTER COLUMN priority_pay TYPE integer;

ALTER TABLE public.express_area_set
ALTER COLUMN first_fee TYPE bigint;

ALTER TABLE public.express_area_set
ALTER COLUMN add_fee TYPE bigint;

ALTER TABLE public.mch_balance_log
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.mch_balance_log
ALTER COLUMN csn_amount TYPE bigint;

ALTER TABLE public.mm_balance_info
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.mm_balance_info
ALTER COLUMN csn_amount TYPE bigint;

ALTER TABLE public.mm_balance_log
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.mm_balance_log
ALTER COLUMN csn_fee TYPE bigint;

ALTER TABLE public.mm_flow_log
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.mm_flow_log
ALTER COLUMN csn_fee TYPE bigint;

ALTER TABLE public.mm_wallet_log
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.mm_wallet_log
ALTER COLUMN csn_fee TYPE bigint;

ALTER TABLE public.order_wholesale_order
ALTER COLUMN item_amount TYPE bigint;

ALTER TABLE public.order_wholesale_order
ALTER COLUMN discount_amount TYPE bigint;

ALTER TABLE public.order_wholesale_order
ALTER COLUMN express_fee TYPE bigint;

ALTER TABLE public.order_wholesale_order
ALTER COLUMN package_fee TYPE bigint;

ALTER TABLE public.order_wholesale_order
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.order_wholesale_item
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.order_wholesale_item
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.order_trade_order
ALTER COLUMN order_amount TYPE bigint;

ALTER TABLE public.order_trade_order
ALTER COLUMN discount_amount TYPE bigint;

ALTER TABLE public.order_trade_order
ALTER COLUMN final_amount TYPE bigint;

ALTER TABLE public.order_trade_order
ALTER COLUMN trade_rate TYPE bigint;

COMMENT ON TABLE public.item_info
  IS '商品信息';

COMMENT ON COLUMN public.item_info.id
    IS '编号';

COMMENT ON COLUMN public.item_info.product_id
    IS '产品编号';

COMMENT ON COLUMN public.item_info.prom_flag
    IS '营销标志';

COMMENT ON COLUMN public.item_info.cat_id
    IS '分类编号';

COMMENT ON COLUMN public.item_info.vendor_id
    IS '供应商编号';

COMMENT ON COLUMN public.item_info.brand_id
    IS '品牌编号';

COMMENT ON COLUMN public.item_info.shop_id
    IS '店铺编号';

COMMENT ON COLUMN public.item_info.shop_cat_id
    IS '店铺分类编号';

COMMENT ON COLUMN public.item_info.express_tid
    IS '快递模板';

COMMENT ON COLUMN public.item_info.title
    IS '商品标题';

COMMENT ON COLUMN public.item_info.short_title
    IS '商品小标题';

COMMENT ON COLUMN public.item_info.code
    IS '商品编码';

COMMENT ON COLUMN public.item_info.image
    IS '商品主图';

COMMENT ON COLUMN public.item_info.is_present
    IS '是否为赠品';

COMMENT ON COLUMN public.item_info.price_range
    IS '价格区间';

COMMENT ON COLUMN public.item_info.stock_num
    IS '库存数量';

COMMENT ON COLUMN public.item_info.sale_num
    IS '销售数量';

COMMENT ON COLUMN public.item_info.sku_num
    IS '规格数量';

COMMENT ON COLUMN public.item_info.sku_id
    IS 'SKU编号';

COMMENT ON COLUMN public.item_info.cost
    IS '成本价';

COMMENT ON COLUMN public.item_info.price
    IS '销售价';

COMMENT ON COLUMN public.item_info.retail_price
    IS '零售价';

COMMENT ON COLUMN public.item_info.weight
    IS '重量';

COMMENT ON COLUMN public.item_info.bulk
    IS '容积';

COMMENT ON COLUMN public.item_info.shelve_state
    IS '上架状态';

COMMENT ON COLUMN public.item_info.review_state
    IS '审核状态';

COMMENT ON COLUMN public.item_info.review_remark
    IS '审核意见';

COMMENT ON COLUMN public.item_info.sort_num
    IS '排列序号';

COMMENT ON COLUMN public.item_info.create_time
    IS '创建时间';

COMMENT ON COLUMN public.item_info.update_time
    IS '更新时间';

ALTER TABLE public.sale_sub_order
ALTER COLUMN item_amount TYPE bigint;

ALTER TABLE public.sale_return
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.sale_refund
ALTER COLUMN amount TYPE bigint;

ALTER TABLE public.sale_after_order
ALTER COLUMN type TYPE smallint USING type::text::smallint;

/** 2012-12-12 */
COMMENT ON TABLE public.portal_nav
  IS '门户导航';

ALTER TABLE IF EXISTS public.portal_nav_type
    RENAME TO portal_nav_group;

COMMENT ON TABLE public.portal_nav_group
  IS '导航分组';

ALTER TABLE IF EXISTS public.portal_nav
    ADD COLUMN nav_group character varying(20) DEFAULT '' NOT NULL;

COMMENT ON COLUMN public.portal_nav.nav_group
    IS '导航分组';

/** 2012-12-13 */

COMMENT ON TABLE public.ad_group
  IS '广告分组';

COMMENT ON COLUMN public.ad_group.id
    IS '编号';

COMMENT ON COLUMN public.ad_group.name
    IS '名称';

COMMENT ON COLUMN public.ad_group.opened
    IS '是否开放';

COMMENT ON COLUMN public.ad_group.enabled
    IS '是否启用';

ALTER TABLE IF EXISTS public.ad_group
    ADD COLUMN flag integer NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.ad_group
    ADD COLUMN flag integer NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.ad_position
    ADD COLUMN flag integer NOT NULL DEFAULT 0;

COMMENT ON COLUMN public.ad_position.flag
    IS '标志';


ALTER TABLE IF EXISTS public.ad_position
    ADD COLUMN group_name character varying(20) NOT NULL DEFAULT '';

COMMENT ON COLUMN public.ad_position.group_name
    IS '分组名称';


COMMENT ON TABLE public.ad_position
  IS '广告位';

COMMENT ON COLUMN public.ad_position.id
    IS '编号';

COMMENT ON COLUMN public.ad_position.key
    IS '广告位编码';

COMMENT ON COLUMN public.ad_position.name
    IS '广告位名称';

ALTER TABLE IF EXISTS public.ad_position
    RENAME default_id TO put_aid;

COMMENT ON COLUMN public.ad_position.put_aid
    IS '投放的广告编号';

COMMENT ON TABLE public.ad_list
  IS '广告列表';

COMMENT ON COLUMN public.ad_list.id
    IS '编号';

COMMENT ON COLUMN public.ad_list.user_id
    IS '用户编号';

COMMENT ON COLUMN public.ad_list.name
    IS '广告名称';

COMMENT ON COLUMN public.ad_list.type_id
    IS '广告类型';

COMMENT ON COLUMN public.ad_list.show_times
    IS '展现次数';

COMMENT ON COLUMN public.ad_list.click_times
    IS '点击次数';

COMMENT ON COLUMN public.ad_list.show_days
    IS '显示天数';

COMMENT ON COLUMN public.ad_list.update_time
    IS '更新时间';

/** 2021-12-14 */
DROP TABLE IF EXISTS public.ad_image_ad;
DROP TABLE IF EXISTS public.ad_group;

COMMENT ON TABLE public.ad_image
  IS '广告图片';

COMMENT ON COLUMN public.ad_image.id
    IS '编号';

COMMENT ON COLUMN public.ad_image.ad_id
    IS '广告编号';

COMMENT ON COLUMN public.ad_image.title
    IS '标题';

COMMENT ON COLUMN public.ad_image.link_url
    IS '链接地址';

COMMENT ON COLUMN public.ad_image.image_url
    IS '图片地址';

COMMENT ON COLUMN public.ad_image.sort_num
    IS '排列序号';

COMMENT ON COLUMN public.ad_image.enabled
    IS '是否启用';

COMMENT ON TABLE public.ad_hyperlink
  IS '文本广告';

COMMENT ON COLUMN public.ad_hyperlink.id
    IS '编号';

COMMENT ON COLUMN public.ad_hyperlink.ad_id
    IS '广告编号';

COMMENT ON COLUMN public.ad_hyperlink.title
    IS '标题';

COMMENT ON COLUMN public.ad_hyperlink.link_url
    IS '链接地址';

ALTER TABLE IF EXISTS public.ad_position DROP COLUMN IF EXISTS group_id;


ALTER TABLE IF EXISTS public.registry
    ADD COLUMN group_name character varying(20) DEFAULT '' NOT NULL;

COMMENT ON COLUMN public.registry.group_name
    IS '分组名称';

/* 2021-12-21 */
COMMENT ON TABLE public.sale_cart
  IS '购物车';
ALTER TABLE public.sale_cart
ALTER COLUMN code TYPE character varying(40) COLLATE pg_catalog."default";

/** 2021-12-30 */
ALTER TABLE public.portal_nav
ALTER COLUMN image TYPE character varying(160) COLLATE pg_catalog."default";
update portal_nav set image=replace(image,'http://','https://')

/** 2021-12-31 */
ALTER TABLE IF EXISTS public.ex_page
    RENAME TO arc_page;

COMMENT ON TABLE public.arc_page
  IS '单页';

ALTER TABLE public.arc_page
ALTER COLUMN id TYPE bigint;
COMMENT ON COLUMN public.arc_page.id
    IS '编号';

ALTER TABLE public.arc_page
ALTER COLUMN user_id TYPE bigint;
COMMENT ON COLUMN public.arc_page.user_id
    IS '用户编号,系统为0';

COMMENT ON COLUMN public.arc_page.title
    IS '标题';

ALTER TABLE IF EXISTS public.arc_page
    RENAME perm_flag TO flag;

COMMENT ON COLUMN public.arc_page.flag
    IS '标志';

COMMENT ON COLUMN public.arc_page.access_key
    IS '访问钥匙';

ALTER TABLE IF EXISTS public.arc_page
    RENAME str_indent TO code;

COMMENT ON COLUMN public.arc_page.code
    IS '页面代码';

COMMENT ON COLUMN public.arc_page.keyword
    IS '关键词';

COMMENT ON COLUMN public.arc_page.description
    IS '描述';

COMMENT ON COLUMN public.arc_page.css_path
    IS '样式表路径';

COMMENT ON COLUMN public.arc_page.enabled
    IS '是否启用';

ALTER TABLE IF EXISTS public.arc_page
    RENAME body TO content;

COMMENT ON COLUMN public.arc_page.content
    IS '内容';

COMMENT ON COLUMN public.arc_page.update_time
    IS '更新时间';
