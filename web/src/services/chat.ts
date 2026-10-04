import { ApiError } from './api'
import type { Usage } from '../types'
import { t } from '../i18n'

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
  onUsage?: (usage: Usage) => void
}

// sendChat proxies a chat request through the EasyTalk backend. The backend
// normalizes OpenAI-compatible SSE into events carrying a `delta` (content
// chunk) or a `usage` object.
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
    if (data.usage) opts.onUsage?.(data.usage)
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
        if (json.usage) {
          opts.onUsage?.(json.usage)
        } else if (typeof json.delta === 'string' && json.delta) {
          opts.onDelta?.(json.delta)
        }
      } catch {
        /* ignore malformed chunk */
      }
    }
  }
}

export function friendlyChatError(provider: string, err: unknown): string {
  const status = (err as { status?: number } | null)?.status
  if (status === 401) return t('error.invalidKey')
  if (status === 404) return t('error.providerNotFound', { provider })
  if (status === 429) return t('error.rateLimit')
  return t('error.connectFailed', { provider })
}