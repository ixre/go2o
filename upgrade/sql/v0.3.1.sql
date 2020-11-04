ALTER TABLE public.express_provider
    ALTER COLUMN enabled TYPE int2 USING enabled::int;

COMMENT ON COLUMN "public".comm_qr_template.id IS '编号';
COMMENT ON COLUMN "public".comm_qr_template.title IS '模板标题';
COMMENT ON COLUMN "public".comm_qr_template.bg_image IS '背景图片';
COMMENT ON COLUMN "public".comm_qr_template.offset_x IS '垂直偏离量';
COMMENT ON COLUMN "public".comm_qr_template.offset_y IS '垂直偏移量';
COMMENT ON COLUMN "public".comm_qr_template.comment IS '二维码模板文本';
COMMENT ON COLUMN "public".comm_qr_template.callback_url IS '回调地址';
COMMENT ON COLUMN "public".comm_qr_template.enabled IS '是否启用';
COMMENT ON COLUMN "public".dlv_coverage.name IS '配送覆盖区域';
COMMENT ON COLUMN "public".mch_account.mch_id IS '商户编号';
COMMENT ON COLUMN "public".mch_account.balance IS '余额';
COMMENT ON COLUMN "public".mch_account.freeze_amount IS '冻结金额';
COMMENT ON COLUMN "public".mch_account.await_amount IS '待入账金额';
COMMENT ON COLUMN "public".mch_account.present_amount IS '平台赠送金额';
COMMENT ON COLUMN "public".mch_account.sales_amount IS '累计销售总额';
COMMENT ON COLUMN "public".mch_account.refund_amount IS '累计退款金额';
COMMENT ON COLUMN "public".mch_account.take_amount IS '已提取金额';
COMMENT ON COLUMN "public".mch_account.offline_sales IS '线下销售金额';
COMMENT ON COLUMN "public".mch_account.update_time IS '更新时间';
COMMENT ON COLUMN "public".mch_api_info.enabled IS '是否启用';
COMMENT ON COLUMN "public".mch_api_info.white_list IS '白名单';
COMMENT ON COLUMN "public".mch_balance_log.mch_id IS '商户编号';
COMMENT ON COLUMN "public".mch_balance_log.kind IS '日志类型';
COMMENT ON COLUMN "public".mch_balance_log.title IS '标题';
COMMENT ON COLUMN "public".mch_balance_log.outer_no IS '外部订单号';
COMMENT ON COLUMN "public".mch_balance_log.amount IS '金额';
COMMENT ON COLUMN "public".mch_balance_log.csn_amount IS '手续费';
COMMENT ON COLUMN "public".mch_balance_log.state IS '状态';
COMMENT ON COLUMN "public".mch_balance_log.create_time IS '创建时间';
COMMENT ON COLUMN "public".mch_balance_log.update_time IS '更新时间';
COMMENT ON COLUMN "public".mch_day_chart.id IS '编号';
COMMENT ON COLUMN "public".mch_day_chart.mch_id IS '商户编号';
COMMENT ON COLUMN "public".mch_day_chart.order_number IS '新增订单数量';
COMMENT ON COLUMN "public".mch_day_chart.order_amount IS '订单额';
COMMENT ON COLUMN "public".mch_day_chart.buyer_number IS '购物会员数';
COMMENT ON COLUMN "public".mch_day_chart.paid_number IS '支付单数量';
COMMENT ON COLUMN "public".mch_day_chart.paid_amount IS '支付总金额';
COMMENT ON COLUMN "public".mch_day_chart.complete_orders IS '完成订单数';
COMMENT ON COLUMN "public".mch_day_chart.in_amount IS '入帐金额';
COMMENT ON COLUMN "public".mch_day_chart.offline_orders IS '线下订单数量';
COMMENT ON COLUMN "public".mch_day_chart.offline_amount IS '线下订单金额';
COMMENT ON COLUMN "public".mch_day_chart."date" IS '日期';
COMMENT ON COLUMN "public".mch_day_chart.date_str IS '日期字符串';
COMMENT ON COLUMN "public".mch_day_chart.update_time IS '更新时间';
COMMENT ON COLUMN "public".mch_enterprise_info.mch_id IS '商户编号';
COMMENT ON COLUMN "public".mch_enterprise_info.company_name IS '公司名称';
COMMENT ON COLUMN "public".mch_enterprise_info.company_no IS '营业执照编号';
COMMENT ON COLUMN "public".mch_enterprise_info.person_name IS '法人姓名';
COMMENT ON COLUMN "public".mch_enterprise_info.person_id IS '法人身份证号';
COMMENT ON COLUMN "public".mch_enterprise_info.tel IS '公司电话';
COMMENT ON COLUMN "public".mch_enterprise_info.province IS '所在省';
COMMENT ON COLUMN "public".mch_enterprise_info.city IS '所在市';
COMMENT ON COLUMN "public".mch_enterprise_info.district IS '所在区';
COMMENT ON COLUMN "public".mch_enterprise_info.location IS '位置';
COMMENT ON COLUMN "public".mch_enterprise_info.address IS '公司地址';
COMMENT ON COLUMN "public".mch_enterprise_info.person_image IS '法人身份证照片';
COMMENT ON COLUMN "public".mch_enterprise_info.company_image IS '营业执照照片';
COMMENT ON COLUMN "public".mch_enterprise_info.auth_doc IS '授权书';
COMMENT ON COLUMN "public".mch_enterprise_info.review_time IS '审核时间';
COMMENT ON COLUMN "public".mch_enterprise_info.review_remark IS '审核备注';
COMMENT ON COLUMN "public".mch_enterprise_info.review_state IS '审核状态';
COMMENT ON COLUMN "public".mch_offline_shop.deliver_radius IS '配送范围';
COMMENT ON COLUMN "public".mch_online_shop.id IS '店铺编号';
COMMENT ON COLUMN "public".mch_online_shop.vendor_id IS '商户编号';
COMMENT ON COLUMN "public".mch_online_shop.shop_name IS '店铺名称';
COMMENT ON COLUMN "public".mch_online_shop.logo IS '店铺标志';
COMMENT ON COLUMN "public".mch_online_shop.host IS '自定义 域名';
COMMENT ON COLUMN "public".mch_online_shop.alias IS '个性化域名';
COMMENT ON COLUMN "public".mch_online_shop.tel IS '电话';
COMMENT ON COLUMN "public".mch_online_shop.addr IS '地址';
COMMENT ON COLUMN "public".mch_online_shop.shop_title IS '店铺标题';
COMMENT ON COLUMN "public".mch_online_shop.shop_notice IS '店铺公告';
COMMENT ON COLUMN "public".mch_online_shop.flag IS '标志';
COMMENT ON COLUMN "public".mch_online_shop.state IS '状态';
COMMENT ON COLUMN "public".mch_online_shop.create_time IS '创建时间';
COMMENT ON COLUMN "public".mch_sale_conf.fx_sales IS '是否启用分销';
COMMENT ON COLUMN "public".mch_sale_conf.cb_percent IS '反现比例,0则不返现';
COMMENT ON COLUMN "public".mch_sale_conf.cb_tg1_percent IS '一级比例';
COMMENT ON COLUMN "public".mch_sale_conf.cb_tg2_percent IS '二级比例';
COMMENT ON COLUMN "public".mch_sale_conf.cb_member_percent IS '会员比例';
COMMENT ON COLUMN "public".mch_sale_conf.oa_open IS '开启自动设置订单';
COMMENT ON COLUMN "public".mch_sale_conf.oa_timeout_minute IS '订单超时取消（分钟）';
COMMENT ON COLUMN "public".mch_sale_conf.oa_confirm_minute IS '订单自动确认（分钟）';
COMMENT ON COLUMN "public".mch_sale_conf.oa_receive_hour IS '超时自动收货（小时）';
COMMENT ON COLUMN "public".pf_riseinfo.settlement_amount IS '结算金额';
COMMENT ON COLUMN "public".pm_cash_back.data_tag IS '自定义数据';
COMMENT ON COLUMN "public".pm_coupon.code IS '优惠码';
COMMENT ON COLUMN "public".pm_coupon.amount IS '优惠码可用数量';
COMMENT ON COLUMN "public".pm_coupon.fee IS '包含金额';
COMMENT ON COLUMN "public".pm_coupon.integral IS '包含积分';
COMMENT ON COLUMN "public".pm_coupon.min_level IS '等级限制';
COMMENT ON COLUMN "public".pm_coupon.min_fee IS '订单金额限制';
COMMENT ON COLUMN "public".pm_coupon.need_bind IS '是否需要绑定';
COMMENT ON COLUMN "public".pm_coupon_bind.member_id IS '会员编号';
COMMENT ON COLUMN "public".pm_coupon_bind.coupon_id IS '优惠券编号';
COMMENT ON COLUMN "public".pm_coupon_bind.is_used IS '是否使用';
COMMENT ON COLUMN "public".pm_coupon_bind.bind_time IS '绑定时间';
COMMENT ON COLUMN "public".pm_coupon_bind.use_time IS '使用时间';
COMMENT ON COLUMN "public".pm_coupon_take.is_apply IS '是否生效,1表示有效';
COMMENT ON COLUMN "public".pm_coupon_take.take_time IS '占用时间';
COMMENT ON COLUMN "public".pm_coupon_take.extra_time IS '释放时间,超过该时间，优惠券释放';
COMMENT ON COLUMN "public".pm_coupon_take.apply_time IS '更新时间';
COMMENT ON COLUMN "public".pt_member_level.value IS '等级值';
COMMENT ON COLUMN "public".pt_member_level.require_exp IS '要求积分';
COMMENT ON COLUMN "public".sale_after_order.image_url IS '商品售后图片凭证';
COMMENT ON COLUMN "public".sale_after_order.rsp_name IS '退货快递名称';
COMMENT ON COLUMN "public".sale_after_order.rsp_order IS '退货快递单号';
COMMENT ON COLUMN "public".sale_order_log.type IS '类型，１:流程,2:调价';
COMMENT ON COLUMN "public".ship_order.id IS '编号';
COMMENT ON COLUMN "public".ship_order.order_id IS '订单编号';
COMMENT ON COLUMN "public".ship_order.sub_orderid IS '子订单编号';
COMMENT ON COLUMN "public".ship_order.sp_id IS '快递SP编号';
COMMENT ON COLUMN "public".ship_order.sp_order IS '快递SP单号';
COMMENT ON COLUMN "public".ship_order.shipment_log IS '物流日志';
COMMENT ON COLUMN "public".ship_order.amount IS '运费';
COMMENT ON COLUMN "public".ship_order.final_amount IS '实际运费';
COMMENT ON COLUMN "public".ship_order.ship_time IS '发货时间';
COMMENT ON COLUMN "public".ship_order.state IS '状态';
COMMENT ON COLUMN "public".ship_order.update_time IS '更新时间';
COMMENT ON COLUMN "public".usr_credential.sign IS '标记凭据类型';
ALTER TABLE "public".item_sku 
  ALTER COLUMN image SET DATA TYPE varchar(256);
