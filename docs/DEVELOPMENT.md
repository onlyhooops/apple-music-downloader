# Apple Music Downloader - 开发文档

**版本**: v1.3.1-dev  
**最后更新**: 2025-11-04

---

## 📖 文档导航

### 用户文档
- [README.md](../README.md) - 项目主文档（中文）
- [README-CN.md](../README-CN.md) - 中文说明文档
- [RELEASE-NOTES-v1.3.0.md](../RELEASE-NOTES-v1.3.0.md) - v1.3.0 发布说明

### 开发文档
- [DEVELOPMENT.md](./DEVELOPMENT.md) - 本文档（开发指南）
- [CODE_QUALITY.md](./CODE_QUALITY.md) - 代码质量报告（整合）

### 配置文档
- [config.yaml.example](../config.yaml.example) - 配置文件示例
- [dev.env.example](../dev.env.example) - 环境变量示例

---

## 🚀 快速开始

### 编译
```bash
# 标准编译
go build -o apple-music-downloader main.go

# 带版本信息编译
./build.sh
```

### 运行测试
```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./internal/core -v
go test ./internal/logger -v

# 运行基准测试
go test ./internal/logger -bench=. -benchmem
```

---

## 📊 项目架构

### 目录结构
```
apple-music-downloader/
├── main.go                 # 主入口
├── internal/               # 内部包
│   ├── api/               # Apple Music API 客户端
│   ├── config/            # 配置验证
│   ├── constants/         # 常量定义
│   ├── core/              # 核心状态管理
│   ├── downloader/        # 下载器
│   ├── logger/            # 日志系统
│   ├── metadata/          # 元数据写入
│   ├── parser/            # URL 解析
│   ├── progress/          # 进度管理
│   └── ui/                # 终端 UI
├── utils/                 # 工具包
│   ├── ampapi/           # Apple Music API 封装
│   ├── lyrics/           # 歌词处理
│   ├── runv3/            # DRM 处理
│   └── structs/          # 数据结构
└── docs/                 # 文档目录
```

### 核心模块

#### 1. API 客户端 (`internal/api`)
- `GetToken()` - 自动获取开发者 token
- `GetMeta()` - 获取专辑/播放列表元数据
- `GetInfoFromAdam()` - 获取曲目详细信息
- `CheckArtist()` - 获取艺术家作品列表

**关键改进**:
- ✅ 修复了循环中的资源泄漏（使用闭包立即关闭 HTTP 连接）
- ✅ 增强错误处理（使用 `%w` 包装错误）
- ✅ 添加调试日志（[API] 前缀）
- ✅ 修复 token 自动获取（支持新旧网站结构）

#### 2. 配置验证 (`internal/config`)
- `ValidateConfig()` - 验证配置文件完整性
- 覆盖 9 大类配置项
- 区分错误（阻止启动）和警告（提示用户）

**验证覆盖**:
1. 账号配置
2. 保存路径
3. 缓存配置
4. 下载性能
5. 工作-休息循环
6. 音质配置
7. MV 配置
8. 路径限制
9. 日志配置

#### 3. 核心状态 (`internal/core`)
- 全局配置管理
- 虚拟 Singles 专辑支持
- RWMutex 并发优化

**关键功能**:
- `IsSingleAlbum()` - 单曲专辑识别
- `GetPrimaryArtist()` - 主艺术家提取
- `GetVirtualSinglesTrackNumber()` - 虚拟 Singles 编号分配

#### 4. 下载器 (`internal/downloader`)
- 缓存机制支持
- 多线程下载
- FFmpeg 修复集成
- Context 取消支持

---

## 🔧 配置系统

### 配置文件 (`config.yaml`)
详细配置说明请参考 `config.yaml.example`

### 环境变量 (`dev.env`)
```bash
APPLE_MUSIC_MEDIA_USER_TOKEN_CN=你的token
# authorization-token 会自动获取，无需配置
```

### 命令行参数
```bash
--atmos              # 启用 Dolby Atmos
--aac                # 启用 AAC 模式
--alac-max 192000    # 设置 ALAC 最大采样率
--mv-max 2160        # 设置 MV 最大分辨率
--config path        # 指定配置文件
--cx                 # 强制覆盖
--debug              # 调试模式
```

