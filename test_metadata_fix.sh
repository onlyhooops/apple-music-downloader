#!/bin/bash

# 测试元数据音质标签修复
# 使用新版本重新下载一首歌，然后检查元数据

echo "=================================================="
echo "🧪 测试元数据音质标签修复"
echo "=================================================="
echo ""

# 检查新版本是否存在
if [ ! -f "apple-music-downloader-v2.6.0-metadata-fix" ]; then
    echo "❌ 未找到新版本：apple-music-downloader-v2.6.0-metadata-fix"
    echo "请先运行: go build -o apple-music-downloader-v2.6.0-metadata-fix"
    exit 1
fi

echo "✅ 找到新版本：apple-music-downloader-v2.6.0-metadata-fix"
echo ""

echo "📝 使用方法："
echo ""
echo "1️⃣ 使用新版本下载专辑："
echo "   ./apple-music-downloader-v2.6.0-metadata-fix \"https://music.apple.com/cn/album/head-hunters/158571524\""
echo ""
echo "2️⃣ 检查下载的文件元数据："
echo "   exiftool \"路径/04. Vein Melter.m4a\" | grep -i album"
echo ""
echo "3️⃣ 期望看到："
echo "   Album                           : Head Hunters Hi-Res Lossless  ✅"
echo "   Album Sort                      : Head Hunters Hi-Res Lossless  ✅"
echo ""