COMMENT ON COLUMN "public".item_sku.id IS '编号';
COMMENT ON COLUMN "public".item_sku.product_id IS '产品编号';
COMMENT ON COLUMN "public".item_sku.item_id IS '商品编号';
COMMENT ON COLUMN "public".item_sku.title IS '标题';
COMMENT ON COLUMN "public".item_sku.image IS '图片';
COMMENT ON COLUMN "public".item_sku.spec_data IS '规格数据';
COMMENT ON COLUMN "public".item_sku.spec_word IS '规格字符';
COMMENT ON COLUMN "public".item_sku.retail_price IS '参考价';
COMMENT ON COLUMN "public".item_sku.code IS '产品编码';
COMMENT ON COLUMN "public".item_sku.price IS '价格（分)';
COMMENT ON COLUMN "public".item_sku.cost IS '成本（分)';
COMMENT ON COLUMN "public".item_sku.weight IS '重量(克)';
COMMENT ON COLUMN "public".item_sku."bulk" IS '体积（毫升)';
COMMENT ON COLUMN "public".item_sku.stock IS '库存';
COMMENT ON COLUMN "public".item_sku.sale_num IS '已销售数量';
COMMENT ON COLUMN "public".portal_nav.id IS '编号';
COMMENT ON COLUMN "public".portal_nav.text IS '文本';
COMMENT ON COLUMN "public".portal_nav.url IS '地址';
COMMENT ON COLUMN "public".portal_nav.target IS '打开目标';
COMMENT ON COLUMN "public".portal_nav.image IS '链接图片';
COMMENT ON COLUMN "public".portal_nav.nav_type IS '导航类型: 1为电脑，2为手机端';
COMMENT ON COLUMN "public".portal_nav_type.id IS '编号';
COMMENT ON COLUMN "public".portal_nav_type.name IS '名称';
COMMENT ON COLUMN "public".sys_kv.id IS '编号';
COMMENT ON COLUMN "public".sys_kv."key" IS '键';
COMMENT ON COLUMN "public".sys_kv.value IS '值';
COMMENT ON COLUMN "public".sys_kv.update_time IS '更新时间';
COMMENT ON COLUMN "public".portal_floor_ad.id IS '编号';
COMMENT ON COLUMN "public".portal_floor_ad.cat_id IS '分类编号';
COMMENT ON COLUMN "public".portal_floor_ad.pos_id IS '广告位编号';
COMMENT ON COLUMN "public".portal_floor_ad.ad_index IS '广告顺序';
COMMENT ON COLUMN "public".portal_floor_link.id IS '编号';
COMMENT ON COLUMN "public".portal_floor_link.cat_id IS '分类编号';
COMMENT ON COLUMN "public".portal_floor_link.text IS '文本';
COMMENT ON COLUMN "public".portal_floor_link.link_url IS '链接地址';
COMMENT ON COLUMN "public".portal_floor_link.target IS '打开方式';
COMMENT ON COLUMN "public".portal_floor_link.sort_num IS '序号';
COMMENT ON COLUMN "public".mm_buyer_group.id IS '编号';
COMMENT ON COLUMN "public".mm_buyer_group.name IS '名称';
COMMENT ON COLUMN "public".mm_buyer_group.is_default IS '是否为默认分组,未设置分组的客户作为该分组。';
COMMENT ON COLUMN "public".mch_buyer_group.id IS '编号';
COMMENT ON COLUMN "public".mch_buyer_group.mch_id IS '商家编号';
COMMENT ON COLUMN "public".mch_buyer_group.group_id IS '客户分组编号';
COMMENT ON COLUMN "public".mch_buyer_group.alias IS '分组别名';
COMMENT ON COLUMN "public".mch_buyer_group.enable_retail IS '是否启用零售';
COMMENT ON COLUMN "public".mch_buyer_group.enable_wholesale IS '是否启用批发';
COMMENT ON COLUMN "public".mch_buyer_group.rebate_period IS '批发返点周期';
COMMENT ON COLUMN "public".ws_cart.id IS '编号';
COMMENT ON COLUMN "public".ws_cart.code IS '购物车编码';
COMMENT ON COLUMN "public".ws_cart.buyer_id IS '买家编号';
COMMENT ON COLUMN "public".ws_cart.deliver_id IS '送货地址';
COMMENT ON COLUMN "public".ws_cart.payment_opt IS '支付选项';
COMMENT ON COLUMN "public".ws_cart.create_time IS '创建时间';
COMMENT ON COLUMN "public".ws_cart.update_time IS '修改时间';
COMMENT ON COLUMN "public".ws_cart_item.id IS '编号';
COMMENT ON COLUMN "public".ws_cart_item.cart_id IS '购物车编号';
COMMENT ON COLUMN "public".ws_cart_item.vendor_id IS '运营商编号';
COMMENT ON COLUMN "public".ws_cart_item.shop_id IS '店铺编号';
COMMENT ON COLUMN "public".ws_cart_item.item_id IS '商品编号';
COMMENT ON COLUMN "public".ws_cart_item.sku_id IS 'SKU编号';
COMMENT ON COLUMN "public".ws_cart_item.quantity IS '数量';
COMMENT ON COLUMN "public".ws_cart_item.checked IS '是否勾选结算';
COMMENT ON COLUMN "public".ship_item.id IS '编号';
COMMENT ON COLUMN "public".ship_item.ship_order IS '发货单编号';
COMMENT ON COLUMN "public".ship_item.snapshot_id IS '商品交易快照编号';
COMMENT ON COLUMN "public".ship_item.quantity IS '商品数量';
COMMENT ON COLUMN "public".ship_item.amount IS '运费';
COMMENT ON COLUMN "public".ship_item.final_amount IS '实际运费';
COMMENT ON COLUMN "public".order_trade_order.id IS '编号';
COMMENT ON COLUMN "public".order_trade_order.order_id IS '订单编号';
COMMENT ON COLUMN "public".order_trade_order.vendor_id IS '商家编号';
COMMENT ON COLUMN "public".order_trade_order.shop_id IS '店铺编号';
COMMENT ON COLUMN "public".order_trade_order.subject IS '订单标题';
COMMENT ON COLUMN "public".order_trade_order.order_amount IS '订单金额';
COMMENT ON COLUMN "public".order_trade_order.discount_amount IS '抵扣金额';;
COMMENT ON COLUMN "public".order_trade_order.final_amount IS '订单最终金额';
COMMENT ON COLUMN "public".order_trade_order.trade_rate IS '交易结算比例（商户)';
COMMENT ON COLUMN "public".order_trade_order.cash_pay IS '是否现金支付';
COMMENT ON COLUMN "public".order_trade_order.ticket_image IS '发票图片';
COMMENT ON COLUMN "public".order_trade_order.remark IS '订单备注';
COMMENT ON COLUMN "public".order_trade_order.state IS '订单状态';
COMMENT ON COLUMN "public".order_trade_order.create_time IS '订单创建时间';
COMMENT ON COLUMN "public".order_trade_order.update_time IS '订单更新时间';
COMMENT ON TABLE "public".mch_merchant IS '商户';
COMMENT ON COLUMN "public".mch_merchant.member_id IS '会员编号';
COMMENT ON COLUMN "public".mch_merchant.login_user IS '登录用户';
COMMENT ON COLUMN "public".mch_merchant.login_pwd IS '登录密码';
COMMENT ON COLUMN "public".mch_merchant.name IS '名称';
COMMENT ON COLUMN "public".mch_merchant.company_name IS '公司名称';
COMMENT ON COLUMN "public".mch_merchant.self_sales IS '是否字营';
COMMENT ON COLUMN "public".mch_merchant.level IS '商户等级';
COMMENT ON COLUMN "public".mch_merchant.logo IS '标志';
COMMENT ON COLUMN "public".mch_merchant.province IS '省';
COMMENT ON COLUMN "public".mch_merchant.city IS '市';
COMMENT ON COLUMN "public".mch_merchant.district IS '区';
COMMENT ON COLUMN "public".mch_merchant.create_time IS '创建时间';
COMMENT ON COLUMN "public".mch_merchant.flag IS '标志';
COMMENT ON COLUMN "public".mch_merchant.enabled IS '是否启用';
COMMENT ON COLUMN "public".mch_merchant.expires_time IS '过期时间';
COMMENT ON COLUMN "public".mch_merchant.update_time IS '更新时间';
COMMENT ON COLUMN "public".mch_merchant.login_time IS '登录时间';
COMMENT ON COLUMN "public".mch_merchant.last_login_time IS '最后登录时间';
COMMENT ON COLUMN "public".mch_shop.id IS '商店编号';
COMMENT ON COLUMN "public".mch_shop.vendor_id IS '商户编号';
COMMENT ON COLUMN "public".mch_shop.name IS '商店名称';
COMMENT ON COLUMN "public".mch_shop.sort_num IS '排序序号';
COMMENT ON COLUMN "public".mch_shop.shop_type IS '店铺类型';
COMMENT ON COLUMN "public".mch_shop.opening_state IS '营业状态';
COMMENT ON COLUMN "public".mch_shop.state IS '状态 1:表示正常,2:表示关闭';
COMMENT ON TABLE "public".wal_wallet IS '钱包';
COMMENT ON COLUMN "public".wal_wallet.id IS '编号';
COMMENT ON COLUMN "public".wal_wallet.hash_code IS '哈希值';
COMMENT ON COLUMN "public".wal_wallet.node_id IS '节点编号';
COMMENT ON COLUMN "public".wal_wallet.user_id IS '用户编号';
COMMENT ON COLUMN "public".wal_wallet.wallet_type IS '钱包类型';
COMMENT ON COLUMN "public".wal_wallet.wallet_flag IS '钱包标志';
COMMENT ON COLUMN "public".wal_wallet.balance IS '余额';
COMMENT ON COLUMN "public".wal_wallet.present_balance IS '赠送余额';
COMMENT ON COLUMN "public".wal_wallet.adjust_amount IS '调整禁遏';
COMMENT ON COLUMN "public".wal_wallet.freeze_amount IS '冻结金额';
COMMENT ON COLUMN "public".wal_wallet.latest_amount IS '结余金额';
COMMENT ON COLUMN "public".wal_wallet.expired_amount IS '失效账户余额';
COMMENT ON COLUMN "public".wal_wallet.total_charge IS '总充值金额';
COMMENT ON COLUMN "public".wal_wallet.total_present IS '累计赠送金额';
COMMENT ON COLUMN "public".wal_wallet.total_pay IS '总支付额';
COMMENT ON COLUMN "public".wal_wallet.remark IS '备注';
COMMENT ON COLUMN "public".wal_wallet.state IS '状态';
COMMENT ON COLUMN "public".wal_wallet.create_time IS '创建时间';
COMMENT ON COLUMN "public".wal_wallet.update_time IS '更新时间';
COMMENT ON COLUMN "public".wal_wallet_log.id IS '编号';
COMMENT ON COLUMN "public".wal_wallet_log.wallet_id IS '钱包编号';
COMMENT ON COLUMN "public".wal_wallet_log.kind IS '业务类型';
COMMENT ON COLUMN "public".wal_wallet_log.title IS '标题';
COMMENT ON COLUMN "public".wal_wallet_log.outer_chan IS '外部通道';
COMMENT ON COLUMN "public".wal_wallet_log.outer_no IS '外部订单号';
COMMENT ON COLUMN "public".wal_wallet_log.value IS '变动金额';
COMMENT ON COLUMN "public".wal_wallet_log.balance IS '余额';
COMMENT ON COLUMN "public".wal_wallet_log.trade_fee IS '交易手续费';
COMMENT ON COLUMN "public".wal_wallet_log.op_uid IS '操作人员用户编号';
COMMENT ON COLUMN "public".wal_wallet_log.op_name IS '操作人员名称';
COMMENT ON COLUMN "public".wal_wallet_log.remark IS '备注';
COMMENT ON COLUMN "public".wal_wallet_log.review_state IS '审核状态';
COMMENT ON COLUMN "public".wal_wallet_log.review_remark IS '审核备注';
COMMENT ON COLUMN "public".wal_wallet_log.review_time IS '审核时间';
COMMENT ON COLUMN "public".wal_wallet_log.create_time IS '创建时间';
COMMENT ON COLUMN "public".wal_wallet_log.update_time IS '更新时间';
COMMENT ON COLUMN "public".mch_trade_conf.id IS '编号';
COMMENT ON COLUMN "public".mch_trade_conf.mch_id IS '商户编号';
COMMENT ON COLUMN "public".mch_trade_conf.trade_type IS '交易类型';
COMMENT ON COLUMN "public".mch_trade_conf.plan_id IS '交易方案，根据方案来自动调整比例';
COMMENT ON COLUMN "public".mch_trade_conf.flag IS '交易标志';
COMMENT ON COLUMN "public".mch_trade_conf.amount_basis IS '交易手续费依据,1:未设置 2:按金额 3:按比例';
COMMENT ON COLUMN "public".mch_trade_conf.trade_fee IS '交易费，按单笔收取';
COMMENT ON COLUMN "public".mch_trade_conf.trade_rate IS '交易手续费比例';
COMMENT ON COLUMN "public".mch_trade_conf.update_time IS '更新时间';
COMMENT ON TABLE "public".registry IS '注册表';
COMMENT ON COLUMN "public".registry."key" IS '键';
COMMENT ON COLUMN "public".registry.description IS '描述';
COMMENT ON COLUMN "public".registry.value IS '值';
COMMENT ON COLUMN "public".registry.flag IS '是否用户定义,0:否,1:是';
COMMENT ON COLUMN "public".registry.default_value IS '默认值';
COMMENT ON COLUMN "public".registry.options IS '可选值';
COMMENT ON TABLE "public".mm_lock_info IS '会员锁定记录';
COMMENT ON COLUMN "public".mm_lock_info.id IS '编号';
COMMENT ON COLUMN "public".mm_lock_info.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_lock_info.lock_time IS '锁定时间';
COMMENT ON COLUMN "public".mm_lock_info.unlock_time IS '解锁时间';
COMMENT ON COLUMN "public".mm_lock_info.remark IS '备注';
COMMENT ON TABLE "public".mm_lock_history IS '会员锁定历史';
COMMENT ON COLUMN "public".mm_lock_history.id IS '编号';
COMMENT ON COLUMN "public".mm_lock_history.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_lock_history.lock_time IS '锁定时间';
COMMENT ON COLUMN "public".mm_lock_history.duration IS '锁定持续分钟数';
COMMENT ON COLUMN "public".mm_lock_history.remark IS '备注';
ALTER TABLE "public".item_trade_snapshot 
  ALTER COLUMN img SET DATA TYPE varchar(256);
