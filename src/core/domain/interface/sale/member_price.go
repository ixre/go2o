/**
 * Copyright 2015 @ S1N1 Team.
 * name : member_price
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

/*
CREATE TABLE `zsdb`.`gs_member_price` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `goods_id` INT NULL,
  `level` INT NULL,
  `price` INT NULL,
  `enabled` TINYINT(1) NULL,
  PRIMARY KEY (`id`));
*/

// 会员价
type MemberPrice struct{
    Id int `db:"id" pk:"yes" auto:"yes"`
    GoodsId int `db:"goods_id"`
    Level int   `db:"level"`
    Price float32 `db:"price"`
    Enabled int `db:"enabled"`
}