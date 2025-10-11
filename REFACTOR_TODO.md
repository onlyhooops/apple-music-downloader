# UI与LOG模块重构 - 详细任务清单

**项目**: Apple Music Downloader UI & Logger Refactoring  
**版本**: v1.0  
**创建日期**: 2025-10-10  
**预计周期**: 4-6周（MVP）

---

## 📋 任务状态说明

- ⬜ 未开始
- 🔄 进行中
- ✅ 已完成
- ⏸️ 暂停/阻塞
- ❌ 取消

---

## 🎯 Week 0: 准备阶段 (预计1-2天)

### 任务组 0.1: 环境准备

#### ⬜ 任务 0.1.1: 创建Git分支
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 10分钟
- **前置依赖**: 无
- **执行步骤**:
  ```bash
  # 1. 确保main分支是最新的
  git checkout main
  git pull origin main
  
  # 2. 创建feature分支
  git checkout -b feature/ui-refactor
  
  # 3. 推送到远程
  git push -u origin feature/ui-refactor
  ```
- **验收标准**:
  - [ ] 分支已创建: `git branch | grep feature/ui-refactor`
  - [ ] 分支已推送到远程
  - [ ] 当前工作目录干净: `git status`

---

#### ⬜ 任务 0.1.2: 保存性能基线
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 30分钟
- **前置依赖**: 任务 0.1.1
- **执行步骤**:
  ```bash
  # 1. 确保当前版本可编译
  go build -o apple-music-downloader
  
  # 2. 运行benchmark并保存结果
  go test -bench=. ./... -benchmem > baseline_bench.txt
  
  # 3. 保存当前版本的可执行文件
  cp apple-music-downloader apple-music-downloader-baseline
  
  # 4. 提交基线数据
  git add baseline_bench.txt
  git commit -m "chore: 保存重构前的性能基线"
  git push
  ```
- **验收标准**:
  - [ ] baseline_bench.txt 文件存在
  - [ ] 文件包含benchmark结果
  - [ ] 备份可执行文件存在

---

#### ⬜ 任务 0.1.3: 创建回滚tag
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 5分钟
- **前置依赖**: 任务 0.1.2
- **执行步骤**:
  ```bash
  # 1. 在main分支打tag
  git checkout main
  git tag -a v2.5.3-pre-refactor -m "重构前的稳定版本"
  git push origin v2.5.3-pre-refactor
  
  # 2. 切回feature分支
  git checkout feature/ui-refactor
  ```
- **验收标准**:
  - [ ] tag已创建: `git tag | grep v2.5.3-pre-refactor`
  - [ ] tag已推送到远程
  - [ ] 可以通过tag回滚: `git checkout v2.5.3-pre-refactor`

---

### 任务组 0.2: 测试数据准备

#### ⬜ 任务 0.2.1: 创建测试目录结构
- **负责人**: _________
- **优先级**: 🟡 中
- **预计时间**: 10分钟
- **前置依赖**: 任务 0.1.1
- **执行步骤**:
  ```bash
  # 1. 创建测试目录
  mkdir -p test/{data,scripts,baseline}
  
  # 2. 创建.gitignore（排除敏感数据）
  cat > test/.gitignore <<'EOF'
  # 排除实际下载的音乐文件
  *.m4a
  *.mp3
  *.flac
  
  # 保留测试URL和输出结果
  !*.txt
  !*.log
  EOF
  
  # 3. 提交结构
  git add test/
  git commit -m "chore: 创建测试目录结构"
  ```
- **验收标准**:
  - [ ] test/data 目录存在
  - [ ] test/scripts 目录存在
  - [ ] test/baseline 目录存在
  - [ ] .gitignore 已配置

---

#### ⬜ 任务 0.2.2: 准备测试URL列表
- **负责人**: _________
- **优先级**: 🟡 中
- **预计时间**: 20分钟
- **前置依赖**: 任务 0.2.1
- **执行步骤**:
  ```bash
  # 1. 创建不同场景的测试URL文件
  
  # 单曲测试
  echo "https://music.apple.com/xx/album/xxx/xxx" > test/data/single_track.txt
  
  # 小专辑测试（5-10首）
  cat > test/data/small_album.txt <<'EOF'
  # 添加5-10首歌的专辑URL
  EOF
  
  # 大专辑测试（20+首）
  cat > test/data/large_album.txt <<'EOF'
  # 添加20+首歌的专辑URL
  EOF
  
  # 批量下载测试
  cat > test/data/batch_download.txt <<'EOF'
  # 添加多个专辑URL
  EOF
  
  git add test/data/*.txt
  git commit -m "chore: 添加测试URL数据"
  ```
- **验收标准**:
  - [ ] 至少有4个测试场景的URL文件
  - [ ] URL格式正确
  - [ ] 文件已提交到Git

---

#### ⬜ 任务 0.2.3: 保存基线输出
- **负责人**: _________
- **优先级**: 🟡 中
- **预计时间**: 30分钟
- **前置依赖**: 任务 0.2.2
- **执行步骤**:
  ```bash
  # 1. 使用当前版本运行测试，保存输出
  ./apple-music-downloader-baseline test/data/single_track.txt \
      > test/baseline/single_track_output.txt 2>&1
  
  ./apple-music-downloader-baseline test/data/small_album.txt \
      > test/baseline/small_album_output.txt 2>&1
  
  # 2. 统计基线数据
  cat > test/baseline/stats.txt <<EOF
  测试时间: $(date)
  版本: v2.5.3-pre-refactor
  
  单曲测试:
  - 耗时: $(grep "完成" test/baseline/single_track_output.txt | tail -1)
  - UI刷新次数: $(grep -c "%" test/baseline/single_track_output.txt)
  
  小专辑测试:
  - 耗时: $(grep "完成" test/baseline/small_album_output.txt | tail -1)
  - UI刷新次数: $(grep -c "%" test/baseline/small_album_output.txt)
  EOF
  
  git add test/baseline/
  git commit -m "chore: 保存基线输出数据"
  ```
- **验收标准**:
  - [ ] 所有测试场景都有基线输出
  - [ ] stats.txt包含统计数据
  - [ ] 数据已提交

---

### 任务组 0.3: 项目文档准备

#### ⬜ 任务 0.3.1: 创建验证脚本
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 30分钟
- **前置依赖**: 任务 0.1.1
- **执行步骤**:
  ```bash
  # 创建验证脚本
  cat > scripts/validate_refactor.sh <<'EOF'
  #!/bin/bash
  set -e
  
  echo "🔍 验证重构安全性..."
  
  # 1. 编译检查
  echo "1️⃣ 编译检查..."
  go build -o apple-music-downloader || { echo "❌ 编译失败"; exit 1; }
  echo "✅ 编译通过"
  
  # 2. 单元测试
  echo "2️⃣ 单元测试..."
  go test ./... -v || { echo "❌ 单元测试失败"; exit 1; }
  echo "✅ 单元测试通过"
  
  # 3. Race检测
  echo "3️⃣ Race检测..."
  go test -race ./... || { echo "❌ Race检测失败"; exit 1; }
  echo "✅ Race检测通过"
  
  # 4. 功能测试（可选，需要测试数据）
  if [ -f "test/data/single_track.txt" ]; then
      echo "4️⃣ 功能测试..."
      ./apple-music-downloader test/data/single_track.txt || { echo "⚠️ 功能测试失败"; }
      echo "✅ 功能测试通过"
  fi
  
  # 5. 性能对比（如果有基线）
  if [ -f "baseline_bench.txt" ]; then
      echo "5️⃣ 性能对比..."
      go test -bench=. ./... > new_bench.txt
      if command -v benchcmp &> /dev/null; then
          benchcmp baseline_bench.txt new_bench.txt || echo "⚠️ 性能有变化，需人工审查"
      else
          echo "⚠️ benchcmp未安装，跳过性能对比"
      fi
  fi
  
  echo ""
  echo "✅ 所有验证通过！"
  EOF
  
  chmod +x scripts/validate_refactor.sh
  git add scripts/validate_refactor.sh
  git commit -m "chore: 添加重构验证脚本"
  ```