COMMENT ON COLUMN "public".item_trade_snapshot.cost IS '供货价';
ALTER TABLE "public".item_snapshot 
  ALTER COLUMN image SET DATA TYPE varchar(256);
COMMENT ON COLUMN "public".item_snapshot.item_id IS '商品编号';
COMMENT ON COLUMN "public".item_snapshot.product_id IS '产品编号';
COMMENT ON COLUMN "public".item_snapshot.snapshot_key IS '快照编码';
COMMENT ON COLUMN "public".item_snapshot.cat_id IS '分类编号';
COMMENT ON COLUMN "public".item_snapshot.vendor_id IS '供货商编号';
COMMENT ON COLUMN "public".item_snapshot.brand_id IS '编号';
COMMENT ON COLUMN "public".item_snapshot.shop_id IS '商铺编号';
COMMENT ON COLUMN "public".item_snapshot.shop_cat_id IS '编号分类编号';
COMMENT ON COLUMN "public".item_snapshot.express_tid IS '运费模板';
COMMENT ON COLUMN "public".item_snapshot.title IS '商品标题';
COMMENT ON COLUMN "public".item_snapshot.short_title IS '短标题';
COMMENT ON COLUMN "public".item_snapshot.code IS '商户编码';
COMMENT ON COLUMN "public".item_snapshot.image IS '商品图片';
COMMENT ON COLUMN "public".item_snapshot.is_present IS '是否为赠品';
COMMENT ON COLUMN "public".item_snapshot.price_range IS '价格区间';
COMMENT ON COLUMN "public".item_snapshot.sku_id IS '默认SKU';
COMMENT ON COLUMN "public".item_snapshot.cost IS '成本';
COMMENT ON COLUMN "public".item_snapshot.price IS '售价';
COMMENT ON COLUMN "public".item_snapshot.retail_price IS '零售价';
COMMENT ON COLUMN "public".item_snapshot.weight IS '重量(g)';
COMMENT ON COLUMN "public".item_snapshot."bulk" IS '体积(ml)';
COMMENT ON COLUMN "public".item_snapshot.level_sales IS '会员价';
COMMENT ON COLUMN "public".item_snapshot.shelve_state IS '上架状态';
COMMENT ON COLUMN "public".item_snapshot.update_time IS '更新时间';
COMMENT ON COLUMN "public".sale_cart_item.id IS '编号';
COMMENT ON COLUMN "public".sale_cart_item.cart_id IS '购物车编号';
COMMENT ON COLUMN "public".sale_cart_item.vendor_id IS '运营商编号';
COMMENT ON COLUMN "public".sale_cart_item.shop_id IS '店铺编号';
COMMENT ON COLUMN "public".sale_cart_item.item_id IS '商品编号';
COMMENT ON COLUMN "public".sale_cart_item.sku_id IS 'SKU编号';
COMMENT ON COLUMN "public".sale_cart_item.quantity IS '数量';
COMMENT ON COLUMN "public".sale_cart_item.checked IS '是否勾选结算';
COMMENT ON COLUMN "public".sale_cart.id IS '编号';
COMMENT ON COLUMN "public".sale_cart.code IS '购物车编码';
COMMENT ON COLUMN "public".sale_cart.buyer_id IS '买家编号';
COMMENT ON COLUMN "public".sale_cart.deliver_id IS '送货地址';
COMMENT ON COLUMN "public".sale_cart.payment_opt IS '支付选项';
COMMENT ON COLUMN "public".sale_cart.create_time IS '创建时间';
COMMENT ON COLUMN "public".sale_cart.update_time IS '修改时间';
COMMENT ON COLUMN "public".sale_order.id IS '编号';
COMMENT ON COLUMN "public".sale_order.order_id IS '订单编号';
COMMENT ON COLUMN "public".sale_order.item_amount IS '商品金额';
COMMENT ON COLUMN "public".sale_order.discount_amount IS '抵扣金额';
COMMENT ON COLUMN "public".sale_order.express_fee IS '物流费';
COMMENT ON COLUMN "public".sale_order.package_fee IS '包装费';
COMMENT ON COLUMN "public".sale_order.final_amount IS '订单最终金额';
COMMENT ON COLUMN "public".sale_order.consignee_person IS '收货人姓名';
COMMENT ON COLUMN "public".sale_order.consignee_phone IS '收货人电话';
COMMENT ON COLUMN "public".sale_order.shipping_address IS '收货人地址';
COMMENT ON COLUMN "public".sale_order.is_break IS '是否拆分';
COMMENT ON COLUMN "public".sale_order.update_time IS '更新时间';
ALTER TABLE "public".item_info 
  ALTER COLUMN image SET DATA TYPE varchar(256);
