CREATE TABLE `ssh_info` (
            `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
            `user_id` bigint(20) NOT NULL COMMENT '用户ID',
            `host` varchar(255) NOT NULL COMMENT '主机地址',
            `port` int(11) NOT NULL COMMENT '端口号',
            `username` varchar(255) NOT NULL COMMENT '用户名',
            `password` varchar(255) NOT NULL COMMENT '密码',
            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
            PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='SSH信息表';
