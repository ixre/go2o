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

CREATE TABLE `comm_qr_template` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `title` varchar(45) COLLATE utf8_unicode_ci NOT NULL COMMENT '模板标题',
  `bg_image` varchar(120) COLLATE utf8_unicode_ci NOT NULL COMMENT '背景图片',
  `offset_x` int(11) NOT NULL COMMENT '垂直偏离量',
  `offset_y` int(11) NOT NULL COMMENT '垂直偏移量',
  `comment` varchar(120) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '二维码模板文本',
  `enabled` int(1) DEFAULT NULL COMMENT '是否启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='二维码模板';

ALTER TABLE `go2o`.`comm_qr_template`
ADD COLUMN `callback_url` VARCHAR(120) NULL COMMENT '回调地址' AFTER `comment`;

