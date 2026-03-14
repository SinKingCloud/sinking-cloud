#!/bin/bash

# 遍历当前目录下的所有子目录并执行命令
# 用法: ./run.sh <command>
# 示例: ./run.sh go mod tidy

if [ $# -eq 0 ]; then
    echo "用法: $0 <command>"
    echo "示例: $0 go mod tidy"
    exit 1
fi

# 遍历当前目录下的所有子目录
for dir in */; do
    # 检查是否为目录（排除非目录文件）
    if [ -d "$dir" ]; then
        echo "=== 进入目录: $dir ==="
        cd "$dir" || continue
        # 执行传入的命令
        eval "$@"
        # 返回上级目录
        cd ..
        echo ""
    fi
done

echo "=== 所有目录执行完成 ==="