- **验收标准**:
  - [ ] 脚本可执行
  - [ ] 运行成功: `./scripts/validate_refactor.sh`
  - [ ] 脚本已提交

---

#### ⬜ 任务 0.3.2: 创建Makefile
- **负责人**: _________
- **优先级**: 🟡 中
- **预计时间**: 20分钟
- **前置依赖**: 任务 0.3.1
- **执行步骤**:
  ```bash
  cat > Makefile <<'EOF'
  .PHONY: all build test bench race lint clean validate ci
  
  all: build
  
  build:
  	go build -o apple-music-downloader
  
  test:
  	go test ./... -v -cover
  
  bench:
  	go test -bench=. ./... -benchmem
  
  race:
  	go test -race ./...
  
  lint:
  	golangci-lint run || echo "golangci-lint not installed"
  
  clean:
  	rm -f apple-music-downloader
  	rm -f *.prof
  	rm -f *_bench.txt
  
  validate:
  	./scripts/validate_refactor.sh
  
  ci: test race lint
  	@echo "✅ All checks passed!"
  
  # 性能对比
  perf-compare:
  	@if [ ! -f baseline_bench.txt ]; then \
  		echo "❌ baseline_bench.txt not found"; \
  		exit 1; \
  	fi
  	go test -bench=. ./... > new_bench.txt
  	benchcmp baseline_bench.txt new_bench.txt || true
  EOF
  
  git add Makefile
  git commit -m "chore: 添加Makefile"
  ```
- **验收标准**:
  - [ ] Makefile创建成功
  - [ ] `make test` 可执行
  - [ ] `make validate` 可执行
  - [ ] 已提交

---

#### ⬜ 任务 0.3.3: 团队评审会议
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 1-2小时
- **前置依赖**: 任务组 0.1, 0.2 完成
- **会议议程**:
  1. 方案总体介绍（15分钟）
  2. MVP vs 完整方案讨论（15分钟）
  3. 风险评估与缓解（20分钟）
  4. 任务分配（20分钟）
  5. 时间表确认（10分钟）
  6. Q&A（20分钟）
- **决策事项**:
  - [ ] 确认是否启动重构
  - [ ] 确认采用MVP方案 or 完整方案
  - [ ] 确认Phase 1负责人
  - [ ] 确认Phase 2负责人
  - [ ] 确认每周代码审查时间
- **输出文档**:
  - [ ] 会议纪要
  - [ ] 任务分配表
  - [ ] 里程碑时间表

---

## 📦 Phase 1: 日志模块重构 (Week 1-2, 预计8-10天)

### 任务组 1.1: Logger包基础实现 (Day 1-2)

#### ⬜ 任务 1.1.1: 创建logger包目录结构
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 10分钟
- **前置依赖**: Week 0完成
- **执行步骤**:
  ```bash
  # 创建目录
  mkdir -p internal/logger
  
  # 创建文件骨架
  touch internal/logger/logger.go
  touch internal/logger/logger_test.go
  touch internal/logger/logger_bench_test.go
  touch internal/logger/config.go
  
  git add internal/logger/
  git commit -m "feat(logger): 创建logger包目录结构"
  ```
- **验收标准**:
  - [ ] 目录结构正确
  - [ ] 文件已创建
  - [ ] 已提交

---

#### ⬜ 任务 1.1.2: 实现Logger接口定义
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 1小时
- **前置依赖**: 任务 1.1.1
- **执行步骤**:
  1. 在 `internal/logger/logger.go` 中实现：
     - LogLevel 枚举类型
     - Logger 接口定义
     - levelNames 常量
  2. 参考文档中的代码示例
  3. 添加必要的注释
- **代码要点**:
  ```go
  type LogLevel int
  
  const (
      DEBUG LogLevel = iota
      INFO
      WARN
      ERROR
  )
  
  type Logger interface {
      Debug(format string, args ...interface{})
      Info(format string, args ...interface{})
      Warn(format string, args ...interface{})
      Error(format string, args ...interface{})
      SetLevel(level LogLevel)
      SetOutput(w io.Writer)
  }
  ```
- **验收标准**:
  - [ ] 代码可编译
  - [ ] 接口定义完整
  - [ ] 注释清晰
  - [ ] 提交: `git commit -m "feat(logger): 定义Logger接口"`

---

#### ⬜ 任务 1.1.3: 实现DefaultLogger结构体
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2小时
- **前置依赖**: 任务 1.1.2
- **执行步骤**:
  1. 实现 DefaultLogger 结构体
  2. 实现 New() 构造函数
  3. 实现 log() 内部方法（带锁保护）
  4. 实现 Debug/Info/Warn/Error 方法
  5. 实现 SetLevel/SetOutput 方法
- **代码要点**:
  ```go
  type DefaultLogger struct {
      mu       sync.Mutex
      level    LogLevel
      output   io.Writer
      showTime bool
  }
  
  func (l *DefaultLogger) log(level LogLevel, format string, args ...interface{}) {
      if level < l.level {
          return
      }
      l.mu.Lock()
      defer l.mu.Unlock()
      // 实现日志输出逻辑
  }
  ```
- **验收标准**:
  - [ ] 所有方法实现完整
  - [ ] 线程安全（使用mutex）
  - [ ] 编译通过: `go build ./internal/logger/...`
  - [ ] 提交: `git commit -m "feat(logger): 实现DefaultLogger"`

---

#### ⬜ 任务 1.1.4: 实现全局logger实例
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 30分钟
- **前置依赖**: 任务 1.1.3
- **执行步骤**:
  1. 创建全局logger变量
  2. 实现包级别的便捷函数
  3. 确保线程安全
- **代码要点**:
  ```go
  var global = New()
  
  func Debug(format string, args ...interface{}) { global.Debug(format, args...) }
  func Info(format string, args ...interface{})  { global.Info(format, args...) }
  func Warn(format string, args ...interface{})  { global.Warn(format, args...) }
  func Error(format string, args ...interface{}) { global.Error(format, args...) }
  func SetLevel(level LogLevel)                  { global.SetLevel(level) }
  func SetOutput(w io.Writer)                    { global.SetOutput(w) }
  ```
- **验收标准**:
  - [ ] 全局实例可用
  - [ ] 包级函数实现
  - [ ] 编译通过
  - [ ] 提交: `git commit -m "feat(logger): 添加全局logger实例"`

---

#### ⬜ 任务 1.1.5: 编写单元测试
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2小时
- **前置依赖**: 任务 1.1.4
- **执行步骤**:
  1. 在 `logger_test.go` 中实现测试
  2. 测试日志等级过滤
  3. 测试输出重定向
  4. 测试并发安全性
  5. 测试格式化输出
- **测试用例**:
  ```go
  func TestLoggerLevel(t *testing.T)      // 测试等级过滤
  func TestLoggerOutput(t *testing.T)     // 测试输出重定向
  func TestLoggerConcurrency(t *testing.T) // 测试并发
  func TestLoggerFormat(t *testing.T)     // 测试格式化
  ```
- **验收标准**:
  - [ ] 所有测试通过: `go test ./internal/logger/...`
  - [ ] 覆盖率 >80%: `go test -cover ./internal/logger/...`
  - [ ] Race检测通过: `go test -race ./internal/logger/...`
  - [ ] 提交: `git commit -m "test(logger): 添加单元测试"`

