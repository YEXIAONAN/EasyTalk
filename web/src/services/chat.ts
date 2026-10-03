import { ApiError } from './api'

export interface ChatParams {
  provider: string
  model: string
  messages: Array<{ role: string; content: string }>
  stream: boolean
  temperature?: number
  max_tokens?: number
  top_p?: number
}

interface ChatOptions {
  signal?: AbortSignal
  onDelta?: (delta: string) => void
}

// sendChat proxies a chat request through the EasyTalk backend. For streaming
// requests it parses OpenAI-compatible SSE and invokes onDelta for each content
// chunk; for non-streaming requests onDelta receives the full reply once.
export async function sendChat(params: ChatParams, opts: ChatOptions = {}): Promise<void> {
  const res = await fetch('/api/chat', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(params),
    signal: opts.signal,
  })

  if (!res.ok) {
    let message = res.statusText
    try {
      const body = await res.json()
      if (body?.error) message = body.error
    } catch {
      /* ignore */
    }
    throw new ApiError(res.status, message)
  }

  if (!params.stream) {
    const data = await res.json()
    opts.onDelta?.(data.content ?? '')
    return
  }

  const reader = res.body?.getReader()
  if (!reader) return

  const decoder = new TextDecoder()
  let buffer = ''

  for (;;) {
    const { value, done } = await reader.read()
    if (done) break

    buffer += decoder.decode(value, { stream: true })
    const lines = buffer.split('\n')
    buffer = lines.pop() ?? ''

    for (const line of lines) {
      const trimmed = line.trim()
      if (!trimmed.startsWith('data:')) continue

      const payload = trimmed.slice(5).trim()
      if (payload === '[DONE]') return

      try {
        const json = JSON.parse(payload)
        const delta = json.choices?.[0]?.delta?.content
        if (typeof delta === 'string' && delta) opts.onDelta?.(delta)
      } catch {
        /* ignore malformed chunk */
      }
    }
  }
}

export function friendlyChatError(provider: string, err: unknown): string {
  const status = (err as { status?: number } | null)?.status
  if (status === 401) return 'API Key 无效，请检查 Provider 配置。'
  if (status === 404) return `未找到 Provider「${provider}」。`
  if (status === 429) return '请求过于频繁，请稍后再试。'
  return `无法连接到 ${provider}，请检查 Base URL、API Key 和网络。`
}