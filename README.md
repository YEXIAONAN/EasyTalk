# EasyTalk

EasyTalk 是一个简单、轻量、自托管的 AI 对话工具，通过统一界面连接你自己的 AI API。

EasyTalk is a simple and lightweight self-hosted AI chat client.

## 简介

- 一个可执行文件即可运行，无需 Node.js / Python / Docker / 数据库
- 通过 OpenAI 兼容接口连接自己的 AI 服务（DeepSeek、OpenRouter、本地 Ollama 等）
- 蓝白配色、简洁界面，支持流式输出
- 聊天历史保存在浏览器本地（IndexedDB）

## Features

- Provider 配置管理（Base URL / API Key / Models）
- Provider 连接测试
- 流式 / 非流式对话
- Markdown 渲染与代码高亮
- 会话历史（本地存储）
- 浅色 / 暗色主题
- 局域网访问

## Quick Start

```bash
./easytalk
```

启动后访问：

```text
Local: http://127.0.0.1:8080
LAN:   http://<你的局域网 IP>:8080
```

首次运行会生成 `config.json`，在其中配置 Provider 后即可开始对话。

## Configuration

见 [config.example.json](./config.example.json)。

## Providers

支持任意 OpenAI 兼容的 `/chat/completions` 接口，例如：

- DeepSeek `https://api.deepseek.com/v1`
- OpenRouter `https://openrouter.ai/api/v1`
- 本地 Ollama `http://127.0.0.1:11434/v1`

## Development

```bash
# 后端
go run ./cmd/easytalk

# 前端
cd web && npm install && npm run dev
```

## Build

```bash
make build
```

## License

MIT