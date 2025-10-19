# 📚 文档索引

> 本索引帮助您快速找到所需文档。项目已采用分类目录结构，便于维护和查找。

---

## 🚀 快速开始

| 文档 | 说明 |
|------|------|
| [README](./README.md) | 项目介绍、功能特性、快速开始指南（英文） |
| [README-CN](./README-CN.md) | 项目介绍、功能特性、快速开始指南（中文） |
| [BUILD_GUIDE](./BUILD_GUIDE.md) | 编译构建指南 |
| [CHANGELOG](./CHANGELOG.md) | 版本变更历史 |

---

## 📖 用户文档

### 功能特性 (`docs/features/`)

| 文档 | 说明 |
|------|------|
| [功能概览](./docs/features/overview.md) | 所有功能特性总览 |
| [历史记录系统](./docs/features/history-system.md) | 全局历史记录功能详解 |
| [工作-休息循环](./docs/features/work-rest-cycle.md) | 批量下载工作休息机制 |
| [缓存机制](./docs/features/cache-mechanism.md) | NFS优化缓存系统 |
| [MV质量显示](./docs/features/mv-quality-display.md) | MV视频质量信息显示 |

### 使用指南 (`docs/guides/`)

#### 快速开始 (`docs/guides/quick-start/`)
| 文档 | 说明 |
|------|------|
| [全局历史记录快速开始](./docs/guides/quick-start/global-history.md) | 5分钟上手全局历史记录 |
| [工作-休息循环快速开始](./docs/guides/quick-start/work-rest.md) | 快速启用工作休息功能 |

#### 高级指南 (`docs/guides/advanced/`)
| 文档 | 说明 |
|------|------|
| [全局历史记录详细指南](./docs/guides/advanced/global-history-guide.md) | 深入了解全局历史记录系统 |

### 问题排查 (`docs/troubleshooting/`)

| 文档 | 说明 |
|------|------|
| [历史记录保存问题](./docs/troubleshooting/history-save-issue.md) | 历史记录保存机制和常见问题 |

---

## 📦 版本发布 (`releases/`)

### v2.7.0 (Latest)
| 文档 | 说明 |
|------|------|
| [发布说明](./releases/v2.7.0/RELEASE_NOTES.md) | v2.7.0完整发布文档 |
| [归档文档](./releases/v2.7.0/archive/) | v2.7.0开发过程文档（已归档） |

---

## 🔧 技术文档 (`technical/`)

### 问题修复 (`technical/fixes/`)

| 文档 | 说明 |
|------|------|
| [ilst box修复](./technical/fixes/ilst-box-fix.md) | ilst box缺失自动修复 |
| [并发Map修复](./technical/fixes/concurrent-map-fix.md) | 并发写入Map崩溃修复 |
| [标签错误处理](./technical/fixes/tag-error-handling.md) | 标签写入失败优化 |

### 技术报告 (`technical/reports/`)

| 文档 | 说明 |
|------|------|
| [安全清理报告](./technical/reports/security-cleanup.md) | Git历史清理完成报告 |
| [功能验证报告](./technical/reports/feature-verification.md) | 工作-休息功能验证 |

---

## 📋 归档文档 (`archive/`)

历史文档和一次性报告的归档目录。

| 类别 | 路径 | 说明 |
|------|------|------|
| 开发总结 | `archive/development-summaries/` | 功能开发完成报告 |
| 迁移报告 | `archive/migration-reports/` | 系统迁移完成报告 |
| 一次性报告 | `archive/one-time-reports/` | 一次性技术报告 |
| 测试报告 | `archive/testing-reports/` | 功能测试验证报告 |
| 演示文档 | `archive/demos/` | 功能演示和示例 |

---

## 📊 文档分类说明

### 目录结构

```
项目根目录/
├── README.md, CHANGELOG.md          # 核心文档
├── BUILD_GUIDE.md                   # 构建指南
├── DOCUMENTATION_INDEX.md           # 本索引文件
│
├── docs/                            # 📁 用户文档
│   ├── features/                    # 功能特性文档
│   ├── guides/                      # 使用指南
│   │   ├── quick-start/             # 快速开始
│   │   └── advanced/                # 高级指南
│   └── troubleshooting/             # 问题排查
│
├── releases/                        # 📁 版本发布
│   └── v2.7.0/                      # 按版本分类
│       ├── RELEASE_NOTES.md         # 发布说明
│       └── archive/                 # 开发过程文档
│
├── technical/                       # 📁 技术文档
│   ├── fixes/                       # 问题修复
│   └── reports/                     # 技术报告
│
└── archive/                         # 📁 归档文档
    ├── development-summaries/       # 开发完成报告
    ├── migration-reports/           # 迁移完成报告
    ├── one-time-reports/            # 一次性报告
    ├── testing-reports/             # 测试报告
    └── demos/                       # 演示文档
```

