-- Adminer 4.8.1 MySQL 8.0.12 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

CREATE TABLE `wa_credential`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `uid`        int(10) unsigned NOT NULL COMMENT '用户ID',
    `cid`        varbinary(255)   NOT NULL COMMENT '凭据ID',
    `credential` text             NOT NULL COMMENT '凭据内容',
    `created_at` datetime         NOT NULL COMMENT '创建时间',
    `updated_at` datetime         NOT NULL COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uid_cid` (`uid`, `cid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='凭据表';


CREATE TABLE `wa_user`
(
    `id`           int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name`         varchar(64)      NOT NULL COMMENT '用户名',
    `display_name` varchar(64)      NOT NULL COMMENT '展示用户名',
    `created_at`   datetime         NOT NULL COMMENT '创建时间',
    `updated_at`   datetime         NOT NULL COMMENT '更新时间',
    `deleted_at`   datetime            DEFAULT NULL COMMENT '删除时间',
    `status`       tinyint(3) unsigned DEFAULT '1' COMMENT '1:未激活 2:注册完成',
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户表';


-- 2022-09-05 06:09:08


# ALTER TABLE `wa_credential`
#     ADD `public_key` varbinary(1024) NOT NULL COMMENT '公钥' AFTER `cid`;

CREATE TABLE `wa_admin`
(
    `id`         int unsigned     NOT NULL COMMENT 'ID' AUTO_INCREMENT PRIMARY KEY,
    `username`   char(20)         NOT NULL COMMENT '账号',
    `password`   varchar(255)     NOT NULL COMMENT '密码',
    `avatar`     varchar(255)     NOT NULL DEFAULT '' COMMENT '头像',
    `created_at` datetime         NOT NULL COMMENT '创建时间',
    `updated_at` datetime         NOT NULL COMMENT '更新时间',
    `deleted_at` datetime         NULL COMMENT '删除时间',
    `status`     tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态 1:正常 2:禁用 3:正常+webauthn'
) COMMENT ='管理员表' ENGINE = 'InnoDB'
                  COLLATE 'utf8mb4_general_ci';


INSERT INTO `wa_admin` (`username`, `password`, `avatar`, `created_at`, `updated_at`, `deleted_at`, `status`)
VALUES ('srun', sha1('Srun@4000'), '', now(), now(), NULL, '1');

CREATE TABLE `wa_admin_credential`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `uid`        int(10) unsigned NOT NULL COMMENT '管理员ID',
    `cid`        varbinary(255)   NOT NULL COMMENT '凭据ID',
    `credential` text             NOT NULL COMMENT '凭据内容',
    `created_at` datetime         NOT NULL COMMENT '创建时间',
    `updated_at` datetime         NOT NULL COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uid_cid` (`uid`, `cid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='管理员凭据表';

ALTER TABLE `wa_admin`
    ADD UNIQUE `username` (`username`);

CREATE TABLE `wa_app`
(
    `id`         int(10) unsigned                        NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `app_id`     char(20) COLLATE utf8mb4_general_ci     NOT NULL COMMENT 'APP ID',
    `app_secret` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '密钥',
    `status`     tinyint(3) unsigned                     NOT NULL DEFAULT '1' COMMENT '1:正常 2:禁用',
    `created_at` datetime                                NOT NULL COMMENT '创建时间',
    `updated_at` datetime                                NOT NULL COMMENT '更新时间',
    `deleted_at` datetime                                         DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `app_id` (`app_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='APP表';

ALTER TABLE `wa_user`
    ADD `aid` int(10) unsigned NOT NULL COMMENT 'wa_app表对应ID' AFTER `id`;

ALTER TABLE `wa_user`
    ADD UNIQUE `aid_name` (`aid`, `name`),
    DROP INDEX `name`;
