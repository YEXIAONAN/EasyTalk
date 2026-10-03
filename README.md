# EasyTalk

> EasyTalk 是一个简单、轻量、自托管的 AI 对话工具，通过统一界面连接你自己的 AI API。

> EasyTalk is a simple and lightweight self-hosted AI chat client.

EasyTalk 不是平台，也不是模型管理工具。它的目标只有一个：**让你用自己的 API，在简洁舒服的网页里和不同模型对话。**

## Features

- 支持任意 OpenAI 兼容服务（`POST /chat/completions`）
- 流式 / 非流式对话，支持中途停止生成
- Markdown 渲染：标题、列表、表格、引用、行内代码、代码块（语法高亮 + 一键复制）
- 左侧 Current Session 实时面板：Provider / Model / Status、Token Usage、Request 指标
- Token Usage 归一化（产物统一 `input_tokens`/`output_tokens`/`total_tokens`/`cached_tokens`，未知显示 `—`）
- 浅色 / 暗色 / 跟随系统 三种主题
- Temperature / Max Tokens / Top P 高级参数
- 运行时重载配置（Reload Config）
- 局域网访问，手机 / 平板可通过 LAN 打开
- 响应式布局，移动端侧边栏抽屉
- 单文件部署：前端资源通过 `go:embed` 编译进一个可执行文件

## Quick Start

```bash
cp .env.example .env
cp config.example.json config.json
```

编辑 `config.json`，填入你自己的 Provider（Base URL 与 API Key）：

```json
{
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

然后启动：

```bash
./easytalk
```

打开浏览器访问：

```text
http://localhost:8080
```

## Configuration

EasyTalk 使用两个职责分离的配置文件：

### `.env` —— 如何启动

保存运行环境参数，示例见 [`.env.example`](.env.example)：

```text
EASYTALK_HOST=0.0.0.0
EASYTALK_PORT=8080
EASYTALK_CONFIG=./config.json
```

配置优先级：命令行参数 > 环境变量 > `.env` > 默认值。

### `config.json` —— 可以调用哪些 Provider

保存 Provider 列表，示例见 [`config.example.json`](config.example.json)。

> 两个真实文件（`.env` 与 `config.json`）均已加入 `.gitignore`，绝不要把真实 API Key 提交到 Git。

如果启动时 `config.json` 不存在，EasyTalk 会输出清晰提示并安全退出，不会 panic。

修改 `config.json` 后，无需重启：在 Settings → Configuration 点击 **Reload Config** 即可刷新 Provider 与 Model。

## Providers

任何实现了 OpenAI 兼容接口的服务都可以使用，例如：

| Provider | Base URL |
| --- | --- |
| DeepSeek | `https://api.deepseek.com/v1` |
| OpenRouter | `https://openrouter.ai/api/v1` |
| 本地 Ollama | `http://127.0.0.1:11434/v1` |

**API Key 永远只保存在后端。** 浏览器只拿到 Provider 名称与模型列表，请求始终经 EasyTalk 后端转发，Key 不会出现在 HTML / JavaScript / LocalStorage / 网络响应中。

## Chat Session

聊天不持久化。当前会话只保存在内存中，**刷新页面即清空**——这是设计行为，不是 Bug。

`Clear Session` 会清空当前消息与用量统计，不产生任何历史记录。

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

EasyTalk 默认监听 `0.0.0.0`，同一局域网内的设备可通过 `http://<主机局域网 IP>:8080` 访问。启动时终端会打印本机局域网 IP。

## Project Structure

```text
EasyTalk/
├── cmd/easytalk/main.go        # 入口：加载 .env/config、启动服务、优雅退出
├── internal/
│   ├── config/                 # .env 与 config.json 加载
│   ├── server/                 # HTTP 路由、Provider / Config / Chat 处理器
│   ├── provider/               # OpenAI 兼容客户端（Chat / Stream / Usage 归一化）
│   └── version/                # 版本号（可用 ldflags 覆盖）
├── web/                        # Vue 3 + Vite + TypeScript 前端
│   ├── public/                 # favicon.svg、logo-512.png、apple-touch-icon.png
│   └── src/
│       ├── assets/branding/    # Logo SVG
│       ├── components/         # 组件
│       ├── views/              # 页面（Chat / Settings / About）
│       ├── composables/        # 状态逻辑（providers / chat / settings）
│       ├── services/           # API、Markdown 渲染
│       ├── config/             # 品牌资源统一引用
│       └── types/              # 类型定义
├── embed.go                    # go:embed 嵌入 web/dist
├── .env.example
├── config.example.json
├── Makefile
└── README.md
```

## License

[MIT](LICENSE)