# 🚀 goo 命令别名配置说明

## 📋 基本信息

| 项目 | 内容 |
|------|------|
| **命令别名** | `goo` |
| **指向版本** | apple-music-downloader-v2.2.0 |
| **配置文件** | `~/.zshrc` (第148行) |
| **二进制路径** | `/root/apple-music-downloader/apple-music-downloader-v2.2.0` |
| **文件大小** | 37MB |
| **编译时间** | 2025-10-09 06:57 |
| **版本特性** | 包含中文帮助、emoji美化、所有最新功能 |

---

## ✨ 特性说明

### 🔧 完整功能集成

✅ **日志UI治理** - OutputMutex + SafePrintf  
✅ **中文帮助菜单** - 所有参数说明已汉化  
✅ **Emoji美化** - 终端输出更直观美观  
✅ **--no-ui 模式** - 纯日志输出支持  
✅ **缓存机制** - NFS性能优化  
✅ **所有v2.2.0功能** - 完整的里程碑版本

### 📤 Emoji 输出示例

```bash
📌 配置文件中未设置 'txtDownloadThreads'，自动设为默认值 5
📌 缓存中转机制已启用，缓存路径: ./Cache

📋 开始下载任务
📝 总数: 2, 并发数: 1

🧾 [1/2] 开始处理: https://music.apple.com/...
🎤 歌手: Anaïs Reno
💽 专辑: Lovesome Thing
🔬 正在进行版权预检，请稍候...
📡 音源: Lossless | 5 线程 | CN | 1 个账户并行下载

📤 正在从缓存转移文件到目标位置...
📥 文件转移完成！

📦 已完成: 25/25 | 警告: 0 | 错误: 0
```

---

## 🚀 使用方法

### 基本命令

```bash
# 查看帮助（完整中文说明）
goo --help

# 下载专辑
goo https://music.apple.com/cn/album/...

# 下载杜比全景声
goo --atmos https://music.apple.com/cn/album/...

# 下载单曲
goo --song https://music.apple.com/cn/album/...?i=...

# 下载播放列表
goo https://music.apple.com/cn/playlist/...

# 批量下载
goo url1 url2 url3
```

### 高级选项

```bash
# 纯日志模式（无动态UI）
goo --no-ui <url>

# 选择性下载
goo --select <album-url>

# 交互式搜索
goo --search song "歌曲名"
goo --search album "专辑名"
goo --search artist "歌手名"

# 下载歌手所有专辑
goo --all-album <artist-url>

# 调试模式（查看音质信息）
goo --debug <url>
```

### 音质选项

```bash
# ALAC 高解析度
goo --alac-max 192000 <url>

# Dolby Atmos
goo --atmos --atmos-max 2768 <url>

# AAC 格式
goo --aac --aac-type aac-lc <url>

# MV 下载
goo --mv-max 2160 <mv-url>
```

---

## ⚙️ 配置详情

### Shell 配置

**文件位置**: `~/.zshrc`

**配置内容**:
```bash
# 第148行
alias goo='/root/apple-music-downloader/apple-music-downloader-v2.2.0'
```

### 生效方式

#### 方式1: 重新加载配置（当前终端）
```bash
source ~/.zshrc
```

#### 方式2: 新开终端
新开的终端会自动加载配置，`goo` 命令直接可用

#### 方式3: 手动测试
```bash
alias goo='/root/apple-music-downloader/apple-music-downloader-v2.2.0'
```

---

## 🔄 版本对比

### 传统方式 vs goo 命令

| 对比项 | 传统方式 | goo 命令 |
|--------|---------|----------|
| **启动方式** | `go run main.go` | `goo` |
| **启动速度** | 慢（每次编译） | 快（直接运行） |
| **命令长度** | 15个字符 | 3个字符 |
| **版本** | 开发版 | v2.2.0 稳定版 |
| **包含功能** | 当前代码 | 完整v2.2.0功能 |

### 性能提升

- ⚡ **启动时间**: 从 ~2秒 → <0.1秒
- 📝 **命令简化**: 减少 80% 字符输入
- 🎯 **版本稳定**: 使用编译版本，避免代码变动影响

---

## 📝 维护说明

### 更新 goo 版本

当有新版本时，重新编译并更新别名：

