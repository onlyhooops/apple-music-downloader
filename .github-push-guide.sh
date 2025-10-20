#!/bin/bash
# GitHub 推送指南脚本
# 这个脚本将引导你安全地推送项目到 GitHub

echo "======================================================================"
echo "  Apple Music Downloader - GitHub 推送指南"
echo "======================================================================"
echo ""

# 最后安全检查
echo "🔒 执行最后安全检查..."
echo ""

echo "检查配置文件是否被正确忽略..."
IGNORED_COUNT=0
for file in config.yaml config_test.yaml dev.env; do
    if git check-ignore "$file" > /dev/null 2>&1; then
        echo "✅ $file 已被正确忽略"
        ((IGNORED_COUNT++))
    else
        echo "❌ 警告: $file 未被忽略！"
    fi
done

if [ $IGNORED_COUNT -ne 3 ]; then
    echo "❌ 错误: 配置文件未全部被忽略！"
    exit 1
fi

echo ""
echo "检查是否有敏感信息..."
# 只检查即将提交的文件
if git ls-files | xargs grep -l "YOUR_MEDIA_USER_TOKEN_HERE" 2>/dev/null | grep -v ".example"; then
    echo "❌ 警告: 发现可能的敏感信息在非示例文件中！"
    exit 1
else
    echo "✅ 未发现敏感信息泄露"
fi

echo ""
echo "======================================================================"
echo "  准备推送到 GitHub"
echo "======================================================================"
echo ""
echo "请按照以下步骤操作:"
echo ""
echo "1. 添加所有更改:"
echo "   git add ."
echo ""
echo "2. 提交更改:"
echo "   git commit -F COMMIT_MESSAGE.txt"
echo ""
echo "3. 设置远程仓库（如果是新仓库）:"
echo "   git remote add origin https://github.com/onlyhooops/apple-music-downloader.git"
echo ""
echo "4. 推送到 GitHub:"
echo "   git push -u origin main"
echo ""
echo "======================================================================"
echo "✅ 安全检查通过！可以安全推送到 GitHub"
echo "======================================================================"
