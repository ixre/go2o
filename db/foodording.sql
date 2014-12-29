CREATE DATABASE  IF NOT EXISTS `foodording` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
USE `foodording`;
-- MySQL dump 10.13  Distrib 5.6.17, for Linux (x86_64)
--
-- Host: localhost    Database: foodording
-- ------------------------------------------------------
-- Server version	5.1.73

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
-- Table structure for table `fd_itemprop`
--

DROP TABLE IF EXISTS `fd_itemprop`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fd_itemprop` (
  `id` int(11) NOT NULL,
  `description` text COLLATE utf8_unicode_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='食物属性,ID与fd_items.id关联';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fd_itemprop`
--

LOCK TABLES `fd_itemprop` WRITE;
/*!40000 ALTER TABLE `fd_itemprop` DISABLE KEYS */;
INSERT INTO `fd_itemprop` VALUES (153,''),(154,''),(155,''),(160,''),(161,''),(162,''),(163,''),(164,''),(165,''),(166,'');
/*!40000 ALTER TABLE `fd_itemprop` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `it_category`
--

DROP TABLE IF EXISTS `it_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `it_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pid` int(11) DEFAULT NULL COMMENT '父分类',
  `ptid` int(11) DEFAULT NULL COMMENT '商家ID(pattern ID);如果为空，则表示模式分类',
  `name` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `order_index` int(11) DEFAULT '0' COMMENT '序号',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `enabled` bit(1) DEFAULT NULL COMMENT '是否可用',
  `descript` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '描述',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=20 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='food categories';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `it_category`
--

LOCK TABLES `it_category` WRITE;
/*!40000 ALTER TABLE `it_category` DISABLE KEYS */;
INSERT INTO `it_category` VALUES (13,0,666888,'小炒',0,'2012-03-05 00:00:00','',''),(14,0,666888,'面食',0,'2012-03-05 00:00:00','',NULL),(15,0,666888,'套餐',0,'2012-03-05 00:00:00','',''),(16,0,666888,'油炸',0,'2012-03-05 00:00:00','',NULL),(17,0,666888,'海鲜',0,'2012-03-06 00:00:00','',NULL),(18,15,666888,'营养套餐',5,'2012-03-30 17:10:30','','');
/*!40000 ALTER TABLE `it_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `it_item`
--

DROP TABLE IF EXISTS `it_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `it_item` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cid` int(11) DEFAULT NULL COMMENT '分类',
  `name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `img` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `cost` decimal(5,2) DEFAULT '0.00' COMMENT ' 成本价',
  `price` decimal(5,2) DEFAULT '0.00' COMMENT '售价(市场价)',
  `sale_price` decimal(5,2) DEFAULT NULL COMMENT '实际销售价',
  `apply_subs` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '供应分店,用'',''隔开',
  `note` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '备注，如新菜色，特价优惠等',
  `description` varchar(500) COLLATE utf8_unicode_ci DEFAULT NULL,
  `state` int(11) DEFAULT '1',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=31 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='食物项';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `it_item`
--

LOCK TABLES `it_item` WRITE;
/*!40000 ALTER TABLE `it_item` DISABLE KEYS */;
INSERT INTO `it_item` VALUES (1,18,'韭黄炒蛋饭-2-2','666888/item_pic/20141022090923.png',15.00,20.00,18.00,'1','1','2',1,'2014-10-22 00:49:51','2014-10-26 05:43:39'),(2,15,'鱼香茄子饭',NULL,15.00,20.00,18.00,NULL,NULL,NULL,1,'2014-10-22 00:49:51','2014-06-12 00:00:00'),(3,15,'尖椒回锅肉',NULL,15.00,20.00,18.00,NULL,NULL,NULL,1,'2014-10-22 00:49:51','2014-06-12 00:00:00'),(4,15,'香菇焖鸡饭',NULL,15.00,20.00,18.00,NULL,NULL,NULL,1,'2014-10-22 00:49:51','2014-06-12 00:00:00'),(5,15,'野山椒炒肉',NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,'2014-10-22 00:49:51','2014-06-12 00:00:00'),(6,15,'极品红烧肉',NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,'2014-10-22 00:49:51','2014-06-12 00:00:00'),(7,15,'花生猪手饭',NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,'2014-10-22 00:49:51','2014-06-12 00:00:00'),(8,15,'红烧鱼腩饭',NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,'2014-10-22 00:49:51','2014-06-12 00:00:00'),(9,15,'香卤鸡腿饭',NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,'2014-10-22 00:49:51','2014-06-12 00:00:00'),(10,15,'酸甜排骨饭',NULL,1.00,0.00,1.00,NULL,NULL,NULL,1,'2014-10-22 00:49:51','2014-06-12 00:00:00'),(29,18,'营养套餐B','666888/item_pic/20141023090944.png',10.00,20.00,18.00,'1',NULL,NULL,1,'2014-10-22 01:53:03','2014-10-26 05:43:29'),(28,18,'营养套餐A','666888/item_pic/20141022090951.png',5.00,15.00,12.00,'1',NULL,NULL,1,'2014-10-22 01:51:23','2014-10-26 05:43:34');
/*!40000 ALTER TABLE `it_item` ENABLE KEYS */;
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
  `update_time` datetime DEFAULT NULL COMMENT '积分',
  PRIMARY KEY (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_account`
--

LOCK TABLES `mm_account` WRITE;
/*!40000 ALTER TABLE `mm_account` DISABLE KEYS */;
INSERT INTO `mm_account` VALUES (1,13990,11.42,183.92,2310.42,0.00,2346.90,'2014-12-22 03:20:58'),(2,0,2.50,NULL,2.50,0.00,0.00,'2013-04-01 22:12:59'),(28,180,2.79,1.44,20.79,0.00,36.00,'2014-12-19 08:57:24'),(29,0,4.14,NULL,4.14,0.00,36.00,'2012-09-24 09:52:09'),(30,0,2.70,NULL,2.70,0.00,27.00,'2012-09-24 09:52:09'),(31,0,0.00,NULL,0.00,0.00,0.00,'2012-10-01 03:03:30'),(32,0,0.00,NULL,0.00,0.00,0.00,'2012-12-19 03:54:43'),(33,0,0.00,NULL,0.00,0.00,0.00,'2012-12-21 22:32:45'),(34,0,0.00,NULL,0.00,0.00,0.00,'2013-03-07 04:10:00'),(35,0,0.00,NULL,0.00,0.00,0.00,'2013-03-07 04:12:31'),(36,0,0.00,NULL,0.00,0.00,0.00,'2013-03-07 04:19:58'),(37,0,0.00,NULL,0.00,0.00,0.00,'2013-03-07 04:20:54'),(38,0,1.00,NULL,1.00,0.00,10.00,'2013-03-11 03:40:36'),(39,0,2.00,NULL,2.00,0.00,20.00,'2013-03-11 03:50:30'),(40,0,3.80,NULL,3.80,0.00,38.00,'2013-04-01 22:12:59'),(41,0,0.00,NULL,0.00,0.00,0.00,'2013-03-31 21:16:12');
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
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_bank`
--

LOCK TABLES `mm_bank` WRITE;
/*!40000 ALTER TABLE `mm_bank` DISABLE KEYS */;
INSERT INTO `mm_bank` VALUES (1,'中国工商银行','513701198801105317','张三 ','上海分行漕溪路支行 ',0,'2012-06-14 04:43:20'),(28,'中国邮政储蓄银行','123486855651',NULL,NULL,1,'2012-09-16 04:25:12');
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
  `record_time` datetime DEFAULT NULL,
  `state` int(11) DEFAULT NULL COMMENT '状态(如：无效），默认为1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=83 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='进账日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_income_log`
--

LOCK TABLES `mm_income_log` WRITE;
/*!40000 ALTER TABLE `mm_income_log` DISABLE KEYS */;
INSERT INTO `mm_income_log` VALUES (8,NULL,28,'backcash',0.27,'来自订单:20120924-9603(商家:万绿园,会员:sumian)收入￥0.27元.','2012-09-24 09:52:09',NULL),(9,NULL,1,'backcash',7.40,'订单:20130308-2932返现￥7.4元','2013-03-08 05:26:20',NULL),(10,NULL,38,'backcash',1.00,'订单:20130307-5447返现￥1.0元','2013-03-11 03:40:36',NULL),(11,NULL,2,'backcash',0.20,'来自订单:20130307-5447(商家:蚁族,会员:yuan)收入￥0.2元.','2013-03-11 03:40:36',NULL),(12,NULL,39,'backcash',2.00,'订单:20130311-9708返现￥2.0元','2013-03-11 03:50:30',NULL),(13,NULL,2,'backcash',0.40,'来自订单:20130311-9708(商家:蚁族,会员:yuanyuan)收入￥0.4元.','2013-03-11 03:50:30',NULL),(14,NULL,40,'backcash',0.90,'订单:20130331-4922返现￥0.9元','2013-03-31 12:16:28',NULL),(15,NULL,2,'backcash',0.18,'来自订单:20130331-4922(商家:蚁族,会员:13924886758)收入￥0.18元.','2013-03-31 12:16:28',NULL),(16,NULL,40,'backcash',2.90,'订单:20130401-2503返现￥2.9元','2013-04-01 22:12:59',NULL),(17,NULL,2,'backcash',0.58,'来自订单:20130401-2503(商家:蚁族,会员:13924886758)收入￥0.58元.','2013-04-01 22:12:59',NULL),(18,NULL,1,'backcash',4.80,'订单:20130401-2503(商家:蚁族)返现￥4.80元','2014-10-26 04:06:38',1),(19,NULL,1,'backcash',1.20,'订单:20130401-2503(商家:蚁族,会员:刘铭)收入￥%!s(float32=1.2)元','2014-10-26 04:06:38',1),(20,NULL,1,'backcash',0.60,'订单:20130401-2503(商家:蚁族,会员:刘铭)收入￥%!s(float32=0.6)元','2014-10-26 04:06:38',1),(21,NULL,1,'backcash',4.80,'订单:20130401-7643(商家:蚁族)返现￥4.80元','2014-10-26 04:16:29',1),(22,NULL,1,'backcash',1.20,'订单:20130401-7643(商家:蚁族,会员:刘铭)收入￥1.20元','2014-10-26 04:16:29',1),(23,NULL,1,'backcash',0.60,'订单:20130401-7643(商家:蚁族,会员:刘铭)收入￥0.60元','2014-10-26 04:16:29',1),(24,NULL,1,'backcash',4.80,'订单:20130401-5283(商家:蚁族)返现￥4.80元','2014-10-26 04:22:56',1),(25,NULL,1,'backcash',1.20,'订单:20130401-5283(商家:蚁族,会员:刘铭)收入￥1.20元','2014-10-26 04:22:56',1),(26,NULL,1,'backcash',0.60,'订单:20130401-5283(商家:蚁族,会员:刘铭)收入￥0.60元','2014-10-26 04:22:56',1),(27,NULL,1,'backcash',4.80,'订单:20130401-8432(商家:蚁族)返现￥4.80元','2014-10-26 04:31:42',1),(28,NULL,1,'backcash',1.20,'订单:20130401-8432(商家:蚁族,会员:刘铭)收入￥1.20元','2014-10-26 04:31:42',1),(29,NULL,1,'backcash',0.60,'订单:20130401-8432(商家:蚁族,会员:刘铭)收入￥0.60元','2014-10-26 04:31:42',1),(30,NULL,1,'backcash',4.80,'订单:20130331-4998(商家:蚁族)返现￥4.80元','2014-10-26 04:31:50',1),(31,NULL,1,'backcash',1.20,'订单:20130331-4998(商家:蚁族,会员:刘铭)收入￥1.20元','2014-10-26 04:31:50',1),(32,NULL,1,'backcash',0.60,'订单:20130331-4998(商家:蚁族,会员:刘铭)收入￥0.60元','2014-10-26 04:31:50',1),(33,NULL,1,'backcash',4.80,'订单:20130331-2185(商家:蚁族)返现￥4.80元','2014-10-26 05:18:01',1),(34,NULL,1,'backcash',1.20,'订单:20130331-2185(商家:蚁族,会员:刘铭)收入￥1.20元','2014-10-26 05:18:01',1),(35,NULL,1,'backcash',0.60,'订单:20130331-2185(商家:蚁族,会员:刘铭)收入￥0.60元','2014-10-26 05:18:01',1),(36,NULL,1,'backcash',4.80,'订单:20130401-4602(商家:蚁族)返现￥4.80元','2014-10-26 05:18:08',1),(37,NULL,1,'backcash',1.20,'订单:20130401-4602(商家:蚁族,会员:刘铭)收入￥1.20元','2014-10-26 05:18:08',1),(38,NULL,1,'backcash',0.60,'订单:20130401-4602(商家:蚁族,会员:刘铭)收入￥0.60元','2014-10-26 05:18:08',1),(39,NULL,1,'backcash',4.80,'订单:20130331-4922(商家:蚁族)返现￥4.80元','2014-10-26 05:47:24',1),(40,NULL,1,'backcash',1.20,'订单:20130331-4922(商家:蚁族,会员:刘铭)收入￥1.20元','2014-10-26 05:47:24',1),(41,NULL,1,'backcash',0.60,'订单:20130331-4922(商家:蚁族,会员:刘铭)收入￥0.60元','2014-10-26 05:47:24',1),(42,NULL,1,'backcash',4.80,'订单:20130315-1432(商家:蚁族)返现￥4.80元','2014-10-26 06:50:38',1),(43,NULL,1,'backcash',1.20,'订单:20130315-1432(商家:蚁族,会员:刘铭)收入￥1.20元','2014-10-26 06:50:38',1),(44,NULL,1,'backcash',0.60,'订单:20130315-1432(商家:蚁族,会员:刘铭)收入￥0.60元','2014-10-26 06:50:38',1),(45,NULL,1,'backcash',4.80,'订单:20130311-9708(商家:蚁族)返现￥4.80元','2014-10-26 07:01:18',1),(46,NULL,1,'backcash',1.20,'订单:20130311-9708(商家:蚁族,会员:刘铭)收入￥1.20元','2014-10-26 07:01:18',1),(47,NULL,1,'backcash',0.60,'订单:20130311-9708(商家:蚁族,会员:刘铭)收入￥0.60元','2014-10-26 07:01:18',1),(48,NULL,1,'backcash',4.80,'订单:20130308-2932(商家:蚁族)返现￥4.80元','2014-10-26 07:01:27',1),(49,NULL,1,'backcash',1.20,'订单:20130308-2932(商家:蚁族,会员:刘铭)收入￥1.20元','2014-10-26 07:01:27',1),(54,129,1,'backcash',4.80,'订单:20130203-1190(商家:美味道)返现￥4.80元','2014-11-13 14:39:40',1),(55,127,1,'backcash',4.80,'订单:20130115-7948(商家:美味道)返现￥4.80元','2014-11-13 14:41:28',1),(58,153,1,'backcash',4.40,'订单:689182668(商家:美味道)返现￥4.40元','2014-11-13 15:07:52',1),(59,0,1,'backcash',1.10,'订单:689182668(商家:美味道,会员:刘铭)收入￥1.10元','2014-11-13 15:07:52',1),(60,0,1,'backcash',0.55,'订单:689182668(商家:美味道,会员:刘铭)收入￥0.55元','2014-11-13 15:07:52',1),(61,154,1,'backcash',1.44,'订单:685815333(商家:美味道)返现￥1.44元','2014-11-14 07:34:14',1),(62,0,1,'backcash',0.36,'订单:685815333(商家:美味道,会员:刘铭)收入￥0.36元','2014-11-14 07:34:14',1),(63,0,1,'backcash',0.18,'订单:685815333(商家:美味道,会员:刘铭)收入￥0.18元','2014-11-14 07:34:14',1),(64,152,1,'backcash',2.88,'订单:684452597(商家:美味道)返现￥2.88元','2014-11-14 07:37:28',1),(65,0,1,'backcash',0.72,'订单:684452597(商家:美味道,会员:刘铭)收入￥0.72元','2014-11-14 07:37:28',1),(66,0,1,'backcash',0.36,'订单:684452597(商家:美味道,会员:刘铭)收入￥0.36元','2014-11-14 07:37:28',1),(67,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元','2014-12-16 01:30:03',1),(68,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元','2014-12-16 01:30:11',1),(69,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元','2014-12-16 01:30:19',1),(70,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元','2014-12-16 01:34:58',1),(71,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元','2014-12-16 01:36:00',1),(72,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元','2014-12-16 01:40:34',1),(73,181,1,'backcash',10.08,'订单:688664284(商家:美味道)返现￥10.08元','2014-12-16 01:46:48',1),(74,180,1,'backcash',4.32,'订单:687153988(商家:美味道)返现￥4.32元','2014-12-16 01:47:01',1),(75,179,1,'backcash',2.88,'订单:681645404(商家:美味道)返现￥2.88元','2014-12-16 01:47:08',1),(76,177,1,'backcash',2.88,'订单:687025732(商家:美味道)返现￥2.88元','2014-12-16 01:47:16',1),(77,183,28,'backcash',1.44,'订单:689469831(商家:美味道)返现￥1.44元','2014-12-19 08:57:24',1),(78,184,1,'backcash',1.44,'订单:689205996(商家:美味道)返现￥1.44元','2014-12-22 02:56:09',1),(79,182,1,'backcash',2.88,'订单:684148199(商家:美味道)返现￥2.88元','2014-12-22 02:58:46',1),(80,178,1,'backcash',2.88,'订单:683878066(商家:美味道)返现￥2.88元','2014-12-22 03:02:34',1),(81,176,1,'backcash',2.88,'订单:681710378(商家:美味道)返现￥2.88元','2014-12-22 03:16:22',1),(82,175,1,'backcash',2.88,'订单:681009722(商家:美味道)返现￥2.88元','2014-12-22 03:20:58',1);
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
  `record_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_integral_log`
--

LOCK TABLES `mm_integral_log` WRITE;
/*!40000 ALTER TABLE `mm_integral_log` DISABLE KEYS */;
INSERT INTO `mm_integral_log` VALUES (1,666888,1,3,600,'订单返积分600个','2014-11-13 14:41:28'),(2,666888,1,3,550,'订单返积分550个','2014-11-13 15:07:52'),(3,666888,1,3,180,'订单返积分180个','2014-11-14 07:34:14'),(4,666888,1,3,360,'订单返积分360个','2014-11-14 07:37:28'),(5,666888,1,3,1260,'订单返积分1260个','2014-12-16 01:30:03'),(6,666888,1,3,1260,'订单返积分1260个','2014-12-16 01:30:11'),(7,666888,1,3,1260,'订单返积分1260个','2014-12-16 01:30:19'),(8,666888,1,3,1260,'订单返积分1260个','2014-12-16 01:34:58'),(9,666888,1,3,1260,'订单返积分1260个','2014-12-16 01:36:01'),(10,666888,1,3,1260,'订单返积分1260个','2014-12-16 01:40:34'),(11,666888,1,3,1260,'订单返积分1260个','2014-12-16 01:46:48'),(12,666888,1,3,540,'订单返积分540个','2014-12-16 01:47:01'),(13,666888,1,3,360,'订单返积分360个','2014-12-16 01:47:08'),(14,666888,1,3,360,'订单返积分360个','2014-12-16 01:47:16'),(15,666888,28,3,180,'订单返积分180个','2014-12-19 08:57:24'),(16,666888,1,3,180,'订单返积分180个','2014-12-22 02:56:09'),(17,666888,1,3,360,'订单返积分360个','2014-12-22 02:58:46'),(18,666888,1,3,360,'订单返积分360个','2014-12-22 03:02:34'),(19,666888,1,3,360,'订单返积分360个','2014-12-22 03:16:22'),(20,666888,1,3,360,'订单返积分360个','2014-12-22 03:20:58');
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
  `reg_time` datetime DEFAULT NULL,
  `reg_ip` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `state` int(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=29 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mm_member`
--

LOCK TABLES `mm_member` WRITE;
/*!40000 ALTER TABLE `mm_member` DISABLE KEYS */;
INSERT INTO `mm_member` VALUES (1,'newmin','768dd5c54e40bcd412f7277cbe77423e','刘铭',414,3,1,'','2012-09-05','18616999822','上海市徐汇区江安路15号1','18867761','q121276@26.com','2012-05-13 10:17:20','112.65.35.191','2014-12-22 05:27:06',1),(2,'lin','22c5b45021a538af170f330bf6d6e46c','林德意',NULL,1,0,'','1970.05。29','13924886758','佛山禅城人民路99号',NULL,NULL,'2012-05-14 21:41:04','61.145.69.27','2012-05-10 00:00:00',1),(3,'hshx','c99425b0a379ac6621adbc0ce4170af5','黄升鑫',NULL,1,0,'','1979-12-01','15602817110','广东省佛山市禅城区人民路鹤园路81号','809987822','809987822@qq.com','2012-10-01 03:03:30','27.36.72.124','0001-01-01 00:00:00',1),(4,'lindeyi','0f167fd9c5f48d81820b544e312e8592','林德意',NULL,1,0,'','1970-05-29','13924886758','佛山市禅城区人民路99号','569101942','lindeyi158@yahoo.com.cn','2012-09-16 03:58:06','183.27.197.170','0001-01-01 00:00:00',1),(5,'sumian','4397d538520a9a645aa456e60744c1e0','',NULL,1,0,'','','13924886758','福建省海天大夏二栋201','','','2012-09-24 09:34:19','14.157.18.39','0001-01-01 00:00:00',1),(6,'linsu','464bf5d58f4e8818671d525cf1530459','',NULL,1,0,'','','13924886758','广州市黄埔去电子大夏3楼301室','','','2012-09-24 09:46:11','14.157.18.39','0001-01-01 00:00:00',1),(7,'sonven','b99831c5e69ac900fdfcfd4c7d0bf89e','',NULL,1,0,'','','18616999822','上海市浦东新区浦电路123弄','','','2012-12-19 03:54:43','183.250.3.128','0001-01-01 00:00:00',1),(8,'yangbo','1df2ee43288507769a15da6cb1cf0dba','',NULL,1,0,'','','18616888888','佛山市禅城区汾江中路20号','','','2012-12-21 22:32:45','183.250.3.128','0001-01-01 00:00:00',1),(9,'xiaoyuan','db5896d7e1951418a6fe0de4ea86b45b','',NULL,1,0,'','','18616999822','上海市徐汇区飞掉国际大厦','','','2013-03-07 04:10:00','218.85.143.146','0001-01-01 00:00:00',1),(10,'liuming','06f267d8e85c3478e00a8b9d2bae5df4','',NULL,1,0,'','','18616999822','上海市徐汇区飞掉国际大厦','','','2013-03-07 04:12:31','218.85.143.146','0001-01-01 00:00:00',1),(11,'xiaoyuanyaun','dda4f29b5f09313383fcfc02c0ce2753','',NULL,1,0,'','','18616999822','上海市徐汇区飞掉国际大厦','','','2013-03-07 04:19:58','218.85.143.146','0001-01-01 00:00:00',1),(12,'liuxiaoyuan','b201ab396944544760bd6b19d356cc8f','',NULL,1,0,'','','18616999822','上海市徐汇区飞掉国际大厦','','','2013-03-07 04:20:54','218.85.143.146','0001-01-01 00:00:00',1),(13,'yuan','86ad59c09a3f4fe980b67b9dedea7329','',NULL,1,0,'','','13728501775','佛山市禅城区张槎四路东大街3号5楼','','','2013-03-07 04:53:47','183.27.195.23','0001-01-01 00:00:00',1),(14,'yuanyuan','40c39d9211d17d1e6732059427c9ee76','',NULL,1,0,'','','13728501775','佛山禅城区张槎四路岗头东大街3号5楼','','','2013-03-11 03:44:00','183.27.199.29','0001-01-01 00:00:00',1),(15,'13924886758','c545e63671045c96669b814901bf0d37','',NULL,1,0,'','','13924886758','佛山市禅城区人民路99号','','','2013-03-31 12:07:47','183.28.79.121','0001-01-01 00:00:00',1),(16,'13728501775','9f86434fc3c081f7548e633c7ccdc5d2','',NULL,1,0,'','','13728501775','佛山市禅城区张槎四路（东海明珠后方）岗头东大街3号5楼','','','2013-03-31 21:16:12','183.27.46.207','0001-01-01 00:00:00',1),(25,'sa','123','刘铭',NULL,1,2,NULL,'1970-11-20','18616999822',NULL,NULL,NULL,'2014-10-22 12:50:14','127.0.0.1','2014-10-22 12:50:14',1),(26,'test','a50f4d2b5d08eca0ff83448fc346dbd6','测试员',NULL,1,1,NULL,'1988-11-09','18616999822',NULL,NULL,NULL,'2014-10-28 10:04:32','127.0.0.1','2014-10-28 10:04:32',1),(27,'test001','4dca21a567d5ae25316f9e8d37d8df1b','刘大炮',NULL,1,0,'share/noavatar.gif','1970-01-01','18616999822',NULL,NULL,NULL,'2014-11-27 02:00:37','127.0.0.1','2014-11-27 02:00:50',1),(28,'newmin123','9315871be89146db634ef0d0e5e181f9','刘大也',0,1,0,'share/noavatar.gif','1970-01-01','18616999822',NULL,NULL,NULL,'2014-12-19 08:55:55','127.0.0.1','2014-12-19 08:56:13',1);
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
  `begin_time` datetime DEFAULT NULL,
  `over_time` datetime DEFAULT NULL,
  `allow_enable` tinyint(1) DEFAULT NULL COMMENT '是否允许使用',
  `need_bind` tinyint(1) DEFAULT NULL COMMENT '是否需要绑定',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL COMMENT '共计数量',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pm_coupon`
--

LOCK TABLES `pm_coupon` WRITE;
/*!40000 ALTER TABLE `pm_coupon` DISABLE KEYS */;
INSERT INTO `pm_coupon` VALUES (1,666888,'30off',NULL,0,10,1,70,30,0,0,'2014-12-01 00:00:00','2014-12-12 00:00:00',1,1,'2014-12-04 02:45:15','2014-12-09 06:15:43'),(2,666888,'1WEEK',NULL,10,10,10,100,0,0,100,'2015-01-01 00:00:00','2015-01-01 00:00:00',0,0,'2014-12-04 11:13:42','2014-12-10 07:23:11'),(3,666888,'10off',NULL,10,10,4,100,0,0,20,'2015-01-01 00:00:00','2015-01-01 00:00:00',0,0,'2014-12-04 23:38:53','2014-12-10 07:21:46'),(4,666888,'dsss',NULL,10,10,0,100,10,0,0,'2015-01-02 00:00:00','2015-01-02 00:00:00',0,0,'2014-12-04 23:40:28','2014-12-10 07:22:10'),(5,666888,'95off','95折全场通用',12,12,0,95,0,0,0,'2014-12-05 00:00:00','2014-12-26 00:00:00',1,0,'2014-12-04 23:51:48','2014-12-10 08:12:26');
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
  `bind_time` int(11) DEFAULT NULL COMMENT '绑定时间',
  `is_used` tinyint(1) DEFAULT '0' COMMENT '是否使用',
  `use_time` int(11) DEFAULT NULL COMMENT '使用时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pm_coupon_bind`
--

LOCK TABLES `pm_coupon_bind` WRITE;
/*!40000 ALTER TABLE `pm_coupon_bind` DISABLE KEYS */;
INSERT INTO `pm_coupon_bind` VALUES (1,2,1,2014,0,2014),(2,1,1,2014,1,1418227815);
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
  `take_time` int(11) DEFAULT NULL COMMENT '占用时间',
  `extra_time` int(11) DEFAULT NULL COMMENT '释放时间,超过该时间，优惠券释放',
  `is_apply` tinyint(1) DEFAULT NULL COMMENT '是否生效,1表示有效',
  `apply_time` int(11) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pm_coupon_take`
--

LOCK TABLES `pm_coupon_take` WRITE;
/*!40000 ALTER TABLE `pm_coupon_take` DISABLE KEYS */;
INSERT INTO `pm_coupon_take` VALUES (3,1,5,1418225350,1418239750,1,1418225361),(2,1,5,1418224519,1418238919,1,1418225318),(4,1,5,1418225381,1418239781,1,1418225391),(5,1,5,1418225657,1418240057,0,1418222057),(6,1,5,1418260980,1418275380,1,1418260983);
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
  `items` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '商品信息,\n17*1|18*2|50',
  `items_info` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `total_fee` decimal(10,2) DEFAULT NULL COMMENT '订单总额',
  `fee` decimal(10,2) DEFAULT NULL COMMENT '订单实际金额',
  `discount_fee` decimal(10,2) DEFAULT NULL COMMENT '优惠金额',
  `coupon_fee` decimal(10,2) DEFAULT NULL COMMENT '优惠券优惠金额',
  `pay_fee` decimal(10,2) DEFAULT '0.00' COMMENT '支付金额',
  `pay_method` int(11) DEFAULT NULL COMMENT '1:餐到付款 2:网上支付  ',
  `is_payed` int(11) DEFAULT NULL COMMENT '是否支付(0:未支付 ，1：已支付)',
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `deliver_name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `deliver_phone` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `deliver_address` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `deliver_time` datetime DEFAULT NULL COMMENT '送餐时间',
  `status` int(11) DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=185 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_order`
--

LOCK TABLES `pt_order` WRITE;
/*!40000 ALTER TABLE `pt_order` DISABLE KEYS */;
INSERT INTO `pt_order` VALUES (1,'20121103-1971',1,666888,'0','%u5C16%u6912%u56DE%u9505%u8089*164*10*1|10','尖椒回锅肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-11-03 08:57:24',0,'2012-11-03 08:57:25','2012-11-03 08:57:25'),(2,'20120924-9603',1,666888,'1','%u97ED%u9EC4%u7092%u86CB%u996D*162*9*1|%u9999%u83C7%u7116%u9E21%u996D*161*9*1|%u9178%u751C%u6392%u9A','韭黄炒蛋饭() * 9\n香菇焖鸡饭() * 9\n酸甜排骨饭() * 9\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-09-24 09:47:08',0,'2012-09-24 09:47:09','2012-09-24 09:47:09'),(3,'20121027-8653',1,666888,'0','%u9178%u751C%u6392%u9AA8%u996D*160*9*1|%u6781%u54C1%u7EA2%u70E7%u8089*166*9*2|27','酸甜排骨饭() * 9\n极品红烧肉() * 9\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'xiexie',NULL,NULL,NULL,'2012-10-27 18:45:50',0,'2012-10-27 18:45:51','2012-10-27 18:45:51'),(4,'20120924-3892',1,666888,'1','%u9178%u751C%u6392%u9AA8%u996D*160*9*1|%u9999%u83C7%u7116%u9E21%u996D*161*9*1|%u97ED%u9EC4%u7092%u86','酸甜排骨饭() * 9\n香菇焖鸡饭() * 9\n韭黄炒蛋饭() * 9\n极品红烧肉() * 9\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-09-24 09:36:00',0,'2012-09-24 09:36:03','2012-09-24 09:36:03'),(5,'20120627-2436',15,666888,'1','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 04:03:55',3,'2012-06-27 04:03:00','2014-10-26 04:04:15'),(6,'20120923-4551',1,666888,'1','%u9178%u751C%u6392%u9AA8%u996D*160*9*1|%u9999%u83C7%u7116%u9E21%u996D*161*9*1|18','酸甜排骨饭() * 9\n香菇焖鸡饭() * 9\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-09-23 20:53:49',0,'2012-09-23 20:53:50','2012-09-23 20:53:50'),(7,'20120922-8308',1,666888,'0','%u9178%u751C%u6392%u9AA8%u996D*160*10*1|%u9999%u83C7%u7116%u9E21%u996D*161*10*1|%u97ED%u9EC4%u7092%u','酸甜排骨饭() * 10\n香菇焖鸡饭() * 10\n韭黄炒蛋饭() * 10\n极品红烧肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-09-22 05:59:25',0,'2012-09-22 05:59:26','2012-09-22 05:59:26'),(8,'20120627-1100',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:04:06',0,'2012-06-27 04:08:24','2012-06-27 04:08:24'),(9,'20120627-7755',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:04:06',0,'2012-06-27 04:08:24','2012-06-27 04:08:24'),(10,'20120627-6581',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:04:07',-1,'2012-06-27 04:08:24','2012-06-27 04:08:24'),(11,'20120627-8125',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:05:11',-1,'2012-06-27 04:11:25','2012-06-27 04:11:25'),(12,'20120627-4651',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:05:11',0,'2012-06-27 04:11:24','2012-06-27 04:11:24'),(13,'20120627-2526',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:05:11',0,'2012-06-27 04:11:28','2012-06-27 04:11:28'),(14,'20120627-7442',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:05:11',-1,'2012-06-27 04:11:25','2012-06-27 04:11:25'),(15,'20120627-7080',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:06:18',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(16,'20120627-8196',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:06:18',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(17,'20120627-5649',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:06:18',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(18,'20120627-6550',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:06:18',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(19,'20120627-6220',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:07:19',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(20,'20120627-1654',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:07:21',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(21,'20120627-8171',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:07:21',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(22,'20120627-2550',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:07:21',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(23,'20120627-7952',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:08:19',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(24,'20120627-6759',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:08:23',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(25,'20120627-8506',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:08:23',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(26,'20120627-5911',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:08:24',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(27,'20120627-7876',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:09:19',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(28,'20120627-8346',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:09:23',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(29,'20120627-6117',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:09:23',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(30,'20120627-1385',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:09:24',0,'2012-06-27 04:15:30','2012-06-27 04:15:30'),(31,'20120627-4827',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:10:19',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(32,'20120627-7772',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:10:23',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(33,'20120627-1022',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:10:23',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(34,'20120627-1881',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:10:25',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(35,'20120627-9740',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:11:19',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(36,'20120627-3168',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:11:24',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(37,'20120627-4255',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:11:28',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(38,'20120627-4677',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:11:28',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(39,'20120627-4637',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:12:20',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(40,'20120627-7191',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:12:24',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(41,'20120627-1907',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:12:28',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(42,'20120627-1813',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:12:28',0,'2012-06-27 04:15:31','2012-06-27 04:15:31'),(43,'20120627-5496',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:13:20',0,'2012-06-27 04:15:36','2012-06-27 04:15:36'),(44,'20120627-2231',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:13:24',0,'2012-06-27 04:15:40','2012-06-27 04:15:40'),(45,'20120627-3244',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:13:29',0,'2012-06-27 04:15:40','2012-06-27 04:15:40'),(46,'20120627-1174',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:13:29',0,'2012-06-27 04:15:40','2012-06-27 04:15:40'),(47,'20120627-1436',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:14:20',0,'2012-06-27 04:15:40','2012-06-27 04:15:40'),(48,'20120627-7910',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:14:24',0,'2012-06-27 04:15:40','2012-06-27 04:15:40'),(49,'20120627-3308',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:14:29',0,'2012-06-27 04:15:40','2012-06-27 04:15:40'),(50,'20120627-5287',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:14:29',0,'2012-06-27 04:15:40','2012-06-27 04:15:40'),(51,'20120627-5115',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:15:20',0,'2012-06-27 04:15:40','2012-06-27 04:15:40'),(52,'20120627-4901',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:15:25',0,'2012-06-27 04:15:40','2012-06-27 04:15:40'),(53,'20120627-7672',15,666888,'0','%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C46%u9F13*27*12*2|%u7CD6%u918B%u6392%u9A','麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:15:29',0,'2012-06-27 04:15:41','2012-06-27 04:15:41'),(54,'20120627-6181',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:27:07',0,'2012-06-27 04:27:08','2012-06-27 04:27:08'),(55,'20120627-1763',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:27:08',0,'2012-06-27 04:27:08','2012-06-27 04:27:08'),(56,'20120627-9625',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:27:09',0,'2012-06-27 04:27:14','2012-06-27 04:27:14'),(57,'20120627-4537',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:27:14',0,'2012-06-27 04:28:14','2012-06-27 04:28:14'),(58,'20120627-8357',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:27:14',0,'2012-06-27 04:28:14','2012-06-27 04:28:14'),(59,'20120627-1266',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:27:14',0,'2012-06-27 04:28:14','2012-06-27 04:28:14'),(60,'20120627-8738',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:27:14',0,'2012-06-27 04:28:14','2012-06-27 04:28:14'),(61,'20120627-9911',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:28:17',0,'2012-06-27 04:30:18','2012-06-27 04:30:18'),(62,'20120627-2810',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:28:17',0,'2012-06-27 04:30:22','2012-06-27 04:30:22'),(63,'20120627-4900',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:28:17',0,'2012-06-27 04:30:22','2012-06-27 04:30:22'),(64,'20120627-1501',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:28:22',0,'2012-06-27 04:30:22','2012-06-27 04:30:22'),(65,'20120627-7929',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:29:18',0,'2012-06-27 04:32:22','2012-06-27 04:32:22'),(66,'20120627-5618',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:29:18',0,'2012-06-27 04:32:22','2012-06-27 04:32:22'),(67,'20120627-6252',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:29:22',0,'2012-06-27 04:32:24','2012-06-27 04:32:24'),(68,'20120627-9793',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:29:18',0,'2012-06-27 04:32:22','2012-06-27 04:32:22'),(69,'20120627-5242',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:30:22',0,'2012-06-27 04:34:29','2012-06-27 04:34:29'),(70,'20120627-2459',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:30:22',0,'2012-06-27 04:34:28','2012-06-27 04:34:28'),(71,'20120627-6815',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:30:22',0,'2012-06-27 04:34:28','2012-06-27 04:34:28'),(72,'20120627-8200',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:30:23',0,'2012-06-27 04:34:29','2012-06-27 04:34:29'),(73,'20120627-3302',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:31:24',0,'2012-06-27 04:35:32','2012-06-27 04:35:32'),(74,'20120627-9453',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:31:22',0,'2012-06-27 04:35:31','2012-06-27 04:35:31'),(75,'20120627-2646',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:32:24',0,'2012-06-27 04:35:33','2012-06-27 04:35:33'),(76,'20120627-6516',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:31:22',0,'2012-06-27 04:35:31','2012-06-27 04:35:31'),(77,'20120627-8374',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:31:22',0,'2012-06-27 04:35:31','2012-06-27 04:35:31'),(78,'20120627-5625',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:32:22',0,'2012-06-27 04:35:32','2012-06-27 04:35:32'),(79,'20120627-1217',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:32:24',0,'2012-06-27 04:35:33','2012-06-27 04:35:33'),(80,'20120627-6054',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:32:29',0,'2012-06-27 04:35:33','2012-06-27 04:35:33'),(81,'20120627-3177',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:33:27',0,'2012-06-27 04:35:33','2012-06-27 04:35:33'),(82,'20120627-7430',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:33:27',0,'2012-06-27 04:35:33','2012-06-27 04:35:33'),(83,'20120627-8355',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:33:29',0,'2012-06-27 04:35:33','2012-06-27 04:35:33'),(84,'20120627-8173',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:33:32',0,'2012-06-27 04:35:33','2012-06-27 04:35:33'),(85,'20120627-3528',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:34:30',0,'2012-06-27 04:36:46','2012-06-27 04:36:46'),(86,'20120627-7629',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:34:30',0,'2012-06-27 04:36:46','2012-06-27 04:36:46'),(87,'20120627-9910',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:34:30',0,'2012-06-27 04:36:46','2012-06-27 04:36:46'),(88,'20120627-5774',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:34:32',0,'2012-06-27 04:36:46','2012-06-27 04:36:46'),(89,'20120627-2158',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:35:33',0,'2012-06-27 04:36:46','2012-06-27 04:36:46'),(90,'20120627-6967',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:35:33',0,'2012-06-27 04:36:46','2012-06-27 04:36:46'),(91,'20120627-8772',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:35:34',0,'2012-06-27 04:36:46','2012-06-27 04:36:46'),(92,'20120627-9923',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:35:34',0,'2012-06-27 04:36:46','2012-06-27 04:36:46'),(93,'20120627-2908',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:36:41',0,'2012-06-27 04:36:47','2012-06-27 04:36:47'),(94,'20120627-1039',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:36:41',0,'2012-06-27 04:36:47','2012-06-27 04:36:47'),(95,'20120627-9722',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:36:41',0,'2012-06-27 04:36:47','2012-06-27 04:36:47'),(96,'20120627-2130',9,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:37:38',0,'2012-06-27 04:37:38','2012-06-27 04:37:38'),(97,'20120627-3465',9,666888,'1','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-06-27 04:37:40',0,'2012-06-27 04:37:40','2012-06-27 04:37:40'),(98,'20120702-9893',15,666888,'1','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-02 23:11:39',0,'2012-07-02 23:11:40','2012-07-02 23:11:40'),(99,'20120702-7680',15,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-02 23:11:39',0,'2012-07-02 23:11:40','2012-07-02 23:11:40'),(100,'20120702-9943',15,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-02 23:18:37',0,'2012-07-02 23:18:38','2012-07-02 23:18:38'),(101,'20120702-9760',15,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-02 23:18:37',0,'2012-07-02 23:18:38','2012-07-02 23:18:38'),(102,'20120702-3696',1,666888,'1','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-02 23:46:46',0,'2012-07-02 23:46:47','2012-07-02 23:46:47'),(103,'20120703-4504',8,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-03 02:56:58',0,'2012-07-03 02:57:02','2012-07-03 02:57:02'),(104,'20120703-8136',8,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-03 02:57:02',0,'2012-07-03 02:57:02','2012-07-03 02:57:02'),(105,'20120703-3632',8,666888,'1','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-03 02:57:08',0,'2012-07-03 02:57:08','2012-07-03 02:57:08'),(106,'20120703-4933',8,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-03 02:57:08',0,'2012-07-03 02:57:08','2012-07-03 02:57:08'),(107,'20120703-2793',8,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-03 02:57:08',0,'2012-07-03 02:57:08','2012-07-03 02:57:08'),(108,'20120703-5358',1,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|20','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-03 03:30:04',0,'2012-07-03 03:30:04','2012-07-03 03:30:04'),(109,'20120703-1458',1,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|20','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-03 03:30:04',0,'2012-07-03 03:30:04','2012-07-03 03:30:04'),(110,'20120703-9432',19,666888,'1','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-03 22:09:08',0,'2012-07-03 22:09:08','2012-07-03 22:09:08'),(111,'20120703-7117',19,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-03 22:09:08',0,'2012-07-03 22:09:08','2012-07-03 22:09:08'),(112,'20120707-2673',1,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u6CB9%u5356%u83','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-07 06:22:58',0,'2012-07-07 06:23:02','2012-07-07 06:23:02'),(113,'20120707-7558',22,666888,'1','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u6CB9%u5356%u83','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-07 11:24:33',0,'2012-07-07 11:24:34','2012-07-07 11:24:34'),(114,'20120707-7906',21,666888,'1','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9178%u8C46%u89D2%u7092%u8089%u7C73*17*10*1|%u9EBB%u5A46%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n酸豆角炒肉米(酸爽美味) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-07 11:25:20',0,'2012-07-07 11:25:20','2012-07-07 11:25:20'),(115,'20120707-8395',16,666888,'0','%u8425%u517B%u5957%u9910 *144*10*1|%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|20','营养套餐 () * 10\n木耳炒黄瓜(优惠仅限3日) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-07 11:45:52',0,'2012-07-07 11:45:52','2012-07-07 11:45:52'),(116,'20120707-6956',23,666888,'0','%u8425%u517B%u5957%u9910 *144*10*1|%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|20','营养套餐 () * 10\n木耳炒黄瓜(优惠仅限3日) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-07 11:56:21',0,'2012-07-07 11:56:22','2012-07-07 11:56:22'),(117,'20120707-4722',23,666888,'0','%u8FC7%u6865%u7C73%u7EBF*19*10*1|10','过桥米线() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-07 11:56:49',0,'2012-07-07 11:56:49','2012-07-07 11:56:49'),(118,'20120707-9913',16,666888,'0','%u9178%u8FA3%u7C89%u4E1D*18*10*1|10','酸辣粉丝() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-07 11:59:52',0,'2012-07-07 11:59:53','2012-07-07 11:59:53'),(119,'20120707-1485',16,666888,'0','%u9178%u8FA3%u7C89%u4E1D*18*10*1|10','酸辣粉丝() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-07 12:02:31',0,'2012-07-07 12:02:31','2012-07-07 12:02:31'),(120,'20120707-8310',23,666888,'0','%u9178%u8FA3%u7C89%u4E1D*18*10*1|%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|20','酸辣粉丝() * 10\n木耳炒黄瓜(优惠仅限3日) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-07 12:04:58',0,'2012-07-07 12:04:59','2012-07-07 12:04:59'),(121,'20120707-9986',23,666888,'0','%u8425%u517B%u5957%u9910 *144*10*1|10','营养套餐 () * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-07 12:05:20',0,'2012-07-07 12:05:20','2012-07-07 12:05:20'),(122,'20120709-2715',1,666888,'0','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9EBB%u5A46%u8C46%u8150*26*10*1|%u756A%u8304%u7092%u86CB*29*','木耳炒黄瓜(优惠仅限3日) * 10\n麻婆豆腐() * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n营养套餐 () * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-09 10:24:23',0,'2012-07-09 10:24:24','2012-07-09 10:24:24'),(123,'20120710-4124',1,666888,'0','%u7CD6%u918B%u6392%u9AA8*28*10*1|10','糖醋排骨(本店招牌特色) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-10 01:37:16',0,'2012-07-10 01:37:17','2012-07-10 01:37:17'),(124,'20120710-7520',1,666888,'0','%u5546%u52A1%u8425%u517B%u5957%u9910*145*168*1|168','商务营养套餐(：\"优惠仅限3日) * 168\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-10 01:57:32',0,'2012-07-10 01:57:32','2012-07-10 01:57:32'),(125,'20120710-9269',1,666888,'1','%u6728%u8033%u7092%u9EC4%u74DC*16*10*1|%u9EBB%u5A46%u8C46%u8150*26*10*1|%u6CB9%u5356%u83DC%u7092%u8C','木耳炒黄瓜(优惠仅限3日) * 10\n麻婆豆腐() * 10\n油卖菜炒豆鼓() * 12\n糖醋排骨(本店招牌特色) * 10\n番茄炒蛋() * 6\n流氓图(麻辣味) * 9.6\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-07-10 04:40:10',0,'2012-07-10 04:40:11','2012-07-10 04:40:11'),(126,'20120810-9477',1,666888,'4','%u9F13%u6C41%u6392%u9AA8%u996D*153*10.4*6|%u897F%u5170%u82B1%u7092%u8089%u7247 *155*10*4|%u8377%u517','鼓汁排骨饭((仅限今日8折优惠)) * 10.4\n西兰花炒肉片 () * 10\n荷兰豆炒肉片() * 10\n木耳炒黄瓜(优惠仅限3日) * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2012-08-10 04:30:56',0,'2012-08-10 04:30:57','2012-08-10 04:30:57'),(127,'20130115-7948',1,666888,'1','%u91CE%u5C71%u6912%u7092%u8089*165*10*1|10','野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,0,'',NULL,NULL,NULL,'2014-11-13 14:41:27',3,'2013-01-15 22:16:46','2014-11-13 14:41:28'),(128,'20130115-3856',1,666888,'0','%u91CE%u5C71%u6912%u7092%u8089*165*10*1|10','野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2013-01-15 22:16:39',0,'2013-01-15 22:16:46','2013-01-15 22:16:46'),(129,'20130203-1190',1,666888,'2','%u91CE%u5C71%u6912%u7092%u8089*165*10*1|%u5C16%u6912%u56DE%u9505%u8089*164*10*1|20','野山椒炒肉() * 10\n尖椒回锅肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,0,'',NULL,NULL,NULL,'2014-11-13 14:39:39',3,'2013-02-03 11:14:00','2014-11-13 14:39:40'),(130,'20130307-5447',1,666888,'1','%u5C16%u6912%u56DE%u9505%u8089*164*10*1|10','尖椒回锅肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'测试',NULL,NULL,NULL,'2014-10-29 03:04:34',3,'2013-03-07 04:55:15','2014-10-29 03:04:36'),(131,'20130308-2932',1,666888,'1','%u9178%u751C%u6392%u9AA8%u996D*160*9*2|%u6781%u54C1%u7EA2%u70E7%u8089*166*9*3|%u9999%u83C7%u7116%u9E','酸甜排骨饭() * 9\n极品红烧肉() * 9\n香菇焖鸡饭() * 9\n花生猪手饭() * 10\n鱼香茄子饭() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'要辣',NULL,NULL,NULL,'2014-10-26 07:01:25',3,'2013-03-08 05:24:36','2014-10-26 07:01:27'),(132,'20130311-9708',1,666888,'1','%u5C16%u6912%u56DE%u9505%u8089*164*10*1|%u91CE%u5C71%u6912%u7092%u8089*165*10*1|20','尖椒回锅肉() * 10\n野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 07:01:18',3,'2013-03-11 03:45:24','2014-10-26 07:01:18'),(133,'20130315-1432',1,666888,'7','%u5C16%u6912%u56DE%u9505%u8089*164*10*1|10','尖椒回锅肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 06:50:35',3,'2013-03-15 05:12:06','2014-10-26 06:50:38'),(134,'20130331-4922',1,666888,'1','%u5C16%u6912%u56DE%u9505%u8089*164*9*1|9','尖椒回锅肉() * 9\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'要辣一点的',NULL,NULL,NULL,'2014-10-26 05:47:23',3,'2013-03-31 12:09:42','2014-10-26 05:47:24'),(135,'20130331-4998',1,666888,'1','%u5C16%u6912%u56DE%u9505%u8089*164*9*1|%u91CE%u5C71%u6912%u7092%u8089*165*10*1|19','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 04:31:49',3,'2013-03-31 23:28:44','2014-10-26 04:31:50'),(136,'20130331-2185',1,666888,'1','%u5C16%u6912%u56DE%u9505%u8089*164*9*1|%u91CE%u5C71%u6912%u7092%u8089*165*10*1|19','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 04:43:41',3,'2013-03-31 23:28:44','2014-10-26 05:18:01'),(137,'20130401-8432',1,666888,'1','%u5C16%u6912%u56DE%u9505%u8089*164*9*1|%u91CE%u5C71%u6912%u7092%u8089*165*10*1|19','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 04:31:40',3,'2013-04-01 11:19:23','2014-10-26 04:31:42'),(138,'20130401-5283',1,666888,'1','%u5C16%u6912%u56DE%u9505%u8089*164*9*1|%u91CE%u5C71%u6912%u7092%u8089*165*10*1|19','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 04:22:55',3,'2013-04-01 11:19:33','2014-10-26 04:22:56'),(139,'20130401-4602',1,666888,'1','%u5C16%u6912%u56DE%u9505%u8089*164*9*1|%u91CE%u5C71%u6912%u7092%u8089*165*10*1|19','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 04:31:14',3,'2013-04-01 11:19:33','2014-10-26 05:18:08'),(140,'20130401-7643',1,666888,'1','%u91CE%u5C71%u6912%u7092%u8089*165*10*1|%u5C16%u6912%u56DE%u9505%u8089*164*9*1|%u6781%u54C1%u7EA2%u7','野山椒炒肉() * 10\n尖椒回锅肉() * 9\n极品红烧肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 04:16:28',3,'2013-04-01 11:34:15','2014-10-26 04:16:29'),(141,'20130401-8214',1,666888,'1','%u91CE%u5C71%u6912%u7092%u8089*165*10*1|%u5C16%u6912%u56DE%u9505%u8089*164*9*1|%u6781%u54C1%u7EA2%u7','野山椒炒肉() * 10\n尖椒回锅肉() * 9\n极品红烧肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 04:01:03',3,'2013-04-01 11:34:26','2014-10-26 04:01:16'),(142,'20130401-8460',1,666888,'1','%u91CE%u5C71%u6912%u7092%u8089*165*10*1|%u6781%u54C1%u7EA2%u70E7%u8089*166*10*1|20','野山椒炒肉() * 10\n极品红烧肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 04:00:28',3,'2013-04-01 11:37:08','2014-10-26 04:00:33'),(143,'20130401-2503',1,666888,'2','%u9999%u5364%u9E21%u817F%u996D*154*10*1|%u82B1%u751F%u732A%u624B%u996D*155*10*1|%u7EA2%u70E7%u9C7C%u','香卤鸡腿饭() * 10\n花生猪手饭() * 10\n红烧鱼腩饭() * 9\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2013-04-01 22:10:18',3,'2013-04-01 22:10:29','2014-10-26 04:06:38'),(144,'20130408-9481',1,666888,'1','%u5C16%u6912%u56DE%u9505%u8089*164*9*1|%u91CE%u5C71%u6912%u7092%u8089*165*10*1|19','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2014-10-26 03:54:37',3,'2013-04-08 01:44:59','2014-10-26 03:54:44'),(145,'20130408-9473',1,666888,'0','%u5C16%u6912%u56DE%u9505%u8089*164*9*1|%u91CE%u5C71%u6912%u7092%u8089*165*10*1|19','尖椒回锅肉() * 9\n野山椒炒肉() * 10\n',100.00,60.00,NULL,NULL,60.00,1,NULL,'',NULL,NULL,NULL,'2013-04-08 01:44:50',-1,'2013-04-08 01:45:00','2013-04-08 01:45:00'),(146,'685450864',1,666888,NULL,'2*1|3*2|4*1','鱼香茄子饭*1\n尖椒回锅肉*2\n香菇焖鸡饭*1',80.00,72.00,NULL,NULL,72.00,1,NULL,NULL,NULL,NULL,NULL,'2014-10-29 09:47:48',NULL,'2014-10-29 09:47:48','2014-10-29 09:47:48'),(147,'683515854',1,666888,NULL,'2*1|3*2|4*1','鱼香茄子饭*1\n尖椒回锅肉*2\n香菇焖鸡饭*1',80.00,72.00,NULL,NULL,72.00,1,NULL,NULL,NULL,NULL,NULL,'2014-10-29 09:50:21',NULL,'2014-10-29 09:50:21','2014-10-29 09:50:21'),(148,'685952634',1,666888,NULL,'2*18|3*18|4*18','鱼香茄子饭*18\n尖椒回锅肉*18\n香菇焖鸡饭*18',1080.00,972.00,NULL,NULL,972.00,2,NULL,NULL,NULL,NULL,NULL,'2014-11-11 02:52:11',NULL,'2014-11-11 02:52:11','2014-11-11 02:52:11'),(149,'681859957',1,666888,NULL,'2*2','鱼香茄子饭*2',40.00,36.00,NULL,NULL,36.00,2,NULL,NULL,NULL,NULL,NULL,'2014-11-13 13:02:35',NULL,'2014-11-13 13:02:35','2014-11-13 13:02:35'),(150,'686802891',1,666888,NULL,'2*2','鱼香茄子饭*2',40.00,36.00,NULL,NULL,36.00,2,NULL,NULL,NULL,NULL,NULL,'2014-11-13 13:14:25',NULL,'2014-11-13 13:14:25','2014-11-13 13:14:25'),(151,'686267421',1,666888,'1','2*2','鱼香茄子饭*2',40.00,36.00,NULL,NULL,36.00,2,0,NULL,NULL,NULL,NULL,'2014-11-13 14:37:31',3,'2014-11-13 14:01:56','2014-11-13 14:37:34'),(152,'684452597',1,666888,'7','2*2','鱼香茄子饭*2',40.00,36.00,NULL,NULL,36.00,2,0,NULL,NULL,NULL,NULL,'2014-11-14 07:37:27',3,'2014-11-13 14:56:43','2014-11-14 07:37:28'),(153,'689182668',1,666888,'7','2*1|3*1|4*1|5*1','鱼香茄子饭*1\n尖椒回锅肉*1\n香菇焖鸡饭*1\n野山椒炒肉*1',60.00,55.00,NULL,NULL,55.00,2,0,NULL,NULL,NULL,NULL,'2014-11-13 15:07:51',3,'2014-11-13 15:07:24','2014-11-13 15:07:52'),(154,'685815333',1,666888,'1','2*1','鱼香茄子饭*1',20.00,18.00,NULL,NULL,18.00,2,0,NULL,NULL,NULL,NULL,'2014-11-14 07:34:14',3,'2014-11-14 04:01:06','2014-11-14 07:34:14'),(155,'689681862',1,666888,'0','2*1','鱼香茄子饭*1',20.00,18.00,NULL,NULL,18.00,2,0,NULL,'刘铭','18616999822','上海市','2014-11-27 01:26:11',0,'2014-11-27 01:26:11','2014-11-27 01:26:11'),(156,'686862923',1,666888,'0','2*1|4*1|5*1|6*1|7*1','鱼香茄子饭*1\n香菇焖鸡饭*1\n野山椒炒肉*1\n极品红烧肉*1\n花生猪手饭*1',40.00,39.00,NULL,NULL,39.00,2,0,NULL,'刘铭','18616999822','上海市','2014-11-27 01:35:03',0,'2014-11-27 01:35:03','2014-11-27 01:35:03'),(157,'687643214',1,666888,'0','3*1|4*1','尖椒回锅肉*1\n香菇焖鸡饭*1',40.00,36.00,NULL,NULL,36.00,2,0,NULL,'刘铭','18616999822','上海市','2014-11-28 01:50:26',0,'2014-11-28 01:50:26','2014-11-28 01:50:26'),(158,'687137477',1,666888,'0','2*1|3*1','鱼香茄子饭*1\n尖椒回锅肉*1',40.00,36.00,NULL,NULL,36.00,2,0,NULL,'刘铭','18616999822','上海市','2014-11-28 01:50:38',0,'2014-11-28 01:50:38','2014-11-28 01:50:38'),(159,'686200508',1,666888,'0','2*1|3*1|4*1','鱼香茄子饭*1\n尖椒回锅肉*1\n香菇焖鸡饭*1',60.00,54.00,NULL,NULL,54.00,2,0,NULL,'李讪讪','18616999822','上海市黄埔区','2014-11-28 02:24:13',0,'2014-11-28 02:24:13','2014-11-28 02:24:13'),(160,'682658763',1,666888,'0','2*1','鱼香茄子饭*1',20.00,18.00,NULL,NULL,18.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-08 09:35:47',0,'2014-12-08 09:35:47','2014-12-08 09:35:47'),(161,'688563456',1,666888,'0','4*1','香菇焖鸡饭*1',20.00,18.00,NULL,NULL,18.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-08 09:36:24',0,'2014-12-08 09:36:24','2014-12-08 09:36:24'),(162,'686432330',1,666888,'0','3*1','尖椒回锅肉*1',20.00,18.00,NULL,NULL,18.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-08 09:36:41',0,'2014-12-08 09:36:41','2014-12-08 09:36:41'),(163,'682539015',1,666888,'0','2*1','鱼香茄子饭*1',20.00,18.00,NULL,NULL,18.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-09 06:14:51',0,'2014-12-09 06:14:51','2014-12-09 06:14:51'),(164,'683135658',1,666888,'0','3*1','尖椒回锅肉*1',20.00,18.00,NULL,NULL,18.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-09 06:57:00',0,'2014-12-09 06:57:00','2014-12-09 06:57:00'),(165,'689424176',1,666888,'0',NULL,'尖椒回锅肉*1',20.00,18.00,7.40,5.40,13.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-10 01:35:20',0,'2014-12-10 02:35:20','2014-12-10 02:35:20'),(166,'682964508',1,666888,'0',NULL,'尖椒回锅肉*1\n尖椒回锅肉*1',40.00,36.00,4.00,0.00,36.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-10 04:18:34',0,'2014-12-10 05:18:34','2014-12-10 05:18:34'),(167,'689386565',1,666888,'0',NULL,'香菇焖鸡饭*2',40.00,36.00,14.80,10.80,25.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-10 04:21:08',0,'2014-12-10 05:21:08','2014-12-10 05:21:08'),(168,'685350377',1,666888,'0',NULL,'香菇焖鸡饭*2',40.00,36.00,14.80,10.80,25.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-10 04:23:01',0,'2014-12-10 05:23:01','2014-12-10 05:23:01'),(169,'683401193',1,666888,'0',NULL,'尖椒回锅肉*2',40.00,36.00,14.80,10.80,25.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-10 04:27:07',0,'2014-12-10 05:27:07','2014-12-10 05:27:07'),(170,'681554163',1,666888,'0','2*1','鱼香茄子饭*1',20.00,18.00,7.40,5.40,13.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-10 04:52:11',0,'2014-12-10 05:52:11','2014-12-10 05:52:11'),(171,'689760500',1,666888,'0','2*1|3*1','尖椒回锅肉*1\n尖椒回锅肉*1',40.00,36.00,14.80,10.80,25.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-10 05:34:51',0,'2014-12-10 06:34:51','2014-12-10 06:34:51'),(172,'687956073',1,666888,'0','4*1|2*2','香菇焖鸡饭*1\n香菇焖鸡饭*1',40.00,36.00,4.00,0.00,36.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-10 05:43:34',0,'2014-12-10 06:43:34','2014-12-10 06:43:34'),(173,'684691894',1,666888,'0','4*2','香菇焖鸡饭*2',40.00,36.00,4.00,0.00,36.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-10 05:45:14',0,'2014-12-10 06:45:14','2014-12-10 06:45:14'),(174,'682466077',1,666888,'0','2*2','鱼香茄子饭*2',40.00,36.00,5.80,1.80,34.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-10 13:01:26',0,'2014-12-10 14:01:26','2014-12-10 14:01:26'),(175,'681009722',1,666888,'1','2*2','鱼香茄子饭*2',40.00,36.00,4.00,0.00,36.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-22 03:20:58',3,'2014-12-10 14:03:02','2014-12-22 03:20:58'),(176,'681710378',1,666888,'1','2*2','鱼香茄子饭*2',40.00,36.00,17.80,13.80,22.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-22 03:16:21',3,'2014-12-10 15:28:38','2014-12-22 03:16:22'),(177,'687025732',1,666888,'1','2*2','鱼香茄子饭*2',40.00,36.00,16.80,12.80,23.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-16 01:47:15',3,'2014-12-10 15:29:21','2014-12-16 01:47:16'),(178,'683878066',1,666888,'1','2*2','鱼香茄子饭*2',40.00,36.00,15.80,11.80,24.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-22 03:02:33',3,'2014-12-10 15:29:51','2014-12-22 03:02:34'),(179,'681645404',1,666888,'7','2*2','鱼香茄子饭*2',40.00,36.00,15.80,11.80,24.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-16 01:47:08',3,'2014-12-10 16:08:29','2014-12-16 01:47:08'),(180,'687153988',1,666888,'1','2*3','鱼香茄子饭*3',60.00,54.00,23.20,17.20,37.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-16 01:47:00',3,'2014-12-10 16:10:15','2014-12-16 01:47:01'),(181,'688664284',1,666888,'1','2*7','鱼香茄子饭*7',140.00,126.00,20.30,6.30,119.70,2,0,NULL,'刘铭','18616999822','上海市','2014-12-16 01:25:45',3,'2014-12-11 01:23:03','2014-12-16 01:46:48'),(182,'684148199',1,666888,'1','4*2','香菇焖鸡饭*2',40.00,36.00,4.00,0.00,36.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-22 02:58:35',3,'2014-12-19 04:55:36','2014-12-22 02:58:46'),(183,'689469831',28,666888,'1','2*1','鱼香茄子饭*1',20.00,18.00,2.00,0.00,18.00,2,0,NULL,'刘大爷','18616999822','佛山市','2014-12-19 08:57:23',3,'2014-12-19 08:56:53','2014-12-19 08:57:24'),(184,'689205996',1,666888,'1','2*1','鱼香茄子饭*1',20.00,18.00,2.00,0.00,18.00,2,0,NULL,'刘铭','18616999822','上海市','2014-12-22 02:55:57',3,'2014-12-22 02:54:42','2014-12-22 02:56:09');
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
-- Table structure for table `pt_orderlog`
--

DROP TABLE IF EXISTS `pt_orderlog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_orderlog` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `orderid` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `state` int(11) DEFAULT NULL,
  `description` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `recordtime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_orderlog`
--

LOCK TABLES `pt_orderlog` WRITE;
/*!40000 ALTER TABLE `pt_orderlog` DISABLE KEYS */;
/*!40000 ALTER TABLE `pt_orderlog` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pt_page`
--

DROP TABLE IF EXISTS `pt_page`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pt_page` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ptid` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `type` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `content` varchar(5000) COLLATE utf8_unicode_ci DEFAULT NULL,
  `updatetime` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='合作商页面';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_page`
--

LOCK TABLES `pt_page` WRITE;
/*!40000 ALTER TABLE `pt_page` DISABLE KEYS */;
INSERT INTO `pt_page` VALUES (1,'666888','notice','','2013-03-31 12:19:35');
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
  `expires` datetime DEFAULT NULL,
  `tel` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `phone` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `join_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `login_time` datetime DEFAULT NULL,
  `last_login_time` datetime DEFAULT NULL COMMENT '标志',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_partner`
--

LOCK TABLES `pt_partner` WRITE;
/*!40000 ALTER TABLE `pt_partner` DISABLE KEYS */;
INSERT INTO `pt_partner` VALUES (666888,'wly','97ccd376043aedb077fc6336d8c5a27c','d435a520e50e960b','美味道',NULL,'2014-12-30 00:00:00','0757-82255311',NULL,'佛山市禅城区亲仁路白燕街9号201','2012-03-10 00:00:00','2012-03-12 00:00:00',NULL,NULL);
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
  PRIMARY KEY (`pt_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_saleconf`
--

LOCK TABLES `pt_saleconf` WRITE;
/*!40000 ALTER TABLE `pt_saleconf` DISABLE KEYS */;
INSERT INTO `pt_saleconf` VALUES (666888,0.10,0.10,0.20,0.80,10,NULL);
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
  `order_index` int(11) DEFAULT '0',
  `state` int(11) DEFAULT NULL COMMENT '0:表示禁用   1:表示正常',
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pt_shop`
--

LOCK TABLES `pt_shop` WRITE;
/*!40000 ALTER TABLE `pt_shop` DISABLE KEYS */;
INSERT INTO `pt_shop` VALUES (1,666888,'百花店','佛山市禅城区百花广场12楼103号','0757-08323123',3,1,'2014-07-23 12:42:27'),(7,666888,'鸿运配送点','汾江中路12号a座B号铺','0757-21211122',1,1,'2014-10-26 05:46:44');
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

-- Dump completed on 2014-12-25 11:02:56
