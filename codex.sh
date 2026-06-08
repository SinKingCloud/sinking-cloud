#!/bin/bash

# ====================== 内置账号配置（已修复：用 | 分隔） ======================
# 格式：序号|BaseURL|API_KEY   【重要：不要用 : 分隔】
ACCOUNTS=(
  "1|https://api.freemodel.dev|fe_oa_69a59f2bb01180b785c085a9cf289e74ba8b1eb0657c51bd"
  "2|https://api.freemodel.dev|fe_oa_ed26000741b1ff0fc76dfcd31433a010fc11cb8d7e206e99"
  "3|https://api.freemodel.dev|fe_oa_3763f9d7646dde4bac1f4f8e741fd8b7b68d0a1720d66123"
  "4|https://new.sharedchat.cc/codex|sk-2038bac55b3211f18dcf00163e012d40"
  "5|https://www.mlyapi.com|sk-c0ebb12417a0947adf449a4f6d6bb6d3f942246cb9ad030cf307526b234882e9"
)
# ======================================================================

# 配置文件路径
AUTH_JSON="$HOME/.codex/auth.json"
CONFIG_TOML="$HOME/.codex/config.toml"

# ====================== 核心函数：写入配置 ======================
update_config() {
  local NEW_URL="$1"
  local NEW_KEY="$2"

  # 1. 写入 auth.json（不存在自动创建）
  mkdir -p "$HOME/.codex"
  cat > "$AUTH_JSON" << EOF
{
  "OPENAI_API_KEY": "${NEW_KEY}"
}
EOF

  # 2. 替换 config.toml 的 base_url
  if [ -f "$CONFIG_TOML" ]; then
    sed -i.bak "s#^base_url = \".*\"#base_url = \"$NEW_URL\"#" "$CONFIG_TOML"
    rm -f "$CONFIG_TOML.bak"
  else
    echo -e "❌ 警告：未找到 $CONFIG_TOML，仅更新了 API Key"
  fi

  echo -e "\n✅ 配置更新成功！"
  echo "🔗 Base URL: $NEW_URL"
  echo "🔑 API Key: ${NEW_KEY:0:10}...${NEW_KEY: -10}"
}

# ====================== 命令：show（查看当前配置） ======================
cmd_show() {
  echo "📋 当前 Codex 配置"
  echo "──────────────────────────"
  # 读取当前URL
  if [ -f "$CONFIG_TOML" ]; then
    CURRENT_URL=$(grep '^base_url' "$CONFIG_TOML" | awk -F '"' '{print $2}')
    echo "🔗 当前 Base URL: $CURRENT_URL"
  else
    echo "🔗 当前 Base URL: 未配置"
  fi
  # 读取当前Key
  if [ -f "$AUTH_JSON" ]; then
    CURRENT_KEY=$(grep 'OPENAI_API_KEY' "$AUTH_JSON" | awk -F '"' '{print $4}')
    echo "🔑 当前 API Key: ${CURRENT_KEY:0:10}...${CURRENT_KEY: -10}"
  else
    echo "🔑 当前 API Key: 未配置"
  fi

  echo -e "\n📇 内置账号列表"
  echo "──────────────────────────"
  for acc in "${ACCOUNTS[@]}"; do
    IFS='|' read -r num url key <<< "$acc"
    echo "$num: $url"
  done
}

# ====================== 命令：set（更新配置） ======================
cmd_set() {
  # 情况1：set 序号
  if [ $# -eq 1 ] && [[ $1 =~ ^[0-9]+$ ]]; then
    local SELECTED_NUM="$1"
    local FOUND=0
    for acc in "${ACCOUNTS[@]}"; do
      IFS='|' read -r num url key <<< "$acc"
      if [ "$num" = "$SELECTED_NUM" ]; then
        update_config "$url" "$key"
        FOUND=1
        break
      fi
    done
    [ $FOUND -ne 1 ] && echo "❌ 错误：未找到序号 $SELECTED_NUM" && exit 1
    return
  fi

  # 情况2：set url key
  if [ $# -eq 2 ]; then
    update_config "$1" "$2"
    return
  fi

  # 用法错误
  echo "❌ set 命令用法错误！"
  echo "  切换内置账号：codex.sh set 序号"
  echo "  自定义配置：codex.sh set url key"
  exit 1
}

# ====================== 主命令入口 ======================
case "$1" in
  show)
    cmd_show
    ;;
  set)
    shift
    cmd_set "$@"
    ;;
  *)
    echo "📖 可用命令"
    echo "──────────────────────────"
    echo "codex.sh show          # 查看当前配置 + 内置账号"
    echo "codex.sh set 序号      # 一键切换内置账号"
    echo "codex.sh set url key   # 手动设置自定义配置"
    ;;
esac