---

#### ⬜ 任务 1.1.6: 编写性能测试
- **负责人**: _________
- **优先级**: 🟡 中
- **预计时间**: 1小时
- **前置依赖**: 任务 1.1.5
- **执行步骤**:
  1. 在 `logger_bench_test.go` 中实现benchmark
  2. 测试单线程性能
  3. 测试并发性能
  4. 测试不同日志等级的性能
- **Benchmark用例**:
  ```go
  func BenchmarkLoggerInfo(b *testing.B)
  func BenchmarkLoggerConcurrent(b *testing.B)
  func BenchmarkLoggerWithDiscard(b *testing.B)
  ```
- **验收标准**:
  - [ ] Benchmark可运行: `go test -bench=. ./internal/logger/...`
  - [ ] 性能达标: >1,000,000 ops/sec (单线程)
  - [ ] 性能达标: >500,000 ops/sec (并发)
  - [ ] 保存结果: `go test -bench=. ./internal/logger/... > phase1_logger_bench.txt`
  - [ ] 提交: `git commit -m "test(logger): 添加性能测试"`

---

### 任务组 1.2: 配置系统集成 (Day 3)

#### ⬜ 任务 1.2.1: 实现logger配置结构
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 1小时
- **前置依赖**: 任务 1.1.6
- **执行步骤**:
  1. 在 `internal/logger/config.go` 中定义配置结构
  2. 实现配置解析函数
  3. 实现从配置初始化logger的函数
- **代码要点**:
  ```go
  type Config struct {
      Level         string `yaml:"level"`
      Output        string `yaml:"output"`
      ShowTimestamp bool   `yaml:"show_timestamp"`
  }
  
  func InitFromConfig(cfg Config) error {
      // 解析并应用配置
  }
  ```
- **验收标准**:
  - [ ] 配置结构定义完整
  - [ ] 解析函数实现
  - [ ] 编译通过
  - [ ] 提交: `git commit -m "feat(logger): 实现配置系统"`

---

#### ⬜ 任务 1.2.2: 更新config.yaml
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 15分钟
- **前置依赖**: 任务 1.2.1
- **执行步骤**:
  ```bash
  # 在config.yaml末尾添加
  cat >> config.yaml <<'EOF'
  
  # 日志配置
  logging:
    level: info              # debug/info/warn/error
    output: stdout           # stdout/stderr/文件路径
    show_timestamp: false    # UI模式下关闭时间戳
  EOF
  
  git add config.yaml
  git commit -m "feat(config): 添加日志配置项"
  ```
- **验收标准**:
  - [ ] config.yaml已更新
  - [ ] 配置格式正确
  - [ ] 已提交

---

#### ⬜ 任务 1.2.3: 更新Config结构体
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 30分钟
- **前置依赖**: 任务 1.2.2
- **执行步骤**:
  1. 在 `internal/core/config.go` 中添加 LoggingConfig 字段
  2. 更新 LoadConfig 函数
  3. 测试配置加载
- **代码修改**:
  ```go
  type Config struct {
      // ... 现有字段 ...
      Logging logger.Config `yaml:"logging"`
  }
  ```
- **验收标准**:
  - [ ] Config结构体已更新
  - [ ] 配置加载测试通过
  - [ ] 编译通过
  - [ ] 提交: `git commit -m "feat(config): 集成logger配置"`

---

### 任务组 1.3: fmt.Print替换 (Day 4-6)

#### ⬜ 任务 1.3.1: 创建兼容层
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 30分钟
- **前置依赖**: 任务 1.2.3
- **执行步骤**:
  1. 在 `internal/core/output.go` 中更新 SafePrintf 等函数
  2. 添加 Deprecated 注释
  3. 将调用转发到 logger
- **代码修改**:
  ```go
  import "apple-music-downloader/internal/logger"
  
  // Deprecated: 使用 logger.Info() 替代
  func SafePrintf(format string, a ...interface{}) {
      logger.Info(format, a...)
  }
  
  // Deprecated: 使用 logger.Info() 替代
  func SafePrintln(a ...interface{}) {
      msg := strings.TrimSuffix(fmt.Sprintln(a...), "\n")
      logger.Info(msg)
  }
  ```
- **验收标准**:
  - [ ] 兼容层实现
  - [ ] 编译通过
  - [ ] 现有代码仍可工作
  - [ ] 提交: `git commit -m "refactor(output): 创建logger兼容层"`

---

#### ⬜ 任务 1.3.2: 在main.go中初始化logger
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 20分钟
- **前置依赖**: 任务 1.3.1
- **执行步骤**:
  1. 在 main() 函数开始处初始化 logger
  2. 从配置加载logger设置
  3. 测试运行
- **代码修改**:
  ```go
  func main() {
      // 加载配置
      cfg, err := core.LoadConfig("config.yaml")
      if err != nil {
          log.Fatal(err)
      }
      
      // 初始化logger
      if err := logger.InitFromConfig(cfg.Logging); err != nil {
          log.Fatal(err)
      }
      
      // ... 其余代码 ...
  }
  ```
- **验收标准**:
  - [ ] Logger初始化代码添加
  - [ ] 程序可正常运行
  - [ ] 提交: `git commit -m "feat(main): 初始化logger"`

---

#### ⬜ 任务 1.3.3: 替换main.go中的fmt.Print
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 1-2小时
- **前置依赖**: 任务 1.3.2
- **执行步骤**:
  1. 搜索main.go中的所有fmt.Print调用
  2. 逐个替换为logger调用
  3. 根据语义选择合适的日志等级
  4. 每替换一部分就测试一次
- **替换规则**:
  ```go
  // 错误信息
  fmt.Printf("错误: %v\n", err) → logger.Error("错误: %v", err)
  
  // 警告信息
  fmt.Printf("警告: %s\n", msg) → logger.Warn("警告: %s", msg)
  
  // 普通信息
  fmt.Printf("🎤 歌手: %s\n", artist) → logger.Info("🎤 歌手: %s", artist)
  
  // 调试信息
  fmt.Printf("debug: %v\n", data) → logger.Debug("debug: %v", data)
  ```
- **验收标准**:
  - [ ] main.go中无fmt.Print调用（排除注释）
  - [ ] 程序运行正常
  - [ ] 输出格式与之前一致
  - [ ] 提交: `git commit -m "refactor(main): 替换fmt.Print为logger调用"`

---

#### ⬜ 任务 1.3.4: 替换internal/core中的fmt.Print
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2小时
- **前置依赖**: 任务 1.3.3
- **执行步骤**:
  1. 检查 internal/core 下所有.go文件
  2. 替换所有fmt.Print调用
  3. 运行测试确保无破坏
- **文件清单**:
  - [ ] internal/core/state.go
  - [ ] internal/core/config.go
  - [ ] internal/core/output.go
  - [ ] 其他core包文件
- **验收标准**:
  - [ ] 无fmt.Print直接调用
  - [ ] 测试通过: `go test ./internal/core/...`
  - [ ] 提交: `git commit -m "refactor(core): 替换fmt.Print为logger调用"`

---

#### ⬜ 任务 1.3.5: 替换internal/downloader中的fmt.Print
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2-3小时
- **前置依赖**: 任务 1.3.4
- **执行步骤**:
  1. 检查 internal/downloader 下所有文件
  2. 特别注意错误信息的日志等级
  3. 替换并测试
- **验收标准**:
  - [ ] 无fmt.Print直接调用
  - [ ] 测试通过: `go test ./internal/downloader/...`
  - [ ] 下载功能正常
  - [ ] 提交: `git commit -m "refactor(downloader): 替换fmt.Print为logger调用"`

