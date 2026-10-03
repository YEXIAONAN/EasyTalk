import { ref } from 'vue'
import type { ChatMessage } from '../types'
import { sendChat, friendlyChatError } from '../services/chat'

const messages = ref<ChatMessage[]>([])
const sending = ref(false)

function newId(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    return crypto.randomUUID()
  }
  return Date.now().toString(36) + Math.random().toString(36).slice(2, 10)
}

async function send(provider: string, model: string, content: string) {
  if (sending.value || !provider || !model) return

  sending.value = true

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
        onDelta: (delta) => {
          assistant.content += delta
        },
      },
    )
  } catch (err) {
    assistant.content = friendlyChatError(provider, err)
  } finally {
    sending.value = false
  }
}

export function useChat() {
  return { messages, sending, send }
}