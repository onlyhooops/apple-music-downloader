#!/bin/bash
# Git 推送命令集合

echo "======================================================================"
echo "  Git 推送和标签发布"
echo "======================================================================"
echo ""

# 1. 提交所有更改
echo "📝 1. 提交更改..."
git add .
git commit -F COMMIT_MESSAGE.txt

echo ""
echo "✅ 代码已提交"
echo ""

# 2. 创建标签
echo "🏷️  2. 创建版本标签..."
git tag -a v1.1.0 -F TAG_v1.1.0.txt

echo ""
echo "✅ 标签 v1.1.0 已创建"
echo ""

# 3. 推送代码和标签
echo "🚀 3. 推送到 GitHub..."
git push origin main
git push origin v1.1.0

echo ""
echo "======================================================================"
echo "✅ 推送完成！"
echo "======================================================================"
echo ""
echo "查看标签: https://github.com/onlyhooops/apple-music-downloader/releases/tag/v1.1.0"
echo ""