```bash
# 1. 编译新版本
cd /root/apple-music-downloader
go build -o apple-music-downloader-v2.3.0 main.go

# 2. 更新别名
sed -i "s|apple-music-downloader-v2.2.0|apple-music-downloader-v2.3.0|g" ~/.zshrc

# 3. 重新加载
source ~/.zshrc

# 4. 验证
goo --help
```

### 回退到开发模式

如需使用 `go run` 方式测试代码：

```bash
# 临时回退
alias goo='go run main.go'

# 永久回退（修改 ~/.zshrc）
sed -i "s|/root/apple-music-downloader/apple-music-downloader-v2.2.0|go run main.go|g" ~/.zshrc
source ~/.zshrc
```

### 清理旧版本

```bash
# 查看所有二进制文件
ls -lh /root/apple-music-downloader/apple-music-downloader*

# 删除旧版本（保留当前版本）
rm /root/apple-music-downloader/apple-music-downloader.baseline
rm /root/apple-music-downloader/apple-music-downloader  # 如果不需要

# 只保留 v2.2.0
ls -lh /root/apple-music-downloader/apple-music-downloader-v2.2.0
```

---

## 🎯 快速参考

### 常用命令速查

```bash
# 帮助
goo --help

# 标准下载
goo <url>

# Atmos
goo --atmos <url>

# 纯日志
goo --no-ui <url>

# 搜索
goo --search song "关键词"
```

### 配置文件

```bash
# 编辑配置
nano ~/apple-music-downloader/config.yaml

# 查看示例
cat ~/apple-music-downloader/config.yaml.example
```

### 日志输出

```bash
# 保存日志
goo --no-ui <url> > download.log 2>&1

# 实时查看日志
goo <url> | tee download.log
```

---

## 💡 提示和技巧

### 1. 命令补全
```bash
# zsh 通常支持别名补全
goo --[Tab][Tab]  # 显示所有可用选项
```

### 2. 多任务下载
```bash
# 并行下载多个专辑
goo url1 url2 url3 url4 url5

# 配置文件设置并发数
# config.yaml: txtDownloadThreads: 5
```

### 3. CI/CD 集成
```bash
# Jenkins / GitLab CI
goo --no-ui <url> > output.log
if [ $? -eq 0 ]; then
  echo "下载成功"
else
  echo "下载失败"
  exit 1
fi
```

### 4. 快捷函数
```bash
# 添加到 ~/.zshrc
function goo-atmos() {
  goo --atmos "$@"
}

function goo-search() {
  goo --search song "$@"
}

# 使用
goo-atmos <url>
goo-search "歌名"
```

---

## 🔗 相关资源

### 文档

- [README.md](./README.md) - 主文档
- [README-CN.md](./README-CN.md) - 中文文档  
- [CHANGELOG.md](./CHANGELOG.md) - 版本历史
- [RELEASE_v2.2.0.md](./RELEASE_v2.2.0.md) - 发布说明
- [EMOJI_DEMO.md](./EMOJI_DEMO.md) - Emoji 演示

### 工具

- [wrapper](https://github.com/zhaarey/wrapper) - 解密服务
- [MP4Box](https://gpac.io/downloads/gpac-nightly-builds/) - 必需工具
- [mp4decrypt](https://www.bento4.com/downloads/) - MV 下载

---

## ✅ 验证清单

使用前请确认：

- [ ] MP4Box 已安装并在 PATH 中
- [ ] wrapper 解密服务正在运行
- [ ] config.yaml 已正确配置
- [ ] media-user-token 已填写（如需歌词/MV）
- [ ] goo 命令可正常执行
- [ ] 帮助菜单显示中文
- [ ] Emoji 正常显示

---

## 🎉 总结

**goo 命令别名配置成功！**

现在您可以：
- ✅ 使用简短的 `goo` 命令启动下载器
- ✅ 享受快速的启动速度（<0.1秒）
- ✅ 体验完整的中文帮助菜单
- ✅ 查看美观的 Emoji 输出
- ✅ 使用所有 v2.2.0 新功能

**开始您的高品质音乐下载之旅吧！** 🎵

---

**配置时间**: 2025-10-09  
**版本**: v2.2.0  
**状态**: ✅ 已生效

