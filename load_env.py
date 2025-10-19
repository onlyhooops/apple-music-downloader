#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Apple Music 下载器环境变量加载脚本
使用Python进行更安全的环境变量替换
"""

import os
import sys
import re
from pathlib import Path

def load_env_file(env_file):
    """加载环境变量文件"""
    if not env_file.exists():
        print(f"错误: 未找到环境变量文件 {env_file}")
        print("请复制 dev.env.example 为 dev.env 并填入真实配置信息")
        sys.exit(1)

    # 读取环境变量文件
    with open(env_file, 'r', encoding='utf-8') as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith('#'):
                continue

            if '=' in line:
                key, value = line.split('=', 1)
                key = key.strip()
                value = value.strip()

                # 移除引号
                if value.startswith('"') and value.endswith('"'):
                    value = value[1:-1]
                elif value.startswith("'") and value.endswith("'"):
                    value = value[1:-1]

                os.environ[key] = value

def replace_env_vars(config_file, output_file):
    """替换配置文件中的环境变量引用"""
    with open(config_file, 'r', encoding='utf-8') as f:
        content = f.read()

    # 替换环境变量引用 ${VAR_NAME} 或 $VAR_NAME
    def replace_var(match):
        var_name = match.group(1)
        return os.environ.get(var_name, match.group(0))

    # 使用正则表达式替换 ${VAR_NAME} 格式
    pattern = r'\$\{([^}]+)\}'
    result = re.sub(pattern, replace_var, content)

    # 替换 $VAR_NAME 格式（如果有的话）
    pattern2 = r'\$([A-Za-z_][A-Za-z0-9_]*)'
    result = re.sub(pattern2, replace_var, result)

    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(result)

def main():
    """主函数"""
    project_root = Path(__file__).parent
    env_file = project_root / 'dev.env'
    config_file = project_root / 'config.yaml'
    temp_config_file = project_root / 'config.yaml.tmp'

    print("正在加载环境变量...")
    load_env_file(env_file)

    # 检查必需的环境变量
    required_vars = [
        'APPLE_MUSIC_MEDIA_USER_TOKEN_CN',
        'APPLE_MUSIC_AUTH_TOKEN_CN'
    ]

    for var in required_vars:
        if var not in os.environ:
            print(f"错误: 必需的环境变量 {var} 未设置")
            print(f"请检查 {env_file} 文件中的配置")
            sys.exit(1)

    print("正在处理配置文件...")
    replace_env_vars(config_file, temp_config_file)

    print("环境变量加载完成")
    print(f"临时配置文件已创建: {temp_config_file}")
    print("程序可以使用此临时配置文件运行")
    print()
    print("已加载的环境变量:")
    for var in required_vars:
        print(f"- {var}: 已设置")

    local_dev = os.environ.get('LOCAL_DEVELOPMENT', 'false')
    print(f"- LOCAL_DEVELOPMENT: {local_dev}")

    print()
    print("使用说明:")
    print(f"1. 程序启动时请使用临时配置文件: {temp_config_file}")
    print("2. 或者直接在程序中引用环境变量")
    print("3. 程序运行结束后可删除临时配置文件")
if __name__ == '__main__':
    main()
