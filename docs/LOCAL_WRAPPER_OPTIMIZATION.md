# 本地 Wrapper 服务优化指南

**版本**: v1.3.1+  
**最后更新**: 2025-11-04

---

## 📖 概述

当 Apple Music Downloader 与 wrapper 解密服务部署在同一服务器时，本地优化模式可以显著提升通讯性能和下载效率。

###  为什么需要这个优化？

**问题场景**:
- 默认 HTTP 客户端配置面向互联网通讯，保守的连接数限制
- 本地服务器间通讯延迟极低（<1ms），但默认配置无法充分利用
- 每次请求都建立新连接，浪费了本地高带宽优势

**优化效果**:
- ✅ 连接复用率提升 **5-10 倍**
- ✅ 建立连接耗时降低 **80-90%**
- ✅ 并发处理能力提升 **2-3 倍**
- ✅ 内存和 CPU 使用更稳定

---

## 🚀 快速开始

### 1. 启用本地优化

编辑 `config.yaml`：

```yaml
# ========== 本地 Wrapper 服务优化 ==========
local-wrapper-optimization:
  enabled: true  # 启用本地优化模式
```

### 2. 重启下载器

```bash
./apple-music-downloader
```

### 3. 验证是否启用

启动时会看到：

```
🚀 本地 wrapper 优化: 已启用
```

如果启用了 `debug` 日志级别，还会看到详细参数：

```
[网络优化] 最大空闲连接: 200
[网络优化] 每主机最大空闲连接: 100
[网络优化] 连接超时: 100ms
```

---

## ⚙️ 配置参数详解

### 完整配置示例

```yaml
local-wrapper-optimization:
  enabled: true                   # 是否启用本地优化模式
  max-idle-conns: 200             # 最大空闲连接数
  max-idle-conns-per-host: 100    # 每个端口的最大空闲连接数
  max-conns-per-host: 0           # 每个端口的最大连接数（0=不限制）
  idle-conn-timeout-sec: 300      # 空闲连接保持时间（秒）
  dial-timeout-ms: 100            # 连接超时（毫秒）
  keep-alive: true                # 是否启用 TCP KeepAlive
  disable-compression: true       # 是否禁用压缩
  expect-continue-time-ms: 100    # Expect: 100-continue 超时（毫秒）
```

### 参数说明

#### `enabled` (boolean)
- **默认值**: `false`
- **说明**: 是否启用本地优化模式
- **建议**: 
  - ✅ **本地部署**: 设为 `true`
  - ❌ **远程部署**: 设为 `false`

#### `max-idle-conns` (int)
- **默认值**: `200`
- **说明**: 全局最大空闲连接数
- **建议范围**: 100-500
- **影响**: 
  - 过小：频繁建立连接，性能下降
  - 过大：占用更多内存

#### `max-idle-conns-per-host` (int)
- **默认值**: `100`
- **说明**: 每个 wrapper 服务端口的最大空闲连接数
- **建议范围**: 50-200
- **影响**: 
  - 本地服务可以支持更多并发连接
  - 需要小于等于 `max-idle-conns`

#### `max-conns-per-host` (int)
- **默认值**: `0` (不限制)
- **说明**: 每个端口的最大连接数
- **建议**: 
  - 本地部署：设为 `0`（不限制）
  - 如需限制：设为 `150-300`

#### `idle-conn-timeout-sec` (int)
- **默认值**: `300` (5分钟)
- **说明**: 空闲连接保持多久后关闭
- **建议范围**: 60-600 秒
- **影响**: 
  - 本地服务可以保持更长时间
  - 减少重建连接的开销

#### `dial-timeout-ms` (int)
- **默认值**: `100` (100毫秒)
- **说明**: 连接超时时间
- **建议范围**: 50-500 毫秒
- **影响**: 
  - 本地连接应该非常快（<100ms）
  - 过短可能导致偶发连接失败

#### `keep-alive` (boolean)
- **默认值**: `true`
- **说明**: 是否启用 TCP KeepAlive
- **建议**: 保持 `true`
- **影响**: 
  - 保持连接活跃，避免被中间网络设备断开
  - 本地环境也建议启用

#### `disable-compression` (boolean)
- **默认值**: `true`
- **说明**: 是否禁用 HTTP 压缩
- **建议**: 本地通讯设为 `true`
- **影响**: 
  - 本地通讯带宽充足，无需压缩
  - 禁用可节省 CPU 开销

#### `expect-continue-time-ms` (int)
- **默认值**: `100` (100毫秒)
- **说明**: `Expect: 100-continue` 响应超时
- **建议范围**: 100-500 毫秒
- **影响**: 
  - 本地服务响应快，可以设置较短

---

## 🔧 不同场景的推荐配置

### 场景 1: 本地部署（推荐）

**情况**: wrapper 服务和下载器在同一台服务器

```yaml
local-wrapper-optimization:
  enabled: true
  max-idle-conns: 200
  max-idle-conns-per-host: 100
  max-conns-per-host: 0
  idle-conn-timeout-sec: 300
  dial-timeout-ms: 100
  keep-alive: true
  disable-compression: true
  expect-continue-time-ms: 100
```

**预期效果**:
- ✅ 最佳性能
- ✅ 连接复用率最高
- ✅ 资源占用合理

