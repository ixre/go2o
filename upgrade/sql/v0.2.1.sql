update pro_product set sort_num=0 WHERE sort_num IS NULL;
 ALTER TABLE pro_product modify column sort_num int(11) NOT NULL;

ALTER TABLE `mm_levelup`
CHANGE COLUMN `upgrade_type` `upgrade_mode` INT(10) NOT NULL AFTER `payment_id`,
CHANGE COLUMN `review_state` `review_state` TINYINT(1) NOT NULL ;

ALTER TABLE mm_trusted_info
   ADD COLUMN country_code varchar(10) NOT NULL comment '证件区域' after real_name,
   ADD COLUMN card_image varchar(120) NOT NULL comment '证件图像' after card_id,
    ADD COLUMN card_type int(1) NOT NULL comment '证件类型' after country_code,
  ADD COLUMN manual_review int(1) NOT NULL comment '是否人工认证' after trust_image;


ALTER TABLE `mch_offline_shop`
CHANGE COLUMN `tel` `tel` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `addr` `addr` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `lng` `lng` FLOAT(5,2) NOT NULL ,
CHANGE COLUMN `lat` `lat` FLOAT(5,2) NOT NULL ,
CHANGE COLUMN `deliver_radius` `deliver_radius` INT(11) NOT NULL COMMENT '配送范围' ,
CHANGE COLUMN `province` `province` INT(11) NOT NULL ,
CHANGE COLUMN `city` `city` INT(11) NOT NULL ,
CHANGE COLUMN `district` `district` INT(11) NOT NULL ;

ALTER TABLE `mch_online_shop`
CHANGE COLUMN `alias` `alias` VARCHAR(20) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `tel` `tel` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `addr` `addr` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `host` `host` VARCHAR(40) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `logo` `logo` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ;


ALTER TABLE `mch_shop`
CHANGE COLUMN `vendor_id` `vendor_id` INT(11) NOT NULL COMMENT '商户编号' ,
CHANGE COLUMN `name` `name` VARCHAR(50) NOT NULL COMMENT '商店名称' ,
CHANGE COLUMN `sort_number` `sort_num` INT(11) NOT NULL COMMENT '排序序号' ,
CHANGE COLUMN `create_time` `create_time` INT(11) NOT NULL ;


ALTER TABLE `article_list`
CHANGE COLUMN `cat_id` `cat_id` INT(11) NOT NULL ,
CHANGE COLUMN `title` `title` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `small_title` `small_title` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `thumbnail` `thumbnail` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `publisher_id` `publisher_id` INT(11) NOT NULL ,
CHANGE COLUMN `location` `location` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `priority` `priority` INT(2) NOT NULL COMMENT '优先级' ,
CHANGE COLUMN `access_key` `access_key` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL COMMENT '访问钥匙' ,
CHANGE COLUMN `content` `content` TEXT CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `tags` `tags` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `view_count` `view_count` INT(11) NOT NULL ,
CHANGE COLUMN `sort_number` `sort_num` INT(11) NOT NULL ,
CHANGE COLUMN `create_time` `create_time` INT(11) NOT NULL ;

ALTER TABLE `article_category`
CHANGE COLUMN `parent_id` `parent_id` INT(11) NOT NULL ,
CHANGE COLUMN `perm_flag` `perm_flag` INT(2) NOT NULL COMMENT '访问权限' ,
CHANGE COLUMN `name` `name` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `cat_alias` `cat_alias` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `title` `title` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `keywords` `keywords` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `describe` `describe` VARCHAR(250) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `sort_number` `sort_num` INT(11) NOT NULL ,
CHANGE COLUMN `location` `location` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `update_time` `update_time` INT(11) NOT NULL ;

ALTER TABLE `ex_page`
CHANGE COLUMN `enabled` `enabled` INT(1) NOT NULL AFTER `css_path`,
CHANGE COLUMN `user_id` `user_id` INT(11) NOT NULL ,
CHANGE COLUMN `title` `title` VARCHAR(100) NOT NULL COMMENT '标题' ,
CHANGE COLUMN `perm_flag` `perm_flag` INT(2) NOT NULL COMMENT '访问权限' ,
CHANGE COLUMN `access_key` `access_key` VARCHAR(45) NOT NULL COMMENT '访问钥匙' ,
CHANGE COLUMN `str_indent` `str_indent` VARCHAR(50) NOT NULL ,
CHANGE COLUMN `keyword` `keyword` VARCHAR(100) NOT NULL ,
CHANGE COLUMN `description` `description` VARCHAR(150) NOT NULL ,
CHANGE COLUMN `css_path` `css_path` VARCHAR(100) NOT NULL ,
CHANGE COLUMN `body` `body` TEXT NOT NULL ,
CHANGE COLUMN `update_time` `update_time` INT(11) NOT NULL ;

