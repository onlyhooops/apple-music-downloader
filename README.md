# Apple Music Downloader

[English](#) | [简体中文](./README-CN.md)

> **A powerful Apple Music high-quality audio downloader**
>
> Supports ALAC, Hi-Res Lossless, Dolby Atmos and other lossless formats, as well as music video downloads.

[![Version](https://img.shields.io/badge/version-v1.1.0-blue.svg)](https://github.com/onlyhooops/apple-music-downloader)
[![Go Version](https://img.shields.io/badge/Go-1.23.1+-00ADD8.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-Personal%20Use-green.svg)](./LICENSE)

---

> [!WARNING]
> **⚠️ Experimental Project Warning**
>
> This project is an experimental fork derived from [@sky8282/apple-music-downloader](https://github.com/sky8282/apple-music-downloader).
>
> **Important Notes:**
> - 🔧 **Personal Customization**: Extensive custom features added based on personal preferences
> - ⚠️ **Experimental Nature**: Some new features have not been widely tested and robustness cannot be guaranteed
> - 🎯 **Specific Environment**: Developed and tested on **Privileged LXC containers in Proxmox VE platform**
> - 📋 **Usage Recommendation**: Please carefully read the documentation and evaluate before use - may not suit all scenarios
>
> **For new features, refer to this README. For issues, check [Troubleshooting](#-troubleshooting) first.**
>
> For a stable version, please use the upstream project: https://github.com/sky8282/apple-music-downloader

---

## 🎉 What's New in v1.1.0

### 🔧 Important Fixes
- **✅ Fixed AAC Binaural/Downmix Download** - Corrected parameter matching bug, now works correctly
- **✅ Album MV Download Optimization** - Fixed path issues and improved AAC independent path configuration

### ⚡ Feature Enhancements
- **✨ New `--cx` Force Download Parameter** - Overwrite existing files, simplify re-download workflow
- **🧹 Removed History Feature** - Deleted 558 lines of code, simplified logic, improved performance
- **⚡ Enhanced File Validation Efficiency** - Optimized existence checks, faster batch downloads

### 🔬 Audio Quality Validation
- **✅ 100% Validation Pass Rate** - 8 audio quality parameter combinations professionally verified
- **✅ Parameter Consistency Check** - Command-line parameters match downloaded files perfectly
- **📊 40+ Technical Parameters Verified** - Using FFprobe 7.1 + MediaInfo 24.12

### 🧹 Configuration Optimization
- **Removed 4 Invalid Config Items** - skip-existing-validation, clean-choice, max-memory-limit, txtDownloadThreads
- **Net Optimization -464 Lines** - Cleaner, more efficient code

[View Complete Changelog](#-changelog)

---

## 📖 Table of Contents

- [What's New in v1.1.0](#-whats-new-in-v110)
- [Core Features](#-core-features)
- [System Requirements](#-system-requirements)
- [Quick Start](#-quick-start)
- [Usage Guide](#-usage-guide)
- [Advanced Features](#-advanced-features)
- [Configuration](#-configuration)
- [Troubleshooting](#-troubleshooting)
- [Changelog](#-changelog)
- [Acknowledgments](#-acknowledgments)

---

## ✨ Core Features

### 🎵 Multi-Format Audio Support

#### Lossless Audio Formats
- **ALAC (Apple Lossless)** - Native lossless audio with perfect quality
- **Hi-Res Lossless** - High-resolution lossless audio up to 24-bit/192kHz
- **Dolby Atmos** - Immersive 3D audio experience
- **Dolby EC-3** - Enhanced Dolby Digital audio

#### Lossy Audio Formats
- **AAC 256kbps** - High-quality AAC-LC encoding
- **AAC Binaural** - Binaural spatial audio
- **AAC Downmix** - Downmixed audio streams
- **AAC Stereo** - Standard stereo format

### 📹 Music Video Downloads

- **Multiple Resolutions** - 4K (2160p) / 1080p / 720p / 480p
- **Multi-Track Audio** - Atmos / AC3 / AAC tracks
- **Media Server Compatible** - Emby / Jellyfin / Plex naming conventions
- **Separate Save Path** - Customizable MV storage location

### 🎼 Complete Metadata Support

- **Album Artwork** - Up to 5000x5000 resolution
- **Lyrics Embedding** - Synchronized LRC lyrics (word-by-word/syllable-by-syllable)
- **Translation Lyrics** - Multi-language translation support (Beta)
- **Animated Artwork** - Supports animated artwork (requires FFmpeg)
- **Full Tags** - Artist, album, track number, release date, etc.

### ⚡ Performance Optimization

- **Cache Mechanism** - Optimized for NFS/SMB network storage, 50-70% speed improvement
- **Parallel Downloads** - Multi-threaded chunk downloads for maximum bandwidth utilization
- **Smart Resume** - Supports download resumption after interruption
- **Batch Processing** - Bulk download from TXT files

### 🛠️ Advanced Features

- **Force Download** - Use `--cx` parameter to overwrite existing files
- **Multi-Account Management** - Support multiple region accounts with auto-selection
- **Interactive Search** - Built-in search for songs/albums/artists
- **Custom Naming** - Flexible folder and file naming formats
- **Quality Tags** - Optional quality tags in filenames and metadata
- **FFmpeg Repair** - Automatic audio encoding issue detection and repair
- **Dynamic UI** - Real-time progress display with multi-task support
- **Pure Log Mode** - `--no-ui` mode suitable for CI/CD environments

---

## 📋 System Requirements

### Required Dependencies

1. **[Go 1.23.1+](https://golang.org/dl/)** - Compilation and runtime environment
2. **[MP4Box](https://gpac.io/downloads/gpac-nightly-builds/)** - Audio file processing (required)
3. **[mp4decrypt](https://www.bento4.com/downloads/)** - Music video decryption (required for MV downloads)
4. **[FFmpeg](https://ffmpeg.org/)** - Animated artwork and auto-repair features (optional)

### System Configuration

- **Operating System**: Linux / macOS / Windows
- **Memory**: 8GB+ recommended
- **Disk Space**: 
  - Without cache: Only space for downloaded files
  - With cache mechanism: Additional 50GB+ local temporary space
  - Large-scale batch downloads: 100GB+ recommended

### Apple Music Account

- **Basic Features**: Free account can download ALAC and Dolby Atmos
- **Advanced Features**: AAC 256, MV downloads, and lyrics require active subscription

---

## 🚀 Quick Start

### 1. Installation

```bash
# Clone repository
git clone https://github.com/onlyhooops/apple-music-downloader.git
cd apple-music-downloader

# Install Go dependencies
go mod tidy

# Build binary
go build -o apple-music-downloader main.go
```

### 2. Configuration

#### Create Configuration Files

```bash
# Copy configuration templates
cp config.yaml.example config.yaml
cp dev.env.example dev.env

# Edit configuration files
nano config.yaml
nano dev.env
```

#### Get Apple Music Token

1. Open [Apple Music Web](https://music.apple.com) and log in
2. Press **F12** to open Developer Tools
3. Go to **Application** tab
4. Navigate to **Cookies** → `https://music.apple.com` in the left sidebar
5. Find the cookie named `media-user-token`
6. Copy its value and paste into `dev.env`:

```bash
# dev.env
APPLE_MUSIC_MEDIA_USER_TOKEN_CN=your-token-here
```

#### Configure Save Paths

Edit `config.yaml`:

```yaml
# Save path configuration
alac-save-folder: "/media/Music/AppleMusic/Alac"
atmos-save-folder: "/media/Music/AppleMusic/Atmos"
mv-save-folder: "/media/Music/AppleMusic/MusicVideos"
```

### 3. Basic Usage

```bash
# Download album (default ALAC lossless)
./apple-music-downloader https://music.apple.com/us/album/album-name/123456789

# Download Dolby Atmos
./apple-music-downloader --atmos https://music.apple.com/us/album/album-name/123456789

# Download AAC 256 format
./apple-music-downloader --aac https://music.apple.com/us/album/album-name/123456789

# Download single track
./apple-music-downloader --song https://music.apple.com/us/album/album/123?i=456

# Download playlist
./apple-music-downloader https://music.apple.com/us/playlist/playlist-name/pl.xxxxx

# Download all works from an artist
./apple-music-downloader https://music.apple.com/us/artist/artist-name/123456
```

---

## 📖 Usage Guide

### Interactive Search

```bash
# Search for songs
./apple-music-downloader --search song "song name"

# Search for albums
./apple-music-downloader --search album "album name"

# Search for artists
./apple-music-downloader --search artist "artist name"
```

After searching, a results list will be displayed. Use arrow keys to select and Enter to confirm download.

### Batch Downloads

Create a `urls.txt` file with one link per line:

```text
https://music.apple.com/us/album/album1/123456789
https://music.apple.com/us/album/album2/987654321
https://music.apple.com/us/playlist/playlist/pl.xxxxx
# This is a comment line and will be ignored
```

Execute batch download:

```bash
./apple-music-downloader urls.txt
```

### Command Line Options

| Option | Description |
|--------|-------------|
| `--atmos` | Download Dolby Atmos format |
| `--aac` | Download AAC 256 format |
| `--aac-type <type>` | Specify AAC type: `aac-lc`, `binaural`, `downmix` |
| `--alac-max <rate>` | Specify ALAC max sample rate: `192000`, `96000`, `48000` |
| `--atmos-max <bitrate>` | Specify Atmos max bitrate: `2768`, `2448` |
| `--song` | Download single track mode |
| `--select` | Interactive track selection |
| `--all-album` | Download all albums from an artist |
| `--search <type> "keyword"` | Search: `song`, `album`, `artist` |
| `--mv-max <resolution>` | MV max resolution: `2160`, `1080`, `720` |
| `--mv-audio-type <type>` | MV audio track type: `atmos`, `ac3`, `aac` |
| `--debug` | Display available quality information (no download) |
| `--no-ui` | Disable dynamic UI, pure log output |
| `--config <path>` | Specify configuration file path |
| `--output <path>` | Specify output directory for this task |
| `--start <number>` | Start from specific link in TXT file (for resume) |

---

## 🎯 Advanced Features

### Cache Mechanism (NFS/SMB Optimization)

When downloading to network storage (NFS/SMB), enabling cache mechanism significantly improves performance:

#### Configuration

Edit `config.yaml`:

```yaml
# Enable cache
enable-cache: true
cache-folder: "./Cache"  # Recommend using local SSD path
```

#### How It Works

1. Files are first downloaded to local cache folder (high-speed disk)
2. All processing (decryption, merging, metadata) completed locally
3. Batch transfer to target network path when finished
4. Automatic cache cleanup to free space

#### Performance Improvement

- Download time: **50-70%** faster
- Network I/O: **90%+** reduction
- Stability: Atomic operations with automatic rollback on failure

#### Important Notes

- ⚠️ Cache folder must be on local fast disk (SSD)
- ⚠️ Do not set cache on NFS or network paths
- ⚠️ Ensure at least 50GB free space
- ⚠️ Program auto-cleans; can also manually delete `Cache` folder

### Force Download Mode

Use the `--cx` parameter to force overwrite existing files, suitable for re-downloading or updating files.

#### Usage Examples

```bash
# Force download album (overwrite existing files)
./apple-music-downloader --cx https://music.apple.com/us/album/xxx/123

# Force download AAC Binaural
./apple-music-downloader --cx --aac --aac-type aac-binaural <url>

# Force download Dolby Atmos
./apple-music-downloader --cx --atmos <url>

# Batch force download
./apple-music-downloader --cx urls.txt
```

#### Use Cases

- 🔄 Corrupted files need re-download
- 🎵 Want to replace with different quality version
- 🆕 Apple Music updated audio quality
- 🔧 Changed naming format, need to regenerate files

### Quality Tag Configuration

From v1.1.0, flexible control over quality tag display:

```yaml
# config.yaml
add-quality-tag-to-folder: true      # Include quality tag in folder name
add-quality-tag-to-metadata: true    # Include quality tag in metadata
```

#### Configuration Effects

| Folder Tag | Metadata Tag | Folder Name | Metadata ALBUM | Use Case |
|:---:|:---:|---|---|---|
| ✅ | ✅ | `Head Hunters Alac/` | `Head Hunters Alac` | **Recommended**: Perfect sync, music software correctly identifies |
| ✅ | ❌ | `Head Hunters Alac/` | `Head Hunters` | File classification clear, metadata concise |
| ❌ | ✅ | `Head Hunters/` | `Head Hunters Alac` | Folder concise, quality in metadata |
| ❌ | ❌ | `Head Hunters/` | `Head Hunters` | Not recommended: Cannot distinguish quality versions |

#### Recommendations

- 🎵 **Plex/Emby/Jellyfin Users**: Enable both
- 💿 **Collecting Multiple Quality Versions**: Enable both
- 🗂️ **File Classification Only**: Enable folder tag only
- ✨ **Minimalist**: Enable metadata tag only

### Custom Naming Formats

#### Available Variables

**Album-related**:
- `{AlbumId}` - Album ID
- `{AlbumName}` - Album name
- `{ArtistName}` - Artist name
- `{ReleaseDate}` - Release date
- `{ReleaseYear}` - Release year
- `{Tag}` - Quality tag (e.g., "Alac", "Hi-Res Lossless")
- `{Quality}` - Quality description
- `{Codec}` - Codec format
- `{UPC}` - UPC code
- `{Copyright}` - Copyright information
- `{RecordLabel}` - Record label

**Track-related**:
- `{SongId}` - Track ID
- `{SongName}` - Track name
- `{SongNumer}` - Track number (two digits)
- `{TrackNumber}` - Track number (original)
- `{DiscNumber}` - Disc number

**Artist-related**:
- `{ArtistId}` - Artist ID
- `{ArtistName}` - Artist name
- `{UrlArtistName}` - Artist name from URL

#### Naming Examples

```yaml
# config.yaml

# Album folder: "{Album Name} {Quality Tag}"
album-folder-format: "{AlbumName} {Tag}"
# Result: Head Hunters Hi-Res Lossless/

# Track file: "{Number}. {Track Name}"
song-file-format: "{SongNumer}. {SongName}"
# Result: 01. Chameleon.m4a

# Artist folder: "{Artist Name}"
artist-folder-format: "{ArtistName}"
# Result: Herbie Hancock/

# Playlist folder: "{Playlist Name} {Quality Tag}"
playlist-folder-format: "{PlaylistName} {Tag}"
# Result: Jazz Classics Alac/
```

### Multi-Account Configuration

Support multiple Apple Music region accounts, program auto-selects based on link region:

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

### Logging Configuration

```yaml
# config.yaml
logging:
  level: info                  # Log level: debug/info/warn/error
  output: stdout               # Output target: stdout/stderr/file path
  show_timestamp: false        # Recommend off for UI mode
```

**Log Levels**:
- `debug` - Show all debug information (for development and troubleshooting)
- `info` - Show general information (default, recommended)
- `warn` - Show warnings and errors only
- `error` - Show errors only

**Recommendations**:
- Dynamic UI mode: `show_timestamp: false`
- Pure log mode (`--no-ui`): `show_timestamp: true`
- CI/CD environment: Use `--no-ui` + log file output

---

## ⚙️ Configuration

### Performance Tuning

#### Network Storage (NFS/SMB)

```yaml
enable-cache: true
cache-folder: "/ssd/cache/apple-music"  # Use local SSD
chunk_downloadthreads: 30
```

#### Batch Download Optimization

```yaml
batch-size: 20                           # Items per batch
skip-existing-validation: true           # Auto skip existing files
work-rest-enabled: true                  # Work-rest cycle
work-duration-minutes: 30                # Work for 30 minutes
rest-duration-minutes: 2                 # Rest for 2 minutes
```

#### Download Thread Configuration

```yaml
# M3U8 chunk downloads
chunk_downloadthreads: 30                # Audio chunk threads
mv_chunk_downloadthreads: 30             # Video chunk threads

# Audio format threads
aac_downloadthreads: 5                   # AAC format
lossless_downloadthreads: 5              # Lossless format
hires_downloadthreads: 5                 # Hi-Res format

# MV downloads
mv_downloadthreads: 3                    # Parallel MV downloads
```

### FFmpeg Configuration

```yaml
ffmpeg-fix: true                         # Auto detect and repair
ffmpeg-check-args: "-map 0:a:0 -f wav -hide_banner -loglevel error -"
ffmpeg-encode-args: "-c:v copy -c:a alac -avoid_negative_ts make_zero -f mp4 -y"
```

---

## 📂 Output Structure

### Album Structure (Emby Compatible)

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

### Music Video Structure (Emby/Jellyfin Compatible)

```
/media/Music/AppleMusic/MusicVideos/
└── Morgan James/
    └── Thunderstruck (2024)/
        └── Thunderstruck (2024).mp4
```

---

## 🐛 Troubleshooting

### 1. "MP4Box not found"

**Cause**: MP4Box not installed or not in system PATH

**Solution**:
```bash
# Install MP4Box
# Linux (Ubuntu/Debian):
sudo apt-get install gpac

# macOS:
brew install gpac

# Verify installation
MP4Box -version
```

### 2. "No media-user-token"

**Cause**:
- AAC-LC, MV, and lyrics features require valid subscription token
- ALAC and Dolby Atmos work with basic token

**Solution**: Ensure token is correctly configured (see "Get Apple Music Token" section)

### 3. UI Output Chaos

**Cause**: Terminal doesn't support dynamic updates or output redirection

**Solution**:
```bash
# Use pure log mode
./apple-music-downloader --no-ui <url>

# Save log to file
./apple-music-downloader --no-ui <url> > download.log 2>&1
```

### 4. Slow NFS Downloads

**Cause**: Network filesystem latency and frequent small file writes

**Solution**: Enable cache mechanism (see "Cache Mechanism" section)

### 5. How to Resume After Interruption

**Method 1**: Program automatically skips existing files
```bash
# Re-run the same command
./apple-music-downloader urls.txt
```

**Method 2**: Use `--start` parameter
```bash
# Start from link 44
./apple-music-downloader --start 44 urls.txt
```

### 6. How to Clean Cache

```bash
# Manually delete cache folder
rm -rf ./Cache

# Program will auto-rebuild on next run
```

---

## 📊 Performance Reference

### Test Environment

- **Server**: Proxmox VE 6.8.12
- **CPU**: 8 Core @ 2.4GHz
- **Memory**: 16GB
- **Storage**: NFS network storage
- **Network**: 1Gbps

### Performance Data

#### Without Cache (Direct NFS Write)

| Item | Data |
|------|------|
| Single Album Download | 8-12 minutes |
| Network I/O | High-frequency small file writes |
| CPU Usage | 30-40% |

#### With Cache Mechanism

| Item | Data | Improvement |
|------|------|-------------|
| Single Album Download | 3-5 minutes | **50-70%** |
| Network I/O | Batch large file transfers | **90%+** |
| CPU Usage | 20-30% | 25% |

---

## 🙏 Acknowledgments

### Original Authors & Core Contributors

- **Sorrow** - Original script author and architecture design
- **chocomint** - Created ARM support
- **Sendy McSenderson** - Stream decryption code

### Upstream Dependencies

- **[mp4ff](https://github.com/Eyevinn/mp4ff)** by Eyevinn - MP4 file processing
- **[mp4ff (fork)](https://github.com/itouakirai/mp4ff)** by itouakirai - Enhanced MP4 support
- **[progressbar/v3](https://github.com/schollz/progressbar)** by schollz - Progress display
- **[requests](https://github.com/sky8282/requests)** by sky8282 - HTTP client
- **[m3u8](https://github.com/grafov/m3u8)** by grafov - M3U8 parsing
- **[pflag](https://github.com/spf13/pflag)** by spf13 - Command line flags
- **[tablewriter](https://github.com/olekukonko/tablewriter)** by olekukonko - Table formatting
- **[color](https://github.com/fatih/color)** by fatih - Colored output

### External Tools

- **[FFmpeg](https://ffmpeg.org/)** - Audio/video processing
- **[MP4Box](https://gpac.io/)** - GPAC multimedia framework
- **[mp4decrypt](https://www.bento4.com/)** - Bento4 decryption tool

### Special Thanks

- **[@sky8282](https://github.com/sky8282)** - Excellent requests library and continuous support
- All contributors and testers
- Apple Music API researchers and reverse engineering community
- Open source community for various libraries and tools

---

## ⚠️ Disclaimer

This tool is for educational and personal use only. Please comply with copyright laws and Apple Music Terms of Service. Do not distribute downloaded content.

Downloaded music files are copyrighted by the original authors and Apple Inc. Content downloaded using this tool is for personal enjoyment and learning only. Commercial use or public distribution is strictly prohibited.

Users are responsible for any legal issues arising from the use of this tool. Developers are not responsible for any legal problems arising from the use of this tool.

---

## 📝 License

This project uses a Personal Use License. See [LICENSE](./LICENSE) file for details.

All rights to downloaded content belong to their respective owners.

---

## 🔗 Related Resources

- [Apple Music for Artists](https://artists.apple.com/)
- [Emby Naming Conventions](https://emby.media/support/articles/Movie-Naming.html)
- [FFmpeg Documentation](https://ffmpeg.org/documentation.html)
- [Go Official Documentation](https://golang.org/doc/)

---

---

## 📈 Changelog

### v1.1.0 (2025-10-20)

#### 🔧 Important Fixes
- **Fixed AAC Binaural/Downmix Parameter Matching Bug**
  - Issue: Parameters defined as `aac-binaural/aac-downmix`, but code checked `binaural/downmix`
  - Impact: Unable to download AAC Binaural and Downmix quality correctly
  - Fix: Corrected 4 check logics in downloader, parser, and metadata
  - Result: ✅ Now correctly downloads AAC Binaural/Downmix

- **Fixed Album MV Download Path Issues**
  - Optimized AAC independent path configuration logic
  - Improved log output format

#### ⚡ Feature Optimizations
- **New `--cx` Force Download Parameter**
  - Function: Force download mode, overwrite existing files
  - Scenarios: Re-download, update files, fix corrupted files
  - Usage: `./apple-music-downloader --cx <url>`

- **Removed History Feature**
  - Deleted `internal/history/global.go` (558 lines of code)
  - Simplified download logic, improved performance
  - Reduced maintenance costs

- **Enhanced File Validation Efficiency**
  - Optimized file existence check logic
  - Reduced unnecessary file system operations
  - Improved batch download performance

#### 🧹 Configuration Cleanup
- **Removed 4 Invalid Config Items** (fully validated)
  - `skip-existing-validation` - Obsolete configuration
  - `clean-choice` - Not actually used
  - `max-memory-limit` - Only defined, not used
  - `txtDownloadThreads` - Only validated, not actually used
- Added `CONFIG_VERIFICATION.md` configuration verification report
  - Verified all 44 configuration items
  - Valid configurations: 40 items (90.9%)

#### 🔬 Audio Quality Validation
- **Completed 8 Audio Quality Parameter Professional Verification**
  - Validation pass rate: **100%** ⭐⭐⭐⭐⭐
  - Validation tools: FFprobe 7.1 + MediaInfo 24.12
  - Test track: Tyla - IS IT (Apple Music ID: 1825447073)

- **Command-line Parameters vs Downloaded Files Consistency Check**
  - Verified encoding format, bitrate, sample rate, channels, and 40+ technical parameters
  - Confirmed intelligent parameter downgrade mechanism works correctly
  - Verified metadata tags accurately identify quality versions
  - Validation result: ✅ **100% Consistent**

- **Technical Validation Results**
  - ALAC Lossless: 48kHz/24bit, 1743 kbps
  - Dolby Atmos: E-AC-3 JOC, 5.1 channels, 15 dynamic objects
  - AAC Binaural/Downmix: 48kHz sample rate, complete metadata
  - Complies with MPEG-4, ATSC A/52B industry standards

#### 📚 Documentation Updates
- **README-CN.md** added "Audio Quality Validation" section
  - Validation results overview table
  - Detailed description of verified quality parameters
  - File size comparison analysis
- Added audio quality test report (5.3 KB)
- Added professional verification report (17 KB, parameter consistency analysis)
- Added project structure documentation
- Created `docs/验证报告/` and `docs/开发文档/` directories

#### 📊 Code Statistics
- Code deleted: 558 lines (history feature)
- Code added: 94 lines (optimized validation logic)
- Bug fixes: 16 modifications (AAC Binaural/Downmix)
- **Net optimization: -464 lines** (cleaner code)

---

View [CHANGELOG.md](./CHANGELOG.md) for complete version history.

---

**Version**: v1.1.0  
**Last Updated**: 2025-10-20  
**Required Go Version**: 1.23.1+

---

**⭐ If this project helps you, please give it a Star!**
