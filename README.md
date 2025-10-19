# Apple Music ALAC / Dolby Atmos Downloader

[English](#) / [简体中文](./README-CN.md)

> [!WARNING]
> **⚠️ Experimental Branch Warning**
>
> This branch is an experimental version for personal use only, containing extensive modifications with numerous unknown bugs and risks. Please use with caution!
>
> This branch contains experimental features that have not been thoroughly tested, which may lead to data loss, download failures, or other unforeseen issues. For production environments, please use the official stable version.

**A powerful Apple Music download tool supporting ALAC, Dolby Atmos, Hi-Res Lossless and other high-quality audio formats, as well as music video downloads.**

Original script by Sorrow, extensively improved and optimized.

---

## 🎉 Latest Updates (v2.8.2)

### 🚀 v2.8.2 - Audio Stream Selection Enhancement (2025-10-19)

#### 🔧 Core Improvements
- **Smart Quality Selection** - Intelligently selects the best available quality while respecting user preferences
- **Binaural Audio Support** - Full support for binaural audio streams (`--aac --aac-type binaural`)
- **Enhanced Audio Traits Detection** - Improved detection of audio characteristics for better quality selection

#### ⚡ Performance Optimizations
- **Stream Selection Algorithm** - Optimized M3U8 parsing for accurate audio stream selection
- **Quality Tag Generation** - Fixed quality tag generation to reflect actual downloaded content
- **History System Enhancement** - Improved quality hash calculation for better duplicate detection

#### 🐛 Critical Fixes
- **Audio Stream Selection** - Fixed issue where binaural streams were not properly selected
- **Quality Tag Consistency** - Resolved inconsistency between downloaded quality and displayed tags
- **Empty Folder Creation** - Eliminated creation of unnecessary empty folders

### 🎯 v2.7.0 - Global History & Architecture Refactor (2025-10-12)

#### 🏗️ Architecture Refactoring
- **Unified Logging System** - New `internal/logger` module with 4-level log control (DEBUG/INFO/WARN/ERROR)
- **Progress Event System** - Observer pattern event architecture, decoupling UI updates from business logic
- **Smart UI Simplification** - Adaptive terminal width display (FullMode/CompactMode/MinimalMode)
- **Album Metadata Enhancement** - Album and AlbumSort fields now include quality tags (e.g., "Head Hunters Hi-Res Lossless")

#### ⚡ Performance Optimizations
- **Reduced UI Refresh Rate** - Changed from 200ms to intelligent refresh, reducing CPU usage
- **Logger Performance Optimization** - Supports high-concurrency logging output
- **Improved Progress Event Distribution** - Enhanced multi-task concurrency efficiency

#### 🐛 Critical Fixes
- **Logger Duplication Issue** - Fixed logger output interfering with UI cursor positioning (redirected logger to stderr)
- **UI Rendering Misalignment** - Fixed line overlapping and scrolling issues caused by long track names
- **Metadata Quality Tags** - Fixed issue where music management software couldn't distinguish different quality versions of the same album

---

## ✨ Core Features

### 🎵 Audio Format Support
- **ALAC (Lossless Audio)** - `audio-alac-stereo`
- **Dolby Atmos / EC3** - `audio-atmos` / `audio-ec3`
- **Hi-Res Lossless (High-Resolution Lossless)** - Up to 24-bit/192kHz
- **AAC Formats** - `audio-stereo`, `audio-stereo-binaural`, `audio-stereo-downmix`
- **AAC-LC** - `audio-stereo` (requires subscription)

### 📹 Music Video Support
- **4K/1080p/720p** resolution options
- **Emby/Jellyfin compatible** naming structure
- **Multi-audio track support** (Atmos/AC3/AAC)
- Independent save path configuration

### 🎼 Rich Metadata and Lyrics
- **Embedded artwork** and artwork (up to 5000x5000)
- **Synchronized LRC lyrics** (word-by-word / syllable-by-syllable)
- **Translation and pronunciation lyrics** support (Beta)
- **Dynamic artwork** (Emby/Jellyfin support)
- Complete metadata tags

### ⚡ Performance Optimizations
- **Cache transfer mechanism** - 50-70% speed improvement for NFS/network storage
- **Parallel downloads** - Multi-threaded segment downloads
- **Smart file checking** - Skip already downloaded files
- **Batch downloads** - Read links from TXT files with configurable thread count

### 🛠️ Advanced Features
- **Multi-account rotation** - Automatic account selection based on region
- **FFmpeg auto-repair** - Detect and repair encoding issues
- **Interactive mode** - Arrow key navigation for search results
- **Artist downloads** - Download all albums/MVs from an artist page
- **Custom naming** - Flexible folder and file naming formats
- **Output modes** - Dynamic UI or pure log mode (`--no-ui`)

---

## 📋 System Requirements

### Required Dependencies

1. **[MP4Box](https://gpac.io/downloads/gpac-nightly-builds/)** - Must be installed and added to system PATH
2. **[mp4decrypt](https://www.bento4.com/downloads/)** - Required for music video downloads
3. **FFmpeg** (optional) - For dynamic artwork and auto-repair features

### System Requirements
- Go 1.23.1 or higher
- Recommended 8GB+ memory
- 50GB+ available disk space (if using cache mechanism)

> [!NOTE]
> **💡 Disk Space Recommendations**
>
> - **Without cache**: Only need space for downloaded files
> - **With cache mechanism**: Additional 50GB+ local temporary space required
> - **Large-scale batch downloads**: Recommend reserving 100GB+ space for optimal performance

---

## 🚀 Quick Start

### 1. Installation

```bash
# Clone repository
git clone https://github.com/onlyhooops/apple-music-downloader.git
cd apple-music-downloader

# Install dependencies
go mod tidy

# Build binary
go build -o apple-music-downloader main.go
```

### 2. Configuration

```bash
# Copy example configuration
cp config.yaml.example config.yaml

# Edit configuration file
nano config.yaml
```

**Get `media-user-token`:**
1. Open [Apple Music](https://music.apple.com) and log in
2. Press `F12` to open Developer Tools
3. Go to `Application` (应用程序) → `Cookies` → `https://music.apple.com`
4. Find the cookie named `media-user-token` and copy its value
5. Paste it into `config.yaml`

### 3. Basic Usage

```bash
# Download album
./apple-music-downloader https://music.apple.com/cn/album/album-name/123456789

# Download Dolby Atmos
./apple-music-downloader --atmos https://music.apple.com/cn/album/album-name/123456789

# Download single track
./apple-music-downloader --song https://music.apple.com/cn/album/album/123?i=456

# Download playlist
./apple-music-downloader https://music.apple.com/cn/playlist/playlist-name/pl.xxxxx

# Download all content from an artist
./apple-music-downloader https://music.apple.com/cn/artist/artist-name/123456

# Interactive search
./apple-music-downloader --search song "search terms"
./apple-music-downloader --search album "album name"
./apple-music-downloader --search artist "artist name"

# Batch download from TXT file
./apple-music-downloader urls.txt

# Pure log mode (for CI/debugging)
./apple-music-downloader --no-ui https://music.apple.com/...
```

---

## 📖 Advanced Usage

### Environment Variables Setup

Create a `dev.env` file based on `dev.env.example`:

```bash
# Copy template
cp dev.env.example dev.env

# Edit and add your credentials
nano dev.env
```

**Required environment variables:**
```bash
APPLE_MUSIC_MEDIA_USER_TOKEN_CN=your-media-user-token-here
APPLE_MUSIC_AUTH_TOKEN_CN=your-auth-token
```

### Cache Mechanism (NFS Optimization)

> [!IMPORTANT]
> **⚠️ Cache Mechanism Important Notes**
>
> **Recommended scenarios:**
> - ✅ Target path is NFS/SMB or other network file systems
> - ✅ Local disk has 50GB+ available space (SSD recommended)
> - ✅ Need frequent batch download tasks
>
> **Key considerations:**
> - ⚠️ **Disk space**: Cache folder requires sufficient temporary storage space, recommend at least 50GB
> - ⚠️ **Cache path**: Must use local fast disk (SSD), don't set on NFS or other network paths
> - ⚠️ **File system**: Cross-filesystem transfer will use copy method, speed will be reduced
> - ⚠️ **Cleanup mechanism**: Program automatically cleans successfully transferred cache, also auto-rollback on failure
> - ⚠️ **Manual cleanup**: Can manually delete `Cache` folder anytime, program will auto-rebuild
>
> **Performance improvement data (measured):**
> - Download time improvement: **50-70%**
> - Network I/O reduction: **90%+**
> - Better stability: Atomic operations, automatic rollback on failure

When downloading to network storage (NFS/SMB), performance can be significantly improved:

```yaml
# config.yaml
enable-cache: true
cache-folder: "./Cache"  # Recommend using local SSD path
```

**Configuration recommendations:**
- ⚡ **Local SSD cache** - Set `cache-folder` to local SSD path (like `/ssd/cache/apple-music`)
- ⚡ **Network storage target** - Set `alac-save-folder` and `atmos-save-folder` to NFS/SMB paths
- ⚡ **Sufficient space** - Ensure cache path has at least 50GB available space

**Working principle:**
1. Files are first downloaded to local cache folder
2. All processing (decryption, merging, metadata) is done locally
3. Completed files are transferred to target network path in batches
4. Cache is automatically cleaned, freeing up space

[📚 Read complete cache mechanism documentation](./CACHE_MECHANISM.md)

### Global History & Resume (v2.7.0+ Enhanced)

> [!TIP]
> **🔄 Global Link-Level History**
>
> All download tasks (single link, multi-link, text batch) now share a global history system. Any successfully downloaded link will be automatically skipped, regardless of which task or file it came from.

**Core features:**
- 🌐 **Global deduplication**: History records are link-based, not task-batch-based
- 📁 **Auto-recording**: All modes automatically record to `history` folder
- 🔍 **Cross-task detection**: Automatically skip links downloaded in previous tasks
- 🎵 **Quality-aware**: Detect quality parameter changes, re-download when necessary
- ⏸️ **Smart resume**: Support resuming from breakpoint after task interruption

**Usage examples:**
```bash
# Scenario 1: Download classic albums list
./apple-music-downloader classic_albums.txt
# Complete download of 50 albums

# Scenario 2: Download jazz music list (10 overlap with classic albums)
./apple-music-downloader jazz_music.txt
# Output: 📜 Global history detection: Found 10 completed tasks
#         ⏭️  Auto-skipped 10, remaining 40 tasks

# Scenario 3: Download individual album (already in previous lists)
./apple-music-downloader https://music.apple.com/cn/album/xxx/123
# Output: ✅ All tasks completed, no duplicate downloads!
```

**Supported modes:**
- ✅ Single link mode
- ✅ Multi-link mode
- ✅ TXT file batch mode
- ✅ Mixed mode (URL + TXT)
- ✅ Interactive mode

**Advanced usage:**
```bash
# View all history records
ls -lh history/

# Clear all history (start fresh)
rm -rf history/

# Clear specific task history
rm history/classic_albums.txt_*.json
```

[📚 Read complete history feature documentation](./HISTORY_FEATURE.md)

### Logging Configuration (v2.6.0+)

**Unified logging system**, supporting 4-level log control and flexible configuration:

```yaml
# config.yaml
logging:
  level: info                  # debug/info/warn/error
  output: stdout               # stdout/stderr/file path
  show_timestamp: false        # UI mode suggests disabling
```

**Log levels:**
- `debug` - Show all debug information (for development and troubleshooting)
- `info` - Show general information (default, recommended)
- `warn` - Show only warnings and errors
- `error` - Show only error information

**Output targets:**
- `stdout` - Standard output (default)
- `stderr` - Standard error output (automatically used in UI mode)
- File path - Such as `./logs/download.log`

**Usage recommendations:**
- Dynamic UI mode: `show_timestamp: false`, avoid timestamp interfering with UI
- Pure log mode (`--no-ui`): `show_timestamp: true`, convenient for tracing
- CI/CD environment: Use `--no-ui` + log file output

### Custom Naming Format

> [!TIP]
> **🏷️ Quality Tag Configuration (v2.5.0+)**
>
> From v2.5.0, you can flexibly control where quality tags are displayed:

```yaml
# config.yaml - Quality tag configuration
add-quality-tag-to-folder: true      # Folder name includes quality tag
add-quality-tag-to-metadata: true    # Metadata includes quality tag
```

**Configuration combination effects:**

| Folder Tag | Metadata Tag | Folder Name | Metadata ALBUM | Use Cases |
|:---:|:---:|---|---|---|
| ✅ | ✅ | `Head Hunters Alac/` | `Head Hunters Alac` | **Recommended**: Perfect sync, music software can correctly identify |
| ✅ | ❌ | `Head Hunters Alac/` | `Head Hunters` | File classification clear, metadata concise |
| ❌ | ✅ | `Head Hunters/` | `Head Hunters Alac` | Folder concise, quality info in metadata |
| ❌ | ❌ | `Head Hunters/` | `Head Hunters` | Not recommended: Cannot distinguish different quality versions |

**Usage recommendations:**
- 🎵 **Plex/Emby/Jellyfin users**: Enable both (true)
- 💿 **Collecting multiple quality versions**: Enable both (true)
- 🗂️ **Only need file classification**: Enable folder tag only
- ✨ **Pursue conciseness**: Enable metadata tag only

```yaml
# Album folder: "Album Name Quality Tag"
album-folder-format: "{AlbumName} {Tag}"

# Song file: "01. Song Name"
song-file-format: "{SongNumer}. {SongName}"

# Artist folder: "Artist Name"
artist-folder-format: "{ArtistName}"

# Playlist folder: "Playlist Name"
playlist-folder-format: "{PlaylistName}"
```

**Available variables:**
- Album: `{AlbumId}`, `{AlbumName}`, `{ArtistName}`, `{ReleaseDate}`, `{ReleaseYear}`, `{Tag}`, `{Quality}`, `{Codec}`, `{UPC}`, `{Copyright}`, `{RecordLabel}`
- Song: `{SongId}`, `{SongNumer}`, `{SongName}`, `{DiscNumber}`, `{TrackNumber}`, `{Tag}`, `{Quality}`, `{Codec}`
- Playlist: `{PlaylistId}`, `{PlaylistName}`, `{ArtistName}`, `{Tag}`, `{Quality}`, `{Codec}`
- Artist: `{ArtistId}`, `{ArtistName}`, `{UrlArtistName}`

### Multi-Account Configuration

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

The program automatically selects the corresponding account based on the URL region (such as `/cn/`, `/us/`).

---

## 🔧 Command Line Options

| Option | Description |
|--------|-------------|
| `--atmos` | Download Dolby Atmos format |
| `--aac` | Download AAC 256 format |
| `--song` | Download single track |
| `--select` | Interactive track selection |
| `--search [type] "keywords"` | Search (song/album/artist) |
| `--debug` | Show available quality information |
| `--no-ui` | Disable dynamic UI, pure log output |
| `--config path` | Specify custom config file |
| `--output path` | Override save folder |

---

## 📂 Output Structure

### Albums (Emby Compatible Naming)

```
/media/Music/AppleMusic/Alac/
└── Taylor Swift/
    └── 1989 (Taylor's Version) Hi-Res Lossless/
        ├── cover.jpg
        ├── 01. Welcome To New York.m4a
        ├── 02. Blank Space.m4a
        └── ...
```

### Music Videos (Emby/Jellyfin Compatible)

```
/media/Music/AppleMusic/MusicVideos/
└── Morgan James/
    └── Thunderstruck (2024)/
        └── Thunderstruck (2024).mp4
```

---

## 🐛 Troubleshooting

### Common Issues

**1. "MP4Box not found"**
- Install [MP4Box](https://gpac.io/downloads/gpac-nightly-builds/)
- Ensure it's added to system PATH
- Test: `MP4Box -version`

**2. "No media-user-token"**
- AAC-LC, MV and lyrics require valid subscription token
- ALAC/Dolby Atmos work with basic token

**3. UI output chaos**
- Use `--no-ui` flag for pure log output
- More suitable for CI/CD processes or output redirection

**4. NFS downloads slow**
- Enable cache mechanism in config.yaml
- Refer to [Cache Quick Start Guide](./QUICKSTART_CACHE.md)

### FFmpeg Auto-Repair

If downloaded files have encoding issues:

```yaml
ffmpeg-fix: true  # Auto-detect after download
```

The program will:
1. Detect damaged/incomplete files
2. Prompt for confirmation
3. Re-encode using FFmpeg and ALAC codec

---

## 📊 Performance Recommendations

### For Network Storage (NFS/SMB)
- ✅ Enable cache mechanism
- ✅ Use local SSD as cache folder
- ✅ Increase segment download thread count

### For Batch Downloads
```yaml
txtDownloadThreads: 5  # Parallel album downloads
chunk_downloadthreads: 30  # Parallel segment downloads
```

### For Large Music Libraries
- ✅ Enable `ffmpeg-fix` for quality assurance
- ✅ Use `--no-ui` for clearer logs
- ✅ Save output to file: `./app --no-ui url > download.log 2>&1`

---

## 📚 Documentation

### User Guides
- [README.md](./README.md) - English documentation
- [QUICKSTART_CACHE.md](./QUICKSTART_CACHE.md) - Cache mechanism quick start
- [CACHE_UPDATE.md](./CACHE_UPDATE.md) - Cache update guide
- [GOO_ALIAS.md](./GOO_ALIAS.md) - Command alias configuration guide
- [EMOJI_DEMO.md](./EMOJI_DEMO.md) - Emoji output demo

### Technical Documentation
- [CHANGELOG.md](./CHANGELOG.md) - Complete version history and changes
- [CACHE_MECHANISM.md](./CACHE_MECHANISM.md) - Complete cache technical documentation
- [MV_QUALITY_DISPLAY.md](./MV_QUALITY_DISPLAY.md) - MV quality detection feature
- [MV_PROGRESS_FIX.md](./MV_PROGRESS_FIX.md) - MV progress tracking improvements
- [MV_LOG_FIX.md](./MV_LOG_FIX.md) - MV download log enhancement

---

## 🙏 Acknowledgments

### 🎖️ Original Authors & Core Contributors
- **Sorrow** - Original script author and architecture design
- **chocomint** - Created `agent-arm64.js` ARM support
- **Sendy McSenderson** - Stream decryption code

### 🔧 Upstream Dependencies & Tools
- **[mp4ff](https://github.com/Eyevinn/mp4ff)** by Eyevinn - MP4 file processing
- **[mp4ff (fork)](https://github.com/itouakirai/mp4ff)** by itouakirai - Enhanced MP4 support
- **[progressbar/v3](https://github.com/schollz/progressbar)** by schollz - Progress display
- **[requests](https://github.com/sky8282/requests)** by sky8282 - HTTP client wrapper
- **[m3u8](https://github.com/grafov/m3u8)** by grafov - M3U8 playlist parsing
- **[pflag](https://github.com/spf13/pflag)** by spf13 - Command line parameters
- **[tablewriter](https://github.com/olekukonko/tablewriter)** by olekukonko - Table formatting
- **[color](https://github.com/fatih/color)** by fatih - Colored terminal output

### 🛠️ External Tools
- **[FFmpeg](https://ffmpeg.org/)** - Audio/video processing
- **[MP4Box](https://gpac.io/)** - GPAC multimedia framework
- **[mp4decrypt](https://www.bento4.com/)** - Bento4 decryption tool

### 💝 Special Thanks
- **[@sky8282](https://github.com/sky8282)** - Provided excellent requests library and continuous support
- All contributors and testers who helped improve this project
- Apple Music API researchers and reverse engineering community
- Open source community providing various libraries and tools

---

## ⚠️ Disclaimer

This tool is for educational and personal use only. Please comply with copyright laws and Apple Music Terms of Service. Please do not distribute downloaded content.

---

## 📝 License

This project is for personal use only. All rights to downloaded content belong to their respective owners.

---

## 🔗 Resources

- [Apple Music for Artists](https://artists.apple.com/)
- [Emby Naming Conventions](https://emby.media/support/articles/Movie-Naming.html)
- [FFmpeg Documentation](https://ffmpeg.org/documentation.html)
- [Chinese Tutorial](https://telegra.ph/Apple-Music-Alac高解析度无损音乐下载教程-04-02-2)

---

**Version:** v2.8.2
**Last Updated:** 2025-10-19
**Required Go Version:** 1.23.1+
