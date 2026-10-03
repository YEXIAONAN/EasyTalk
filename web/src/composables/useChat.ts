import { computed, ref } from 'vue'
import type { ChatMessage, Conversation } from '../types'
import { sendChat, friendlyChatError } from '../services/chat'
import * as db from '../services/db'

const conversations = ref<Conversation[]>([])
const activeId = ref<string | null>(null)
const sending = ref(false)

let controller: AbortController | null = null

const activeConversation = computed(
  () => conversations.value.find((c) => c.id === activeId.value) ?? null,
)

const messages = computed<ChatMessage[]>(() => activeConversation.value?.messages ?? [])

function newId(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    return crypto.randomUUID()
  }
  return Date.now().toString(36) + Math.random().toString(36).slice(2, 10)
}

function isAbort(err: unknown): boolean {
  return (err as { name?: string } | null)?.name === 'AbortError'
}

function makeTitle(content: string): string {
  const first = content.trim().split('\n')[0].trim()
  return first.length > 30 ? first.slice(0, 30) + '…' : first
}

async function load() {
  const list = await db.getAllConversations()
  list.sort((a, b) => b.updatedAt - a.updatedAt)
  conversations.value = list
}

function select(id: string) {
  activeId.value = id
}

async function persist(conv: Conversation) {
  conv.updatedAt = Date.now()
  await db.saveConversation(conv)
}

async function send(provider: string, model: string, content: string) {
  if (sending.value || !provider || !model) return

  sending.value = true
  controller = new AbortController()

  // Create a conversation on the first message when none is active.
  if (!activeConversation.value) {
    conversations.value.unshift({
      id: newId(),
      title: makeTitle(content),
      createdAt: Date.now(),
      updatedAt: Date.now(),
      messages: [],
    })
    activeId.value = conversations.value[0].id
  }

  const conv = activeConversation.value!

  conv.messages.push({ id: newId(), role: 'user', content, createdAt: Date.now() })
  conv.messages.push({ id: newId(), role: 'assistant', content: '', createdAt: Date.now() })

  // Mutate the reactive element so the UI updates while deltas stream in.
  const assistant = conv.messages[conv.messages.length - 1]

  const apiMessages = conv.messages
    .filter((m) => m.id !== assistant.id)
    .map((m) => ({ role: m.role, content: m.content }))

  void persist(conv)

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
    if (!isAbort(err)) {
      assistant.content = friendlyChatError(provider, err)
    }
  } finally {
    sending.value = false
    controller = null
    void persist(conv)
  }
}

function stop() {
  controller?.abort()
}

export function useChat() {
  return {
    conversations,
    activeId,
    activeConversation,
    messages,
    sending,
    load,
    select,
    send,
    stop,
  }
}