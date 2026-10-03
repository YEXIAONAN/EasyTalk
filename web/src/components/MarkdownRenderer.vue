<script setup lang="ts">
import { computed } from 'vue'
import { renderMarkdown } from '../services/markdown'

const props = defineProps<{ content: string }>()

const html = computed(() => renderMarkdown(props.content))

async function onClick(e: MouseEvent) {
  const btn = (e.target as HTMLElement).closest('.copy-btn') as HTMLElement | null
  if (!btn) return

  const code = btn.closest('.code-block')?.querySelector('pre code')?.textContent ?? ''
  await copyText(code)

  btn.textContent = 'Copied'
  setTimeout(() => {
    btn.textContent = 'Copy'
  }, 1200)
}

async function copyText(text: string) {
  if (navigator.clipboard?.writeText) {
    try {
      await navigator.clipboard.writeText(text)
      return
    } catch {
      /* fall through to legacy path */
    }
  }

  const ta = document.createElement('textarea')
  ta.value = text
  ta.style.position = 'fixed'
  ta.style.opacity = '0'
  document.body.appendChild(ta)
  ta.select()
  document.execCommand('copy')
  document.body.removeChild(ta)
}
</script>

<template>
  <div class="markdown" v-html="html" @click="onClick"></div>
</template>