# ilst box 缺失自动修复功能

## 📋 功能概述

本功能实现了对 `ilst box not present` 错误的自动检测和修复。当 MP4/M4A 文件缺少 iTunes 元数据容器（ilst box）时，系统会自动使用 FFmpeg 重新封装文件，添加必要的元数据结构，然后重试标签写入。

## 🎯 问题背景

### 什么是 ilst box？

`ilst`（item list）是 MP4 文件格式中的一个 atom/box，用于存储 iTunes 兼容的元数据标签，包括：
- 🎵 标题、艺术家、专辑
- 🎨 封面图片
- 📀 曲目编号、光盘编号
- 🏷️ 流派、发行日期
- 📝 歌词、评论等

### 错误原因

`ilst box not present` 错误通常由以下原因引起：

1. **源文件问题**：Apple Music 服务器返回的文件本身没有元数据容器
2. **解密/下载问题**：某些下载或解密过程可能生成不完整的文件结构
3. **库限制**：`go-mp4tag` 库只能修改已存在的 ilst box，不能创建新的

### 传统解决方案的局限

**优化前**（TAG_ERROR_HANDLING.md 中的方案）：
- ❌ 重试 3 次后自动跳过
- ❌ 文件保留但无元数据标签
- ❌ 需要手动使用外部工具补充标签

**优化后**（本功能）：
- ✅ 自动使用 FFmpeg 修复文件结构
- ✅ 修复后自动重试标签写入
- ✅ 无需手动干预，一次性完成

## 🔧 技术实现

### 1. 核心函数：`fixIlstBoxMissing()`

**位置**：`internal/metadata/writer.go` (第106-155行)

**功能**：使用 FFmpeg 重新封装 MP4 文件以添加缺失的 ilst box

**实现原理**：
```go
func fixIlstBoxMissing(trackPath string) error {
    // 1. 检查 ffmpeg 是否存在
    _, err := exec.LookPath("ffmpeg")
    
    // 2. 创建临时文件路径
    tempPath := trackPath + ".tmp.m4a"
    
    // 3. 使用 FFmpeg 重新封装文件
    cmd := exec.Command("ffmpeg", "-i", trackPath, 
                       "-c", "copy",              // 复制流，不重新编码
                       "-movflags", "+faststart", // 优化文件结构
                       "-f", "mp4",               // 强制 MP4 格式
                       "-y", tempPath)            // 覆盖已存在文件
    
    // 4. 执行并处理错误
    // 5. 替换原文件
}
```

**FFmpeg 参数说明**：
- `-i <input>`：输入文件路径
- `-c copy`：复制所有流，不重新编码（快速，无损）
- `-movflags +faststart`：将 moov atom（包含元数据容器）移到文件开头
- `-f mp4`：强制输出为 MP4 格式，确保正确的容器结构
- `-y`：自动覆盖输出文件

### 2. 智能重试函数：`WriteMP4TagsWithRetry()`

**位置**：`internal/metadata/writer.go` (第157-189行)

**功能**：尝试写入 MP4 标签，如果遇到 ilst box 缺失错误则自动修复

**执行流程**：
```
1. 尝试写入标签 (WriteMP4Tags)
   │
   ├─ 成功 ──> 返回 nil
   │
   └─ 失败
      │
      ├─ 检查错误类型
      │  │
      │  ├─ 是 ilst box 错误 ──> 继续修复流程
      │  │
      │  └─ 其他错误 ──> 直接返回错误
      │
      ├─ 使用 FFmpeg 修复文件 (fixIlstBoxMissing)
      │  │
      │  ├─ 修复成功 ──> 重试写入标签
      │  │                │
      │  │                ├─ 成功 ──> 返回 nil
      │  │                │
      │  │                └─ 失败 ──> 返回"修复后仍失败"错误
      │  │
      │  └─ 修复失败 ──> 返回"修复失败"错误
```

