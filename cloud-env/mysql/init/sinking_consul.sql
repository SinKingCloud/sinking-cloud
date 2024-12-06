CREATE DATABASE IF NOT EXISTS `sinking_consul`;
USE sinking_consul;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sinking_apps
-- ----------------------------
CREATE TABLE `sinking_apps`  (
  `id` bigint(32) NOT NULL DEFAULT 0,
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '标题',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标识',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT 0 COMMENT '是否删除(0:否/1:是)',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_all`(`title`, `name`, `create_time`, `update_time`, `is_delete`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '应用列表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sinking_apps
-- ----------------------------
INSERT INTO `sinking_apps` VALUES (1, '沉沦网络', 'cloud-server', '2021-10-16 15:35:00', '2021-10-16 15:35:00', 0);

-- ----------------------------
-- Table structure for sinking_configs
-- ----------------------------
CREATE TABLE `sinking_configs`  (
  `id` bigint(32) NOT NULL DEFAULT 0,
  `app_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '应用标识',
  `env_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '环境标识',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '标题',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标识',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '配置内容',
  `hash` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '配置md5 hash',
  `group_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分组名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '类型(txt|json|yaml|properties)',
  `status` int(11) NOT NULL DEFAULT 0 COMMENT '状态(0:已发布/1:未发布)',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT 0 COMMENT '是否删除(0:否/1:是)',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_all`(`title`, `app_name`, `env_name`, `group_name`, `name`, `hash`, `type`, `status`, `create_time`, `update_time`, `is_delete`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '配置表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sinking_configs
-- ----------------------------
INSERT INTO `sinking_configs` VALUES (1, 'cloud-server', 'dev', '网关配置', 'gateway', '{\r\n    \"/ws\":{\r\n        \"/*\":{\r\n            \"group\":\"cloud-common\",\r\n            \"name\":\"cloud-websocket\",\r\n            \"protocol\":\"ws://\"\r\n        }\r\n    },\r\n    \"/auth\":{\r\n        \"/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-auth\",\r\n            \"protocol\":\"http://\"\r\n        }\r\n    },\r\n    \"/open\":{\r\n        \"/api/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-open\",\r\n            \"protocol\":\"http://\"\r\n        },\r\n        \"/user/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-open\",\r\n            \"protocol\":\"http://\",\r\n            \"login\":\"user\"\r\n        },\r\n        \"/admin/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-open\",\r\n            \"protocol\":\"http://\",\r\n            \"login\":\"admin\"\r\n        }\r\n    },\r\n    \"/user\":{\r\n        \"/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-user\",\r\n            \"protocol\":\"http://\",\r\n            \"login\":\"user\"\r\n        }\r\n    },\r\n    \"/admin\":{\r\n        \"/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-admin\",\r\n            \"protocol\":\"http://\",\r\n            \"login\":\"admin\"\r\n        }\r\n    }\r\n}', '69522e095ec0044a79f4438329f81de0', 'cloud-gateway', 'json', 0, '2021-10-14 17:49:58', '2021-10-14 17:49:58', 0);
INSERT INTO `sinking_configs` VALUES (2, 'cloud-server', 'dev', 'redis配置', 'redis-common', '{\"host\":\"222.187.238.201\",\"port\":\"6379\",\"password\":\"Qq@1178710004...\",\"db\":0}', '20646ad99be7fa776664e4a6e7701900', 'cloud-common', 'json', 0, '2021-10-14 17:49:58', '2021-10-14 17:49:58', 0);
INSERT INTO `sinking_configs` VALUES (3, 'cloud-server', 'dev', 'mysql配置', 'mysql-common', '{\"host\":\"222.187.238.201\",\"port\":\"3306\",\"user\":\"root\",\"password\":\"Qq@1178710004...\",\"database\":\"sinking_cloud\"}', '198ee80a6a9352b0e097b23f9d3d0264', 'cloud-common', 'json', 0, '2021-10-14 17:49:58', '2021-10-14 17:49:58', 0);
INSERT INTO `sinking_configs` VALUES (4, 'cloud-server', 'dev', 'rocketmq配置', 'rocketmq-common', '{\"server\":\"222.187.238.201:9876\"}', '9ef9d379fa1495c1b02322d346011a0f', 'cloud-common', 'json', 0, '2021-10-14 17:49:58', '2021-10-14 17:49:58', 0);
INSERT INTO `sinking_configs` VALUES (5, 'cloud-server', 'dev', '开放服务mysql配置', 'open-mysql', '{\"host\":\"222.187.238.201\",\"port\":\"3306\",\"user\":\"root\",\"password\":\"Qq@1178710004...\",\"database\":\"sinking_open\"}', '198ee80a6a9352b0e097b23f9d3d0267', 'cloud-service', 'json', 0, '2022-03-24 18:01:51', '2022-03-24 18:01:54', 0);
INSERT INTO `sinking_configs` VALUES (6, 'cloud-server', 'release', '网关配置', 'gateway', '{\r\n    \"/ws\":{\r\n        \"/*\":{\r\n            \"group\":\"cloud-common\",\r\n            \"name\":\"cloud-websocket\",\r\n            \"protocol\":\"ws://\"\r\n        }\r\n    },\r\n    \"/auth\":{\r\n        \"/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-auth\",\r\n            \"protocol\":\"http://\"\r\n        }\r\n    },\r\n    \"/open\":{\r\n        \"/api/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-open\",\r\n            \"protocol\":\"http://\"\r\n        },\r\n        \"/user/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-open\",\r\n            \"protocol\":\"http://\",\r\n            \"login\":\"user\"\r\n        },\r\n        \"/admin/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-open\",\r\n            \"protocol\":\"http://\",\r\n            \"login\":\"admin\"\r\n        }\r\n    },\r\n    \"/user\":{\r\n        \"/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-user\",\r\n            \"protocol\":\"http://\",\r\n            \"login\":\"user\"\r\n        }\r\n    },\r\n    \"/admin\":{\r\n        \"/*\":{\r\n            \"group\":\"cloud-service\",\r\n            \"name\":\"service-admin\",\r\n            \"protocol\":\"http://\",\r\n            \"login\":\"admin\"\r\n        }\r\n    }\r\n}', '69522e095ec0044a79f4438329f81de0', 'cloud-gateway', 'json', 0, '2021-10-14 17:49:58', '2021-10-14 17:49:58', 0);
INSERT INTO `sinking_configs` VALUES (7, 'cloud-server', 'release', 'redis配置', 'redis-common', '{\"host\":\"222.187.238.201\",\"port\":\"6379\",\"password\":\"Qq@1178710004...\",\"db\":0}', '20646ad99be7fa776664e4a6e7701900', 'cloud-common', 'json', 0, '2021-10-14 17:49:58', '2021-10-14 17:49:58', 0);
INSERT INTO `sinking_configs` VALUES (8, 'cloud-server', 'release', 'mysql配置', 'mysql-common', '{\"host\":\"222.187.238.201\",\"port\":\"3306\",\"user\":\"root\",\"password\":\"Qq@1178710004...\",\"database\":\"sinking_cloud\"}', '198ee80a6a9352b0e097b23f9d3d0264', 'cloud-common', 'json', 0, '2021-10-14 17:49:58', '2021-10-14 17:49:58', 0);
INSERT INTO `sinking_configs` VALUES (9, 'cloud-server', 'release', 'rocketmq配置', 'rocketmq-common', '{\"server\":\"222.187.238.201:9876\"}', '9ef9d379fa1495c1b02322d346011a0f', 'cloud-common', 'json', 0, '2021-10-14 17:49:58', '2021-10-14 17:49:58', 0);
INSERT INTO `sinking_configs` VALUES (10, 'cloud-server', 'release', '开放服务mysql配置', 'open-mysql', '{\"host\":\"222.187.238.201\",\"port\":\"3306\",\"user\":\"root\",\"password\":\"Qq@1178710004...\",\"database\":\"sinking_open\"}', '198ee80a6a9352b0e097b23f9d3d0267', 'cloud-service', 'json', 0, '2022-03-24 18:01:51', '2022-03-24 18:01:54', 0);

-- ----------------------------
-- Table structure for sinking_envs
-- ----------------------------
CREATE TABLE `sinking_envs`  (
  `id` bigint(32) NOT NULL DEFAULT 0 COMMENT '应用环境',
  `app_id` bigint(32) NOT NULL DEFAULT 0 COMMENT '应用ID',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '标题',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标识',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT 0 COMMENT '是否删除(0:否/1:是)',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_all`(`app_id`, `title`, `name`, `create_time`, `update_time`, `is_delete`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '应用环境' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sinking_envs
-- ----------------------------
INSERT INTO `sinking_envs` VALUES (1, 1, '开发环境', 'dev', '2021-10-16 15:35:00', '2021-10-16 15:35:00', 0);
INSERT INTO `sinking_envs` VALUES (2, 1, '生产环境', 'release', '2021-10-16 15:35:00', '2021-10-16 15:35:00', 0);

-- ----------------------------
-- Table structure for sinking_roles
-- ----------------------------
CREATE TABLE `sinking_roles`  (
  `id` bigint(32) NOT NULL DEFAULT 0 COMMENT '角色id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '角色名称',
  `auths` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '授权访问路径',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT 0 COMMENT '是否删除(0:否/1:是)',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_all`(`name`, `auths`(10), `create_time`, `update_time`, `is_delete`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '角色表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sinking_roles
-- ----------------------------

-- ----------------------------
-- Table structure for sinking_users
-- ----------------------------
CREATE TABLE `sinking_users`  (
  `id` bigint(32) NOT NULL DEFAULT 0,
  `role_id` bigint(32) NOT NULL DEFAULT 0 COMMENT '角色ID',
  `user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '账号',
  `pwd` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '账户名称',
  `login_time` datetime NULL DEFAULT NULL COMMENT '登陆时间',
  `login_ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '登陆ip',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT 0 COMMENT '是否删除(0:否/1:是)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sinking_users
-- ----------------------------
INSERT INTO `sinking_users` VALUES (1, 0, 'admin', '$2a$04$nK49S2fPoyZdw727GBvd6upZN5mC/MPJfpY3t9qsSmAOADdjhZuV.', '管理员', '2021-11-02 06:46:33', '152.136.19.96', NULL, '2021-11-02 06:46:33', 0);

SET FOREIGN_KEY_CHECKS = 1;
