# Apple Music Downloader - 文档中心

欢迎来到 Apple Music Downloader 的文档中心！

---

## 📚 文档导航

### 用户指南
- **[README.md](../README.md)** - 项目主文档（English）
- **[README-CN.md](../README-CN.md)** - 中文说明文档
- **[config.yaml.example](../config.yaml.example)** - 配置文件示例
- **[dev.env.example](../dev.env.example)** - 环境变量示例

### 开发者文档
- **[DEVELOPMENT.md](./DEVELOPMENT.md)** - 开发指南
  - 项目架构说明
  - 编译和测试指南
  - 代码规范和最佳实践
  - 贡献指南

- **[CODE_QUALITY.md](./CODE_QUALITY.md)** - 代码质量报告
  - 代码健康度评估
  - 近期改进总结
  - 待优化事项
  - 性能指标

### 版本信息
- **[CHANGELOG.md](../CHANGELOG.md)** - 完整更新日志
- **[VERSION](../VERSION)** - 当前版本号
- **[LICENSE](../LICENSE)** - MIT 开源许可

---

## 🚀 快速开始

### 新用户
1. 阅读 [README-CN.md](../README-CN.md) 了解项目功能
2. 参考 [config.yaml.example](../config.yaml.example) 配置参数
3. 查看 [dev.env.example](../dev.env.example) 设置环境变量

### 开发者
1. 阅读 [DEVELOPMENT.md](./DEVELOPMENT.md#快速开始) 搭建开发环境
2. 查看 [CODE_QUALITY.md](./CODE_QUALITY.md) 了解代码规范
3. 参考 [CHANGELOG.md](../CHANGELOG.md) 了解版本演进

---

## 📖 文档分类

### 按主题分类

#### 🎯 配置相关
- [config.yaml.example](../config.yaml.example) - 完整配置示例
- [dev.env.example](../dev.env.example) - 环境变量配置
- [DEVELOPMENT.md#配置系统](./DEVELOPMENT.md#配置系统) - 配置系统详解

#### 🔧 开发相关
- [DEVELOPMENT.md#项目架构](./DEVELOPMENT.md#项目架构) - 项目架构
- [DEVELOPMENT.md#测试](./DEVELOPMENT.md#测试) - 测试指南
- [CODE_QUALITY.md#代码规范](./CODE_QUALITY.md) - 代码规范

#### 📦 版本相关
- [CHANGELOG.md](../CHANGELOG.md) - 所有版本更新日志
- [VERSION](../VERSION) - 当前版本号

#### 🤝 贡献相关
- [DEVELOPMENT.md#贡献指南](./DEVELOPMENT.md#贡献指南) - 如何贡献代码
- [CODE_QUALITY.md#待优化事项](./CODE_QUALITY.md#待优化事项) - 可以改进的地方

---

## 🔍 常见问题

### Q1: 如何配置多账号？
**A**: 参考 [config.yaml.example](../config.yaml.example) 中的多账号配置示例（CN/US/JP）

### Q2: 如何启用调试模式？
**A**: 在配置文件中设置 `debug: true` 或使用命令行参数 `--debug`

### Q3: MV 分辨率如何设置？
**A**: 在配置文件中设置 `mv-max` (上限) 和 `mv-min` (下限)，详见 [CHANGELOG.md](../CHANGELOG.md#未发布---v131)

### Q4: 如何贡献代码？
**A**: 阅读 [DEVELOPMENT.md#贡献指南](./DEVELOPMENT.md#贡献指南)，了解提交规范和开发流程

### Q5: 在哪里查看最新更新？
**A**: 查看 [CHANGELOG.md](../CHANGELOG.md) 了解所有版本的更新内容

---

## 📝 文档维护

### 文档更新原则
1. **同步更新**: 代码变更时同步更新相关文档
2. **清晰简洁**: 使用简洁的语言和清晰的示例
3. **分类明确**: 按主题和用户类型组织文档
4. **避免冗余**: 避免重复内容，使用链接引用

### 文档贡献
如果您发现文档有误或需要改进，欢迎：
- 提交 Issue 指出问题
- 提交 Pull Request 改进文档
- 在讨论区分享您的使用经验

---

## 🔗 外部资源

### 官方文档
- [Go 官方文档](https://golang.org/doc/)
- [Apple Music API](https://developer.apple.com/documentation/applemusicapi)

### 相关技术
- [FFmpeg 文档](https://ffmpeg.org/documentation.html)
- [M3U8 格式说明](https://datatracker.ietf.org/doc/html/rfc8216)

---

## 📧 联系方式

- **项目地址**: https://github.com/onlyhooops/apple-music-downloader
- **问题反馈**: https://github.com/onlyhooops/apple-music-downloader/issues
- **讨论区**: https://github.com/onlyhooops/apple-music-downloader/discussions

---

**最后更新**: 2025-11-04  
**维护者**: Apple Music Downloader Team
