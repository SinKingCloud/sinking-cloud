#!/bin/bash
# 获取脚本所在目录的绝对路径
CURRENT_DIR="$(dirname "$(realpath "$0")")"

# 定义需要执行go mod tidy的目录列表
DIRS=(
  "cloud-common"
  "cloud-gateway"
  "cloud-websocket"
  "cloud-service/service-admin"
  "cloud-service/service-auth"
  "cloud-service/service-open"
  "cloud-service/service-user"
)

# 遍历目录列表并执行go mod tidy命令
for DIR in "${DIRS[@]}"; do
  cd "$CURRENT_DIR/$DIR" || exit 1  # 如果切换目录失败则退出脚本
  go mod tidy  # 执行go mod tidy命令清理模块依赖关系
done
