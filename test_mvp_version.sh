#!/bin/bash

echo "🧪 Apple Music Downloader v2.6.0 MVP 测试脚本"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 检查二进制文件是否存在
if [ ! -f "apple-music-downloader-v2.6.0-mvp" ]; then
    echo "❌ 找不到 apple-music-downloader-v2.6.0-mvp"
    echo "正在构建..."
    make build
    exit 1
fi

echo "✅ 找到MVP版本二进制文件"
echo ""

# 显示版本信息
echo "📦 版本信息:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./apple-music-downloader-v2.6.0-mvp --version 2>&1 | head -10 || echo "运行版本检查..."
echo ""

# 显示文件大小
echo "📊 文件大小:"
ls -lh apple-music-downloader-v2.6.0-mvp | awk '{print "  大小: " $5}'
echo ""

# 显示帮助信息（测试基本功能）
echo "📖 帮助信息测试:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./apple-music-downloader-v2.6.0-mvp --help 2>&1 | head -15
echo ""

# 测试Logger配置
echo "🔍 Logger配置测试:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "当前config.yaml中的logging配置:"
grep -A 3 "logging:" config.yaml || echo "  未找到logging配置"
echo ""

# 运行基本功能测试
echo "🧪 基本功能测试:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "测试选项:"
echo "  1) 运行帮助命令测试（已完成）"
echo "  2) 测试Logger不同等级"
echo "  3) 运行实际下载测试（需要URL）"
echo ""

# Logger等级测试
echo "🔬 Logger等级测试示例:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "# 测试DEBUG等级（显示所有日志）:"
echo "  1. 编辑config.yaml，设置 level: debug"
echo "  2. 运行: ./apple-music-downloader-v2.6.0-mvp <url>"
echo ""
echo "# 测试ERROR等级（仅显示错误）:"
echo "  1. 编辑config.yaml，设置 level: error"
echo "  2. 运行: ./apple-music-downloader-v2.6.0-mvp <url>"
echo ""

# Progress系统测试
echo "🎨 Progress系统测试:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "运行任何下载都会自动测试Progress系统:"
echo "  - 进度通过Progress事件更新"
echo "  - UI监听器自动格式化显示"
echo "  - 适配器模式自动转换"
echo ""

# 对比测试建议
if [ -f "apple-music-downloader-baseline" ]; then
    echo "📊 对比测试建议:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "基线版本: apple-music-downloader-baseline"
    echo "MVP版本:  apple-music-downloader-v2.6.0-mvp"
    echo ""
    echo "对比测试方法:"
    echo "  1) 使用相同的测试URL"
    echo "  2) 观察输出差异"
    echo "  3) 对比性能和UI表现"
    echo ""
fi

# 验证重构
echo "✅ 验证重构结果:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "运行完整验证:"
echo "  ./scripts/validate_refactor.sh"
echo ""
echo "运行测试:"
echo "  make test"
echo ""
echo "运行Race检测:"
echo "  make race"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✨ MVP版本已就绪，可以开始测试！"
echo ""
echo "快速开始:"
echo "  ./apple-music-downloader-v2.6.0-mvp <your_url>"
echo ""
echo "查看文档:"
echo "  cat FINAL_SUMMARY.md"
echo "  cat MVP_COMPLETE.md"
echo "  cat CHANGELOG_v2.6.0.md"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

