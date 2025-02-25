CREATE TABLE IF NOT EXISTS  `user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '用户id',
    `version` BIGINT NOT NULL,
    `username` VARCHAR(25) DEFAULT NULL COMMENT '用户名',
    `password` VARCHAR(255) DEFAULT NULL COMMENT '密码',
    `mobile` BIGINT(11) DEFAULT NULL COMMENT '手机号',
    `nickname` VARCHAR(50) NOT NULL COMMENT '昵称',
    `gender` TINYINT(1) NOT NULL DEFAULT '1' COMMENT '性别，1：男，0：女，默认为1',
    `avatar` VARCHAR(255) DEFAULT NULL COMMENT '用户头像',
    `birthday` DATE DEFAULT NULL COMMENT '生日',
    `email` VARCHAR(254) DEFAULT NULL COMMENT '电子邮箱',
    `brief` VARCHAR(255) DEFAULT NULL COMMENT '简介|个性签名',
    `info` TEXT,
    `del_state` INt COMMENT '删除状态，0: 未删除，1：已删除',
    `delete_time` TIMESTAMP COMMENT '删除时间',
    `status` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '状态，0：正常，1：锁定，2：禁用，3：过期',
    `admin` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '是否管理员，1：是，0：否',
    `now_volume` INT(11) NOT NULL DEFAULT '0' COMMENT '当前存储容量',
    `total_volume` INT(11) NOT NULL DEFAULT '1000000000' COMMENT '最大存储容量',
    `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name_unique` (`username`),
    UNIQUE KEY `idx_mobile_unique` (`mobile`)
);

-- 执行命令: goctl model mysql ddl --src user.sql --dir .