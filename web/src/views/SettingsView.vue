<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import type { Provider } from '../types'
import { useProviders } from '../composables/useProviders'
import * as api from '../services/api'
import ProviderForm from '../components/ProviderForm.vue'

const { providers, loading, fetchProviders, addProvider, editProvider, removeProvider } =
  useProviders()

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
    testResults[p.name] = {
      status: 'error',
      message: friendlyTestError(p.name, err),
    }
  }
}

function friendlyTestError(name: string, err: unknown): string {
  const status = (err as { status?: number } | null)?.status
  if (status === 401) return 'API Key 无效'
  return `无法连接到 ${name}，请检查 Base URL 和网络`
}
</script>

<template>
  <div class="settings">
    <h1>Settings</h1>

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
            <span
              v-if="testResults[p.name]"
              class="test-result"
              :class="testResults[p.name].status"
            >
              {{ testResults[p.name].message }}
            </span>
            <button class="btn" @click="onTest(p)">Test</button>
            <button class="btn" @click="openEdit(p)">Edit</button>
            <button class="btn danger" @click="onRemove(p)">Delete</button>
          </div>
        </li>
      </ul>
    </section>

    <ProviderForm v-if="showForm" :provider="editing" @close="showForm = false" @save="onSave" />
  </div>
</template>

<style scoped>
.settings {
  padding: 28px 32px;
  max-width: 820px;
  height: 100%;
  overflow-y: auto;
}

h1 {
  font-size: 20px;
  margin: 0 0 20px;
}

.section {
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
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