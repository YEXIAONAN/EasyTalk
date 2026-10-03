import { ref } from 'vue'
import type { ChatMessage } from '../types'
import { sendChat, friendlyChatError } from '../services/chat'

const messages = ref<ChatMessage[]>([])
const sending = ref(false)

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

  messages.value.push({ id: newId(), role: 'user', content, createdAt: Date.now() })
  messages.value.push({ id: newId(), role: 'assistant', content: '', createdAt: Date.now() })

  // Mutate the reactive element (not a detached copy) so the UI updates while
  // deltas stream in.
  const assistant = messages.value[messages.value.length - 1]

  const apiMessages = messages.value
    .filter((m) => m.id !== assistant.id)
    .map((m) => ({ role: m.role, content: m.content }))

  try {
    await sendChat(
      { provider, model, messages: apiMessages, stream: true },
      {
        signal: controller.signal,
        onDelta: (delta) => {
          assistant.content += delta
        },
      },
    )
  } catch (err) {
    // A user-initiated stop keeps whatever partial content was produced.
    if (!isAbort(err)) {
      assistant.content = friendlyChatError(provider, err)
    }
  } finally {
    sending.value = false
    controller = null
  }
}

function stop() {
  controller?.abort()
}

export function useChat() {
  return { messages, sending, send, stop }
}