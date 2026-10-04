import { ref } from 'vue'
import { en } from './en'
import { zhCN } from './zh-CN'

export type Language = 'en' | 'zh-CN'

export type MessageKey = keyof typeof en

const STORAGE_KEY = 'easytalk.language'

const dictionaries: Record<Language, Record<MessageKey, string>> = {
  en,
  'zh-CN': zhCN,
}

// Default to English; only switch to Chinese when the user has explicitly
// chosen it before. We never auto-detect from the browser locale.
function load(): Language {
  try {
    if (localStorage.getItem(STORAGE_KEY) === 'zh-CN') return 'zh-CN'
  } catch {
    /* ignore */
  }
  return 'en'
}

const language = ref<Language>(load())

function applyHtmlLang(): void {
  document.documentElement.lang = language.value
}

export function setLanguage(lang: Language): void {
  language.value = lang
  try {
    localStorage.setItem(STORAGE_KEY, lang)
  } catch {
    /* ignore */
  }
  applyHtmlLang()
}

export function currentLanguage(): Language {
  return language.value
}

export function t(key: MessageKey, params?: Record<string, string | number>): string {
  let text: string = dictionaries[language.value][key]
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      text = text.replace(`{${k}}`, String(v))
    }
  }
  return text
}

export function useI18n() {
  return { language, t, setLanguage }
}

applyHtmlLang()