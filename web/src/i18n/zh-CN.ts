import type { en } from './en'

// 简体中文 locale。
export const zhCN: Record<keyof typeof en, string> = {
  // Sidebar
  'sidebar.currentSession': '当前会话',
  'sidebar.provider': '提供商',
  'sidebar.model': '模型',
  'sidebar.status': '状态',
  'sidebar.connected': '已连接',
  'sidebar.idle': '空闲',
  'sidebar.tokenUsage': 'Token 用量',
  'sidebar.input': '输入',
  'sidebar.output': '输出',
  'sidebar.cached': '缓存',
  'sidebar.total': '总计',
  'sidebar.request': '请求',
  'sidebar.requests': '请求数',
  'sidebar.lastResponse': '最后响应',
  'sidebar.firstToken': '首字延迟',

  // Navigation
  'nav.clearSession': '清空会话',
  'nav.settings': '设置',
  'nav.about': '关于',

  // Chat header
  'header.noProvider': '未配置 Provider',
  'header.noModel': '未配置 Model',
  'header.menu': '菜单',

  // Empty state
  'chat.emptyTitle': '开始新的对话',
  'chat.emptyDesc': '选择提供商和模型，然后在下方输入内容开始对话。',

  // Composer
  'composer.placeholder': '输入消息…（Shift + Enter 换行）',
  'composer.send': '发送',
  'composer.stop': '停止',

  // Message
  'message.you': '你',

  // Code block
  'code.copy': '复制',
  'code.copied': '已复制',

  // Settings
  'settings.title': '设置',
  'settings.chat': '对话',
  'settings.streamResponse': '流式输出',
  'settings.appearance': '外观',
  'settings.theme': '主题',
  'settings.light': '浅色',
  'settings.dark': '暗色',
  'settings.system': '跟随系统',
  'settings.language': '语言',
  'settings.english': 'English',
  'settings.chinese': '简体中文',
  'settings.advanced': '高级',
  'settings.default': '默认',
  'settings.configuration': '配置',
  'settings.reloadHint': '修改 config.json 后点击重新加载，刷新 Provider 和 Model 列表。',
  'settings.reload': '重新加载配置',
  'settings.reloading': '加载中…',
  'settings.reloaded': '已重新加载',
  'settings.reloadFailed': '重载失败',

  // About
  'about.tagline': '一个简单、轻量、自托管的 AI 对话工具。',
  'about.version': '版本',
  'about.build': '构建',
  'about.github': 'GitHub',
  'about.releases': 'Release',
  'about.viewReleases': '查看 Release',
  'about.viewRelease': '查看此版本',
  'about.license': '许可证',
  'about.description':
    '一个自托管的 AI 对话工具。无需注册、无需数据库、无需云服务——通过你自己的 API 与不同模型对话。',
  'about.apiKeyNote': 'API Key 仅保存在 EasyTalk 服务端，不会发送到浏览器。',

  // Errors
  'error.invalidKey': 'API Key 无效，请检查 Provider 配置。',
  'error.providerNotFound': '未找到 Provider「{provider}」。',
  'error.rateLimit': '请求过于频繁，请稍后再试。',
  'error.connectFailed': '无法连接到 {provider}，请检查 Base URL、API Key 和网络。',
}