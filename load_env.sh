#!/bin/bash

# Apple Music 下载器环境变量加载脚本
# 此脚本用于加载环境变量并处理配置文件中的环境变量引用

set -e

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$PROJECT_ROOT/dev.env"
CONFIG_FILE="$PROJECT_ROOT/config.yaml"
TEMP_CONFIG_FILE="$PROJECT_ROOT/config.yaml.tmp"

# 检查环境变量文件是否存在
if [[ ! -f "$ENV_FILE" ]]; then
    echo "警告: 未找到环境变量文件 $ENV_FILE"
    echo "请复制 dev.env.example 为 dev.env 并填入真实配置信息"
    exit 1
fi

# 加载环境变量
echo "正在加载环境变量..."
source "$ENV_FILE"

# 检查必要的环境变量是否存在
required_vars=(
    "APPLE_MUSIC_MEDIA_USER_TOKEN_CN"
    "APPLE_MUSIC_AUTH_TOKEN_CN"
)

for var in "${required_vars[@]}"; do
    if [[ -z "${!var:-}" ]]; then
        echo "错误: 必需的环境变量 $var 未设置"
        echo "请检查 $ENV_FILE 文件中的配置"
        exit 1
    fi
done

# 创建临时配置文件，替换环境变量引用
echo "正在处理配置文件..."
cp "$CONFIG_FILE" "$TEMP_CONFIG_FILE"

# 使用perl进行安全的替换，避免特殊字符问题
perl -i -pe "s/\${APPLE_MUSIC_MEDIA_USER_TOKEN_CN}/$APPLE_MUSIC_MEDIA_USER_TOKEN_CN/g" "$TEMP_CONFIG_FILE"
perl -i -pe "s/\${APPLE_MUSIC_AUTH_TOKEN_CN}/$APPLE_MUSIC_AUTH_TOKEN_CN/g" "$TEMP_CONFIG_FILE"

echo "环境变量加载完成"
echo "临时配置文件已创建: $TEMP_CONFIG_FILE"
echo "程序可以使用此临时配置文件运行"

# 输出加载的环境变量信息（不显示敏感内容）
echo ""
echo "已加载的环境变量:"
echo "- APPLE_MUSIC_MEDIA_USER_TOKEN_CN: 已设置"
echo "- APPLE_MUSIC_AUTH_TOKEN_CN: 已设置"
echo "- LOCAL_DEVELOPMENT: ${LOCAL_DEVELOPMENT:-false}"

echo ""
echo "使用说明:"
echo "1. 程序启动时请使用临时配置文件: $TEMP_CONFIG_FILE"
echo "2. 或者直接在程序中引用环境变量"
echo "3. 程序运行结束后可删除临时配置文件"
