-- MySQL dump 10.13  Distrib 5.6.19, for Linux (x86_64)
--
-- Host: localhost    Database: foodording
-- ------------------------------------------------------
-- Server version	5.5.40-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `conf_integral_type`
--

DROP TABLE IF EXISTS `conf_integral_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `conf_integral_type` (
  `id` int(11) NOT NULL,
  `name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `conf_integral_type`
--

LOCK TABLES `conf_integral_type` WRITE;
/*!40000 ALTER TABLE `conf_integral_type` DISABLE KEYS */;
INSERT INTO `conf_integral_type` VALUES (1,'系统赠送'),(2,'登录'),(3,'订单'),(4,'退还'),(11,'系统扣减'),(12,'兑换');
/*!40000 ALTER TABLE `conf_integral_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `conf_member_level`
--

DROP TABLE IF EXISTS `conf_member_level`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `conf_member_level` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `value` int(11) DEFAULT '1' COMMENT '等级值',
  `require_exp` int(11) DEFAULT NULL COMMENT '要求积分',
  `enabled` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `conf_member_level`
--

LOCK TABLES `conf_member_level` WRITE;
/*!40000 ALTER TABLE `conf_member_level` DISABLE KEYS */;
INSERT INTO `conf_member_level` VALUES (1,'一星会员',1,0,1),(2,'二星会员',2,100,1),(3,'三星会员',3,300,1),(4,'四星会员',4,700,1),(5,'五星会员',5,1200,1);
/*!40000 ALTER TABLE `conf_member_level` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gs_category`
--

DROP TABLE IF EXISTS `gs_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gs_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL COMMENT '父分类',
  `partner_id` int(11) DEFAULT NULL COMMENT '商家ID(pattern ID);如果为空，则表示模式分类',
  `name` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `enabled` bit(1) DEFAULT NULL COMMENT '是否可用',
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '描述',
  `order_index` int(11) DEFAULT '0' COMMENT '序号',
  `create_time` int(11) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=20 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='food categories';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gs_category`
--

LOCK TABLES `gs_category` WRITE;
/*!40000 ALTER TABLE `gs_category` DISABLE KEYS */;
INSERT INTO `gs_category` VALUES (13,0,666888,'小炒','','',0,2012),(14,0,666888,'面食','',NULL,0,2012),(15,0,666888,'套餐','','',0,2012),(16,0,666888,'油炸','',NULL,0,2012),(17,0,666888,'海鲜','',NULL,0,2012),(18,15,666888,'营养套餐','','',5,2012);
/*!40000 ALTER TABLE `gs_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gs_goods`
--

DROP TABLE IF EXISTS `gs_goods`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gs_goods` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) DEFAULT NULL COMMENT '分类',
  `name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `small_title` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `on_shelves` tinyint(4) DEFAULT NULL COMMENT '是否上架',
  `img` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `cost` decimal(5,2) DEFAULT '0.00' COMMENT ' 成本价',
  `price` decimal(5,2) DEFAULT '0.00' COMMENT '售价(市场价)',
  `sale_price` decimal(5,2) DEFAULT NULL COMMENT '实际销售价',
  `apply_subs` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '供应分店,用'',''隔开',
  `note` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '备注，如新菜色，特价优惠等',
  `description` varchar(500) COLLATE utf8_unicode_ci DEFAULT NULL,
  `state` int(11) DEFAULT '1',
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=31 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='食物项';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gs_goods`
--

LOCK TABLES `gs_goods` WRITE;
/*!40000 ALTER TABLE `gs_goods` DISABLE KEYS */;
INSERT INTO `gs_goods` VALUES (1,18,'韭黄炒蛋饭-2-2',NULL,1,'666888/item_pic/20141022090923.png',15.00,20.00,18.00,'1','1','2',1,1421419030,1421419030),(2,15,'鱼香茄子饭','1',1,NULL,15.00,20.00,18.00,'1,7,10,11',NULL,NULL,1,1421419031,1421419031),(3,15,'尖椒回锅肉',NULL,1,NULL,15.00,20.00,18.00,NULL,NULL,NULL,1,1421419032,1421419032),(4,15,'香菇焖鸡饭',NULL,1,NULL,15.00,20.00,18.00,NULL,NULL,NULL,1,1421419033,1421419033),(5,15,'野山椒炒肉',NULL,1,NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,1421419034,1421419034),(6,15,'极品红烧肉',NULL,1,NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,1421419035,1421419035),(7,15,'花生猪手饭',NULL,1,NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,1421419036,1421419036),(8,15,'红烧鱼腩饭',NULL,1,NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,1421419037,1421419037),(9,15,'香卤鸡腿饭',NULL,1,NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,1421419038,1421419038),(10,15,'酸甜排骨饭',NULL,1,NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,1421419039,1421419039),(29,18,'营养套餐B',NULL,1,'666888/item_pic/20141023090944.png',10.00,20.00,18.00,'1',NULL,NULL,1,1421419058,1421419058),(28,18,'营养套餐A','鸡蛋＋香肠双拼',1,'666888/item_pic/20141022090951.png',5.00,15.00,12.00,'1,7,10,11',NULL,'d',1,1421419057,1423030454);
/*!40000 ALTER TABLE `gs_goods` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gs_snapshot`
--

DROP TABLE IF EXISTS `gs_snapshot`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `gs_snapshot` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `snapshot_key` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `goods_id` int(11) DEFAULT NULL,
  `goods_name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `small_title` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `category_name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `img` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `cost` decimal(5,2) DEFAULT '0.00' COMMENT ' 成本价',
  `price` decimal(5,2) DEFAULT '0.00' COMMENT '售价(市场价)',
  `sale_price` decimal(5,2) DEFAULT NULL COMMENT '实际销售价',
  `create_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=35 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='食物项';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gs_snapshot`
--

LOCK TABLES `gs_snapshot` WRITE;
/*!40000 ALTER TABLE `gs_snapshot` DISABLE KEYS */;
INSERT INTO `gs_snapshot` VALUES (31,'666888-g4-1423022472',4,'香菇焖鸡饭',NULL,'套餐',NULL,15.00,20.00,18.00,1423022472),(32,'666888-g28-1423028841',28,'营养套餐A','鸡蛋＋香肠双拼','营养套餐','666888/item_pic/20141022090951.png',5.00,15.00,12.00,1423028841),(33,'666888-g28-1423029098',28,'营养套餐A','-鸡蛋＋香肠双拼','营养套餐','666888/item_pic/20141022090951.png',5.00,15.00,12.00,1423029098),(34,'666888-g28-1423029140',28,'营养套餐A','鸡蛋＋香肠双拼','营养套餐','666888/item_pic/20141022090951.png',5.00,15.00,12.00,1423029140);
/*!40000 ALTER TABLE `gs_snapshot` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mm_account`
--

DROP TABLE IF EXISTS `mm_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mm_account` (
  `member_id` int(11) NOT NULL,
  `integral` int(11) DEFAULT '0',
  `balance` float(8,2) DEFAULT NULL,
  `present_balance` float(8,2) DEFAULT NULL,
  `total_fee` float(8,2) DEFAULT NULL,
  `total_charge` float(8,2) DEFAULT NULL,
  `total_pay` float(8,2) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL COMMENT '积分',
  PRIMARY KEY (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_account`
--

LOCK TABLES `mm_account` WRITE;
/*!40000 ALTER TABLE `mm_account` DISABLE KEYS */;
INSERT INTO `mm_account` VALUES (1,34930,11.42,351.44,4404.42,0.00,4303.60,1422947691),(2,0,2.50,NULL,2.50,0.00,0.00,2013),(28,360,2.79,2.88,38.79,0.00,54.00,1421419030),(29,0,4.14,NULL,4.14,0.00,36.00,2012),(30,0,2.70,NULL,2.70,0.00,27.00,2012),(31,0,0.00,NULL,0.00,0.00,0.00,2012),(32,0,0.00,NULL,0.00,0.00,0.00,2012),(33,0,0.00,NULL,0.00,0.00,0.00,2012),(34,0,0.00,NULL,0.00,0.00,0.00,2013),(35,0,0.00,NULL,0.00,0.00,0.00,2013),(36,0,0.00,NULL,0.00,0.00,0.00,2013),(37,0,0.00,NULL,0.00,0.00,0.00,2013),(38,0,1.00,NULL,1.00,0.00,10.00,2013),(39,0,2.00,NULL,2.00,0.00,20.00,2013),(40,0,3.80,NULL,3.80,0.00,38.00,2013),(41,0,0.00,NULL,0.00,0.00,0.00,2013);
/*!40000 ALTER TABLE `mm_account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mm_bank`
--

DROP TABLE IF EXISTS `mm_bank`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mm_bank` (
  `member_id` int(11) NOT NULL,
  `name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '银行名称',
  `account` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '银行账号',
  `account_name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '银行户名',
  `network` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '银行网点',
  `state` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_bank`
--

LOCK TABLES `mm_bank` WRITE;
/*!40000 ALTER TABLE `mm_bank` DISABLE KEYS */;
INSERT INTO `mm_bank` VALUES (1,'中国工商银行','513701198801105317','张三 ','上海分行漕溪路支行 ',0,2012),(28,'中国邮政储蓄银行','123486855651',NULL,NULL,1,2012);
/*!40000 ALTER TABLE `mm_bank` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mm_deliver_addr`
--

DROP TABLE IF EXISTS `mm_deliver_addr`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mm_deliver_addr` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) DEFAULT NULL,
  `real_name` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `phone` varchar(11) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_default` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_deliver_addr`
--

LOCK TABLES `mm_deliver_addr` WRITE;
/*!40000 ALTER TABLE `mm_deliver_addr` DISABLE KEYS */;
INSERT INTO `mm_deliver_addr` VALUES (1,1,'刘铭','18616999822','上海市',0),(2,27,'刘铭','18616999822','上海市',NULL),(4,28,'刘大爷','18616999822','佛山市',0),(5,1,'谢环子','18867889090','上海市徐汇区三江路23号西雅图小区1902室',0);
/*!40000 ALTER TABLE `mm_deliver_addr` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mm_income_log`
--

DROP TABLE IF EXISTS `mm_income_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mm_income_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` int(11) DEFAULT NULL,
  `member_id` int(11) DEFAULT NULL,
  `type` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `fee` float(6,2) DEFAULT NULL,
  `log` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `record_time` int(11) DEFAULT NULL,
  `state` int(11) DEFAULT NULL COMMENT '状态(如：无效），默认为1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=143 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='进账日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_income_log`
--

LOCK TABLES `mm_income_log` WRITE;
/*!40000 ALTER TABLE `mm_income_log` DISABLE KEYS */;
INSERT INTO `mm_income_log` VALUES (8,NULL,28,'backcash',0.27,'来自订单:20120924-9603(商家:万绿园,会员:sumian)收入￥0.27元.',2012,NULL),(9,NULL,1,'backcash',7.40,'订单:20130308-2932返现￥7.4元',2013,NULL),(10,NULL,38,'backcash',1.00,'订单:20130307-5447返现￥1.0元',2013,NULL),(11,NULL,2,'backcash',0.20,'来自订单:20130307-5447(商家:蚁族,会员:yuan)收入￥0.2元.',2013,NULL),(12,NULL,39,'backcash',2.00,'订单:20130311-9708返现￥2.0元',2013,NULL),(13,NULL,2,'backcash',0.40,'来自订单:20130311-9708(商家:蚁族,会员:yuanyuan)收入￥0.4元.',2013,NULL),(14,NULL,40,'backcash',0.90,'订单:20130331-4922返现￥0.9元',2013,NULL),(15,NULL,2,'backcash',0.18,'来自订单:20130331-4922(商家:蚁族,会员:13924886758)收入￥0.18元.',2013,NULL),(16,NULL,40,'backcash',2.90,'订单:20130401-2503返现￥2.9元',2013,NULL),(17,NULL,2,'backcash',0.58,'来自订单:20130401-2503(商家:蚁族,会员:13924886758)收入￥0.58元.',2013,NULL),(18,NULL,1,'backcash',4.80,'订单:20130401-2503(商家:蚁族)返现￥4.80元',2014,1),(19,NULL,1,'backcash',1.20,'订单:20130401-2503(商家:蚁族,会员:刘铭)收入￥%!s(float32=1.2)元',2014,1),(20,NULL,1,'backcash',0.60,'订单:20130401-2503(商家:蚁族,会员:刘铭)收入￥%!s(float32=0.6)元',2014,1),(21,NULL,1,'backcash',4.80,'订单:20130401-7643(商家:蚁族)返现￥4.80元',2014,1),(22,NULL,1,'backcash',1.20,'订单:20130401-7643(商家:蚁族,会员:刘铭)收入￥1.20元',2014,1),(23,NULL,1,'backcash',0.60,'订单:20130401-7643(商家:蚁族,会员:刘铭)收入￥0.60元',2014,1),(24,NULL,1,'backcash',4.80,'订单:20130401-5283(商家:蚁族)返现￥4.80元',2014,1),(25,NULL,1,'backcash',1.20,'订单:20130401-5283(商家:蚁族,会员:刘铭)收入￥1.20元',2014,1),(26,NULL,1,'backcash',0.60,'订单:20130401-5283(商家:蚁族,会员:刘铭)收入￥0.60元',2014,1),(27,NULL,1,'backcash',4.80,'订单:20130401-8432(商家:蚁族)返现￥4.80元',2014,1),(28,NULL,1,'backcash',1.20,'订单:20130401-8432(商家:蚁族,会员:刘铭)收入￥1.20元',2014,1),(29,NULL,1,'backcash',0.60,'订单:20130401-8432(商家:蚁族,会员:刘铭)收入￥0.60元',2014,1),(30,NULL,1,'backcash',4.80,'订单:20130331-4998(商家:蚁族)返现￥4.80元',2014,1),(31,NULL,1,'backcash',1.20,'订单:20130331-4998(商家:蚁族,会员:刘铭)收入￥1.20元',2014,1),(32,NULL,1,'backcash',0.60,'订单:20130331-4998(商家:蚁族,会员:刘铭)收入￥0.60元',2014,1),(33,NULL,1,'backcash',4.80,'订单:20130331-2185(商家:蚁族)返现￥4.80元',2014,1),(34,NULL,1,'backcash',1.20,'订单:20130331-2185(商家:蚁族,会员:刘铭)收入￥1.20元',2014,1),(35,NULL,1,'backcash',0.60,'订单:20130331-2185(商家:蚁族,会员:刘铭)收入￥0.60元',2014,1),(36,NULL,1,'backcash',4.80,'订单:20130401-4602(商家:蚁族)返现￥4.80元',2014,1),(37,NULL,1,'backcash',1.20,'订单:20130401-4602(商家:蚁族,会员:刘铭)收入￥1.20元',2014,1),(38,NULL,1,'backcash',0.60,'订单:20130401-4602(商家:蚁族,会员:刘铭)收入￥0.60元',2014,1),(39,NULL,1,'backcash',4.80,'订单:20130331-4922(商家:蚁族)返现￥4.80元',2014,1),(40,NULL,1,'backcash',1.20,'订单:20130331-4922(商家:蚁族,会员:刘铭)收入￥1.20元',2014,1),(41,NULL,1,'backcash',0.60,'订单:20130331-4922(商家:蚁族,会员:刘铭)收入￥0.60元',2014,1),(42,NULL,1,'backcash',4.80,'订单:20130315-1432(商家:蚁族)返现￥4.80元',2014,1),(43,NULL,1,'backcash',1.20,'订单:20130315-1432(商家:蚁族,会员:刘铭)收入￥1.20元',2014,1),(44,NULL,1,'backcash',0.60,'订单:20130315-1432(商家:蚁族,会员:刘铭)收入￥0.60元',2014,1),(45,NULL,1,'backcash',4.80,'订单:20130311-9708(商家:蚁族)返现￥4.80元',2014,1),(46,NULL,1,'backcash',1.20,'订单:20130311-9708(商家:蚁族,会员:刘铭)收入￥1.20元',2014,1),(47,NULL,1,'backcash',0.60,'订单:20130311-9708(商家:蚁族,会员:刘铭)收入￥0.60元',2014,1),(48,NULL,1,'backcash',4.80,'订单:20130308-2932(商家:蚁族)返现￥4.80元',2014,1),(49,NULL,1,'backcash',1.20,'订单:20130308-2932(商家:蚁族,会员:刘铭)收入￥1.20元',2014,1),(54,129,1,'backcash',4.80,'订单:20130203-1190(商家:美味道)返现￥4.80元',2014,1),(55,127,1,'backcash',4.80,'订单:20130115-7948(商家:美味道)返现￥4.80元',2014,1),(58,153,1,'backcash',4.40,'订单:689182668(商家:美味道)返现￥4.40元',2014,1),(59,0,1,'backcash',1.10,'订单:689182668(商家:美味道,会员:刘铭)收入￥1.10元',2014,1),(60,0,1,'backcash',0.55,'订单:689182668(商家:美味道,会员:刘铭)收入￥0.55元',2014,1),(61,154,1,'backcash',1.44,'订单:685815333(商家:美味道)返现￥1.44元',2014,1),(62,0,1,'backcash',0.36,'订单:685815333(商家:美味道,会员:刘铭)收入￥0.36元',2014,1),(63,0,1,'backcash',0.18,'订单:685815333(商家:美味道,会员:刘铭)收入￥0.18元',2014,1),(64,152,1,'backcash',2.88,'订单:684452597(商家:美味道)返现￥2.88元',2014,1),(65,0,1,'backcash',0.72,'订单:684452597(商家:美味道,会员:刘铭)收入￥0.72元',2014,1),(66,0,1,'backcash',0.36,'订单:684452597(商家:美味道,会员:刘铭)收入￥0.36元',2014,1),(67,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元',2014,1),(68,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元',2014,1),(69,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元',2014,1),(70,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元',2014,1),(71,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元',2014,1),(72,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元',2014,1),(73,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元',2014,1),(74,180,1,'backcash',4.32,'订单:687153988(商家:美味道)返现￥4.32元',2014,1),(75,179,1,'backcash',2.88,'订单:681645404(商家:美味道)返现￥2.88元',2014,1),(76,177,1,'backcash',2.88,'订单:687025732(商家:美味道)返现￥2.88元',2014,1),(77,183,28,'backcash',1.44,'订单:689469831(商家:美味道)返现￥1.44元',2014,1),(78,184,1,'backcash',1.44,'订单:689205996(商家:美味道)返现￥1.44元',2014,1),(79,182,1,'backcash',2.88,'订单:684148199(商家:美味道)返现￥2.88元',2014,1),(80,178,1,'backcash',2.88,'订单:683878066(商家:美味道)返现￥2.88元',2014,1),(81,176,1,'backcash',2.88,'订单:681710378(商家:美味道)返现￥2.88元',2014,1),(82,175,1,'backcash',2.88,'订单:681009722(商家:美味道)返现￥2.88元',2014,1),(83,174,1,'backcash',2.88,'订单:682466077(商家:美味汇)返现￥2.88元',NULL,1),(84,173,1,'backcash',2.88,'订单:684691894(商家:美味汇)返现￥2.88元',NULL,1),(85,171,1,'backcash',2.88,'订单:689760500(商家:美味汇)返现￥2.88元',NULL,1),(86,172,1,'backcash',2.88,'订单:687956073(商家:美味汇)返现￥2.88元',NULL,1),(87,170,1,'backcash',1.44,'订单:681554163(商家:美味汇)返现￥1.44元',NULL,1),(88,167,1,'backcash',2.88,'订单:689386565(商家:美味汇)返现￥2.88元',NULL,1),(89,164,1,'backcash',1.44,'订单:683135658(商家:美味汇)返现￥1.44元',NULL,1),(90,190,1,'backcash',2.88,'订单:685994267(商家:美味汇)返现￥2.88元',2015,1),(91,190,1,'backcash',2.88,'订单:685994267(商家:美味汇)返现￥2.88元',1421414462,1),(92,151,1,'backcash',2.88,'订单:686267421(商家:美味汇)返现￥2.88元',1421419029,1),(93,152,1,'backcash',2.88,'订单:684452597(商家:美味汇)返现￥2.88元',1421419029,1),(94,153,1,'backcash',4.40,'订单:689182668(商家:美味汇)返现￥4.40元',1421419029,1),(95,154,1,'backcash',1.44,'订单:685815333(商家:美味汇)返现￥1.44元',1421419029,1),(96,164,1,'backcash',1.44,'订单:683135658(商家:美味汇)返现￥1.44元',1421419029,1),(97,165,1,'backcash',1.44,'订单:689424176(商家:美味汇)返现￥1.44元',1421419029,1),(98,166,1,'backcash',2.88,'订单:682964508(商家:美味汇)返现￥2.88元',1421419029,1),(99,167,1,'backcash',2.88,'订单:689386565(商家:美味汇)返现￥2.88元',1421419030,1),(100,170,1,'backcash',1.44,'订单:681554163(商家:美味汇)返现￥1.44元',1421419030,1),(101,171,1,'backcash',2.88,'订单:689760500(商家:美味汇)返现￥2.88元',1421419030,1),(102,172,1,'backcash',2.88,'订单:687956073(商家:美味汇)返现￥2.88元',1421419030,1),(103,173,1,'backcash',2.88,'订单:684691894(商家:美味汇)返现￥2.88元',1421419030,1),(104,174,1,'backcash',2.88,'订单:682466077(商家:美味汇)返现￥2.88元',1421419030,1),(105,175,1,'backcash',2.88,'订单:681009722(商家:美味汇)返现￥2.88元',1421419030,1),(106,176,1,'backcash',2.88,'订单:681710378(商家:美味汇)返现￥2.88元',1421419030,1),(107,177,1,'backcash',2.88,'订单:687025732(商家:美味汇)返现￥2.88元',1421419030,1),(108,178,1,'backcash',2.88,'订单:683878066(商家:美味汇)返现￥2.88元',1421419030,1),(109,179,1,'backcash',2.88,'订单:681645404(商家:美味汇)返现￥2.88元',1421419030,1),(110,180,1,'backcash',4.32,'订单:687153988(商家:美味汇)返现￥4.32元',1421419030,1),(111,181,1,'backcash',10.08,'订单:688664284(商家:美味汇)返现￥10.08元',1421419030,1),(112,182,1,'backcash',2.88,'订单:684148199(商家:美味汇)返现￥2.88元',1421419030,1),(113,183,28,'backcash',1.44,'订单:689469831(商家:美味汇)返现￥1.44元',1421419030,1),(114,184,1,'backcash',1.44,'订单:689205996(商家:美味汇)返现￥1.44元',1421419031,1),(115,191,1,'backcash',1.44,'订单:683249100(商家:美味汇)返现￥1.44元',1421419079,1),(116,155,1,'backcash',1.44,'订单:689681862(商家:美味汇)返现￥1.44元',1421419138,1),(117,156,1,'backcash',3.12,'订单:686862923(商家:美味汇)返现￥3.12元',1421419138,1),(118,157,1,'backcash',2.88,'订单:687643214(商家:美味汇)返现￥2.88元',1421419138,1),(119,158,1,'backcash',2.88,'订单:687137477(商家:美味汇)返现￥2.88元',1421419138,1),(120,159,1,'backcash',4.32,'订单:686200508(商家:美味汇)返现￥4.32元',1421419138,1),(121,160,1,'backcash',1.44,'订单:682658763(商家:美味汇)返现￥1.44元',1421419138,1),(122,161,1,'backcash',1.44,'订单:688563456(商家:美味汇)返现￥1.44元',1421419139,1),(123,162,1,'backcash',1.44,'订单:686432330(商家:美味汇)返现￥1.44元',1421419139,1),(124,163,1,'backcash',1.44,'订单:682539015(商家:美味汇)返现￥1.44元',1421419139,1),(125,185,1,'backcash',4.32,'订单:683958933(商家:美味汇)返现￥4.32元',1421419139,1),(126,186,1,'backcash',4.32,'订单:686018039(商家:美味汇)返现￥4.32元',1421419139,1),(127,187,1,'backcash',2.88,'订单:681426294(商家:美味汇)返现￥2.88元',1421419139,1),(128,188,1,'backcash',4.40,'订单:686396270(商家:美味汇)返现￥4.40元',1421419139,1),(129,189,1,'backcash',4.32,'订单:686660811(商家:美味汇)返现￥4.32元',1421419139,1),(130,192,1,'backcash',2.88,'订单:681879385(商家:美味汇)返现￥2.88元',1421419139,1),(131,193,1,'backcash',2.88,'订单:684439119(商家:美味汇)返现￥2.88元',1421419139,1),(132,194,1,'backcash',1.44,'订单:682733252(商家:美味汇)返现￥1.44元',1421419139,1),(133,195,1,'backcash',2.88,'订单:687082275(商家:美味汇)返现￥2.88元',1421419139,1),(134,196,1,'backcash',2.88,'订单:681830704(商家:美味汇)返现￥2.88元',1421419139,1),(135,197,1,'backcash',2.88,'订单:689460669(商家:美味汇)返现￥2.88元',1421419139,1),(136,198,1,'backcash',4.32,'订单:689127006(商家:美味汇)返现￥4.32元',1421419139,1),(137,199,1,'backcash',4.32,'订单:681246014(商家:美味汇)返现￥4.32元',1421419139,1),(138,200,1,'backcash',2.88,'订单:687298508(商家:美味汇)返现￥2.88元',1421419139,1),(139,201,1,'backcash',2.96,'订单:683810426(商家:美味汇)返现￥2.96元',1421419139,1),(140,202,1,'backcash',2.88,'订单:687535196(商家:美味汇)返现￥2.88元',1421419139,1),(141,203,1,'backcash',1.44,'订单:688886497(商家:美味汇)返现￥1.44元',1422947630,1),(142,204,1,'backcash',1.44,'订单:686076719(商家:美味汇)返现￥1.44元',1422947691,1);
/*!40000 ALTER TABLE `mm_income_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mm_integral_log`
--

DROP TABLE IF EXISTS `mm_integral_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mm_integral_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pt_id` int(11) DEFAULT NULL,
  `member_id` int(11) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `integral` int(11) DEFAULT NULL,
  `log` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `record_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=81 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_integral_log`
--

LOCK TABLES `mm_integral_log` WRITE;
/*!40000 ALTER TABLE `mm_integral_log` DISABLE KEYS */;
INSERT INTO `mm_integral_log` VALUES (1,666888,1,3,600,'订单返积分600个',2014),(2,666888,1,3,550,'订单返积分550个',2014),(3,666888,1,3,180,'订单返积分180个',2014),(4,666888,1,3,360,'订单返积分360个',2014),(5,666888,1,3,1260,'订单返积分1260个',2014),(6,666888,1,3,1260,'订单返积分1260个',2014),(7,666888,1,3,1260,'订单返积分1260个',2014),(8,666888,1,3,1260,'订单返积分1260个',2014),(9,666888,1,3,1260,'订单返积分1260个',2014),(10,666888,1,3,1260,'订单返积分1260个',2014),(11,666888,1,3,1260,'订单返积分1260个',2014),(12,666888,1,3,540,'订单返积分540个',2014),(13,666888,1,3,360,'订单返积分360个',2014),(14,666888,1,3,360,'订单返积分360个',2014),(15,666888,28,3,180,'订单返积分180个',2014),(16,666888,1,3,180,'订单返积分180个',2014),(17,666888,1,3,360,'订单返积分360个',2014),(18,666888,1,3,360,'订单返积分360个',2014),(19,666888,1,3,360,'订单返积分360个',2014),(20,666888,1,3,360,'订单返积分360个',2014),(21,666888,1,3,360,'订单返积分360个',NULL),(22,666888,1,3,360,'订单返积分360个',NULL),(23,666888,1,3,360,'订单返积分360个',NULL),(24,666888,1,3,360,'订单返积分360个',NULL),(25,666888,1,3,180,'订单返积分180个',NULL),(26,666888,1,3,360,'订单返积分360个',NULL),(27,666888,1,3,180,'订单返积分180个',NULL),(28,666888,1,3,360,'订单返积分360个',2015),(29,666888,1,3,360,'订单返积分360个',1421414462),(30,666888,1,3,360,'订单返积分360个',1421419029),(31,666888,1,3,360,'订单返积分360个',1421419029),(32,666888,1,3,550,'订单返积分550个',1421419029),(33,666888,1,3,180,'订单返积分180个',1421419029),(34,666888,1,3,180,'订单返积分180个',1421419029),(35,666888,1,3,180,'订单返积分180个',1421419029),(36,666888,1,3,360,'订单返积分360个',1421419029),(37,666888,1,3,360,'订单返积分360个',1421419030),(38,666888,1,3,180,'订单返积分180个',1421419030),(39,666888,1,3,360,'订单返积分360个',1421419030),(40,666888,1,3,360,'订单返积分360个',1421419030),(41,666888,1,3,360,'订单返积分360个',1421419030),(42,666888,1,3,360,'订单返积分360个',1421419030),(43,666888,1,3,360,'订单返积分360个',1421419030),(44,666888,1,3,360,'订单返积分360个',1421419030),(45,666888,1,3,360,'订单返积分360个',1421419030),(46,666888,1,3,360,'订单返积分360个',1421419030),(47,666888,1,3,360,'订单返积分360个',1421419030),(48,666888,1,3,540,'订单返积分540个',1421419030),(49,666888,1,3,1260,'订单返积分1260个',1421419030),(50,666888,1,3,360,'订单返积分360个',1421419030),(51,666888,28,3,180,'订单返积分180个',1421419031),(52,666888,1,3,180,'订单返积分180个',1421419031),(53,666888,1,3,180,'订单返积分180个',1421419080),(54,666888,1,3,180,'订单返积分180个',1421419138),(55,666888,1,3,390,'订单返积分390个',1421419138),(56,666888,1,3,360,'订单返积分360个',1421419138),(57,666888,1,3,360,'订单返积分360个',1421419138),(58,666888,1,3,540,'订单返积分540个',1421419138),(59,666888,1,3,180,'订单返积分180个',1421419138),(60,666888,1,3,180,'订单返积分180个',1421419139),(61,666888,1,3,180,'订单返积分180个',1421419139),(62,666888,1,3,180,'订单返积分180个',1421419139),(63,666888,1,3,540,'订单返积分540个',1421419139),(64,666888,1,3,540,'订单返积分540个',1421419139),(65,666888,1,3,360,'订单返积分360个',1421419139),(66,666888,1,3,550,'订单返积分550个',1421419139),(67,666888,1,3,540,'订单返积分540个',1421419139),(68,666888,1,3,360,'订单返积分360个',1421419139),(69,666888,1,3,360,'订单返积分360个',1421419139),(70,666888,1,3,180,'订单返积分180个',1421419139),(71,666888,1,3,360,'订单返积分360个',1421419139),(72,666888,1,3,360,'订单返积分360个',1421419139),(73,666888,1,3,360,'订单返积分360个',1421419139),(74,666888,1,3,540,'订单返积分540个',1421419139),(75,666888,1,3,540,'订单返积分540个',1421419139),(76,666888,1,3,360,'订单返积分360个',1421419139),(77,666888,1,3,370,'订单返积分370个',1421419139),(78,666888,1,3,360,'订单返积分360个',1421419140),(79,666888,1,3,180,'订单返积分180个',1422947630),(80,666888,1,3,180,'订单返积分180个',1422947691);
/*!40000 ALTER TABLE `mm_integral_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mm_member`
--

DROP TABLE IF EXISTS `mm_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mm_member` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `usr` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '用户名',
  `pwd` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '密码',
  `name` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '名字',
  `exp` int(11) unsigned DEFAULT '0',
  `level` int(11) DEFAULT '1',
  `sex` int(1) DEFAULT NULL COMMENT '性别(0: 未知,1:男,2：女)',
  `avatar` varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL,
  `birthday` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `phone` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '送餐地址',
  `qq` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  `email` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `reg_time` int(11) DEFAULT NULL,
  `reg_ip` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `last_login_time` int(11) DEFAULT NULL COMMENT '最后登录时间',
  `state` int(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=29 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_member`
--

LOCK TABLES `mm_member` WRITE;
/*!40000 ALTER TABLE `mm_member` DISABLE KEYS */;
INSERT INTO `mm_member` VALUES (1,'newmin','768dd5c54e40bcd412f7277cbe77423e','刘铭',3553,5,1,'','2012-09-05','18616999822','上海市徐汇区江安路15号1','18867761','q121276@26.com',2012,'112.65.35.191',2015,1),(2,'lin','22c5b45021a538af170f330bf6d6e46c','林德意',NULL,1,0,'','1970.05。29','13924886758','佛山禅城人民路99号',NULL,NULL,2012,'61.145.69.27',2012,1),(3,'hshx','c99425b0a379ac6621adbc0ce4170af5','黄升鑫',NULL,1,0,'','1979-12-01','15602817110','广东省佛山市禅城区人民路鹤园路81号','809987822','809987822@qq.com',2012,'27.36.72.124',1,1),(4,'lindeyi','0f167fd9c5f48d81820b544e312e8592','林德意',NULL,1,0,'','1970-05-29','13924886758','佛山市禅城区人民路99号','569101942','lindeyi158@yahoo.com.cn',2012,'183.27.197.170',1,1),(5,'sumian','4397d538520a9a645aa456e60744c1e0','',NULL,1,0,'','','13924886758','福建省海天大夏二栋201','','',2012,'14.157.18.39',1,1),(6,'linsu','464bf5d58f4e8818671d525cf1530459','',NULL,1,0,'','','13924886758','广州市黄埔去电子大夏3楼301室','','',2012,'14.157.18.39',1,1),(7,'sonven','b99831c5e69ac900fdfcfd4c7d0bf89e','',NULL,1,0,'','','18616999822','上海市浦东新区浦电路123弄','','',2012,'183.250.3.128',1,1),(8,'yangbo','1df2ee43288507769a15da6cb1cf0dba','',NULL,1,0,'','','18616888888','佛山市禅城区汾江中路20号','','',2012,'183.250.3.128',1,1),(9,'xiaoyuan','db5896d7e1951418a6fe0de4ea86b45b','',NULL,1,0,'','','18616999822','上海市徐汇区飞掉国际大厦','','',2013,'218.85.143.146',1,1),(10,'liuming','06f267d8e85c3478e00a8b9d2bae5df4','',NULL,1,0,'','','18616999822','上海市徐汇区飞掉国际大厦','','',2013,'218.85.143.146',1,1),(11,'xiaoyuanyaun','dda4f29b5f09313383fcfc02c0ce2753','',NULL,1,0,'','','18616999822','上海市徐汇区飞掉国际大厦','','',2013,'218.85.143.146',1,1),(12,'liuxiaoyuan','b201ab396944544760bd6b19d356cc8f','',NULL,1,0,'','','18616999822','上海市徐汇区飞掉国际大厦','','',2013,'218.85.143.146',1,1),(13,'yuan','86ad59c09a3f4fe980b67b9dedea7329','',NULL,1,0,'','','13728501775','佛山市禅城区张槎四路东大街3号5楼','','',2013,'183.27.195.23',1,1),(14,'yuanyuan','40c39d9211d17d1e6732059427c9ee76','',NULL,1,0,'','','13728501775','佛山禅城区张槎四路岗头东大街3号5楼','','',2013,'183.27.199.29',1,1),(15,'13924886758','c545e63671045c96669b814901bf0d37','',NULL,1,0,'','','13924886758','佛山市禅城区人民路99号','','',2013,'183.28.79.121',1,1),(16,'13728501775','9f86434fc3c081f7548e633c7ccdc5d2','',NULL,1,0,'','','13728501775','佛山市禅城区张槎四路（东海明珠后方）岗头东大街3号5楼','','',2013,'183.27.46.207',1,1),(25,'sa','123','刘铭',NULL,1,2,NULL,'1970-11-20','18616999822',NULL,NULL,NULL,2014,'127.0.0.1',2014,1),(26,'test','a50f4d2b5d08eca0ff83448fc346dbd6','测试员',NULL,1,1,NULL,'1988-11-09','18616999822',NULL,NULL,NULL,2014,'127.0.0.1',2014,1),(27,'test001','4dca21a567d5ae25316f9e8d37d8df1b','刘大炮',NULL,1,0,'share/noavatar.gif','1970-01-01','18616999822',NULL,NULL,NULL,2014,'127.0.0.1',2014,1),(28,'newmin123','9315871be89146db634ef0d0e5e181f9','刘大也',27,1,0,'share/noavatar.gif','1970-01-01','18616999822',NULL,NULL,NULL,2014,'127.0.0.1',2014,1);
/*!40000 ALTER TABLE `mm_member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `mm_relation`
--

DROP TABLE IF EXISTS `mm_relation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mm_relation` (
  `member_id` int(11) NOT NULL,
  `card_id` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `tg_id` int(11) DEFAULT NULL COMMENT '推广会员ID',
  `reg_ptid` int(11) DEFAULT NULL,
  PRIMARY KEY (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_relation`
--

LOCK TABLES `mm_relation` WRITE;
/*!40000 ALTER TABLE `mm_relation` DISABLE KEYS */;
INSERT INTO `mm_relation` VALUES (1,'',0,666888),(2,'',0,666888),(28,'201412191655',0,666888),(29,'',28,666888),(30,'',29,666888),(31,'',2,666888),(32,'',2,666888),(33,'',2,666888),(34,'',2,666888),(35,'',2,666888),(36,'',2,666888),(37,'',2,666888),(38,'',2,666888),(39,'',2,666888),(40,'',2,666888),(41,'',39,666888);
/*!40000 ALTER TABLE `mm_relation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pm_coupon`
--

DROP TABLE IF EXISTS `pm_coupon`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pm_coupon` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pt_id` int(11) DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '优惠码',
  `description` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `amount` int(11) DEFAULT NULL COMMENT '优惠码可用数量',
  `total_amount` int(11) DEFAULT NULL,
  `fee` int(11) DEFAULT NULL COMMENT '包含金额',
  `discount` int(11) DEFAULT NULL,
  `integral` int(11) DEFAULT NULL COMMENT '包含积分',
  `min_level` int(11) DEFAULT NULL COMMENT '等级限制',
  `min_fee` int(11) DEFAULT NULL COMMENT '订单金额限制',
  `begin_time` int(11) DEFAULT NULL,
  `over_time` int(11) DEFAULT NULL,
  `allow_enable` tinyint(1) DEFAULT NULL COMMENT '是否允许使用',
  `need_bind` tinyint(1) DEFAULT NULL COMMENT '是否需要绑定',
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL COMMENT '共计数量',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pm_coupon`
--

LOCK TABLES `pm_coupon` WRITE;
/*!40000 ALTER TABLE `pm_coupon` DISABLE KEYS */;
INSERT INTO `pm_coupon` VALUES (1,666888,'30off',NULL,0,10,1,70,30,0,0,2014,2014,1,1,2014,2014),(2,666888,'1WEEK',NULL,10,10,10,100,0,0,100,2015,2015,0,0,2014,2014),(3,666888,'10off',NULL,10,10,4,100,0,0,20,2015,2015,0,0,2014,2014),(4,666888,'dsss',NULL,10,10,0,100,10,0,0,2015,2015,0,0,2014,2014),(5,666888,'95off','95折全场通用',12,12,0,95,0,0,0,2014,2014,1,0,2014,2014);
/*!40000 ALTER TABLE `pm_coupon` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pm_coupon_bind`
--

DROP TABLE IF EXISTS `pm_coupon_bind`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pm_coupon_bind` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) DEFAULT NULL COMMENT '会员编号',
  `coupon_id` int(11) DEFAULT NULL COMMENT '优惠券编号',
  `is_used` tinyint(1) DEFAULT '0' COMMENT '是否使用',
  `bind_time` int(11) DEFAULT NULL COMMENT '绑定时间',
  `use_time` int(11) DEFAULT NULL COMMENT '使用时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pm_coupon_bind`
--

LOCK TABLES `pm_coupon_bind` WRITE;
/*!40000 ALTER TABLE `pm_coupon_bind` DISABLE KEYS */;
/*!40000 ALTER TABLE `pm_coupon_bind` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pm_coupon_take`
--

DROP TABLE IF EXISTS `pm_coupon_take`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pm_coupon_take` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) DEFAULT NULL,
  `coupon_id` int(11) DEFAULT NULL,
  `is_apply` tinyint(1) DEFAULT NULL COMMENT '是否生效,1表示有效',
  `take_time` int(11) DEFAULT NULL COMMENT '占用时间',
  `extra_time` int(11) DEFAULT NULL COMMENT '释放时间,超过该时间，优惠券释放',
  `apply_time` int(11) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pm_coupon_take`
--

LOCK TABLES `pm_coupon_take` WRITE;
/*!40000 ALTER TABLE `pm_coupon_take` DISABLE KEYS */;
/*!40000 ALTER TABLE `pm_coupon_take` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_order`
--

DROP TABLE IF EXISTS `pt_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_order` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_no` varchar(15) COLLATE utf8_unicode_ci NOT NULL,
  `member_id` int(11) DEFAULT NULL COMMENT '-1代表游客订餐',
  `pt_id` int(11) DEFAULT NULL COMMENT '商家ID',
  `shop_id` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '商家分店ID, 0为未指定，需管理指定',
  `items_info` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `total_fee` decimal(10,2) DEFAULT NULL COMMENT '订单总额',
  `fee` decimal(10,2) DEFAULT NULL COMMENT '订单实际金额',
  `discount_fee` decimal(10,2) DEFAULT NULL COMMENT '优惠金额',
  `coupon_fee` decimal(10,2) DEFAULT NULL COMMENT '优惠券优惠金额',
  `pay_fee` decimal(10,2) DEFAULT '0.00' COMMENT '支付金额',
  `pay_method` int(11) DEFAULT NULL COMMENT '1:餐到付款 2:网上支付  ',
  `is_suspend` tinyint(4) DEFAULT '0',
  `is_paid` int(11) DEFAULT NULL COMMENT '是否支付(0:未支付 ，1：已支付)',
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `remark` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '订单备注(可为取消理由)',
  `deliver_name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `deliver_phone` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `deliver_address` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `deliver_time` int(11) DEFAULT NULL COMMENT '送餐时间',
  `paid_time` int(11) DEFAULT NULL COMMENT '支付时间',
  `status` tinyint(4) DEFAULT NULL,
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=210 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_order`
--

LOCK TABLES `pt_order` WRITE;
/*!40000 ALTER TABLE `pt_order` DISABLE KEYS */;
INSERT INTO `pt_order` VALUES (1,'20121103-1971',1,666888,'0','尖椒回锅肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(2,'20120924-9603',1,666888,'1','韭黄炒蛋饭() * 9\n香菇焖鸡饭() * 9\n酸甜排骨饭() * 9\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(3,'20121027-8653',1,666888,'0','酸甜排骨饭() * 9\n极品红烧肉() * 9\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'xiexie取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(4,'20120924-3892',1,666888,'1','酸甜排骨饭() * 9\n香菇焖鸡饭() * 9\n韭黄炒蛋饭() * 9\n极品红烧肉() * 9\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(5,'20120627-2436',15,666888,'1','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2012,1421419029),(6,'20120923-4551',1,666888,'1','酸甜排骨饭() * 9\n香菇焖鸡饭() * 9\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(7,'20120922-8308',1,666888,'0','酸甜排骨饭() * 10\n香菇焖鸡饭() * 10\n韭黄炒蛋饭() * 10\n极品红烧肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(8,'20120627-1100',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(9,'20120627-7755',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(10,'20120627-6581',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,0,NULL,'',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(11,'20120627-8125',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,0,NULL,'',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(12,'20120627-4651',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(13,'20120627-2526',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(14,'20120627-7442',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,0,NULL,'',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(15,'20120627-7080',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(16,'20120627-8196',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(17,'20120627-5649',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(18,'20120627-6550',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(19,'20120627-6220',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(20,'20120627-1654',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(21,'20120627-8171',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(22,'20120627-2550',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(23,'20120627-7952',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(24,'20120627-6759',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(25,'20120627-8506',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(26,'20120627-5911',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(27,'20120627-7876',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(28,'20120627-8346',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(29,'20120627-6117',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(30,'20120627-1385',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(31,'20120627-4827',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(32,'20120627-7772',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(33,'20120627-1022',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(34,'20120627-1881',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(35,'20120627-9740',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(36,'20120627-3168',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(37,'20120627-4255',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(38,'20120627-4677',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(39,'20120627-4637',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(40,'20120627-7191',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(41,'20120627-1907',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(42,'20120627-1813',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(43,'20120627-5496',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(44,'20120627-2231',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(45,'20120627-3244',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(46,'20120627-1174',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(47,'20120627-1436',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(48,'20120627-7910',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(49,'20120627-3308',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(50,'20120627-5287',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(51,'20120627-5115',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(52,'20120627-4901',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(53,'20120627-7672',15,666888,'0','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(54,'20120627-6181',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(55,'20120627-1763',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(56,'20120627-9625',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(57,'20120627-4537',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(58,'20120627-8357',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(59,'20120627-1266',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(60,'20120627-8738',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(61,'20120627-9911',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(62,'20120627-2810',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(63,'20120627-4900',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(64,'20120627-1501',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(65,'20120627-7929',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(66,'20120627-5618',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(67,'20120627-6252',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(68,'20120627-9793',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(69,'20120627-5242',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(70,'20120627-2459',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(71,'20120627-6815',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(72,'20120627-8200',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(73,'20120627-3302',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(74,'20120627-9453',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(75,'20120627-2646',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(76,'20120627-6516',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(77,'20120627-8374',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(78,'20120627-5625',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(79,'20120627-1217',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(80,'20120627-6054',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(81,'20120627-3177',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(82,'20120627-7430',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(83,'20120627-8355',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(84,'20120627-8173',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(85,'20120627-3528',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(86,'20120627-7629',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(87,'20120627-9910',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(88,'20120627-5774',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(89,'20120627-2158',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(90,'20120627-6967',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(91,'20120627-8772',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(92,'20120627-9923',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(93,'20120627-2908',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(94,'20120627-1039',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(95,'20120627-9722',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(96,'20120627-2130',9,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(97,'20120627-3465',9,666888,'1','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(98,'20120702-9893',15,666888,'1','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(99,'20120702-7680',15,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(100,'20120702-9943',15,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(101,'20120702-9760',15,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(102,'20120702-3696',1,666888,'1','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(103,'20120703-4504',8,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(104,'20120703-8136',8,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(105,'20120703-3632',8,666888,'1','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(106,'20120703-4933',8,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(107,'20120703-2793',8,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(108,'20120703-5358',1,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(109,'20120703-1458',1,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(110,'20120703-9432',19,666888,'1','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(111,'20120703-7117',19,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(112,'20120707-2673',1,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(113,'20120707-7558',22,666888,'1','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(114,'20120707-7906',21,666888,'1','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(115,'20120707-8395',16,666888,'0','营养套餐 () * 10\n木耳炒黄瓜(优惠仅限3日) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(116,'20120707-6956',23,666888,'0','营养套餐 () * 10\n木耳炒黄瓜(优惠仅限3日) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(117,'20120707-4722',23,666888,'0','过桥米线() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(118,'20120707-9913',16,666888,'0','酸辣粉丝() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(119,'20120707-1485',16,666888,'0','酸辣粉丝() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(120,'20120707-8310',23,666888,'0','酸辣粉丝() * 10\n木耳炒黄瓜(优惠仅限3日) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(121,'20120707-9986',23,666888,'0','营养套餐 () * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(122,'20120709-2715',1,666888,'0','木耳炒黄瓜(优惠仅限3日) * 10\n麻婆豆腐() * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n营养套餐 () * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(123,'20120710-4124',1,666888,'0','糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(124,'20120710-7520',1,666888,'0','商务营养套餐(：\"优惠仅限3日) * 168\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(125,'20120710-9269',1,666888,'1','木耳炒黄瓜(优惠仅限3日) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(126,'20120810-9477',1,666888,'4','鼓汁排骨饭((仅限今日8折优惠)) * 10.4\n西兰花炒肉片 () * 10\n荷兰豆炒肉片() * 10\n木耳炒黄瓜(优惠仅限3日) * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2012,NULL,0,2012,2012),(127,'20130115-7948',1,666888,'1','野山椒炒肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(128,'20130115-3856',1,666888,'0','野山椒炒肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'取消原因:超时未付款，系统取消',NULL,NULL,NULL,NULL,2013,NULL,0,2013,2013),(129,'20130203-1190',1,666888,'2','野山椒炒肉() * 10\n尖椒回锅肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(130,'20130307-5447',1,666888,'1','尖椒回锅肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'测试',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(131,'20130308-2932',1,666888,'1','酸甜排骨饭() * 9\n极品红烧肉() * 9\n香菇焖鸡饭() * 9\n花生猪手饭() * 10\n鱼香茄子饭() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'要辣',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(132,'20130311-9708',1,666888,'1','尖椒回锅肉() * 10\n野山椒炒肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(133,'20130315-1432',1,666888,'7','尖椒回锅肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(134,'20130331-4922',1,666888,'1','尖椒回锅肉() * 9\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'要辣一点的',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(135,'20130331-4998',1,666888,'1','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(136,'20130331-2185',1,666888,'1','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(137,'20130401-8432',1,666888,'1','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(138,'20130401-5283',1,666888,'1','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(139,'20130401-4602',1,666888,'1','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(140,'20130401-7643',1,666888,'1','野山椒炒肉() * 10\n尖椒回锅肉() * 9\n极品红烧肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(141,'20130401-8214',1,666888,'1','野山椒炒肉() * 10\n尖椒回锅肉() * 9\n极品红烧肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(142,'20130401-8460',1,666888,'1','野山椒炒肉() * 10\n极品红烧肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(143,'20130401-2503',1,666888,'2','香卤鸡腿饭() * 10\n花生猪手饭() * 10\n红烧鱼腩饭() * 9\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2013,0,0,2013,1421419029),(144,'20130408-9481',1,666888,'1','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,0.00,0.00,60.00,1,0,0,'',NULL,NULL,NULL,NULL,2014,0,0,2013,1421419029),(145,'20130408-9473',1,666888,'0','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,0,NULL,'',NULL,NULL,NULL,NULL,2013,NULL,0,2013,2013),(146,'685450864',1,666888,NULL,'鱼香茄子饭*1\n尖椒回锅肉*2\n香菇焖鸡饭*1',80.00,72.00,NULL,NULL,72.00,1,0,NULL,NULL,NULL,NULL,NULL,NULL,2014,NULL,NULL,2014,2014),(147,'683515854',1,666888,NULL,'鱼香茄子饭*1\n尖椒回锅肉*2\n香菇焖鸡饭*1',80.00,72.00,NULL,NULL,72.00,1,0,NULL,NULL,NULL,NULL,NULL,NULL,2014,NULL,NULL,2014,2014),(148,'685952634',1,666888,NULL,'鱼香茄子饭*18\n尖椒回锅肉*18\n香菇焖鸡饭*18',1080.00,972.00,NULL,NULL,972.00,2,0,NULL,NULL,NULL,NULL,NULL,NULL,2014,NULL,NULL,2014,2014),(149,'681859957',1,666888,NULL,'鱼香茄子饭*2',40.00,36.00,NULL,NULL,36.00,2,0,NULL,NULL,NULL,NULL,NULL,NULL,2014,NULL,NULL,2014,2014),(150,'686802891',1,666888,NULL,'鱼香茄子饭*2',40.00,36.00,NULL,NULL,36.00,2,0,NULL,NULL,NULL,NULL,NULL,NULL,2014,NULL,NULL,2014,2014),(151,'686267421',1,666888,'1','鱼香茄子饭*2',40.00,36.00,0.00,0.00,36.00,2,0,1,NULL,NULL,NULL,NULL,NULL,2014,1421418460,6,2014,1421419029),(152,'684452597',1,666888,'7','鱼香茄子饭*2',40.00,36.00,0.00,0.00,36.00,2,0,1,NULL,NULL,NULL,NULL,NULL,2014,1421418526,6,2014,1421419029),(153,'689182668',1,666888,'7','鱼香茄子饭*1\n尖椒回锅肉*1\n香菇焖鸡饭*1\n野山椒炒肉*1',60.00,55.00,0.00,0.00,55.00,2,0,1,NULL,NULL,NULL,NULL,NULL,2014,1421418626,6,2014,1421419029),(154,'685815333',1,666888,'1','鱼香茄子饭*1',20.00,18.00,0.00,0.00,18.00,2,0,1,NULL,NULL,NULL,NULL,NULL,2014,1421418633,6,2014,1421419029),(155,'689681862',1,666888,'0','鱼香茄子饭*1',20.00,18.00,0.00,0.00,18.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2014,1421419138),(156,'686862923',1,666888,'0','鱼香茄子饭*1\n香菇焖鸡饭*1\n野山椒炒肉*1\n极品红烧肉*1\n花生猪手饭*1',40.00,39.00,0.00,0.00,39.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2014,1421419138),(157,'687643214',1,666888,'0','尖椒回锅肉*1\n香菇焖鸡饭*1',40.00,36.00,0.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2014,1421419138),(158,'687137477',1,666888,'0','鱼香茄子饭*1\n尖椒回锅肉*1',40.00,36.00,0.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2014,1421419138),(159,'686200508',1,666888,'0','鱼香茄子饭*1\n尖椒回锅肉*1\n香菇焖鸡饭*1',60.00,54.00,0.00,0.00,54.00,2,0,0,NULL,NULL,'李讪讪','18616999822','上海市黄埔区',1421419091,0,6,2014,1421419138),(160,'682658763',1,666888,'0','鱼香茄子饭*1',20.00,18.00,0.00,0.00,18.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2014,1421419138),(161,'688563456',1,666888,'0','香菇焖鸡饭*1',20.00,18.00,0.00,0.00,18.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2014,1421419139),(162,'686432330',1,666888,'0','尖椒回锅肉*1',20.00,18.00,0.00,0.00,18.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2014,1421419139),(163,'682539015',1,666888,'0','鱼香茄子饭*1',20.00,18.00,0.00,0.00,18.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2014,1421419139),(164,'683135658',1,666888,'1','尖椒回锅肉*1',20.00,18.00,0.00,0.00,18.00,2,0,1,NULL,NULL,'刘铭','18616999822','上海市',2014,1421418639,6,2014,1421419029),(165,'689424176',1,666888,'1','尖椒回锅肉*1',20.00,18.00,7.40,5.40,13.00,2,0,1,NULL,NULL,'刘铭','18616999822','上海市',1421413341,1421418380,6,2014,1421419029),(166,'682964508',1,666888,'1','尖椒回锅肉*1\n尖椒回锅肉*1',40.00,36.00,4.00,0.00,36.00,2,0,1,NULL,NULL,'刘铭','18616999822','上海市',1421413353,1421418267,6,2014,1421419029),(167,'689386565',1,666888,'1','香菇焖鸡饭*2',40.00,36.00,14.80,10.80,25.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(168,'685350377',1,666888,'1','香菇焖鸡饭*2',40.00,36.00,14.80,10.80,25.00,2,0,0,'取消原因:无法联系送达',NULL,'刘铭','18616999822','上海市',2014,NULL,0,2014,2014),(169,'683401193',1,666888,'1','尖椒回锅肉*2',40.00,36.00,14.80,10.80,25.00,2,0,0,'取消原因:会员申请取消',NULL,'刘铭','18616999822','上海市',2014,NULL,0,2014,2014),(170,'681554163',1,666888,'1','鱼香茄子饭*1',20.00,18.00,7.40,5.40,13.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(171,'689760500',1,666888,'1','尖椒回锅肉*1\n尖椒回锅肉*1',40.00,36.00,14.80,10.80,25.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(172,'687956073',1,666888,'1','香菇焖鸡饭*1\n香菇焖鸡饭*1',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(173,'684691894',1,666888,'1','香菇焖鸡饭*2',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(174,'682466077',1,666888,'1','鱼香茄子饭*2',40.00,36.00,5.80,1.80,34.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(175,'681009722',1,666888,'1','鱼香茄子饭*2',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(176,'681710378',1,666888,'1','鱼香茄子饭*2',40.00,36.00,17.80,13.80,22.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(177,'687025732',1,666888,'1','鱼香茄子饭*2',40.00,36.00,16.80,12.80,23.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(178,'683878066',1,666888,'1','鱼香茄子饭*2',40.00,36.00,15.80,11.80,24.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(179,'681645404',1,666888,'7','鱼香茄子饭*2',40.00,36.00,15.80,11.80,24.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(180,'687153988',1,666888,'1','鱼香茄子饭*3',60.00,54.00,23.20,17.20,37.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(181,'688664284',1,666888,'1','鱼香茄子饭*7',140.00,126.00,20.30,6.30,119.70,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(182,'684148199',1,666888,'1','香菇焖鸡饭*2',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419030),(183,'689469831',28,666888,'1','鱼香茄子饭*1',20.00,18.00,2.00,0.00,18.00,2,0,0,NULL,NULL,'刘大爷','18616999822','佛山市',2014,0,6,2014,1421419030),(184,'689205996',1,666888,'1','鱼香茄子饭*1',20.00,18.00,2.00,0.00,18.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',2014,0,6,2014,1421419031),(185,'683958933',1,666888,'0','鱼香茄子饭*3',60.00,54.00,6.00,0.00,54.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,0,1421419139),(186,'686018039',1,666888,'0','香菇焖鸡饭*1\n香菇焖鸡饭*1\n香菇焖鸡饭*1',60.00,54.00,6.00,0.00,54.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2015,1421419139),(187,'681426294',1,666888,'0','尖椒回锅肉*2',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2015,1421419139),(188,'686396270',1,666888,'0','鱼香茄子饭*1\n尖椒回锅肉*1\n香菇焖鸡饭*1\n野山椒炒肉*1',60.00,55.00,5.00,0.00,55.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2015,1421419139),(189,'686660811',1,666888,'0','鱼香茄子饭*1\n尖椒回锅肉*1\n香菇焖鸡饭*1',60.00,54.00,6.00,0.00,54.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,2015,1421419139),(190,'685994267',1,666888,'1','鱼香茄子饭*1\n尖椒回锅肉*1',40.00,36.00,4.00,0.00,36.00,2,0,1,NULL,NULL,'刘铭','18616999822','上海市',2015,NULL,6,2015,1421414462),(191,'683249100',1,666888,'1','鱼香茄子饭*1',20.00,18.00,2.00,0.00,18.00,2,0,1,NULL,NULL,'刘铭','18616999822','上海市',1421419078,1421419071,6,2015,1421419079),(192,'681879385',1,666888,'0','鱼香茄子饭*2',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421115773,1421419139),(193,'684439119',1,666888,'0','鱼香茄子饭*2',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421116017,1421419139),(194,'682733252',1,666888,'0','鱼香茄子饭*1',20.00,18.00,2.00,0.00,18.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421116540,1421419139),(195,'687082275',1,666888,'0','鱼香茄子饭*2',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421117123,1421419139),(196,'681830704',1,666888,'0','鱼香茄子饭*1\n尖椒回锅肉*1',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421117237,1421419139),(197,'689460669',1,666888,'0','鱼香茄子饭*2',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421117290,1421419139),(198,'689127006',1,666888,'0','鱼香茄子饭*1\n尖椒回锅肉*1\n香菇焖鸡饭*1',60.00,54.00,6.00,0.00,54.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421117367,1421419139),(199,'681246014',1,666888,'0','鱼香茄子饭*3',60.00,54.00,6.00,0.00,54.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421126427,1421419139),(200,'687298508',1,666888,'0','鱼香茄子饭*1\n尖椒回锅肉*1',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421126653,1421419139),(201,'683810426',1,666888,'0','尖椒回锅肉*1\n香菇焖鸡饭*1\n野山椒炒肉*1',40.00,37.00,3.00,0.00,37.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421126896,1421419139),(202,'687535196',1,666888,'0','鱼香茄子饭*1\n尖椒回锅肉*1',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1421419091,0,6,1421127721,1421419139),(203,'688886497',1,666888,'0','鱼香茄子饭*1',20.00,18.00,2.00,0.00,18.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1422947510,0,6,1422945619,1422947630),(204,'686076719',1,666888,'0','鱼香茄子饭*1',20.00,18.00,2.00,0.00,18.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1422947570,0,6,1422945675,1422947691),(205,'681656585',1,666888,'0','鱼香茄子饭(1)*32\n尖椒回锅肉*3\n香菇焖鸡饭*1\n野山椒炒肉*1\n极品红烧肉*1',720.00,650.00,70.00,0.00,650.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1423386000,0,4,1423384138,1423386000),(206,'686262415',1,666888,'0','鱼香茄子饭(1)*33\n尖椒回锅肉*4\n香菇焖鸡饭*2\n野山椒炒肉*2\n极品红烧肉*1',780.00,705.00,75.00,0.00,705.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1423382578,0,2,1423386178,1423386444),(207,'688714535',1,666888,'0','鱼香茄子饭(1)*33\n尖椒回锅肉*4\n香菇焖鸡饭*2\n野山椒炒肉*2\n极品红烧肉*1',780.00,705.00,75.00,0.00,705.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1423382676,0,2,1423386276,1423386609),(208,'685733757',1,666888,'0','花生猪手饭*40\n红烧鱼腩饭*4\n香卤鸡腿饭*2\n酸甜排骨饭*2\n',920.00,831.00,89.00,0.00,831.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1423382794,0,2,1423386394,1423386770),(209,'681884165',1,666888,'0','韭黄炒蛋饭-2-2*1\n鱼香茄子饭(1)*1',40.00,36.00,4.00,0.00,36.00,2,0,0,NULL,NULL,'刘铭','18616999822','上海市',1423382950,0,2,1423386550,1423386826);
/*!40000 ALTER TABLE `pt_order` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_order_coupon`
--

DROP TABLE IF EXISTS `pt_order_coupon`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_order_coupon` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` int(11) DEFAULT NULL,
  `coupon_id` int(11) DEFAULT NULL,
  `coupon_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `coupon_fee` decimal(10,2) DEFAULT NULL,
  `coupon_describe` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `send_integral` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_order_coupon`
--

LOCK TABLES `pt_order_coupon` WRITE;
/*!40000 ALTER TABLE `pt_order_coupon` DISABLE KEYS */;
INSERT INTO `pt_order_coupon` VALUES (1,168,1,'30off',11.00,'任意订单,7.0折优惠,赠送积分30点',30),(2,169,1,'30off',11.00,'任意订单,7.0折优惠,赠送积分30点',30),(3,170,1,'30off',5.00,'任意订单,7.0折优惠,赠送积分30点',30),(4,171,1,'30off',11.00,'任意订单,7.0折优惠,赠送积分30点',30),(5,174,5,'95off',2.00,'任意订单,9.5折优惠',0),(6,176,5,'95off',14.00,'任意订单,9.5折优惠,另减12元',0),(7,177,5,'95off',13.00,'任意订单,9.5折优惠,另减11元',0),(8,178,5,'95off',12.00,'任意订单,9.5折优惠,另减10元',0),(9,179,1,'30off',12.00,'任意订单,7.0折优惠,另减1元,赠送积分30点',30),(10,180,1,'30off',17.00,'任意订单,7.0折优惠,另减1元,赠送积分30点',30),(11,181,5,'95off',6.00,'任意订单,9.5折优惠',0);
/*!40000 ALTER TABLE `pt_order_coupon` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_order_item`
--

DROP TABLE IF EXISTS `pt_order_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_order_item` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `snapshot_id` int(11) DEFAULT NULL,
  `quantity` int(11) DEFAULT NULL,
  `sku` varchar(100) DEFAULT NULL,
  `fee` decimal(10,0) DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_order_item`
--

LOCK TABLES `pt_order_item` WRITE;
/*!40000 ALTER TABLE `pt_order_item` DISABLE KEYS */;
/*!40000 ALTER TABLE `pt_order_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_order_log`
--

DROP TABLE IF EXISTS `pt_order_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_order_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `order_id` int(11) DEFAULT NULL,
  `type` tinyint(1) DEFAULT '1' COMMENT '类型，１:流程,2:调价',
  `message` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `is_system` tinyint(4) DEFAULT NULL,
  `record_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=206 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_order_log`
--

LOCK TABLES `pt_order_log` WRITE;
/*!40000 ALTER TABLE `pt_order_log` DISABLE KEYS */;
INSERT INTO `pt_order_log` VALUES (1,187,1,'订单已经确认',NULL,2015),(2,188,1,'订单已经确认',NULL,2015),(3,190,1,'订单已经确认',NULL,2015),(4,190,1,'订单处理中',NULL,2015),(5,190,1,'订单开始配送',NULL,2015),(6,190,1,'订单已完成',NULL,2015),(7,189,1,'订单已经确认',NULL,2015),(8,191,1,'订单已经确认',0,1421412981),(9,191,1,'订单已锁定智能分配门店失败！原因：无法识别的地址：上海市',1,1421413006),(10,192,1,'订单已经确认',0,1421413006),(11,192,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413011),(12,193,1,'订单已经确认',0,1421413011),(13,193,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413011),(14,194,1,'订单已经确认',0,1421413011),(15,194,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413011),(16,195,1,'订单已经确认',0,1421413011),(17,195,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413012),(18,196,1,'订单已经确认',0,1421413012),(19,196,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413012),(20,197,1,'订单已经确认',0,1421413012),(21,197,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413012),(22,198,1,'订单已经确认',0,1421413012),(23,198,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413012),(24,199,1,'订单已经确认',0,1421413012),(25,199,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413012),(26,200,1,'订单已经确认',0,1421413012),(27,200,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413012),(28,201,1,'订单已经确认',0,1421413012),(29,201,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413013),(30,202,1,'订单已经确认',0,1421413013),(31,202,1,'自动分配门店:南庄店,电话：0757-36668888',0,1421413013),(32,165,1,'订单开始配送',0,1421413341),(33,165,1,'已收货',0,1421413345),(34,166,1,'订单开始配送',0,1421413353),(35,166,1,'已收货',0,1421413353),(36,190,1,'订单已完成',0,1421414462),(37,5,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(38,127,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(39,129,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(40,130,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(41,131,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(42,132,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(43,133,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(44,134,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(45,135,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(46,136,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(47,137,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(48,138,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(49,139,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(50,140,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(51,141,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(52,142,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(53,143,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(54,144,1,'订单已取消,原因：超时未付款，系统取消',1,1421419029),(55,151,1,'订单已完成',0,1421419029),(56,152,1,'订单已完成',0,1421419029),(57,153,1,'订单已完成',0,1421419029),(58,154,1,'订单已完成',0,1421419029),(59,155,1,'订单处理中',0,1421419029),(60,156,1,'订单处理中',0,1421419029),(61,157,1,'订单处理中',0,1421419029),(62,158,1,'订单处理中',0,1421419029),(63,159,1,'订单处理中',0,1421419029),(64,160,1,'订单处理中',0,1421419029),(65,161,1,'订单处理中',0,1421419029),(66,162,1,'订单处理中',0,1421419029),(67,163,1,'订单处理中',0,1421419029),(68,164,1,'订单已完成',0,1421419029),(69,165,1,'订单已完成',0,1421419029),(70,166,1,'订单已完成',0,1421419030),(71,167,1,'订单已完成',0,1421419030),(72,170,1,'订单已完成',0,1421419030),(73,171,1,'订单已完成',0,1421419030),(74,172,1,'订单已完成',0,1421419030),(75,173,1,'订单已完成',0,1421419030),(76,174,1,'订单已完成',0,1421419030),(77,175,1,'订单已完成',0,1421419030),(78,176,1,'订单已完成',0,1421419030),(79,177,1,'订单已完成',0,1421419030),(80,178,1,'订单已完成',0,1421419030),(81,179,1,'订单已完成',0,1421419030),(82,180,1,'订单已完成',0,1421419030),(83,181,1,'订单已完成',0,1421419030),(84,182,1,'订单已完成',0,1421419030),(85,183,1,'订单已完成',0,1421419031),(86,184,1,'订单已完成',0,1421419031),(87,185,1,'订单处理中',0,1421419031),(88,186,1,'订单处理中',0,1421419031),(89,187,1,'订单处理中',0,1421419031),(90,188,1,'订单处理中',0,1421419031),(91,189,1,'订单处理中',0,1421419031),(92,192,1,'订单处理中',0,1421419031),(93,193,1,'订单处理中',0,1421419031),(94,194,1,'订单处理中',0,1421419031),(95,195,1,'订单处理中',0,1421419031),(96,196,1,'订单处理中',0,1421419031),(97,197,1,'订单处理中',0,1421419031),(98,198,1,'订单处理中',0,1421419031),(99,199,1,'订单处理中',0,1421419031),(100,200,1,'订单处理中',0,1421419031),(101,201,1,'订单处理中',0,1421419031),(102,202,1,'订单处理中',0,1421419031),(103,191,1,'订单处理中',0,1421419077),(104,191,1,'订单开始配送',0,1421419078),(105,191,1,'已收货',0,1421419078),(106,191,1,'订单已完成',0,1421419080),(107,155,1,'订单开始配送',0,1421419091),(108,156,1,'订单开始配送',0,1421419091),(109,157,1,'订单开始配送',0,1421419091),(110,158,1,'订单开始配送',0,1421419091),(111,159,1,'订单开始配送',0,1421419091),(112,160,1,'订单开始配送',0,1421419091),(113,161,1,'订单开始配送',0,1421419091),(114,162,1,'订单开始配送',0,1421419091),(115,163,1,'订单开始配送',0,1421419091),(116,185,1,'订单开始配送',0,1421419091),(117,186,1,'订单开始配送',0,1421419091),(118,187,1,'订单开始配送',0,1421419091),(119,188,1,'订单开始配送',0,1421419091),(120,189,1,'订单开始配送',0,1421419091),(121,192,1,'订单开始配送',0,1421419091),(122,193,1,'订单开始配送',0,1421419091),(123,194,1,'订单开始配送',0,1421419091),(124,195,1,'订单开始配送',0,1421419091),(125,196,1,'订单开始配送',0,1421419091),(126,197,1,'订单开始配送',0,1421419091),(127,198,1,'订单开始配送',0,1421419091),(128,199,1,'订单开始配送',0,1421419091),(129,200,1,'订单开始配送',0,1421419091),(130,201,1,'订单开始配送',0,1421419091),(131,202,1,'订单开始配送',0,1421419091),(132,155,1,'已收货',0,1421419099),(133,156,1,'已收货',0,1421419099),(134,157,1,'已收货',0,1421419099),(135,158,1,'已收货',0,1421419099),(136,159,1,'已收货',0,1421419099),(137,160,1,'已收货',0,1421419099),(138,161,1,'已收货',0,1421419099),(139,162,1,'已收货',0,1421419099),(140,163,1,'已收货',0,1421419099),(141,185,1,'已收货',0,1421419099),(142,186,1,'已收货',0,1421419099),(143,187,1,'已收货',0,1421419099),(144,188,1,'已收货',0,1421419099),(145,189,1,'已收货',0,1421419099),(146,192,1,'已收货',0,1421419099),(147,193,1,'已收货',0,1421419099),(148,194,1,'已收货',0,1421419099),(149,195,1,'已收货',0,1421419099),(150,196,1,'已收货',0,1421419099),(151,197,1,'已收货',0,1421419099),(152,198,1,'已收货',0,1421419099),(153,199,1,'已收货',0,1421419099),(154,200,1,'已收货',0,1421419099),(155,201,1,'已收货',0,1421419099),(156,202,1,'已收货',0,1421419099),(157,155,1,'订单已完成',0,1421419138),(158,156,1,'订单已完成',0,1421419138),(159,157,1,'订单已完成',0,1421419138),(160,158,1,'订单已完成',0,1421419138),(161,159,1,'订单已完成',0,1421419138),(162,160,1,'订单已完成',0,1421419139),(163,161,1,'订单已完成',0,1421419139),(164,162,1,'订单已完成',0,1421419139),(165,163,1,'订单已完成',0,1421419139),(166,185,1,'订单已完成',0,1421419139),(167,186,1,'订单已完成',0,1421419139),(168,187,1,'订单已完成',0,1421419139),(169,188,1,'订单已完成',0,1421419139),(170,189,1,'订单已完成',0,1421419139),(171,192,1,'订单已完成',0,1421419139),(172,193,1,'订单已完成',0,1421419139),(173,194,1,'订单已完成',0,1421419139),(174,195,1,'订单已完成',0,1421419139),(175,196,1,'订单已完成',0,1421419139),(176,197,1,'订单已完成',0,1421419139),(177,198,1,'订单已完成',0,1421419139),(178,199,1,'订单已完成',0,1421419139),(179,200,1,'订单已完成',0,1421419139),(180,201,1,'订单已完成',0,1421419139),(181,202,1,'订单已完成',0,1421419140),(182,203,1,'订单已经确认',0,1422945902),(183,203,1,'自动分配门店:南庄店,电话：0757-36668888',0,1422945903),(184,204,1,'订单已经确认',0,1422945963),(185,204,1,'自动分配门店:南庄店,电话：0757-36668888',0,1422945963),(186,203,1,'订单处理中',0,1422946323),(187,204,1,'订单处理中',0,1422946383),(188,203,1,'订单开始配送',0,1422947510),(189,203,1,'已收货',0,1422947570),(190,204,1,'订单开始配送',0,1422947570),(191,203,1,'订单已完成',0,1422947631),(192,204,1,'已收货',0,1422947631),(193,204,1,'订单已完成',0,1422947691),(194,205,1,'订单已经确认',0,1423385535),(195,205,1,'自动分配门店:南庄店,电话：0757-36668888',0,1423385535),(196,205,1,'订单处理中',0,1423385574),(197,205,1,'订单开始配送',0,1423386000),(198,206,1,'订单已经确认',0,1423386444),(199,206,1,'自动分配门店:南庄店,电话：0757-36668888',0,1423386445),(200,207,1,'订单已经确认',0,1423386609),(201,207,1,'自动分配门店:南庄店,电话：0757-36668888',0,1423386610),(202,208,1,'订单已经确认',0,1423386770),(203,208,1,'自动分配门店:南庄店,电话：0757-36668888',0,1423386776),(204,209,1,'订单已经确认',0,1423386826),(205,209,1,'自动分配门店:南庄店,电话：0757-36668888',0,1423386826);
/*!40000 ALTER TABLE `pt_order_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_page`
--

DROP TABLE IF EXISTS `pt_page`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_page` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pt_id` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `type` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `content` varchar(5000) COLLATE utf8_unicode_ci DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='合作商页面';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_page`
--

LOCK TABLES `pt_page` WRITE;
/*!40000 ALTER TABLE `pt_page` DISABLE KEYS */;
INSERT INTO `pt_page` VALUES (1,'666888','notice','',2013);
/*!40000 ALTER TABLE `pt_page` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_partner`
--

DROP TABLE IF EXISTS `pt_partner`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_partner` (
  `id` int(11) NOT NULL,
  `usr` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `pwd` varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  `secret` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `logo` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `tel` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `phone` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `join_time` int(11) DEFAULT NULL,
  `expires_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  `login_time` int(11) DEFAULT NULL,
  `last_login_time` int(11) DEFAULT NULL COMMENT '标志',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_partner`
--

LOCK TABLES `pt_partner` WRITE;
/*!40000 ALTER TABLE `pt_partner` DISABLE KEYS */;
INSERT INTO `pt_partner` VALUES (666888,'wly','97ccd376043aedb077fc6336d8c5a27c','d435a520e50e960b','美味汇',NULL,'0757-82255311','18616999822','佛山市禅城区亲仁路白燕街9号201',2012,1466666666,1421890499,0,0);
/*!40000 ALTER TABLE `pt_partner` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_positions`
--

DROP TABLE IF EXISTS `pt_positions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_positions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `psid` int(11) DEFAULT NULL COMMENT '分店ID',
  `name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_positions`
--

LOCK TABLES `pt_positions` WRITE;
/*!40000 ALTER TABLE `pt_positions` DISABLE KEYS */;
INSERT INTO `pt_positions` VALUES (3,666888,'电器大厦','禅城区汾江中路20号'),(4,666888,'东建世纪广场','禅城区季华三路');
/*!40000 ALTER TABLE `pt_positions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_saleconf`
--

DROP TABLE IF EXISTS `pt_saleconf`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_saleconf` (
  `pt_id` int(11) NOT NULL,
  `cb_percent` float(4,2) DEFAULT NULL COMMENT '反现比例,0则不返现',
  `cb_tg1_percent` float(4,2) DEFAULT NULL COMMENT '一级比例',
  `cb_tg2_percent` float(4,2) DEFAULT NULL COMMENT '二级比例',
  `cb_member_percent` float(4,2) DEFAULT NULL COMMENT '会员比例',
  `ib_num` int(11) DEFAULT NULL COMMENT '每一元返多少积分',
  `ib_extra` int(11) DEFAULT NULL COMMENT '额外赠送积分',
  `auto_setup_order` tinyint(4) DEFAULT '0',
  PRIMARY KEY (`pt_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_saleconf`
--

LOCK TABLES `pt_saleconf` WRITE;
/*!40000 ALTER TABLE `pt_saleconf` DISABLE KEYS */;
INSERT INTO `pt_saleconf` VALUES (666888,0.10,0.10,0.20,0.80,10,NULL,1);
/*!40000 ALTER TABLE `pt_saleconf` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_shop`
--

DROP TABLE IF EXISTS `pt_shop`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_shop` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pt_id` int(11) DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `phone` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `location` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '坐标',
  `deliver_radius` int(11) DEFAULT NULL COMMENT '配送范围',
  `order_index` int(11) DEFAULT '0',
  `state` int(11) DEFAULT NULL COMMENT '0:表示禁用   1:表示正常',
  `create_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_shop`
--

LOCK TABLES `pt_shop` WRITE;
/*!40000 ALTER TABLE `pt_shop` DISABLE KEYS */;
INSERT INTO `pt_shop` VALUES (1,666888,'百花店','佛山市禅城区汾江中路20号','0757-08323123','113.116159,23.044202',5,3,1,1421890019),(7,666888,'鸿运配送点','汾江中路12号a座B号铺','0757-21211122',NULL,0,1,1,2014),(10,666888,'天马店','佛山禅城区天马大厦','0757-22226666','113.114432,23.049800',3,2,1,2015),(11,666888,'南庄店','佛山市禅城区南庄镇吉利市场','0757-36668888','113.017224,22.987949',3,5,1,1421138355);
/*!40000 ALTER TABLE `pt_shop` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_siteconf`
--

DROP TABLE IF EXISTS `pt_siteconf`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_siteconf` (
  `pt_id` int(11) NOT NULL,
  `host` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
  `logo` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `index_title` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '首页标题',
  `sub_title` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '子页面标题',
  `state` int(11) DEFAULT NULL COMMENT '状态: 0:暂停  1：正常',
  `state_html` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`pt_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='合作商域名绑定';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_siteconf`
--

LOCK TABLES `pt_siteconf` WRITE;
/*!40000 ALTER TABLE `pt_siteconf` DISABLE KEYS */;
INSERT INTO `pt_siteconf` VALUES (666888,'','share/logo.gif','味道美 & 美味到_网上订餐系统',NULL,1,'技术故障，正在抢修');
/*!40000 ALTER TABLE `pt_siteconf` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sale_cart`
--

DROP TABLE IF EXISTS `sale_cart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sale_cart` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cart_key` varchar(32) DEFAULT NULL,
  `buyer_id` int(11) DEFAULT NULL,
  `order_no` varchar(45) DEFAULT NULL,
  `is_bought` tinyint(1) DEFAULT NULL COMMENT '是否已经购买',
  `create_time` int(11) DEFAULT NULL,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sale_cart`
--

LOCK TABLES `sale_cart` WRITE;
/*!40000 ALTER TABLE `sale_cart` DISABLE KEYS */;
INSERT INTO `sale_cart` VALUES (6,'a381dce294bb6f7c',1,'681884165',1,1423386532,1423386550),(7,'2f33ad05fd841e53',1,NULL,0,1423386773,1423386773);
/*!40000 ALTER TABLE `sale_cart` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sale_cart_item`
--

DROP TABLE IF EXISTS `sale_cart_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sale_cart_item` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cart_id` int(11) DEFAULT NULL,
  `goods_id` int(11) DEFAULT NULL COMMENT '商品快照编号',
  `num` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sale_cart_item`
--

LOCK TABLES `sale_cart_item` WRITE;
/*!40000 ALTER TABLE `sale_cart_item` DISABLE KEYS */;
INSERT INTO `sale_cart_item` VALUES (1,6,2,1),(2,6,3,1);
/*!40000 ALTER TABLE `sale_cart_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_ips`
--

DROP TABLE IF EXISTS `t_ips`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_ips` (
  `ip` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_ips`
--

LOCK TABLES `t_ips` WRITE;
/*!40000 ALTER TABLE `t_ips` DISABLE KEYS */;
INSERT INTO `t_ips` VALUES ('124.115.0.28'),('220.181.125.107'),('124.115.0.107'),('220.181.108.166'),('211.97.128.142');
/*!40000 ALTER TABLE `t_ips` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_members`
--

DROP TABLE IF EXISTS `t_members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_members` (
  `member` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_members`
--

LOCK TABLES `t_members` WRITE;
/*!40000 ALTER TABLE `t_members` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_members` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_usrcount`
--

DROP TABLE IF EXISTS `t_usrcount`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `t_usrcount` (
  `id` int(11) NOT NULL,
  `viewcount` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `member` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `guest` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updatetime` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_usrcount`
--

LOCK TABLES `t_usrcount` WRITE;
/*!40000 ALTER TABLE `t_usrcount` DISABLE KEYS */;
INSERT INTO `t_usrcount` VALUES (1,'5','0','5','2012-10-27');
/*!40000 ALTER TABLE `t_usrcount` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2015-02-08 17:14:44
