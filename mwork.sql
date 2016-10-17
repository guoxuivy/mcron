/*
SQLyog Ultimate v11.24 (32 bit)
MySQL - 5.5.47 : Database - mwork
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`mwork` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `mwork`;

/*Table structure for table `job_list` */

DROP TABLE IF EXISTS `job_list`;

CREATE TABLE `job_list` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `schedule_expr` varchar(200) NOT NULL COMMENT '时间规则',
  `desc` varchar(200) NOT NULL COMMENT '描述',
  `shell` varchar(200) NOT NULL COMMENT '脚本命令',
  `ip` varchar(50) NOT NULL DEFAULT '127.0.0.1' COMMENT '执行客户端IP',
  `status` tinyint(2) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8;

/*Table structure for table `job_log` */

DROP TABLE IF EXISTS `job_log`;

CREATE TABLE `job_log` (
  `id` int(20) NOT NULL AUTO_INCREMENT,
  `job_id` int(11) NOT NULL COMMENT '任务id',
  `action` varchar(100) NOT NULL COMMENT '任务动作',
  `log` varchar(200) DEFAULT NULL COMMENT '日志内容',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1839 DEFAULT CHARSET=utf8;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
