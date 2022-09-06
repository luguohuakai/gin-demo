-- Adminer 4.8.1 MySQL 8.0.12 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

CREATE DATABASE `webauthn` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `webauthn`;

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `credential`;
CREATE TABLE `credential` (
                              `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
                              `uid` int(10) unsigned NOT NULL COMMENT '用户ID',
                              `cid` varbinary(255) NOT NULL COMMENT '凭据ID',
                              `credential` text NOT NULL COMMENT '凭据内容',
                              `created_at` datetime NOT NULL COMMENT '创建时间',
                              `updated_at` datetime NOT NULL COMMENT '更新时间',
                              `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                              PRIMARY KEY (`id`),
                              UNIQUE KEY `uid_cid` (`uid`,`cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='凭据表';


DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                        `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
                        `name` varchar(64) NOT NULL COMMENT '用户名',
                        `display_name` varchar(64) NOT NULL COMMENT '展示用户名',
                        `created_at` datetime NOT NULL COMMENT '创建时间',
                        `updated_at` datetime NOT NULL COMMENT '更新时间',
                        `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                        `status` tinyint(3) unsigned DEFAULT '1' COMMENT '1:未激活 2:注册完成',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';


-- 2022-09-05 06:09:08
