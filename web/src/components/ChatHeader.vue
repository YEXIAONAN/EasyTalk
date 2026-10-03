<script setup lang="ts">
import type { Provider } from '../types'

defineProps<{
  providers: Provider[]
  provider: string
  model: string
  hasMessages: boolean
}>()

const emit = defineEmits<{
  'update:provider': [name: string]
  'update:model': [name: string]
  settings: []
  clear: []
  toggleMenu: []
}>()
</script>

<template>
  <header class="header">
    <button class="icon-btn menu-btn" title="菜单" @click="emit('toggleMenu')">
      <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
        <path d="M2.5 4h11M2.5 8h11M2.5 12h11" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
      </svg>
    </button>

    <select
      class="select"
      :value="provider"
      :disabled="providers.length === 0"
      @change="emit('update:provider', ($event.target as HTMLSelectElement).value)"
    >
      <option v-if="providers.length === 0" value="">未配置 Provider</option>
      <option v-for="p in providers" :key="p.name" :value="p.name">{{ p.name }}</option>
    </select>

    <select
      class="select"
      :value="model"
      :disabled="providers.length === 0"
      @change="emit('update:model', ($event.target as HTMLSelectElement).value)"
    >
      <option v-if="providers.length === 0" value="">未配置 Model</option>
      <template v-else>
        <option
          v-for="m in providers.find((p) => p.name === provider)?.models ?? []"
          :key="m"
          :value="m"
        >
          {{ m }}
        </option>
      </template>
    </select>

    <div class="spacer"></div>

    <button v-if="hasMessages" class="icon-btn" title="清空当前对话" @click="emit('clear')">
      <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
        <path
          d="M3 5h10M6.5 5V3.5h3V5M5 5l.5 8h5l.5-8"
          stroke="currentColor"
          stroke-width="1.4"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </button>

    <button class="icon-btn" title="Settings" @click="emit('settings')">
      <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
        <circle cx="8" cy="8" r="2.4" stroke="currentColor" stroke-width="1.4" />
        <path
          d="M8 1.8v1.4M8 12.8v1.4M1.8 8h1.4M12.8 8h1.4M3.6 3.6l1 1M11.4 11.4l1 1M12.4 3.6l-1 1M4.6 11.4l-1 1"
          stroke="currentColor"
          stroke-width="1.4"
          stroke-linecap="round"
        />
      </svg>
    </button>
  </header>
</template>

<style scoped>
.header {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 52px;
  padding: 0 16px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.select {
  max-width: 220px;
  height: 32px;
  padding: 0 28px 0 10px;
  font-size: 13px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--bg-input);
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 16 16'%3E%3Cpath d='M4 6l4 4 4-4' stroke='%236b7280' stroke-width='1.6' fill='none' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 8px center;
}

.select:focus {
  outline: none;
  border-color: var(--primary);
}

.spacer {
  flex: 1;
}

.menu-btn {
  display: none;
}

@media (max-width: 768px) {
  .menu-btn {
    display: inline-flex;
    flex-shrink: 0;
  }

  .select {
    max-width: none;
    flex: 1;
    min-width: 0;
  }
}
</style>