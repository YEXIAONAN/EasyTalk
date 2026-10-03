<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { Provider } from '../types'

const props = defineProps<{
  provider: Provider | null
}>()

const emit = defineEmits<{
  close: []
  save: [p: Provider]
}>()

const form = reactive({
  name: props.provider?.name ?? '',
  base_url: props.provider?.base_url ?? '',
  api_key: props.provider?.api_key ?? '',
  models: (props.provider?.models ?? []).join('\n'),
})

const error = ref('')

function submit() {
  const name = form.name.trim()
  const baseUrl = form.base_url.trim()
  if (!name || !baseUrl) {
    error.value = 'Name 和 Base URL 不能为空'
    return
  }

  const models = form.models
    .split('\n')
    .map((m) => m.trim())
    .filter(Boolean)

  emit('save', {
    name,
    base_url: baseUrl,
    api_key: form.api_key,
    models,
  })
}
</script>

<template>
  <div class="overlay" @click.self="emit('close')">
    <div class="modal">
      <div class="modal-head">
        <h2>{{ provider ? 'Edit Provider' : 'Add Provider' }}</h2>
        <button class="icon-btn" title="关闭" @click="emit('close')">
          <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
            <path d="M4 4l8 8M12 4l-8 8" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
          </svg>
        </button>
      </div>

      <div class="modal-body">
        <label>Name</label>
        <input v-model="form.name" type="text" placeholder="DeepSeek" />

        <label>Base URL</label>
        <input v-model="form.base_url" type="text" placeholder="https://api.deepseek.com/v1" />

        <label>API Key</label>
        <input v-model="form.api_key" type="password" placeholder="sk-..." />

        <label>Models（一行一个）</label>
        <textarea v-model="form.models" rows="4" placeholder="deepseek-chat&#10;deepseek-reasoner"></textarea>

        <p v-if="error" class="error">{{ error }}</p>
      </div>

      <div class="modal-foot">
        <button class="btn" @click="emit('close')">Cancel</button>
        <button class="btn primary" @click="submit">Save</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(17, 24, 39, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
}

.modal {
  width: 440px;
  max-width: calc(100vw - 32px);
  background: var(--bg);
  border-radius: var(--radius);
  border: 1px solid var(--border);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.14);
}

.modal-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border);
}

.modal-head h2 {
  margin: 0;
  font-size: 15px;
  font-weight: 650;
}

.modal-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.modal-body label {
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-top: 6px;
}

.modal-body input,
.modal-body textarea {
  width: 100%;
  padding: 8px 10px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 13.5px;
  resize: vertical;
}

.modal-body input:focus,
.modal-body textarea:focus {
  outline: none;
  border-color: var(--primary);
}

.error {
  margin: 4px 0 0;
  font-size: 12.5px;
  color: #ef4444;
}

.modal-foot {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--border);
}

.btn {
  padding: 7px 16px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  font-size: 13px;
  color: var(--text);
}

.btn:hover {
  background: var(--bg-hover);
}

.btn.primary {
  background: var(--primary);
  color: #fff;
  border-color: var(--primary);
}

.btn.primary:hover {
  background: var(--primary-hover);
}
</style>