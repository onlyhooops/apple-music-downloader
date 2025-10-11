.PHONY: all build test bench race lint clean validate ci help

# 默认目标
all: build

# 构建
build:
	@echo "🔨 构建项目..."
	go build -o apple-music-downloader
	@echo "✅ 构建完成"

# 运行测试
test:
	@echo "🧪 运行单元测试..."
	go test ./... -v -cover

# 性能测试
bench:
	@echo "⚡ 运行性能测试..."
	go test -bench=. ./... -benchmem

# 竞态检测
race:
	@echo "🔍 运行race检测..."
	go test -race ./...

# 代码检查
lint:
	@echo "📝 运行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint未安装"; \
		echo "安装命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# 清理
clean:
	@echo "🧹 清理构建文件..."
	rm -f apple-music-downloader
	rm -f apple-music-downloader-baseline
	rm -f *.prof
	rm -f new_bench.txt
	@echo "✅ 清理完成"

# 验证重构
validate:
	@echo "✅ 运行重构验证..."
	./scripts/validate_refactor.sh

# CI流程
ci: test race
	@echo "✅ 所有CI检查通过！"

# 性能对比
perf-compare:
	@echo "📊 性能对比..."
	@if [ ! -f baseline_bench.txt ]; then \
		echo "❌ baseline_bench.txt不存在"; \
		exit 1; \
	fi
	go test -bench=. ./... > new_bench.txt
	@if command -v benchcmp >/dev/null 2>&1; then \
		benchcmp baseline_bench.txt new_bench.txt; \
	else \
		echo "⚠️  benchcmp未安装"; \
		echo "安装命令: go install golang.org/x/tools/cmd/benchcmp@latest"; \
	fi

# 帮助信息
help:
	@echo "Apple Music Downloader - Makefile帮助"
	@echo ""
	@echo "可用命令:"
	@echo "  make build        - 构建项目"
	@echo "  make test         - 运行单元测试"
	@echo "  make bench        - 运行性能测试"
	@echo "  make race         - 运行竞态检测"
	@echo "  make lint         - 代码检查"
	@echo "  make clean        - 清理构建文件"
	@echo "  make validate     - 验证重构进度"
	@echo "  make ci           - 运行CI流程"
	@echo "  make perf-compare - 性能对比"
	@echo "  make help         - 显示此帮助信息"
	@echo ""

