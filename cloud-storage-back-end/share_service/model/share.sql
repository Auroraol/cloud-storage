CREATE TABLE IF NOT EXISTS `share_basic` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL DEFAULT '0',
  `repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '公共池中的唯一标识',
  `user_repository_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户池子中的唯一标识',
  `expired_time` int(11) NOT NULL DEFAULT '0' COMMENT '失效时间，单位秒, 【0-永不失效】',
  `click_num` int(11) NOT NULL DEFAULT '0' COMMENT '点击次数',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `code` VARCHAR(30) NOT NULL DEFAULT '' COMMENT '提取码',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_repository_id` (`repository_id`),
  KEY `idx_user_repository_id` (`user_repository_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- 执行命令: goctl model mysql ddl --src share.sql --dir . cache --cache
--
-- SELECT
--     share_basic.id,
--     share_basic.repository_id,
--     user_repository.name,
--     repository_pool.ext,
--     repository_pool.path,
--     repository_pool.size,
--     share_basic.click_num,
--     User.username AS owner,
--     User.avatar,
--     share_basic.expired_time,
--     share_basic.update_time
-- FROM
--     share_basic
--         LEFT JOIN
--     repository_pool ON repository_pool.identity = share_basic.repository_id
--         LEFT JOIN
--     user_repository ON user_repository.id = share_basic.user_repository_id
--         LEFT JOIN
--     User ON share_basic.user_id = User.id
-- WHERE
--     share_basic.user_id = ?
--   AND share_basic.deleted_at IS NULL
-- ORDER BY
--     share_basic.update_time DESC;