ALTER TABLE `ad_group`
CHANGE COLUMN `name` `name` VARCHAR(10) NOT NULL ,
CHANGE COLUMN `opened` `opened` INT(1) NOT NULL ,
CHANGE COLUMN `enabled` `enabled` INT(1) NOT NULL ;

ALTER TABLE `ad_hyperlink`
CHANGE COLUMN `ad_id` `ad_id` INT(11) NOT NULL ,
CHANGE COLUMN `title` `title` VARCHAR(50) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `link_url` `link_url` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ;

ALTER TABLE `ad_image`
CHANGE COLUMN `ad_id` `ad_id` INT(11) NOT NULL ,
CHANGE COLUMN `title` `title` VARCHAR(45) NOT NULL ,
CHANGE COLUMN `link_url` `link_url` VARCHAR(100) NOT NULL ,
CHANGE COLUMN `image_url` `image_url` VARCHAR(150) NOT NULL ,
CHANGE COLUMN `sort_number` `sort_num` INT(11) NOT NULL ,
CHANGE COLUMN `enabled` `enabled` INT(1) NOT NULL ;

ALTER TABLE `ad_image_ad`
CHANGE COLUMN `ad_id` `ad_id` INT(10) NOT NULL ,
CHANGE COLUMN `title` `title` VARCHAR(45) NOT NULL ,
CHANGE COLUMN `link_url` `link_url` VARCHAR(100) NOT NULL ,
CHANGE COLUMN `image_url` `image_url` VARCHAR(150) NOT NULL ,
CHANGE COLUMN `sort_number` `sort_num` INT(10) NOT NULL ,
CHANGE COLUMN `enabled` `enabled` INT(1) NOT NULL ;


ALTER TABLE `ad_list`
CHANGE COLUMN `user_id` `user_id` INT(11) NOT NULL ,
CHANGE COLUMN `name` `name` VARCHAR(45) NOT NULL ,
CHANGE COLUMN `type_id` `type_id` TINYINT(1) NOT NULL ,
CHANGE COLUMN `show_times` `show_times` INT(11) NOT NULL COMMENT '展现数量' ,
CHANGE COLUMN `click_times` `click_times` INT(11) NOT NULL COMMENT '点击次数' ,
CHANGE COLUMN `show_days` `show_days` INT(11) NOT NULL COMMENT '投放天数' ,
CHANGE COLUMN `update_time` `update_time` INT(11) NOT NULL ;

ALTER TABLE `ad_position`
DROP FOREIGN KEY `id`;
ALTER TABLE `ad_position`
CHANGE COLUMN `group_id` `group_id` INT(11) NOT NULL ,
CHANGE COLUMN `key` `key` VARCHAR(45) NOT NULL ,
CHANGE COLUMN `name` `name` VARCHAR(45) NOT NULL ,
CHANGE COLUMN `default_id` `default_id` INT(11) NOT NULL ,
CHANGE COLUMN `opened` `opened` INT(1) NOT NULL ,
CHANGE COLUMN `enabled` `enabled` INT(1) NOT NULL ;

ALTER TABLE `ad_position` 
DROP INDEX `id_idx` ;

ALTER TABLE `ad_userset`
CHANGE COLUMN `pos_id` `pos_id` INT(11) NOT NULL ,
CHANGE COLUMN `user_id` `user_id` INT(11) NOT NULL ,
CHANGE COLUMN `ad_id` `ad_id` INT(11) NOT NULL ;


ALTER TABLE `mm_trusted_info`
CHANGE COLUMN `real_name` `real_name` VARCHAR(10) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `card_id` `card_id` VARCHAR(20) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `trust_image` `trust_image` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `review_state` `review_state` TINYINT(1) NOT NULL ,
CHANGE COLUMN `review_time` `review_time` INT(11) NOT NULL ,
CHANGE COLUMN `remark` `remark` VARCHAR(120) CHARACTER SET 'utf8' COLLATE 'utf8_unicode_ci' NOT NULL ,
CHANGE COLUMN `update_time` `update_time` INT(11) NOT NULL ;






