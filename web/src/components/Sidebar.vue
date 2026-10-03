<script setup lang="ts">
import { computed } from 'vue'
import { BRANDING } from '../config/branding'

const props = defineProps<{
  active: string
  mobileOpen: boolean
  provider: string
  model: string
  requestCount: number
  inputTokens: number
  outputTokens: number
  cachedTokens: number | null
  totalTokens: number
  lastResponseMs: number | null
  firstTokenMs: number | null
}>()

const emit = defineEmits<{
  navigate: [view: 'chat' | 'settings' | 'about']
  clear: []
}>()

const connected = computed(() => props.provider !== '')

const numberFmt = new Intl.NumberFormat('en-US')

function fmt(n: number): string {
  return numberFmt.format(n)
}

function fmtDuration(ms: number | null): string {
  if (ms === null) return '—'
  if (ms < 1000) return `${Math.round(ms)}ms`
  return `${(ms / 1000).toFixed(1)}s`
}

function fmtTokens(n: number | null): string {
  return n === null ? '—' : numberFmt.format(n)
}
</script>

<template>
  <aside class="sidebar" :class="{ open: mobileOpen }">
    <div class="brand">
      <img class="logo" :src="BRANDING.logoMark" alt="EasyTalk" />
      <span class="name">{{ BRANDING.name }}</span>
    </div>

    <div class="scroll">
      <!-- Current session -->
      <div class="panel">
        <div class="panel-title">Current Session</div>
        <div class="stat">
          <span class="stat-label">Provider</span>
          <span class="stat-value">{{ provider || '—' }}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Model</span>
          <span class="stat-value">{{ model || '—' }}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Status</span>
          <span class="stat-value status" :class="{ connected }">
            <i class="dot"></i>
            {{ connected ? 'Connected' : 'Idle' }}
          </span>
        </div>
      </div>

      <!-- Token usage -->
      <div class="panel">
        <div class="panel-title">Token Usage</div>
        <div class="stat">
          <span class="stat-label">Input</span>
          <span class="stat-value">{{ fmt(inputTokens) }}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Output</span>
          <span class="stat-value">{{ fmt(outputTokens) }}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Cached</span>
          <span class="stat-value">{{ fmtTokens(cachedTokens) }}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Total</span>
          <span class="stat-value">{{ fmt(totalTokens) }}</span>
        </div>
      </div>

      <!-- Request -->
      <div class="panel">
        <div class="panel-title">Request</div>
        <div class="stat">
          <span class="stat-label">Requests</span>
          <span class="stat-value">{{ requestCount }}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Last Response</span>
          <span class="stat-value">{{ fmtDuration(lastResponseMs) }}</span>
        </div>
        <div class="stat">
          <span class="stat-label">First Token</span>
          <span class="stat-value">{{ fmtDuration(firstTokenMs) }}</span>
        </div>
      </div>
    </div>

    <nav class="bottom">
      <button class="nav-item" @click="emit('clear')">
        <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
          <path
            d="M13.5 8A5.5 5.5 0 1 1 8 2.5M13.5 8V3.5M13.5 8h-4.5"
            stroke="currentColor"
            stroke-width="1.4"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
        <span>Clear Session</span>
      </button>
      <button
        class="nav-item"
        :class="{ active: active === 'settings' }"
        @click="emit('navigate', 'settings')"
      >
        <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
          <circle cx="8" cy="8" r="2.4" stroke="currentColor" stroke-width="1.4" />
          <path
            d="M8 1.8v1.4M8 12.8v1.4M1.8 8h1.4M12.8 8h1.4M3.6 3.6l1 1M11.4 11.4l1 1M12.4 3.6l-1 1M4.6 11.4l-1 1"
            stroke="currentColor"
            stroke-width="1.4"
            stroke-linecap="round"
          />
        </svg>
        <span>Settings</span>
      </button>
      <button
        class="nav-item"
        :class="{ active: active === 'about' }"
        @click="emit('navigate', 'about')"
      >
        <svg viewBox="0 0 16 16" width="16" height="16" fill="none">
          <circle cx="8" cy="8" r="6.2" stroke="currentColor" stroke-width="1.4" />
          <path d="M8 7.2v4M8 5v.2" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
        </svg>
        <span>About</span>
      </button>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  width: var(--sidebar-width);
  height: 100%;
  background: var(--bg-sidebar);
  border-right: 1px solid var(--border);
  padding: 14px 12px;
  flex-shrink: 0;
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 2px 6px 14px;
}

.logo {
  display: block;
  width: 28px;
  height: 28px;
  border-radius: 7px;
}

.name {
  font-size: 15px;
  font-weight: 650;
}

.scroll {
  flex: 1;
  overflow-y: auto;
}

.panel {
  padding: 6px 2px 16px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 14px;
}

.panel-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.4px;
  padding: 0 6px 8px;
}

.stat {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 6px;
  font-size: 13px;
}

.stat-label {
  color: var(--text-secondary);
}

.stat-value {
  color: var(--text);
  font-weight: 550;
  text-align: right;
  max-width: 60%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--text-muted);
}

.status.connected .dot {
  background: #22c55e;
}

.bottom {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 8px 10px;
  border-radius: var(--radius-sm);
  font-size: 13.5px;
  color: var(--text-secondary);
  transition: background 0.12s ease, color 0.12s ease;
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text);
}

.nav-item.active {
  background: var(--bg-hover);
  color: var(--text);
  font-weight: 550;
}

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 40;
    transform: translateX(-100%);
    transition: transform 0.2s ease;
    box-shadow: none;
  }

  .sidebar.open {
    transform: translateX(0);
    box-shadow: 0 0 32px rgba(0, 0, 0, 0.24);
  }
}
</style>