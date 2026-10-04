<script setup lang="ts">
import { ref } from 'vue'
import { useProviders } from '../composables/useProviders'
import { useSettings } from '../composables/useSettings'
import { useI18n } from '../i18n'
import type { Language } from '../i18n'
import * as api from '../services/api'

const { fetchProviders } = useProviders()
const { settings } = useSettings()
const { t, language, setLanguage } = useI18n()

const reloadState = ref<'settings.reloading' | 'settings.reloaded' | 'settings.reloadFailed' | ''>('')

async function onReload() {
  reloadState.value = 'settings.reloading'
  try {
    await api.reloadConfig()
    await fetchProviders()
    reloadState.value = 'settings.reloaded'
  } catch {
    reloadState.value = 'settings.reloadFailed'
  }
}

function onNumber(e: Event, key: 'temperature' | 'maxTokens' | 'topP') {
  const value = (e.target as HTMLInputElement).value
  settings[key] = value === '' ? null : Number(value)
}

function onLanguage(e: Event) {
  setLanguage((e.target as HTMLInputElement).value as Language)
}
</script>

<template>
  <div class="settings">
    <h1>{{ t('settings.title') }}</h1>

    <!-- Chat -->
    <section class="section">
      <div class="section-head"><h2>{{ t('settings.chat') }}</h2></div>
      <div class="fields">
        <label class="row">
          <input v-model="settings.streamResponse" type="checkbox" />
          <span>{{ t('settings.streamResponse') }}</span>
        </label>
      </div>
    </section>

    <!-- Appearance -->
    <section class="section">
      <div class="section-head"><h2>{{ t('settings.appearance') }}</h2></div>
      <div class="fields">
        <div class="field">
          <label>{{ t('settings.theme') }}</label>
          <div class="theme-options">
            <label class="row">
              <input v-model="settings.theme" type="radio" value="light" />
              <span>{{ t('settings.light') }}</span>
            </label>
            <label class="row">
              <input v-model="settings.theme" type="radio" value="dark" />
              <span>{{ t('settings.dark') }}</span>
            </label>
            <label class="row">
              <input v-model="settings.theme" type="radio" value="system" />
              <span>{{ t('settings.system') }}</span>
            </label>
          </div>
        </div>
      </div>
    </section>

    <!-- Language -->
    <section class="section">
      <div class="section-head"><h2>{{ t('settings.language') }}</h2></div>
      <div class="fields">
        <div class="field">
          <div class="theme-options">
            <label class="row">
              <input
                type="radio"
                name="language"
                value="en"
                :checked="language === 'en'"
                @change="onLanguage"
              />
              <span>{{ t('settings.english') }}</span>
            </label>
            <label class="row">
              <input
                type="radio"
                name="language"
                value="zh-CN"
                :checked="language === 'zh-CN'"
                @change="onLanguage"
              />
              <span>{{ t('settings.chinese') }}</span>
            </label>
          </div>
        </div>
      </div>
    </section>

    <!-- Advanced -->
    <section class="section">
      <div class="section-head"><h2>{{ t('settings.advanced') }}</h2></div>
      <div class="fields">
        <div class="field">
          <label>Temperature</label>
          <input
            type="number"
            step="0.1"
            min="0"
            max="2"
            :value="settings.temperature ?? ''"
            :placeholder="t('settings.default')"
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
            :placeholder="t('settings.default')"
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
            :placeholder="t('settings.default')"
            @input="onNumber($event, 'topP')"
          />
        </div>
      </div>
    </section>

    <!-- Configuration -->
    <section class="section">
      <div class="section-head"><h2>{{ t('settings.configuration') }}</h2></div>
      <div class="fields">
        <div class="field">
          <p class="hint">{{ t('settings.reloadHint') }}</p>
          <div class="reload-row">
            <button class="btn primary" @click="onReload">{{ t('settings.reload') }}</button>
            <span v-if="reloadState" class="reload-state">{{ t(reloadState) }}</span>
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