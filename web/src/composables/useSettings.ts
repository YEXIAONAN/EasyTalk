import { reactive, watch } from 'vue'
import type { Settings } from '../types'

const STORAGE_KEY = 'easytalk.settings'

const defaults: Settings = {
  streamResponse: true,
  theme: 'system',
  temperature: null,
  maxTokens: null,
  topP: null,
}

function load(): Settings {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return { ...defaults, ...JSON.parse(raw) }
  } catch {
    /* ignore */
  }
  return { ...defaults }
}

const settings = reactive<Settings>(load())

watch(
  settings,
  () => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings))
  },
  { deep: true },
)

function applyTheme() {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const dark = settings.theme === 'dark' || (settings.theme === 'system' && prefersDark)
  document.documentElement.dataset.theme = dark ? 'dark' : 'light'
}

export function useSettings() {
  return { settings, applyTheme }
}