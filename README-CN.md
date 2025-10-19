# Apple Music 下载器

[English](./README.md) | [简体中文](#)

> **一个功能强大的 Apple Music 高品质音频下载工具**
>
> 支持 ALAC、Hi-Res Lossless、Dolby Atmos 等多种无损格式，以及音乐视频下载。

[![版本](https://img.shields.io/badge/版本-v1.0.0-blue.svg)](https://github.com/onlyhooops/apple-music-downloader)
[![Go版本](https://img.shields.io/badge/Go-1.23.1+-00ADD8.svg)](https://golang.org/)
[![许可证](https://img.shields.io/badge/许可证-个人使用-green.svg)](./LICENSE)

---

## 📖 目录

- [核心特性](#-核心特性)
- [系统要求](#-系统要求)
- [快速开始](#-快速开始)
- [使用指南](#-使用指南)
- [高级功能](#-高级功能)
- [配置说明](#-配置说明)
- [常见问题](#-常见问题)
- [致谢](#-致谢)

---

## ✨ 核心特性

### 🎵 多格式音频支持

#### 无损音频格式
- **ALAC (Apple Lossless)** - 原生无损音频，完美音质
- **Hi-Res Lossless** - 高解析度无损音频，最高支持 24-bit/192kHz
- **Dolby Atmos** - 杜比全景声，沉浸式3D音频体验
- **Dolby EC-3** - 杜比数字增强音频

#### 有损音频格式
- **AAC 256kbps** - 高品质 AAC-LC 编码
- **AAC Binaural** - 双声道空间音频
- **AAC Downmix** - 降混音频流
- **AAC Stereo** - 标准立体声

### 📹 音乐视频下载

- **多分辨率支持** - 4K (2160p) / 1080p / 720p / 480p
- **多音轨选择** - Atmos / AC3 / AAC 音轨
- **媒体服务器兼容** - Emby / Jellyfin / Plex 命名规范
- **独立保存路径** - 可自定义 MV 存储位置

### 🎼 完整的元数据支持

- **专辑封面** - 高达 5000x5000 分辨率
- **歌词嵌入** - 同步 LRC 歌词（逐词/逐音节）
- **翻译歌词** - 支持多语言翻译歌词（Beta）
- **动态封面** - 支持动画封面图（需要 FFmpeg）
- **完整标签** - 艺术家、专辑、曲目编号、发行日期等

### ⚡ 性能优化

- **缓存机制** - 针对 NFS/SMB 网络存储优化，下载速度提升 50-70%
- **并行下载** - 多线程分片下载，充分利用带宽
- **智能断点续传** - 支持下载中断后继续
- **批量处理** - 从 TXT 文件批量下载链接

### 🛠️ 高级功能

- **全局历史记录** - 自动跟踪已下载内容，避免重复下载
- **多账号管理** - 支持多区域账号配置，自动选择
- **交互式搜索** - 内置搜索功能，支持歌曲/专辑/艺术家搜索
- **自定义命名** - 灵活的文件夹和文件命名格式
- **音质标签** - 可选择是否在文件名和元数据中添加音质标签
- **FFmpeg 修复** - 自动检测并修复音频编码问题
- **动态 UI** - 实时显示下载进度，支持多任务并发
- **纯日志模式** - 适用于 CI/CD 环境的 `--no-ui` 模式

---

## 📋 系统要求

### 必需依赖

1. **[Go 1.23.1+](https://golang.org/dl/)** - 编译和运行环境
2. **[MP4Box](https://gpac.io/downloads/gpac-nightly-builds/)** - 音频文件处理（必需）
3. **[mp4decrypt](https://www.bento4.com/downloads/)** - 音乐视频解密（下载 MV 时需要）
4. **[FFmpeg](https://ffmpeg.org/)** - 动态封面和自动修复功能（可选）

### 系统配置建议

- **操作系统**: Linux / macOS / Windows
- **内存**: 建议 8GB 以上
- **磁盘空间**: 
  - 不使用缓存：只需下载文件的存储空间
  - 使用缓存机制：额外需要 50GB+ 本地临时空间
  - 大规模批量下载：建议预留 100GB+ 空间

### Apple Music 账号

- **基础功能**: 免费账号即可下载 ALAC 和 Dolby Atmos
- **高级功能**: AAC 256、MV 下载和歌词需要有效的订阅

---

## 🚀 快速开始

### 1. 安装

```bash
# 克隆仓库
git clone https://github.com/onlyhooops/apple-music-downloader.git
cd apple-music-downloader

# 安装 Go 依赖
go mod tidy

# 编译程序
go build -o apple-music-downloader main.go
```

### 2. 配置

#### 创建配置文件

```bash
# 复制配置文件模板
cp config.yaml.example config.yaml
cp dev.env.example dev.env

# 编辑配置文件
nano config.yaml
nano dev.env
```

#### 获取 Apple Music Token

1. 打开 [Apple Music 网页版](https://music.apple.com) 并登录
2. 按 **F12** 打开开发者工具
3. 选择 **Application**（应用程序）标签
4. 在左侧导航中找到 **Cookies** → `https://music.apple.com`
5. 找到名为 `media-user-token` 的 Cookie
6. 复制其值，粘贴到 `dev.env` 文件中：

```bash
# dev.env
APPLE_MUSIC_MEDIA_USER_TOKEN_CN=你的token值
```

#### 配置保存路径

编辑 `config.yaml`：

```yaml
# 保存路径配置
alac-save-folder: "/media/Music/AppleMusic/Alac"
atmos-save-folder: "/media/Music/AppleMusic/Atmos"
mv-save-folder: "/media/Music/AppleMusic/MusicVideos"
```

### 3. 基本使用

```bash
# 下载专辑（默认 ALAC 无损）
./apple-music-downloader https://music.apple.com/cn/album/专辑名/123456789

# 下载杜比全景声
./apple-music-downloader --atmos https://music.apple.com/cn/album/专辑名/123456789

# 下载 AAC 256 格式
./apple-music-downloader --aac https://music.apple.com/cn/album/专辑名/123456789

# 下载单曲
./apple-music-downloader --song https://music.apple.com/cn/album/专辑/123?i=456

# 下载播放列表
./apple-music-downloader https://music.apple.com/cn/playlist/歌单名/pl.xxxxx

# 下载艺术家的所有作品
./apple-music-downloader https://music.apple.com/cn/artist/歌手名/123456
```

---

## 📖 使用指南

### 交互式搜索

```bash
# 搜索歌曲
./apple-music-downloader --search song "歌曲名"

# 搜索专辑
./apple-music-downloader --search album "专辑名"

# 搜索艺术家
./apple-music-downloader --search artist "歌手名"
```

搜索后会显示结果列表，可使用箭头键选择，Enter 键确认下载。

### 批量下载

创建 `urls.txt` 文件，每行一个链接：

```text
https://music.apple.com/cn/album/专辑1/123456789
https://music.apple.com/cn/album/专辑2/987654321
https://music.apple.com/cn/playlist/歌单/pl.xxxxx
# 这是注释行，会被忽略
```

执行批量下载：

```bash
./apple-music-downloader urls.txt
```

### 命令行参数

| 参数 | 说明 |
|------|------|
| `--atmos` | 下载杜比全景声格式 |
| `--aac` | 下载 AAC 256 格式 |
| `--aac-type <类型>` | 指定 AAC 类型：`aac-lc`、`binaural`、`downmix` |
| `--alac-max <采样率>` | 指定 ALAC 最大采样率：`192000`、`96000`、`48000` |
| `--atmos-max <码率>` | 指定 Atmos 最大码率：`2768`、`2448` |
| `--song` | 下载单曲模式 |
| `--select` | 交互式选择曲目 |
| `--all-album` | 下载艺术家的所有专辑 |
| `--search <类型> "关键词"` | 搜索：`song`、`album`、`artist` |
| `--mv-max <分辨率>` | MV 最大分辨率：`2160`、`1080`、`720` |
| `--mv-audio-type <类型>` | MV 音轨类型：`atmos`、`ac3`、`aac` |
| `--debug` | 显示可用音质信息（不下载） |
| `--no-ui` | 禁用动态 UI，纯日志输出 |
| `--config <路径>` | 指定配置文件路径 |
| `--output <路径>` | 指定本次任务的输出目录 |
| `--start <编号>` | 从 TXT 文件的第几个链接开始（用于断点续传） |

---

## 🎯 高级功能

### 缓存机制（NFS/SMB 优化）

当下载到网络存储（NFS/SMB）时，启用缓存机制可显著提升性能：

#### 配置方法

编辑 `config.yaml`：

```yaml
# 启用缓存
enable-cache: true
cache-folder: "./Cache"  # 建议使用本地 SSD 路径
```

#### 工作原理

1. 文件先下载到本地缓存文件夹（高速磁盘）
2. 所有处理（解密、合并、元数据）在本地完成
3. 完成后批量传输到目标网络路径
4. 自动清理缓存，释放空间

#### 性能提升

- 下载时间：提升 **50-70%**
- 网络 I/O：减少 **90%+**
- 稳定性：原子操作，失败自动回滚

#### 注意事项

- ⚠️ 缓存文件夹必须在本地快速磁盘（SSD）上
- ⚠️ 不要将缓存设置在 NFS 等网络路径
- ⚠️ 确保至少有 50GB 可用空间
- ⚠️ 程序会自动清理，也可手动删除 `Cache` 文件夹

### 全局历史记录系统

#### 核心特性

- **全局去重** - 所有下载任务共享历史记录
- **自动跳过** - 已下载的内容自动跳过
- **音质感知** - 检测音质参数变化，必要时重新下载
- **断点续传** - 支持任务中断后继续

#### 使用示例

```bash
# 第一次：下载经典专辑列表
./apple-music-downloader classic_albums.txt
# 完成 50 个专辑

# 第二次：下载爵士音乐列表（有 10 个重复）
./apple-music-downloader jazz_music.txt
# 输出：📜 全局历史记录检测: 发现 10 个已完成的任务
#       ⏭️  已自动跳过 10 个，剩余 40 个任务

# 第三次：单独下载某个专辑（已在列表中）
./apple-music-downloader https://music.apple.com/cn/album/xxx/123
# 输出：✅ 所有任务都已完成，无需重复下载！
```

#### 管理历史记录

```bash
# 查看所有历史记录
ls -lh history/

# 清空所有历史（重新开始）
rm -rf history/

# 清空特定任务的历史
rm history/经典专辑.txt_*.json
```

### 音质标签配置

从 v1.0.0 开始，可以灵活控制音质标签的显示：

```yaml
# config.yaml
add-quality-tag-to-folder: true      # 文件夹名称包含音质标签
add-quality-tag-to-metadata: true    # 元数据包含音质标签
```

#### 配置效果对比

| 文件夹标签 | 元数据标签 | 文件夹名称 | 元数据 ALBUM | 适用场景 |
|:---:|:---:|---|---|---|
| ✅ | ✅ | `Head Hunters Alac/` | `Head Hunters Alac` | **推荐**：完美同步，音乐软件正确识别 |
| ✅ | ❌ | `Head Hunters Alac/` | `Head Hunters` | 文件分类明确，元数据简洁 |
| ❌ | ✅ | `Head Hunters/` | `Head Hunters Alac` | 文件夹简洁，音质在元数据中 |
| ❌ | ❌ | `Head Hunters/` | `Head Hunters` | 不推荐：无法区分不同音质版本 |

#### 使用建议

- 🎵 **Plex/Emby/Jellyfin 用户**：两项都启用
- 💿 **收藏多音质版本**：两项都启用
- 🗂️ **仅需文件分类**：仅启用文件夹标签
- ✨ **追求简洁**：仅启用元数据标签

### 自定义命名格式

#### 可用变量

**专辑相关**：
- `{AlbumId}` - 专辑 ID
- `{AlbumName}` - 专辑名称
- `{ArtistName}` - 艺术家名称
- `{ReleaseDate}` - 发行日期
- `{ReleaseYear}` - 发行年份
- `{Tag}` - 音质标签（如 "Alac"、"Hi-Res Lossless"）
- `{Quality}` - 音质描述
- `{Codec}` - 编码格式
- `{UPC}` - UPC 码
- `{Copyright}` - 版权信息
- `{RecordLabel}` - 唱片公司

**曲目相关**：
- `{SongId}` - 曲目 ID
- `{SongName}` - 曲目名称
- `{SongNumer}` - 曲目编号（两位数）
- `{TrackNumber}` - 曲目编号（原始）
- `{DiscNumber}` - 碟片编号

**艺术家相关**：
- `{ArtistId}` - 艺术家 ID
- `{ArtistName}` - 艺术家名称
- `{UrlArtistName}` - URL 中的艺术家名称

#### 命名示例

```yaml
# config.yaml

# 专辑文件夹："{专辑名} {音质标签}"
album-folder-format: "{AlbumName} {Tag}"
# 结果：Head Hunters Hi-Res Lossless/

# 曲目文件："{编号}. {曲名}"
song-file-format: "{SongNumer}. {SongName}"
# 结果：01. Chameleon.m4a

# 艺术家文件夹："{艺术家名}"
artist-folder-format: "{ArtistName}"
# 结果：Herbie Hancock/

# 播放列表文件夹："{列表名} {音质标签}"
playlist-folder-format: "{PlaylistName} {Tag}"
# 结果：Jazz Classics Alac/
```

### 多账号配置

支持配置多个区域的 Apple Music 账号，程序会根据链接的区域自动选择：

```yaml
# config.yaml
accounts:
  - name: "CN"
    storefront: "cn"
    media-user-token: "${APPLE_MUSIC_MEDIA_USER_TOKEN_CN}"
    decrypt-m3u8-port: "127.0.0.1:10020"
    get-m3u8-port: "127.0.0.1:10021"

  - name: "US"
    storefront: "us"
    media-user-token: "${APPLE_MUSIC_MEDIA_USER_TOKEN_US}"
    decrypt-m3u8-port: "127.0.0.1:20020"
    get-m3u8-port: "127.0.0.1:20021"
```

### 日志配置

```yaml
# config.yaml
logging:
  level: info                  # 日志级别: debug/info/warn/error
  output: stdout               # 输出目标: stdout/stderr/文件路径
  show_timestamp: false        # UI 模式建议关闭时间戳
```

**日志级别说明**：
- `debug` - 显示所有调试信息（开发和故障排查）
- `info` - 显示常规信息（默认，推荐）
- `warn` - 仅显示警告和错误
- `error` - 仅显示错误信息

**使用建议**：
- 动态 UI 模式：`show_timestamp: false`
- 纯日志模式（`--no-ui`）：`show_timestamp: true`
- CI/CD 环境：使用 `--no-ui` + 日志文件输出

---

## ⚙️ 配置说明

### 性能调优

#### 网络存储（NFS/SMB）

```yaml
enable-cache: true
cache-folder: "/ssd/cache/apple-music"  # 使用本地 SSD
chunk_downloadthreads: 30
```

#### 批量下载优化

```yaml
batch-size: 20                           # 每批处理数量
skip-existing-validation: true           # 自动跳过已存在文件
work-rest-enabled: true                  # 工作-休息循环
work-duration-minutes: 30                # 工作 30 分钟
rest-duration-minutes: 2                 # 休息 2 分钟
```

#### 下载线程配置

```yaml
# M3U8 切片下载
chunk_downloadthreads: 30                # 音频切片线程
mv_chunk_downloadthreads: 30             # 视频切片线程

# 音频格式线程
aac_downloadthreads: 5                   # AAC 格式
lossless_downloadthreads: 5              # 无损格式
hires_downloadthreads: 5                 # Hi-Res 格式

# MV 下载
mv_downloadthreads: 3                    # MV 并行下载数
```

### FFmpeg 配置

```yaml
ffmpeg-fix: true                         # 自动检测并修复
ffmpeg-check-args: "-map 0:a:0 -f wav -hide_banner -loglevel error -"
ffmpeg-encode-args: "-c:v copy -c:a alac -avoid_negative_ts make_zero -f mp4 -y"
```

---

## 📂 输出结构

### 专辑结构（Emby 兼容）

```
/media/Music/AppleMusic/Alac/
└── Taylor Swift/
    └── 1989 (Taylor's Version) Hi-Res Lossless/
        ├── cover.jpg
        ├── 01. Welcome To New York.m4a
        ├── 02. Blank Space.m4a
        ├── 03. Style.m4a
        └── ...
```

### 音乐视频结构（Emby/Jellyfin 兼容）

```
/media/Music/AppleMusic/MusicVideos/
└── Morgan James/
    └── Thunderstruck (2024)/
        └── Thunderstruck (2024).mp4
```

---

## 🐛 常见问题

### 1. "MP4Box not found"

**原因**：未安装 MP4Box 或未添加到系统 PATH

**解决方法**：
```bash
# 安装 MP4Box
# Linux (Ubuntu/Debian):
sudo apt-get install gpac

# macOS:
brew install gpac

# 验证安装
MP4Box -version
```

### 2. "No media-user-token"

**原因**：
- AAC-LC、MV 和歌词功能需要有效的订阅 token
- ALAC 和 Dolby Atmos 可使用基础 token

**解决方法**：确保正确配置 token（参见"获取 Apple Music Token"章节）

### 3. UI 输出混乱

**原因**：终端不支持动态更新或重定向输出

**解决方法**：
```bash
# 使用纯日志模式
./apple-music-downloader --no-ui <链接>

# 保存日志到文件
./apple-music-downloader --no-ui <链接> > download.log 2>&1
```

### 4. NFS 下载速度慢

**原因**：网络文件系统的延迟和频繁的小文件写入

**解决方法**：启用缓存机制（参见"缓存机制"章节）

### 5. 下载中断后如何继续

**方法一**：全局历史记录会自动跳过已完成的任务
```bash
# 重新运行相同命令即可
./apple-music-downloader urls.txt
```

**方法二**：使用 `--start` 参数
```bash
# 从第 44 个链接开始
./apple-music-downloader --start 44 urls.txt
```

### 6. 如何清理缓存

```bash
# 手动删除缓存文件夹
rm -rf ./Cache

# 程序会在下次运行时自动重建
```

---

## 📊 性能参考

### 测试环境

- **服务器**: Proxmox VE 6.8.12
- **CPU**: 8 Core @ 2.4GHz
- **内存**: 16GB
- **存储**: NFS 网络存储
- **网络**: 1Gbps

### 性能数据

#### 不启用缓存（直接写入 NFS）

| 项目 | 数据 |
|------|------|
| 单专辑下载时间 | 8-12 分钟 |
| 网络 I/O | 高频小文件写入 |
| CPU 占用 | 30-40% |

#### 启用缓存机制

| 项目 | 数据 | 提升 |
|------|------|------|
| 单专辑下载时间 | 3-5 分钟 | **50-70%** |
| 网络 I/O | 批量大文件传输 | **90%+** |
| CPU 占用 | 20-30% | 25% |

---

## 🙏 致谢

### 原作者与核心贡献者

- **Sorrow** - 原始脚本作者和架构设计
- **chocomint** - 创建 ARM 支持
- **Sendy McSenderson** - 流解密代码

### 上游依赖

- **[mp4ff](https://github.com/Eyevinn/mp4ff)** by Eyevinn - MP4 文件处理
- **[mp4ff (fork)](https://github.com/itouakirai/mp4ff)** by itouakirai - 增强 MP4 支持
- **[progressbar/v3](https://github.com/schollz/progressbar)** by schollz - 进度显示
- **[requests](https://github.com/sky8282/requests)** by sky8282 - HTTP 客户端
- **[m3u8](https://github.com/grafov/m3u8)** by grafov - M3U8 解析
- **[pflag](https://github.com/spf13/pflag)** by spf13 - 命令行参数
- **[tablewriter](https://github.com/olekukonko/tablewriter)** by olekukonko - 表格格式化
- **[color](https://github.com/fatih/color)** by fatih - 彩色输出

### 外部工具

- **[FFmpeg](https://ffmpeg.org/)** - 音视频处理
- **[MP4Box](https://gpac.io/)** - GPAC 多媒体框架
- **[mp4decrypt](https://www.bento4.com/)** - Bento4 解密工具

### 特别感谢

- **[@sky8282](https://github.com/sky8282)** - 提供优秀的 requests 库和持续支持
- 所有贡献者和测试者
- Apple Music API 研究者和逆向工程社区
- 开源社区提供的各种库和工具

---

## ⚠️ 免责声明

本工具仅用于教育和个人使用。请遵守版权法和 Apple Music 服务条款。请勿分发下载的内容。

下载的音乐文件版权归原作者和 Apple Inc. 所有。使用本工具下载的内容仅限个人欣赏和学习，严禁用于商业用途或公开传播。

用户需自行承担使用本工具的法律责任。开发者不对任何因使用本工具而产生的法律问题负责。

---

## 📝 许可证

本项目采用个人使用许可证。详见 [LICENSE](./LICENSE) 文件。

所有下载内容的权利归其各自所有者所有。

---

## 🔗 相关资源

- [Apple Music for Artists](https://artists.apple.com/)
- [Emby 命名规范](https://emby.media/support/articles/Movie-Naming.html)
- [FFmpeg 文档](https://ffmpeg.org/documentation.html)
- [Go 官方文档](https://golang.org/doc/)

---

## 📈 更新日志

查看 [CHANGELOG.md](./CHANGELOG.md) 了解详细的版本历史和更新内容。

---

**版本**: v1.0.0  
**最后更新**: 2025-10-19  
**需要 Go 版本**: 1.23.1+

---

**⭐ 如果这个项目对您有帮助，请给个 Star！**
