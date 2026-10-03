<script setup lang="ts">
import type { Provider } from '../types'
import ChatHeader from './ChatHeader.vue'
import Composer from './Composer.vue'

defineProps<{
  providers: Provider[]
  provider: string
  model: string
  sending: boolean
}>()

const emit = defineEmits<{
  'update:provider': [name: string]
  'update:model': [name: string]
  send: [content: string]
  stop: []
  settings: []
}>()
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
    />

    <div class="messages">
      <div class="empty">
        <div class="empty-mark">ET</div>
        <h2>开始新的对话</h2>
        <p>配置 Provider 并选择模型后，在下方输入内容开始对话。</p>
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
  display: flex;
  flex-direction: column;
}

.empty {
  margin: auto;
  text-align: center;
  padding: 24px;
}

.empty-mark {
  width: 44px;
  height: 44px;
  margin: 0 auto 16px;
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
</style>