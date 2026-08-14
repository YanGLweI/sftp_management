-- SFTP 模块配置表创建脚本
-- 用于存储标签上传和中国联通模块的登录方式、可登录角色以及双控开关配置

CREATE TABLE IF NOT EXISTS `t_sftp_module_config` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `module_name` varchar(64) NOT NULL COMMENT '模块名称：hotlabel-标签上传，chinaunicom-中国联通',
  `login_type` varchar(32) NOT NULL DEFAULT 'ldap' COMMENT '登录类型：local-本地登录，ldap-LDAP 登录',
  `enabled_roles` json DEFAULT NULL COMMENT '允许登录的角色 ID 列表（JSON 数组格式）',
  `dual_auth_enabled` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用双控：0-关闭，1-启用',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_module_name` (`module_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='SFTP 模块配置表';

-- 插入初始数据
INSERT INTO `t_sftp_module_config` (`module_name`, `login_type`, `enabled_roles`, `dual_auth_enabled`) 
VALUES 
  ('hotlabel', 'ldap', '[]', 0),
  ('chinaunicom', 'ldap', '[]', 1)
ON DUPLICATE KEY UPDATE `updated_at` = CURRENT_TIMESTAMP;