### 文档类型说明

| 类型 | 位置 | 目标受众 | 更新频率 |
|------|------|----------|---------|
| 核心文档 | 根目录 | 所有用户 | 每个版本 |
| 用户文档 | `docs/` | 最终用户 | 功能变更时 |
| 版本文档 | `releases/` | 所有用户 | 每次发布 |
| 技术文档 | `technical/` | 开发者 | 问题修复时 |
| 归档文档 | `archive/` | 参考用 | 不再更新 |

---

## 🔍 如何查找文档？

### 按需求查找

| 我想... | 查看文档 |
|---------|---------|
| 了解项目和快速开始 | [README](./README.md) |
| 查看某个功能怎么用 | [docs/features/](./docs/features/) |
| 5分钟上手某个功能 | [docs/guides/quick-start/](./docs/guides/quick-start/) |
| 深入学习某个功能 | [docs/guides/advanced/](./docs/guides/advanced/) |
| 解决遇到的问题 | [docs/troubleshooting/](./docs/troubleshooting/) |
| 查看版本更新内容 | [CHANGELOG](./CHANGELOG.md) 或 [releases/](./releases/) |
| 了解技术实现细节 | [technical/](./technical/) |

### 按版本查找

- **最新版本**：查看 [README](./README.md) 或 [CHANGELOG](./CHANGELOG.md)
- **特定版本**：进入 `releases/v[版本号]/` 查看该版本的发布说明
- **版本历史**：查看 [CHANGELOG](./CHANGELOG.md)

---

## 💡 文档维护

### 贡献文档

如果您想贡献或改进文档：
1. 遵循现有的目录结构
2. 使用清晰的文档命名（小写+连字符）
3. 更新本索引文件
4. 提交Pull Request

### 文档规范

参见 [DOCUMENTATION_GOVERNANCE_REPORT](./DOCUMENTATION_GOVERNANCE_REPORT.md) 了解：
- 文档分类标准
- 命名规范
- 维护流程
- 质量要求

---

## 📞 反馈与帮助

- **文档问题**：提交GitHub Issue并标注 `documentation` 标签
- **功能问题**：查看对应功能文档或提交Issue
- **建议改进**：欢迎提交Pull Request

---

## 📋 文档治理

### 核心原则

本文档体系遵循以下治理原则：

#### 🎯 分类标准
- **根目录文档**：核心入口文档（README、CHANGELOG、BUILD_GUIDE等）
- **功能文档**：用户文档，按功能分类（`docs/features/`）
- **指南文档**：使用说明，按难度分层（`docs/guides/`）
- **技术文档**：开发者文档，问题修复和技术报告（`technical/`）
- **归档文档**：历史和临时文档，按类型归档（`archive/`）

#### 📝 命名规范
- ✅ **小写 + 连字符**：`history-system.md`
- ✅ **语义化命名**：文件名应清晰表达内容
- ✅ **英文优先**：主要使用英文（避免中文文件名）
- ✅ **版本无关**：除非是版本专属文档

#### 🔧 维护流程
1. **新增功能**：创建功能文档 + 快速开始指南 + 更新索引
2. **版本发布**：创建发布文档 + 更新变更日志 + 归档临时文档
3. **文档归档**：按类型放入对应archive子目录
4. **定期审查**：每年检查文档时效性，清理过时内容

### 📚 历史文档

如需查看详细的文档治理报告和重组过程文档，请访问：
- [文档治理报告](./archive/one-time-reports/DOCS_REORGANIZATION_SUCCESS.md)
- [重组完成报告](./archive/one-time-reports/REORGANIZATION_COMPLETE.md)
- [文档结构总览（历史版本）](./archive/one-time-reports/DOCUMENTATION_STRUCTURE.md)

---

**最后更新**: 2025-10-12
**文档版本**: v2.7.0
**维护状态**: 🟢 活跃维护中

