#!/usr/bin/env bash
# EasyTalk 一键启动脚本（Linux）
set -euo pipefail

# 定位脚本所在目录与项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$ROOT_DIR"

# 首次运行：从示例创建 config.json
if [ ! -f config.json ]; then
  if [ -f config.example.json ]; then
    cp config.example.json config.json
    echo "[EasyTalk] 已从 config.example.json 创建 config.json，请填入你的 API Key。"
  else
    echo "[EasyTalk] 未找到 config.example.json，请手动创建 config.json。"
    exit 1
  fi
fi

# 首次运行：从示例创建 .env（可选）
if [ ! -f .env ]; then
  if [ -f .env.example ]; then
    cp .env.example .env
    echo "[EasyTalk] 已从 .env.example 创建 .env。"
  fi
fi

# 定位可执行文件
BIN=""
for c in "$ROOT_DIR/easytalk" "$ROOT_DIR/bin/easytalk" "$SCRIPT_DIR/easytalk"; do
  if [ -x "$c" ]; then
    BIN="$c"
    break
  fi
done

if [ -z "$BIN" ]; then
  echo "[EasyTalk] 未找到 easytalk 可执行文件，请先构建（make build）或下载 Release。"
  exit 1
fi

exec "$BIN"