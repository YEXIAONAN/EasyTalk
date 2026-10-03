<script setup lang="ts">
import { useProviders } from '../composables/useProviders'
import { useSettings } from '../composables/useSettings'

const { providers } = useProviders()
const { settings } = useSettings()

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
</style>