COMMENT ON TABLE "public".product_model_attr IS '模型商品属性';
COMMENT ON COLUMN "public".product_model_attr.id IS '编号';
COMMENT ON COLUMN "public".product_model_attr.pro_model IS '商品模型';
COMMENT ON COLUMN "public".product_model_attr.name IS '属性名称';
COMMENT ON COLUMN "public".product_model_attr.is_filter IS '是否作为筛选条件';
COMMENT ON COLUMN "public".product_model_attr.multi_chk IS '是否多选';
COMMENT ON COLUMN "public".product_model_attr.item_values IS '属性项值';
COMMENT ON COLUMN "public".product_model_attr.sort_num IS '排列序号';
COMMENT ON TABLE "public".product_attr_info IS '产品属性信息';
COMMENT ON COLUMN "public".product_attr_info.id IS '编号';
COMMENT ON COLUMN "public".product_attr_info.product_id IS '产品编号';
COMMENT ON COLUMN "public".product_attr_info.attr_id IS '属性编号';
COMMENT ON COLUMN "public".product_attr_info.attr_data IS '属性数据';
COMMENT ON COLUMN "public".product_attr_info.attr_word IS '属性文本';
COMMENT ON TABLE "public".product_model_attr_item IS '商品模型属性项';
COMMENT ON COLUMN "public".product_model_attr_item.id IS '编号';
COMMENT ON COLUMN "public".product_model_attr_item.attr_id IS '属性编号';
COMMENT ON COLUMN "public".product_model_attr_item.pro_model IS '商品模型';
COMMENT ON COLUMN "public".product_model_attr_item.value IS '属性值';
COMMENT ON COLUMN "public".product_model_attr_item.sort_num IS '排列序号';
COMMENT ON TABLE "public".product_brand IS '品牌';
COMMENT ON COLUMN "public".product_brand.id IS '编号';
COMMENT ON COLUMN "public".product_brand.name IS '品牌名称';
COMMENT ON COLUMN "public".product_brand.image IS '品牌图片';
COMMENT ON COLUMN "public".product_brand.site_url IS '品牌官网';
COMMENT ON COLUMN "public".product_brand.intro IS '品牌介绍';
COMMENT ON COLUMN "public".product_brand.review_state IS '审核';
COMMENT ON COLUMN "public".product_brand.create_time IS '创建时间';
COMMENT ON TABLE "public".product_category IS '产品分类';
COMMENT ON COLUMN "public".product_category.id IS '编号';
COMMENT ON COLUMN "public".product_category.parent_id IS '上级分类';
COMMENT ON COLUMN "public".product_category.prod_model IS '产品模型';
COMMENT ON COLUMN "public".product_category.priority IS '优先级';
COMMENT ON COLUMN "public".product_category.name IS '分类名称';
COMMENT ON COLUMN "public".product_category.virtual_cat IS '是否为虚拟分类';
COMMENT ON COLUMN "public".product_category.cat_url IS '分类链接地址';
COMMENT ON COLUMN "public".product_category.icon IS '图标';
COMMENT ON COLUMN "public".product_category.icon_xy IS '图标坐标';
COMMENT ON COLUMN "public".product_category.level IS '分类层级';
COMMENT ON COLUMN "public".product_category.sort_num IS '序号';
COMMENT ON COLUMN "public".product_category.floor_show IS '是否楼层显示';
COMMENT ON COLUMN "public".product_category.enabled IS '是否启用';
COMMENT ON COLUMN "public".product_category.create_time IS '创建时间';
COMMENT ON TABLE "public".product_model IS '产品模型';
COMMENT ON COLUMN "public".product_model.id IS '编号';
COMMENT ON COLUMN "public".product_model.name IS '模型名称';
COMMENT ON COLUMN "public".product_model.attr_str IS '属性值';
COMMENT ON COLUMN "public".product_model.spec_str IS '规格值';
COMMENT ON COLUMN "public".product_model.enabled IS '是否启用';
COMMENT ON TABLE "public".product_model_brand IS '商品模型关联品牌';
COMMENT ON COLUMN "public".product_model_brand.id IS '编号';
COMMENT ON COLUMN "public".product_model_brand.brand_id IS '品牌编号';
COMMENT ON COLUMN "public".product_model_brand.pro_model IS '商品模型';
ALTER TABLE "public".product 
  ALTER COLUMN img SET DATA TYPE varchar(256);
