#!/bin/bash

# 统计项目代码行数（mac 专用，自动忽略无用目录）

# 项目根目录（默认脚本所在文件夹为根目录，可自行修改）
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "项目目录：$PROJECT_DIR"
echo "------------------------------------"

# 统计文件数 & 总行数
RESULT=$(find "$PROJECT_DIR" \
  -type f \
  -not -path '*/node_modules/*' \
  -not -path '*/.umi/*' \
  -not -path '*/.git/*' \
  -not -path '*/dist/*' \
  -not -path '*/build/*' \
  -not -path '*/.DS_Store' \
  -not -path '*/yarn.lock' \
  -not -path '*/go.mod' \
  -not -path '*/go.sum' \
  -exec file --mime-type {} + | \
  grep -E 'text/|application/json|application/javascript|application/xml|application/x-sh' | \
  cut -d: -f1 | \
  xargs wc -l | \
  tail -n1)

FILES=$(echo "$RESULT" | awk '{print $2}')
LINES=$(echo "$RESULT" | awk '{print $1}')

echo "有效代码文件数：$FILES"
echo "总代码行数：$LINES"