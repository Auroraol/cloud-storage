CREATE TABLE IF NOT EXISTS `upload_history`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`       bigint unsigned NOT NULL DEFAULT '0' comment '用户id',
    `file_name` varchar(255) NOT NULL DEFAULT '' COMMENT '文件名',
    `size` int(11) NOT NULL DEFAULT '0' COMMENT '文件大小',
    `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '文件id',
    `status`       tinyint(1) NOT NULL DEFAULT '0' COMMENT '上传状态，0：上传中，1：上传成功，2：上传失败',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_repository_id_unique` (`repository_id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

# -- 执行命令: goctl model mysql ddl --src upload_history.sql --dir .  cache --cache  //cache --cache带缓存