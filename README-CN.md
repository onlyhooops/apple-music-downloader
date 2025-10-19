# Apple Music ALAC / 杜比全景声下载器

[English](./README.md) / [简体中文](#)

> [!WARNING]
> **⚠️ 实验性分支警告**
> 
> 该分支仅个人使用的实验分支，增加大量魔改内容，存在大量未知BUG和风险，请谨慎使用！
> 
> 本分支包含未经充分测试的实验性功能，可能导致数据丢失、下载失败或其他不可预见的问题。生产环境请使用官方稳定版本。

**强大的 Apple Music 下载工具，支持 ALAC、Dolby Atmos、Hi-Res Lossless 等高品质音频格式及音乐视频下载。**

原脚本由 Sorrow 编写，经过大量改进和优化。

---

## 🎉 最新更新 (v2.6.0)

### 🚀 v2.6.0 - 架构重构与系统优化 (2025-10-11)

#### 🏗️ 架构重构
- **统一日志系统** - 全新的 `internal/logger` 模块，支持 4 级日志控制（DEBUG/INFO/WARN/ERROR）
- **Progress 事件系统** - 观察者模式的事件架构，解耦 UI 更新与业务逻辑
- **UI 智能简化** - 自适应终端宽度显示（FullMode/CompactMode/MinimalMode）
- **专辑元数据修复** - Album 和 AlbumSort 字段现在包含音质标签（如 "Head Hunters Hi-Res Lossless"）

#### ⚡ 性能优化
- **减少 UI 刷新频率** - 从 200ms 改为智能刷新，降低 CPU 占用
- **优化 logger 性能** - 支持高并发场景的日志输出
- **改进进度事件分发** - 提升多任务并发效率

#### 🐛 关键修复
- **日志重复问题** - 修复 logger 输出干扰 UI 光标定位（重定向 logger 到 stderr）
- **UI 渲染错位** - 修复长曲目名称导致的行覆盖和滚动问题
- **元数据音质标签** - 修复音乐管理软件无法区分同一专辑不同音质版本的问题

### 🎯 v2.5.0 - 元数据与命名规范完善 (2025-10-10)

> [!IMPORTANT]
> **⭐ 核心功能：专辑文件夹命名与元数据音质标签**
> 
> 这是一个重要功能，专门解决音乐管理软件（如 Plex/Emby/Jellyfin）无法区分同一专辑不同音质版本的问题。
> 
> **功能特点：**
> - ✅ **专辑文件夹标签** - 文件夹名称末尾添加音质标签（如 `Head Hunters Alac/`）
> - ✅ **元数据音质标签** - ALBUM 和 ALBUMSORT 字段包含音质标签（如 `ALBUM = "Head Hunters Alac"`）
> - ✅ **完美同步** - 文件夹名称和元数据保持一致，确保完整性
> - ✅ **可配置** - 通过配置文件灵活控制是否启用
> 
> **适用场景：**
> - 📚 收藏同一专辑的多个音质版本（Alac、Hi-Res、Atmos、AAC）
> - 🎵 使用 Plex/Emby/Jellyfin 等音乐服务器
> - 💿 需要精确管理和区分专辑版本
> 
> **配置方法：**
> ```yaml
> # config.yaml
> add-quality-tag-to-folder: true      # 文件夹名称包含音质标签
> add-quality-tag-to-metadata: true    # 元数据包含音质标签
> ```

#### 专辑元数据音质标签
- **Album 和 AlbumSort 字段** - 为元数据中的 ALBUM 和 ALBUMSORT 字段添加音质标签
- **修复识别问题** - 确保音乐管理软件能正确识别不同音质版本
- **全音质支持** - Alac / Hi-Res Lossless / Dolby Atmos / Aac 256
- **兼容性** - iTunes / Plex / Emby / Jellyfin 等完美支持

**效果示例：**
```
专辑文件夹: Head Hunters Hi-Res Lossless/
曲目元数据: ALBUM = "Head Hunters Hi-Res Lossless"
           ALBUMSORT = "Head Hunters Hi-Res Lossless"
```

### 📊 近期重大更新

#### v2.3.0 - MV 下载增强 (2025-10-09)
- **🎬 MV 质量显示** - 自动检测并显示视频质量（4K/1080P/720P/480P）
- **📈 真实进度追踪** - 修复 MV 下载进度显示真实总大小，而非分片大小
- **🎨 UI 优化** - 精简进度条，明确标注视频/音频流

#### v2.2.0 - UI 与日志治理 (2025-10-09)
- **🌏 中文帮助菜单** - `--help` 参数说明完全汉化
- **✨ Emoji 美化** - 美观的终端输出，配合上下文 emoji 图标
- **🔧 线程安全日志** - OutputMutex + SafePrintf 实现清晰并发日志

#### v2.1.0 - 性能与用户体验改进 (2025-10-09)
- **⚡ 缓存中转机制** - NFS/网络存储下载速度提升 50-70%
- **🔍 交互式文件检查** - 智能提示已存在文件，提供跳过选项
- **🎯 质量标签标准化** - Emby 兼容 MV 路径和统一质量标签

### 📈 改进总结
- **代码质量**：全面重构，新增单元测试，并发安全性验证（go test -race）
- **用户体验**：Emoji 丰富输出、中文本地化、智能 UI 简化、更清晰的进度指示
- **性能提升**：并发下载、智能缓存、减少网络开销、优化 UI 刷新
- **文档完善**：15 个技术文档、完整变更日志、详细配置说明

---

## ✨ 核心特性

### 🎵 音频格式支持
- **ALAC（无损音频）** - `audio-alac-stereo`
- **Dolby Atmos / EC3** - `audio-atmos` / `audio-ec3`
- **Hi-Res Lossless（高解析度无损）** - 最高 24-bit/192kHz
- **AAC 格式** - `audio-stereo`、`audio-stereo-binaural`、`audio-stereo-downmix`
- **AAC-LC** - `audio-stereo`（需要订阅）

### 📹 音乐视频支持
- **4K/1080p/720p** 分辨率选项
- **Emby/Jellyfin 兼容**命名结构
- **多音轨支持**（Atmos/AC3/AAC）
- 独立保存路径配置

### 🎼 丰富的元数据和歌词
- **内嵌封面**和艺术作品（最高 5000x5000）
- **同步 LRC 歌词**（逐词 / 逐音节）
- **翻译和发音歌词**支持（Beta）
- **动态封面**（Emby/Jellyfin 支持）
- 完整的元数据标签

### ⚡ 性能优化
- **缓存中转机制** - NFS/网络存储下载速度提升 50-70%
- **并行下载** - 多线程分片下载
- **智能文件检查** - 跳过已下载文件
- **批量下载** - 从 TXT 文件读取链接并可配置线程数

### 🛠️ 高级功能
- **多账号轮换** - 根据区域自动选择账号
- **FFmpeg 自动修复** - 检测并修复编码问题
- **交互模式** - 箭头键导航搜索结果
- **歌手下载** - 下载歌手页面的所有专辑/MV
- **自定义命名** - 灵活的文件夹和文件命名格式
- **输出模式** - 动态 UI 或纯日志模式（`--no-ui`）

---

## 📋 环境要求

### 必需依赖

1. **[MP4Box](https://gpac.io/downloads/gpac-nightly-builds/)** - 必须安装并添加到系统 PATH
2. **[mp4decrypt](https://www.bento4.com/downloads/)** - 下载音乐视频时需要
3. **FFmpeg**（可选）- 用于动态封面和自动修复功能

### 系统要求

### 系统要求
- Go 1.23.1 或更高版本
- 建议 8GB+ 内存
- 50GB+ 可用磁盘空间（如果使用缓存机制）

> [!NOTE]
> **💡 磁盘空间建议**
> 
> - **不使用缓存**：只需足够存储下载文件的空间
> - **使用缓存机制**：额外需要 50GB+ 本地临时空间
> - **大规模批量下载**：建议预留 100GB+ 空间以获得最佳性能

---

## 🚀 快速开始

### 1. 安装

```bash
# 克隆仓库
git clone https://github.com/your-repo/apple-music-downloader.git
cd apple-music-downloader

# 安装依赖
go mod tidy

# 编译二进制文件
go build -o apple-music-downloader main.go
```

### 2. 配置

```bash
# 复制示例配置
cp config.yaml.example config.yaml

# 编辑配置文件
nano config.yaml
```

**获取 `media-user-token`：**
1. 打开 [Apple Music](https://music.apple.com) 并登录
2. 按 `F12` 打开开发者工具
3. 进入 `Application`（应用程序）→ `Cookies` → `https://music.apple.com`
4. 找到名为 `media-user-token` 的 Cookie 并复制其值
5. 粘贴到 `config.yaml` 中

### 3. 基本使用

```bash
# 下载专辑
./apple-music-downloader https://music.apple.com/cn/album/专辑名/123456789

# 下载杜比全景声
./apple-music-downloader --atmos https://music.apple.com/cn/album/专辑名/123456789

# 下载单曲
./apple-music-downloader --song https://music.apple.com/cn/album/专辑/123?i=456

# 下载播放列表
./apple-music-downloader https://music.apple.com/cn/playlist/歌单名/pl.xxxxx

# 下载歌手的所有内容
./apple-music-downloader https://music.apple.com/cn/artist/歌手名/123456

# 交互式搜索
./apple-music-downloader --search song "搜索词"
./apple-music-downloader --search album "专辑名"
./apple-music-downloader --search artist "歌手名"

# 从 TXT 文件批量下载
./apple-music-downloader urls.txt

# 纯日志模式（用于 CI/调试）
./apple-music-downloader --no-ui https://music.apple.com/...
```

---

## 📖 高级用法

### 缓存机制（NFS 优化）

> [!IMPORTANT]
> **⚠️ 缓存机制重要说明**
> 
> **建议启用场景：**
> - ✅ 目标路径为 NFS/SMB 等网络文件系统
> - ✅ 本地有 50GB+ 可用磁盘空间（建议 SSD）
> - ✅ 需要频繁进行批量下载任务
> 
> **关键注意事项：**
> - ⚠️ **磁盘空间**：缓存文件夹需要足够的临时存储空间，建议至少 50GB
> - ⚠️ **缓存路径**：务必使用本地快速磁盘（SSD），不要设置在 NFS 等网络路径上
> - ⚠️ **文件系统**：跨文件系统传输会使用拷贝方式，速度会有所降低
> - ⚠️ **清理机制**：程序会自动清理成功传输的缓存，失败时也会自动回滚
> - ⚠️ **手动清理**：可以随时手动删除 `Cache` 文件夹，程序会自动重建
> 
> **性能提升数据（实测）：**
> - 下载时间提升：**50-70%**
> - 网络 I/O 减少：**90%+**
> - 更好的稳定性：原子操作，失败自动回滚

当下载到网络存储（NFS/SMB）时，可显著提升性能：

```yaml
# config.yaml
enable-cache: true
cache-folder: "./Cache"  # 建议使用本地 SSD 路径
```

**配置建议：**
- ⚡ **本地 SSD 缓存** - 将 `cache-folder` 设置为本地 SSD 路径（如 `/ssd/cache/apple-music`）
- ⚡ **网络存储目标** - `alac-save-folder` 和 `atmos-save-folder` 设置为 NFS/SMB 路径
- ⚡ **充足空间** - 确保缓存路径有至少 50GB 可用空间

**工作原理：**
1. 文件先下载到本地缓存文件夹
2. 所有处理（解密、合并、元数据）在本地完成
3. 完成后一次性批量传输到目标网络路径
4. 自动清理缓存，释放空间

[📚 阅读缓存机制完整文档](./CACHE_MECHANISM.md)

### 历史记录与断点续传（v2.7.0+ 增强版）

> [!TIP]
> **🔄 全局链接级历史记录**
> 
> 所有下载任务（单链接、多链接、文本批量）现在共享全局历史记录系统。任何已成功下载的链接都会被自动跳过，无论它来自哪个任务或文件。

**核心特性：**
- 🌐 **全局去重**：历史记录以链接为单位，而非任务文件
- 📁 **自动记录**：所有模式自动记录到 `history` 文件夹
- 🔍 **跨任务检测**：自动跳过在之前任务中已下载的链接
- 🎵 **音质感知**：检测音质参数变化，必要时重新下载
- ⏸️ **智能续传**：支持任务中断后从断点继续

**使用示例：**
```bash
# 场景1: 下载经典专辑列表
./apple-music-downloader 经典专辑.txt
# 完成50个专辑的下载

# 场景2: 下载爵士音乐列表（其中10个与经典专辑重复）
./apple-music-downloader 爵士音乐.txt
# 输出: 📜 全局历史记录检测: 发现 10 个已完成的任务
#       ⏭️  已自动跳过 10 个，剩余 40 个任务

# 场景3: 单独下载某个专辑（该专辑已在之前的列表中）
./apple-music-downloader https://music.apple.com/cn/album/xxx/123
# 输出: ✅ 所有任务都已完成，无需重复下载！
```

**支持的模式：**
- ✅ 单链接模式
- ✅ 多链接模式
- ✅ TXT文件批量模式
- ✅ 混合模式（URL + TXT）
- ✅ 交互模式

**高级用法：**
```bash
# 查看所有历史记录
ls -lh history/

# 清空所有历史记录（重新开始）
rm -rf history/

# 清空特定任务的历史记录
rm history/经典专辑.txt_*.json
```

[📚 阅读历史记录功能完整文档](./HISTORY_FEATURE.md)

### 日志配置（v2.6.0+）

**统一日志系统**，支持 4 级日志控制和灵活配置：

```yaml
# config.yaml
logging:
  level: info                  # debug/info/warn/error
  output: stdout               # stdout/stderr/文件路径
  show_timestamp: false        # UI模式下建议关闭
```

**日志级别说明：**
- `debug` - 显示所有调试信息（用于开发和故障排查）
- `info` - 显示常规信息（默认，推荐）
- `warn` - 仅显示警告和错误
- `error` - 仅显示错误信息

**输出目标：**
- `stdout` - 标准输出（默认）
- `stderr` - 标准错误输出（UI 模式下自动使用）
- 文件路径 - 如 `./logs/download.log`

**使用建议：**
- 动态 UI 模式：`show_timestamp: false`，避免时间戳干扰 UI
- 纯日志模式（`--no-ui`）：`show_timestamp: true`，便于追溯
- CI/CD 环境：使用 `--no-ui` + 日志文件输出

### 自定义命名格式

> [!TIP]
> **🏷️ 音质标签配置 (v2.5.0+)**
> 
> 从 v2.5.0 开始，您可以灵活控制音质标签的显示位置：
> 
> ```yaml
> # config.yaml - 音质标签配置
> add-quality-tag-to-folder: true      # 文件夹名称包含音质标签
> add-quality-tag-to-metadata: true    # 元数据包含音质标签
> ```
> 
> **配置组合效果：**
> 
> | 文件夹标签 | 元数据标签 | 文件夹名称 | 元数据 ALBUM | 适用场景 |
> |:---:|:---:|---|---|---|
> | ✅ | ✅ | `Head Hunters Alac/` | `Head Hunters Alac` | **推荐**：完美同步，音乐软件能正确识别 |
> | ✅ | ❌ | `Head Hunters Alac/` | `Head Hunters` | 文件分类明确，元数据简洁 |
> | ❌ | ✅ | `Head Hunters/` | `Head Hunters Alac` | 文件夹简洁，音质信息在元数据中 |
> | ❌ | ❌ | `Head Hunters/` | `Head Hunters` | 不推荐：无法区分不同音质版本 |
> 
> **使用建议：**
> - 🎵 **Plex/Emby/Jellyfin 用户**：两项都启用（`true`）
> - 💿 **收藏多音质版本**：两项都启用（`true`）
> - 🗂️ **仅需文件分类**：仅启用文件夹标签
> - ✨ **追求简洁**：仅启用元数据标签

```yaml
# 专辑文件夹："专辑名 Dolby Atmos"
album-folder-format: "{AlbumName} {Tag}"

# 歌曲文件："01. 歌曲名"
song-file-format: "{SongNumer}. {SongName}"

# 歌手文件夹："歌手名"
artist-folder-format: "{ArtistName}"

# 播放列表文件夹："播放列表名"
playlist-folder-format: "{PlaylistName}"
```

**可用变量：**
- 专辑：`{AlbumId}`、`{AlbumName}`、`{ArtistName}`、`{ReleaseDate}`、`{ReleaseYear}`、`{Tag}`、`{Quality}`、`{Codec}`、`{UPC}`、`{Copyright}`、`{RecordLabel}`
- 歌曲：`{SongId}`、`{SongNumer}`、`{SongName}`、`{DiscNumber}`、`{TrackNumber}`、`{Tag}`、`{Quality}`、`{Codec}`
- 播放列表：`{PlaylistId}`、`{PlaylistName}`、`{ArtistName}`、`{Tag}`、`{Quality}`、`{Codec}`
- 歌手：`{ArtistId}`、`{ArtistName}`、`{UrlArtistName}`

### 多账号配置

```yaml
accounts:
  - name: "CN"
    storefront: "cn"
    media-user-token: "你的中国区token"
    decrypt-m3u8-port: "127.0.0.1:10020"
    get-m3u8-port: "127.0.0.1:10021"
    
  - name: "US"
    storefront: "us"
    media-user-token: "你的美国区token"
    decrypt-m3u8-port: "127.0.0.1:20020"
    get-m3u8-port: "127.0.0.1:20021"
```

程序会根据 URL 的区域（如 `/cn/`、`/us/`）自动选择对应的账号。

### 翻译和发音歌词（Beta）

1. 打开 [Apple Music Beta](https://beta.music.apple.com) 并登录
2. 按 `F12` → `Network`（网络）标签页
3. 搜索并播放一首 K-Pop 歌曲（或任何有翻译的歌曲）
4. 点击歌词按钮
5. 在网络标签页中找到 `syllable-lyrics` 请求
6. 从 URL 中复制 `l=` 参数的值
7. 粘贴到 `config.yaml`：

```yaml
language: "zh-CN%2Cko-KR%5Bttml%3Aruby%5D"
```

---

## 🔧 命令行选项

| 选项 | 说明 |
|------|------|
| `--atmos` | 下载杜比全景声格式 |
| `--aac` | 下载 AAC 256 格式 |
| `--song` | 下载单曲 |
| `--select` | 交互式选择曲目 |
| `--search [类型] "关键词"` | 搜索（song/album/artist） |
| `--debug` | 显示可用音质信息 |
| `--no-ui` | 禁用动态 UI，纯日志输出 |
| `--config 路径` | 指定自定义配置文件 |
| `--output 路径` | 覆盖保存文件夹 |

---

## 📂 输出结构

### 专辑（Emby 兼容命名）

```
/media/Music/AppleMusic/Alac/
└── Taylor Swift/
    └── 1989 (Taylor's Version) Hi-Res Lossless/
        ├── cover.jpg
        ├── 01. Welcome To New York.m4a
        ├── 02. Blank Space.m4a
        └── ...
```

### 音乐视频（Emby/Jellyfin 兼容）

```
/media/Music/AppleMusic/MusicVideos/
└── Morgan James/
    └── Thunderstruck (2024)/
        └── Thunderstruck (2024).mp4
```

---

## 🐛 故障排除

### 常见问题

**1. "MP4Box not found"（找不到 MP4Box）**
- 安装 [MP4Box](https://gpac.io/downloads/gpac-nightly-builds/)
- 确保已添加到系统 PATH
- 测试：`MP4Box -version`

**2. "No media-user-token"（没有 media-user-token）**
- AAC-LC、MV 和歌词需要有效的订阅 token
- ALAC/Dolby Atmos 使用基本 token 即可

**3. UI 输出混乱**
- 使用 `--no-ui` 标志获得纯日志输出
- 更适合 CI/CD 流程或输出重定向

**4. NFS 下载缓慢**
- 在 config.yaml 中启用缓存机制
- 参阅[缓存快速入门指南](./QUICKSTART_CACHE.md)

### FFmpeg 自动修复

如果下载的文件有编码问题：

```yaml
ffmpeg-fix: true  # 下载完成后自动检测
```

程序将会：
1. 检测损坏/不完整的文件
2. 提示确认
3. 使用 FFmpeg 和 ALAC 编解码器重新编码

---

## 📊 性能建议

### 对于网络存储（NFS/SMB）
- ✅ 启用缓存机制
- ✅ 使用本地 SSD 作为缓存文件夹
- ✅ 增加分片下载线程数

### 对于批量下载
```yaml
txtDownloadThreads: 5  # 并行专辑下载数
chunk_downloadthreads: 30  # 并行分片下载数
```

### 对于大型音乐库
- ✅ 启用 `ffmpeg-fix` 进行质量保证
- ✅ 使用 `--no-ui` 获得更清晰的日志
- ✅ 保存输出到文件：`./app --no-ui url > download.log 2>&1`

---

## 📚 文档

### 用户指南
- [README.md](./README.md) - English documentation
- [QUICKSTART_CACHE.md](./QUICKSTART_CACHE.md) - 缓存机制快速开始
- [CACHE_UPDATE.md](./CACHE_UPDATE.md) - 缓存更新指南
- [GOO_ALIAS.md](./GOO_ALIAS.md) - 命令别名配置指南
- [EMOJI_DEMO.md](./EMOJI_DEMO.md) - Emoji 输出演示

### 技术文档
- [CHANGELOG.md](./CHANGELOG.md) - 完整版本历史和变更记录
- [CACHE_MECHANISM.md](./CACHE_MECHANISM.md) - 完整缓存技术文档
- [MV_QUALITY_DISPLAY.md](./MV_QUALITY_DISPLAY.md) - MV 质量检测功能
- [MV_PROGRESS_FIX.md](./MV_PROGRESS_FIX.md) - MV 进度追踪改进
- [MV_LOG_FIX.md](./MV_LOG_FIX.md) - MV 下载日志增强

---

## 🙏 致谢与鸣谢

### 🎖️ 原作者与核心贡献者
- **Sorrow** - 原始脚本作者和架构设计
- **chocomint** - 创建了 `agent-arm64.js` ARM 支持
- **Sendy McSenderson** - 流解密代码

### 🔧 上游依赖与工具
- **[mp4ff](https://github.com/Eyevinn/mp4ff)** by Eyevinn - MP4 文件处理
- **[mp4ff (fork)](https://github.com/itouakirai/mp4ff)** by itouakirai - 增强 MP4 支持
- **[progressbar/v3](https://github.com/schollz/progressbar)** by schollz - 进度显示
- **[requests](https://github.com/sky8282/requests)** by sky8282 - HTTP 客户端封装
- **[m3u8](https://github.com/grafov/m3u8)** by grafov - M3U8 播放列表解析
- **[pflag](https://github.com/spf13/pflag)** by spf13 - 命令行参数
- **[tablewriter](https://github.com/olekukonko/tablewriter)** by olekukonko - 表格格式化
- **[color](https://github.com/fatih/color)** by fatih - 彩色终端输出

### 🛠️ 外部工具
- **[FFmpeg](https://ffmpeg.org/)** - 音视频处理
- **[MP4Box](https://gpac.io/)** - GPAC 多媒体框架
- **[mp4decrypt](https://www.bento4.com/)** - Bento4 解密工具

### 💝 特别感谢
- **[@sky8282](https://github.com/sky8282)** - 提供优秀的 requests 库和持续支持
- 所有帮助改进本项目的贡献者和测试者
- Apple Music API 研究者和逆向工程社区
- 提供各种库和工具的开源社区

---

## ⚠️ 免责声明

此工具仅用于教育和个人用途。请遵守版权法和 Apple Music 服务条款。请勿分发下载的内容。

---

## 📝 许可证

本项目仅供个人使用。所有下载内容的权利归其各自所有者所有。

---

## 🔗 资源

- [Apple Music for Artists](https://artists.apple.com/)
- [Emby 命名规范](https://emby.media/support/articles/Movie-Naming.html)
- [FFmpeg 文档](https://ffmpeg.org/documentation.html)
- [中文教程](https://telegra.ph/Apple-Music-Alac高解析度无损音乐下载教程-04-02-2)

---

**版本：** v2.6.0  
**最后更新：** 2025-10-11  
**需要 Go 版本：** 1.23.1+
