# 🔨 构建指南

本文档说明如何使用 `build.sh` 脚本构建 Apple Music Downloader 二进制文件。

---

## 📦 构建脚本

### 文件位置

```
项目根目录/build.sh
```

### 功能特性

✅ **自动化构建流程**
- 自动清理旧版本二进制文件
- 收集版本信息（Git标签、提交哈希、构建时间）
- 编译生成优化的二进制文件
- 自动设置执行权限
- 验证构建结果

✅ **版本信息嵌入**
- Git标签（如 v2.6.0）
- Git提交哈希（短格式）
- 构建时间（UTC）
- Go版本

✅ **输出信息**
- 彩色终端输出
- 详细的构建进度
- 文件大小和路径信息
- 使用方法提示

---

## 🚀 使用方法

### 基本使用

```bash
# 在项目根目录执行
./build.sh
```

### 输出

构建完成后会在项目根目录生成：
- **文件名**: `apple-music-downloader`
- **大小**: 约 27 MB
- **权限**: 自动设置为可执行（755）

---

## 📋 构建流程

脚本会依次执行以下步骤：

### 1. 清理旧版本
```
🧹 清理旧版本...
✅ 已删除旧版本二进制文件
```

### 2. 收集版本信息
```
📝 收集版本信息...
   版本标签:   v2.6.0
   提交哈希:   cc03b2c
   构建时间:   2025-10-11 04:54:16 UTC
   Go 版本:    go1.24.8
```

### 3. 编译二进制
```
🔨 开始编译...
   目标平台:   linux/amd64
   输出文件:   apple-music-downloader

✅ 编译成功！
```

### 4. 显示构建结果
```
📦 二进制文件信息:
   文件名:     apple-music-downloader
   文件大小:   27M
   完整路径:   /root/apple-music-downloader/apple-music-downloader
   可执行:     ✓
```

### 5. 验证构建
```
🧪 验证构建...
✅ 验证成功，程序可正常运行
```

---

## ⚙️ 构建配置

### 编译参数

```bash
LDFLAGS="-s -w"
```

- `-s`: 去除符号表
- `-w`: 去除DWARF调试信息
- **效果**: 减小二进制文件大小（约减少30%）

### 版本信息注入

脚本会自动注入以下变量（需要在 `main.go` 中定义）：

```go
var (
    Version    string  // Git标签
    CommitHash string  // Git提交哈希
    BuildTime  string  // 构建时间
)
```

### 目标平台

- **操作系统**: Linux（当前环境）
- **架构**: amd64
- **自动检测**: 使用 `go env GOOS` 和 `go env GOARCH`

---

## 📁 输出文件

### 文件详情

| 属性 | 值 |
|------|-----|
| 文件名 | `apple-music-downloader` |
| 路径 | 项目根目录 |
| 大小 | ~27 MB |
| 权限 | 755 (rwxr-xr-x) |
| 格式 | ELF 64-bit LSB executable |

### 运行方式

```bash
# 直接运行
./apple-music-downloader

# 查看帮助
./apple-music-downloader --help

# 查看版本
./apple-music-downloader --version
```

---

## 🔧 高级用法

### 修改输出文件名

编辑 `build.sh`，修改以下变量：

```bash
PROJECT_NAME="apple-music-downloader"
BUILD_OUTPUT="${PROJECT_NAME}"
```

### 修改构建参数

编辑 `build.sh`，修改 LDFLAGS：

```bash
LDFLAGS="-s -w"  # 当前配置（优化大小）
# 或
LDFLAGS=""       # 保留调试信息
```

### 添加自定义标志

```bash
LDFLAGS="${LDFLAGS} -X 'main.CustomVar=value'"
```

---

## 🐛 故障排除

### 问题: 构建失败

**可能原因**:
- Go未安装或版本不兼容
- 缺少依赖包

**解决方法**:
```bash
# 检查Go版本
go version

# 下载依赖
go mod download

# 更新依赖
go mod tidy
```

### 问题: 权限被拒绝

**解决方法**:
```bash
# 添加执行权限
chmod +x build.sh

# 重新运行
./build.sh
```

### 问题: Git信息未找到

**现象**:
```
版本标签:   unknown
提交哈希:   unknown
```

**解决方法**:
```bash
# 确保在Git仓库中
git status

# 确保有提交记录
git log --oneline -1

# 创建标签（如果需要）
git tag v2.6.0
```

---

## 📊 构建统计

### 编译时间

| 环境 | CPU | 内存 | 时间 |
|------|-----|------|------|
| 示例 | 8核 | 16GB | ~5秒 |

### 文件大小对比

| 配置 | 大小 | 说明 |
|------|------|------|
| 带调试信息 | ~38 MB | `-ldflags ""` |
| 优化后 | ~27 MB | `-ldflags "-s -w"` |
| 减少 | 29% | 推荐使用优化配置 |

---

## ✅ 验证构建

### 检查文件

```bash
# 查看文件信息
ls -lh apple-music-downloader

# 查看文件类型
file apple-music-downloader

# 检查依赖
ldd apple-music-downloader
```

### 运行测试

```bash
# 查看帮助（验证程序可运行）
./apple-music-downloader --help

# 查看版本（验证版本信息）
./apple-music-downloader --version
```

---

## 📝 脚本维护

### 更新日志

| 版本 | 日期 | 更新内容 |
|------|------|----------|
| 1.0 | 2025-10-11 | 初始版本 |

### 兼容性

- ✅ Linux (amd64)
- ✅ macOS (需修改部分命令)
- ✅ Windows (WSL/Git Bash)

---

## 📞 支持

如遇问题，请：
1. 检查本文档的故障排除部分
2. 确认Go环境配置正确
3. 查看构建脚本输出的错误信息

---

**构建脚本版本**: 1.0  
**最后更新**: 2025-10-11  
**维护者**: Apple Music Downloader Team

