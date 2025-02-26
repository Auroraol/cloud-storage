CREATE TABLE IF NOT EXISTS `audit`(
    `id` INT UNSIGNED AUTO_INCREMENT COMMENT '操作记录ID',
    `user_id` bigint unsigned NOT NULL DEFAULT '0' comment '用户id',
    `content` TEXT NOT NULL COMMENT '操作内容',
    `flag` tinyint(1) NOT NULL default 0 comment '操作类型，0：上传，1：下载，2：删除，3.恢复 4：重命名，5：移动，6：复制，7：创建文件夹，8：修改文件',
    `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '0则为文件夹, 其他为文件id',
    `file_name`    varchar(255) NOT NULL DEFAULT '' COMMENT '文件夹名称',
    `create_time` int(11) NOT NULL default 0 comment '创建时间戳',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 comment '操作记录表';

-- 执行命令: goctl model mysql ddl --src audit.sql --dir . cache --cache