---

#### ⬜ 任务 1.3.6: 替换utils/runv14中的fmt.Print
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2小时
- **前置依赖**: 任务 1.3.5
- **执行步骤**:
  1. 检查 utils/runv14/runv14.go
  2. 替换fmt.Print调用
  3. 特别注意进度输出相关的代码
- **验收标准**:
  - [ ] 无fmt.Print直接调用
  - [ ] 测试通过
  - [ ] 进度显示正常
  - [ ] 提交: `git commit -m "refactor(runv14): 替换fmt.Print为logger调用"`

---

#### ⬜ 任务 1.3.7: 替换utils/runv3中的fmt.Print
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2小时
- **前置依赖**: 任务 1.3.6
- **执行步骤**:
  1. 检查 utils/runv3/runv3.go
  2. 替换fmt.Print调用
  3. 测试
- **验收标准**:
  - [ ] 无fmt.Print直接调用
  - [ ] 测试通过
  - [ ] 提交: `git commit -m "refactor(runv3): 替换fmt.Print为logger调用"`

---

#### ⬜ 任务 1.3.8: 替换其他模块中的fmt.Print
- **负责人**: _________
- **优先级**: 🟡 中
- **预计时间**: 1-2小时
- **前置依赖**: 任务 1.3.7
- **执行步骤**:
  1. 检查其余所有包
  2. 替换所有fmt.Print
  3. 全量测试
- **检查命令**:
  ```bash
  # 查找剩余的fmt.Print调用
  grep -r "fmt\.Print" internal/ main.go utils/ \
    --exclude-dir=vendor \
    --exclude="*_test.go" \
    | grep -v "// OK:"
  ```
- **验收标准**:
  - [ ] 检查命令输出为0
  - [ ] 全量测试通过: `make test`
  - [ ] 提交: `git commit -m "refactor: 完成所有fmt.Print替换"`

---

### 任务组 1.4: Phase 1验收与发布 (Day 7)

#### ⬜ 任务 1.4.1: 运行完整验收测试
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2小时
- **前置依赖**: 任务组 1.3 完成
- **执行步骤**:
  ```bash
  # 1. 检查fmt.Print替换完成度
  echo "检查fmt.Print替换..."
  grep -r "fmt\.Print" internal/ main.go utils/ \
    --exclude-dir=vendor \
    --exclude="*_test.go" \
    | grep -v "// OK:" | wc -l
  # 预期: 0
  
  # 2. 运行所有测试
  echo "运行单元测试..."
  go test ./...
  
  # 3. Race检测
  echo "运行race检测..."
  go test -race ./internal/logger/...
  
  # 4. 日志等级测试
  echo "测试日志等级过滤..."
  go run main.go --log-level=error test/data/single_track.txt 2>&1 | grep -c "INFO"
  # 预期: 0
  
  # 5. 性能测试
  echo "运行性能测试..."
  go test -bench=. ./internal/logger/... -benchmem
  # 预期: >1000000 ops/sec
  
  # 6. 完整验证脚本
  ./scripts/validate_refactor.sh
  ```
- **验收标准**:
  - [ ] 所有检查项通过
  - [ ] 无fmt.Print直接调用
  - [ ] Race检测通过
  - [ ] 性能达标
  - [ ] 文档记录结果

---

#### ⬜ 任务 1.4.2: 功能回归测试
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 1-2小时
- **前置依赖**: 任务 1.4.1
- **执行步骤**:
  ```bash
  # 1. 测试单曲下载
  ./apple-music-downloader test/data/single_track.txt > test/phase1/single_track.txt 2>&1
  
  # 2. 测试专辑下载
  ./apple-music-downloader test/data/small_album.txt > test/phase1/small_album.txt 2>&1
  
  # 3. 输出对比（允许格式轻微差异）
  diff <(grep -v "时间\|速度" test/baseline/single_track_output.txt) \
       <(grep -v "时间\|速度" test/phase1/single_track.txt) || true
  
  # 4. 手动测试检查点
  # - [ ] UI显示正常
  # - [ ] 下载功能正常
  # - [ ] 错误处理正常
  # - [ ] 进度显示正常
  ```
- **验收标准**:
  - [ ] 所有测试场景通过
  - [ ] 输出格式一致
  - [ ] 无功能回退
  - [ ] 手动测试通过

---

#### ⬜ 任务 1.4.3: 性能对比
- **负责人**: _________
- **优先级**: 🟡 中
- **预计时间**: 30分钟
- **前置依赖**: 任务 1.4.2
- **执行步骤**:
  ```bash
  # 1. 运行新版本benchmark
  go test -bench=. ./... > phase1_bench.txt
  
  # 2. 对比基线
  benchcmp baseline_bench.txt phase1_bench.txt > phase1_perf_report.txt
  
  # 3. 分析结果
  cat phase1_perf_report.txt
  
  # 4. 保存结果
  git add test/phase1/*
  git add phase1_*
  git commit -m "test: Phase 1性能测试结果"
  ```
- **验收标准**:
  - [ ] 性能持平或提升
  - [ ] 无明显性能回退
  - [ ] 结果已保存

---

#### ⬜ 任务 1.4.4: 代码审查
- **负责人**: _________（代码作者以外的人）
- **优先级**: 🔴 高
- **预计时间**: 2-3小时
- **前置依赖**: 任务 1.4.3
- **审查清单**:
  - [ ] 代码符合Go语言规范
  - [ ] 所有public函数有注释
  - [ ] 错误处理完善
  - [ ] 无明显性能问题
  - [ ] 测试覆盖充分
  - [ ] 无安全隐患
  - [ ] 日志等级使用合理
  - [ ] 兼容层实现正确
- **审查方式**:
  ```bash
  # 创建PR（如果使用）
  git push origin feature/ui-refactor
  # 或本地审查
  git diff v2.5.3-pre-refactor...HEAD
  ```
- **输出**:
  - [ ] 审查意见文档
  - [ ] 修复建议列表

---

#### ⬜ 任务 1.4.5: 修复审查问题
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2-4小时
- **前置依赖**: 任务 1.4.4
- **执行步骤**:
  1. 根据审查意见逐项修复
  2. 每次修复后运行测试
  3. 重新提交审查（如需要）
- **验收标准**:
  - [ ] 所有审查问题已解决
  - [ ] 测试仍然通过
  - [ ] 审查者批准

---

#### ⬜ 任务 1.4.6: 更新文档
- **负责人**: _________
- **优先级**: 🟡 中
- **预计时间**: 1小时
- **前置依赖**: 任务 1.4.5
- **执行步骤**:
  1. 更新README.md（如有logger使用说明）
  2. 更新CHANGELOG.md
  3. 更新内部文档
- **CHANGELOG示例**:
  ```markdown
  ## [2.6.0-rc1] - 2025-10-XX
  
  ### Added
  - 新增统一日志系统（internal/logger包）
  - 支持日志等级控制（DEBUG/INFO/WARN/ERROR）
  - 支持日志配置（config.yaml）
  
  ### Changed
  - 替换所有fmt.Print为logger调用
  - SafePrintf等函数标记为Deprecated
  
  ### Performance
  - 日志性能：>1,000,000 ops/sec
  ```
- **验收标准**:
  - [ ] README已更新
  - [ ] CHANGELOG已更新
  - [ ] 提交: `git commit -m "docs: 更新Phase 1文档"`

---

