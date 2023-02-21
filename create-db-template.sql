-- Active: 1664461804421@@127.0.0.1@3306
DROP DATABASE IF EXISTS `simpledouyin`;
CREATE DATABASE simpledouyin
    DEFAULT CHARACTER SET = 'utf8mb4';

USE simpledouyin;

DROP TABLE IF EXISTS `favorites`;
CREATE TABLE
    `favorites`(
        `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
        `user_id` BIGINT NOT NULL,
        `video_id` BIGINT NOT NULL,
        `action_type` INT NOT NULL
    ) ENGINE=InnoDB CHARSET=utf8 COMMENT 'favorites table';