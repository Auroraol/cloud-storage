CREATE TABLE `user_repository`
(
    `id`                  bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`       bigint unsigned NOT NULL DEFAULT '0',
    `parent_id`           bigint unsigned NOT NULL DEFAULT '0' ,
    `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '0则为文件夹, 其他为文件id',
    `name`                varchar(255) NOT NULL DEFAULT '' COMMENT '文件夹名称',
    `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_repository_id` (`repository_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 8
  DEFAULT CHARSET = utf8;

-- 执行命令: goctl model mysql ddl --src user_repository.sql --dir . cache --cache  //cache --cache带缓存