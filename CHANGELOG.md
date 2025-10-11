# Changelog

All notable changes to this project will be documented in this file.

## [v2.6.0] - 2025-10-11

### 🎉 重大改进

#### 架构重构
本次版本进行了UI与LOG模块的全面重构，解决了长期存在的日志竞争、UI刷新性能等问题。

### ✨ 新增功能

#### 1. 统一日志系统 (internal/logger)
- ✅ 4级日志控制（DEBUG/INFO/WARN/ERROR）
- ✅ 配置化日志输出（stdout/stderr/文件）
- ✅ 时间戳显示控制
- ✅ 线程安全的日志记录

**配置示例**:
```yaml
logging:
  level: info              # debug/info/warn/error
  output: stdout           # stdout/stderr/文件路径
  show_timestamp: false    # UI模式下关闭时间戳
```

#### 2. Progress事件系统 (internal/progress)
- ✅ 观察者模式的事件架构
- ✅ 解耦UI更新与业务逻辑
- ✅ 适配器模式实现平滑迁移
- ✅ 支持多个监听器扩展

#### 3. UI智能简化与优化 (internal/ui)
- ✅ 自适应终端宽度显示（FullMode/CompactMode/MinimalMode）
- ✅ 智能简化长曲目名称和质量信息
- ✅ 优化进度条和状态指示器
- ✅ 修复UI渲染错位和行覆盖问题

### 🐛 Bug修复

#### 元数据修复
- ✅ **专辑元数据音质标签** - Album和AlbumSort字段现在包含音质标签（如"Head Hunters Hi-Res Lossless"）
  - 修复音乐管理软件无法区分同一专辑不同音质版本的问题
  - 影响所有音质类型：Alac/Hi-Res Lossless/Dolby Atmos/Aac 256

#### UI修复
- ✅ **日志重复问题** - 修复logger输出干扰UI光标定位（重定向logger到stderr）
- ✅ **UI渲染错位** - 修复长曲目名称导致的行覆盖和滚动问题
- ✅ **进度更新优化** - 减少不必要的UI刷新，提升性能

### 🔧 配置变更

#### 新增配置项
```yaml
# 日志配置
logging:
  level: info
  output: stdout
  show_timestamp: false

# 工作-休息循环（仅批量模式生效）
work-rest-enabled: false
work-duration-minutes: 5
rest-duration-minutes: 1
```

### 📚 技术改进

#### 代码质量
- ✅ 系统化替换 fmt.Print 为 logger 调用
- ✅ 移除 UI 直接调用，使用事件驱动架构
- ✅ 完整的单元测试覆盖（logger + progress）
- ✅ 并发安全性验证（go test -race）

#### 性能优化
- ✅ 减少UI刷新频率（200ms → 智能刷新）
- ✅ 优化logger性能（支持高并发）
- ✅ 改进进度事件分发效率

### 🚧 已知问题

无重大已知问题。

### 📝 升级说明

1. **配置文件更新**：请参考 `config.yaml.example` 添加新的 `logging` 配置项
2. **元数据更新**：新下载的文件会自动包含音质标签，旧文件不受影响
3. **UI体验**：动态UI现在更稳定，支持 `--no-ui` 禁用动态UI回退到日志模式

---

## [v2.2.0] - 2025-10-09

### ✨ New Features

#### 🔧 日志与UI治理方案
- **全局OutputMutex + SafePrintf封装** - 解决输出交织和光标错位问题
- **--no-ui 开关** - 支持纯日志输出模式，适合CI/调试环境
- **UI Suspend/Resume API** - 交互式输入时暂停UI更新，避免冲突

#### 🎵 音频与视频增强
- **MV下载Emby兼容路径** - 完全符合Emby/Jellyfin媒体服务器命名规范
- **音质标签统一处理** - `FormatQualityTag()` 统一所有质量标签格式化
- **独立MV保存路径** - 支持 `mv-save-folder` 配置

#### 💡 用户体验优化
- **交互式文件检查** - 检测到文件已存在时询问是否校验
- **智能提示信息** - 区分"文件转移完成"和"校验完成"
- **改进的统计输出** - 简化格式，修复换行符问题

### ⚡ Performance

#### 🚀 缓存中转机制 (v1.1.0+)
- **下载速度提升 50-70%** - 针对NFS等网络文件系统优化
- **网络I/O减少 90%+** - 本地处理后批量传输
- **原子性保证** - 失败自动回滚，不留垃圾文件

