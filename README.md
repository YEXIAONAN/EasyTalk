# EasyTalk

[![CI](https://github.com/YEXIAONAN/EasyTalk/actions/workflows/ci.yml/badge.svg)](https://github.com/YEXIAONAN/EasyTalk/actions/workflows/ci.yml)

> EasyTalk is a simple and lightweight self-hosted AI chat client.

> EasyTalk 是一个简单、轻量、自托管的 AI 对话工具，通过统一界面连接你自己的 AI API。

EasyTalk 不是平台，也不是模型管理工具。它的目标只有一个：**让你用自己的 API，在简洁舒服的网页里和不同模型对话。** 单文件部署，无需数据库、无需 Docker、无需云服务。

## Features

- 支持任意 OpenAI 兼容服务（`POST /chat/completions`）
- 流式 / 非流式对话，支持中途停止生成
- Markdown 渲染：标题、列表、表格、引用、行内代码、代码块（语法高亮 + 一键复制）
- 左侧 Current Session 实时面板：Provider / Model / Status、Token Usage、Request 指标
- Token Usage 归一化（统一 `input_tokens` / `output_tokens` / `total_tokens` / `cached_tokens`，未知显示 `—`）
- 浅色 / 暗色 / 跟随系统 三种主题
- English / 简体中文 双语界面，默认 English，可在 Settings → Language 切换
- 语言偏好保存在本地（localStorage）
- Temperature / Max Tokens / Top P 高级参数
- 运行时重载配置（Reload Config）
- 局域网访问，手机 / 平板可通过 LAN 打开
- 响应式布局，移动端侧边栏抽屉
- 单文件部署：前端资源通过 `go:embed` 编译进一个可执行文件

## Quick Start

### 方式一：下载 Release（推荐）

1. 前往 [Releases](https://github.com/YEXIAONAN/EasyTalk/releases) 下载对应平台的压缩包。
2. 解压后进入目录，首次运行会自动从示例生成配置：

   ```bash
   cp config.example.json config.json       # macOS / Linux
   # 或 Windows（PowerShell）：
   #   Copy-Item config.example.json config.json
   ```

3. 编辑 `config.json`，填入你自己的 Provider（Base URL 与 API Key）。
4. 启动：

   ```bash
   ./easytalk                                # macOS / Linux
   # Windows：双击或运行 easytalk.exe
   ```

5. 浏览器访问 `http://localhost:8080`。

也可以直接使用仓库内的**一键启动脚本**（会自动复制配置并拉起服务）：

- Windows：双击 `start/start.bat`
- macOS：双击 `start/start.command`
- Linux：`./start/start.sh`

### 方式二：从源码构建

见下方 [Development](#development) 与 [Build](#build)。最终 Release 用户**不需要**安装 Go / Node.js / npm。

## Provider Configuration

编辑 `config.json`（示例见 [`config.example.json`](config.example.json)）：

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

字段说明：

| 字段 | 含义 |
| --- | --- |
| `name` | 显示在界面里的 Provider 名称 |
| `base_url` | OpenAI 兼容接口的 Base URL |
| `api_key` | 你的 API Key，**只保存在后端** |
| `models` | 该 Provider 可选的模型列表 |

任何实现 OpenAI 兼容接口的服务都可以使用，例如：

| Provider | Base URL |
| --- | --- |
| DeepSeek | `https://api.deepseek.com/v1` |
| OpenRouter | `https://openrouter.ai/api/v1` |
| 本地 Ollama | `http://127.0.0.1:11434/v1` |

修改 `config.json` 后无需重启：在 Settings → Configuration 点击 **Reload Config** 即可刷新 Provider 与 Model。

## `.env` Configuration

`.env` 只决定 EasyTalk **如何运行**（示例见 [`.env.example`](.env.example)）：

```text
EASYTALK_HOST=0.0.0.0
EASYTALK_PORT=8080
EASYTALK_CONFIG=./config.json
```

- `EASYTALK_HOST`：监听地址，`0.0.0.0` 表示允许局域网访问。
- `EASYTALK_PORT`：监听端口。
- `EASYTALK_CONFIG`：Provider 配置文件路径。

`config.json` 决定**可以调用哪些 AI Provider**，二者职责不同，不要混淆。

> 两个真实文件（`.env` 与 `config.json`）均已加入 `.gitignore`，绝不要把真实 API Key 提交到 Git。配置优先级：命令行参数 > 环境变量 > `.env` > 默认值。

## LAN Usage

EasyTalk 默认监听 `0.0.0.0`，同一局域网内的设备可通过 `http://<主机局域网 IP>:8080` 访问。启动时终端会打印本机局域网 IP。

## Session Behavior

EasyTalk **不保存聊天历史**。聊天只存在于当前浏览器页面内存中，刷新或关闭页面即清空；服务端同样不保存 Conversation History。

这是 Privacy / Simplicity 的**设计选择**，不是 Bug。`Clear Session` 会清空当前消息与用量统计。

> 注意：聊天数据（消息、Token 统计、请求指标）不持久化；但 **Language 与 Theme** 属于界面偏好，会保存在浏览器本地（`localStorage`），与聊天不持久化并不冲突。

## Token Usage

Current Session 面板可展示：

- Input Tokens / Output Tokens / Total Tokens
- Cached Tokens
- Requests
- Response Time / TTFT（首字延迟）

以上数据依赖 Provider 实际返回。若 Provider 未提供缓存信息，则 `Cached Tokens` 显示 `—`——并非所有 Provider 都支持 Cache Metrics。

## Development

```bash
# 后端（需要先构建前端产物 web/dist）
cd web && npm ci && npm run build
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
cd web && npm ci && npm run build        # 构建前端到 web/dist（含 vue-tsc 类型检查）
go build -o easytalk ./cmd/easytalk      # 将前端资源嵌入单个二进制
```

产物 `easytalk` 是一个自包含可执行文件。Vue 仅属于开发 / 构建阶段，最终 Release 用户无需 Node.js。

## CI

[`.github/workflows/ci.yml`](.github/workflows/ci.yml) 在每次 **Push** 与 **Pull Request** 时自动运行，分为两个 Job：

- **Frontend**：`npm ci` → `npm run build`（含 TypeScript 类型检查）
- **Backend**：构建前端 → `go test ./...` → `go vet ./...` → `CGO_ENABLED=0 go build`

普通 Push 只做质量校验，**不会**创建 Tag、Release 或上传任何二进制。

## Release

[`.github/workflows/release.yml`](.github/workflows/release.yml) 只由 **`vX.Y.Z` 格式的 Git Tag** 触发。发布一个版本只需：

```bash
git tag v0.1.0
git push origin v0.1.0
```

GitHub 会自动完成：重新验证 → Vue 生产构建 → Go Embed → 六平台交叉编译 → 打包 → SHA256 校验 → 创建 Draft Release → 上传全部资产 → 发布为 Latest。整个过程无需手动编译或上传二进制。

本地也可复现完整构建（不发布）：

```bash
./scripts/build-release.sh v0.1.0
# 或
make release-local
```

## Supported Platforms

| Platform | Architecture | Archive |
| --- | --- | --- |
| Windows | AMD64 | `easytalk-vX.Y.Z-windows-amd64.zip` |
| Windows | ARM64 | `easytalk-vX.Y.Z-windows-arm64.zip` |
| Linux | AMD64 | `easytalk-vX.Y.Z-linux-amd64.tar.gz` |
| Linux | ARM64 | `easytalk-vX.Y.Z-linux-arm64.tar.gz` |
| macOS | Intel | `easytalk-vX.Y.Z-darwin-amd64.tar.gz` |
| macOS | Apple Silicon | `easytalk-vX.Y.Z-darwin-arm64.tar.gz` |

每次 Release 同时提供 `SHA256SUMS.txt`，供下载后校验完整性。

## Project Structure

```text
EasyTalk/
├── cmd/easytalk/main.go        # 入口：加载 .env/config、启动服务、优雅退出
├── internal/
│   ├── config/                 # .env 与 config.json 加载
│   ├── server/                 # HTTP 路由、Provider / Config / Chat 处理器
│   ├── provider/               # OpenAI 兼容客户端（Chat / Stream / Usage 归一化）
│   └── buildinfo/              # 版本 / Commit / 仓库信息（ldflags 注入，CLI / API / About 统一）
├── web/                        # Vue 3 + Vite + TypeScript 前端
│   ├── public/                 # favicon.svg、logo-512.png、apple-touch-icon.png
│   └── src/
│       ├── assets/branding/    # Logo SVG
│       ├── components/         # 组件
│       ├── views/              # 页面（Chat / Settings / About）
│       ├── composables/        # 状态逻辑（providers / chat / settings）
│       ├── services/           # API、Markdown 渲染
│       ├── i18n/               # 多语言（en / zh-CN）
│       ├── config/             # 品牌资源统一引用
│       └── types/              # 类型定义
├── .github/workflows/          # ci.yml + release.yml
├── scripts/build-release.sh    # 六平台交叉编译 + 打包 + SHA256
├── start/                      # Windows / macOS / Linux 一键启动脚本
├── embed.go                    # go:embed 嵌入 web/dist
├── config.example.json
├── .env.example
├── Makefile
└── README.md
```

## Security

- API Key 只保存在后端 `config.json`，浏览器只拿到 Provider 名称与模型列表。
- 请求始终经 EasyTalk 后端转发，Key 不会出现在 HTML / JavaScript / LocalStorage / 网络响应中。
- `.env` 与 `config.json` 已加入 `.gitignore`；CI / Release 不读取任何真实 Key，测试仅使用 Mock 数据，不调用真实收费 API。
- Release Workflow 仅使用 `${{ github.token }}` 完成发布，不引入第三方 Secrets。

## License

[MIT](LICENSE)