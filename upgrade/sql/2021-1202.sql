ALTER TABLE "public".mch_merchant
    ADD COLUMN salt varchar(10) DEFAULT '' NOT NULL;
COMMENT
ON COLUMN "public".mch_merchant.salt IS '加密盐';