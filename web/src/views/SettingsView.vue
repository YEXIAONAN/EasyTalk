<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import type { Provider } from '../types'
import { useProviders } from '../composables/useProviders'
import { useSettings } from '../composables/useSettings'
import * as api from '../services/api'
import ProviderForm from '../components/ProviderForm.vue'

const { providers, loading, fetchProviders, addProvider, editProvider, removeProvider } =
  useProviders()
const { settings } = useSettings()

const showForm = ref(false)
const editing = ref<Provider | null>(null)

type TestState = { status: 'testing' | 'ok' | 'error'; message: string }
const testResults = reactive<Record<string, TestState>>({})

onMounted(fetchProviders)

function openAdd() {
  editing.value = null
  showForm.value = true
}

function openEdit(p: Provider) {
  editing.value = p
  showForm.value = true
}

async function onSave(p: Provider) {
  if (editing.value) {
    await editProvider(editing.value.name, p)
  } else {
    await addProvider(p)
  }
  showForm.value = false
}

async function onRemove(p: Provider) {
  if (!confirm(`删除 Provider「${p.name}」？`)) return
  await removeProvider(p.name)
  delete testResults[p.name]
}

async function onTest(p: Provider) {
  testResults[p.name] = { status: 'testing', message: '测试中…' }
  try {
    await api.testProvider(p.name)
    testResults[p.name] = { status: 'ok', message: '连接成功' }
  } catch (err) {
    testResults[p.name] = { status: 'error', message: friendlyTestError(p.name, err) }
  }
}

function friendlyTestError(name: string, err: unknown): string {
  const status = (err as { status?: number } | null)?.status
  if (status === 401) return 'API Key 无效'
  return `无法连接到 ${name}，请检查 Base URL 和网络`
}

function onNumber(e: Event, key: 'temperature' | 'maxTokens' | 'topP') {
  const value = (e.target as HTMLInputElement).value
  settings[key] = value === '' ? null : Number(value)
}

const modelsOfSelectedProvider = () =>
  providers.value.find((p) => p.name === settings.defaultProvider)?.models ?? []
</script>

<template>
  <div class="settings">
    <h1>Settings</h1>

    <!-- Providers -->
    <section class="section">
      <div class="section-head">
        <h2>Providers</h2>
        <button class="btn primary" @click="openAdd">Add Provider</button>
      </div>

      <div v-if="loading" class="hint">加载中…</div>
      <div v-else-if="providers.length === 0" class="hint">
        尚未配置任何 Provider。点击「Add Provider」添加你的第一个 AI 服务。
      </div>
      <ul v-else class="list">
        <li v-for="p in providers" :key="p.name" class="item">
          <div class="info">
            <div class="name">{{ p.name }}</div>
            <div class="meta">{{ p.base_url }}</div>
            <div class="meta models">{{ p.models.join('，') || '无模型' }}</div>
          </div>
          <div class="actions">
            <span v-if="testResults[p.name]" class="test-result" :class="testResults[p.name].status">
              {{ testResults[p.name].message }}
            </span>
            <button class="btn" @click="onTest(p)">Test</button>
            <button class="btn" @click="openEdit(p)">Edit</button>
            <button class="btn danger" @click="onRemove(p)">Delete</button>
          </div>
        </li>
      </ul>
    </section>

    <!-- General -->
    <section class="section">
      <div class="section-head"><h2>General</h2></div>
      <div class="fields">
        <div class="field">
          <label>Default Provider</label>
          <select v-model="settings.defaultProvider" class="select-full">
            <option value="">（无）</option>
            <option v-for="p in providers" :key="p.name" :value="p.name">{{ p.name }}</option>
          </select>
        </div>
        <div class="field">
          <label>Default Model</label>
          <select v-model="settings.defaultModel" class="select-full">
            <option value="">（无）</option>
            <option v-for="m in modelsOfSelectedProvider()" :key="m" :value="m">{{ m }}</option>
          </select>
        </div>
      </div>
    </section>

    <!-- Chat -->
    <section class="section">
      <div class="section-head"><h2>Chat</h2></div>
      <div class="fields">
        <label class="row">
          <input v-model="settings.streamResponse" type="checkbox" />
          <span>流式输出（Stream Response）</span>
        </label>
        <label class="row">
          <input v-model="settings.saveHistory" type="checkbox" />
          <span>保存聊天历史（Save History）</span>
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

    <ProviderForm v-if="showForm" :provider="editing" @close="showForm = false" @save="onSave" />
  </div>
</template>

<style scoped>
.settings {
  padding: 24px 32px 40px;
  max-width: 760px;
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

.hint {
  padding: 20px 16px;
  font-size: 13.5px;
  color: var(--text-secondary);
}

.list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border);
}

.item:last-child {
  border-bottom: none;
}

.info {
  min-width: 0;
}

.name {
  font-size: 14px;
  font-weight: 600;
}

.meta {
  font-size: 12.5px;
  color: var(--text-secondary);
  margin-top: 2px;
  word-break: break-all;
}

.models {
  color: var(--text-muted);
}

.actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.test-result {
  font-size: 12px;
  margin-right: 4px;
  max-width: 220px;
  word-break: break-all;
}

.test-result.ok {
  color: #16a34a;
}

.test-result.error {
  color: #ef4444;
}

.test-result.testing {
  color: var(--text-muted);
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

.field input:not([type='radio']):not([type='checkbox']),
.select-full {
  width: 100%;
  max-width: 320px;
  padding: 7px 10px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 13.5px;
  background: var(--bg-input);
}

.field input:focus,
.select-full:focus {
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

.btn {
  padding: 6px 12px;
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

.btn.danger {
  color: #ef4444;
}

.btn.danger:hover {
  background: #fef2f2;
  border-color: #fecaca;
}
</style>