**配置示例：**
```yaml
enable-cache: true
cache-folder: "./Cache"
```

### 🐛 Bug Fixes

#### 缓存机制修复
- **修复缓存跳过已下载文件问题** (v1.1.1) - 智能检查最终目标路径
- **修复ffmpeg检测路径错误** (v1.1.2) - 区分工作路径和返回路径

#### 输出优化
- **修复统计输出缺少换行符** - 解决shell提示符显示异常
- **简化统计输出格式** - 移除书名号，提升可读性

### 📚 Documentation

#### 新增文档
- `CACHE_MECHANISM.md` - 完整缓存机制技术文档
- `QUICKSTART_CACHE.md` - 缓存快速开始指南
- `CACHE_UPDATE.md` - 缓存更新说明
- `CHANGELOG.md` - 版本变更日志（本文档）

#### 更新文档
- `README.md` - 全面更新功能说明和使用指南
- `README-CN.md` - 中文文档同步更新
- `config.yaml.example` - 添加新配置项说明

### 🔧 Configuration Changes

#### 新增配置项
```yaml
# 缓存中转机制
enable-cache: true
cache-folder: "./Cache"

# MV保存路径
mv-save-folder: "/media/Music/AppleMusic/MusicVideos"

# 专辑和文件命名格式
album-folder-format: "{AlbumName} {Tag}"
song-file-format: "{SongNumer}. {SongName}"
artist-folder-format: "{ArtistName}"
playlist-folder-format: "{PlaylistName}"
```

### 📊 Performance Benchmarks

#### 测试场景：Hi-Res专辑（12首）

| 指标 | 未启用缓存 | 启用缓存 | 提升 |
|------|-----------|---------|------|
| 下载时间 | 510秒 | 190秒 | **63%** |
| 网络请求 | 450次 | 15次 | **97%** |

---

## [v2.1.0] - 2025-10-09

### ✨ Features
- 缓存中转机制初始实现
- MV下载Emby兼容路径
- 音质标签统一处理机制

### 🐛 Fixes
- 多项性能优化和稳定性改进

---

## [v2.0.0-logger-system] - 2025-10-09

### ✨ Features
- 统一日志系统重构
- 静态展示和固定位置刷新
- 解决并发竞争问题

---

## [v1.x] - Earlier Versions

### Core Features
- ALAC/Dolby Atmos/Hi-Res Lossless 下载
- 内嵌封面和LRC歌词
- 逐词与未同步歌词支持
- 歌手专辑批量下载
- MV下载功能
- 交互式搜索
- FFmpeg自动修复
- 多账号轮换
- TXT批量下载

---

## Upgrade Guide

### From v2.1.0 to v2.2.0

1. **更新配置文件**
   ```bash
   # 备份现有配置
   cp config.yaml config.yaml.backup
   
   # 添加新配置项（可选）
   # 参考 config.yaml.example
   ```

2. **启用缓存机制**（推荐用于NFS）
   ```yaml
   enable-cache: true
   cache-folder: "./Cache"
   ```

3. **使用新的输出模式**
   ```bash
   # 动态UI（默认）
   ./apple-music-downloader <url>
   
   # 纯日志模式
   ./apple-music-downloader --no-ui <url>
   ```

### Breaking Changes

**无破坏性变更** - 所有新功能都是可选的，向后兼容。

---

## Known Issues

### Current Limitations
- 缓存机制需要额外磁盘空间（建议50GB+）
- 动态UI在某些终端上可能显示异常（使用`--no-ui`解决）
- 翻译和发音歌词功能仍在Beta阶段

### Workarounds
- NFS慢速问题 → 启用缓存机制
- UI显示混乱 → 使用 `--no-ui` 标志
- 编码问题 → 启用 `ffmpeg-fix`

---

## Future Plans

### v2.3.0 (计划中)
- [ ] 缓存大小限制和自动清理
- [ ] 缓存统计信息面板
- [ ] 更多音质格式支持

### v3.0.0 (长期)
- [ ] Web UI 界面
- [ ] 分布式缓存支持
- [ ] 插件系统

---

## Credits

- **Sorrow** - Original script author
- **chocomint** - Created `agent-arm64.js`
- **zhaarey** - wrapper decryption service
- **Sendy McSenderson** - Stream decryption code
- All contributors and testers

---

**Note:** This changelog follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format and uses [Semantic Versioning](https://semver.org/).