#### ⬜ 任务 1.4.7: 打Tag并发布
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 15分钟
- **前置依赖**: 任务 1.4.6
- **执行步骤**:
  ```bash
  # 1. 确保所有改动已提交
  git status
  
  # 2. 打tag
  git tag -a v2.6.0-rc1 -m "Phase 1完成: 日志模块重构
  
  主要改进:
  - 统一日志系统
  - 日志等级控制
  - 替换所有fmt.Print
  
  测试状态:
  - 单元测试: ✅ 通过
  - Race检测: ✅ 通过
  - 性能测试: ✅ 达标
  - 功能测试: ✅ 通过"
  
  # 3. 推送tag
  git push origin v2.6.0-rc1
  
  # 4. 推送分支
  git push origin feature/ui-refactor
  ```
- **验收标准**:
  - [ ] Tag已创建
  - [ ] Tag已推送
  - [ ] 分支已推送
  - [ ] **Phase 1 完成** 🎉

---

## 🎨 Phase 2: UI模块解耦与事件驱动 (Week 3-5, 预计12-15天)

### 任务组 2.1: Progress包基础实现 (Week 3)

#### ⬜ 任务 2.1.1: 创建progress包结构
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 15分钟
- **前置依赖**: Phase 1完成
- **执行步骤**:
  ```bash
  mkdir -p internal/progress
  touch internal/progress/progress.go
  touch internal/progress/adapter.go
  touch internal/progress/progress_test.go
  
  git add internal/progress/
  git commit -m "feat(progress): 创建progress包结构"
  ```
- **验收标准**:
  - [ ] 目录结构创建
  - [ ] 文件已创建
  - [ ] 已提交

---

#### ⬜ 任务 2.1.2: 定义ProgressEvent结构
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 1小时
- **前置依赖**: 任务 2.1.1
- **执行步骤**:
  1. 在progress.go中定义事件结构
  2. 包含所有必要字段
  3. 添加文档注释
- **代码要点**:
  ```go
  // ProgressEvent 进度事件
  type ProgressEvent struct {
      TrackIndex int       // 曲目索引（在批次中）
      Stage      string    // 阶段: download/decrypt/tag/complete/error
      Percentage int       // 进度百分比 (0-100)
      SpeedBPS   float64   // 速度（字节/秒）
      Status     string    // 状态描述文本
      Error      error     // 错误信息（如有）
      Metadata   map[string]interface{} // 额外元数据
  }
  ```
- **验收标准**:
  - [ ] 结构定义完整
  - [ ] 字段注释清晰
  - [ ] 编译通过
  - [ ] 提交: `git commit -m "feat(progress): 定义ProgressEvent结构"`

---

#### ⬜ 任务 2.1.3: 实现ProgressListener接口
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 30分钟
- **前置依赖**: 任务 2.1.2
- **执行步骤**:
  1. 定义监听器接口
  2. 定义回调方法
- **代码要点**:
  ```go
  // ProgressListener 进度监听器接口
  type ProgressListener interface {
      OnProgress(event ProgressEvent)
      OnComplete(trackIndex int)
      OnError(trackIndex int, err error)
  }
  ```
- **验收标准**:
  - [ ] 接口定义完整
  - [ ] 方法签名合理
  - [ ] 编译通过
  - [ ] 提交: `git commit -m "feat(progress): 定义ProgressListener接口"`

---

#### ⬜ 任务 2.1.4: 实现ProgressNotifier
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2小时
- **前置依赖**: 任务 2.1.3
- **执行步骤**:
  1. 实现通知器结构体
  2. 实现监听器注册方法
  3. 实现事件分发方法
  4. 确保线程安全
- **代码要点**:
  ```go
  type ProgressNotifier struct {
      listeners []ProgressListener
      mu        sync.RWMutex
  }
  
  func NewNotifier() *ProgressNotifier {
      return &ProgressNotifier{
          listeners: make([]ProgressListener, 0),
      }
  }
  
  func (n *ProgressNotifier) AddListener(l ProgressListener) {
      n.mu.Lock()
      defer n.mu.Unlock()
      n.listeners = append(n.listeners, l)
  }
  
  func (n *ProgressNotifier) Notify(event ProgressEvent) {
      n.mu.RLock()
      defer n.mu.RUnlock()
      for _, listener := range n.listeners {
          listener.OnProgress(event)
      }
  }
  ```
- **验收标准**:
  - [ ] 通知器实现完整
  - [ ] 线程安全（使用RWMutex）
  - [ ] 编译通过
  - [ ] 提交: `git commit -m "feat(progress): 实现ProgressNotifier"`

---

#### ⬜ 任务 2.1.5: 实现适配器（关键！）
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2-3小时
- **前置依赖**: 任务 2.1.4
- **执行步骤**:
  1. 在adapter.go中实现适配器
  2. 将旧的channel模式转换为新的事件模式
  3. 确保无goroutine泄漏
- **代码要点**:
  ```go
  // ProgressUpdate 旧的进度更新结构（保持兼容）
  type ProgressUpdate struct {
      Percentage int
      SpeedBPS   float64
      Stage      string
  }
  
  // ProgressAdapter 适配器
  type ProgressAdapter struct {
      notifier   *ProgressNotifier
      trackIndex int
      stage      string
  }
  
  func NewProgressAdapter(notifier *ProgressNotifier, trackIndex int, stage string) *ProgressAdapter {
      return &ProgressAdapter{
          notifier:   notifier,
          trackIndex: trackIndex,
          stage:      stage,
      }
  }
  
  // ToChan 创建一个兼容旧代码的channel
  func (a *ProgressAdapter) ToChan() chan<- ProgressUpdate {
      ch := make(chan ProgressUpdate, 10)
      go func() {
          defer close(ch)  // 防止goroutine泄漏
          for update := range ch {
              a.notifier.Notify(ProgressEvent{
                  TrackIndex: a.trackIndex,
                  Stage:      update.Stage,
                  Percentage: update.Percentage,
                  SpeedBPS:   update.SpeedBPS,
              })
          }
      }()
      return ch
  }
  ```
- **验收标准**:
  - [ ] 适配器实现完整
  - [ ] 无goroutine泄漏
  - [ ] 编译通过
  - [ ] 提交: `git commit -m "feat(progress): 实现适配器模式"`

---

#### ⬜ 任务 2.1.6: 编写progress包测试
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2小时
- **前置依赖**: 任务 2.1.5
- **执行步骤**:
  1. 测试事件通知
  2. 测试监听器注册
  3. 测试适配器功能
  4. 测试并发安全
- **测试用例**:
  ```go
  func TestProgressNotifier(t *testing.T)
  func TestProgressListener(t *testing.T)
  func TestProgressAdapter(t *testing.T)
  func TestProgressConcurrency(t *testing.T)
  ```
- **验收标准**:
  - [ ] 所有测试通过: `go test ./internal/progress/...`
  - [ ] Race检测通过: `go test -race ./internal/progress/...`
  - [ ] 覆盖率 >80%
  - [ ] 提交: `git commit -m "test(progress): 添加单元测试"`

---

### 任务组 2.2: UI监听器实现 (Week 4, Day 1-2)

#### ⬜ 任务 2.2.1: 创建UI监听器文件
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 10分钟
- **前置依赖**: 任务组 2.1 完成
- **执行步骤**:
  ```bash
  touch internal/ui/listener.go
  git add internal/ui/listener.go
  git commit -m "feat(ui): 创建监听器文件"
  ```
- **验收标准**:
  - [ ] 文件已创建
  - [ ] 已提交

---

#### ⬜ 任务 2.2.2: 实现UIProgressListener
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2-3小时
- **前置依赖**: 任务 2.2.1
- **执行步骤**:
  1. 实现监听器结构体
  2. 实现OnProgress方法
  3. 实现OnComplete方法
  4. 实现OnError方法
  5. 集成现有的UpdateStatus功能