COMMENT ON TABLE "public".product IS '产品';
COMMENT ON COLUMN "public".product.id IS '编号';
COMMENT ON COLUMN "public".product.cat_id IS '分类编号';
COMMENT ON COLUMN "public".product.supplier_id IS '供货商编号';
COMMENT ON COLUMN "public".product.brand_id IS '品牌编号';
COMMENT ON COLUMN "public".product.name IS '名称';
COMMENT ON COLUMN "public".product.code IS '产品编码';
COMMENT ON COLUMN "public".product.img IS '产品图片';
COMMENT ON COLUMN "public".product.sort_num IS '排序编号';
COMMENT ON COLUMN "public".product.description IS '描述';
COMMENT ON COLUMN "public".product.state IS '产品状态';
COMMENT ON COLUMN "public".product.remark IS '备注';
COMMENT ON COLUMN "public".product.create_time IS '创建时间';
COMMENT ON COLUMN "public".product.update_time IS '更新时间';
COMMENT ON TABLE "public".product_model_spec IS '商品模型规格';
COMMENT ON COLUMN "public".product_model_spec.id IS '编号';
COMMENT ON COLUMN "public".product_model_spec.pro_model IS '商品模型';
COMMENT ON COLUMN "public".product_model_spec.name IS '规格名称';
COMMENT ON COLUMN "public".product_model_spec.item_values IS '规格项值';
COMMENT ON COLUMN "public".product_model_spec.sort_num IS '排列序号';
COMMENT ON TABLE "public".product_model_spec_item IS '产品模型规格项';
COMMENT ON COLUMN "public".product_model_spec_item.id IS '编号';
COMMENT ON COLUMN "public".product_model_spec_item.spec_id IS '规格编号';
COMMENT ON COLUMN "public".product_model_spec_item.pro_model IS '商品模型';
COMMENT ON COLUMN "public".product_model_spec_item.value IS '规格值';
COMMENT ON COLUMN "public".product_model_spec_item.color IS '规格颜色';
COMMENT ON COLUMN "public".product_model_spec_item.sort_num IS '排列序号';
ALTER TABLE "public".ws_item 
  alter column price_range set default ''::character varying;
