import { ref } from 'vue'
import type { ChatMessage } from '../types'
import { sendChat, friendlyChatError } from '../services/chat'
import { useSettings } from './useSettings'

const { settings } = useSettings()

// The current session lives only in memory. A page refresh clears it.
const messages = ref<ChatMessage[]>([])
const sending = ref(false)

// Session metrics (in-memory only; reset on refresh / clear session).
const requestCount = ref(0)
const inputTokens = ref(0)
const outputTokens = ref(0)
const totalTokens = ref(0)
const cachedTokens = ref<number | null>(null)
const lastResponseMs = ref<number | null>(null)
const firstTokenMs = ref<number | null>(null)

let controller: AbortController | null = null

function newId(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    return crypto.randomUUID()
  }
  return Date.now().toString(36) + Math.random().toString(36).slice(2, 10)
}

function isAbort(err: unknown): boolean {
  return (err as { name?: string } | null)?.name === 'AbortError'
}

async function send(provider: string, model: string, content: string) {
  if (sending.value || !provider || !model) return

  sending.value = true
  controller = new AbortController()

  requestCount.value += 1
  const startedAt = performance.now()
  const streamMode = settings.streamResponse
  let gotFirstToken = false
  firstTokenMs.value = null

  messages.value.push({ id: newId(), role: 'user', content, createdAt: Date.now() })
  messages.value.push({ id: newId(), role: 'assistant', content: '', createdAt: Date.now() })

  // Mutate the reactive element so the UI updates while deltas stream in.
  const assistant = messages.value[messages.value.length - 1]

  const apiMessages = messages.value
    .filter((m) => m.id !== assistant.id)
    .map((m) => ({ role: m.role, content: m.content }))

  try {
    await sendChat(
      {
        provider,
        model,
        messages: apiMessages,
        stream: streamMode,
        temperature: settings.temperature ?? undefined,
        max_tokens: settings.maxTokens ?? undefined,
        top_p: settings.topP ?? undefined,
      },
      {
        signal: controller.signal,
        onDelta: (delta) => {
          if (streamMode && !gotFirstToken) {
            firstTokenMs.value = Math.round(performance.now() - startedAt)
            gotFirstToken = true
          }
          assistant.content += delta
        },
        onUsage: (usage) => {
          inputTokens.value += usage.input_tokens
          outputTokens.value += usage.output_tokens
          totalTokens.value += usage.total_tokens
          if (usage.cached_tokens !== undefined) {
            cachedTokens.value = (cachedTokens.value ?? 0) + usage.cached_tokens
          }
        },
      },
    )
  } catch (err) {
    if (!isAbort(err)) {
      assistant.content = friendlyChatError(provider, err)
    }
  } finally {
    lastResponseMs.value = Math.round(performance.now() - startedAt)
    sending.value = false
    controller = null
  }
}

function stop() {
  controller?.abort()
}

function clear() {
  if (sending.value) return
  messages.value = []
  requestCount.value = 0
  inputTokens.value = 0
  outputTokens.value = 0
  totalTokens.value = 0
  cachedTokens.value = null
  lastResponseMs.value = null
  firstTokenMs.value = null
}

export function useChat() {
  return {
    messages,
    sending,
    requestCount,
    inputTokens,
    outputTokens,
    totalTokens,
    cachedTokens,
    lastResponseMs,
    firstTokenMs,
    send,
    stop,
    clear,
  }
}