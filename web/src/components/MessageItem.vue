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
    <div class="meta">
      <span class="role">{{ label }}</span>
      <span class="time">{{ time }}</span>
    </div>

    <MarkdownRenderer v-if="message.role === 'assistant'" :content="message.content" />
    <div v-else class="content">{{ message.content }}</div>
  </div>
</template>

<style scoped>
.message {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 0;
}

.meta {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.role {
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text);
}

.time {
  font-size: 11.5px;
  color: var(--text-muted);
}

.content {
  font-size: 14px;
  line-height: 1.65;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>