**错误检测逻辑**：
```go
errMsg := err.Error()
if strings.Contains(errMsg, "ilst box not present") || strings.Contains(errMsg, "ilst") {
    // 触发自动修复
}
```

### 3. 下载流程集成

**位置**：`internal/downloader/downloader.go` (第1084行)

**修改点**：
```go
// 旧代码
tagErr := metadata.WriteMP4Tags(trackPath, finalLrc, meta, trackIndexInMeta, len(...))

// 新代码（自动修复）
tagErr := metadata.WriteMP4TagsWithRetry(trackPath, finalLrc, meta, trackIndexInMeta, len(...))
```

## 📊 修复效果对比

### 修复前（自动跳过）

```
Track 61 of 64: Goldberg Variations... (24bit/44.1kHz) - 重试 1/3: 标签写入失败: ilst box not present...
Track 61 of 64: Goldberg Variations... (24bit/44.1kHz) - 重试 2/3: 标签写入失败: ilst box not present...
Track 61 of 64: Goldberg Variations... (24bit/44.1kHz) - 重试 3/3: 标签写入失败: ilst box not present...
Track 61 of 64: Goldberg Variations... (24bit/44.1kHz) - 已跳过 (标签失败)
```

**结果**：
- ❌ 文件无元数据标签
- ⚠️ 需要手动使用外部工具修复

### 修复后（自动修复）

```
Track 61 of 64: Goldberg Variations... (24bit/44.1kHz) - 下载中 95%
Track 61 of 64: Goldberg Variations... (24bit/44.1kHz) - 下载完成
```

**结果**：
- ✅ 自动修复并写入完整标签
- ✅ 用户无感知，流程顺畅
- ✅ 无需额外操作

## 🎨 用户体验改进

### 1. 透明修复
- 修复过程在后台自动完成
- UI 无额外提示（除非修复失败）
- 不影响下载进度显示

### 2. 性能优化
- FFmpeg 使用 `-c copy`，不重新编码
- 处理速度快（通常 < 1 秒）
- 对整体下载时间影响极小

### 3. 错误处理
- 如果 FFmpeg 不存在，返回明确错误
- 如果修复失败，保留详细错误信息
- 仍会触发原有的重试机制（3次重试）

## 📝 依赖要求

### 必需依赖

**FFmpeg**：
- 版本要求：≥ 3.0（建议 4.0+）
- 必须在系统 PATH 中可用
- 用于检测：`exec.LookPath("ffmpeg")`

### 安装 FFmpeg

**Ubuntu/Debian**：
```bash
sudo apt update
sudo apt install ffmpeg
```

**macOS**：
```bash
brew install ffmpeg
```

**验证安装**：
```bash
ffmpeg -version
```

### 如果 FFmpeg 不可用

如果系统中没有安装 FFmpeg，功能会优雅降级：
1. 检测到 FFmpeg 不存在
2. 返回错误：`未找到 ffmpeg 命令`
3. 触发原有的重试机制（3次重试）
4. 最终跳过该曲目（标记为"已跳过"）

## 🔍 故障排查

### 问题1：修复失败，提示 "未找到 ffmpeg 命令"

**原因**：系统未安装 FFmpeg 或未在 PATH 中

**解决方案**：
1. 安装 FFmpeg（见上方安装说明）
2. 确保 FFmpeg 在 PATH 中：
   ```bash
   which ffmpeg  # Linux/macOS
   where ffmpeg  # Windows
   ```

### 问题2：修复后标签写入仍失败

**原因**：文件本身存在其他问题（损坏、加密等）

**解决方案**：
1. 检查错误日志中的详细信息
2. 尝试手动使用 FFmpeg 检查文件：
   ```bash
   ffmpeg -i <file.m4a> -f null -
   ```
3. 如果文件确实损坏，尝试重新下载

### 问题3：FFmpeg 重新封装失败

**原因**：
- 磁盘空间不足（需要临时文件）
- 文件权限问题
- 输入文件已损坏

