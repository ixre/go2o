-- 更换为POSTGRESQL --
ALTER TABLE "public".mm_member ADD COLUMN flag int4 DEFAULT 0 NOT NULL;