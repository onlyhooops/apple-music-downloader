#!/bin/bash

# 批量更新 M4A 文件的 Album 和 AlbumSort 元数据
# 从文件夹名称中提取音质标签并添加到专辑名称中

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=================================================="
echo "🔧 批量更新专辑元数据 - 添加音质标签"
echo -e "==================================================${NC}"
echo ""

# 检查是否安装了 exiftool
if ! command -v exiftool &> /dev/null; then
    echo -e "${RED}❌ 未找到 exiftool 命令${NC}"
    echo ""
    echo "请先安装 exiftool:"
    echo "  macOS:  brew install exiftool"
    echo "  Ubuntu: sudo apt install libimage-exiftool-perl"
    echo "  CentOS: sudo yum install perl-Image-ExifTool"
    exit 1
fi

# 获取目标目录
TARGET_DIR="${1:-.}"

if [ ! -d "$TARGET_DIR" ]; then
    echo -e "${RED}❌ 目录不存在: $TARGET_DIR${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 目标目录: $TARGET_DIR${NC}"
echo ""

# 统计
total_files=0
updated_files=0
skipped_files=0
error_files=0

# 递归查找所有 m4a 文件
while IFS= read -r -d '' file; do
    ((total_files++))
    
    # 获取文件所在目录名称
    dir_name=$(basename "$(dirname "$file")")
    
    # 提取音质标签（假设格式为：专辑名 音质标签）
    # 支持的标签: Hi-Res Lossless, Alac, Dolby Atmos, Aac 256
    quality_tag=""
    
    if [[ "$dir_name" =~ " Hi-Res Lossless"$ ]]; then
        quality_tag="Hi-Res Lossless"
        album_name="${dir_name% Hi-Res Lossless}"
    elif [[ "$dir_name" =~ " Alac"$ ]]; then
        quality_tag="Alac"
        album_name="${dir_name% Alac}"
    elif [[ "$dir_name" =~ " Dolby Atmos"$ ]]; then
        quality_tag="Dolby Atmos"
        album_name="${dir_name% Dolby Atmos}"
    elif [[ "$dir_name" =~ " Aac 256"$ ]]; then
        quality_tag="Aac 256"
        album_name="${dir_name% Aac 256}"
    else
        # 目录名称中没有音质标签，跳过
        echo -e "${YELLOW}⏭  跳过 (无音质标签): $(basename "$file")${NC}"
        ((skipped_files++))
        continue
    fi
    
    # 读取当前的 Album 字段
    current_album=$(exiftool -s -s -s -Album "$file" 2>/dev/null || echo "")
    
    # 检查是否已经包含音质标签
    if [[ "$current_album" =~ " $quality_tag"$ ]]; then
        echo -e "${YELLOW}⏭  跳过 (已有标签): $(basename "$file")${NC}"
        ((skipped_files++))
        continue
    fi
    
    # 如果当前 Album 为空，使用从文件夹名提取的专辑名
    if [ -z "$current_album" ]; then
        current_album="$album_name"
    fi
    
    # 构建新的专辑名称（添加音质标签）
    new_album="$current_album $quality_tag"
    
    # 更新元数据
    echo -e "${BLUE}📝 更新: $(basename "$file")${NC}"
    echo "   旧: $current_album"
    echo "   新: $new_album"
    
    if exiftool -overwrite_original \
        -Album="$new_album" \
        -AlbumSort="$new_album" \
        "$file" &>/dev/null; then
        echo -e "${GREEN}   ✅ 更新成功${NC}"
        ((updated_files++))
    else
        echo -e "${RED}   ❌ 更新失败${NC}"
        ((error_files++))
    fi
    echo ""
    
done < <(find "$TARGET_DIR" -type f -name "*.m4a" -print0)

# 显示统计结果
echo -e "${BLUE}=================================================="
echo "📊 处理统计"
echo -e "==================================================${NC}"
echo -e "总文件数:   ${BLUE}$total_files${NC}"
echo -e "已更新:     ${GREEN}$updated_files${NC}"
echo -e "已跳过:     ${YELLOW}$skipped_files${NC}"
echo -e "失败:       ${RED}$error_files${NC}"
echo ""

if [ $updated_files -gt 0 ]; then
    echo -e "${GREEN}✅ 批量更新完成！${NC}"
    echo ""
    echo "验证方法："
    echo "  exiftool \"文件路径.m4a\" | grep -i album"
else
    echo -e "${YELLOW}⚠️  没有文件需要更新${NC}"
fi

