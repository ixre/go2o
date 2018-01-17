

/*ALTER TABLE mm_trusted_info
  DROP COLUMN reviewed;*/
ALTER TABLE mm_trusted_info
  modify column real_name varchar(10) NOT NULL;
ALTER TABLE mm_trusted_info
  ADD COLUMN card_type int(1) NOT NULL comment '证件类型'
ALTER TABLE mm_trusted_info
  ADD COLUMN card_area varchar(5) NOT NULL comment '证件区域';
ALTER TABLE mm_trusted_info
  modify column card_id varchar(20) NOT NULL;
ALTER TABLE mm_trusted_info
  ADD COLUMN card_image varchar(120) NOT NULL comment '证件图像';
ALTER TABLE mm_trusted_info
  modify column trust_image varchar(120) NOT NULL;
ALTER TABLE mm_trusted_info
  ADD COLUMN review_state int(1) NOT NULL comment '审核状态';
ALTER TABLE mm_trusted_info
  modify column review_time int(11) NOT NULL;
ALTER TABLE mm_trusted_info
  ADD COLUMN manual_review int(1) NOT NULL comment '是否人工认证';
ALTER TABLE mm_trusted_info
  modify column remark varchar(120) NOT NULL;
ALTER TABLE mm_trusted_info
  modify column update_time int(11) NOT NULL;