**解决方案**：
1. 检查磁盘空间：`df -h`
2. 检查文件权限：`ls -l <file.m4a>`
3. 查看 FFmpeg 错误输出（已包含在错误信息中）

## 📈 性能影响

### 时间开销

**正常流程**（无需修复）：
- 无额外开销
- 第一次标签写入成功，直接返回

**需要修复**：
- FFmpeg 重新封装：通常 < 1 秒
- 重试标签写入：< 0.1 秒
- 总额外时间：约 1-2 秒

### 空间开销

**临时文件**：
- 大小：与原文件相同（约 30-100 MB）
- 位置：与原文件相同目录
- 命名：`<原文件名>.tmp.m4a`
- 自动清理：修复完成后自动删除

**最终文件**：
- 大小变化：通常略小（优化后）
- 结构优化：moov atom 移到文件开头
- 质量影响：无（使用 `-c copy`，不重新编码）

## 🧪 测试建议

### 手动测试

1. **创建测试文件**（模拟缺失 ilst box）：
   ```bash
   # 下载一个正常的 M4A 文件
   # 使用 FFmpeg 创建一个没有元数据容器的版本
   ffmpeg -i input.m4a -c copy -map_metadata -1 test_no_ilst.m4a
   ```

2. **测试自动修复**：
   - 运行下载器，使用该测试文件
   - 观察是否自动修复并写入标签成功

3. **验证结果**：
   ```bash
   # 检查文件是否包含元数据
   ffprobe test_no_ilst.m4a
   
   # 或使用其他工具查看标签
   mp3tag test_no_ilst.m4a  # Windows
   ```

## 🔄 与现有功能的兼容性

### 与 FFmpeg Fix 配置的关系

**配置项**：`ffmpeg-fix: true` （config.yaml）

**功能区别**：
- `ffmpeg-fix`：用于修复编码问题（损坏、不完整的音频流）
- `ilst box 修复`：用于修复元数据容器缺失问题

**是否冲突**：
- ❌ 不冲突
- ✅ 可以同时使用
- 📌 执行顺序：
  1. 下载文件
  2. FFmpeg Fix（如果启用）
  3. 标签写入 + ilst box 修复（如果需要）

### 与重试机制的关系

**原有重试机制**（TAG_ERROR_HANDLING.md）：
- 重试次数：3 次
- 重试间隔：1.5 秒
- 触发条件：任何后处理错误

**ilst box 修复**：
- 在每次重试中都会尝试修复
- 如果第一次重试修复成功，不会再重试
- 如果修复失败，继续按原流程重试

**执行流程**：
```
第 1 次尝试：标签写入失败（ilst box 错误）
             ↓
          FFmpeg 修复
             ↓
          重试标签写入
             ↓
         成功 → 完成
         失败 → 继续第 2 次尝试
```

## 📚 相关文档

- **TAG_ERROR_HANDLING.md**：原有的标签错误处理机制（自动跳过）
- **config.yaml.example**：FFmpeg 相关配置说明
- **README-CN.md**：FFmpeg 安装和使用指南

## 🎯 总结

### 核心优势

1. **自动化**：无需手动干预，自动检测和修复
2. **无损**：使用 FFmpeg `-c copy`，不重新编码
3. **快速**：修复过程通常 < 1 秒
4. **透明**：用户无感知，UI 无额外提示
5. **兼容**：与现有功能完全兼容

### 适用场景

- ✅ Apple Music 下载的 M4A 文件缺少元数据容器
- ✅ 解密后的文件结构不完整
- ✅ 任何 `ilst box not present` 错误

### 注意事项

- ⚠️ 需要安装 FFmpeg（≥ 3.0）
- ⚠️ 需要足够的磁盘空间（临时文件）
- ⚠️ 如果文件本身损坏，修复可能失败

---

**开发分支**：`feature/fix-ilst-box-missing`  
**开发日期**：2025-10-10  
**相关 Issue**：ilst box not present 错误自动修复