---

## 🧪 测试

### 单元测试覆盖

#### core 模块
- `TestIsSingleAlbum` - 单曲识别逻辑（8 个子测试）
- `TestGetPrimaryArtist` - 艺术家提取（6 个子测试）
- `TestGetVirtualSinglesTrackNumber` - 编号分配

#### logger 模块
- 8 个测试函数
- 完整的并发测试
- 基准测试支持

#### progress 模块
- 8 个测试函数
- 并发场景覆盖

**测试通过率**: 100% (19/19)

---

## 🎯 代码质量

### 最近改进（v1.3.1）

#### 1. 资源泄漏修复
```go
// 修复前：defer 在循环外部
for ... {
    resp, _ := http.Do(req)
    defer resp.Body.Close()  // ❌ 泄漏
}

// 修复后：使用闭包
for ... {
    err := func() error {
        defer resp.Body.Close()  // ✅ 立即释放
        // ... 处理逻辑
        return nil
    }()
}
```

#### 2. Context 取消机制
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// 监听中断信号
go func() {
    <-sigChan
    cancel()
}()

// 在循环中检查
select {
case <-ctx.Done():
    return
default:
    // 继续处理
}
```

#### 3. 并发优化
```go
// 从 sync.Mutex 升级到 sync.RWMutex
var virtualSinglesLock sync.RWMutex

// 读操作使用 RLock
virtualSinglesLock.RLock()
defer virtualSinglesLock.RUnlock()

// 写操作使用 Lock
virtualSinglesLock.Lock()
defer virtualSinglesLock.Unlock()
```

#### 4. 常量定义
所有魔法数字已移至 `internal/constants/constants.go`

---

## 📝 代码规范

### 命名规范
- 包名：小写单词
- 导出函数：大写开头驼峰
- 私有函数：小写开头驼峰
- 常量：大写开头驼峰或全大写

### 错误处理
```go
// ✅ 推荐：使用 %w 包装错误
if err != nil {
    return fmt.Errorf("操作失败: %w", err)
}

// ❌ 避免：丢失错误链
if err != nil {
    return errors.New("操作失败")
}
```

### 日志规范
```go
// 使用统一的模块前缀
logger.Debug("[API] 请求成功: %s", url)
logger.Debug("[虚拟Singles] 识别为单曲: %s", name)
logger.Debug("[工作-休息] 循环已启用")
```

---

## 🐛 已知问题

### 已修复
- ✅ 循环中的资源泄漏
- ✅ Token 自动获取失败
- ✅ 虚拟 Singles 合作艺术家处理
- ✅ 工作-休息循环日志缺失
- ✅ 配置参数验证缺失

### 待处理
- ⏸️ 长期重构：拆分过长函数
- ⚠️ MV 下载分辨率门槛（部分场景）
- ⚠️ 专辑中 MV 同时下载的稳定性

---

## 🔄 发布流程

### 版本号规则
遵循语义化版本 (SemVer): `major.minor.patch`

### 发布检查清单
- [ ] 所有测试通过
- [ ] 代码无 linter 错误
- [ ] 更新 CHANGELOG
- [ ] 更新版本号
- [ ] 创建 Git tag
- [ ] 编译发布版本

```bash
# 1. 运行测试
go test ./...

# 2. 编译
./build.sh

# 3. 创建 tag
git tag v1.3.1 -m "版本描述"

# 4. 推送
git push origin v1.3.1
```

---

## 🤝 贡献指南

### 提交规范
使用 Conventional Commits 格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型 (type)**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建/工具变动

**示例**:
```
feat(config): 添加 MV 分辨率下限配置

- 新增 mv-min 配置项
- 支持区间过滤（如 1080p ~ 2160p）
- 添加配置验证

Closes #123
```

---

## 📚 参考资源

### 内部文档
- [代码质量报告](./CODE_QUALITY.md)
- [配置参数说明](../config.yaml.example)

### 外部资源
- [Go 官方文档](https://golang.org/doc/)
- [Apple Music API](https://developer.apple.com/documentation/applemusicapi)

---

**维护者**: Apple Music Downloader Team  
**最后更新**: 2025-11-04

