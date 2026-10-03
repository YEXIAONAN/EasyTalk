<script setup lang="ts">
import { ref } from 'vue'
import { useProviders } from '../composables/useProviders'
import { useSettings } from '../composables/useSettings'
import * as api from '../services/api'

const { fetchProviders } = useProviders()
const { settings } = useSettings()

const reloadState = ref('')

async function onReload() {
  reloadState.value = '加载中…'
  try {
    await api.reloadConfig()
    await fetchProviders()
    reloadState.value = '已重新加载'
  } catch {
    reloadState.value = '重载失败'
  }
}

function onNumber(e: Event, key: 'temperature' | 'maxTokens' | 'topP') {
  const value = (e.target as HTMLInputElement).value
  settings[key] = value === '' ? null : Number(value)
}
</script>

<template>
  <div class="settings">
    <h1>Settings</h1>

    <!-- Chat -->
    <section class="section">
      <div class="section-head"><h2>Chat</h2></div>
      <div class="fields">
        <label class="row">
          <input v-model="settings.streamResponse" type="checkbox" />
          <span>流式输出（Stream Response）</span>
        </label>
      </div>
    </section>

    <!-- Appearance -->
    <section class="section">
      <div class="section-head"><h2>Appearance</h2></div>
      <div class="fields">
        <div class="field">
          <label>Theme</label>
          <div class="theme-options">
            <label class="row">
              <input v-model="settings.theme" type="radio" value="light" />
              <span>Light</span>
            </label>
            <label class="row">
              <input v-model="settings.theme" type="radio" value="dark" />
              <span>Dark</span>
            </label>
            <label class="row">
              <input v-model="settings.theme" type="radio" value="system" />
              <span>System</span>
            </label>
          </div>
        </div>
      </div>
    </section>

    <!-- Advanced -->
    <section class="section">
      <div class="section-head"><h2>Advanced</h2></div>
      <div class="fields">
        <div class="field">
          <label>Temperature</label>
          <input
            type="number"
            step="0.1"
            min="0"
            max="2"
            :value="settings.temperature ?? ''"
            placeholder="默认"
            @input="onNumber($event, 'temperature')"
          />
        </div>
        <div class="field">
          <label>Max Tokens</label>
          <input
            type="number"
            step="1"
            min="1"
            :value="settings.maxTokens ?? ''"
            placeholder="默认"
            @input="onNumber($event, 'maxTokens')"
          />
        </div>
        <div class="field">
          <label>Top P</label>
          <input
            type="number"
            step="0.1"
            min="0"
            max="1"
            :value="settings.topP ?? ''"
            placeholder="默认"
            @input="onNumber($event, 'topP')"
          />
        </div>
      </div>
    </section>

    <!-- Configuration -->
    <section class="section">
      <div class="section-head"><h2>Configuration</h2></div>
      <div class="fields">
        <div class="field">
          <p class="hint">修改 config.json 后点击重新加载，刷新 Provider 和 Model 列表。</p>
          <div class="reload-row">
            <button class="btn primary" @click="onReload">Reload Config</button>
            <span v-if="reloadState" class="reload-state">{{ reloadState }}</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.settings {
  padding: 24px 32px 40px;
  max-width: 680px;
  height: 100%;
  overflow-y: auto;
}

h1 {
  font-size: 20px;
  margin: 0 0 18px;
}

.section {
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
  margin-bottom: 18px;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 13px 16px;
  border-bottom: 1px solid var(--border);
  background: var(--bg-sidebar);
}

.section-head h2 {
  margin: 0;
  font-size: 14px;
  font-weight: 650;
}

.fields {
  padding: 8px 16px 14px;
  display: flex;
  flex-direction: column;
}

.field {
  padding: 8px 0;
}

.field label {
  display: block;
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.field input:not([type='radio']):not([type='checkbox']) {
  width: 100%;
  max-width: 320px;
  padding: 7px 10px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 13.5px;
  background: var(--bg-input);
}

.field input:focus {
  outline: none;
  border-color: var(--primary);
}

.row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13.5px;
  padding: 7px 0;
  cursor: pointer;
}

.row input {
  margin: 0;
}

.theme-options {
  display: flex;
  gap: 20px;
}

.theme-options .row {
  padding: 0;
}

.hint {
  margin: 0 0 10px;
  font-size: 12.5px;
  color: var(--text-muted);
}

.reload-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.reload-state {
  font-size: 12.5px;
  color: var(--text-secondary);
}

.btn {
  padding: 6px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  font-size: 12.5px;
  color: var(--text);
  white-space: nowrap;
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