ALTER TABLE `zxdb`.`pt_merchant`
RENAME TO  `zxdb`.`pt_merchant` ;
ALTER TABLE `zxdb`.`pt_page`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`dlv_partner_bind`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL , RENAME TO  `zxdb`.`dlv_merchant_bind` ;
ALTER TABLE `zxdb`.`gs_category`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL COMMENT '商家ID(pattern ID);如果为空，则表示模式分类' ;
ALTER TABLE `zxdb`.`gs_sale_tag`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`mm_relation`
CHANGE COLUMN `reg_partner_id` `reg_merchant_id` INT(11) NULL DEFAULT NULL COMMENT '注册商户编号' ;
ALTER TABLE `zxdb`.`pm_info`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_ad`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_api`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NOT NULL ;
ALTER TABLE `zxdb`.`pt_kvset`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_kvset_member`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_mail_queue`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_mail_template`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_member_level`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_order`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL COMMENT '商家ID' ;
ALTER TABLE `zxdb`.`pt_saleconf`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NOT NULL ;
ALTER TABLE `zxdb`.`pt_shop`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NULL DEFAULT NULL ;
ALTER TABLE `zxdb`.`pt_siteconf`
CHANGE COLUMN `partner_id` `merchant_id` INT(11) NOT NULL ;


