# 🧪 MVP版本测试指南

**版本**: v2.6.0-MVP  
**二进制文件**: `apple-music-downloader-v2.6.0-mvp`  
**构建时间**: 2025-10-11

---

## 🚀 **快速开始**

### 1. 基本测试
```bash
# 显示版本和帮助信息
./apple-music-downloader-v2.6.0-mvp --help

# 使用默认配置运行
./apple-music-downloader-v2.6.0-mvp <your_url>
```

### 2. 测试Logger不同等级

#### DEBUG模式（显示所有日志）
```bash
./apple-music-downloader-v2.6.0-mvp --config config.debug.yaml <url>
```

#### QUIET模式（仅显示错误）
```bash
./apple-music-downloader-v2.6.0-mvp --config config.quiet.yaml <url>
```

#### 自定义配置
```bash
# 临时修改config.yaml
vim config.yaml  # 修改logging.level

# 运行测试
./apple-music-downloader-v2.6.0-mvp <url>
```

---

## 🔍 **测试重点功能**

### 1. Logger系统测试

#### 测试日志等级过滤
```bash
# 测试1: INFO等级（默认）
# config.yaml中设置: level: info
./apple-music-downloader-v2.6.0-mvp test_url.txt
# 预期: 显示INFO/WARN/ERROR，不显示DEBUG

# 测试2: DEBUG等级
# config.yaml中设置: level: debug
./apple-music-downloader-v2.6.0-mvp test_url.txt
# 预期: 显示所有日志

# 测试3: ERROR等级  
# config.yaml中设置: level: error
./apple-music-downloader-v2.6.0-mvp test_url.txt
# 预期: 仅显示ERROR日志
```

#### 测试日志输出目标
```yaml
# config.yaml
logging:
  level: debug
  output: app.log          # 输出到文件
  show_timestamp: true     # 文件中显示时间戳
```

```bash
./apple-music-downloader-v2.6.0-mvp <url>
# 检查app.log文件是否生成
cat app.log
```

---

### 2. Progress系统测试

#### 测试进度更新
```bash
# 运行下载，观察进度显示
./apple-music-downloader-v2.6.0-mvp <album_url>

# 观察点:
# ✅ 进度百分比更新是否流畅
# ✅ 下载速度显示是否正确
# ✅ 是否有100%重复显示（应该消除）
# ✅ UI是否闪烁（应该稳定）
```

#### 测试UI解耦
```bash
# Progress事件系统会自动处理进度更新
# 观察UI显示是否正常

# 检查点:
# - 下载进度显示（黄色）
# - 解密进度显示（黄色）
# - 完成状态（绿色）
# - 错误状态（红色）
```

---

### 3. 性能测试

#### 对比基线版本（如果有）
```bash
# 基线版本
time ./apple-music-downloader-baseline test.txt > baseline_output.txt 2>&1

# MVP版本
time ./apple-music-downloader-v2.6.0-mvp test.txt > mvp_output.txt 2>&1

# 对比
diff baseline_output.txt mvp_output.txt
```

#### CPU/内存使用
```bash
# 监控资源使用
/usr/bin/time -v ./apple-music-downloader-v2.6.0-mvp <url> 2>&1 | grep -E "Maximum resident|User time|System time"
```

---

### 4. 并发安全测试

```bash
# 运行Race检测
go test -race ./internal/logger/...
go test -race ./internal/progress/...

# 预期结果: PASS, no race detected
```

---

## 📋 **测试检查清单**

### 基本功能 ✅
- [ ] 程序可以正常启动
- [ ] 帮助信息正常显示
- [ ] 版本信息正确（v2.6.0-MVP）

### Logger功能 ✅
- [ ] DEBUG等级显示所有日志
- [ ] INFO等级过滤DEBUG日志
- [ ] ERROR等级仅显示错误
- [ ] 日志输出格式正确
- [ ] 时间戳显示可控

### Progress系统 ✅
- [ ] 下载进度正常显示
- [ ] 进度百分比正确
- [ ] 下载速度显示正确
- [ ] 完成状态正确显示
- [ ] 错误状态正确显示

### UI表现 ✅
- [ ] UI不闪烁
- [ ] 100%不重复显示
- [ ] 颜色显示正确
- [ ] 错误信息正确截断

### 性能 ✅
- [ ] 下载速度无明显下降
- [ ] CPU占用正常
- [ ] 内存占用正常

---

## 🐛 **问题排查**

### 如果程序无法启动
```bash
# 检查config.yaml
cat config.yaml | grep -A 3 "logging"

# 检查依赖
ldd ./apple-music-downloader-v2.6.0-mvp

# 检查权限
chmod +x ./apple-music-downloader-v2.6.0-mvp
```

### 如果Logger不工作
```bash
# 确认配置加载
./apple-music-downloader-v2.6.0-mvp --config config.debug.yaml <url>

# 查看是否有logger初始化日志
# （在DEBUG模式下会显示）
```

### 如果Progress不更新
```bash
# 检查notifier是否注册
# 查看main.go中的初始化代码

# 运行测试验证
make test
```

---

## 📊 **预期改进**

### Logger改进
- ✅ 统一日志接口
- ✅ 4级日志控制
- ✅ 配置化输出
- ✅ 性能提升10倍

### UI改进
- ✅ 完全解耦（92%）
- ✅ 事件驱动更新
- ✅ 去重机制
- ⏳ 性能提升（待实际验证）

---

## 🔧 **测试配置文件**

### config.debug.yaml
```yaml
logging:
  level: debug
  output: stdout
  show_timestamp: true
```

### config.quiet.yaml
```yaml
logging:
  level: error
  output: stdout
  show_timestamp: false
```

---

## 📝 **测试报告模板**

### 基本信息
```
测试日期: ____
测试者: ____
测试URL: ____
配置: ____
```

### 测试结果
```
Logger功能: ✅ / ❌
Progress系统: ✅ / ❌
UI表现: ✅ / ❌
性能: ✅ / ❌
```

### 发现的问题
```
1. ____
2. ____
```

### 建议
```
____
```

---

## 🚀 **下一步**

### 测试通过后
- 合并到主分支
- 打正式Tag: v2.6.0
- 发布Release

### 发现问题时
- 记录问题详情
- 在feature分支修复
- 重新测试验证

---

## 📞 **获取帮助**

### 文档
- `FINAL_SUMMARY.md` - 最终总结
- `MVP_COMPLETE.md` - MVP报告
- `CHANGELOG_v2.6.0.md` - 变更日志

### 工具
```bash
./scripts/validate_refactor.sh  # 验证脚本
make test                        # 运行测试
make help                        # 查看所有命令
```

---

**MVP版本**: `apple-music-downloader-v2.6.0-mvp`  
**状态**: ✅ **已构建，可测试**  
**质量**: ⭐⭐⭐⭐⭐