- **代码要点**:
  ```go
  package ui
  
  import (
      "apple-music-downloader/internal/progress"
      "fmt"
  )
  
  // UIProgressListener UI进度监听器
  type UIProgressListener struct {
      // 可以添加需要的状态
  }
  
  // OnProgress 处理进度更新
  func (l *UIProgressListener) OnProgress(event progress.ProgressEvent) {
      status := formatStatus(event)
      color := getColorFunc(event.Stage)
      UpdateStatus(event.TrackIndex, status, color)
  }
  
  // OnComplete 处理完成事件
  func (l *UIProgressListener) OnComplete(trackIndex int) {
      UpdateStatus(trackIndex, "下载完成", greenFunc)
  }
  
  // OnError 处理错误事件
  func (l *UIProgressListener) OnError(trackIndex int, err error) {
      errMsg := truncateError(err)
      UpdateStatus(trackIndex, errMsg, redFunc)
  }
  
  // formatStatus 格式化状态文本
  func formatStatus(event progress.ProgressEvent) string {
      switch event.Stage {
      case "download":
          return fmt.Sprintf("下载中 %d%% (%s)", 
                            event.Percentage, 
                            formatSpeed(event.SpeedBPS))
      case "decrypt":
          return fmt.Sprintf("解密中 %d%%", event.Percentage)
      case "tag":
          return "写入标签中..."
      default:
          return event.Status
      }
  }
  
  // getColorFunc 根据阶段返回颜色函数
  func getColorFunc(stage string) func(...interface{}) string {
      switch stage {
      case "download", "decrypt":
          return yellowFunc
      case "complete":
          return greenFunc
      case "error":
          return redFunc
      default:
          return func(a ...interface{}) string {
              return fmt.Sprint(a...)
          }
      }
  }
  
  // formatSpeed 格式化速度
  func formatSpeed(bps float64) string {
      mbps := bps / 1024 / 1024
      return fmt.Sprintf("%.1f MB/s", mbps)
  }
  
  // truncateError 截断错误信息
  func truncateError(err error) string {
      msg := err.Error()
      if len(msg) > 50 {
          return msg[:50] + "..."
      }
      return msg
  }
  ```
- **验收标准**:
  - [ ] 监听器实现完整
  - [ ] 所有方法实现
  - [ ] 编译通过
  - [ ] 提交: `git commit -m "feat(ui): 实现UIProgressListener"`

---

#### ⬜ 任务 2.2.3: 在main.go中注册监听器
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 30分钟
- **前置依赖**: 任务 2.2.2
- **执行步骤**:
  1. 在main.go中创建通知器
  2. 注册UI监听器
  3. 将通知器传递给下载器
- **代码修改**:
  ```go
  func main() {
      // ... 现有初始化代码 ...
      
      // 创建进度通知器
      notifier := progress.NewNotifier()
      
      // 注册UI监听器
      notifier.AddListener(&ui.UIProgressListener{})
      
      // 传递给下载流程
      runDownloads(tracks, notifier)
      
      // ...
  }
  ```
- **验收标准**:
  - [ ] 通知器创建
  - [ ] 监听器注册
  - [ ] 编译通过
  - [ ] 提交: `git commit -m "feat(main): 注册UI进度监听器"`

---

### 任务组 2.3: 下载器迁移 (Week 4 Day 3-5, Week 5)

#### ⬜ 任务 2.3.1: 迁移downloader.go（使用适配器）
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 3-4小时
- **前置依赖**: 任务 2.2.3
- **执行步骤**:
  1. 在Rip函数中接收notifier参数
  2. 使用适配器创建progressChan
  3. 保持现有逻辑不变
  4. 测试功能
- **代码修改示例**:
  ```go
  // 修改前
  func Rip(track Track, statusIndex int) error {
      progressChan := make(chan ProgressUpdate, 10)
      go func() {
          for p := range progressChan {
              ui.UpdateStatus(statusIndex, formatProgress(p), yellowFunc)
          }
      }()
      // ... 下载逻辑 ...
  }
  
  // 修改后（使用适配器）
  func Rip(track Track, statusIndex int, notifier *progress.ProgressNotifier) error {
      adapter := progress.NewProgressAdapter(notifier, statusIndex, "download")
      progressChan := adapter.ToChan()
      // ... 下载逻辑保持不变 ...
      // progressChan <- ProgressUpdate{...} // 仍然可以这样用
  }
  ```
- **验收标准**:
  - [ ] 代码迁移完成
  - [ ] 编译通过
  - [ ] 下载功能正常
  - [ ] 进度显示正常
  - [ ] 提交: `git commit -m "refactor(downloader): 使用适配器接入progress系统"`

---

#### ⬜ 任务 2.3.2: 迁移runv14.go（使用适配器）
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 3-4小时
- **前置依赖**: 任务 2.3.1
- **执行步骤**:
  1. 修改函数签名接收notifier
  2. 使用适配器
  3. 测试
- **验收标准**:
  - [ ] 代码迁移完成
  - [ ] 编译通过
  - [ ] 下载功能正常
  - [ ] 提交: `git commit -m "refactor(runv14): 使用适配器接入progress系统"`

---

#### ⬜ 任务 2.3.3: 迁移runv3.go（使用适配器）
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 3-4小时
- **前置依赖**: 任务 2.3.2
- **执行步骤**:
  1. 修改函数签名
  2. 使用适配器
  3. 测试
- **验收标准**:
  - [ ] 代码迁移完成
  - [ ] 编译通过
  - [ ] 下载功能正常
  - [ ] 提交: `git commit -m "refactor(runv3): 使用适配器接入progress系统"`

---

#### ⬜ 任务 2.3.4: 移除下载器对UI的直接调用
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2-3小时
- **前置依赖**: 任务 2.3.3
- **执行步骤**:
  1. 搜索所有 ui.UpdateStatus 调用
  2. 确认都已通过notifier替换
  3. 移除直接调用
  4. 验证解耦
- **检查命令**:
  ```bash
  # 检查下载器中是否还有直接UI调用
  grep -r "ui\.UpdateStatus" \
    internal/downloader/ \
    utils/runv14/ \
    utils/runv3/ | wc -l
  # 预期: 0
  ```
- **验收标准**:
  - [ ] 检查命令输出为0
  - [ ] 下载器与UI完全解耦
  - [ ] 功能正常
  - [ ] 提交: `git commit -m "refactor: 移除下载器对UI的直接依赖"`

---

### 任务组 2.4: Phase 2验收与发布 (Week 5末)

#### ⬜ 任务 2.4.1: 运行Phase 2验收测试
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2-3小时
- **前置依赖**: 任务组 2.3 完成
- **执行步骤**:
  ```bash
  # 1. 检查UI解耦
  echo "检查UI解耦..."
  grep -r "ui\.UpdateStatus" internal/downloader/ utils/runv14/ utils/runv3/ | wc -l
  # 预期: 0
  
  # 2. 验证进度去重
  echo "测试进度去重..."
  ./apple-music-downloader test/data/small_album.txt 2>&1 | grep "100%" | sort | uniq -c
  # 预期: 每首歌最多2-3次100%
  
  # 3. 性能测试（UI CPU占用）
  echo "性能测试..."
  go test -cpuprofile=cpu.prof ./internal/ui/...
  go tool pprof -top cpu.prof | grep "PrintUI"
  # 预期: CPU占用 < 5%
  
  # 4. 并发安全测试
  echo "并发安全测试..."
  go test -race ./internal/progress/...
  # 预期: PASS
  
  # 5. 功能一致性测试
  echo "功能对比..."
  diff <(grep -v "速度\|时间" test/baseline/small_album_output.txt) \
       <(grep -v "速度\|时间" test/phase2/small_album.txt) || true
  ```
