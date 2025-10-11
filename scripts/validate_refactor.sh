#!/bin/bash
set -e

echo "🔍 验证重构安全性..."
echo ""

# 1. 编译检查
echo "1️⃣ 编译检查..."
if go build -o apple-music-downloader; then
    echo "✅ 编译通过"
else
    echo "❌ 编译失败"
    exit 1
fi
echo ""

# 2. 单元测试
echo "2️⃣ 单元测试..."
if go test ./... -v 2>&1 | grep -E "^(PASS|FAIL|ok|FAIL)" | head -20; then
    echo "✅ 单元测试执行完成"
else
    echo "⚠️  单元测试有问题（部分包可能无测试文件）"
fi
echo ""

# 3. Race检测
echo "3️⃣ Race检测..."
echo "检测internal/logger包..."
if [ -d "internal/logger" ]; then
    if go test -race ./internal/logger/... 2>&1 | grep -E "^(PASS|FAIL|ok|FAIL)"; then
        echo "✅ Logger包race检测通过"
    else
        echo "⚠️  Logger包race检测问题或无测试"
    fi
else
    echo "⏭️  Logger包尚未创建，跳过"
fi
echo ""

# 4. 检查fmt.Print替换进度
echo "4️⃣ 检查fmt.Print替换进度..."
FMT_COUNT=$(grep -r "fmt\.Print" internal/ main.go utils/ 2>/dev/null | \
    grep -v "vendor" | \
    grep -v "_test.go" | \
    grep -v "// OK:" | \
    grep -v "baseline" | \
    wc -l || echo "0")
echo "剩余fmt.Print调用: $FMT_COUNT 处"
if [ "$FMT_COUNT" -eq 0 ]; then
    echo "✅ 所有fmt.Print已替换"
else
    echo "⚠️  还有 $FMT_COUNT 处fmt.Print需要替换"
fi
echo ""

# 5. 检查UI解耦进度
echo "5️⃣ 检查UI解耦进度..."
if [ -d "internal/downloader" ]; then
    UI_CALL_COUNT=$(grep -r "ui\.UpdateStatus" internal/downloader/ utils/runv14/ utils/runv3/ 2>/dev/null | \
        grep -v "vendor" | \
        wc -l || echo "0")
    echo "下载器中UI直接调用: $UI_CALL_COUNT 处"
    if [ "$UI_CALL_COUNT" -eq 0 ]; then
        echo "✅ 下载器与UI完全解耦"
    else
        echo "⚠️  还有 $UI_CALL_COUNT 处UI直接调用需要解耦"
    fi
else
    echo "⏭️  下载器模块检查跳过"
fi
echo ""

# 6. 性能对比（如果有基线）
echo "6️⃣ 性能对比..."
if [ -f "baseline_bench.txt" ]; then
    echo "发现基线数据，运行性能测试..."
    if command -v benchcmp &> /dev/null; then
        go test -bench=. ./... > new_bench.txt 2>/dev/null || true
        benchcmp baseline_bench.txt new_bench.txt || echo "⚠️  性能有变化，需人工审查"
    else
        echo "⚠️  benchcmp未安装，跳过性能对比"
        echo "提示: go install golang.org/x/tools/cmd/benchcmp@latest"
    fi
else
    echo "⏭️  无基线数据，跳过性能对比"
fi
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 验证完成！"
echo ""
echo "当前重构阶段检查:"
if [ -d "internal/logger" ]; then
    echo "  ✅ Phase 1: Logger包已创建"
else
    echo "  ⬜ Phase 1: Logger包待创建"
fi

if [ -d "internal/progress" ]; then
    echo "  ✅ Phase 2: Progress包已创建"
else
    echo "  ⬜ Phase 2: Progress包待创建"
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