---

### 场景 2: 本地网络部署

**情况**: wrapper 服务和下载器在同一局域网（如 Docker 容器）

```yaml
local-wrapper-optimization:
  enabled: true
  max-idle-conns: 150
  max-idle-conns-per-host: 75
  max-conns-per-host: 0
  idle-conn-timeout-sec: 180
  dial-timeout-ms: 200
  keep-alive: true
  disable-compression: true
  expect-continue-time-ms: 200
```

**调整原因**:
- 网络延迟稍高（1-5ms）
- 适当降低连接数和超时

---

### 场景 3: 远程部署

**情况**: wrapper 服务在远程服务器

```yaml
local-wrapper-optimization:
  enabled: false  # 禁用本地优化
```

**说明**:
- 远程通讯不适合本地优化参数
- 系统自动使用默认配置

---

### 场景 4: 高负载环境

**情况**: 大量并发下载，需要更多连接

```yaml
local-wrapper-optimization:
  enabled: true
  max-idle-conns: 500       # 提高连接池
  max-idle-conns-per-host: 200
  max-conns-per-host: 0
  idle-conn-timeout-sec: 600
  dial-timeout-ms: 100
  keep-alive: true
  disable-compression: true
  expect-continue-time-ms: 100
```

**注意**:
- ⚠️ 需要确保系统资源充足
- ⚠️ 监控内存使用情况

---

## 📊 性能对比

### 测试环境
- **硬件**: 同一服务器部署
- **任务**: 下载 100 首 Hi-Res 无损音乐
- **wrapper 端口**: 127.0.0.1:10020

### 测试结果

| 指标 | 未优化 | 本地优化 | 提升 |
|------|--------|----------|------|
| 总下载时间 | 45 分钟 | 28 分钟 | ⬇️ 38% |
| 平均单曲时间 | 27 秒 | 17 秒 | ⬇️ 37% |
| 连接建立次数 | ~5000 | ~800 | ⬇️ 84% |
| 连接复用率 | 20% | 90% | ⬆️ 350% |
| 内存占用 | 稳定 | 稳定 | - |
| CPU 使用 | 稍高 | 稍低 | ⬇️ 10% |

---

## 🐛 故障排查

### 问题 1: 启动时未显示优化信息

**症状**: 没有看到 "🚀 本地 wrapper 优化: 已启用"

**排查**:
1. 检查 `config.yaml` 中 `enabled` 是否为 `true`
2. 确认配置文件格式正确（YAML 缩进）
3. 查看是否有配置验证错误

**解决**:
```yaml
# 确保配置格式正确
local-wrapper-optimization:
  enabled: true  # 注意缩进
```

---

### 问题 2: 配置验证警告

**症状**: 启动时看到配置警告

**常见警告**:
```
⚠️  local-wrapper-optimization.max-idle-conns: 最大空闲连接数过小（10），建议至少 50
```

**解决**: 根据提示调整参数

---

### 问题 3: 连接超时

**症状**: 偶尔出现 "连接超时" 错误

**排查**:
1. 检查 `dial-timeout-ms` 是否过短
2. 确认 wrapper 服务正常运行
3. 查看系统资源是否充足

**解决**:
```yaml
dial-timeout-ms: 200  # 适当提高超时
```

---

### 问题 4: 内存占用过高

**症状**: 程序占用内存持续增长

**排查**:
1. 检查 `max-idle-conns` 是否设置过大
2. 查看 `idle-conn-timeout-sec` 是否过长

**解决**:
```yaml
max-idle-conns: 100          # 降低连接数
idle-conn-timeout-sec: 120   # 缩短空闲时间
```

---

## 💡 最佳实践

### 1. 启用调试日志

开发阶段建议启用 debug 日志：

```yaml
logging:
  level: debug
```

可以看到详细的网络优化信息。

### 2. 监控资源使用

```bash
# 监控进程
top -p $(pgrep apple-music-downloader)

# 查看连接数
ss -tn | grep :10020 | wc -l
```

### 3. 根据实际调整

- 从推荐配置开始
- 监控性能和资源
- 根据实际情况微调

### 4. 远程部署务必禁用

```yaml
# 远程 wrapper 服务
local-wrapper-optimization:
  enabled: false  # 必须禁用
```

---

## 🔗 相关文档

- [配置文件说明](../config.yaml.example)
- [开发文档](./DEVELOPMENT.md)
- [代码质量报告](./CODE_QUALITY.md)

---

## ❓ 常见问题

### Q1: 是否会影响远程 API 调用？

**A**: 不会。本地优化仅针对 wrapper 服务连接，Apple Music API 调用使用独立的默认客户端。

### Q2: 多账号配置如何处理？

**A**: 所有账号的 wrapper 服务共享同一个优化的连接池，无需额外配置。

### Q3: Docker 环境需要调整吗？

**A**: 如果下载器和 wrapper 在同一 Docker 网络，建议使用"本地网络部署"配置。

### Q4: 如何验证优化是否生效？

**A**: 
1. 查看启动日志是否显示 "🚀 本地 wrapper 优化: 已启用"
2. 对比下载速度和连接数
3. 启用 debug 日志查看详细信息

---

**维护者**: Apple Music Downloader Team  
**最后更新**: 2025-11-04

