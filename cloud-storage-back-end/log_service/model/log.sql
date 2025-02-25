CREATE TABLE IF NOT EXISTS `logfile`(
    `id` INT UNSIGNED AUTO_INCREMENT COMMENT '日志文件ID',
    `user_id` INT NOT NULL,
    `name` VARCHAR(200) NOT NULL COMMENT '日志文件名',
    `host` LONGTEXT NOT NULL COMMENT '主机信息',
    `path` VARCHAR(1024) NOT NULL COMMENT '日志文件路径',
    `create_time` DATETIME COMMENT '创建时间',
    `comment` VARCHAR(200) COMMENT '备注',
    `monitor_choice` INT COMMENT '监控选择',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;