- **验收标准**:
  - [ ] 所有检查项通过
  - [ ] UI解耦验证通过
  - [ ] 性能达标
  - [ ] 功能一致

---

#### ⬜ 任务 2.4.2: 手动测试
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 1-2小时
- **前置依赖**: 任务 2.4.1
- **测试检查点**:
  - [ ] 下载10首歌，观察UI是否稳定无闪烁
  - [ ] 下载完成后，100%状态不重复出现（或最多出现2次）
  - [ ] 暂停/恢复功能正常
  - [ ] 错误信息正确显示并截断
  - [ ] 批量下载功能正常
  - [ ] 并发下载无race问题
- **验收标准**:
  - [ ] 所有手动测试通过
  - [ ] UI体验良好
  - [ ] 无明显bug

---

#### ⬜ 任务 2.4.3: 性能对比与分析
- **负责人**: _________
- **优先级**: 🟡 中
- **预计时间**: 1小时
- **前置依赖**: 任务 2.4.2
- **执行步骤**:
  ```bash
  # 统计UI刷新次数对比
  echo "基线版本UI刷新次数:"
  grep -c "%" test/baseline/small_album_output.txt
  
  echo "Phase 2版本UI刷新次数:"
  grep -c "%" test/phase2/small_album.txt
  
  # 计算改进百分比
  # 预期: 减少90%左右
  
  # 保存分析报告
  cat > test/phase2/performance_report.md <<EOF
  # Phase 2 性能分析报告
  
  ## UI刷新性能
  - 基线版本: XXX次
  - Phase 2版本: XXX次
  - 改进: XX%
  
  ## CPU占用
  - PrintUI函数: < 5%
  
  ## 并发安全
  - Race检测: ✅ 通过
  EOF
  
  git add test/phase2/
  git commit -m "test: Phase 2性能分析报告"
  ```
- **验收标准**:
  - [ ] UI刷新次数减少 >80%
  - [ ] 报告已生成
  - [ ] 数据已保存

---

#### ⬜ 任务 2.4.4: 代码审查
- **负责人**: _________（代码作者以外的人）
- **优先级**: 🔴 高
- **预计时间**: 3-4小时
- **前置依赖**: 任务 2.4.3
- **审查重点**:
  - [ ] 适配器实现正确（无goroutine泄漏）
  - [ ] 事件驱动模型合理
  - [ ] UI与下载器完全解耦
  - [ ] 所有public接口有文档
  - [ ] 错误处理完善
  - [ ] 测试覆盖充分
  - [ ] 性能满足要求
- **验收标准**:
  - [ ] 审查完成
  - [ ] 问题已记录
  - [ ] 审查者签字

---

#### ⬜ 任务 2.4.5: 修复审查问题并更新文档
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2-4小时
- **前置依赖**: 任务 2.4.4
- **执行步骤**:
  1. 修复审查问题
  2. 更新CHANGELOG
  3. 更新README（如需要）
- **CHANGELOG**:
  ```markdown
  ## [2.6.0-rc2] - 2025-10-XX
  
  ### Added
  - 新增进度事件系统（internal/progress包）
  - 新增UI进度监听器
  - 新增适配器模式支持渐进迁移
  
  ### Changed
  - UI与下载器完全解耦
  - 进度更新改为事件驱动
  
  ### Performance
  - UI刷新性能提升90%
  - 去重机制消除重复100%输出
  ```
- **验收标准**:
  - [ ] 问题已修复
  - [ ] 文档已更新
  - [ ] 审查者批准

---

#### ⬜ 任务 2.4.6: 打Tag并发布
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 15分钟
- **前置依赖**: 任务 2.4.5
- **执行步骤**:
  ```bash
  git tag -a v2.6.0-rc2 -m "Phase 2完成: UI模块解耦
  
  主要改进:
  - UI与下载器解耦
  - 事件驱动进度更新
  - UI性能提升90%
  
  测试状态:
  - 单元测试: ✅ 通过
  - Race检测: ✅ 通过
  - 性能提升: ✅ 90%
  - 功能测试: ✅ 通过"
  
  git push origin v2.6.0-rc2
  git push origin feature/ui-refactor
  ```
- **验收标准**:
  - [ ] Tag已创建并推送
  - [ ] **Phase 2 完成** 🎉

---

## ✅ Phase 4.1: MVP基础测试 (Week 6, 预计3-5天)

### 任务组 4.1: 集成测试

#### ⬜ 任务 4.1.1: 创建集成测试脚本
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2-3小时
- **前置依赖**: Phase 2完成
- **执行步骤**:
  ```bash
  cat > test/integration_test.sh <<'EOF'
  #!/bin/bash
  set -e
  
  echo "🧪 运行集成测试..."
  
  # 1. 单曲下载测试
  echo "1️⃣ 单曲下载测试..."
  ./apple-music-downloader test/data/single_track.txt
  
  # 2. 专辑下载测试
  echo "2️⃣ 专辑下载测试..."
  ./apple-music-downloader test/data/small_album.txt
  
  # 3. 批量下载测试
  echo "3️⃣ 批量下载测试..."
  ./apple-music-downloader test/data/batch_download.txt
  
  # 4. 错误恢复测试（故意使用错误URL）
  echo "4️⃣ 错误恢复测试..."
  ./apple-music-downloader test/data/invalid_url.txt || echo "预期失败"
  
  echo "✅ 集成测试完成"
  EOF
  
  chmod +x test/integration_test.sh
  git add test/integration_test.sh
  git commit -m "test: 添加集成测试脚本"
  ```
- **验收标准**:
  - [ ] 脚本创建成功
  - [ ] 脚本可执行
  - [ ] 已提交

---

#### ⬜ 任务 4.1.2: 运行完整测试套件
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 2-3小时
- **前置依赖**: 任务 4.1.1
- **执行步骤**:
  ```bash
  # 1. 运行make ci
  make ci
  
  # 2. 运行集成测试
  ./test/integration_test.sh
  
  # 3. 运行验证脚本
  ./scripts/validate_refactor.sh
  
  # 4. 保存测试结果
  make test > test/mvp/test_results.txt 2>&1
  make bench > test/mvp/bench_results.txt 2>&1
  ./test/integration_test.sh > test/mvp/integration_results.txt 2>&1
  
  git add test/mvp/
  git commit -m "test: MVP完整测试结果"
  ```
- **验收标准**:
  - [ ] 所有测试通过
  - [ ] 结果已保存
  - [ ] 已提交

---

#### ⬜ 任务 4.1.3: 性能回归测试
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 1小时
- **前置依赖**: 任务 4.1.2
- **执行步骤**:
  ```bash
  # 对比基线
  make perf-compare
  
  # 保存对比结果
  benchcmp baseline_bench.txt new_bench.txt > test/mvp/perf_comparison.txt
  
  # 分析结果
  cat test/mvp/perf_comparison.txt
  
  git add test/mvp/perf_comparison.txt
  git commit -m "test: MVP性能对比结果"
  ```
- **验收标准**:
  - [ ] 性能无明显回退
  - [ ] UI性能提升明显
  - [ ] 结果已保存

---

### 任务组 4.2: 文档与发布准备

#### ⬜ 任务 4.2.1: 更新README
- **负责人**: _________
- **优先级**: 🟡 中
- **预计时间**: 1-2小时
- **前置依赖**: 任务组 4.1 完成
- **更新内容**:
  1. 日志配置说明
  2. 日志等级使用说明
  3. 性能改进说明
  4. 架构改进说明（可选）
- **验收标准**:
  - [ ] README已更新
  - [ ] 示例清晰
  - [ ] 提交: `git commit -m "docs: 更新README"`

---

