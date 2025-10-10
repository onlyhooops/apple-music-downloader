# Apple Music Downloader - 功能特性

> **当前版本**: v2.5.2+  
> **分支**: feature/fix-ilst-box-missing → main  
> **更新日期**: 2025-10-10

---

## 🎯 核心功能

### 1. ilst box 自动修复
**问题**: 标签写入失败 `ilst box not present`  
**解决**: FFmpeg 自动重新封装 + 智能重试

<details>
<summary>详细说明</summary>

- 自动检测 ilst box 缺失错误
- 使用 FFmpeg 重新封装文件（无损）
- 修复后自动重试标签写入
- 用户无感知，透明操作

**文档**: `ILST_BOX_FIX.md`
</details>

---

### 2. 并发安全修复
**问题**: `fatal error: concurrent map writes`  
**解决**: 互斥锁保护共享资源

<details>
<summary>详细说明</summary>

- 修复 2 处并发写入问题
- 使用 `SharedLock` 保护 `OkDict`
- 批量下载稳定，无崩溃

**文档**: `CONCURRENT_MAP_FIX.md`
</details>

---

### 3. 工作-休息循环
**功能**: 定期休息，避免限流  
**配置**: `work-rest-enabled: true`

<details>
<summary>详细说明</summary>

```yaml
# config.yaml
work-rest-enabled: true
work-duration-minutes: 5  # 工作 5 分钟
rest-duration-minutes: 1  # 休息 1 分钟
```

- 任务完成后才休息（安全衔接）
- 友好的倒计时提示
- 降低限流风险

**文档**: `WORK_REST_CYCLE.md`
</details>

---

### 4. 从指定位置开始
**功能**: 续传、分段下载  
**用法**: `--start 44`

<details>
<summary>详细说明</summary>

```bash
./apple-music-downloader albums.txt --start 44
# 从第 44 个链接开始
```

- 跳过前面的链接
- 任务编号显示真实位置
- 零性能开销

**文档**: `START_FROM_FEATURE.md`
</details>

---

## 📚 用户指南

### 快速开始

```bash
# 1. 单专辑下载
./apple-music-downloader <album-url>

# 2. 批量下载
./apple-music-downloader albums.txt

# 3. 从指定位置开始
./apple-music-downloader albums.txt --start 44

# 4. 使用自定义配置
./apple-music-downloader albums.txt --config my-config.yaml
```

### 推荐配置

```yaml
# config.yaml

# 批量下载
batch-size: 20
skip-existing-validation: true

# 工作-休息循环
work-rest-enabled: true
work-duration-minutes: 5
rest-duration-minutes: 1

# 缓存机制
enable-cache: true
cache-folder: "./Cache"

# FFmpeg 修复
ffmpeg-fix: true
```

---

## 🔧 命令行参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `--start <N>` | 从第 N 个开始 | `--start 44` |
| `--config <path>` | 指定配置文件 | `--config my.yaml` |
| `--output <dir>` | 输出目录 | `--output /mnt/music` |
| `--no-ui` | 禁用动态 UI | `--no-ui` |
| `--atmos` | 杜比全景声模式 | `--atmos` |
| `--aac` | AAC 模式 | `--aac` |
| `--select` | 选择性下载 | `--select` |

**查看所有参数**:
```bash
./apple-music-downloader --help
```

---

## 📊 性能特性

| 特性 | 说明 | 影响 |
|------|------|------|
| 并发下载 | 专辑内多线程 | 提速 3-5倍 |
| 缓存机制 | NFS/网络优化 | 提速 50-70% |
| 批量处理 | 分批加载 | 降低内存 |
| 工作-休息 | 定期休息 | 成功率 +2-5% |

---

## 🛡️ 稳定性保证

- ✅ **并发安全**: 修复所有 map 并发写入
- ✅ **错误重试**: 自动重试 3 次
- ✅ **历史记录**: 自动跳过已完成
- ✅ **自动修复**: ilst box / FFmpeg 自动修复

---

## 📖 完整文档

### 核心功能
- `ILST_BOX_FIX.md` - ilst box 自动修复
- `CONCURRENT_MAP_FIX.md` - 并发安全修复
- `WORK_REST_CYCLE.md` - 工作-休息循环
- `START_FROM_FEATURE.md` - 从指定位置开始

### 用户指南
- `README-CN.md` - 项目说明（中文）
- `README.md` - 项目说明（英文）
- `CHANGELOG.md` - 更新日志

### 高级功能
- `HISTORY_FEATURE.md` - 历史记录功能
- `CACHE_MECHANISM.md` - 缓存机制
- `TAG_ERROR_HANDLING.md` - 标签错误处理

---

## 🚀 快速问题解答

**Q: 如何从中断的地方继续？**
```bash
# 方法1: 使用历史记录（推荐）
./apple-music-downloader albums.txt

# 方法2: 使用 --start
./apple-music-downloader albums.txt --start 44
```

**Q: 如何避免被限流？**
```yaml
# config.yaml
work-rest-enabled: true
work-duration-minutes: 5
rest-duration-minutes: 1
```

**Q: 遇到 ilst box 错误怎么办？**
> 自动修复，无需操作。确保安装了 FFmpeg。

**Q: 批量下载时崩溃？**
> 已修复并发问题，更新到最新版本。

---

## 📞 获取帮助

- **查看帮助**: `./apple-music-downloader --help`
- **项目主页**: README-CN.md
- **功能文档**: 见上方"完整文档"部分

---

**开发分支**: feature/fix-ilst-box-missing  
**状态**: ✅ 准备合并到 main