COMMENT ON TABLE "public".mm_balance_log IS '余额日志';
COMMENT ON COLUMN "public".mm_balance_log.id IS '编号';
COMMENT ON COLUMN "public".mm_balance_log.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_balance_log.kind IS '类型';
COMMENT ON COLUMN "public".mm_balance_log.title IS '标题';
COMMENT ON COLUMN "public".mm_balance_log.outer_no IS '外部交易号';
COMMENT ON COLUMN "public".mm_balance_log.amount IS '金额';
COMMENT ON COLUMN "public".mm_balance_log.csn_fee IS '手续费';
COMMENT ON COLUMN "public".mm_balance_log.review_state IS '审核状态';
COMMENT ON COLUMN "public".mm_balance_log.rel_user IS '关联用户';
COMMENT ON COLUMN "public".mm_balance_log.remark IS '备注';
COMMENT ON COLUMN "public".mm_balance_log.create_time IS '创建时间';
COMMENT ON COLUMN "public".mm_balance_log.update_time IS '更新时间';
COMMENT ON TABLE "public".mm_deliver_addr IS '会员收货地址';
COMMENT ON COLUMN "public".mm_deliver_addr.id IS '编号';
COMMENT ON COLUMN "public".mm_deliver_addr.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_deliver_addr.consignee_name IS '收货人姓名';
COMMENT ON COLUMN "public".mm_deliver_addr.consignee_phone IS '收货人电话';
COMMENT ON COLUMN "public".mm_deliver_addr.province IS '数字编码(省)';
COMMENT ON TABLE "public".mm_favorite IS '会员收藏';
COMMENT ON COLUMN "public".mm_favorite.id IS '编号';
COMMENT ON COLUMN "public".mm_favorite.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_favorite.fav_type IS '收藏类型';
COMMENT ON COLUMN "public".mm_favorite.refer_id IS '关联编号';
COMMENT ON COLUMN "public".mm_favorite.create_time IS '收藏时间';
COMMENT ON TABLE "public".mm_flow_log IS '活动账户明细';
COMMENT ON COLUMN "public".mm_flow_log.id IS '编号';
COMMENT ON COLUMN "public".mm_flow_log.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_flow_log.kind IS '类型';
COMMENT ON COLUMN "public".mm_flow_log.title IS '标题';
COMMENT ON COLUMN "public".mm_flow_log.outer_no IS '外部交易号';
COMMENT ON COLUMN "public".mm_flow_log.amount IS '金额';
COMMENT ON COLUMN "public".mm_flow_log.csn_fee IS '手续费';
COMMENT ON COLUMN "public".mm_flow_log.review_state IS '审核状态';
COMMENT ON COLUMN "public".mm_flow_log.rel_user IS '关联用户';
COMMENT ON COLUMN "public".mm_flow_log.remark IS '备注';
COMMENT ON COLUMN "public".mm_flow_log.create_time IS '创建时间';
COMMENT ON COLUMN "public".mm_flow_log.update_time IS '更新时间';
ALTER TABLE "public".mm_integral_log 
  alter column title set default '""'::character varying;
ALTER TABLE "public".mm_integral_log 
  alter column outer_no set default ''::character varying;
COMMENT ON TABLE "public".mm_integral_log IS '积分明细';
COMMENT ON COLUMN "public".mm_integral_log.id IS '编号';
COMMENT ON COLUMN "public".mm_integral_log.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_integral_log.kind IS '类型';
COMMENT ON COLUMN "public".mm_integral_log.title IS '标题';
COMMENT ON COLUMN "public".mm_integral_log.outer_no IS '关联的编号';
COMMENT ON COLUMN "public".mm_integral_log.value IS '积分值';
COMMENT ON COLUMN "public".mm_integral_log.remark IS '备注';
COMMENT ON COLUMN "public".mm_integral_log.rel_user IS '关联用户';
COMMENT ON COLUMN "public".mm_integral_log.review_state IS '审核状态';
COMMENT ON COLUMN "public".mm_integral_log.create_time IS '创建时间';
COMMENT ON COLUMN "public".mm_integral_log.update_time IS '更新时间';
COMMENT ON TABLE "public".mm_levelup IS '会员升级日志表';
COMMENT ON COLUMN "public".mm_levelup.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_levelup.origin_level IS '原来等级';
COMMENT ON COLUMN "public".mm_levelup.target_level IS '现在等级';
COMMENT ON COLUMN "public".mm_levelup.is_free IS '是否为免费升级的会员';
COMMENT ON COLUMN "public".mm_levelup.payment_id IS '支付单编号';
COMMENT ON COLUMN "public".mm_levelup.create_time IS '升级时间';
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
COMMENT ON COLUMN "public".mm_member."user" IS '用户名';
COMMENT ON COLUMN "public".mm_member.flag IS '会员标志';
COMMENT ON COLUMN "public".mm_member.name IS '昵称';
COMMENT ON COLUMN "public".mm_member.real_name IS '真实姓名';
COMMENT ON TABLE "public".mm_receipts_code IS '收款码';
COMMENT ON COLUMN "public".mm_receipts_code.id IS '编号';
COMMENT ON COLUMN "public".mm_receipts_code.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_receipts_code."identity" IS '账户标识,如:alipay';
COMMENT ON COLUMN "public".mm_receipts_code.name IS '账户名称';
COMMENT ON COLUMN "public".mm_receipts_code.account_id IS '账号';
COMMENT ON COLUMN "public".mm_receipts_code.code_url IS '收款码地址';
COMMENT ON COLUMN "public".mm_receipts_code.state IS '是否启用';
COMMENT ON COLUMN "public".mm_relation.inviter_id IS '邀请会员编号';
COMMENT ON COLUMN "public".mm_relation.inviter_d2 IS '邀请会员编号(depth2)';
COMMENT ON COLUMN "public".mm_relation.inviter_d3 IS '邀请会员编号(depth3)';
ALTER TABLE "public".mm_trusted_info 
  alter column card_reverse_image set default ''::character varying;
