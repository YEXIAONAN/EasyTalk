# EasyTalk

> EasyTalk 是一个简单、轻量、自托管的 AI 对话工具，通过统一界面连接你自己的 AI API。

> EasyTalk is a simple and lightweight self-hosted AI chat client.

EasyTalk 不是 AI 平台，也不是模型管理工具。它的目标只有一个：**让你用自己的 API，在一个简洁舒服的网页里和不同模型对话**。

## Features

- Provider 配置管理：Base URL / API Key / Models 的添加、编辑、删除、连接测试
- 兼容任意 OpenAI 兼容接口（`POST /chat/completions`）
- 流式与非流式对话，支持中途停止生成
- Markdown 渲染：标题、列表、表格、引用、行内代码、代码块（语法高亮 + 一键复制）
- 会话历史保存在浏览器本地（IndexedDB），支持新建、切换、重命名、删除、清空
- 浅色 / 暗色 / 跟随系统 三种主题
- 默认 Provider / Model、Temperature / Max Tokens / Top P 等本地设置
- 局域网访问，手机 / 平板可通过 LAN 打开
- 响应式布局，移动端侧边栏抽屉
- 单文件部署：前端资源通过 `go:embed` 编译进一个可执行文件

## Quick Start

```bash
./easytalk
```

启动后终端会显示本地和局域网地址：

```text
EasyTalk v0.1.0

Local:
http://127.0.0.1:8080

LAN:
http://192.168.1.10:8080

Config:
./config.json
```

首次运行会自动在同目录生成 `config.json`。打开网页，在 Settings → Providers 中添加你的 Provider 即可开始对话。

## Configuration

配置文件示例见 [`config.example.json`](config.example.json)。核心结构如下：

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8080
  },
  "providers": [
    {
      "name": "DeepSeek",
      "base_url": "https://api.deepseek.com/v1",
      "api_key": "YOUR_API_KEY",
      "models": ["deepseek-chat", "deepseek-reasoner"]
    }
  ]
}
```

- `server.host` / `server.port`：监听地址与端口，默认 `0.0.0.0:8080` 便于局域网访问
- `providers`：Provider 列表，API Key 以 `0600` 权限保存，日志与接口响应中仅显示掩码（`****xxxx`）

> 也可以用 `./easytalk -config /path/to/config.json` 指定配置文件位置。

## Providers

只要实现了 OpenAI 兼容接口的服务都可以使用，例如：

| Provider | Base URL |
| --- | --- |
| DeepSeek | `https://api.deepseek.com/v1` |
| OpenRouter | `https://openrouter.ai/api/v1` |
| 本地 Ollama | `http://127.0.0.1:11434/v1` |

API Key 仅保存在你的设备上，浏览器请求始终经过 EasyTalk 后端转发，不会在前端代码中暴露。

## Development

```bash
# 后端（需要先构建前端产物 web/dist）
cd web && npm install && npm run build
cd ..
go run ./cmd/easytalk

# 前端热更新开发（Vite dev server + 代理到 8080）
cd web && npm run dev
```

## Build

```bash
make build
```

等价于：

```bash
cd web && npm install && npm run build   # 构建前端到 web/dist
go build -o easytalk ./cmd/easytalk       # 将前端资源嵌入单个二进制
```

产物 `easytalk` 是一个自包含可执行文件，无需 Node.js / Python / Docker / 数据库。

## LAN Usage

EasyTalk 默认监听 `0.0.0.0`，同一局域网内的设备（手机、平板、其他电脑）可直接通过 `http://<主机局域网 IP>:8080` 访问。启动时终端会打印本机局域网 IP。

## Project Structure

```text
EasyTalk/
├── cmd/easytalk/main.go        # 入口：加载配置、启动服务、打印启动信息
├── internal/
│   ├── config/                 # 配置加载与保存
│   ├── server/                 # HTTP 路由、静态资源、Provider、Chat 处理器
│   ├── provider/               # OpenAI 兼容客户端（Test / Chat / Stream）
│   └── version/                # 版本号（可用 ldflags 覆盖）
├── web/                        # Vue 3 + Vite + TypeScript 前端
│   └── src/
│       ├── components/         # 组件
│       ├── views/              # 页面（Chat / Settings / About）
│       ├── composables/        # 状态逻辑（providers / chat / settings）
│       ├── services/           # API、IndexedDB、Markdown 渲染
│       └── types/              # 类型定义
├── embed.go                    # go:embed 嵌入 web/dist
├── config.example.json
├── Makefile
└── README.md
```

## License

[MIT](LICENSE)