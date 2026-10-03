<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import type { ChatMessage, Provider } from '../types'
import ChatHeader from './ChatHeader.vue'
import MessageItem from './MessageItem.vue'
import Composer from './Composer.vue'

const props = defineProps<{
  providers: Provider[]
  provider: string
  model: string
  sending: boolean
  messages: ChatMessage[]
}>()

const emit = defineEmits<{
  'update:provider': [name: string]
  'update:model': [name: string]
  send: [content: string]
  stop: []
  settings: []
  toggleMenu: []
}>()

const listEl = ref<HTMLElement | null>(null)

// Auto-scroll to the bottom as messages arrive and stream in.
watch(
  () => [props.messages.length, props.messages[props.messages.length - 1]?.content],
  () => {
    nextTick(() => {
      listEl.value?.scrollTo({ top: listEl.value.scrollHeight })
    })
  },
)
</script>

<template>
  <div class="chat">
    <ChatHeader
      :providers="providers"
      :provider="provider"
      :model="model"
      @update:provider="emit('update:provider', $event)"
      @update:model="emit('update:model', $event)"
      @settings="emit('settings')"
      @toggle-menu="emit('toggleMenu')"
    />

    <div ref="listEl" class="messages">
      <div v-if="messages.length === 0" class="empty">
        <div class="empty-mark">ET</div>
        <h2>开始新的对话</h2>
        <p>配置 Provider 并选择模型后，在下方输入内容开始对话。</p>
      </div>

      <div v-else class="list">
        <MessageItem
          v-for="m in messages"
          :key="m.id"
          :message="m"
          :provider-name="provider"
          :model-name="model"
        />
      </div>
    </div>

    <Composer :sending="sending" @send="emit('send', $event)" @stop="emit('stop')" />
  </div>
</template>

<style scoped>
.chat {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.messages {
  flex: 1;
  overflow-y: auto;
}

.empty {
  margin: auto;
  text-align: center;
  padding: 24px;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.empty-mark {
  width: 44px;
  height: 44px;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: var(--primary-soft);
  color: var(--primary);
  font-size: 16px;
  font-weight: 700;
}

.empty h2 {
  font-size: 18px;
  font-weight: 650;
  margin: 0 0 6px;
}

.empty p {
  margin: 0;
  font-size: 13.5px;
  color: var(--text-secondary);
}

.list {
  max-width: 820px;
  margin: 0 auto;
  padding: 16px 24px 8px;
}
</style>