COMMENT ON COLUMN "public".mm_trusted_info.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_trusted_info.real_name IS '真实姓名';
COMMENT ON COLUMN "public".mm_trusted_info.country_code IS '国家代码';
COMMENT ON COLUMN "public".mm_trusted_info.card_type IS '证件类型';
COMMENT ON COLUMN "public".mm_trusted_info.card_id IS '证件编号';
COMMENT ON COLUMN "public".mm_trusted_info.card_image IS '证件图片';
COMMENT ON COLUMN "public".mm_trusted_info.card_reverse_image IS '证件反面图片';
COMMENT ON COLUMN "public".mm_trusted_info.trust_image IS '认证图片,人与身份证的图像等';
COMMENT ON COLUMN "public".mm_trusted_info.manual_review IS '人工审核';
COMMENT ON COLUMN "public".mm_trusted_info.review_state IS '审核状态';
COMMENT ON COLUMN "public".mm_trusted_info.review_time IS '审核时间';
COMMENT ON COLUMN "public".mm_trusted_info.remark IS '备注';
COMMENT ON COLUMN "public".mm_trusted_info.update_time IS '更新时间';
COMMENT ON TABLE "public".mm_wallet_log IS '钱包日志';
COMMENT ON COLUMN "public".mm_wallet_log.id IS '编号';
COMMENT ON COLUMN "public".mm_wallet_log.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_wallet_log.kind IS '类型';
COMMENT ON COLUMN "public".mm_wallet_log.title IS '标题';
COMMENT ON COLUMN "public".mm_wallet_log.outer_no IS '外部交易号';
COMMENT ON COLUMN "public".mm_wallet_log.amount IS '金额';
COMMENT ON COLUMN "public".mm_wallet_log.csn_fee IS '手续费';
COMMENT ON COLUMN "public".mm_wallet_log.review_state IS '审核状态';
COMMENT ON COLUMN "public".mm_wallet_log.rel_user IS '关联用户';
COMMENT ON COLUMN "public".mm_wallet_log.remark IS '备注';
COMMENT ON COLUMN "public".mm_wallet_log.create_time IS '创建时间';
COMMENT ON COLUMN "public".mm_wallet_log.update_time IS '更新时间';
COMMENT ON COLUMN "public".shop_site_conf.index_title IS '首页标题';
COMMENT ON COLUMN "public".shop_site_conf.sub_title IS '子页面标题';
COMMENT ON COLUMN "public".shop_site_conf.state IS '状态: 0:暂停  1：正常';
COMMENT ON TABLE "public".mch_express_template IS '商户运费模板';
COMMENT ON COLUMN "public".mch_express_template.id IS '编号';
COMMENT ON COLUMN "public".mch_express_template.vendor_id IS '运营商编号';
COMMENT ON COLUMN "public".mch_express_template.name IS '运费模板名称';
COMMENT ON COLUMN "public".mch_express_template.is_free IS '是否卖价承担运费';
COMMENT ON COLUMN "public".mch_express_template.basis IS '运费计价依据';
COMMENT ON COLUMN "public".mch_express_template.first_unit IS '首次计价单位,如首重为2kg';
COMMENT ON COLUMN "public".mch_express_template.first_fee IS '首次计价单价,如续重1kg';
COMMENT ON COLUMN "public".mch_express_template.add_unit IS '超过首次计价计算单位,如续重1kg';
COMMENT ON COLUMN "public".mch_express_template.add_fee IS '超过首次计价单价，如续重1kg';
COMMENT ON COLUMN "public".mch_express_template.enabled IS '是否启用';


-- 银行卡 --
DROP TABLE mm_bank_card;
CREATE TABLE "public".mm_bank_card (
   id           BIGSERIAL NOT NULL,
   member_id    int8 NOT NULL,
   bank_account varchar(20) NOT NULL,
   account_name varchar(20) NOT NULL,
   bank_id      int4 NOT NULL,
   bank_name    varchar(45) NOT NULL,
   bank_code    varchar(10) NOT NULL,
   network      varchar(20) NOT NULL,
   auth_code    varchar(40) NOT NULL,
   state        int2 NOT NULL,
   create_time  int8 NOT NULL,
   CONSTRAINT mm_bank_card_pkey
       PRIMARY KEY (id));
COMMENT ON TABLE "public".mm_bank_card IS '银行卡';
COMMENT ON COLUMN "public".mm_bank_card.id IS '编号';
COMMENT ON COLUMN "public".mm_bank_card.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_bank_card.bank_account IS '银行账号';
COMMENT ON COLUMN "public".mm_bank_card.account_name IS '户名';
COMMENT ON COLUMN "public".mm_bank_card.bank_id IS '银行编号';
COMMENT ON COLUMN "public".mm_bank_card.bank_name IS '银行名称';
COMMENT ON COLUMN "public".mm_bank_card.bank_code IS '银行卡代码';
COMMENT ON COLUMN "public".mm_bank_card.network IS '网点';
COMMENT ON COLUMN "public".mm_bank_card.auth_code IS '快捷支付授权码';
COMMENT ON COLUMN "public".mm_bank_card.state IS '状态';
COMMENT ON COLUMN "public".mm_bank_card.create_time IS '添加时间';

-- 删除钱包表 --
DROP TABLE IF EXISTS "public".wal_wallet CASCADE;
DROP TABLE IF EXISTS "public".wal_wallet_log CASCADE;



CREATE TABLE "public".wal_wallet (
     id              bigserial NOT NULL,
     hash_code       varchar(40) NOT NULL,
     node_id         int4 NOT NULL,
     user_id         int8 NOT NULL,
     wallet_type     int4 NOT NULL,
     wallet_flag     int4 NOT NULL,
     balance         int4 DEFAULT 0 NOT NULL,
     present_balance int4 NOT NULL,
     adjust_amount   int4 NOT NULL,
     freeze_amount   int4 NOT NULL,
     latest_amount   int4 NOT NULL,
     expired_amount  int4 NOT NULL,
     total_charge    int4 DEFAULT 0 NOT NULL,
     total_present   int4 NOT NULL,
     total_pay       int4 DEFAULT 0 NOT NULL,
     state           int2 NOT NULL,
     remark          varchar(40) NOT NULL,
     create_time     int8 NOT NULL,
     update_time     int8 DEFAULT 0 NOT NULL,
     CONSTRAINT wal_wallet_pkey
         PRIMARY KEY (id));
COMMENT ON TABLE "public".wal_wallet IS '钱包';
COMMENT ON COLUMN "public".wal_wallet.id IS '编号';
COMMENT ON COLUMN "public".wal_wallet.hash_code IS '哈希值';
COMMENT ON COLUMN "public".wal_wallet.node_id IS '节点编号';
COMMENT ON COLUMN "public".wal_wallet.user_id IS '用户编号';
COMMENT ON COLUMN "public".wal_wallet.wallet_type IS '钱包类型';
COMMENT ON COLUMN "public".wal_wallet.wallet_flag IS '钱包标志';
COMMENT ON COLUMN "public".wal_wallet.balance IS '余额';
COMMENT ON COLUMN "public".wal_wallet.present_balance IS '赠送余额';
COMMENT ON COLUMN "public".wal_wallet.adjust_amount IS '调整禁遏';
COMMENT ON COLUMN "public".wal_wallet.freeze_amount IS '冻结金额';
COMMENT ON COLUMN "public".wal_wallet.latest_amount IS '结余金额';
COMMENT ON COLUMN "public".wal_wallet.expired_amount IS '失效账户余额';
COMMENT ON COLUMN "public".wal_wallet.total_charge IS '总充值金额';
COMMENT ON COLUMN "public".wal_wallet.total_present IS '累计赠送金额';
COMMENT ON COLUMN "public".wal_wallet.total_pay IS '总支付额';
COMMENT ON COLUMN "public".wal_wallet.state IS '状态';
COMMENT ON COLUMN "public".wal_wallet.remark IS '备注';
COMMENT ON COLUMN "public".wal_wallet.create_time IS '创建时间';
COMMENT ON COLUMN "public".wal_wallet.update_time IS '更新时间';
CREATE TABLE "public".wal_wallet_log (
     id            bigserial NOT NULL,
     wallet_id     int8 NOT NULL,
     kind          int4 NOT NULL,
     title         varchar(45) NOT NULL,
     outer_chan    varchar(20) NOT NULL,
     outer_no      varchar(45) NOT NULL,
     value         int4 NOT NULL,
     balance       int4 NOT NULL,
     trade_fee     int4 NOT NULL,
     opr_uid       int4 NOT NULL,
     opr_name      varchar(20) NOT NULL,
     bank_name     varchar(20) NOT NULL,
     bank_account  varchar(20) NOT NULL,
     account_name  varchar(20) NOT NULL,
     review_state  int4 NOT NULL,
     review_remark varchar(120) NOT NULL,
     review_time   int8 NOT NULL,
     remark        varchar(40) NOT NULL,
     create_time   int8 NOT NULL,
     update_time   int8 NOT NULL,
     CONSTRAINT wal_wallet_log_pkey
         PRIMARY KEY (id));
