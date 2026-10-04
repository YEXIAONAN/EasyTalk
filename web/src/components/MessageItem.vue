<script setup lang="ts">
import { computed } from 'vue'
import type { ChatMessage } from '../types'
import { useI18n } from '../i18n'
import MarkdownRenderer from './MarkdownRenderer.vue'

const { t, language } = useI18n()

const props = defineProps<{
  message: ChatMessage
  providerName: string
  modelName: string
}>()

const label = computed(() => {
  return props.message.role === 'user' ? t('message.you') : `${props.providerName} · ${props.modelName}`
})

const time = computed(() => {
  const locale = language.value === 'zh-CN' ? 'zh-CN' : 'en-US'
  return new Date(props.message.createdAt).toLocaleTimeString(locale, {
    hour: '2-digit',
    minute: '2-digit',
  })
})
</script>

<template>
  <div class="message" :class="message.role">
    <div class="bubble">
      <div class="meta">
        <span class="role">{{ label }}</span>
        <span class="time">{{ time }}</span>
      </div>

      <MarkdownRenderer v-if="message.role === 'assistant'" :content="message.content" />
      <div v-else class="content">{{ message.content }}</div>
    </div>
  </div>
</template>

<style scoped>
.message {
  display: flex;
  flex-direction: column;
  padding: 8px 0;
}

.message.assistant {
  align-items: flex-start;
}

.message.user {
  align-items: flex-end;
}

.bubble {
  max-width: 76%;
  padding: 10px 14px;
  border-radius: 14px;
  font-size: 14px;
  line-height: 1.6;
}

.message.assistant .bubble {
  background: var(--bg-sidebar);
  border: 1px solid var(--border);
  border-top-left-radius: 4px;
  color: var(--text);
}

.message.user .bubble {
  background: var(--primary);
  color: #fff;
  border-top-right-radius: 4px;
}

.meta {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 4px;
  font-size: 11.5px;
}

.role {
  font-weight: 600;
}

.message.assistant .role {
  color: var(--text-secondary);
}

.message.assistant .time {
  color: var(--text-muted);
}

.message.user .role {
  color: rgba(255, 255, 255, 0.9);
}

.message.user .time {
  color: rgba(255, 255, 255, 0.65);
}

.content {
  white-space: pre-wrap;
  word-break: break-word;
}

/* Markdown inside the assistant bubble should inherit the bubble text color. */
.message.assistant :deep(.markdown) {
  color: var(--text);
}
</style>