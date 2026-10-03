<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  sending: boolean
}>()

const emit = defineEmits<{
  send: [content: string]
  stop: []
}>()

const text = ref('')
const el = ref<HTMLTextAreaElement | null>(null)

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    submit()
  }
}

function submit() {
  const content = text.value.trim()
  if (!content || props.sending) return
  emit('send', content)
  text.value = ''
  if (el.value) el.value.style.height = 'auto'
}

function autoResize() {
  const ta = el.value
  if (!ta) return
  ta.style.height = 'auto'
  ta.style.height = Math.min(ta.scrollHeight, 200) + 'px'
}
</script>

<template>
  <footer class="composer">
    <div class="box">
      <button class="icon-btn" title="附件">
        <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
          <path
            d="M8 11.5V4.5M5.5 7l2.5-2.5L10.5 7"
            stroke="currentColor"
            stroke-width="1.4"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
          <path d="M3.5 12.5h9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
        </svg>
      </button>

      <textarea
        ref="el"
        v-model="text"
        rows="1"
        class="input"
        placeholder="输入消息...（Shift + Enter 换行）"
        @keydown="onKeydown"
        @input="autoResize"
      ></textarea>

      <button
        v-if="!sending"
        class="send-btn"
        title="发送"
        :disabled="!text.trim()"
        @click="submit"
      >
        <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
          <path
            d="M2.5 8L13 3l-3.5 10-2-4-5-1z"
            stroke="currentColor"
            stroke-width="1.4"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </button>
      <button v-else class="send-btn stop" title="停止" @click="emit('stop')">
        <svg viewBox="0 0 16 16" width="12" height="12" fill="currentColor">
          <rect x="3" y="3" width="10" height="10" rx="1.5" />
        </svg>
      </button>
    </div>
  </footer>
</template>

<style scoped>
.composer {
  padding: 12px 16px 16px;
  flex-shrink: 0;
}

.box {
  display: flex;
  align-items: flex-end;
  gap: 6px;
  max-width: 820px;
  margin: 0 auto;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 8px 10px;
  background: var(--bg-input);
  transition: border-color 0.12s ease;
}

.box:focus-within {
  border-color: var(--primary);
}

.input {
  flex: 1;
  resize: none;
  border: none;
  outline: none;
  font-size: 14px;
  line-height: 1.5;
  max-height: 200px;
  padding: 6px 0;
  background: transparent;
}

.send-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: var(--radius-sm);
  background: var(--primary);
  color: #fff;
  transition: background 0.12s ease;
  flex-shrink: 0;
}

.send-btn:hover:not(:disabled) {
  background: var(--primary-hover);
}

.send-btn.stop {
  background: #ef4444;
}
</style>