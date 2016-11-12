CREATE TABLE `mm_levelup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) NOT NULL COMMENT '会员编号',
  `origin_level` tinyint(2) NOT NULL COMMENT '原来等级',
  `target_level` tinyint(2) NOT NULL COMMENT '现在等级',
  `is_free` int(1) NOT NULL COMMENT '是否为免费升级的会员',
  `payment_id` int(11) NOT NULL COMMENT '支付单编号',
  `reviewed` int(1) NOT NULL COMMENT '是否审核及处理',
  `create_time` int(11) NOT NULL COMMENT '升级时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='升级日志表';