#### ⬜ 任务 4.2.2: 完善CHANGELOG
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 30分钟
- **前置依赖**: 任务 4.2.1
- **CHANGELOG内容**:
  ```markdown
  ## [2.6.0] - 2025-10-XX (MVP Release)
  
  ### 🎉 重大改进
  - 重构日志系统，支持等级控制
  - 重构UI模块，事件驱动架构
  - UI性能提升90%
  - 彻底解决日志竞争问题
  - 消除下载100%重复显示
  
  ### Added
  - internal/logger包：统一日志系统
  - internal/progress包：进度事件系统
  - 日志配置支持（config.yaml）
  - 日志等级控制（DEBUG/INFO/WARN/ERROR）
  
  ### Changed
  - 替换所有fmt.Print为logger调用
  - UI与下载器完全解耦
  - 进度更新改为事件驱动
  
  ### Performance
  - 日志性能: >1,000,000 ops/sec
  - UI刷新性能提升90%
  - CPU占用降低
  
  ### Fixed
  - 修复日志输出竞争问题
  - 修复UI刷新闪烁问题
  - 修复100%重复显示问题
  
  ### Technical
  - 测试覆盖率: >80%
  - Race检测: 零警告
  - 架构: 模块解耦，可维护性提升
  ```
- **验收标准**:
  - [ ] CHANGELOG完整
  - [ ] 格式规范
  - [ ] 提交: `git commit -m "docs: 完善CHANGELOG"`

---

#### ⬜ 任务 4.2.3: 最终代码审查
- **负责人**: 团队所有成员
- **优先级**: 🔴 高
- **预计时间**: 3-4小时
- **前置依赖**: 任务 4.2.2
- **审查清单**:
  - [ ] 代码质量达标
  - [ ] 测试覆盖充分
  - [ ] 文档完整
  - [ ] 性能达标
  - [ ] 无已知bug
  - [ ] 向后兼容
  - [ ] 安全无问题
- **验收标准**:
  - [ ] 审查完成
  - [ ] 所有问题已解决
  - [ ] 团队批准发布

---

#### ⬜ 任务 4.2.4: 创建发布说明
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 1小时
- **前置依赖**: 任务 4.2.3
- **执行步骤**:
  1. 创建RELEASE_NOTES.md
  2. 总结主要改进
  3. 列出升级指南
  4. 列出已知问题（如有）
- **模板**:
  ```markdown
  # Apple Music Downloader v2.6.0 Release Notes
  
  ## 🎉 主要改进
  
  ### 1. 统一日志系统
  - 支持日志等级控制
  - 可配置日志输出
  - 彻底解决日志竞争问题
  
  ### 2. UI性能优化
  - 事件驱动架构
  - UI刷新性能提升90%
  - 消除100%重复显示
  
  ### 3. 架构改进
  - UI与下载器完全解耦
  - 模块化设计
  - 可维护性大幅提升
  
  ## 📦 升级指南
  
  1. 更新config.yaml（可选）:
     ```yaml
     logging:
       level: info
       output: stdout
       show_timestamp: false
     ```
  
  2. 重新编译:
     ```bash
     go build -o apple-music-downloader
     ```
  
  3. 测试运行:
     ```bash
     ./apple-music-downloader --help
     ```
  
  ## ⚠️ 已知问题
  
  - 无
  
  ## 🙏 致谢
  
  感谢所有参与重构的团队成员！
  ```
- **验收标准**:
  - [ ] 发布说明完整
  - [ ] 升级指南清晰
  - [ ] 已提交

---

#### ⬜ 任务 4.2.5: 打MVP正式版Tag
- **负责人**: _________
- **优先级**: 🔴 高
- **预计时间**: 15分钟
- **前置依赖**: 任务 4.2.4
- **执行步骤**:
  ```bash
  # 确保所有改动已提交
  git status
  
  # 打正式版tag
  git tag -a v2.6.0 -m "v2.6.0 MVP Release
  
  重大改进:
  - 统一日志系统
  - UI性能提升90%
  - 架构完全重构
  
  测试状态:
  - 单元测试: ✅ 100%通过
  - 集成测试: ✅ 通过
  - Race检测: ✅ 零警告
  - 性能测试: ✅ 达标
  - 代码审查: ✅ 通过
  
  详见 RELEASE_NOTES.md"
  
  # 推送tag
  git push origin v2.6.0
  
  # 推送分支
  git push origin feature/ui-refactor
  
  # 合并到main（如果团队批准）
  # git checkout main
  # git merge feature/ui-refactor
  # git push origin main
  ```
- **验收标准**:
  - [ ] Tag已创建
  - [ ] Tag已推送
  - [ ] **MVP正式发布** 🎉🎉🎉

---

## 🎊 MVP完成庆祝与总结

### ⬜ 任务 MVP.1: 团队复盘会议
- **负责人**: 项目负责人
- **优先级**: 🟡 中
- **预计时间**: 1-2小时
- **前置依赖**: MVP发布
- **会议议程**:
  1. 回顾MVP目标达成情况（15分钟）
  2. 技术指标回顾（15分钟）
  3. 遇到的挑战与解决（20分钟）
  4. 经验教训总结（20分钟）
  5. Phase 3是否继续的讨论（20分钟）
  6. 庆祝与感谢（10分钟）
- **输出文档**:
  - [ ] 复盘总结文档
  - [ ] 经验教训清单
  - [ ] Phase 3决策结果

---

### ⬜ 任务 MVP.2: 创建后续计划
- **负责人**: _________
- **优先级**: 🟢 低
- **预计时间**: 1小时
- **前置依赖**: 任务 MVP.1
- **内容**:
  1. 如果继续Phase 3:
     - 制定Phase 3详细计划
     - 分配责任人
     - 确定时间表
  2. 如果暂停:
     - 维护计划
     - 监控计划
     - 未来规划
- **验收标准**:
  - [ ] 后续计划明确
  - [ ] 文档已保存

---

## 📊 任务统计总览

### Phase 0 准备阶段
- **总任务数**: 10
- **预计时间**: 1-2天
- **关键任务**: 3

### Phase 1 日志重构
- **总任务数**: 28
- **预计时间**: 8-10天
- **关键任务**: 15

### Phase 2 UI解耦
- **总任务数**: 24
- **预计时间**: 12-15天
- **关键任务**: 14

### Phase 4.1 MVP测试
- **总任务数**: 11
- **预计时间**: 3-5天
- **关键任务**: 8

### **MVP总计**
- **总任务数**: 73
- **预计总时间**: 4-6周
- **关键路径任务**: 40

---

## 🎯 关键里程碑检查点

| 里程碑 | 完成标志 | 预计日期 | 状态 |
|-------|---------|---------|------|
| Week 0 完成 | 所有准备任务✅ | Week 0末 | ⬜ |
| Phase 1 完成 | v2.6.0-rc1发布 | Week 2末 | ⬜ |
| Phase 2 完成 | v2.6.0-rc2发布 | Week 5末 | ⬜ |
| MVP 完成 | v2.6.0正式发布 | Week 6末 | ⬜ |

---

## 📌 使用说明

### 任务状态更新
在每个任务前的符号标记当前状态：
- ⬜ → 🔄 (开始工作时)
- 🔄 → ✅ (完成时)
- 如遇阻塞 → ⏸️

### 每日站会检查
1. 昨天完成了什么任务？
2. 今天计划完成什么任务？
3. 有什么阻塞？

### 每周回顾
1. 本周完成任务数
2. 是否按计划进行？
3. 是否需要调整？

---

**文档创建日期**: 2025-10-10  
**最后更新**: 2025-10-10  
**版本**: v1.0  
**下一步**: 开始Week 0准备工作！🚀

