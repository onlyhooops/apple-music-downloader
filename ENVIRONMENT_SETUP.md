# 本地开发环境配置指南

本项目使用环境变量来隔离敏感配置信息，避免将账号令牌、密码等敏感数据提交到版本控制系统中。

## 文件结构

```
apple-music-downloader/
├── dev.env              # 本地环境变量文件（包含敏感信息）
├── dev.env.example      # 环境变量模板文件
├── config.yaml          # 主配置文件（引用环境变量）
├── config.yaml.example  # 配置模板文件
├── load_env.sh         # 环境变量加载脚本
└── ENVIRONMENT_SETUP.md # 本指南文档
```

## 快速开始

### 1. 配置环境变量

```bash
# 复制环境变量模板
cp dev.env.example dev.env

# 编辑环境变量文件，填入真实信息
nano dev.env
```

在 `dev.env` 文件中填入你的真实配置：

```bash
# Apple Music 账号令牌
APPLE_MUSIC_MEDIA_USER_TOKEN_CN=你的真实媒体用户令牌
APPLE_MUSIC_AUTH_TOKEN_CN=你的真实授权令牌
```

### 2. 加载环境变量

```bash
# 运行环境变量加载脚本
./load_env.sh
```

此脚本将：
- 加载环境变量
- 检查必需的环境变量是否已设置
- 创建临时的配置文件 `config.yaml.tmp`，其中环境变量引用已被替换为实际值

### 3. 运行程序

使用临时配置文件运行程序：

```bash
# 使用临时配置文件运行程序
./apple-music-downloader -config config.yaml.tmp

# 或者直接在程序中引用环境变量（取决于程序的实现）
```

## 环境变量说明

| 环境变量 | 说明 | 是否必需 |
|---------|------|---------|
| `APPLE_MUSIC_MEDIA_USER_TOKEN_CN` | Apple Music 中国区媒体用户令牌 | 是 |
| `APPLE_MUSIC_AUTH_TOKEN_CN` | Apple Music 中国区授权令牌 | 否（可选） |

## 获取令牌指南

### Apple Music 媒体用户令牌

1. 打开 Apple Music 网页版并登录你的账号
2. 按 F12 打开开发者工具
3. 转到 Application -> Cookies -> https://music.apple.com
4. 复制 `media-user-token` 的值
5. 将此值填入 `APPLE_MUSIC_MEDIA_USER_TOKEN_CN` 环境变量

## 安全注意事项

1. **永远不要**将 `dev.env` 文件提交到版本控制系统
2. **永远不要**在公开场合分享你的令牌和密码
3. 定期更换你的令牌和密码
4. 考虑使用密码管理器来管理这些凭证

## 故障排除

### 环境变量未加载

```bash
# 检查环境变量文件是否存在且格式正确
ls -la dev.env
cat dev.env

# 手动加载环境变量
source dev.env

# 检查环境变量是否设置成功
echo $APPLE_MUSIC_MEDIA_USER_TOKEN_CN
```

### 程序无法读取配置

- 确保运行了 `load_env.sh` 脚本
- 检查临时配置文件 `config.yaml.tmp` 是否存在
- 使用 `-config config.yaml.tmp` 参数指定临时配置文件

### 令牌无效

- 重新从 Apple Music 网站获取令牌
- 检查令牌是否完整复制（令牌通常很长）
- 确保令牌没有过期

## 最佳实践

1. 将 `dev.env` 添加到 `.gitignore` 文件中（如果使用 Git）
2. 使用强密码和定期更换令牌
3. 仅在受信任的环境中运行程序
4. 考虑使用容器化部署来进一步隔离敏感信息

## 高级配置

如果你需要为不同的环境（开发、测试、生产）使用不同的配置，可以创建多个环境文件：

- `dev.env` - 开发环境
- `test.env` - 测试环境
- `prod.env` - 生产环境

相应地修改加载脚本以支持不同的环境。