COMMENT ON COLUMN "public".wal_wallet_log.id IS '编号';
COMMENT ON COLUMN "public".wal_wallet_log.wallet_id IS '钱包编号';
COMMENT ON COLUMN "public".wal_wallet_log.kind IS '业务类型';
COMMENT ON COLUMN "public".wal_wallet_log.title IS '标题';
COMMENT ON COLUMN "public".wal_wallet_log.outer_chan IS '外部通道';
COMMENT ON COLUMN "public".wal_wallet_log.outer_no IS '外部订单号';
COMMENT ON COLUMN "public".wal_wallet_log.value IS '变动金额';
COMMENT ON COLUMN "public".wal_wallet_log.balance IS '余额';
COMMENT ON COLUMN "public".wal_wallet_log.trade_fee IS '交易手续费';
COMMENT ON COLUMN "public".wal_wallet_log.opr_uid IS '操作人员用户编号';
COMMENT ON COLUMN "public".wal_wallet_log.opr_name IS '操作人员名称';
COMMENT ON COLUMN "public".wal_wallet_log.bank_name IS '提现银行';
COMMENT ON COLUMN "public".wal_wallet_log.bank_account IS '提现银行账号';
COMMENT ON COLUMN "public".wal_wallet_log.account_name IS '提现银行账户名称';
COMMENT ON COLUMN "public".wal_wallet_log.review_state IS '审核状态';
COMMENT ON COLUMN "public".wal_wallet_log.review_remark IS '审核备注';
COMMENT ON COLUMN "public".wal_wallet_log.review_time IS '审核时间';
COMMENT ON COLUMN "public".wal_wallet_log.remark IS '备注';
COMMENT ON COLUMN "public".wal_wallet_log.create_time IS '创建时间';
COMMENT ON COLUMN "public".wal_wallet_log.update_time IS '更新时间';
CREATE INDEX wal_wallet_hash_code
    ON "public".wal_wallet (hash_code);

-- 关联钱包 --
ALTER TABLE "public".mm_account
    ADD COLUMN wallet_code varchar(32) default '' NOT NULL;
COMMENT ON TABLE "public".mm_account IS '会员账户';
COMMENT ON COLUMN "public".mm_account.member_id IS '会员编号';
COMMENT ON COLUMN "public".mm_account.integral IS '积分';
COMMENT ON COLUMN "public".mm_account.freeze_integral IS '冻结积分';
COMMENT ON COLUMN "public".mm_account.balance IS '余额';
COMMENT ON COLUMN "public".mm_account.freeze_balance IS '冻结余额';
COMMENT ON COLUMN "public".mm_account.expired_balance IS '失效的余额';
COMMENT ON COLUMN "public".mm_account.wallet_balance IS '钱包余额';
COMMENT ON COLUMN "public".mm_account.wallet_code IS '钱包代码';
COMMENT ON COLUMN "public".mm_account.freeze_wallet IS '冻结钱包余额,作废';
COMMENT ON COLUMN "public".mm_account.expired_wallet IS ',作废';
COMMENT ON COLUMN "public".mm_account.total_charge IS '累计充值';
COMMENT ON COLUMN "public".mm_account.total_pay IS '累计支付';
COMMENT ON COLUMN "public".mm_account.total_expense IS '累计消费';


ALTER TABLE "public".wal_wallet
    ADD COLUMN wallet_name varchar(40) NOT NULL;
COMMENT ON TABLE "public".wal_wallet IS '钱包';
COMMENT ON COLUMN "public".wal_wallet.id IS '编号';
COMMENT ON COLUMN "public".wal_wallet.hash_code IS '哈希值';
COMMENT ON COLUMN "public".wal_wallet.node_id IS '节点编号';
COMMENT ON COLUMN "public".wal_wallet.user_id IS '用户编号';
COMMENT ON COLUMN "public".wal_wallet.wallet_type IS '钱包类型';
COMMENT ON COLUMN "public".wal_wallet.wallet_flag IS '钱包标志';
COMMENT ON COLUMN "public".wal_wallet.wallet_name IS '钱包名称';
COMMENT ON COLUMN "public".wal_wallet.balance IS '余额';
COMMENT ON COLUMN "public".wal_wallet.present_balance IS '赠送余额';
COMMENT ON COLUMN "public".wal_wallet.adjust_amount IS '调整禁遏';
COMMENT ON COLUMN "public".wal_wallet.freeze_amount IS '冻结金额';
COMMENT ON COLUMN "public".wal_wallet.latest_amount IS '结余金额';
COMMENT ON COLUMN "public".wal_wallet.expired_amount IS '失效账户余额';
COMMENT ON COLUMN "public".wal_wallet.total_charge IS '总充值金额';
COMMENT ON COLUMN "public".wal_wallet.total_present IS '累计赠送金额';
COMMENT ON COLUMN "public".wal_wallet.total_pay IS '总支付额';
COMMENT ON COLUMN "public".wal_wallet.state IS '状态';
COMMENT ON COLUMN "public".wal_wallet.create_time IS '创建时间';
COMMENT ON COLUMN "public".wal_wallet.update_time IS '更新时间';
CREATE INDEX wal_wallet_hash_code
    ON "public".wal_wallet (hash_code);


ALTER TABLE "public".wal_wallet_log
    ADD COLUMN account_no varchar(20) NOT NULL;
COMMENT ON COLUMN "public".wal_wallet_log.id IS '编号';
COMMENT ON COLUMN "public".wal_wallet_log.wallet_id IS '钱包编号';
COMMENT ON COLUMN "public".wal_wallet_log.kind IS '业务类型';
COMMENT ON COLUMN "public".wal_wallet_log.title IS '标题';
COMMENT ON COLUMN "public".wal_wallet_log.outer_chan IS '外部通道';
COMMENT ON COLUMN "public".wal_wallet_log.outer_no IS '外部订单号';
COMMENT ON COLUMN "public".wal_wallet_log.value IS '变动金额';
COMMENT ON COLUMN "public".wal_wallet_log.balance IS '余额';
COMMENT ON COLUMN "public".wal_wallet_log.trade_fee IS '交易手续费';
COMMENT ON COLUMN "public".wal_wallet_log.opr_uid IS '操作人员用户编号';
COMMENT ON COLUMN "public".wal_wallet_log.opr_name IS '操作人员名称';
COMMENT ON COLUMN "public".wal_wallet_log.account_no IS '提现账号';
COMMENT ON COLUMN "public".wal_wallet_log.account_name IS '提现银行账户名称';
COMMENT ON COLUMN "public".wal_wallet_log.bank_name IS '提现银行';
COMMENT ON COLUMN "public".wal_wallet_log.review_state IS '审核状态';
COMMENT ON COLUMN "public".wal_wallet_log.review_remark IS '审核备注';
COMMENT ON COLUMN "public".wal_wallet_log.review_time IS '审核时间';
COMMENT ON COLUMN "public".wal_wallet_log.remark IS '备注';
COMMENT ON COLUMN "public".wal_wallet_log.create_time IS '创建时间';
COMMENT ON COLUMN "public".wal_wallet_log.update_time IS '更新时间';

